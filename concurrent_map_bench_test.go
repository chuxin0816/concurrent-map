package cmap

import (
	"strconv"
	"sync"
	"testing"
)

func BenchmarkItems(b *testing.B) {
	m := New[Animal]()

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		m.Items()
	}
}

func BenchmarkItemsInt(b *testing.B) {
	m := New[Animal]()

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		m.Items()
	}
}

func BenchmarkMarshalJson(b *testing.B) {
	m := New[Animal]()

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		_, err := m.MarshalJSON()
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Itoa(i)
	}
}

func BenchmarkSingleInsertAbsent(b *testing.B) {
	m := New[string]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), "value")
	}
}

func BenchmarkSingleInsertAbsentSyncMap(b *testing.B) {
	var m sync.Map
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store(strconv.Itoa(i), "value")
	}
}

func BenchmarkSingleInsertPresent(b *testing.B) {
	m := New[string]()
	m.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("key", "value")
	}
}

func BenchmarkSingleInsertPresentSyncMap(b *testing.B) {
	var m sync.Map
	m.Store("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store("key", "value")
	}
}

func BenchmarkMultiInsertDifferent(b *testing.B) {
	m := New[string]()
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferentSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, b.N)
	_, set := GetSetSyncMap[string](&m, finished)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_1_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](1))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_2_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](2))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_4_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](4))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_8_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](8))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_16_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](16))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_32_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](32))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_64_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](64))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_128_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](128))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_256_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](256))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_512_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](512))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_1024_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](1024))
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertSame(b *testing.B) {
	m := New[string]()
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	m.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertSameSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, b.N)
	_, set := GetSetSyncMap[string](&m, finished)
	m.Store("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSame(b *testing.B) {
	m := New[string]()
	finished := make(chan struct{}, b.N)
	get, _ := GetSet(m, finished)
	m.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go get("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSameSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, b.N)
	get, _ := GetSetSyncMap[string](&m, finished)
	m.Store("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go get("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent(b *testing.B) {
	m := New[string]()
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferentSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSetSyncMap[string](&m, finished)
	m.Store("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_1_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](1))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_16_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](16))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_32_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](32))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_64_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](64))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_128_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](128))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_256_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](256))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_512_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](512))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_1024_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](1024))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlock(b *testing.B) {
	m := New[string]()
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i%100), "value")
		go get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlockSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSetSyncMap[string](&m, finished)
	for i := 0; i < b.N; i++ {
		m.Store(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i%100), "value")
		go get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlock_1_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](1))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i%100), "value")
		go get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlock_16_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](16))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i%100), "value")
		go get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlock_32_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](32))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i%100), "value")
		go get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlock_256_Shard(b *testing.B) {
	m := New[string](WithShardCount[string](256))
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i%100), "value")
		go get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func GetSet[V any](m *ConcurrentMap[V], finished chan struct{}) (set func(key string, value V), get func(key string, value V)) {
	return func(key string, value V) {
			for i := 0; i < 10; i++ {
				m.Get(key)
			}
			finished <- struct{}{}
		}, func(key string, value V) {
			for i := 0; i < 10; i++ {
				m.Set(key, value)
			}
			finished <- struct{}{}
		}
}

func GetSetSyncMap[V any](m *sync.Map, finished chan struct{}) (get func(key string, value V), set func(key string, value V)) {
	get = func(key string, _ V) {
		for i := 0; i < 10; i++ {
			m.Load(key)
		}
		finished <- struct{}{}
	}
	set = func(key string, value V) {
		for i := 0; i < 10; i++ {
			m.Store(key, value)
		}
		finished <- struct{}{}
	}
	return
}

func BenchmarkKeys(b *testing.B) {
	m := New[Animal]()

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		m.Keys()
	}
}
