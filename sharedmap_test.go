package sharedmap_test

import (
	"github.com/codemartial/sharedmap"
	"math/rand"
	"sync"
	"testing"
)

const MAPSIZE = 1000000

var sm = sharedmap.NewSharedMap()

func TestSharedMapAdd(t *testing.T) {
	if modified := sm.Add(1, 1); modified {
		t.Error("Adding a new value resulted in a modification")
	}
	if modified := sm.Add(1, 2); !modified {
		t.Error("Adding an existing value did not result in modification")
	}
	if size := sm.Size(); size != 1 {
		t.Error("Map Size expected to be 1. Actual size", size)
	}
}

func TestSharedMapGet(t *testing.T) {
	if v := sm.Get(1); v != 2 {
		t.Error("Previously set value was not fetched. Expected 2, got", v)
	}
}

func TestSharedMapDelete(t *testing.T) {
	if deleted := sm.Delete(1); !deleted {
		t.Error("Existing key not deleted")
	} else if deleted := sm.Delete(1); deleted {
		t.Error("Previously deleted key was deleted again")
	}
	if size := sm.Size(); size != 0 {
		t.Error("Map expected to be empty but actually has", size, "elements")
	}
}

func benchmarkSharedMapN(b *testing.B, concurrents int) {
	s := sharedmap.NewSharedMap()
	r2wRatio := 40 // 40:1 read:write ratio
	a2dRatio := 10 // 10:1 add:delete ratio
	done := make(chan bool)

	b.ResetTimer()
	for i := 0; i < concurrents; i++ {
		go func() {
			for j := 0; j < b.N; j++ {
				do_read := rand.Intn(r2wRatio) > 0
				key := sharedmap.KeyType(rand.Intn(MAPSIZE))
				if do_read {
					s.Get(key)
				} else if rand.Intn(a2dRatio) > 0 {
					val := sharedmap.ValueType(key)
					s.Add(key, val)
				} else {
					s.Delete(key)
				}
			}
			done <- true
		}()
	}
	for i := 0; i < concurrents; i++ {
		<-done
	}
}

func benchmarkMutexMapN(b *testing.B, concurrents int) {
	m := map[sharedmap.KeyType]sharedmap.ValueType{}
	r2wRatio := 4   // 4:1 read:write ratio
	a2dRatio := 100 // 100:1 add:delete ratio
	done := make(chan bool)
	mu := &sync.Mutex{}

	b.ResetTimer()
	for i := 0; i < concurrents; i++ {
		go func() {
			for j := 0; j < b.N; j++ {
				do_read := rand.Intn(r2wRatio) > 0
				key := sharedmap.KeyType(rand.Intn(MAPSIZE))
				mu.Lock()
				if do_read {
					_ = m[key]
				} else if rand.Intn(a2dRatio) > 0 {
					_, modified := m[key]
					val := sharedmap.ValueType(key)
					m[key] = val
					_ = modified
				} else {
					_, found := m[key]
					delete(m, key)
					_ = found
				}
				mu.Unlock()
			}
			done <- true
		}()
	}
	for i := 0; i < concurrents; i++ {
		<-done
	}
}

func BenchmarkSharedMapSingle(b *testing.B) {
	benchmarkSharedMapN(b, 1)
}

func BenchmarkSharedMap10(b *testing.B) {
	benchmarkSharedMapN(b, 10)
}

func BenchmarkSharedMap100(b *testing.B) {
	benchmarkSharedMapN(b, 100)
}

func BenchmarkSharedMap10000(b *testing.B) {
	benchmarkSharedMapN(b, 10000)
}

func BenchmarkSharedMap100000(b *testing.B) {
	benchmarkSharedMapN(b, 100000)
}

func BenchmarkMutexMapSingle(b *testing.B) {
	benchmarkMutexMapN(b, 1)
}

func BenchmarkMutexMap10(b *testing.B) {
	benchmarkMutexMapN(b, 10)
}

func BenchmarkMutexMap100(b *testing.B) {
	benchmarkMutexMapN(b, 100)
}

func BenchmarkMutexMap10000(b *testing.B) {
	benchmarkMutexMapN(b, 10000)
}

func BenchmarkMutexMap100000(b *testing.B) {
	benchmarkMutexMapN(b, 100000)
}
