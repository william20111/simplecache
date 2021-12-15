package simplecache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestCache_Get(t *testing.T) {
	cache := newLRU(100)
	cache.set("test", "test", 100)
	type fields struct {
		cache    cacher
		maxItems int64
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantRes bool
	}{
		{
			name: "test get hit",
			fields: fields{
				cache: &lru{
					size:       100,
					mu:         sync.Mutex{},
					items:      cache.items,
					cacheOrder: nil,
				},
				maxItems: 100,
			},
			args:    args{key: "test"},
			want:    "test",
			wantRes: true,
		},
		{
			name: "test get miss",
			fields: fields{
				cache: &lru{
					size:       100,
					mu:         sync.Mutex{},
					items:      cache.items,
					cacheOrder: nil,
				},
				maxItems: 100,
			},
			args:    args{key: "doesntexist"},
			want:    nil,
			wantRes: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				cache: tt.fields.cache,
			}
			got, got1 := c.Get(tt.args.key)
			assert.Equalf(t, tt.want, got, "Get(%v)", tt.args.key)
			assert.Equalf(t, tt.wantRes, got1, "Get(%v)", tt.args.key)
		})
	}
}

func TestCache_Purge(t *testing.T) {
	cache := newLRU(100)
	cache.set("test", "test", 100)
	type fields struct {
		cache    cacher
		maxItems int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test get hit",
			fields: fields{
				cache: &lru{
					size:       100,
					mu:         sync.Mutex{},
					items:      cache.items,
					cacheOrder: cache.cacheOrder,
				},
				maxItems: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				cache: tt.fields.cache,
			}
			c.Purge()
			assert.Equal(t, 0, c.Len())
		})
	}
}

func TestCache_Set(t *testing.T) {
	cache := newLRU(100)
	type fields struct {
		cache    cacher
		maxItems int64
	}
	type args struct {
		key    string
		value  interface{}
		expiry time.Duration
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     bool
		populate int
	}{
		{
			name: "test set",
			fields: fields{
				cache: &lru{
					size:       100,
					mu:         sync.Mutex{},
					items:      cache.items,
					cacheOrder: cache.cacheOrder,
				},
				maxItems: 100,
			},
			want:     false,
			populate: 1,
		},
		{
			name: "test set max",
			fields: fields{
				cache: &lru{
					size:       100,
					mu:         sync.Mutex{},
					items:      cache.items,
					cacheOrder: cache.cacheOrder,
				},
				maxItems: 100,
			},
			want:     true,
			populate: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				cache: tt.fields.cache,
			}
			for i := 0; i < tt.populate; i++ {
				c.Set(fmt.Sprintf("test%d", i), fmt.Sprintf("test%d", i), 100)
			}
			assert.Equal(t, tt.populate, c.Len())
			assert.Equal(t, tt.want, c.Set(fmt.Sprintf("test%d", tt.populate+1), fmt.Sprintf("test%d", tt.populate+1), 100))
		})
	}
}

func TestNewCache(t *testing.T) {
	type args struct {
		maxItems int64
	}
	tests := []struct {
		name string
		args args
		want *Cache
	}{
		{
			name: "new test",
			args: args{maxItems: 100},
			want: &Cache{cache: newLRU(100)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, New(tt.args.maxItems), "New(%v)", tt.args.maxItems)
		})
	}
}

type BenchMarkObj struct {
	Name  string
	Blah  string
	Blah2 string
}

func BenchmarkCacheSet(b *testing.B) {
	cache := New(1000)
	for i := 0; i < b.N; i++ {
		cache.Set(
			fmt.Sprintf("test%d", i),
			BenchMarkObj{
				Name:  fmt.Sprintf("test%d", i),
				Blah:  fmt.Sprintf("test%d", i),
				Blah2: fmt.Sprintf("test%d", i),
			},
			100,
		)
	}
}

func Benchmark_CacheSynthetic(b *testing.B) {
	cache := New(1000)
	for i := 0; i < b.N; i++ {
		cache.Set(
			fmt.Sprintf("test%d", i),
			BenchMarkObj{
				Name:  fmt.Sprintf("test%d", i),
				Blah:  fmt.Sprintf("test%d", i),
				Blah2: fmt.Sprintf("test%d", i),
			},
			100,
		)
	}
}
