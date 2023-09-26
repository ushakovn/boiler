package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func Test_CamelCaseToSnakeCase(t *testing.T) {
	type testCase struct {
		s        string
		expected string
	}
	testCases := []*testCase{
		{
			s:        "camelCasedString",
			expected: "camel_cased_string",
		},
		{
			s:        "StringWithCamelCase",
			expected: "string_with_camel_case",
		},
	}
	for _, tt := range testCases {
		tt := tt

		t.Run(tt.s, func(t *testing.T) {
			actual := CamelCaseToSnakeCase(tt.s)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_SnakeCaseToUpperCamelCase(t *testing.T) {
	type testCase struct {
		s        string
		expected string
	}
	testCases := []*testCase{
		{
			s:        "camel_cased_string",
			expected: "CamelCasedString",
		},
		{
			s:        "string_with_camel_case",
			expected: "StringWithCamelCase",
		},
	}
	for _, tt := range testCases {
		tt := tt

		t.Run(tt.s, func(t *testing.T) {
			actual := SnakeCaseToUpperCamelCase(tt.s)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
