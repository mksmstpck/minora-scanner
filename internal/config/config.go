package config

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Token             string
	ChatID            int64
	BinanceTickerUrl  string
	BinanceCexInfoUrl string
	BybitTickerUrl    string
	BybitCexInfoUrl   string
	GateTickerUrl     string
	GateCexInfoUrl    string
	KucoinTickerUrl   string
	KucoinCexInfoUrl  string
	MexcTickerUrl     string
	MexcCexInfoUrl    string
	MinPairSpread     float64
	MaxPairSpread     float64
	CoingeckoBulkData string
	CacheExpMin       time.Duration
}

func NewConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Error(err)
	}
	return Config{
		Token:             viper.GetString("TOKEN"),
		ChatID:            viper.GetInt64("CHAT_ID"),
		BinanceTickerUrl:  viper.GetString("BINANCE_URL"),
		BinanceCexInfoUrl: viper.GetString("BINANCE1_URL"),
		BybitTickerUrl:    viper.GetString("BYBIT_URL"),
		BybitCexInfoUrl:   viper.GetString("BYBIT_URL1"),
		GateTickerUrl:     viper.GetString("GATE_URL"),
		GateCexInfoUrl:    viper.GetString("GATE1_URL"),
		KucoinTickerUrl:   viper.GetString("KUCOIN_URL"),
		KucoinCexInfoUrl:  viper.GetString("KUCOIN1_URL"),
		MexcTickerUrl:     viper.GetString("MEXC_URL"),
		MexcCexInfoUrl:    viper.GetString("MEXC1_URL"),
		MinPairSpread:     viper.GetFloat64("MIN_PAIR_SPREAD"),
		MaxPairSpread:     viper.GetFloat64("MAX_PAIR_SPREAD"),
		CoingeckoBulkData: viper.GetString("COINGECKO_BULK_DATA"),
		CacheExpMin:       time.Duration(viper.GetInt("CACHE_EXP_MIN")) * time.Minute,
	}
}
