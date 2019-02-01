package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"gitlab.corp.cloudsimple.com/cloudsimple/csos/incubator/shark-tank/storage"
	"gitlab.corp.cloudsimple.com/cloudsimple/csos/incubator/shark-tank/types"
)

type ExampleServer struct {
	repo   *storage.CassandraRepository
	router *httprouter.Router
}

func NewServer(repo *storage.CassandraRepository) *ExampleServer {
	s := &ExampleServer{
		router: httprouter.New(),
		repo:   repo,
	}
	s.router.GET("/test/entity", WithJSONDecorator(s.List))
	s.router.POST("/test/entity", WithJSONDecorator(s.Create))
	s.router.POST("/test/inject", WithJSONDecorator(s.Inject))
	s.router.GET("/test/entity/:id", WithJSONDecorator(s.Get))
	return s
}

func (s *ExampleServer) Run(port string) error {
	log.Printf("start server on %s port", port)
	return http.ListenAndServe(port, s.router)
}

func CorrelationIDHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		next(rw, r)
	}
}

func (s *ExampleServer) List(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) (interface{}, int, error) {
	resp, err := s.repo.List(r.Context())
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return resp, http.StatusOK, nil
}

func (s *ExampleServer) Get(rw http.ResponseWriter, r *http.Request, params httprouter.Params) (interface{}, int, error) {
	uuid, err := gocql.ParseUUID(params.ByName("id"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	resp, err := s.repo.Get(r.Context(), uuid)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return resp, http.StatusOK, nil
}

func (s *ExampleServer) Create(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) (interface{}, int, error) {
	item := &types.ExampleEntity{}
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		return nil, http.StatusBadRequest, err
	}
	if err := s.repo.Create(r.Context(), item); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return item, http.StatusOK, nil
}

func (s *ExampleServer) Inject(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) (interface{}, int, error) {
	item := &types.ExampleEntity{}
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		return nil, http.StatusBadRequest, err
	}
	if err := s.repo.Inject(r.Context(), item); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return item, http.StatusOK, nil
}
