package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	result, _ := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	res := model.CreateTODOResponse{}
	res.TODO.ID = result.ID
	res.TODO.Subject = result.Subject
	res.TODO.Description = result.Description
	res.TODO.CreatedAt = result.CreatedAt
	res.TODO.UpdatedAt = result.UpdatedAt
	return &res, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == "POST" {
		request := model.CreateTODORequest{}
		len := r.ContentLength
		body := make([]byte, len)
		r.Body.Read(body)
		err := json.Unmarshal(body, &request)
		if err != nil {
			log.Fatal(err)
		}
		if request.Subject == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else {
			ctx := r.Context()
			response, err := h.Create(ctx, &request)
			if err != nil {
				log.Fatal(err)
				return
			}
			e := json.NewEncoder(w)
			if err := e.Encode(response); err != nil {
				fmt.Println(err)
				return
			}
		}

	}
}
