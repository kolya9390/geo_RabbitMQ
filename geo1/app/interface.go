package app

import "github.com/kolya9390/gRPC_GeoProvider/server_rpc/service/dadata"

type GeoProvider interface {
	AddressSearch(input string) ([]*dadata.Address, error)
	GeoCode(lat, lng string) ([]*dadata.Address, error)
}