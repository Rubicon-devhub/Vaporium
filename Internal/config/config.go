package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"development"`
	HttpServer `yaml:"httpServer"`
}
type HttpServer struct  {
	Address string `yaml:"address" env-default:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"60s"`

}


func MustLoad() Config{
	configPath := os.Getenv("CONFIG_PATH")
	
	
	if configPath == ""{
		log.Fatal("Config path not found ")
	}

	if _,err:= os.Stat(configPath);os.IsNotExist(err){
		log.Fatalf("config file %s not exist",configPath)
	}

	
	var cfg Config

	if err:= cleanenv.ReadConfig(configPath,&cfg);err != nil{
		log.Fatalf("cannot read config: %s",err)
	}
	return  cfg

}