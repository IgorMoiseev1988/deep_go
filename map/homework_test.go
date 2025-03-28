package main

import (
	"reflect"
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type OrderedMap struct {
	head* OrderedMapNode
	size int
}

type OrderedMapNode struct {

	key int
	value int
	parent* OrderedMapNode
	lnode* OrderedMapNode
	rnode* OrderedMapNode 
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func NewOrderedMapNode(key, value int) *OrderedMapNode {
	return &OrderedMapNode { key: key, value: value }
}

func (m *OrderedMap) InsertImpl(n *OrderedMapNode, key int, value int) {
	if m.head == nil {
		m.head = NewOrderedMapNode(key, value)
		m.size++
		return
	}
	if key < n.key {
		if n.lnode == nil {
			n.lnode = NewOrderedMapNode(key, value)
			m.size++
		} else {
			m.InsertImpl(n.lnode, key, value)
		}
	} else if key > n.key {
		if n.rnode == nil {
			n.rnode = NewOrderedMapNode(key, value)
			m.size++
		} else {
			m.InsertImpl(n.rnode, key, value)
		} 
	} else {
		n.value = value
	}
}

func (m *OrderedMap) Insert(key, value int) {
	m.InsertImpl(m.head, key, value)
}

func (m *OrderedMap) Erase(key int) {
	// need to implement
}

func (n *OrderedMapNode) ContainsImpl(key int) bool {
	if (n == nil) {
		return false
	}
	if n.key == key {
		return true
	}
	if key < n.key {
		return n.lnode.ContainsImpl(key)
	} else {
		return n.rnode.ContainsImpl(key)
	}
}

func (m *OrderedMap) Contains(key int) bool {
	return m.head.ContainsImpl(key)	
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (n *OrderedMapNode) ForEachImpl(action func(int, int)) {
	if (n == nil) {
		return
	}
	n.lnode.ForEachImpl(action)
	action(n.key, n.value)
	n.rnode.ForEachImpl(action)
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.head.ForEachImpl(action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	fmt.Printf("Map: head %v, size %d\n", data.head, data.size)
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	fmt.Printf("%v\n", data)
	data.Insert(10, 15)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})
	fmt.Printf("%v == %v\n", keys, expectedKeys)
	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
