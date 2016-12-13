package main

import (
	"encoding/gob"
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/thibran/pubip"
)

// Cache result of the ip addresses
type Cache struct {
	V6ip   string    `json:"v6ip"`
	V6last time.Time `json:"v6last"`
	V4ip   string    `json:"v4ip"`
	V4last time.Time `json:"v4last"`
}

var (
	errIPv6Empty   = errors.New("no IPv6 address in cache")
	errIPv4Empty   = errors.New("no IPv4 address in cache")
	errNotInCache  = errors.New("value not in cache")
	cacheTimeLimit = time.Duration(15 * time.Minute)
)

const errTempDirUnknown = "unknown tmp dir location"

func loadCache(cacheFile string) (*Cache, error) {
	f, err := os.Open(cacheFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data Cache
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Cache) save(cacheFile string) error {
	f, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	return enc.Encode(c)
}

// maybeIP the matching IP address for the given IPType t, if the
// cache entry is not older than cacheTimeLimit (15 min).
func (c *Cache) maybeIP(t pubip.IPType) (string, error) {
	useCache := func(last time.Time) bool {
		now := time.Now()
		other := last.Add(cacheTimeLimit)
		return other.After(now)
	}
	if t == pubip.IPv6 || t == pubip.IPv6orIPv4 {
		if useCache(c.V6last) && len(c.V6ip) != 0 {
			return c.V6ip, nil
			// case IPv6 only
		} else if t == pubip.IPv6 {
			return "", errIPv6Empty
		}
	}
	if t == pubip.IPv4 || t == pubip.IPv6orIPv4 {
		if useCache(c.V4last) && len(c.V4ip) != 0 {
			return c.V4ip, nil
			// case IPv4 only
		} else if t == pubip.IPv4 {
			return "", errIPv4Empty
		}
	}
	return "", errNotInCache
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
	if d := os.Getenv("SNAP_USER_COMMON"); len(d) != 0 {
		cdir = d
		// XDG cache
	} else if d := os.Getenv("XDG_CACHE_HOME"); len(d) != 0 {
		cdir = d
		// ~/.cache
	} else if d, err := dotCache(); err == nil {
		cdir = d
		// tmp
	} else if d := os.TempDir(); len(d) != 0 {
		cdir = d
	}
	if len(cdir) == 0 {
		log.Fatalln(errTempDirUnknown)
	}
	return cdir + string(filepath.Separator) + "pubip.cache"
}

// returns maybe the path to the linux ~/.cache dir
func dotCache() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	cache := user.HomeDir + "/.cache"
	fi, err := os.Stat(cache)
	if err == nil && fi.IsDir() {
		return cache, nil
	}
	return "", err
}

// func loadCache(cacheFile string) (*Cache, error) {
// 	buf, err := ioutil.ReadFile(cacheFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var data Cache
// 	if err := json.Unmarshal(buf, &data); err != nil {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (c *Cache) save(cacheDir string) error {
// 	buf, err := json.Marshal(c)
// 	if err != nil {
// 		return err
// 	}
// 	f, err := os.Create(cacheDir)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	_, err = f.Write(buf)
// 	return err
// }