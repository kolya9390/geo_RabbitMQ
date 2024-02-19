package geo_client

type GeoClient interface {
	SearchGeoAdres(query RequestAddressSearch) ([]Address,error)
	GeoCoder(geocode RequestAddressGeocode) ([]Address,error)
}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}