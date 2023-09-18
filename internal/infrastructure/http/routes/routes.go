package routes

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/antsrp/fio_service/internal/domain"
	"github.com/antsrp/fio_service/internal/domain/graphql"
	"github.com/antsrp/fio_service/internal/infrastructure/messages"
	igql "github.com/antsrp/fio_service/internal/interfaces/graphql"
	"github.com/antsrp/fio_service/internal/mapper"
	"github.com/antsrp/fio_service/internal/usecases/service"
	"github.com/go-chi/chi"

	"github.com/graphql-go/handler"
)

type Handler struct {
	srv *service.Service
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h Handler) Routes(api igql.PersonAPI) chi.Router {

	//fileServer := http.FileServer(http.Dir(service.REPORTS_RELATIVE_PATH))

	gqlHandler := handler.New(&handler.Config{
		Schema:     graphql.InitSchema(api),
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/api/v1/fio", h.get)
		r.Delete("/api/v1/fio", h.delete)
		r.Patch("/api/v1/fio", h.update)
		r.Post("/api/v1/fio", h.add)
		r.Handle("/graphql", gqlHandler)
	})

	return r
}

func (h Handler) readBody(r *http.Request) []byte {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.srv.Logger.Errorf("cannot read body: %v", err.Error())
		return nil
	}

	return body
}

func (h Handler) getQueryParam(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

func (h Handler) parseID(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		h.srv.Logger.Errorf("can't parse id from request: %v", err.Error())
		val = -1
	}
	return val
}

func (h Handler) sendResult(w http.ResponseWriter, code int, data []byte) {
	w.WriteHeader(code)
	w.Write(data)
}

func (h Handler) sendError(w http.ResponseWriter, isInternal bool) {
	if isInternal {
		h.sendResult(w, http.StatusInternalServerError, []byte(messages.InternalError))
	} else {
		h.sendResult(w, http.StatusBadRequest, []byte(messages.InvalidInput))
	}
}

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	where, order := h.getQueryParam(r, "where"), h.getQueryParam(r, "order")
	var page int
	if s := h.getQueryParam(r, "page"); s != "" {
		page = h.parseID(s)
	}

	var wheres, orders []string
	if where != "" {
		wheres = strings.Split(where, ",")
	}
	if order != "" {
		orders = strings.Split(order, ",")
	}

	persons, err := h.srv.GetPersons(wheres, orders, page)
	if err != nil {
		h.srv.Logger.Errorf("can't get persons data: %v", err.Cause.Error())
		h.sendError(w, err.IsInternal)
		return
	}

	data, mapErr := mapper.ToJSON[[]domain.Person](persons, mapper.NewIndent("", "\t"))
	if mapErr != nil {
		h.srv.Logger.Errorf("can't map persons data to json: %v", mapErr.Error())
		h.sendResult(w, http.StatusBadRequest, []byte(messages.InvalidInput))
		return
	}

	w.Header().Add("Content-type", "application/json")
	h.sendResult(w, http.StatusOK, data)
}

func (h Handler) update(w http.ResponseWriter, r *http.Request) {
	data := h.readBody(r)
	if data == nil {
		h.sendResult(w, http.StatusBadRequest, []byte(messages.InvalidInput))
		return
	}

	person, mapErr := mapper.FromJSON[domain.Person](data)
	if mapErr != nil {
		h.srv.Logger.Errorf("can't map data from json to person: %v", mapErr.Error())
		h.sendResult(w, http.StatusBadRequest, []byte(messages.InvalidInput))
		return
	}

	rowsAffected, err := h.srv.UpdatePerson(*person)
	if err != nil {
		h.srv.Logger.Errorf("can't update person: %v", err.Cause.Error())
		h.sendError(w, err.IsInternal)
		return
	}
	if rowsAffected == 0 {
		h.sendResult(w, http.StatusOK, []byte(messages.NoRowsAffected))
	} else {
		h.sendResult(w, http.StatusOK, []byte(fmt.Sprintf(messages.SuccessfulUpdateMsg, person.Id)))
	}
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := -1
	if str := h.getQueryParam(r, "id"); str != "" {
		id = h.parseID(str)
	}
	if id == -1 {
		h.sendResult(w, http.StatusBadRequest, []byte(messages.InvalidInput))
		return
	}

	rowsAffected, err := h.srv.DeletePerson(id)
	if err != nil {
		h.srv.Logger.Errorf("can't delete person: %v", err.Cause.Error())
		h.sendError(w, err.IsInternal)
		return
	}
	if rowsAffected == 0 {
		h.sendResult(w, http.StatusOK, []byte(messages.NoRowsAffected))
	} else {
		h.sendResult(w, http.StatusOK, []byte(fmt.Sprintf(messages.SuccessfulDeleteMsg, id)))
	}
}

func (h Handler) add(w http.ResponseWriter, r *http.Request) {
	data := h.readBody(r)
	if data == nil {
		h.sendResult(w, http.StatusBadRequest, []byte(messages.InvalidInput))
		return
	}

	person, err := mapper.FromJSON[domain.PersonCommon](data)
	if err != nil {
		h.srv.Logger.Errorf("can't map data from json to person: %v", err.Error())
		h.sendResult(w, http.StatusBadRequest, []byte(messages.InvalidInput))
		return
	}
	if err := h.srv.SendToBroker(*person); err != nil {
		h.srv.Logger.Errorf("can't send data to broker: %v", err.Error())
		h.sendResult(w, http.StatusInternalServerError, []byte(messages.InternalError))
		return
	}

	h.sendResult(w, http.StatusAccepted, []byte(messages.SendToAddMsg))
}
