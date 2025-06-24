package in_memory

import (
	"hash/fnv"
	"sync"
)

type HashMap struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewHashMap() *HashMap {
	return &HashMap{
		make(map[string]string),
		sync.RWMutex{},
	}
}

func (hm *HashMap) Set(k string, v string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.data[k] = v
}

func (hm *HashMap) Get(k string) (string, bool) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	val, ok := hm.data[k]
	return val, ok
}

func (hm *HashMap) Delete(k string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	delete(hm.data, k)
}

type SharedMap struct {
	data     []*HashMap
	shardCnt uint32
}

func NewSharedMap(shardNum uint32) *SharedMap {
	if shardNum == 0 {
		shardNum = 1
	}
	data := make([]*HashMap, shardNum)
	for i := range data {
		data[i] = NewHashMap()
	}
	return &SharedMap{data: data, shardCnt: shardNum}
}

func (sm *SharedMap) shardNumByKey(key string) uint32 {
	hash32 := fnv.New32()
	hash32.Write([]byte(key))
	return hash32.Sum32() % sm.shardCnt
}

func (sm *SharedMap) Set(k string, v string) {
	shardNum := sm.shardNumByKey(k)
	sm.data[shardNum].Set(k, v)
}

func (sm *SharedMap) Get(k string) (string, bool) {
	shardNum := sm.shardNumByKey(k)
	return sm.data[shardNum].Get(k)
}

func (sm *SharedMap) Delete(k string) {
	shardNum := sm.shardNumByKey(k)
	sm.data[shardNum].Delete(k)
}
