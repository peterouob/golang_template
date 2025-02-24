package grpcpool

import (
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"sync"
)

//REF:https://zhuanlan.zhihu.com/p/486094439

type Host struct {
	Name         string
	HLoadBalance int64
}

type Consistent struct {
	replicaNum         int
	totalLoad          int64
	hashFunc           func(key string) uint64
	hostMap            map[string]*Host
	replicaHostsMap    map[uint64]string
	sortedHostsHashSet []uint64
	sync.RWMutex
}

var (
	defaultReplicaNum = 10
	loadBoundFactor   = 0.25
	defaultHashFunc   = func(key string) uint64 {
		out := sha512.Sum512([]byte(key))
		return binary.LittleEndian.Uint64(out[:])
	}
	ErrHostAlreadyExists = errors.New("host already exists")
	ErrHostNotFound      = errors.New("host not found")
)

func (c *Consistent) RegisterHost(hostName string) error {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.hostMap[hostName]; ok {
		return ErrHostAlreadyExists
	}

	c.hostMap[hostName] = &Host{
		Name:         hostName,
		HLoadBalance: 0,
	}

	for i := 0; i < defaultReplicaNum; i++ {
		hashIndex := c.hashFunc(fmt.Sprintf("%s%d", hostName, i))
		c.replicaHostsMap[hashIndex] = hostName
		c.sortedHostsHashSet = append(c.sortedHostsHashSet, hashIndex)
	}

	sort.Slice(c.sortedHostsHashSet, func(i, j int) bool {
		if c.sortedHostsHashSet[i] < c.sortedHostsHashSet[j] {
			return true
		}
		return false
	})

	return nil
}

func (c *Consistent) RemoveHost(hostName string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.hostMap[hostName]; !ok {
		return ErrHostNotFound
	}
	delete(c.hostMap, hostName)
	for i := 0; i < defaultReplicaNum; i++ {
		hashIndex := c.hashFunc(fmt.Sprintf("%s%d", hostName, i))
		delete(c.replicaHostsMap, hashIndex)
		c.removeNode(hashIndex)
	}
	return nil
}

func (c *Consistent) removeNode(v uint64) {
	idx := -1
	l := 0
	r := len(c.sortedHostsHashSet) - 1
	for l <= r {
		m := (l + r) >> 1
		if c.sortedHostsHashSet[m] == v {
			idx = m
			break
		} else if c.sortedHostsHashSet[m] < v {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	if idx != -1 {
		c.sortedHostsHashSet = append(c.sortedHostsHashSet[:idx], c.sortedHostsHashSet[idx+1:]...)
	}
}

func (c *Consistent) GetHost(hostName string) (string, error) {
	hashKey := c.hashFunc(hostName)
	idx := c.getHost(hashKey)
	return c.replicaHostsMap[c.sortedHostsHashSet[idx]], nil
}

func (c *Consistent) getHost(key uint64) int {
	idx := sort.Search(len(c.sortedHostsHashSet), func(i int) bool {
		return c.sortedHostsHashSet[i] >= key
	})
	if idx >= len(c.sortedHostsHashSet) {
		idx = 0
	}
	return idx
}
