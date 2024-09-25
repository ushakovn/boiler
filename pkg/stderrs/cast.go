package stderrs

import (
  "errors"

  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
)

var grpcCodeToStdErr = map[codes.Code]*Error{
  codes.Canceled:           Canceled,
  codes.Unknown:            Unknown,
  codes.InvalidArgument:    InvalidArgument,
  codes.DeadlineExceeded:   DeadlineExceeded,
  codes.NotFound:           NotFound,
  codes.AlreadyExists:      AlreadyExists,
  codes.PermissionDenied:   PermissionDenied,
  codes.ResourceExhausted:  ResourceExhausted,
  codes.FailedPrecondition: FailedPrecondition,
  codes.Aborted:            Aborted,
  codes.OutOfRange:         OutOfRange,
  codes.Unimplemented:      Unimplemented,
  codes.Internal:           Internal,
  codes.Unavailable:        Unavailable,
  codes.DataLoss:           DataLoss,
  codes.Unauthenticated:    Unauthenticated,
}

var castFunctions = []func(err error) *Error{
  castBuiltin,
  castGRPC,
}

// Cast allows you to restore the standard error from the interface.
// If the error is not standard, it will be returned nil Error.
func Cast(err error) (*Error, bool) {
  if err == nil {
    return nil, false
  }

  for _, handler := range castFunctions {
    if std := handler(err); std != nil {
      return std, true
    }
  }

  return nil, false
}

func castBuiltin(err error) *Error {
  var std *Error

  if errors.As(err, &std) {
    return std
  }

  return nil
}

func castGRPC(err error) *Error {
  parsed, ok := status.FromError(err)
  if !ok {
    return nil
  }

  std, ok := grpcCodeToStdErr[parsed.Code()]
  if !ok {
    return nil
  }

  return std
}
