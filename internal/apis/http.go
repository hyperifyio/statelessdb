// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package apis

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"

	"statelessdb/internal/dtos"
	"statelessdb/internal/metrics"
	"statelessdb/internal/states"
)

func (s *Server) handleComputeStateRequest(w http.ResponseWriter, r *http.Request) {

	metrics.HttpRequestsTotal.WithLabelValues(r.URL.Path).Inc()

	now := states.NewTimeNow()

	// Initialize an instance of ComputeRequestDTO
	var computeRequest dtos.ComputeRequestDTO

	// Decode the request body into computeRequest
	err := json.NewDecoder(r.Body).Decode(&computeRequest)
	if err != nil {
		sendHttpError(w, BadBodyError, http.StatusBadRequest)
		return
	}

	// Read current state

	var state *states.ComputeState
	if computeRequest.Payload != nil {

		if computeRequest.Payload.Private == "" {
			sendHttpError(w, BadPrivateBodyError, http.StatusBadRequest)
			return
		}

		decryptedState, err := states.DecryptComputeState(computeRequest.Payload.Private, s.Decryptor)
		if err != nil {
			log.Errorf("handleComputeStateRequest: decrypting: error: %v", err)
			sendHttpError(w, DecryptionFailedError, http.StatusInternalServerError)
			return
		}

		state = decryptedState

	} else {
		metrics.ResourceCreatedTotal.WithLabelValues().Inc()
		var public map[string]interface{}
		var private map[string]interface{}
		state = states.New(uuid.New(), uuid.New(), now, now, public, private)
	}

	// Perform actions

	// Prepare response
	private, err := state.Encrypt(s.Encryptor)
	if err != nil {
		log.Errorf("handleComputeStateRequest: encrypting: error: %v", err)
		sendHttpError(w, EncryptionFailedError, http.StatusInternalServerError)
		return
	}

	var public map[string]interface{}

	response := dtos.ComputeStateDTO{
		Id:      state.Id.String(),
		Owner:   state.Owner.String(),
		Public:  public,
		Private: private,
	}

	// Set the Content-Type header.
	w.Header().Set("Content-Type", "application/json")

	// Serialize the map to JSON and write it to the response.
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorf("handleComputeStateRequest: encoding: error: %v", err)
		sendHttpError(w, EncodingFailedError, http.StatusInternalServerError)
		return
	}

}

func (s *Server) StartLocalServer(listen string, enablePprof bool) {

	r := mux.NewRouter()

	//// Wrap the file server handler to track requests using Prometheus
	//fileServerHandler := http.FileServer(http.FS(frontends.BuildFS))
	//wrappedFileServerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	metrics.HttpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
	//	fileServerHandler.ServeHTTP(w, r)
	//})

	// Register pprof routes with mux
	if enablePprof {
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
		r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
		log.Infof("Enabled: /debug/pprof/")
	}

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/api/v1", s.handleComputeStateRequest)
	//r.PathPrefix("/").Handler(http.StripPrefix("/", wrappedFileServerHandler))

	err := http.ListenAndServe(listen, r)
	if err != nil {
		panic("failed to start StatelessDB server")
	}

}
