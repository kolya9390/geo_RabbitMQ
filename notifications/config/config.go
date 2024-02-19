package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
)

const (
	AppName = "APP_NAME"
)

type AppConf struct {
	AppName     string   `yaml:"app_name"`
	Environment string   `yaml:"environment"`
	Domain      string   `yaml:"domain"`
	APIUrl      string   `yaml:"api_url"`
	Token       Token    `yaml:"token"`
	DB          DB       `yaml:"db"`
	NoticRPC    NoticRPC  `yaml:"notic_rpc"`
	RPCServer   RPCServer `yaml:"rpc_server"`
	UserRPC     UserRPC   `yaml:"user_rpc"`
	AuthorizationDADATA	AuthorizationDADATA
	Rebbit_host string
}

type NoticRPC struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
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

type Cache struct {
	Address  string `yaml:"address"`
	Password string `json:"-" yaml:"password"`
	Port     string `yaml:"port"`
}


type AuthorizationDADATA struct {
	ApiKeyValue    string
	SecretKeyValue string
}

type RPCServer struct {
	Port         string        `yaml:"port"`
	ShutdownTime time.Duration `yaml:"shutdown_timeout"`
	Type         string        `yaml:"type"` // grpc or jsonrpc or http
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


func NewAppConf(envPath string /*"server_app/.env"*/) AppConf {

	env, err := godotenv.Read(envPath)

	if err != nil {
		log.Println(err)
	}
//TODO COFIG
	return AppConf{
		AppName: env[AppName],
		Rebbit_host: env["REBBTI_HOST"],
		RPCServer: RPCServer{
			Port: env["RPC_PORT"],
			ShutdownTime: 1,
			Type: env["RPC_PROTOCOL"],
		},
		NoticRPC: NoticRPC{
			Host: env["NOTIFIC_HOST"],
			Port: env["NOTIFIC_PORT"],
		},

		UserRPC: UserRPC{
			Host: "",
			Port: "",
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
		AuthorizationDADATA: AuthorizationDADATA{
			ApiKeyValue: env["DADATA_API_KEY"],
			SecretKeyValue: env["DADATA_SECRET_KEY"],
		},
	}
}
