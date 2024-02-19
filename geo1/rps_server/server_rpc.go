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
	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/service/rebbit"
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

	config := config.NewAppConf("server_app/.env")

	// Инициализация подключения к базе данных
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)

	db, err := sqlx.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}
	time.Sleep(time.Second * 3)
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

	conn, err := amqp.Dial(urlRebbit)
	if err != nil {
		log.Fatalf("Error conn: %s", err)
	}
	log.Println("successful conn RebbitMQ")

	defer conn.Close()

	serviceRebbit, err := rebbit.NewRabbitMQ(conn)

	if err != nil {
		log.Fatalf("Error NewRebbit: %s", err)
	}
//
	cache := storage.NewGeoRedis(redisClient)
	storageDB := storage.NewGeoRepositoryProxy(*postgresDB, cache)
	sevisDAdata := dadata.NewDadataService(config.AuthorizationDADATA)

	err = postgresDB.ConnectToDB()

	if err != nil {
		log.Printf("Error conect DB %s", err)
	}

	gs.GeoProviderGRPCServer.geoProvider = app.NewGeoProvider(storageDB,sevisDAdata,serviceRebbit)
	
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
