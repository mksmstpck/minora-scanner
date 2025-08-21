package storage

import (
	"fmt"

	"github.com/mksmstpck/minora-scanner/internal/config"
	"github.com/mksmstpck/minora-scanner/internal/models"
	"github.com/patrickmn/go-cache"
)

type Storage struct {
	cache  *cache.Cache
	config config.Config
}

type Storager interface {
	CheckExists(pair models.Pair) bool
	SetPair(pair models.Pair)
	genKey(pair models.Pair) string
}

func NewStorage(cache *cache.Cache, config config.Config) Storager {
	return &Storage{
		cache:  cache,
		config: config,
	}
}

func (s *Storage) SetPair(pair models.Pair) {
	s.cache.Set(s.genKey(pair), pair, s.config.CacheExpMin)
}

func (s *Storage) CheckExists(pair models.Pair) bool {
	_, res := s.cache.Get(s.genKey(pair))
	return res
}

func (s *Storage) genKey(pair models.Pair) string {
	return fmt.Sprintf("%d%d%s", pair.PriceHigh.CexType, pair.PriceLow.CexType, pair.Coin)
}
