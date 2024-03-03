package main

import "log"

var debug bool

func debugLog(format string, args ...interface{}) {
	if debug {
		log.Printf(format, args...)
	}
}
