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
	lnode* OrderedMapNode
	rnode* OrderedMapNode 
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (n *OrderedMapNode) InsertImpl(key, value int) (node *OrderedMapNode, created bool) {
	if (n == nil) { return &OrderedMapNode { key: key, value: value }, true }
	if n.key == key {
		n.value = value
		return n, false
	}
	if key < n.key {
		n.lnode, created = n.lnode.InsertImpl(key, value)
	} else {
		n.rnode, created = n.rnode.InsertImpl(key, value)
	}
	return n, created

}

func (m *OrderedMap) Insert(key, value int) {
	var created bool
	m.head, created = m.head.InsertImpl(key, value)
	if created { m.size++ }
}

func (n *OrderedMapNode) Find(key int) *OrderedMapNode {
	if n == nil    { return nil }
	if key < n.key { return n.lnode.Find(key) }
	if key > n.key { return n.rnode.Find(key) }
	return n
}

func (m *OrderedMap) Erase(key int) {
    var deleted bool
    m.head, deleted = deleteNode(m.head, key)
    if deleted {
        m.size--
    }
}

// Вспомогательная функция для удаления узла
// Возвращает новый корень поддерева и флаг, был ли удален узел
func deleteNode(root *OrderedMapNode, key int) (*OrderedMapNode, bool) {
    if root == nil {
        return nil, false
    }

    var deleted bool

    // Ищем узел для удаления
    if key < root.key {
        root.lnode, deleted = deleteNode(root.lnode, key)
    } else if key > root.key {
        root.rnode, deleted = deleteNode(root.rnode, key)
    } else {
        // Узел найден, выполняем удаление
        deleted = true
        
        // Случай 1: Узел - лист или имеет только одного потомка
        if root.lnode == nil {
            return root.rnode, deleted
        } else if root.rnode == nil {
            return root.lnode, deleted
        }

        // Случай 2: Узел имеет двух потомков
        // Находим минимальный узел в правом поддереве
        minNode := findMin(root.rnode)
        // Копируем данные
        root.key = minNode.key
        root.value = minNode.value
        // Удаляем дубликат
        root.rnode, _ = deleteNode(root.rnode, minNode.key)
    }

    return root, deleted
}

// findMin находит узел с минимальным ключом в поддереве
func findMin(node *OrderedMapNode) *OrderedMapNode {
    for node.lnode != nil {
        node = node.lnode
    }
    return node
}

func (m *OrderedMap) Contains(key int) bool {
	node := m.head.Find(key)
	if node != nil { return true }
	return false
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
