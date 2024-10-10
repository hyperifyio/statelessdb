// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package apis

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"net/http/pprof"
	"statelessdb/internal/encryptions"
	"statelessdb/pkg/metrics"
	"statelessdb/pkg/requests"
)

type Server struct {
	enablePprof bool
	routes      map[string]requests.ResponseManager
}

func NewServer() *Server {
	return &Server{
		false,
		make(map[string]requests.ResponseManager),
	}
}

func (s *Server) Handle(path string, r requests.ResponseManager) {
	s.routes[path] = r
}

func (s *Server) EnablePprof() {
	s.enablePprof = true
}

func (s *Server) BuildHandler(handler requests.ResponseManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.HttpRequestsTotal.WithLabelValues(r.URL.Path).Inc()

		// Read the request body
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("[Server.BuildHandler]: Failed to read body: %v", err)
			sendHttpError(w, BadBodyError, http.StatusBadRequest)
			return
		}

		dto, err := handler.ProcessBytes(requestBody)
		if err != nil {
			log.Errorf("[Server.BuildHandler]: Failed to read body: %v", err)
			sendHttpError(w, BadBodyError, http.StatusBadRequest)
			return
		}

		// Prepare response as JSON
		encoderState := encryptions.GetJsonEncoderState()
		defer encoderState.Release()
		if err = encoderState.Encoder.Encode(dto); err != nil {
			log.Errorf("[Server.BuildHandler]: encoding: error: %v", err)
			sendHttpError(w, WritingBodyFailedError, http.StatusInternalServerError)
			return
		}

		// Write response bytes to the HTTP request
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(encoderState.Bytes()); err != nil {
			log.Errorf("[Server.BuildHandler]: writing: error: %v", err)
			sendHttpError(w, WritingBodyFailedError, http.StatusInternalServerError)
			return
		}

	}
}

func (s *Server) StartLocalServer(listen string) {

	r := mux.NewRouter()

	//// Wrap the file server handler to track requests using Prometheus
	//fileServerHandler := http.FileServer(http.FS(frontends.BuildFS))
	//wrappedFileServerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	metrics.HttpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
	//	fileServerHandler.ServeHTTP(w, r)
	//})

	// Register pprof routes with mux
	if s.enablePprof {
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
		r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
		log.Infof("Enabled: /debug/pprof/")
	}

	r.Handle("/metrics", promhttp.Handler())

	for path, handler := range s.routes {
		r.HandleFunc(path, s.BuildHandler(handler))
	}

	//r.PathPrefix("/").Handle(http.StripPrefix("/", wrappedFileServerHandler))

	err := http.ListenAndServe(listen, r)
	if err != nil {
		panic("failed to start StatelessDB server")
	}

}
