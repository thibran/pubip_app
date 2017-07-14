package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/thibran/pubip"
)

type app struct {
	showVersion bool
	format      pubip.IPType
	cacheFile   string
}

const appVersion = "0.5"

var verbose = false

func main() {
	app := parseArgs()
	if app.showVersion {
		fmt.Printf("pubip %s   %s\n", appVersion, runtime.Version())
		os.Exit(0)
	}
	app.run()
}

func (ap *app) run() {
	cache, err := loadCache(ap.cacheFile)
	if err != nil {
		logfn("no cache entry")
		ap.handleNotInCache(cache)
		return
	}
	logfn("cache entry exists")
	ip, err := cache.maybeIP(ap.format)
	if err != nil {
		logfn("ip not in cache or too old")
		ap.handleNotInCache(cache)
		return
	}
	fmt.Println(ip)
}

func (ap *app) handleNotInCache(cache *Cache) {
	ip, isipv6, err := ap.ipFromInternet()
	if err != nil {
		log.Fatalf("ip not in cache: %s\n", err)
	}
	fmt.Println(ip)
	ap.writeToCache(ip, isipv6, cache)
}

type isipv6 bool

// fromInternet returns the IP address by requesting it online.
func (ap *app) ipFromInternet() (string, isipv6, error) {
	m := pubip.NewMaster()
	m.Parallel = 2
	m.Format = ap.format
	logf("request %s address", ap.format)
	ip, err := m.Address()
	if err != nil {
		return "", false, err
	}
	v6 := isipv6(!pubip.IsIPv4(ip))
	return ip, v6, nil
}

func (ap *app) writeToCache(ip string, v6 isipv6, cache *Cache) {
	// cache might be nil
	if cache == nil {
		cache = new(Cache)
	}
	if v6 {
		logfn("set IPv6")
		cache.setIPv6(ip)
	} else {
		logfn("set IPv4")
		cache.setIPv4(ip)
	}
	if err := cache.save(ap.cacheFile); err == nil {
		logfn("result cached")
	}
}

func parseArgs() app {
	showVersion := flag.Bool("version", false, "print version")
	v6 := flag.Bool("6", false, "only IPv6")
	v4 := flag.Bool("4", false, "only IPv4")
	flag.BoolVar(&verbose, "v", false, "print verbose info about app execution")
	flag.Parse()
	ipFormat := pubip.IPv6orIPv4
	if *v6 {
		ipFormat = pubip.IPv6
	} else if *v4 {
		ipFormat = pubip.IPv4
	}
	return app{
		showVersion: *showVersion,
		format:      ipFormat,
		cacheFile:   cacheLocation(),
	}
}
