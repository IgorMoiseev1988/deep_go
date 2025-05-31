package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize[T any](t T) string {
	valueT := reflect.ValueOf(t)
	typeT := valueT.Type()
	var sb strings.Builder

	for i := 0; i < typeT.NumField(); i++ {
		// Get type and value of field
		field := typeT.Field(i)
		fieldValue := valueT.Field(i)
		var omitempty bool

		// get tags
		tags := strings.Split(field.Tag.Get("properties"), ",")
		if len(tags) == 0 {
			continue // skip field if it doesn't has tag properties
		}
		name := tags[0]
		if len(tags) > 1 {
			omitempty = tags[1] == "omitempty"
		}

		// check omitempty in tag and value of field is zero-value
		if omitempty {
			zero := reflect.Zero(field.Type).Interface()
			current := fieldValue.Interface()
			if reflect.DeepEqual(zero, current) {
				continue
			}
		}

		// Append tag-value pair to result
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		fmt.Fprintf(&sb, "%s=%v", name, fieldValue)
	}
	return sb.String()
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
