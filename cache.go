package main

import (
	"encoding/gob"
	"errors"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"github.com/thibran/pubip"
)

// Cache result of the ip addresses
type Cache struct {
	V6ip      string    `json:"v6ip"`
	V6last    time.Time `json:"v6last"`
	V4ip      string    `json:"v4ip"`
	V4last    time.Time `json:"v4last"`
	cacheFile string
	mut       sync.Mutex
}

var (
	errIPv6Empty   = errors.New("no IPv6 address in cache")
	errIPv4Empty   = errors.New("no IPv4 address in cache")
	errNotInCache  = errors.New("value not in cache")
	cacheTimeLimit = time.Duration(time.Minute * 15)
)

// loadCache from file. Returns always a non-nil cache object.
func loadCache(cacheFile string) *Cache {
	logln("read:", cacheFile)
	cache := &Cache{cacheFile: cacheFile}
	f, err := os.Open(cacheFile)
	if err != nil {
		logln(err)
		return cache
	}
	defer f.Close()
	decodeFrom(f, cache)
	return cache
}

// loadCache from file. Returns always a non-nil cache object.
func decodeFrom(r io.Reader, cache *Cache) {
	dec := gob.NewDecoder(r)
	if err := dec.Decode(cache); err != nil {
		logf("no cache file: %s\n", err)
		return
	}
	logln("cache file exists")
	return
}

// maybeIP returns the matching IP address for the given IPType t,
// if the cache entry is not older than cacheTimeLimit (15 min).
func (c *Cache) maybeIP(t pubip.IPType) (string, error) {
	c.mut.Lock()
	defer c.mut.Unlock()

	checkCache := func(last time.Time) bool {
		other := last.Add(cacheTimeLimit)
		return other.After(time.Now())
	}

	if t == pubip.IPv6 || t == pubip.IPv6orIPv4 {
		if checkCache(c.V6last) && len(c.V6ip) != 0 {
			return c.V6ip, nil
			// IPv6 only
		} else if t == pubip.IPv6 {
			return "", errIPv6Empty
		}
	}
	if t == pubip.IPv4 || t == pubip.IPv6orIPv4 {
		if checkCache(c.V4last) && len(c.V4ip) != 0 {
			return c.V4ip, nil
			// IPv4 only
		} else if t == pubip.IPv4 {
			return "", errIPv4Empty
		}
	}
	return "", errNotInCache
}

func (c *Cache) save() error {
	c.mut.Lock()
	defer c.mut.Unlock()
	f, err := os.Create(c.cacheFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return encodeTo(f, c)
}

func encodeTo(w io.Writer, c *Cache) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(c)
}

// setIPv6 non-empty IP string.
func (c *Cache) setIPv6(ip string) {
	c.V6ip = ip
	c.V6last = time.Now()
}

// setIPv4 non-empty IP string.
func (c *Cache) setIPv4(ip string) {
	c.V4ip = ip
	c.V4last = time.Now()
}

// cacheLocation path or panics.
func cacheLocation() string {
	var cdir string
	// snap dir
	if d := os.Getenv("SNAP_USER_COMMON"); d != "" {
		cdir = d
		// XDG cache
	} else if d := os.Getenv("XDG_CACHE_HOME"); d != "" {
		cdir = d
		// ~/.cache
	} else if d, err := dotCacheDir(); err == nil {
		cdir = d
		// tmp
	} else if d := os.TempDir(); d != "" {
		cdir = d
	}
	if cdir == "" {
		log.Fatalln("unknown tmp dir location")
	}
	return filepath.Join(cdir, "pubip.cache")
}

// returns maybe the path to the linux ~/.cache directory
func dotCacheDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	cache := filepath.Join(user.HomeDir, ".cache")
	fi, err := os.Stat(cache)
	if err == nil && fi.IsDir() {
		return cache, nil
	}
	return "", err
}
