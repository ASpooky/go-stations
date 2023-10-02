package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

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
	result, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
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
	result, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	res := model.ReadTODOResponse{}
	res.TODOs = result
	return &res, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	result, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	res := model.UpdateTODOResponse{}
	res.TODO.ID = result.ID
	res.TODO.Subject = result.Subject
	res.TODO.Description = result.Description
	res.TODO.CreatedAt = result.CreatedAt
	res.TODO.UpdatedAt = result.UpdatedAt
	return &res, nil
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
			fmt.Println(err.Error())
		}
		if request.Subject == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else {
			ctx := r.Context()
			response, err := h.Create(ctx, &request)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			e := json.NewEncoder(w)
			if err := e.Encode(response); err != nil {
				fmt.Println(err)
				return
			}
		}
	} else if method == "PUT" {
		request := model.UpdateTODORequest{}
		len := r.ContentLength
		body := make([]byte, len)
		r.Body.Read(body)
		err := json.Unmarshal(body, &request)
		if err != nil {
			fmt.Println(err.Error())
		}
		if request.Subject == "" || request.ID == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else {
			ctx := r.Context()
			response, err := h.Update(ctx, &request)
			if err != nil {
				fmt.Println(err.Error())
				if reflect.TypeOf(err) == reflect.TypeOf(&model.ErrNotFound{}) {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				return
			}
			e := json.NewEncoder(w)
			if err := e.Encode(response); err != nil {
				fmt.Println(err)
				return
			}
		}
	} else if method == "GET" {
		request := model.ReadTODORequest{Size: 3}
		prevId := r.FormValue("prev_id")
		size := r.FormValue("size")
		if len(prevId) > 0 {

			intPrevId, err := strconv.Atoi(prevId)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			request.PrevID = int64(intPrevId)
		}
		if len(size) > 0 {
			intSize, err := strconv.Atoi(size)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			request.Size = int64(intSize)
		}

		ctx := r.Context()
		response, err := h.Read(ctx, &request)
		if err != nil {
			fmt.Println(err.Error())
			if reflect.TypeOf(err) == reflect.TypeOf(&model.ErrNotFound{}) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			return
		}
		e := json.NewEncoder(w)
		if err := e.Encode(response); err != nil {
			fmt.Println(err)
			return
		}
		return
	}
}
