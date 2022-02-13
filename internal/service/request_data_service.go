package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/NetworkPy/synergy_test_task/internal/model"
)

type RDSConfig struct {
	Endpoints           []model.MethodUrl
	CacheDataRepository model.CacheDataRepository
}

type requestDataService struct {
	Endpoints           []model.MethodUrl
	CacheDataRepository model.CacheDataRepository
	Client              *http.Client
}

func NewRequestDataService(config *RDSConfig) (model.RequestDataService, error) {
	client := &http.Client{}

	if len(config.Endpoints) == 0 {
		return nil, fmt.Errorf("found zero urls to request")
	}

	return &requestDataService{
		Endpoints:           config.Endpoints,
		CacheDataRepository: config.CacheDataRepository,
		Client:              client,
	}, nil
}

func (s *requestDataService) GetData(key int) (string, error) {
	var data model.Data

	for {
		dataByte, err := s.CacheDataRepository.GetData(key)
		if err != nil {
			log.Println(err)
			return "", err
		}
		if err := json.Unmarshal(dataByte, &data); err != nil {
			log.Println(err)
			return "", fmt.Errorf("bad data format to response")
		}

		if data.Action != nil {
			return *data.Action, nil
		}
	}

}

func (s *requestDataService) Start() {
	for _, m := range s.Endpoints {
		go s.startMonitoring(m.Method, m.Url)
	}
}

func (s *requestDataService) startMonitoring(method string, url string) {
	log.Printf("start: %s", url)
	timeout := RandomNumber()
	ticker := time.NewTicker(time.Second * time.Duration(int64(timeout)))
	defer ticker.Stop()

	for range ticker.C {
		timeout := RandomNumber()
		ticker.Reset(time.Second * time.Duration(int64(timeout)))

		b, err := s.doRequest(method, url)
		if err != nil {
			log.Printf("bad request to other service: %s", url)
			continue
		}
		if len(b) == 0 {
			log.Printf("body is empty: %s", url)
			continue
		}
		s.CacheDataRepository.SetData(0, b)
	}
}

func (s *requestDataService) doRequest(method string, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func RandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 3
	return rand.Intn(max-min+1) + min
}
