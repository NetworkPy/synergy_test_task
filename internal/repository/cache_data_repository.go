package repository

import (
	"fmt"
	"sync"

	"github.com/NetworkPy/synergy_test_task/internal/model"
)

type cacheDataRepository struct {
	Cache map[int][]byte
	Mu    *sync.RWMutex
}

func NewCacheDataRepository() model.CacheDataRepository {
	cache := make(map[int][]byte)
	mu := &sync.RWMutex{}
	return &cacheDataRepository{
		Cache: cache,
		Mu:    mu,
	}
}

func (r *cacheDataRepository) GetData(key int) ([]byte, error) {
	r.Mu.RLock()
	defer r.Mu.RUnlock()

	if dataByte, exist := r.Cache[key]; exist {
		return dataByte, nil
	}

	return nil, fmt.Errorf(fmt.Sprintf("no data with key = %d", key))
}

func (r *cacheDataRepository) SetData(key int, data []byte) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.Cache[key] = data
}
