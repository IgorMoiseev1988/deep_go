package main

import (
	"unsafe"
	"testing"
	"github.com/stretchr/testify/assert"
)

type (
	TestsMap = map[string]struct{number uint32; result uint32}
)

/* Run benchmark: go test -bench=. homework_test.go
/* Run test     : go test -v homework_test.go

/* Benchmark results:
	goos: linux
	goarch: amd64
	cpu: AMD Ryzen 5 5600X 6-Core Processor
	BenchmarkConversion_1-12        1000000000               0.0002078 ns/op
	BenchmarkConversion_2-12        1000000000               0.001603 ns/op
	BenchmarkConversion_3-12        1000000000               0.002320 ns/op
	PASS
	ok      command-line-arguments  0.033s
*/


/* first version: fast, but need a separate function for each type */
func ToLittleEndian_1(number uint32) uint32 {
	return number >> 24 | (number & 0x00FF0000) >> 8 | (number & 0x0000FF00) << 8 | number << 24
}


/* second version: generic */
func ToLittleEndian_2[T uint16 | uint32 | uint64](number T) T {
	var result T
         
	const (
		mask = 0xFF
		byte_bit_sz = 8
	)
	
	type_bit_sz := unsafe.Sizeof(number) * byte_bit_sz

	var r_shift uintptr = type_bit_sz
	var l_shift uintptr

	for l_shift = 0; l_shift < type_bit_sz; l_shift += byte_bit_sz {
		r_shift -= byte_bit_sz
		result |= ((number & (mask << r_shift)) >> r_shift) << l_shift
	
	}
        return result	
} 

/* Third version: simple, but slow */
func ToLittleEndian_3[T uint16 | uint32 | uint64](number T) T {
	p := unsafe.Pointer(&number)
	var type_size uintptr = unsafe.Sizeof(number)
	var offset uintptr

	for offset = 0; offset < type_size / 2; offset++ {
		lhs := (*uint8)(unsafe.Add(p, offset))
		rhs := (*uint8)(unsafe.Add(p, 4 - offset - 1))
		/* swap bytes */
		*lhs ^= *rhs
		*rhs ^= *lhs
		*lhs ^= *rhs
	}
	return number
}

func GetTests() TestsMap {
	return TestsMap {
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0xFFFF0000,
			result: 0x0000FFFF,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}
}

func TestConversion_1(t* testing.T) {
	tests := GetTests()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian_1(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestConversion_2(t* testing.T) {
	tests := GetTests()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian_2(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestConversion_3(t* testing.T) {
	tests := GetTests()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian_3(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

func BenchmarkConversion_1(b *testing.B) {
	var max uint32 = 1 << 20
	var idx uint32
	for idx = 0; idx < max; idx++ {
		_ = ToLittleEndian_1(idx)
	}
}

func BenchmarkConversion_2(b *testing.B) {
	var max uint32 = 1 << 20
	var idx uint32
	for idx = 0; idx < max; idx++ {
		_ = ToLittleEndian_2(idx)
	}
}

func BenchmarkConversion_3(b *testing.B) {
	var max uint32 = 1 << 20
	var idx uint32
	for idx = 0; idx < max; idx++ {
		_ = ToLittleEndian_3(idx)
	}
}
