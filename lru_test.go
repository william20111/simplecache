package simplecache

import (
	"container/list"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type fakeItems struct {
	key    string
	val    interface{}
	expiry time.Duration
}

func Test_lru_get(t *testing.T) {
	type args struct {
		key  string
		size int
	}
	tests := []struct {
		name        string
		args        args
		populate    []fakeItems
		wantElement interface{}
		wantBool    bool
	}{
		{
			name: "test get miss",
			args: args{
				key:  "test",
				size: 100,
			},
			populate: []fakeItems{
				{
					key:    "test20",
					val:    "test20",
					expiry: 1,
				},
			},
			wantElement: nil,
			wantBool:    false,
		},
		{
			name: "test get hit",
			args: args{
				key:  "test20",
				size: 100,
			},
			populate: []fakeItems{
				{
					key:    "test20",
					val:    "test20",
					expiry: 1,
				},
			},
			wantElement: "test20",
			wantBool:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLRU(tt.args.size)
			for i := range tt.populate {
				_ = l.set(tt.populate[i].key, tt.populate[i].val, tt.populate[i].expiry)
			}
			got, ok := l.get(tt.args.key)
			assert.Equal(t, tt.wantBool, ok)
			assert.Equal(t, tt.wantElement, got)
		})
	}
}

func Test_newLRU(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want *lru
	}{
		{
			name: "new lru",
			args: args{size: 100},
			want: &lru{
				size:       100,
				mu:         sync.Mutex{},
				items:      map[string]*list.Element{},
				cacheOrder: list.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, newLRU(tt.args.size), "newLRU(%v)", tt.args.size)
		})
	}
}

func Test_lru_len(t *testing.T) {
	type fields struct {
		size       int
		items      map[string]*list.Element
		cacheOrder *list.List
	}
	tests := []struct {
		name     string
		fields   fields
		populate []fakeItems
		want     int
	}{
		{
			name: "test len 1",
			fields: fields{
				size:       100,
				items:      map[string]*list.Element{},
				cacheOrder: list.New(),
			},
			want: 1,
			populate: []fakeItems{
				{
					val:    "test",
					expiry: 0,
				},
			},
		},
		{
			name: "test len 10",
			fields: fields{
				size:       100,
				items:      map[string]*list.Element{},
				cacheOrder: list.New(),
			},
			want: 8,
			populate: []fakeItems{
				{
					key:    "test",
					val:    "test",
					expiry: 0,
				},
				{
					key:    "test1",
					val:    "test1",
					expiry: 0,
				},
				{
					key:    "test2",
					val:    "test2",
					expiry: 0,
				},
				{
					key:    "test3",
					val:    "test3",
					expiry: 0,
				},
				{
					key:    "test4",
					val:    "test4",
					expiry: 0,
				},
				{
					key:    "test5",
					val:    "test5",
					expiry: 0,
				},
				{
					key:    "test6",
					val:    "test6",
					expiry: 0,
				},
				{
					key:    "test7",
					val:    "test7",
					expiry: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lru{
				size:       tt.fields.size,
				mu:         sync.Mutex{},
				items:      tt.fields.items,
				cacheOrder: tt.fields.cacheOrder,
			}
			for i := range tt.populate {
				l.set(tt.populate[i].key, &tt.populate[i].val, tt.populate[i].expiry)
			}
			assert.Equalf(t, tt.want, l.len(), "len()")
		})
	}
}

func Test_lru_purge(t *testing.T) {
	type fields struct {
		size       int
		items      map[string]*list.Element
		cacheOrder *list.List
	}
	tests := []struct {
		name     string
		fields   fields
		want     int
		populate []fakeItems
	}{
		{
			name: "test purge",
			fields: fields{
				size:       100,
				items:      map[string]*list.Element{},
				cacheOrder: list.New(),
			},
			populate: []fakeItems{
				{
					key:    "test",
					val:    "test",
					expiry: 0,
				},
				{
					key:    "test1",
					val:    "test1",
					expiry: 0,
				},
				{
					key:    "test2",
					val:    "test2",
					expiry: 0,
				},
				{
					key:    "test3",
					val:    "test3",
					expiry: 0,
				},
				{
					key:    "test4",
					val:    "test4",
					expiry: 0,
				},
				{
					key:    "test5",
					val:    "test5",
					expiry: 0,
				},
				{
					key:    "test6",
					val:    "test6",
					expiry: 0,
				},
				{
					key:    "test7",
					val:    "test7",
					expiry: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lru{
				size:       tt.fields.size,
				mu:         sync.Mutex{},
				items:      tt.fields.items,
				cacheOrder: tt.fields.cacheOrder,
			}
			for i := range tt.populate {
				l.set(tt.populate[i].key, &tt.populate[i].val, tt.populate[i].expiry)
			}
			res := l.purge()
			assert.Equal(t, true, res)
			assert.Equal(t, 0, l.cacheOrder.Len())
			assert.Equal(t, 0, len(l.items))
		})
	}
}

func Test_lru_remove(t *testing.T) {
	type fields struct {
		size       int
		items      map[string]*list.Element
		cacheOrder *list.List
	}
	type args struct {
		key string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     bool
		wantLen  int
		populate []fakeItems
		remove   []string
	}{
		{
			name: "test len 1",
			fields: fields{
				size:       100,
				items:      map[string]*list.Element{},
				cacheOrder: list.New(),
			},
			populate: []fakeItems{
				{
					val:    "test",
					expiry: 0,
				},
			},
			remove: []string{
				"test",
			},
			want:    true,
			wantLen: 0,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lru{
				size:       tt.fields.size,
				mu:         sync.Mutex{},
				items:      tt.fields.items,
				cacheOrder: tt.fields.cacheOrder,
			}
			for i := range tt.populate {
				l.set(tt.populate[i].key, &tt.populate[i].val, tt.populate[i].expiry)
			}
			for i := range tt.remove {
				l.remove(tt.remove[i])
			}
			assert.Equalf(t, tt.want, l.remove(tt.args.key), "remove(%v)", tt.args.key)
			assert.Equal(t, tt.wantLen, len(l.items))
			assert.Equal(t, tt.wantLen, l.cacheOrder.Len())
		})
	}
}

func Test_lru_set_eviction(t *testing.T) {
	type fields struct {
		size       int
		items      map[string]*list.Element
		cacheOrder *list.List
	}
	type args struct {
		key    string
		val    interface{}
		expiry time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantLen int
	}{
		{
			name: "test len 1",
			fields: fields{
				size:       100,
				items:      map[string]*list.Element{},
				cacheOrder: list.New(),
			},
			want:    true,
			wantLen: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lru{
				size:       tt.fields.size,
				mu:         sync.Mutex{},
				items:      tt.fields.items,
				cacheOrder: tt.fields.cacheOrder,
			}
			for i := 0; i < 100; i++ {
				l.set(fmt.Sprintf("test%d", i), fmt.Sprintf("test%d", i), 0)
			}
			res := l.set("test101", "test101", 0)
			assert.Equal(t, tt.want, res)
			assert.Equal(t, tt.wantLen, l.len())
		})
	}
}
