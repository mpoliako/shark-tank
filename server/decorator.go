package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.corp.cloudsimple.com/cloudsimple/csos/incubator/shark-tank/types"
)

type JSONDecorator func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) (interface{}, int, error)

func WithJSONDecorator(h JSONDecorator) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		resp, code, err := h(w, req, params)
		if err != nil {
			handleError(err, code, w)
			return
		}
		if resp == nil {
			w.WriteHeader(code)
			return
		}
		data, err := json.Marshal(resp)
		if err != nil {
			handleError(err, http.StatusInternalServerError, w)
			return
		}
		w.WriteHeader(code)
		w.Write(data)
	}
}

func handleError(err error, code int, rw http.ResponseWriter) {
	rw.WriteHeader(code)
	response := types.ErrorResponse{
		Error: err.Error(),
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf(`failed to marshal error response: [err: %s, response: %+v]`, err, response)
	}
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write(data)
}
