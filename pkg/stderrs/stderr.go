package stderrs

import (
  "encoding/json"
  "errors"
  "fmt"
  "net/http"
  "strings"

  "google.golang.org/grpc/codes"
)

// Error standard error.
type Error struct {
  Code    string         `json:"code"`
  Message string         `json:"message,omitempty"`
  Err     error          `json:"error,omitempty"`
  Wraps   []string       `json:"wraps,omitempty"`
  Fields  map[string]any `json:"fields,omitempty"`
  Codes   struct {
    GRPC codes.Code `json:"grpc"`
    HTTP int        `json:"http"`
  } `json:"codes"`
}

// New create new standard Error.
func New(msg string) *Error {
  return &Error{
    Code: msg,
    Codes: struct {
      GRPC codes.Code `json:"grpc"`
      HTTP int        `json:"http"`
    }{
      GRPC: codes.Internal,
      HTTP: http.StatusInternalServerError,
    },
  }
}

// Error implement builtin error interface.
func (e *Error) Error() string {
  if e == nil {
    return ""
  }

  var parts []string

  if len(e.Code) != 0 {
    parts = append(parts, fmt.Sprintf(`"code": "%s"`, e.Code))
  } else {
    parts = append(parts, fmt.Sprintf(`"code": "undefined"`))
  }

  if len(e.Message) != 0 {
    parts = append(parts, fmt.Sprintf(`"message": "%s"`, e.Message))
  }

  if len(e.Fields) != 0 {
    raw, err := json.Marshal(e.Fields)
    if err == nil {
      parts = append(parts, fmt.Sprintf(`"fields": %s`, raw))
    }
  }

  switch v := e.Err.(type) {

  case interface{ Unwrap() []error }:
    var messages []string

    for _, item := range v.Unwrap() {
      var std *Error

      if errors.As(item, &std) {
        messages = append(messages, item.Error())
      }
    }
    parts = append(parts, fmt.Sprintf(`"embed": [%s]`, strings.Join(messages, ", ")))

  case interface{ Unwrap() error }:
    parts = append(parts, fmt.Sprintf(`"embed": [%s]`, e.Err.Error()))

  default:
    if v == nil {
      break
    }
    parts = append(parts, fmt.Sprintf(`"embed": [%s]`, e.Err.Error()))
  }

  message := fmt.Sprintf("{%s}", strings.Join(parts, ", "))

  for i := 0; i < len(e.Wraps); i++ {
    message = fmt.Sprintf("%s > %s", e.Wraps[i], message)
  }

  return message
}

// SetMessage set message to Error copy.
func (e *Error) SetMessage(format string, args ...any) *Error {
  if e == nil {
    return nil
  }

  cp := *e

  cp.Message = fmt.Sprintf(format, args...)

  return &cp
}

// SetHTTPCode set HTTP code to Error copy.
func (e *Error) SetHTTPCode(code int) *Error {
  if e == nil {
    return nil
  }

  cp := *e

  cp.Codes.HTTP = code

  return &cp
}

// SetGRPCCode set gRPC code to Error copy.
func (e *Error) SetGRPCCode(code codes.Code) *Error {
  if e == nil {
    return nil
  }

  cp := *e

  cp.Codes.GRPC = code

  return &cp
}

// Join builtin errors to Error copy.
func (e *Error) Join(errs ...error) *Error {
  if e == nil {
    return nil
  }

  cp := *e

  switch v := e.Err.(type) {

  case interface{ Unwrap() error }:
    join := append([]error{v.Unwrap()}, errs...)

    cp.Err = errors.Join(join...)

  case interface{ Unwrap() []error }:
    join := append(v.Unwrap(), errs...)

    cp.Err = errors.Join(join...)

  default:
    cp.Err = errors.Join(errs...)
  }

  return &cp
}

// Wrap message to Error copy.
func (e *Error) Wrap(msg ...string) *Error {
  if e == nil {
    return nil

  }

  cp := *e
  wraps := make([]string, len(e.Wraps))

  copy(wraps, e.Wraps)

  cp.Wraps = append(wraps, msg...)

  return &cp
}

// Interface return builtin error interface.
func (e *Error) Interface() error {
  if e == nil {
    return nil
  }

  return e
}

// WithField set field to Error copy.
func (e *Error) WithField(key string, value any) *Error {
  res := make(map[string]any)

  for k, v := range e.Fields {
    res[k] = v
  }
  res[key] = value

  cp := *e
  cp.Fields = res

  return &cp
}

// WithFieldIf set field to Error copy.
func (e *Error) WithFieldIf(condition bool, key string, value any) *Error {
  if !condition {
    return e
  }

  res := make(map[string]any)

  for k, v := range e.Fields {
    res[k] = v
  }
  res[key] = value

  cp := *e
  cp.Fields = res

  return &cp
}

// WithFields set fields to Error copy.
func (e *Error) WithFields(fields map[string]any) *Error {
  res := make(map[string]any)

  for key, value := range e.Fields {
    res[key] = value
  }

  for key, value := range fields {
    res[key] = value
  }

  cp := *e
  cp.Fields = res

  return &cp
}

// Unwrap implement errors.Unwrap method.
func (e *Error) Unwrap() error {
  if e == nil {
    return nil
  }

  return e.Err
}

func unwrap(err error) []error {
  switch v := err.(type) {

  case interface{ Unwrap() error }:
    return []error{v.Unwrap()}

  case interface{ Unwrap() []error }:
    return v.Unwrap()

  default:
    return nil
  }
}

// Is implement errors.Is method.
func (e *Error) Is(err error) bool {
  if e == nil || err == nil {
    return false
  }

  var std *Error

  if errors.As(err, &std) {
    switch {
    case e.Code == std.Code && e.Err != nil && std.Err == nil:
      return true

    case e.Code == std.Code && e.Err == nil && std.Err == nil:
      return true

    case e.Code == std.Code && e.Err == nil && std.Err != nil:
      return false

    case e.Code == std.Code && e.Err != nil && std.Err != nil:
      for _, item := range unwrap(std.Err) {
        if !errors.Is(e.Err, item) {
          return false
        }
      }
      return true

    case e.Code != std.Code && e.Err != nil:
      return errors.Is(e.Err, err)
    }
  }

  for _, item := range unwrap(err) {
    if errors.Is(e, item) {
      return true
    }
  }

  if errors.Is(e.Err, err) {
    return true
  }

  return false
}
