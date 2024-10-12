// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"

	"statelessdb/pkg/apis"
	"statelessdb/pkg/dtos"
	"statelessdb/pkg/encodings"
	"statelessdb/pkg/events"
	"statelessdb/pkg/requests"
	"statelessdb/pkg/states"

	statelessdb "statelessdb"
)

import _ "net/http/pprof"

const (
	eventTimeoutTime         = time.Second * 10 // Default request timeout time
	eventExpirationTime      = 20 * time.Second // Time until events expire
	eventCleanupIntervalTime = 30 * time.Second // Interval to clean up expired events
)

func main() {

	var err error

	// Define flags
	addr := flag.String("addr", "", "change default address to listen")
	port := flag.Int("port", parseIntEnv("PORT", 3001), "change default port")
	privateKeyString := flag.String("private-key", parseStringEnv("PRIVATE_KEY", ""), "set private key")
	version := flag.Bool("version", false, "Show version information")
	enablePprof := flag.Bool("pprof", parseBooleanEnv("ENABLE_PPROF", false), "Enable pprof for debugging")
	initPrivateKey := flag.Bool("init-private-key", false, "Create a new private key and print it")

	// Parse flags
	flag.Parse()

	// Handle --addr and --port
	listenTo := fmt.Sprintf("%s:%d", *addr, *port)

	// Handle --version
	if *version {
		fmt.Printf("%s v%s by %s\nURL = %s\n", statelessdb.Name, statelessdb.Version, statelessdb.Author, statelessdb.URL)
		return
	}

	// Handle --init-private-key
	if *initPrivateKey {
		key, err := encodings.GenerateKey(32) // AES-256
		if err != nil {
			log.Errorf("Failed to generate key: %v", err)
			os.Exit(1)
		} else {
			fmt.Printf("PRIVATE_KEY=%s\n", hex.EncodeToString(key))
		}
		return
	}

	// Handle --private-key
	serverKey, err := parsePrivateKeyString(*privateKeyString)
	if err != nil {
		log.Errorf("Private key parsing failed: %v", err)
		os.Exit(1)
	}

	// Define the server

	newState := func() *states.ComputeState {
		return &states.ComputeState{}
	}

	newRequestDTO := func() *requests.ComputeRequest {
		return &requests.ComputeRequest{}
	}

	computeRequestManager, err := requests.NewJsonRequestManager[*states.ComputeState, *requests.ComputeRequest, *dtos.ComputeResponseDTO](
		"ComputeState",
		serverKey,
		newState,
		newRequestDTO,
	)
	if err != nil {
		log.Errorf("Failed to initialize JSON request handler: %v", err)
		os.Exit(1)
	}

	eventBus := events.NewLocalEventBus[uuid.UUID, interface{}]()

	server := apis.NewServer()
	if *enablePprof {
		server.EnablePprof()
	}
	server.Handle("/api/v1", computeRequestManager.HandleWith(ApiRequestHandler(eventBus)).WithResponse(NewComputeResponseDTO(eventBus)).WithMethods("GET", "POST"))
	server.Handle("/api/v1/events", computeRequestManager.HandleWith(ApiEventHandler(eventBus, eventTimeoutTime, eventExpirationTime, eventCleanupIntervalTime)).WithResponse(NewEventResponseDTO(eventBus)).WithMethods("GET", "POST"))

	// Start the server
	log.Infof("Starting server at %s", listenTo)
	server.StartLocalServer(listenTo)

}
