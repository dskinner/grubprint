package client

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/facebookgo/httpcontrol/httpcache"
	"github.com/peterbourgon/diskv"
)

func MemoryCacheTransport(t http.RoundTripper) http.RoundTripper {
	return &httpcache.Transport{
		Config:    httpcache.CacheByPath(24 * time.Hour),
		ByteCache: &memoryCache{bin: make(map[string][]byte)},
		Transport: t,
	}
}

type memoryCache struct {
	sync.RWMutex
	bin map[string][]byte
}

func (c *memoryCache) Store(key string, value []byte, timeout time.Duration) error {
	c.Lock()
	defer c.Unlock()
	c.bin[key] = value
	return nil
}

func (c *memoryCache) Get(key string) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	if v, ok := c.bin[key]; ok {
		return v, nil
	}
	return nil, nil
}

func DiskCacheTransport(path string, t http.RoundTripper) http.RoundTripper {
	if err := os.Mkdir(path, 0700); err != nil && !os.IsExist(err) {
		log.Fatalf("Mkdir(%s) failed: %s", path, err)
	}
	return &httpcache.Transport{
		Config: httpcache.CacheByPath(24 * time.Hour),
		ByteCache: &diskCache{
			diskv.New(diskv.Options{
				BasePath:     path,
				CacheSizeMax: 100 * 1024 * 1024, // 100MB
			}),
		},
		Transport: t,
	}
}

func ktof(key string) (filename string) {
	h := md5.New()
	io.WriteString(h, key)
	return hex.EncodeToString(h.Sum(nil))
}

type diskCache struct {
	d *diskv.Diskv
}

func (c *diskCache) Store(key string, value []byte, timeout time.Duration) error {
	return c.d.WriteStream(ktof(key), bytes.NewReader(value), true)
}

func (c *diskCache) Get(key string) ([]byte, error) {
	return c.d.Read(ktof(key))
}
