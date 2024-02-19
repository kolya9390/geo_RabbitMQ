package rpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"notific/config"
	notific_service "notific/protoc/gen"
	"notific/service"
	"time"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type NotificService struct {
	NotificProviderGRPCServer
}

// GeoProviderGRPCServer must be embedded to have forward compatible implementations.
type NotificProviderGRPCServer struct {
	ServiceRebbit service.MessageQueuer
	notific_service.UnimplementedNotificServiceGRPCServer
	msseges []string
}

func NewGeoServis() *NotificService {
	return &NotificService{}
}

func (ns *NotificService) StartServer() error {

	config := config.NewAppConf("server_app/.env")
	urlRebbit := fmt.Sprintf("amqp://guest:guest@%s:5672/", config.Rebbit_host)

	time.Sleep(15 * time.Second)

	conn, err := amqp.Dial(urlRebbit)
	if err != nil {
		log.Fatalf("Error conn: %s", err)
	}
	log.Println("successful conn RebbitMQ")

	defer conn.Close()

	sevisRebbit, err := service.NewRabbitMQ(conn)

	if err != nil {
		log.Fatalf("Error NewRebbit: %s", err)
	}

	if err := service.CreateExchange(conn, "notification", "direct"); err != nil {
		log.Fatalf("Error Test CreateExchange: %s", err)
	}

	ns.ServiceRebbit = sevisRebbit

	//

	messages, err := ns.ServiceRebbit.Subscribe("notification")

	go func() {
		for d := range messages {
			log.Println(string(d.Data))
			ns.msseges = append(ns.msseges, string(d.Data))
		}
	}()
	if err != nil {
		log.Println(err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", config.NoticRPC.Port))
	if err != nil {
		log.Printf("Eroor Listen %v", err)
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
