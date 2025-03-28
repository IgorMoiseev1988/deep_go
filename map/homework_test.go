package main

import (
	"reflect"
	"testing"
//	"fmt"

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

func (n *OrderedMapNode) Search(key int) (parent, node *OrderedMapNode) {
	if n == nil     { return nil, nil }
	if n.key == key { return nil, n	  }

	if n.lnode != nil && n.lnode.key == key {
		return n, n.lnode
	}
	if n.rnode != nil && n.rnode.key == key {
		return n, n.rnode
	}
	if key < n.key {
		return n.lnode.Search(key)
	} else {
		return n.rnode.Search(key)
	}
	
}

func (m *OrderedMap) Erase(key int) {
	parent, node := m.head.Search(key)
	_ = parent
	if node == nil { return }

	if node.rnode == nil {
		if parent.lnode == node { 
			parent.lnode = node.lnode 
		} else { 
			parent.rnode = node.lnode 
		}
		node.lnode = nil
	} else {
		parent_tmp, node_tmp := node, node.rnode
		for node_tmp.lnode != nil {
			parent_tmp = node_tmp
			node_tmp = node_tmp.lnode
		}
		node.key = node_tmp.key
		if parent_tmp.lnode == node_tmp {
			parent_tmp.lnode = node_tmp.rnode
		} else {
			parent_tmp.rnode = node_tmp.rnode
		}
	}
		
	m.size--
}

func (n *OrderedMapNode) ContainsImpl(key int) bool {
	_, node := n.Search(key)
	if node != nil { return true } else { return false }
}

func (m *OrderedMap) Contains(key int) bool {
	return m.head.ContainsImpl(key)	
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (n *OrderedMapNode) ForEachImpl(action func(int, int)) {
	if n == nil { return }

	n.lnode.ForEachImpl(action)
	action(n.key, n.value)
	n.rnode.ForEachImpl(action)
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.head.ForEachImpl(action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
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
