package cmap

import (
	"encoding/json"
	"sync"
)

const SHARD_COUNT = 128

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (SHARD_COUNT) map shards.
type ConcurrentMap[V any] struct {
	shardCount int
	shards     []*ConcurrentMapShared[V]
	sharding   func(key string) uint64
}

// A "thread" safe string to anything map.
type ConcurrentMapShared[V any] struct {
	items        map[string]V
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

type Option[V any] func(*ConcurrentMap[V])

func WithShardCount[V any](shardCount int) Option[V] {
	if shardCount <= 0 {
		panic("shardCount must be greater than 0")
	}
	if shardCount&(shardCount-1) != 0 {
		panic("shardCount must be a power of 2")
	}

	return func(cm *ConcurrentMap[V]) {
		cm.shardCount = shardCount
		cm.shards = make([]*ConcurrentMapShared[V], shardCount)
	}
}

func WithShardingFunction[V any](sharding func(key string) uint64) Option[V] {
	return func(cm *ConcurrentMap[V]) {
		cm.sharding = sharding
	}
}

// Creates a new concurrent map.
func New[V any](opts ...Option[V]) *ConcurrentMap[V] {
	m := &ConcurrentMap[V]{
		shardCount: SHARD_COUNT,
		sharding:   fnv64a,
		shards:     make([]*ConcurrentMapShared[V], SHARD_COUNT),
	}
	for _, opt := range opts {
		opt(m)
	}

	for i := 0; i < m.shardCount; i++ {
		m.shards[i] = &ConcurrentMapShared[V]{items: make(map[string]V)}
	}
	return m
}

// GetShard returns shard under given key
func (m ConcurrentMap[V]) GetShard(key string) *ConcurrentMapShared[V] {
	return m.shards[uint(m.sharding(key))%uint(m.shardCount)]
}

func (m ConcurrentMap[V]) MSet(data map[string]V) {
	for key, value := range data {
		shard := m.GetShard(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
}

// Sets the given value under the specified key.
func (m ConcurrentMap[V]) Set(key string, value V) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

// Callback to return new element to be inserted into the map
// It is called while lock is held, therefore it MUST NOT
// try to access other keys in same map, as it can lead to deadlock since
// Go sync.RWLock is not reentrant
type UpsertCb[V any] func(exist bool, valueInMap V, newValue V) V

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m ConcurrentMap[V]) Upsert(key string, value V, cb UpsertCb[V]) (res V) {
	shard := m.GetShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	res = cb(ok, v, value)
	shard.items[key] = res
	shard.Unlock()
	return res
}

// Sets the given value under the specified key if no value was associated with it.
func (m ConcurrentMap[V]) SetIfAbsent(key string, value V) bool {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.Unlock()
	return !ok
}

// Get retrieves an element from map under given key.
func (m ConcurrentMap[V]) Get(key string) (V, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Count returns the number of elements within the map.
func (m ConcurrentMap[V]) Count() int {
	count := 0
	for i := 0; i < m.shardCount; i++ {
		shard := m.shards[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Looks up an item under specified key
func (m ConcurrentMap[V]) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Remove removes an element from the map.
func (m ConcurrentMap[V]) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// RemoveCb is a callback executed in a map.RemoveCb() call, while Lock is held
// If returns true, the element will be removed from the map
type RemoveCb[V any] func(key string, v V, exists bool) bool

// RemoveCb locks the shard containing the key, retrieves its current value and calls the callback with those params
// If callback returns true and element exists, it will remove it from the map
// Returns the value returned by the callback (even if element was not present in the map)
func (m ConcurrentMap[V]) RemoveCb(key string, cb RemoveCb[V]) bool {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	remove := cb(key, v, ok)
	if remove && ok {
		delete(shard.items, key)
	}
	shard.Unlock()
	return remove
}

// Pop removes an element from the map and returns it
func (m ConcurrentMap[V]) Pop(key string) (v V, exists bool) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// IsEmpty checks if map is empty.
func (m ConcurrentMap[V]) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type Tuple[V any] struct {
	Key string
	Val V
}

// Iter returns an iterator which could be used in a for range loop.
//
// Deprecated: using IterBuffered() will get a better performence
func (m ConcurrentMap[V]) Iter() <-chan Tuple[V] {
	chans := snapshot(m)
	ch := make(chan Tuple[V])
	go fanIn(chans, ch)
	return ch
}

// IterBuffered returns a buffered iterator which could be used in a for range loop.
func (m ConcurrentMap[V]) IterBuffered() <-chan Tuple[V] {
	chans := snapshot(m)
	total := 0
	for _, c := range chans {
		total += cap(c)
	}
	ch := make(chan Tuple[V], total)
	go fanIn(chans, ch)
	return ch
}

func (m ConcurrentMap[V]) PopAll() <-chan Tuple[V] {
	chans := popAll(m)
	total := 0
	for _, c := range chans {
		total += cap(c)
	}
	ch := make(chan Tuple[V], total)
	go fanIn(chans, ch)
	return ch
}

// Clear removes all items from map.
func (m ConcurrentMap[V]) Clear() {
	for item := range m.IterBuffered() {
		m.Remove(item.Key)
	}
}

// Returns a array of channels that contains elements in each shard,
// which likely takes a snapshot of `m`.
// It returns once the size of each buffered channel is determined,
// before all the channels are populated using goroutines.
func snapshot[V any](m ConcurrentMap[V]) (chans []chan Tuple[V]) {
	// When you access map items before initializing.
	if len(m.shards) == 0 {
		panic(`cmap.ConcurrentMap is not initialized. Should run New() before usage.`)
	}
	chans = make([]chan Tuple[V], m.shardCount)
	wg := sync.WaitGroup{}
	wg.Add(m.shardCount)
	// Foreach shard.
	for index, shard := range m.shards {
		go func(index int, shard *ConcurrentMapShared[V]) {
			// Foreach key, value pair.
			shard.RLock()
			chans[index] = make(chan Tuple[V], len(shard.items))
			wg.Done()
			for key, val := range shard.items {
				chans[index] <- Tuple[V]{key, val}
			}
			shard.RUnlock()
			close(chans[index])
		}(index, shard)
	}
	wg.Wait()
	return chans
}

// Returns a array of channels that contains elements in each shard and clears the map.
func popAll[V any](m ConcurrentMap[V]) (chans []chan Tuple[V]) {
	// When you access map items before initializing.
	if len(m.shards) == 0 {
		panic(`cmap.ConcurrentMap is not initialized. Should run New() before usage.`)
	}
	chans = make([]chan Tuple[V], m.shardCount)
	wg := sync.WaitGroup{}
	wg.Add(m.shardCount)
	// Foreach shard.
	for index, shard := range m.shards {
		go func(index int, shard *ConcurrentMapShared[V]) {
			// Foreach key, value pair.
			shard.Lock()
			chans[index] = make(chan Tuple[V], len(shard.items))
			wg.Done()
			for key, val := range shard.items {
				chans[index] <- Tuple[V]{key, val}
			}
			close(chans[index])
			shard.items = make(map[string]V)
			shard.Unlock()
		}(index, shard)
	}
	wg.Wait()
	return chans
}

// fanIn reads elements from channels `chans` into channel `out`
func fanIn[V any](chans []chan Tuple[V], out chan Tuple[V]) {
	wg := sync.WaitGroup{}
	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch chan Tuple[V]) {
			for t := range ch {
				out <- t
			}
			wg.Done()
		}(ch)
	}
	wg.Wait()
	close(out)
}

// Items returns all items as map[string]V
func (m ConcurrentMap[V]) Items() map[string]V {
	tmp := make(map[string]V)

	// Insert items to temporary map.
	for item := range m.IterBuffered() {
		tmp[item.Key] = item.Val
	}

	return tmp
}

// Iterator callbacalled for every key,value found in
// maps. RLock is held for all calls for a given shard
// therefore callback sess consistent view of a shard,
// but not across the shards
type IterCb[V any] func(key string, v V)

// Callback based iterator, cheapest way to read
// all elements in a map.
func (m ConcurrentMap[V]) IterCb(fn IterCb[V]) {
	for idx := range m.shards {
		shard := (m.shards)[idx]
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// Keys returns all keys as []string
func (m ConcurrentMap[V]) Keys() []string {
	count := m.Count()
	ch := make(chan string, count)
	go func() {
		// Foreach shard.
		wg := sync.WaitGroup{}
		wg.Add(m.shardCount)
		for _, shard := range m.shards {
			go func(shard *ConcurrentMapShared[V]) {
				// Foreach key, value pair.
				shard.RLock()
				for key := range shard.items {
					ch <- key
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	// Generate keys
	keys := make([]string, 0, count)
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

// Reviles ConcurrentMap "private" variables to json marshal.
func (m ConcurrentMap[V]) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[string]V)

	// Insert items to temporary map.
	for item := range m.IterBuffered() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}

func fnv64a(key string) uint64 {
	var hash uint64 = 14695981039346656037
	const prime64 = 1099511628211
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}

	return hash
}

// Reverse process of Marshal.
func (m *ConcurrentMap[V]) UnmarshalJSON(b []byte) (err error) {
	tmp := make(map[string]V)

	// Unmarshal into a single map.
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	// foreach key,value pair in temporary map insert into our concurrent map.
	for key, val := range tmp {
		m.Set(key, val)
	}
	return nil
}
