package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/thibran/pubip"
)

const appVersion = "0.5"

var verbose = false

type app struct {
	showVersion bool
	showBoth    bool
	format      pubip.IPType
	cacheFile   string
}

func main() {
	app := parseFlags()
	if app.showVersion {
		fmt.Printf("pubip %s   %s\n", appVersion, runtime.Version())
		os.Exit(0)
	}
	if app.showBoth {
		app.format = pubip.IPv6
		v6 := app.run()
		app.format = pubip.IPv4
		v4 := app.run()
		fmt.Printf("IPv6: %s\nIPv4: %s\n", v6, v4)
		os.Exit(0)
	}
	fmt.Println(app.run())
}

func (ap *app) run() string {
	cache, err := loadCache(ap.cacheFile)
	if err != nil {
		logln("no cache entry")
		return ap.handleNotInCache(cache)
	}
	logln("cache entry exists")
	ip, err := cache.maybeIP(ap.format)
	if err != nil {
		logln("ip not in cache or too old")
		return ap.handleNotInCache(cache)
	}
	return ip
}

func (ap *app) handleNotInCache(cache *Cache) string {
	ip, isipv6, err := ap.ipFromInternet()
	if err != nil {
		log.Fatalf("ip not in cache: %s\n", err)
	}
	ap.writeToCache(ip, isipv6, cache)
	return ip
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
	return ip, isipv6(!pubip.IsIPv4(ip)), nil
}

func (ap *app) writeToCache(ip string, v6 isipv6, cache *Cache) {
	if cache == nil {
		cache = new(Cache)
	}
	if v6 {
		logln("set IPv6")
		cache.setIPv6(ip)
	} else {
		logln("set IPv4")
		cache.setIPv4(ip)
	}
	if err := cache.save(ap.cacheFile); err == nil {
		logln("result cached")
	}
}

func parseFlags() app {
	showVersion := flag.Bool("version", false, "print version")
	v6 := flag.Bool("6", false, "only IPv6")
	v4 := flag.Bool("4", false, "only IPv4")
	both := flag.Bool("both", false, "IPv6 and IPv4")
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
		showBoth:    *both,
		format:      ipFormat,
		cacheFile:   cacheLocation(),
	}
}
