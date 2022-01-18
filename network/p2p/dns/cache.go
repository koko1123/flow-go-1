package dns

import (
	"net"
	"sync"
	"time"

	"github.com/onflow/flow-go/module/mempool/model"
)

// DefaultTimeToLive is the default duration a dns result is cached.
const (
	// DefaultTimeToLive is the default duration a dns result is cached.
	DefaultTimeToLive = 5 * time.Minute
	cacheEntryExists  = true
	cacheEntryFresh   = true // TTL yet has not reached
)

// cache is a ttl-based cache for dns entries
type cache struct {
	sync.RWMutex
	ttl      time.Duration // time-to-live for cache entry
	ipCache  map[string]*model.IpCacheEntryDNS
	txtCache map[string]*model.TxtCacheEntryDNS
}

func newCache() *cache {
	return &cache{
		ttl:      DefaultTimeToLive,
		ipCache:  make(map[string]*model.IpCacheEntryDNS),
		txtCache: make(map[string]*model.TxtCacheEntryDNS),
	}
}

// resolveIPCache resolves the domain through the cache if it is available.
// First boolean variable determines whether the domain exists in the cache.
// Second boolean variable determines whether the domain cache is fresh, i.e., TTL has not yet reached.
func (c *cache) resolveIPCache(domain string) ([]net.IPAddr, bool, bool) {
	c.RLock()

	entry, ok := c.ipCache[domain]

	c.RUnlock()

	if !ok {
		// does not exist
		return nil, !cacheEntryExists, !cacheEntryFresh
	}

	if time.Duration(runtimeNano()-entry.Timestamp) > c.ttl {
		// exists but expired
		return entry.Addresses, cacheEntryExists, !cacheEntryFresh
	}

	// exists and fresh
	return entry.Addresses, cacheEntryExists, cacheEntryFresh
}

// resolveIPCache resolves the txt through the cache if it is available.
// First boolean variable determines whether the txt exists in the cache.
// Second boolean variable determines whether the txt cache is fresh, i.e., TTL has not yet reached.
func (c *cache) resolveTXTCache(txt string) ([]string, bool, bool) {
	c.RLock()

	entry, ok := c.txtCache[txt]

	c.RUnlock()

	if !ok {
		// does not exist
		return nil, !cacheEntryExists, !cacheEntryFresh
	}

	if time.Duration(runtimeNano()-entry.Timestamp) > c.ttl {
		// exists but expired
		return entry.Addresses, cacheEntryExists, !cacheEntryFresh
	}

	// exists and fresh
	return entry.Addresses, cacheEntryExists, cacheEntryFresh
}

// updateIPCache updates the cache entry for the domain.
func (c *cache) updateIPCache(domain string, addr []net.IPAddr) {
	c.Lock()
	defer c.Unlock()

	c.ipCache[domain] = &model.IpCacheEntryDNS{
		Addresses: addr,
		Timestamp: runtimeNano(),
	}
}

// updateTXTCache updates the cache entry for the txt.
func (c *cache) updateTXTCache(txt string, addr []string) {
	c.Lock()
	defer c.Unlock()

	c.txtCache[txt] = &model.TxtCacheEntryDNS{
		Addresses: addr,
		Timestamp: runtimeNano(),
	}
}

// invalidateIPCacheEntry atomically invalidates ip cache entry. Boolean variable determines whether invalidation
// is successful.
func (c *cache) invalidateIPCacheEntry(domain string) bool {
	c.RLock()
	_, exists := c.ipCache[domain]
	c.RUnlock()

	if !exists {
		return false
	}

	c.Lock()
	delete(c.ipCache, domain)
	c.Unlock()

	return true
}

// invalidateTXTCacheEntry atomically invalidates txt cache entry. Boolean variable determines whether invalidation
// is successful.
func (c *cache) invalidateTXTCacheEntry(txt string) bool {
	c.RLock()
	_, exists := c.txtCache[txt]
	c.RUnlock()

	if !exists {
		return false
	}

	c.Lock()
	delete(c.txtCache, txt)
	c.Unlock()

	return exists
}
