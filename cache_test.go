package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/thibran/pubip"
)

func TestCacheLocation(t *testing.T) {
	dir := cacheLocation()
	fmt.Println(dir)
	if len(dir) == 0 {
		t.Fail()
	}
}

func TestLoadCache_fail(t *testing.T) {
	cacheFile := "/foo"
	if _, err := loadCache(cacheFile); err == nil {
		t.Fail()
	}
}

func TestSave(t *testing.T) {
	cacheDir := "/home/tux/.cache/pubip.cache"
	cache := Cache{
		V6ip:   "ip-6",
		V6last: time.Now(),
		V4ip:   "ip-4",
		V4last: time.Now(),
	}
	if err := cache.save(cacheDir); err != nil {
		t.Error(err)
	}
}

func TestMaybeIP_v6_tooOld(t *testing.T) {
	cache := Cache{
		V6ip:   "ip-6",
		V6last: time.Now().Add(-time.Duration(16 * time.Minute)),
	}
	if _, err := cache.maybeIP(pubip.IPv6); err == nil {
		t.Fail()
	}
}

func TestMaybeIP_v6_ok(t *testing.T) {
	cache := Cache{
		V6ip:   "ip-6",
		V6last: time.Now().Add(-time.Duration(14 * time.Minute)),
	}
	ip, err := cache.maybeIP(pubip.IPv6)
	if err != nil {
		t.Error(err)
	}
	if ip != cache.V6ip {
		t.Fail()
	}
}

func TestMaybeIP_v6_empty(t *testing.T) {
	cache := Cache{
		V6last: time.Now().Add(-time.Duration(14 * time.Minute)),
	}
	if _, err := cache.maybeIP(pubip.IPv6); err == nil {
		t.Fail()
	}
}

func TestMaybeIP_v4_ok(t *testing.T) {
	cache := Cache{
		V4ip:   "ip-4",
		V4last: time.Now().Add(-time.Duration(14 * time.Minute)),
	}
	ip, err := cache.maybeIP(pubip.IPv4)
	if err != nil {
		t.Error(err)
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
		t.Error(err)
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
		t.Error(err)
	}
	if ip != cache.V4ip {
		t.Fail()
	}
}
