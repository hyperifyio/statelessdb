// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package apis

import (
	"io"
	"net/http"
	"net/http/pprof"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"statelessdb/internal/dtos"
	"statelessdb/internal/encryptions"
	"statelessdb/internal/metrics"
	"statelessdb/internal/states"
)

func (s *Server) HandleComputeStateRequest(w http.ResponseWriter, r *http.Request) {
	var err error

	metrics.HttpRequestsTotal.WithLabelValues(r.URL.Path).Inc()

	now := states.NewTimeNow()

	// Initialize an instance of ComputeRequestDTO
	var computeRequest dtos.ComputeRequestDTO

	reader := GetJsonReaderState()
	defer reader.Release()

	// Decode the request body into computeRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Debugf("Failed to read body: %v", err)
		sendHttpError(w, BadBodyError, http.StatusBadRequest)
		return
	}

	reader.Buffer.Reset(body)
	if err = reader.Decoder.Decode(&computeRequest); err != nil {
		log.Errorf("Bad body error: %v", err)
		sendHttpError(w, BadBodyError, http.StatusBadRequest)
		return
	}

	// Read optional current state
	var state *states.ComputeState
	if computeRequest.Private != "" {
		decryptedState, err := states.DecryptComputeState(computeRequest.Private, s.Decryptor)
		if err != nil {
			log.Errorf("HandleComputeStateRequest: decrypting: error: %v", err)
			sendHttpError(w, DecryptionFailedError, http.StatusInternalServerError)
			return
		}
		state = decryptedState
	} else {
		metrics.ResourceCreatedTotal.WithLabelValues().Inc()
		public := computeRequest.Public
		//if public == nil {
		//	public = make(map[string]interface{})
		//}
		var private map[string]interface{}
		//private = make(map[string]interface{})
		state = states.New(uuid.New(), uuid.New(), now, now, public, private)
	}

	// Perform actions
	state.Updated = now

	// Prepare response
	var private string
	if private, err = state.Encrypt(s.Encryptor); err != nil {
		log.Errorf("HandleComputeStateRequest: encrypting: error: %v", err)
		sendHttpError(w, EncryptionFailedError, http.StatusInternalServerError)
		return
	}

	response := dtos.NewComputeStateDTO(
		state.Id,
		state.Owner,
		state.Created,
		state.Updated,
		state.Public,
		private,
	)

	// Set the Content-Type header.
	w.Header().Set("Content-Type", "application/json")

	// Serialize the map to JSON and write it to the response.
	encoderState := encryptions.GetJsonEncoderState()
	defer encoderState.Release()
	if err = encoderState.Encoder.Encode(response); err != nil {
		log.Errorf("HandleComputeStateRequest: encoding: error: %v", err)
		sendHttpError(w, EncodingFailedError, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(encoderState.Bytes()); err != nil {
		log.Errorf("HandleComputeStateRequest: writing: error: %v", err)
		sendHttpError(w, WritingBodyFailedError, http.StatusInternalServerError)
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
	r.HandleFunc("/api/v1", s.HandleComputeStateRequest)
	//r.PathPrefix("/").Handler(http.StripPrefix("/", wrappedFileServerHandler))

	err := http.ListenAndServe(listen, r)
	if err != nil {
		panic("failed to start StatelessDB server")
	}

}
