package pgdump

import (
  "context"
  "fmt"
  "os"

  validation "github.com/go-ozzo/ozzo-validation"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/internal/pkg/executor"
  "github.com/ushakovn/boiler/internal/pkg/filer"
)

type DumpOption struct {
  pgDump      []byte
  customTypes map[string]struct{}
}

func NewDumpOption() DumpOption {
  return DumpOption{}
}

func (o DumpOption) Validate() error {
  return validation.ValidateStruct(&o,
    validation.Field(&o.pgDump, validation.Required),
  )
}

func (o DumpOption) WithPgConfig(ctx context.Context, config PgConfig) DumpOption {
  pgDump, err := execPgDump(ctx, config)
  if err != nil {
    log.Fatalf("pg dump execution failed: %v", err)
  }
  o.pgDump = pgDump
  return o
}

func (o DumpOption) WithPgDumpPath(pgDumpPath string) DumpOption {
  if ext := filer.ExtractFileExtension(pgDumpPath); ext != "sql" {
    log.Fatalf("pg dump file must have sql file extension: %s", ext)
  }
  pgDump, err := os.ReadFile(pgDumpPath)
  if err != nil {
    log.Fatalf("pg dump file reading failed: %v", err)
  }
  o.pgDump = pgDump
  return o
}

func (o DumpOption) WithCustomTypes(customTypes []string) DumpOption {
  m := make(map[string]struct{}, len(customTypes))

  for _, customTyp := range customTypes {
    m[customTyp] = struct{}{}
  }
  o.customTypes = m

  return o
}

type PgConfig struct {
  Host     string
  Port     string
  User     string
  DBName   string
  Password string
}

func (c PgConfig) pgDumpCmd() (name string, args []string, err error) {
  if err = os.Setenv("PGPASSWORD", c.Password); err != nil {
    return "", nil, fmt.Errorf("os.Setenv: %w", err)
  }
  args = []string{
    "--no-owner",
    "-h", c.Host,
    "-p", c.Port,
    "-U", c.User,
    c.DBName,
  }
  name = "pg_dump"

  return name, args, nil
}

func execPgDump(ctx context.Context, config PgConfig) ([]byte, error) {
  name, args, err := config.pgDumpCmd()
  if err != nil {
    return nil, fmt.Errorf("config.pgDumpCmd: %w", err)
  }
  buf, err := executor.ExecCmdCtxWithOut(ctx, name, args...)
  if err != nil {
    return nil, fmt.Errorf("executor.ExecCmdCtxWithOut: %w", err)
  }
  return buf, nil
}
