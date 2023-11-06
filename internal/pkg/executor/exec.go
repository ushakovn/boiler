package executor

import (
  "context"
  "fmt"
  "os/exec"
  "strings"
)

func ExecCommandContext(ctx context.Context, name string, args ...string) error {
  cmd := exec.CommandContext(ctx, name, args...)
  if cmd == nil {
    return fmt.Errorf("exec.CommandContext: cmd is a nil")
  }
  buf, err := cmd.CombinedOutput()
  if err != nil {
    str := strings.TrimSpace(string(buf))
    return fmt.Errorf("exec.CommandContext:\n\toutput: %s\n\terror: %s\n\tcmd.Error: %s", str, err, cmd.Err)
  }
  return nil
}

