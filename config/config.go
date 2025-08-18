package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Token             string
	ChatID            string
	BinanceUrl        string
	BybitUrl          string
	GateUrl           string
	KucoinUrl         string
	MexcUrl           string
	MinPairSpread     float64
	MaxPairSpread     float64
	CoingeckoBulkData string
}

func NewConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	return Config{
		Token:             viper.GetString("TOKEN"),
		ChatID:            viper.GetString("CHAT_ID"),
		BinanceUrl:        viper.GetString("BINANCE_URL"),
		BybitUrl:          viper.GetString("BYBIT_URL"),
		GateUrl:           viper.GetString("GATE_URL"),
		KucoinUrl:         viper.GetString("KUCOIN_URL"),
		MexcUrl:           viper.GetString("MEXC_URL"),
		MinPairSpread:     viper.GetFloat64("MIN_PAIR_SPREAD"),
		MaxPairSpread:     viper.GetFloat64("MAX_PAIR_SPREAD"),
		CoingeckoBulkData: viper.GetString("COINGECKO_BULK_DATA"),
	}
}
