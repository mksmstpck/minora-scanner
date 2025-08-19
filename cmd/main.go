package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/mksmstpck/minora-scanner/internal/config"
	"github.com/mksmstpck/minora-scanner/internal/events"
	"github.com/mksmstpck/minora-scanner/internal/services"
	"github.com/mksmstpck/minora-scanner/internal/storage"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.Function), f.Line)
		},
	}
	logrus.SetFormatter(formatter)
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}

func main() {
	httpClient := http.Client{}

	config := config.NewConfig()

	c := cache.New(config.CacheExpMin, config.CacheExpMin)

	storage := storage.NewStorage(c, config)

	events := events.NewEvents(config, &httpClient)

	services := services.NewServiecs(events, storage)

	pairs, err := services.SeekPairs()
	if err != nil {
		logrus.Error(err)
	}

	for _, pair := range pairs {
		logrus.Infof("coin: %s \n priceHighCex: %d \n priceLowCex: %d \n spread: %f", pair.Coin, pair.PriceHigh.CexType, pair.PriceLow.CexType, pair.SpreadPercents)
	}
}
