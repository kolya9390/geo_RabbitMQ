package config

import (
	"log"
	//	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	AppName = "APP_NAME"
)

type AppConf struct {
	ServepPort  string
	AppName     string    `yaml:"app_name"`
	Environment string    `yaml:"environment"`
	Domain      string    `yaml:"domain"`
	APIUrl      string    `yaml:"api_url"`
	Token       Token     `yaml:"token"`
	DB          DB        `yaml:"db"`
	RPCServer   RPCServer `yaml:"rpc_server"`
	UserRPC     UserRPC   `yaml:"user_rpc"`
	GeoRPC      GeoRPC    `yaml:"geo_rpc"`
	AuthRPC     AuthRPC   `yaml:"auth_rpc"`
	NoticRPC    NoticRPC  `yaml:"notic_rpc"`
}

type DB struct {
	Net      string `yaml:"net"`
	Driver   string `yaml:"driver"`
	Name     string `yaml:"name"`
	User     string `json:"-" yaml:"user"`
	Password string `json:"-" yaml:"password"`
	Host     string `yaml:"host"`
	MaxConn  int    `yaml:"max_conn"`
	Port     string `yaml:"port"`
	Timeout  int    `yaml:"timeout"`
}

type RPCServer struct {
	Port         string `yaml:"port"`
	RPC_SERVER   string
	ShutdownTime time.Duration `yaml:"shutdown_timeout"`
	Type         string        `yaml:"type"` // grpc or jsonrpc or http
}

type GeoRPC struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AuthRPC struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type NoticRPC struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type UserRPC struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Token struct {
	AccessTTL     time.Duration `yaml:"access_ttl"`
	RefreshTTL    time.Duration `yaml:"refresh_ttl"`
	AccessSecret  string        `yaml:"access_secret"`
	RefreshSecret string        `yaml:"refresh_secret"`
}

func NewAppConf(envPath string /*"client_app/.env"*/) AppConf {

	env, err := godotenv.Read(envPath)

	if err != nil {
		log.Println(err)
	}
	//TODO COFIG
	return AppConf{
		AppName:    env[AppName],
		ServepPort: env["SERVER_PORT"],
		RPCServer: RPCServer{
			Port:         env["RPC_PORT"],
			ShutdownTime: 1,
			Type:         env["RPC_PROTOCOL"],
			RPC_SERVER:   env["RPC_SERVER"],
		},
		GeoRPC: GeoRPC{
			Host: env["GEO_HOST"],
			Port: env["PORT_GEO"],
		},
		AuthRPC: AuthRPC{
			Host: env[""],
			Port: env[""],
		},
		NoticRPC: NoticRPC{
			Host: env["NOTIFIC_HOST"],
			Port: env["NOTIFIC_PORT"],
		},
		UserRPC: UserRPC{
			Host: env[""],
			Port: env[""],
		},

		DB: DB{
			Net:      env["DB_NET"],
			Driver:   env["DB_DRIVER"],
			Name:     env["DB_NAME"],
			User:     env["DB_USER"],
			Password: env["DB_PASSWORD"],
			Host:     env["DB_HOST"],
			Port:     env["DB_PORT"],
		},
		Token: Token{
			AccessSecret: env["ACCEC_SECRET"],
			AccessTTL:    90 * time.Minute,
		},
	}
}
