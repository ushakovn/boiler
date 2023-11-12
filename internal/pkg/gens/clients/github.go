package clients

import (
  "context"
  "encoding/base64"
  "encoding/json"
  "fmt"
  "net/http"
  "strings"
)

type githubContentResp struct {
  DownloadUrl string `json:"download_url"`
  Content     string `json:"content"`
  Encoding    string `json:"encoding"`
}

func (r *githubContentResp) Validate() error {
  if r.DownloadUrl == "" && (r.Content == "" || r.Encoding == "") {
    return fmt.Errorf("content not specified")
  }
  return nil
}

func (r *githubContentResp) HasDownloadUrl() bool {
  return r.DownloadUrl != ""
}

func (r *githubContentResp) HasContent() bool {
  return r.Content != "" && r.Encoding != ""
}

func (r *githubContentResp) DecodeContent() ([]byte, error) {
  if !strings.EqualFold(r.Encoding, "base64") {
    return nil, fmt.Errorf("content has unsupported encoding: %s", r.Encoding)
  }
  decoded, err := base64.StdEncoding.DecodeString(r.Content)
  if err != nil {
    return nil, fmt.Errorf("base64.StdEncoding.DecodeString: %w", err)
  }
  return decoded, nil
}

func (g *Clients) doGithubRequest(ctx context.Context, request string, response any) error {
  resp, err := g.githubClient.R().SetContext(ctx).Get(request)
  if err != nil {
    return fmt.Errorf("g.githubClient.Get: %s: %w", request, err)
  }
  if resp.StatusCode() != http.StatusOK || resp.Body() == nil {
    return fmt.Errorf("g.githubClient.Get: invalid response: resp.Status=%s", resp.Status())
  }
  if err = json.Unmarshal(resp.Body(), response); err != nil {
    return fmt.Errorf("json.Unmarshal: %w", err)
  }
  return nil
}

func (g *Clients) doGithubRequestRaw(ctx context.Context, request string) ([]byte, error) {
  resp, err := g.githubClient.R().SetContext(ctx).Get(request)
  if err != nil {
    return nil, fmt.Errorf("g.githubClient.Get: %s: %w", request, err)
  }
  if resp.StatusCode() != http.StatusOK && resp.Body() == nil {
    return nil, fmt.Errorf("g.githubClient.Get: invalid response: resp.Status=%s", resp.Status())
  }
  return resp.Body(), nil
}
