package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/thibran/pubip"
)

func TestCacheLocation(t *testing.T) {
	dir := cacheLocation()
	fmt.Println(dir)
	if len(dir) == 0 {
		t.Fatal()
	}
}

func TestLoadCache_fail(t *testing.T) {
	cache := loadCache("/foo")
	if cache.V6ip != "" || cache.V4ip != "" {
		t.Fatal()
	}
}

func TestEncodeTo(t *testing.T) {
	var buf bytes.Buffer
	cacheFile := "/baz/zot/pubip.cache"
	req := "2003:c8:9bec:4a70:e755:6628:78d1:aa06"
	cache := Cache{cacheFile: cacheFile, V6ip: req}

	if err := encodeTo(&buf, &cache); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal()
	}
}

func TestDecodeFrom(t *testing.T) {
	var buf bytes.Buffer
	cacheFile := "/baz/zot/pubip.cache"
	req := "2003:c8:9bec:4a70:e755:6628:78d1:aa06"
	encodeTo(&buf, &Cache{cacheFile: cacheFile, V6ip: req})

	cache := &Cache{cacheFile: cacheFile}
	if decodeFrom(&buf, cache); cache.V6ip != req {
		t.Errorf("should be %s, but is %s\n", req, cache.V6ip)
	}
}

func TestSave(t *testing.T) {
	tmp, err := ioutil.TempFile("", "pubip_cache")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	now := time.Now()
	cache := Cache{
		V6ip:      "ip-6",
		V6last:    now,
		V4ip:      "ip-4",
		V4last:    now,
		cacheFile: tmp.Name(),
	}
	if err := cache.save(); err != nil {
		t.Fatal(err)
	}
}

func TestMaybeIP_v6_tooOld(t *testing.T) {
	cache := Cache{
		V6ip:   "ip-6",
		V6last: time.Now().Add(-time.Duration(16 * time.Minute)),
	}
	if _, err := cache.maybeIP(pubip.IPv6); err == nil {
		t.Fatal()
	}
}

func TestMaybeIP_v6_ok(t *testing.T) {
	cache := Cache{
		V6ip:   "ip-6",
		V6last: time.Now().Add(-time.Duration(14 * time.Minute)),
	}
	ip, err := cache.maybeIP(pubip.IPv6)
	if err != nil {
		t.Fatal(err)
	}
	if ip != cache.V6ip {
		t.Fatal()
	}
}

func TestMaybeIP_v6_empty(t *testing.T) {
	cache := Cache{
		V6last: time.Now().Add(-time.Duration(14 * time.Minute)),
	}
	if _, err := cache.maybeIP(pubip.IPv6); err == nil {
		t.Fatal()
	}
}

func TestMaybeIP_v4_ok(t *testing.T) {
	cache := Cache{
		V4ip:   "ip-4",
		V4last: time.Now().Add(-time.Duration(14 * time.Minute)),
	}
	ip, err := cache.maybeIP(pubip.IPv4)
	if err != nil {
		t.Fatal(err)
	}
	if ip != cache.V4ip {
		t.Fail()
	}
}

func TestMaybeIP_v4_empty(t *testing.T) {
	cache := Cache{
		V4last: time.Now().Add(-time.Duration(14 * time.Minute)),
	}
	if _, err := cache.maybeIP(pubip.IPv4); err == nil {
		t.Fail()
	}
}

func TestMaybeIP_v6orv4_returnIPv6(t *testing.T) {
	cache := Cache{
		V6ip:   "ip-6",
		V6last: time.Now().Add(-time.Duration(14 * time.Minute)),
	}
	ip, err := cache.maybeIP(pubip.IPv6orIPv4)
	if err != nil {
		t.Fatal(err)
	}
	if ip != cache.V6ip {
		t.Fail()
	}
}

func TestMaybeIP_v6orv4_returnIPv4(t *testing.T) {
	cache := Cache{
		V6last: time.Now().Add(-time.Duration(14 * time.Minute)),
		V4ip:   "ip-4",
		V4last: time.Now().Add(-time.Duration(14 * time.Minute)),
	}
	ip, err := cache.maybeIP(pubip.IPv6orIPv4)
	if err != nil {
		t.Fatal(err)
	}
	if ip != cache.V4ip {
		t.Fatal()
	}
}
