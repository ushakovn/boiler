package config

import "context"

type ctxKey struct{}

func ContextClient(ctx context.Context) Client {
  if ctxClient, ok := ctx.Value(ctxKey{}).(Client); ok {
    return ctxClient
  }
  if client == nil {
    return noopClient
  }
  return client
}

func ContextWithClient(parent context.Context, client Client) context.Context {
  return context.WithValue(parent, ctxKey{}, client)
}
