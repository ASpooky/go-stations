package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HealthzResponse := &model.HealthzResponse{}
	HealthzResponse.Message = "OK"

	e := json.NewEncoder(w)

	hr := HealthzResponse

	if err := e.Encode(hr); err != nil {
		fmt.Println(err)
		return
	}
}
