package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/services/access"
)

type ItemsListHandler struct {
	log           *logger.Logger
	accessService access.AccessService
	manager       repository.Repository
}

func NewItemsListHandler(accessService access.AccessService, manager repository.Repository, log *logger.Logger) *ItemsListHandler {
	return &ItemsListHandler{
		accessService: accessService,
		manager:       manager,
		log:           log,
	}
}

type itemDataResponse struct {
	model_data.ItemDataResponse
}

type listDataItemsResponse struct {
	Items []itemDataResponse `json:"items"`
}

func (hr listDataItemsResponse) Render(res http.ResponseWriter, req *http.Request) error {
	return nil
}

func (hr itemDataResponse) Render(res http.ResponseWriter, req *http.Request) error {
	return nil
}

func (ih *ItemsListHandler) HandleItemsList(res http.ResponseWriter, req *http.Request) {

	userUUID, err := ih.accessService.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		ih.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}
	req.URL.Query().Get("offset")
	offset := req.URL.Query().Get("offset")
	limit := req.URL.Query().Get("limit")
	if offset == "" {
		offset = "0"
	}
	if limit == "" {
		limit = "200"
	}
	o, _ := strconv.Atoi(offset)
	l, _ := strconv.Atoi(limit)
	dataList, err := ih.manager.Owner().AllOwnerData(req.Context(), userUUID, o, l)
	if err != nil {
		ih.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
	}

	items := make([]itemDataResponse, 0)
	n := 1
	for _, data := range dataList {
		item := itemDataResponse{}
		item.UUID = data.DataUUID
		item.Name = data.DataName
		item.Type = data.DataTypeName
		item.Number = strconv.Itoa(n)
		items = append(items, item)
		n++
	}

	response := new(listDataItemsResponse)
	response.Items = items

	err = render.Render(res, req, response)
	if err != nil {
		ih.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
	}
}
