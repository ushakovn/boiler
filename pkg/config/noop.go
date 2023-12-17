package config

var noopClient = newNoopClient()

type noopConfigClient struct{}

func newNoopClient() *noopConfigClient {
  return &noopConfigClient{}
}

func (c *noopConfigClient) GetValue(string) Value {
  return &configValue{}
}

func (c *noopConfigClient) GetAppInfo() AppInfo {
  const (
    name    = "app"
    version = "v0.0.1"
  )
  return AppInfo{
    Name:    name,
    Version: version,
  }
}
