// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"encoding/hex"

	"statelessdb/internal/apis"
	"statelessdb/internal/encryptions"
	"statelessdb/internal/states"

	statelessdb "statelessdb"
)

import _ "net/http/pprof"

func main() {

	var err error

	// Define flags
	addr := flag.String("addr", "", "change default address to listen")
	privateKeyString := flag.String("private-key", parseStringEnv("PRIVATE_KEY", ""), "set private key")
	port := flag.Int("port", parseIntEnv("PORT", 3001), "change default port")
	version := flag.Bool("version", false, "Show version information")
	enablePprof := flag.Bool("pprof", parseBooleanEnv("ENABLE_PPROF", false), "Enable pprof for debugging")
	initPrivateKey := flag.Bool("init-private-key", false, "Create a new private key and print it")

	listenTo := fmt.Sprintf("%s:%d", *addr, *port)

	// Parse the flags
	flag.Parse()

	if *version {
		fmt.Printf("%s v%s by %s\nURL = %s\n", statelessdb.Name, statelessdb.Version, statelessdb.Author, statelessdb.URL)
		return
	}

	if *initPrivateKey {
		key, err := encryptions.GenerateKey(32) // AES-256
		if err != nil {
			log.Errorf("Failed to generate key: %v", err)
			os.Exit(1)
		} else {
			fmt.Printf("PRIVATE_KEY=%s\n", hex.EncodeToString(key))
		}
		return
	}

	var serverKey []byte
	if *privateKeyString == "" {
		key, err := encryptions.GenerateKey(32) // AES-256
		if err != nil {
			log.Errorf("Failed to generate key: %v", err)
			os.Exit(1)
		} else {
			log.Warnf("Initialized with a random private key '%s'. You might want to make this persistent.", hex.EncodeToString(key))
			serverKey = key
		}
	} else {
		serverKey, err = hex.DecodeString(*privateKeyString)
		if err != nil {
			log.Errorf("Failed to decode private key: %v", err)
			os.Exit(1)
		}
	}

	serializer := encryptions.NewJsonSerializer[*states.ComputeState]("ComputeState")
	unserializer := encryptions.NewJsonUnserializer[*states.ComputeState]("ComputeState")

	encryptor := encryptions.NewEncryptor[*states.ComputeState](serializer)
	err = encryptor.Initialize(serverKey)
	if err != nil {
		log.Errorf("Failed to initialize encryptor: %v", err)
		os.Exit(1)
	}

	decryptor := encryptions.NewDecryptor[*states.ComputeState](unserializer)
	err = decryptor.Initialize(serverKey)
	if err != nil {
		log.Errorf("Failed to initialize decryptor: %v", err)
		os.Exit(1)
	}

	server := apis.NewServer(encryptor, decryptor)

	log.Infof("Starting server at %s", listenTo)
	server.StartLocalServer(listenTo, *enablePprof)

}

func parseIntEnv(key string, defaultValue int) int {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return result
}

func parseStringEnv(key string, defaultValue string) string {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}
	return str
}

func parseBooleanEnv(key string, defaultValue bool) bool {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}
	switch str {
	case "0":
		return false
	case "f":
		return false
	case "false":
		return false
	case "off":
		return false
	case "null":
		return false
	case "1":
		return true
	case "t":
		return true
	case "true":
		return true
	case "on":
		return true
	}
	return false
}
