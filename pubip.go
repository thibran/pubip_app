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
	v6 := flag.Bool("6", false, "only IPv6")
	v4 := flag.Bool("4", false, "only IPv4")
	flag.Parse()
	if *showVersion {
		fmt.Printf("pubip 0.2   %s\n", runtime.Version())
		os.Exit(0)
	}
	ipFormat := pubip.IPv6orIPv4
	if *v6 {
		ipFormat = pubip.IPv6
	} else if *v4 {
		ipFormat = pubip.IPv4
	}
	m := pubip.NewMaster()
	m.Parallel = 2
	m.Format = ipFormat
	ip, err := m.Address()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ip)
}
