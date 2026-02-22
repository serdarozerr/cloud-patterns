package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

type Config struct {
	Mode string `json:"mode"`
	Port string `json:"port"`
}

func NewConfig(path string) *Config {
	f,err:=os.Open(path)

	if err!=nil{
		slog.Info("opening config file","error",err)
		panic(1)
	}

	var cfg Config

	err=json.NewDecoder(f).Decode(&cfg)
	if err !=nil{
		slog.Info("decoding config file","error",err)
	}
	return &cfg
}
