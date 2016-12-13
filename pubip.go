package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/thibran/pubip"
)

func main() {
	showVersion := flag.Bool("version", false, "print version")
	flag.Parse()
	if *showVersion {
		fmt.Printf("pubip 0.1   %s\n", runtime.Version())
		os.Exit(0)
	}
	m := pubip.NewMaster()
	m.Parallel = 4
	ip, err := m.Address()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ip)
}
