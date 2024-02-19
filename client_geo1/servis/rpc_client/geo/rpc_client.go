package geo_client

import (
	"context"
	"fmt"
	"log"

	"github.com/kolya9390/gRPC_GeoProvider/client_Proxy/config"
	geo_provider "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/protoc/gen/geo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type grpcGeoClient struct {
	client *grpc.ClientConn
}

func NewGeoClient() GeoClient {

	cfgRPC := config.NewAppConf("client_app/.env").GeoRPC
	log.Println(cfgRPC.Host,cfgRPC.Port)
	host_port := fmt.Sprintf("%s:%s", cfgRPC.Host, cfgRPC.Port)
		client, err := grpc.Dial(host_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect %s", err)
		}
		return &grpcGeoClient{client: client}

}

type Address struct {
	GeoLat string `json:"lat"`
	GeoLon string `json:"lon"`
	Result string `json:"result"`
}


func (c *grpcGeoClient) SearchGeoAdres(query RequestAddressSearch) ([]Address,error) {
	client := geo_provider.NewGeoProviderGRPCClient(c.client)
	// Вызываем метод SearchGeoAdres
	resp, err := client.AddressSearchGRPC(context.Background(), &geo_provider.RequestAddressSearch{
		Query: query.Query,
	})
	if err != nil {
		log.Printf("Error searching for address: %v", err)
		return nil,err
	}
	// Формируем ответ
	return []Address{
		{
			GeoLat: resp.Geolat,
			GeoLon: resp.Geolon,
			Result: resp.Result,
		},
	},nil
}

func (c *grpcGeoClient) GeoCoder(geocode RequestAddressGeocode) ([]Address,error) {
	// Определяем gRPC клиент
	client := geo_provider.NewGeoProviderGRPCClient(c.client)
	// Вызываем метод GeoCoder
	resp, err := client.AddressGeoCodeGRPC(context.Background(), &geo_provider.RequestAddressGeocode{
		Lat: geocode.Lat,
		Lng: geocode.Lng,
	})
	if err != nil {
		log.Printf("Error geocoding address: %v", err)
		return nil,err
	}
	// Формируем ответ
	return []Address{
		{
			GeoLat: resp.Geolat,
			GeoLon: resp.Geolon,
			Result: resp.Result,
		},
	},nil
}
