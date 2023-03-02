package gomal_test

import (
	"reflect"
	"testing"

	"github.com/ItsMalma/gomal"
)

func TestNotNil(t *testing.T) {
	tests := []struct {
		name       string
		nameField  string
		valueField any
		results    []gomal.ValidationResult
	}{
		{
			name:       "success",
			nameField:  "x",
			valueField: "not nil",
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "failed",
			nameField:  "x",
			valueField: nil,
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x must not be empty."}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			results := gomal.Validate(gomal.If(test.nameField, test.valueField).NotNil())
			if !reflect.DeepEqual(results, test.results) {
				tt.Fatalf("expected %#v but got %#v instead", test.results, results)
			}
		})
	}
}

func TestNotEmpty(t *testing.T) {
	tests := []struct {
		name       string
		nameField  string
		valueField any
		results    []gomal.ValidationResult
	}{
		{
			name:       "success",
			nameField:  "x",
			valueField: "not empty",
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "failed because nil",
			nameField:  "x",
			valueField: nil,
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		},
		{
			name:       "failed because empty string",
			nameField:  "x",
			valueField: "",
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		},
		{
			name:       "failed because whitespace",
			nameField:  "x",
			valueField: "    \n\t\r",
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		},
		{
			name:       "failed because default (array)",
			nameField:  "x",
			valueField: [0]string{},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		},
		{
			name:       "failed because default (slice)",
			nameField:  "x",
			valueField: []string{},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		},
		{
			name:       "failed because default (map)",
			nameField:  "x",
			valueField: map[string]string{},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		},
		{
			name:       "failed because default (array's pointer)",
			nameField:  "x",
			valueField: &[0]string{},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		},
		{
			name:       "failed because default (boolean)",
			nameField:  "x",
			valueField: false,
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		}, {
			name:       "failed because default (complex)",
			nameField:  "x",
			valueField: 0 + 0i,
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		},
		{
			name:       "failed because default (float)",
			nameField:  "x",
			valueField: 0.0,
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		},
		{
			name:       "failed because default (int and uint)",
			nameField:  "x",
			valueField: 0,
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be empty."}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			results := gomal.Validate(gomal.If(test.nameField, test.valueField).NotEmpty())
			if !reflect.DeepEqual(results, test.results) {
				tt.Fatalf("expected %#v but got %#v instead", test.results, results)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	tests := []struct {
		name       string
		nameField  string
		valueField any
		another    any
		results    []gomal.ValidationResult
	}{
		{
			name:       "success (primitive)",
			nameField:  "x",
			valueField: "hello",
			another:    "world",
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "success (array)",
			nameField:  "x",
			valueField: [1]string{"hello"},
			another:    [1]string{"world"},
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "success (slice)",
			nameField:  "x",
			valueField: []string{"hello"},
			another:    []string{"world"},
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "success (map)",
			nameField:  "x",
			valueField: map[string]string{"hello": "world"},
			another:    map[string]string{"world": "hello"},
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "failed (primitive)",
			nameField:  "x",
			valueField: "hello",
			another:    "hello",
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be equal to hello."}}},
		},
		{
			name:       "failed (array without element)",
			nameField:  "x",
			valueField: [0]string{},
			another:    [0]string{},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be equal to []."}}},
		},
		{
			name:       "failed (array with element)",
			nameField:  "x",
			valueField: [1]string{"hello"},
			another:    [1]string{"hello"},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be equal to [hello]."}}},
		},
		{
			name:       "failed (slice without element)",
			nameField:  "x",
			valueField: []string{},
			another:    []string{},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be equal to []."}}},
		},
		{
			name:       "failed (slice with element)",
			nameField:  "x",
			valueField: []string{"world"},
			another:    []string{"world"},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be equal to [world]."}}},
		},
		{
			name:       "failed (map without element)",
			nameField:  "x",
			valueField: map[string]string{},
			another:    map[string]string{},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be equal to map[]."}}},
		},
		{
			name:       "failed (map with element)",
			nameField:  "x",
			valueField: map[string]string{"hello": "world"},
			another:    map[string]string{"hello": "world"},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should not be equal to map[hello:world]."}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			results := gomal.Validate(gomal.If(test.nameField, test.valueField).NotEqual(test.another))
			if !reflect.DeepEqual(results, test.results) {
				tt.Fatalf("expected %#v but got %#v instead", test.results, results)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name       string
		nameField  string
		valueField any
		another    any
		results    []gomal.ValidationResult
	}{
		{
			name:       "success (primitive)",
			nameField:  "x",
			valueField: "hello",
			another:    "hello",
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "success (array without element)",
			nameField:  "x",
			valueField: [0]string{},
			another:    [0]string{},
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "success (array with element)",
			nameField:  "x",
			valueField: [1]string{"hello"},
			another:    [1]string{"hello"},
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "success (slice without element)",
			nameField:  "x",
			valueField: []string{},
			another:    []string{},
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "success (slice with element)",
			nameField:  "x",
			valueField: []string{"world"},
			another:    []string{"world"},
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "success (map without element)",
			nameField:  "x",
			valueField: map[string]string{},
			another:    map[string]string{},
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "success (map with element)",
			nameField:  "x",
			valueField: map[string]string{"hello": "world"},
			another:    map[string]string{"hello": "world"},
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "failed (primitive)",
			nameField:  "x",
			valueField: "hello",
			another:    "world",
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should be equal to world."}}},
		},
		{
			name:       "failed (array)",
			nameField:  "x",
			valueField: [1]string{"hello"},
			another:    [1]string{"world"},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should be equal to [world]."}}},
		},
		{
			name:       "failed (slice)",
			nameField:  "x",
			valueField: []string{"hello"},
			another:    []string{"world"},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should be equal to [world]."}}},
		},
		{
			name:       "failed (map)",
			nameField:  "x",
			valueField: map[string]string{"hello": "world"},
			another:    map[string]string{"world": "hello"},
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"x should be equal to map[world:hello]."}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			results := gomal.Validate(gomal.If(test.nameField, test.valueField).Equal(test.another))
			if !reflect.DeepEqual(results, test.results) {
				tt.Fatalf("expected %#v but got %#v instead", test.results, results)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name       string
		nameField  string
		valueField string
		min        int
		results    []gomal.ValidationResult
	}{
		{
			name:       "minimum 0",
			nameField:  "x",
			valueField: "Hello World",
			min:        0,
			results:    []gomal.ValidationResult{},
		},
		{
			name:       "minimum 1 but failed",
			nameField:  "x",
			valueField: "",
			min:        1,
			results:    []gomal.ValidationResult{{Name: "x", Messages: []string{"The length of x must be at least 1 characters. You entered 0 characters."}}},
		},
		{
			name:       "minimum 1 but pass",
			nameField:  "x",
			valueField: "H",
			min:        1,
			results:    []gomal.ValidationResult{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			results := gomal.Validate(gomal.If(test.nameField, test.valueField).MinLength(test.min))
			if !reflect.DeepEqual(results, test.results) {
				tt.Fatalf("expected %#v but got %#v instead", test.results, results)
			}
		})
	}
}
