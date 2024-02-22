package rpcserver

import (
	"context"
	"fmt"
	"gitlab.com/ptflp/gopubsub/kafkamq"
	"gitlab.com/ptflp/gopubsub/queue"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"log"
	"net"
	"notific/config"
	notific_service "notific/protoc/gen"
	"time"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type NotificService struct {
	NotificProviderGRPCServer
}

// GeoProviderGRPCServer must be embedded to have forward compatible implementations.
type NotificProviderGRPCServer struct {
	ServiceBroker queue.MessageQueuer
	notific_service.UnimplementedNotificServiceGRPCServer
	msseges []string
}

func NewGeoServis() *NotificService {
	return &NotificService{}
}

func (ns *NotificService) StartServer() error {

	config := config.NewAppConf()
	urlConRebbit := fmt.Sprintf("amqp://guest:guest@%s:5672/", config.Rebbit_host)
	urlConKafka := fmt.Sprintf("%s:%s", config.Kafka.Host, config.Kafka.Port)

//	log.Println(config)
	time.Sleep(15 * time.Second)

	var err error
	var messages <-chan queue.Message
	switch config.BrokerType {
	case "kafka":
		service_kafka, err := kafkamq.NewKafkaMQ(urlConKafka, "myGroup")

		if err != nil {
			log.Printf("Error NewKafkaMQ: %s", err)
		}
		ns.ServiceBroker = service_kafka
		log.Println("successful conn Kafka")
		messages, err = ns.ServiceBroker.Subscribe("notification")
	case "rebbit":
		conn, err := amqp.Dial(urlConRebbit)
		if err != nil {
			log.Printf("Error conn: %s", err)
		}

		defer conn.Close()

		sevisRebbit, err := rabbitmq.NewRabbitMQ(conn)

		if err != nil {
			log.Printf("Error NewRebbit: %s", err)
		}
		ns.ServiceBroker = sevisRebbit
		log.Println("successful conn RebbitMQ")
		if err := rabbitmq.CreateExchange(conn, "notification", "direct"); err != nil {
			log.Printf("Error CreateExchange: %s", err)
		}
		messages, err = ns.ServiceBroker.Subscribe("notification")
	default:
		log.Fatalf("Unknown broker type: %s", config.BrokerType)
		return fmt.Errorf("unknown broker type: %s", config.BrokerType)
	}
	if err != nil {
		log.Fatalf("Error subscribing to messages: %v", err)
		return err
	}
	//


	go func(mssgs <- chan queue.Message) {
		for {
			select {
			case msg, ok := <-mssgs:
				if !ok {
					return
				}
				if msg.Err != nil {
					log.Printf("Failed to receive message: %s\n", msg.Err)
				}
				fmt.Printf("Received message: %s\n", string(msg.Data))
				ns.msseges=append(ns.msseges, string(msg.Data))
				// Подтверждаем получение сообщения
				err = ns.ServiceBroker.Ack(&msg)
				if err != nil {
					log.Printf("Failed to ack message: %s\n", err)
				}
			}
	

		}

	}(messages)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", config.NoticRPC.Port))
	if err != nil {
		log.Printf("Error Listen %v", err)
		return err
	}
	defer listen.Close()

	log.Printf("RPC типа %s сервер запущен и прослушивает порт :%s", config.RPCServer.Type, config.NoticRPC.Port)

	//
	grpcServer := grpc.NewServer()
	notific_service.RegisterNotificServiceGRPCServer(grpcServer,
		&ns.NotificProviderGRPCServer)
	grpcServer.Serve(listen)


	return nil
}

func (ns *NotificProviderGRPCServer) SmsSend(ctx context.Context, req *notific_service.RequestSMS) (*notific_service.RespListSMS, error) {

	respNotific := &notific_service.RespListSMS{}
	if len(ns.msseges) == 0 {
		respNotific.SMSlist = append(respNotific.SMSlist, &notific_service.RespSMS{Phones: "1234566789", Msg: "0"})
		return respNotific, nil
	}
	for _, msg := range ns.msseges {
		respNotific.SMSlist = append(respNotific.SMSlist, &notific_service.RespSMS{Phones: "1234566789", Msg: msg})
	}

	return respNotific, nil
}

func (ns *NotificProviderGRPCServer) EmailSend(ctx context.Context, req *notific_service.RequestEmail) (*notific_service.RespListEmail, error) {

	respNotific := &notific_service.RespListEmail{}

	if len(ns.msseges) == 0 {
		respNotific.ListEmail = append(respNotific.ListEmail, &notific_service.RespEmail{Email: "test@test.com", Msg: "0"})
		return respNotific, nil
	}

	for _, msg := range ns.msseges {
		respNotific.ListEmail = append(respNotific.ListEmail, &notific_service.RespEmail{Email: "test@test.com", Msg: msg})
	}

	return respNotific, nil
}
