package executor

import (
  "context"
  "fmt"
  "os/exec"
  "strings"
)

func ExecCmdCtx(ctx context.Context, name string, args ...string) error {
  _, err := ExecCmdCtxWithOut(ctx, name, args...)
  return err
}

func ExecCmdCtxWithOut(ctx context.Context, name string, args ...string) ([]byte, error) {
  cmd := exec.CommandContext(ctx, name, args...)
  if cmd == nil {
    return nil, fmt.Errorf("exec.CommandContext: cmd is a nil")
  }
  buf, err := cmd.CombinedOutput()
  if err != nil {
    str := strings.TrimSpace(string(buf))
    return nil, fmt.Errorf("exec.CommandContext:\n\toutput: %s\n\terror: %s\n\tcmd.Error: %s", str, err, cmd.Err)
  }
  return buf, nil
}
