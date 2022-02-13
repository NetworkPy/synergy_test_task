package handler

import (
	"net/http"

	"github.com/NetworkPy/synergy_test_task/internal/model"
)

type DHConfig struct {
	DataService model.RequestDataService
	Mux         *http.ServeMux
}

type DataHandler struct {
	DataService model.RequestDataService
}

func NewDataHandler(config *DHConfig) {
	handler := &DataHandler{
		DataService: config.DataService,
	}

	config.Mux.HandleFunc("/data", handler.GetData)
}

func (h *DataHandler) GetData(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := h.DataService.GetData(0)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
		}
		w.Write([]byte(data))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("resource not found"))
	}
}
