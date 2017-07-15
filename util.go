package main

import "log"

// logln logs if verbose mode is active.
func logln(a ...interface{}) {
	if verbose {
		log.Println(a...)
	}
}

// logf logs if verbose mode is active.
func logf(format string, a ...interface{}) {
	if verbose {
		log.Printf(format, a...)
	}
}
