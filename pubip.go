package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/thibran/pubip"
)

type cfg struct {
	showVersion bool
	format      pubip.IPType
	cacheFile   string
}

const appVersion = "0.3"

var verbose = false

func main() {
	cfg := parseArgs()
	if cfg.showVersion {
		fmt.Printf("pubip %s   %s\n", appVersion, runtime.Version())
		os.Exit(0)
	}
	cfg.run()
}

func (cfg *cfg) run() {
	cache, err := loadCache(cfg.cacheFile)
	if err != nil {
		l("no cache")
		cfg.handleNotInCache(cache)
		return
	}
	l("from cache")
	ip, err := cache.maybeIP(cfg.format)
	if err != nil {
		l("ip not in cache or too old")
		cfg.handleNotInCache(cache)
		return
	}
	fmt.Println(ip)
}

func (cfg *cfg) handleNotInCache(cache *Cache) {
	ip, isipv6, err := cfg.ipFromInternet()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ip)
	cfg.writeToCache(ip, isipv6, cache)
}

type isipv6 bool

// fromInternet returns the IP address by requesting it online.
func (cfg *cfg) ipFromInternet() (string, isipv6, error) {
	m := pubip.NewMaster()
	m.Parallel = 2
	m.Format = cfg.format
	l(fmt.Sprintf("request %s address", cfg.format))
	ip, err := m.Address()
	if err != nil {
		return "", false, err
	}
	v6 := isipv6(!pubip.IsIPv4(ip))
	return ip, v6, nil
}

func (cfg *cfg) writeToCache(ip string, v6 isipv6, cache *Cache) {
	// cache might be nil
	if cache == nil {
		cache = new(Cache)
	}
	if v6 {
		l("set IPv6")
		cache.setIPv6(ip)
	} else {
		l("set IPv4")
		cache.setIPv4(ip)
	}
	if err := cache.save(cfg.cacheFile); err == nil {
		l("result cached")
	}
}

func parseArgs() cfg {
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
	return cfg{
		showVersion: *showVersion,
		format:      ipFormat,
		cacheFile:   cacheLocation(),
	}
}

// l logs string if verbose mode is active.
func l(s string) {
	if verbose {
		log.Println(s)
	}
}
