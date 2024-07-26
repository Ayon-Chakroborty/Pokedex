package cache

import (
	"sync"
	"time"
)

type locationCacheEntry struct {
	createdAt    time.Time
	locationData []byte
}

type pokemonCacheEntry struct {
	createdAt   time.Time
	pokemonData []byte
}

type PokeCache struct {
	locationCache map[string]locationCacheEntry
	pokemonCache  map[string]pokemonCacheEntry
	mux           *sync.Mutex
}

func NewCache(interval time.Duration) *PokeCache {
	cache := PokeCache{
		locationCache: make(map[string]locationCacheEntry),
		pokemonCache: make(map[string]pokemonCacheEntry),
		mux:           &sync.Mutex{},
	}

	go cache.reapLoop(interval)

	return &cache
}

func (p *PokeCache) AddLocation(key string, data []byte) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.locationCache[key] = locationCacheEntry{
		createdAt:    time.Now(),
		locationData: data,
	}
}

func (p *PokeCache) AddPokemon(key string, data []byte){
	p.mux.Lock()
	defer p.mux.Unlock()
	p.pokemonCache[key] = pokemonCacheEntry{
		createdAt: time.Now(),
		pokemonData: data,
	}
}

func (p *PokeCache) GetLocation(key string) (val []byte, exists bool) {
	p.mux.Lock()
	defer p.mux.Unlock()
	data, ok := p.locationCache[key]
	return data.locationData, ok
}

func (p *PokeCache) GetPokemon(key string) (val []byte, exists bool){
	p.mux.Lock()
	defer p.mux.Unlock()
	data, ok := p.pokemonCache[key]
	return data.pokemonData, ok
}

func (p *PokeCache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		p.reap(interval)
	}
}

func (p *PokeCache) reap(interaval time.Duration) {
	p.mux.Lock()
	defer p.mux.Unlock()

	for key, val := range p.locationCache {
		if val.createdAt.Before(time.Now().Add(-interaval)) {
			delete(p.locationCache, key)
		}
	}

	for key, val := range p.pokemonCache{
		if(val.createdAt.Before(time.Now().Add(-interaval))){
			delete(p.pokemonCache, key)
		}
	}

	for key, val := range p.pokemonCache {
		if val.createdAt.Before(time.Now().Add(-interaval)) {
			delete(p.pokemonCache, key)
		}
	}
}
