package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/thibran/pubip"
)

const appVersion = "0.6"

var verbose = false

type app struct {
	showVersion bool
	showBoth    bool
	ipType      pubip.IPType
	cacheFile   string
}

func main() {
	app := parseFlags()
	if app.showVersion {
		fmt.Printf("pubip %s   %s\n", appVersion, runtime.Version())
		os.Exit(0)
	}
	cache := loadCache(app.cacheFile)
	if app.showBoth {
		handleShowBoth(cache)
		os.Exit(0)
	}
	handleDefault(cache, app.ipType)
}

func handleDefault(cache *Cache, ipType pubip.IPType) {
	ip, err := getIP(cache, ipType)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ip)
}

func handleShowBoth(cache *Cache) {
	errc := make(chan error)
	v6c := make(chan string)
	v4c := make(chan string)
	withChan := func(result chan<- string, t pubip.IPType) {
		ip, err := getIP(cache, t)
		if err == nil {
			result <- ip
		} else {
			errc <- err
		}
	}
	go withChan(v6c, pubip.IPv6)
	go withChan(v4c, pubip.IPv4)
	go func() {
		for err := range errc {
			log.Fatalln(err)
		}
		close(v6c)
		close(v4c)
	}()
	fmt.Printf("IPv6: %s\nIPv4: %s\n", <-v6c, <-v4c)
	close(errc)
}

func getIP(cache *Cache, ipType pubip.IPType) (string, error) {
	ip, err := cache.maybeIP(ipType)
	if err != nil {
		logf("%s - not in cache or too old\n", ipType)
		m := pubip.NewMaster()
		m.Parallel = 2
		m.Format = ipType
		logf("request %s address\n", ipType)
		ip, err = m.Address()
		if err == nil {
			writeToCache(cache, ip, ipType)
		}
	}
	return ip, err
}

func writeToCache(cache *Cache, ip string, ipType pubip.IPType) error {
	if ipType == pubip.IPv6 {
		cache.setIPv6(ip)
	} else {
		cache.setIPv4(ip)
	}
	if err := cache.save(); err != nil {
		return err
	}
	logln("wrote to cache")
	return nil
}

func parseFlags() app {
	showVersion := flag.Bool("version", false, "print version")
	v6 := flag.Bool("6", false, "only IPv6")
	v4 := flag.Bool("4", false, "only IPv4")
	both := flag.Bool("both", false, "IPv6 and IPv4")
	flag.BoolVar(&verbose, "v", false, "print verbose info about app execution")
	flag.Parse()
	ipType := pubip.IPv6orIPv4
	if *v6 {
		ipType = pubip.IPv6
	} else if *v4 {
		ipType = pubip.IPv4
	}
	return app{
		showVersion: *showVersion,
		showBoth:    *both,
		ipType:      ipType,
		cacheFile:   cacheLocation(),
	}
}
