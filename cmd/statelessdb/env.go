// Copyright (c) 2024. Jaakko Heusala <jheusala@iki.fi>. All rights reserved.
// Licensed under the FSL-1.1-MIT, see LICENSE.md in the project root for details.

package main

import (
	"os"
	"strconv"
)

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
