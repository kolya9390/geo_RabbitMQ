package rpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/app"
	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/config"
	geo_provider "github.com/kolya9390/gRPC_GeoProvider/server_rpc/gen"
	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/service/dadata"
	"gitlab.com/ptflp/gopubsub/kafkamq"
	"gitlab.com/ptflp/gopubsub/queue"
	"gitlab.com/ptflp/gopubsub/rabbitmq"

	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/storage"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)


type GeoService struct {
	GeoProviderGRPCServer
}

// GeoProviderGRPCServer must be embedded to have forward compatible implementations.
type GeoProviderGRPCServer struct {
	geoProvider app.GeoProvider
	geo_provider.UnimplementedGeoProviderGRPCServer
}


func NewGeoServis() *GeoService {
	return &GeoService{}
}

func (gs *GeoService) StartServer() error {

	config := config.NewAppConf()

	// Инициализация подключения к базе данных
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)

	db, err := sqlx.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}
	time.Sleep(time.Second * 6)
	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %s", err)
	}

	postgresDB := storage.NewGeoRepositoryDB(db)

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Cache.Address, config.Cache.Port),
	})

	defer db.Close()
	defer redisClient.Close()

	// rebbit
	urlRebbit := fmt.Sprintf("amqp://guest:guest@%s:5672/", config.Rebbit_host)
	urlKafkaCon := fmt.Sprintf("%s:%s", config.Kafka.Host, config.Kafka.Port)

	var broker queue.MessageQueuer
	switch config.Broker_type {
	case "kafka":
		service_kafka, err := kafkamq.NewKafkaMQ(urlKafkaCon, "myGroup")

		if err != nil {
			log.Printf("Error NewKafkaMQ: %s", err)
		}
		broker = service_kafka
		log.Println("successful conn Kafka")
	case "rebbit":
		conn, err := amqp.Dial(urlRebbit)
		if err != nil {
			log.Fatalf("Error conn: %s", err)
		}
		defer conn.Close()
		sevisRebbit, err := rabbitmq.NewRabbitMQ(conn)

		if err != nil {
			log.Printf("Error NewRebbit: %s", err)
		}
		broker = sevisRebbit
		log.Println("successful conn RebbitMQ")
	default:
		log.Fatalf("Unknown broker type: %s", config.Broker_type)
		return fmt.Errorf("unknown broker type: %s", config.Broker_type)
	}
	if err != nil {
		log.Printf("Error subscribing to messages: %v", err)
		return err
	}
//
	cache := storage.NewGeoRedis(redisClient)
	storageDB := storage.NewGeoRepositoryProxy(*postgresDB, cache)
	sevisDAdata := dadata.NewDadataService(config.AuthorizationDADATA)

	err = postgresDB.ConnectToDB()

	if err != nil {
		log.Printf("Error conect DB %s", err)
	}

	gs.GeoProviderGRPCServer.geoProvider = app.NewGeoProvider(storageDB,sevisDAdata,broker)
	
	//

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GeoRPC.Port))
	if err != nil {
		log.Printf("Eroor Listen %v", err)
		return err
	}
	defer listen.Close()



	log.Printf("RPC типа %s сервер запущен и прослушивает порт :%s", config.RPCServer.Type, config.GeoRPC.Port)

	//
		grpcServer := grpc.NewServer()
		geo_provider.RegisterGeoProviderGRPCServer(grpcServer, 
			&gs.GeoProviderGRPCServer)
		grpcServer.Serve(listen)
		

	return nil
}




func (gs *GeoProviderGRPCServer) AddressSearchGRPC(ctx context.Context, req *geo_provider.RequestAddressSearch) (*geo_provider.RespAddress, error) {

	addresses, err := gs.geoProvider.AddressSearch(req.Query)
	if err != nil {
		log.Printf("Error AddressSearch: %v", err)
		return nil, err
	}

	return &geo_provider.RespAddress{
		Geolat: addresses[0].GeoLat,
		Geolon: addresses[0].GeoLon,
		Result: addresses[0].Result,
	}, nil
}
func (gs *GeoProviderGRPCServer) AddressGeoCodeGRPC(ctx context.Context, req *geo_provider.RequestAddressGeocode) (*geo_provider.RespAddress, error) {

	addresses, err := gs.geoProvider.GeoCode(req.Lat, req.Lng)
	if err != nil {
		log.Printf("Error AddressGeoCode: %v", err)
		return nil, err
	}

	return &geo_provider.RespAddress{
		Geolat: addresses[0].GeoLat,
		Geolon: addresses[0].GeoLon,
		Result: addresses[0].Result,
	}, nil
}
