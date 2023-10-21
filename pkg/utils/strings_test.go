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
    {
      s:        "BatchUploadPhotosV2",
      expected: "batch_upload_photos_v2",
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
    {
      s:        "batch_upload_photos_v2",
      expected: "BatchUploadPhotosV2",
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

func Test_CamelCaseToUpperCamelCase(t *testing.T) {
  type testCase struct {
    s        string
    expected string
  }
  testCases := []*testCase{
    {
      s:        "camelcasedString",
      expected: "CamelcasedString",
    },
    {
      s:        "camel",
      expected: "Camel",
    },
    {
      s:        "BatchUploadPhotosV2",
      expected: "BatchUploadPhotosV2",
    },
  }
  for _, tt := range testCases {
    tt := tt

    t.Run(tt.s, func(t *testing.T) {
      actual := CamelCaseToUpperCamelCase(tt.s)
      assert.Equal(t, tt.expected, actual)
    })
  }
}

func Test_StringToUpperCamelCase(t *testing.T) {
  type testCase struct {
    s        string
    expected string
  }
  testCases := []*testCase{
    {
      s:        "snake_cased_string",
      expected: "SnakeCasedString",
    },
    {
      s:        "camel",
      expected: "Camel",
    },
    {
      s:        "batchUploadPhotosV2",
      expected: "BatchUploadPhotosV2",
    },
  }
  for _, tt := range testCases {
    tt := tt

    t.Run(tt.s, func(t *testing.T) {
      actual := StringToUpperCamelCase(tt.s)
      assert.Equal(t, tt.expected, actual)
    })
  }
}

func Test_StringToLowerCamelCase(t *testing.T) {
  type testCase struct {
    s        string
    expected string
  }
  testCases := []*testCase{
    {
      s:        "snake_cased_string",
      expected: "snakeCasedString",
    },
    {
      s:        "camel",
      expected: "camel",
    },
    {
      s:        "Camel",
      expected: "camel",
    },
    {
      s:        "BatchUploadPhotosV2",
      expected: "batchUploadPhotosV2",
    },
  }
  for _, tt := range testCases {
    tt := tt

    t.Run(tt.s, func(t *testing.T) {
      actual := StringToLowerCamelCase(tt.s)
      assert.Equal(t, tt.expected, actual)
    })
  }
}
