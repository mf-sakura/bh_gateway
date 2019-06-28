package config

import (
	"github.com/kelseyhightower/envconfig"
)

type GRPCConfig struct {
	HotelPort int    `envconfig:"HOTEL_PORT" default:"5001"`
	HotelHost string `envconfig:"HOTEL_HOST" default:"bh_hotel"`
	UserPort  string `envconfig:"USER_PORT" default:"5002"`
	UserHost  string `envconfig:"USER_PORT" default:"bh_user"`
}

func LoadConifg() (*GRPCConfig, error) {
	var gRPCConf GRPCConfig
	if err := envconfig.Process("", &gRPCConf); err != nil {
		return nil, err
	}

	return &gRPCConf, nil
}
