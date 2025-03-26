package main

import (
	"reflect"
	"testing"
	"unsafe"
	"fmt"
	"runtime"

	"github.com/stretchr/testify/assert"
)

type COWBuffer struct {
	data []byte
	refs *int
}

func NewCOWBuffer(data []byte) COWBuffer {
	var refs int = 1
	newBuffer := COWBuffer{data, &refs}
	runtime.SetFinalizer(&newBuffer, (*COWBuffer).Close)
	return newBuffer
}

func (b *COWBuffer) Clone() COWBuffer {
	*b.refs++
	newBuffer := COWBuffer{b.data, b.refs}
	runtime.SetFinalizer(&newBuffer, (*COWBuffer).Close)
	return newBuffer
}

func (b *COWBuffer) Close() {
	if b.data != nil {
		b.data = nil
		*b.refs--
	}
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if index < 0 || index >= len(b.data) || *b.refs <= 0 {
		return false
	}
	if *b.refs > 1 {
		newData := make([]byte, len(b.data))
		copy(newData, b.data)
		*b.refs--
		b.data = newData
		var refs int = 1
		b.refs = &refs
	}
	b.data[index] = value
	
	return true
}

func (b *COWBuffer) String() string {
	if len(b.data) == 0 {
	return ""
	}

	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()
	
	fmt.Printf("====STATE1===\n")	
	fmt.Printf("%v, %d\n", buffer, *buffer.refs)
	fmt.Printf("%v, %d\n", copy1, *copy1.refs)
	fmt.Printf("%v, %d\n", copy2, *copy2.refs)


	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	fmt.Printf("====STATE2===\n")	
	fmt.Printf("%v, %d\n", buffer, *buffer.refs)
	fmt.Printf("%v, %d\n", copy1, *copy1.refs)
	fmt.Printf("%v, %d\n", copy2, *copy2.refs)
	assert.False(t, buffer.Update(-1, 'g'))
	fmt.Printf("====STATE3===\n")	
	fmt.Printf("%v, %d\n", buffer, *buffer.refs)
	fmt.Printf("%v, %d\n", copy1, *copy1.refs)
	fmt.Printf("%v, %d\n", copy2, *copy2.refs)
	assert.False(t, buffer.Update(4, 'g'))
	fmt.Printf("====STATE4===\n")	
	fmt.Printf("%v, %d\n", buffer, *buffer.refs)
	fmt.Printf("%v, %d\n", copy1, *copy1.refs)
	fmt.Printf("%v, %d\n", copy2, *copy2.refs)

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()
	fmt.Printf("====STATE5===\n")	
	fmt.Printf("%v, %d\n", buffer, *buffer.refs)
	fmt.Printf("%v, %d\n", copy1, *copy1.refs)
	fmt.Printf("%v, %d\n", copy2, *copy2.refs)
	copy1.Close()
	fmt.Printf("====STATE6===\n")	
	fmt.Printf("%v, %d\n", buffer, *buffer.refs)
	fmt.Printf("%v, %d\n", copy1, *copy1.refs)
	fmt.Printf("%v, %d\n", copy2, *copy2.refs)

	previous := copy2.data
	copy2.Update(0, 'f')
	fmt.Printf("====STATE7===\n")	
	fmt.Printf("%v, %d\n", buffer, *buffer.refs)
	fmt.Printf("%v, %d\n", copy1, *copy1.refs)
	fmt.Printf("%v, %d\n", copy2, *copy2.refs)
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
	fmt.Printf("====STATE8===\n")	
	fmt.Printf("%v, %d\n", buffer, *buffer.refs)
	fmt.Printf("%v, %d\n", copy1, *copy1.refs)
	fmt.Printf("%v, %d\n", copy2, *copy2.refs)
}
