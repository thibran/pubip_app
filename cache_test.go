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

func TestMaybeIP(t *testing.T) {
	timeOk := time.Now().Add(-time.Duration(14 * time.Minute))
	timeErr := time.Now().Add(-time.Duration(16 * time.Minute))
	tt := []struct {
		name, exp string
		expErr    error
		ip6, ip4  string
		iptype    pubip.IPType
		time      time.Time
	}{
		{name: "v6 ok", ip6: "ip-6", exp: "ip-6", expErr: nil,
			time: timeOk, iptype: pubip.IPv6},
		{name: "v6 too old", ip6: "ip-6", expErr: errIPv6Empty,
			time: timeErr, iptype: pubip.IPv6},
		{name: "v6 empty", expErr: errIPv6Empty,
			time: timeOk, iptype: pubip.IPv6},
		{name: "v4 ok", ip4: "ip-4", exp: "ip-4", expErr: nil,
			time: timeOk, iptype: pubip.IPv4},
		{name: "v4 too old", ip4: "ip-4", expErr: errIPv4Empty,
			time: timeErr, iptype: pubip.IPv4},
		{name: "v4 empty", expErr: errIPv4Empty,
			time: timeOk, iptype: pubip.IPv4},
		{name: "v6 or v4 - return v6", expErr: nil,
			ip6: "ip-6", exp: "ip-6",
			time: timeOk, iptype: pubip.IPv6orIPv4},
		{name: "v6 or v4 - return v4", expErr: nil,
			ip4: "ip-4", exp: "ip-4",
			time: timeOk, iptype: pubip.IPv6orIPv4},
		{name: "not in cache", expErr: errNotInCache,
			time: timeOk, iptype: pubip.IPv6orIPv4},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cache := Cache{
				V6ip:   tc.ip6,
				V4ip:   tc.ip4,
				V6last: tc.time,
				V4last: tc.time,
			}
			ip, err := cache.maybeIP(tc.iptype)
			if ip != tc.exp {
				t.Fatalf("exp ip: %q, got %q", tc.exp, ip)
			}
			if err != tc.expErr {
				t.Fatalf("exp err: %q, got %q", tc.expErr, err)
			}
			// if ip != "" && err != tc.expErr {
			// 	t.Fatalf("exp err: %v, got %v",
			// 		tc.expErr, err)
			// }
		})
	}
}
