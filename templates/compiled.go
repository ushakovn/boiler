package templates

// Rpc Generator compiled templates
const (
  // RpcContracts const for compiled Boiler build with contracts template
  RpcContracts = "// Code generated by Boiler; DO NOT EDIT.\n\npackage handler\n\n{{range .Requests}}\n// {{.Name}}Request ...\ntype {{.Name}}Request struct {\n  {{range .Fields}}\n  {{.Name}} {{.Type}} `json:\"{{.Tag}}\"`\n  {{- end}}\n}\n{{- end}}\n\n{{range .Responses}}\n// {{.Name}}Response ...\ntype {{.Name}}Response struct {\n  {{range .Fields}}\n  {{.Name}} {{.Type}} `json:\"{{.Tag}}\"`\n  {{- end}}\n}\n{{- end}}\n\n{{range .TypeDefs}}\n// {{.Name}} ...\ntype {{.Name}} struct{\n  {{range .Fields}}\n  {{.Name}} {{.Type}} `json:\"{{.Tag}}\"`\n  {{- end}}\n}\n{{- end}}"
  // RpcHandle const for compiled Boiler build with handle template
  RpcHandle = "// Code generated by Boiler; YOU MAY CHANGE THIS.\n\npackage handler\n\nimport (\n\t\"net/http\"\n\n\t\"github.com/gin-gonic/gin\"\n)\n\n// Handle{{.Name}} stub ...\nfunc (h *Handler) Handle{{.Name}}(ctx *gin.Context) {\n  req := &{{.Name}}Request{}\n\n  if err := ctx.BindJSON(req); err != nil {\n    ctx.JSON(http.StatusBadRequest, err)\n    return\n  }\n\n  resp := &{{.Name}}Response{}\n  ctx.JSON(http.StatusOK, resp)\n}\n"
  // RpcHandler const for compiled Boiler build with handler template
  RpcHandler = "// Code generated by Boiler; DO NOT EDIT.\n\npackage handler\n\nimport (\n  \"github.com/gin-gonic/gin\"\n  log \"github.com/sirupsen/logrus\"\n)\n\n// Handler ...\ntype Handler struct {\n  addr string\n  g    *gin.Engine\n}\n\n// NewHandler ...\nfunc NewHandler() *Handler {\n  return &Handler{}\n}\n\ntype (\n  route  string\n  Routes map[route]gin.IRoutes\n)\n\nconst (\n    {{- range .Handles}}\n\t{{.Name}}Route route = \"{{.Route}}\"\n\t{{- end}}\n)\n\n// RegisterRoutes ...\nfunc (h *Handler) RegisterRoutes() Routes {\n  h.g = gin.New()\n  // Registered routes\n  return Routes{\n    {{- range .Handles}}\n    {{.Name}}Route: h.g.POST(string({{.Name}}Route), h.Handle{{.Name}}),\n    {{- end}}\n  }\n}\n\n// ServeHTTP ...\nfunc (h *Handler) ServeHTTP() {\n  go func() {\n    if err := h.g.Run(h.addr); err != nil {\n      log.Fatalf(\"ServeHTTP error: gin.Run: %v\", err)\n    }\n  }()\n}\n"
)

// Project Generator compiled templates
const (
  // ProjectMain const for compiled Boiler build with main template
  ProjectMain = "// Code generated by Boiler. YOU MAY CHANGE THIS\npackage main\n\nimport \"github.com/ushakovn/boiler/pkg/app\"\n\nfunc main() {\n  a := app.NewApp()\n  a.Run()\n}"
  // ProjectGomod const for compiled Boiler build with go mod template
  ProjectGomod = "module {{.goModName}}\n\ngo {{.goModVersion}}\n"
  // ProjectMakefile const for compiled Boiler build with makefile template
  ProjectMakefile = "# Including generated by Boiler. DO NOT EDIT.\ninclude make.mk\n\n# Target section. YOU MAY CHANGE THIS.\n.PHONY: help\nhelp:\n\tboiler --help\n"

  // ProjectDockerCompose const for compiled Boiler build with docker compose template
  ProjectDockerCompose = "# Docker compose generated by Boiler; DO NOT EDIT.\nversion: '3.7'\n\nservices:\n  # Metrics section\n  prometheus:\n    image: prom/prometheus:latest\n    volumes:\n      - ./prometheus_config.yaml:/etc/prometheus/prometheus.yml\n    ports:\n      - \"9090:9090\"\n    networks:\n      - boiler-network\n  grafana:\n    image: grafana/grafana-oss:9.4.3\n    ports:\n      - \"3000:3000\"\n    volumes:\n      - boiler-volume:/var/lib/grafana\n    networks:\n      - boiler-network\n\n  # Tracing sections\n  jaeger:\n    image: jaegertracing/all-in-one:latest\n    ports:\n      - \"16686:16686\"\n      - \"4317:4317\"\n      - \"4318:4318\"\n    networks:\n      - boiler-network\n\n# Boiler docker network\nnetworks:\n  boiler-network:\n\n# Boiler docker volume\nvolumes:\n  boiler-volume:\n"
  // ProjectPrometheusConfig const for compiled Boiler build with prometheus config template
  ProjectPrometheusConfig = "# Prometheus config generated by Boiler; DO NOT EDIT.\nglobal:\n  scrape_interval:     30s\n  evaluation_interval: 30s\n\nscrape_configs:\n  # Prometheus node job\n  - job_name: 'prometheus'\n    scrape_interval: 15s\n    scrape_timeout: 15s\n    static_configs:\n      - targets: ['localhost:9090']\n\n  # Boiler app job\n  - job_name: 'boiler'\n    scrape_interval: 15s\n    scrape_timeout: 15s\n    static_configs:\n      - targets: ['host.docker.internal:8092']\n"
)

// Project Generator compiled templates names
const (
  // NameMain name for compiled ProjectMain file template
  NameMain = "main"
  // NameGomod name for compiled ProjectGomod file template
  NameGomod = "gomod"
  // NameMakefile name for compiled ProjectMakefile file template
  NameMakefile = "makefile"
  // NameDockerCompose name for compiled ProjectDockerCompose file template
  NameDockerCompose = "docker_compose"
  // NamePrometheusConfig name for compiled ProjectPrometheusConfig file template
  NamePrometheusConfig = "prometheus_config"
)

const (
  // StorageFilterIfStmtWithPtr const for compiled Boiler build with filter if statement for zero typed filters
  StorageFilterIfStmtWithPtr = "filters.%s.Ptr() != nil"
  // StorageFilterIfStmtWithLen const for compiled Boiler build with filter if statement for slice typed filters
  StorageFilterIfStmtWithLen = "len(filters.%s) > 0"

  // StorageInputIfStmtWithPtr const for compiled Boiler build with input statement for zero typed fields
  StorageInputIfStmtWithPtr = "input.%s.Ptr() != nil"
  // StorageInputIfStmtWithLen const for compiled Boiler build with input statement for slice typed fields
  StorageInputIfStmtWithLen = "len(input.%s) > 0"

  // StorageConsts ...
  StorageConsts = "// Code generated by Boiler; DO NOT EDIT.\n\npackage storage\n\n{{- $count := len .Models}}\n{{- if ne $count 0 }}\ntype tableName string\n\nconst (\n  {{- range .Models}}\n  {{.ModelName}}_TableName tableName = \"{{.SqlTableName}}\"\n  {{- end}}\n)\n\ntype (\n  {{- range .Models}}\n  {{toLowerCamelCase .ModelName}}_Field string\n  {{- end}}\n)\n\n{{- range $m := .Models}}\nconst (\n  {{- range .ModelFields}}\n  {{.FieldName}}_{{$m.ModelName}}_Field {{toLowerCamelCase $m.ModelName}}_Field = \"{{.SqlTableFieldName}}\"\n  {{- end}}\n)\n{{end}}\n{{- end}}\n"
  // StorageModelMethods ...
  StorageModelMethods = "// Code generated by Boiler; DO NOT EDIT. {{$modelName := .ModelName}}\n\npackage storage\n\nimport (\n  {{- range .ModelMethodsPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\nfunc (s *Storage) List{{$modelName}}(ctx context.Context, input List{{$modelName}}Input) ([]*models.{{$modelName}}, error) {\n  builder := br.NewSelectBuilder().\n    Columns(\n      {{- range .ModelFields}}\n      string({{.FieldName}}_{{$modelName}}_Field),\n      {{- end}}\n    ).\n    From(string({{$modelName}}_TableName))\n\n  if err := input.Pagination.validate(); err != nil {\n    return nil, fmt.Errorf(\"pagination.Validate: %w\", err)\n  }\n\n  offset, limit := input.Pagination.orDefault().toOffsetLimit()\n  builder = builder.Offset(offset).Limit(limit)\n\n  if input.Sort != nil {\n    builder = builder.OrderBy(input.Sort.{{toLowerCamelCase $modelName}}Sort())\n  }\n\n  if filters := input.Filters; filters != nil {\n    {{- range $modelField := .ModelFields}}\n    {{- range .ModelsFieldFilters}}\n    if {{.FilterIfStmt}} {\n      builder = builder.Where(sq.{{.FilterSqOperator}}{string({{$modelField.FieldName}}_{{$modelName}}_Field): filters.{{.FilterName}}{{.FilterTypeSuffix}}})\n    }\n    {{- end}}\n    {{- end}}\n  }\n\n  return pg.SelectCtx[*models.{{$modelName}}](ctx, s.executor, builder)\n}\n\nfunc (s *Storage) {{$modelName}}(ctx context.Context, input {{$modelName}}Input) (*models.{{$modelName}}, error) {\n  builder := br.NewSelectBuilder().\n    Columns(\n      {{- range .ModelFields}}\n      string({{.FieldName}}_{{$modelName}}_Field),\n      {{- end}}\n    ).\n    From(string({{$modelName}}_TableName))\n\n  if filters := input.Filters; filters != nil {\n    {{- range $modelField := .ModelFields}}\n    {{- range .ModelFieldFilters}}\n    if {{.FilterIfStmt}} {\n      builder = builder.Where(sq.Eq{string({{$modelField.FieldName}}_{{$modelName}}_Field): filters.{{.FilterName}}{{.FilterTypeSuffix}}})\n    }\n    {{- end}}\n    {{- end}}\n  }\n\n  model, err := pg.GetCtx[*models.{{$modelName}}](ctx, s.executor, builder)\n  if err != nil {\n    return nil, err\n  }\n  return model, nil\n}\n\nfunc (s *Storage) Create{{$modelName}}(ctx context.Context, input Create{{$modelName}}Input) (*models.{{$modelName}}, error) {\n  builder := br.NewInsertBuilder().\n    Into(string({{$modelName}}_TableName))\n\n  fields := map[string]any{\n    {{- range .ModelFields}}\n    {{- if eq .NotNullField true}}\n    {{- if eq .WithDefaultField false}}\n    string({{.FieldName}}_{{$modelName}}_Field): input.{{.FieldName}},\n    {{- end}}\n    {{- end}}\n    {{- end}}\n  }\n\n  {{- range $modelField := .ModelFields}}\n  {{- if eq .NotNullField false}}\n  {{- if eq .WithDefaultField false}}\n  if {{.FieldIfStmt}} {\n    fields[string({{$modelField.FieldName}}_{{$modelName}}_Field)] = input.{{.FieldName}}{{.FieldTypeSuffix}}\n  }\n  {{- end}}\n  {{- end}}\n  {{- end}}\n\n  builder = builder.SetMap(fields).\n    Suffix(suffixReturning)\n\n  model, err := pg.GetCtx[*models.{{$modelName}}](ctx, s.executor, builder)\n  if err != nil {\n    return nil, err\n  }\n  return model, nil\n}\n\nfunc (s *Storage) Update{{$modelName}}(ctx context.Context, input Update{{$modelName}}Input) (*models.{{$modelName}}, error) {\n  builder := br.NewUpdateBuilder().\n    Table(string({{$modelName}}_TableName))\n\n  fields := map[string]any{}\n\n  {{- range .ModelFields}}\n  {{- if ne .FieldBadge \"pk\"}}\n  if {{.FieldZeroTypeIfStmt}} {\n    fields[string({{.FieldName}}_{{$modelName}}_Field)] = input.{{.FieldName}}{{.FieldZeroTypeSuffix}}\n  }\n  {{- end}}\n  {{- end}}\n\n  builder = builder.\n    SetMap(fields).\n    {{- range .ModelFields}}\n    {{- if eq .FieldBadge \"pk\"}}\n    Where(sq.Eq{string({{.FieldName}}_{{$modelName}}_Field): input.{{.FieldName}}}).\n    {{- end}}\n    {{- end}}\n    Suffix(suffixReturning)\n\n  model, err := pg.GetCtx[*models.{{$modelName}}](ctx, s.executor, builder)\n  if err != nil {\n    return nil, err\n  }\n  return model, nil\n}\n\nfunc (s *Storage) Delete{{$modelName}}(ctx context.Context, input Delete{{$modelName}}Input) (*models.{{.ModelName}}, error) {\n  builder := br.NewDeleteBuilder().\n    From(string({{$modelName}}_TableName)).\n    {{- range .ModelFields}}\n    {{- if eq .FieldBadge \"pk\"}}\n    Where(sq.Eq{string({{.FieldName}}_{{$modelName}}_Field): input.{{.FieldName}}}).\n    {{- end}}\n    {{- end}}\n    Suffix(suffixReturning)\n\n  model, err := pg.GetCtx[*models.{{$modelName}}](ctx, s.executor, builder)\n  if err != nil {\n    return nil, err\n  }\n  return model, nil\n}\n"
  // StorageModelOptions ...
  StorageModelOptions = "// Code generated by Boiler; DO NOT EDIT.\n\npackage storage\n\nimport (\n  {{- range .ModelOptionsPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\n// Suppress unused imports\nvar (\n    _ = zero.Time{}\n    _ = time.Time{}\n)\n\ntype {{.ModelName}}Input struct {\n  Filters *{{.ModelName}}Filter\n}\n\ntype {{.ModelName}}Filter struct {\n  {{- range .ModelFields}}\n  {{- range .ModelFieldFilters}}\n  {{.FilterName}} {{.FilterType}}\n  {{- end}}\n  {{- end}}\n}\n\ntype Create{{.ModelName}}Input struct {\n  {{- range .ModelFields}}\n  {{- if eq .WithDefaultField false}}\n  {{.FieldName}} {{.FieldType}} {{- if eq .FieldBadge \"pk\"}} // PRIMARY KEY{{- end}}\n  {{- end}}\n  {{- end}}\n}\n\ntype Delete{{.ModelName}}Input struct {\n  {{- range .ModelFields}}\n  {{- if eq .FieldBadge \"pk\"}}\n  {{.FieldName}} {{.FieldType}} // PRIMARY KEY\n  {{- end}}\n  {{- end}}\n}\n\ntype Update{{.ModelName}}Input struct {\n  {{- range .ModelFields}}\n  {{- if eq .FieldBadge \"pk\"}}\n  {{.FieldName}} {{.FieldType}} // PRIMARY KEY\n  {{- else}}\n  {{.FieldName}} {{.FieldZeroType}}\n  {{- end}}\n  {{- end}}\n}\n\ntype List{{.ModelName}}Input struct {\n  Filters    *List{{.ModelName}}Filters\n  Sort       *List{{.ModelName}}Sort\n  Pagination *Pagination\n}\n\ntype List{{.ModelName}}Sort struct {\n  Field {{toLowerCamelCase .ModelName}}_Field\n  Order sortOrder\n}\n\ntype List{{.ModelName}}Filters struct {\n  {{- range .ModelFields}}\n  {{- range .ModelsFieldFilters}}\n  {{.FilterName}} {{.FilterType}}\n  {{- end}}\n  {{- end}}\n}\n\nfunc (p *List{{.ModelName}}Sort) {{toLowerCamelCase .ModelName}}Sort() string {\n  return fmt.Sprintf(\"%s %s\", p.Field, p.Order)\n}\n"
  // StorageModel ...
  StorageModel = "// Code generated by Boiler; DO NOT EDIT.\n\npackage models\n\nimport (\n  {{- range .ModelPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\n// Suppress unused imports\nvar (\n    _ = zero.Time{}\n    _ = time.Time{}\n)\n\ntype {{.ModelName}} struct {\n  {{- range .ModelFields}}\n  {{.FieldName}} {{.FieldType}} `db:\"{{.SqlTableFieldName}}\"` {{- if eq .FieldBadge \"pk\"}} // PRIMARY KEY{{- end}}\n  {{- end}}\n}\n"
  // StorageOptions ...
  StorageOptions = "// Code generated by Boiler; DO NOT EDIT.\n\npackage storage\n\nimport (\n  {{- range .OptionsPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\nconst suffixReturning = \"RETURNING *\"\n\ntype sortOrder string\n\nconst (\n  SortOrderAsc  sortOrder = \"ASC\"\n  SortOrderDesc sortOrder = \"DESC\"\n  SortOrderRand sortOrder = \"RAND()\"\n)\n\ntype Pagination struct {\n  Page    uint64\n  PerPage uint64\n}\n\nfunc (p *Pagination) orDefault() *Pagination {\n  if p != nil {\n    return p\n  }\n  return &Pagination{\n    Page:    0,\n    PerPage: 100,\n  }\n}\n\nfunc (p *Pagination) validate() error {\n  if p == nil {\n    return nil\n  }\n  if p.Page < 0 {\n    return fmt.Errorf(\"pagination.Page=%d must be non-negative\", p.Page)\n  }\n  if p.PerPage < 0 {\n    return fmt.Errorf(\"pagination.PerPage=%d must be positive\", p.PerPage)\n  }\n  return nil\n}\n\nfunc (p *Pagination) toOffsetLimit() (offset uint64, limit uint64) {\n  offset = p.Page * p.PerPage\n  limit = p.PerPage\n  return\n}\n\ntype Storage struct {\n  executor pg.Executor\n}\n\nfunc NewStorage(executor pg.Executor) *Storage {\n  return &Storage{\n    executor: executor,\n  }\n}\n\nfunc (s *Storage) WithTransaction(ctx context.Context, fTx func(*Storage) error) error {\n  defer func() {\n    if rec := recover(); rec != nil {\n      log.Errorf(\"client.WithTransaction: panic recovered: %v\", rec)\n    }\n  }()\n\n  tx, err := s.executor.Begin(ctx)\n  if err != nil {\n    return fmt.Errorf(\"s.executor.BeginTx: %w\", err)\n  }\n  txStorage := NewStorage(tx)\n\n  if err = fTx(txStorage); err != nil {\n    if errTx := tx.Rollback(ctx); errTx != nil {\n      log.Errorf(\"Storage.WithTransaction: tx.Rollback: %v\", errTx)\n    }\n    return err\n  }\n\n  if txErr := tx.Commit(ctx); txErr != nil {\n    return fmt.Errorf(\"tx.Commit: %w\", err)\n  }\n  return nil\n}\n\n"

  // StorageConfig ...
  StorageConfig = "# Config generated by Boiler; YOU MUST CHANGE THIS.\n\n# Boiler storage generator config\nstorage_config:\n  # 1. Fill one of config sections\n\n  # 1.1. Connection config\n  pg_config:\n    host: \"\"\n    port: \"\"\n    user: \"\"\n    db_name: \"\"\n    password: \"\"\n\n  # 1.2. Path to pg dump file\n  #  pg_dump_path: \"\"\n\n  # 2. Fill pg table section optionally\n\n  # 2.1. Config for pg tables\n  pg_table_config:\n\n    # Filters for each table by default\n    pg_column_filter:\n      # Generate all filters / Skip all filters\n      all_by_default: true # false\n\n      string: [ \"In\", \"NotIn\", \"Like\", \"NotLike\" ]\n      numeric: [ \"Lt\", \"LtOrEq\", \"Gt\", \"GtOrEq\", \"Eq\", \"NotEq\", \"In\", \"NotIn\" ]\n      # Bool always has one filter \"Where\"\n      # Time is the same as numeric\n\n      # Fill overrides for specific tables\n      overrides:\n\n        # Example for \"dummy\" table\n        dummy:\n          id: [ \"In\", \"NotIn\" ]\n          string: [ \"In\", \"Like\" ]\n          enum: [ \"In\", \"NotIn\" ] # Enum must be numeric column\n          # Bool always has one filter \"Where\"\n          numeric: [ \"Lt\", \"LtOrEq\", \"Gt\", \"GtOrEq\", \"Eq\", \"NotEq\", \"In\", \"NotIn\" ]\n          time: [ \"GtOrEq\", \"LtOrEq\" ] # Time is the same as numeric\n\n  # 3. Fill pg column types section optionally\n\n  # 3.1. Config for pg column types\n  pg_type_config:\n\n    # Example for \"citext\" column type\n    citext:\n      go_type: \"string\"\n      go_zero_type: \"zero.String\"\n"

  // StorageCustomModel ...
  StorageCustomModel = "// Code generated by Boiler; DO NOT EDIT.\n\npackage models\n\nimport (\n  {{- range .ModelPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\n{{.StructDescription}}\n"

  // StorageRocketLockModelMethods ...
  StorageRocketLockModelMethods = "// Code generated by Boiler; DO NOT EDIT.\n\npackage storage\n\nimport (\n  {{- range .ModelMethodsPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\nfunc (s *Storage) CreateRocketLock(ctx context.Context, input CreateRocketLockInput) (*models.RocketLock, error) {\n  query := `INSERT INTO rocket_locks(lock_id, locked_until)\n            VALUES (%s, NOW() + INTERVAL %s)\n            ON CONFLICT (lock_id)\n            DO UPDATE SET locked_until = NOW() + INTERVAL %[2]s\n            WHERE rocket_locks.locked_until < NOW()\n            RETURNING lock_id, locked_until`\n\n  builder := br.NewBuildedExpr(query, input.LockID, input.LockTTL)\n\n  model, err := pg.GetCtx[*models.RocketLock](ctx, s.client.Executor(ctx), builder)\n  if err != nil {\n    if errors.Is(err, pe.ErrModelNotFound) {\n      return nil, ErrRocketLockConflict\n    }\n    return nil, err\n  }\n\n  return model, nil\n}\n\nfunc (s *Storage) DeleteRocketLock(ctx context.Context, input DeleteRocketLockInput) (*models.RocketLock, error) {\n  query := `DELETE FROM rocket_locks\n            WHERE lock_id = %s\n            RETURNING lock_id, locked_until`\n\n  builder := br.NewBuildedExpr(query, input.LockID)\n\n  model, err := pg.GetCtx[*models.RocketLock](ctx, s.client.Executor(ctx), builder)\n  if err != nil {\n    if errors.Is(err, pe.ErrModelNotFound) {\n      return nil, ErrRocketLockNotFound\n    }\n    return nil, err\n  }\n\n  return model, nil\n}\n\nfunc (s *Storage) WithRocketLock(ctx context.Context, input WithRocketLockInput, fn func(ctx context.Context) error) error {\n  if _, err := s.CreateRocketLock(ctx, CreateRocketLockInput{\n    LockID:  input.LockID,\n    LockTTL: input.LockTTL,\n  }); err != nil {\n    return fmt.Errorf(\"s.CreateLock: %w\", err)\n  }\n\n  defer func() {\n    if _, err := s.DeleteRocketLock(ctx, DeleteRocketLockInput{\n      LockID: input.LockID,\n    }); err != nil {\n      log.Errorf(\"rocketLocksStorage: WithLock: s.DeleteLock: %v\", err)\n    }\n  }()\n\n  return fn(ctx)\n}\n"
  // StorageRocketLockModelOptions ...
  StorageRocketLockModelOptions = "// Code generated by Boiler; DO NOT EDIT.\n\npackage storage\n\nimport (\n  {{- range .ModelOptionsPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\nvar (\n  ErrRocketLockConflict = errors.New(\"rocket lock conflict\")\n  ErrRocketLockNotFound = errors.New(\"rocked lock not found\")\n)\n\ntype CreateRocketLockInput struct {\n  LockID  string\n  LockTTL time.Duration\n}\n\ntype DeleteRocketLockInput struct {\n  LockID string\n}\n\ntype WithRocketLockInput struct {\n  LockID  string\n  LockTTL time.Duration\n}\n"
  // StorageRocketLockMigration ...
  StorageRocketLockMigration = "-- +goose Up\n-- +goose StatementBegin\n\n-- migration created by Boiler; DO NOT EDIT.\nCREATE TABLE rocket_locks(\n    lock_id      TEXT NOT NULL PRIMARY KEY,\n    locked_until TIMESTAMP WITHOUT TIME ZONE\n);\n\n-- +goose StatementEnd\n\n-- +goose Down\n-- +goose StatementBegin\n\n-- migration created by Boiler; DO NOT EDIT.\nDROP TABLE rocket_locks;\n\n-- +goose StatementEnd\n"
  // StorageRocketLockModel ...
  StorageRocketLockModel = "type RocketLock struct {\n  LockID      string    `db:\"lock_id\"`\n  LockedUntil time.Time `db:\"locked_until\"`\n}\n"
)

// Gqlgen Generator compiled templates
const (
  // GqlgenConfig const for compiled Boiler build with gqlgen config file
  GqlgenConfig = "# Code generated by Boiler; YOU MAY CHANGE THIS.\n# Config for 99designs/gqlgen.\n\n# defined schema\nschema:\n  - api/graphql/**/*.graphql\n\n# generated code\nexec:\n  layout: follow-schema\n  dir: internal/app/graph/generated\n  package: generated\n\n# generated models\nmodel:\n  filename: internal/pkg/model/generated.go\n  package: model\n\n# generated resolvers\nresolver:\n  layout: follow-schema\n  dir: internal/app/graph\n  package: graph\n  filename_template: \"{name}.resolvers.go\"\n\n# models binding\nmodels:\n  # ID graphql type resolver\n  ID:\n    model:\n      - github.com/99designs/gqlgen/graphql.ID\n      - github.com/99designs/gqlgen/graphql.Int\n      - github.com/99designs/gqlgen/graphql.Int32\n      - github.com/99designs/gqlgen/graphql.Int64\n\n  # Int graphql type resolver\n  Int:\n    model:\n      - github.com/99designs/gqlgen/graphql.Int\n      - github.com/99designs/gqlgen/graphql.Int32\n      - github.com/99designs/gqlgen/graphql.Int64\n\n  # Int32 graphql type resolver\n  Int32:\n    model:\n      - github.com/99designs/gqlgen/graphql.Int32\n\n  # Int64 graphql type resolver\n  Int64:\n    model:\n      - github.com/99designs/gqlgen/graphql.Int64\n      \n  # UUID graphql type resolver\n  UUID:\n    model:\n      - github.com/99designs/gqlgen/graphql.UUID"
  // GqlgenTools const for compiled Boiler build with gqlgen packages
  GqlgenTools = "//go:build tools\n// +build tools\n\npackage tools\n\nimport (\n\t_ \"github.com/99designs/gqlgen\"\n\t_ \"github.com/99designs/gqlgen/graphql/introspection\"\n)\n"
  // GqlgenSandbox const for compiled Boiler build with apollo graphql sandbox
  GqlgenSandbox = "<head>\n  <meta charset=\"utf-8\">\n  <title>{{.Title}}</title>\n</head>\n<body>\n    <div style=\"width: 100%; height: 100%;\" id='embedded-sandbox'></div>\n    <script src=\"https://embeddable-sandbox.cdn.apollographql.com/_latest/embeddable-sandbox.umd.production.min.js\"></script>\n    <script>\n      new window.EmbeddedSandbox({\n        target: '#embedded-sandbox',\n        initialEndpoint: '{{.InitialEndpoint}}',\n        includeCookies: true,\n      });\n    </script>\n</body>"
  // GqlgenService const for compiled Boiler build with gqlgen service
  GqlgenService = "// Code generated by Boiler; YOU MUST CHANGE THIS.\n\npackage graph\n\nimport (\n  {{- range .ServicePackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\ntype Implementation struct {\n  *Resolver\n}\n\nfunc NewService() *Implementation {\n  return &Implementation{\n    Resolver: &Resolver{},\n  }\n}\n\n// RegisterService Code generated by Boiler; DO NOT EDIT.\nfunc (s *Implementation) RegisterService(params *app.RegisterParams) error {\n  params.SetServiceType(app.GqlgenServiceTyp)\n\n  schema := generated.NewExecutableSchema(generated.Config{\n    Resolvers: s.Resolver,\n  })\n  params.Gqlgen().SetGqlgenSchema(schema)\n  \n  return nil \n}\n"
  // GqlgenRegisterService const for compiled Boiler build with gqlgen service register method
  GqlgenRegisterService = "// RegisterService Code generated by Boiler; DO NOT EDIT.\nfunc (s *Implementation) RegisterService(params *app.RegisterParams) error {\n  params.SetServiceType(app.GqlgenServiceTyp)\n\n  schema := generated.NewExecutableSchema(generated.Config{\n    Resolvers: s.Resolver,\n  })\n  params.Gqlgen().SetGqlgenSchema(schema)\n  \n  return nil \n}\n"
)

// Gqlgen Generator Makefile Targets templates
const (
  // GqlgenMakeMkBinDeps ...
  GqlgenMakeMkBinDeps = "# Target generated by Boiler; DO NOT EDIT.\n.PHONY: bin-deps-gqlgen\nbin-deps-gqlgen:\n\tmkdir -p bin\n\tGOBIN=${PWD}/bin go install github.com/99designs/gqlgen@latest\n"
  // GqlgenMakeMkGenerate ...
  GqlgenMakeMkGenerate = "# Target generated by Boiler; DO NOT EDIT.\n.PHONY: generate-gqlgen\ngenerate-gqlgen:\n\t${PWD}/bin/gqlgen generate --config=.config/gqlgen_config.yaml\n"
)

// Gqlgen Generator Makefile Targets names
const (
  // GqlgenMakeMkBinDepsName ...
  GqlgenMakeMkBinDepsName = "bin-deps-gqlgen"
  // GqlgenMakeMkGenerateName ...
  GqlgenMakeMkGenerateName = "generate-gqlgen"
)

// Gqlgen Generator Graphql Schema compiled templates names
const (
  // NameGqlgenSchema ...
  NameGqlgenSchema = "gqlgen_schema"
  // NameGqlgenMutation ...
  NameGqlgenMutation = "gqlgen_mutation"
  // NameGqlgenQuery ...
  NameGqlgenQuery = "gqlgen_query"
  // NameGqlgenTypes ...
  NameGqlgenTypes = "gqlgen_types"
  // NameGqlgenEnums ...
  NameGqlgenEnums = "gqlgen_enums"
  // NameGqlgenScalars ...
  NameGqlgenScalars = "gqlgen_scalars"
)

// Gqlgen Generator Graphql Schema compiled templates
const (
  // GqlgenSchema ...
  GqlgenSchema = "\"\"\" Query \"\"\"\ntype Query\n\n\"\"\" Mutation \"\"\"\ntype Mutation\n"
  // GqlgenMutation ...
  GqlgenMutation = "extend type Mutation {\n    \"\"\"dummyMutation \"\"\"\n    dummyMutation(input: DummyMutationInput!): DummyMutationPayload!\n}\n\n\"\"\" DummyMutationInput \"\"\"\ninput DummyMutationInput {\n    \"\"\" field \"\"\"\n    field: String!\n}\n\n\"\"\" DummyMutationPayload \"\"\"\ntype DummyMutationPayload {\n    \"\"\" dummy \"\"\"\n    dummy: Dummy!\n}\n"
  // GqlgenQuery ...
  GqlgenQuery = "extend type Query {\n    \"\"\"dummyQuery \"\"\"\n    dummyQuery(input: DummyQueryInput!): DummyQueryPayload!\n}\n\n\"\"\" DummyQueryInput \"\"\"\ninput DummyQueryInput {\n    \"\"\" id \"\"\"\n    id: ID!\n}\n\n\"\"\" DummyQueryPayload \"\"\"\ntype DummyQueryPayload {\n    \"\"\" dummy \"\"\"\n    dummy: Dummy!\n}\n"
  // GqlgenTypes ...
  GqlgenTypes = "\"\"\" Dummy \"\"\"\ntype Dummy {\n    \"\"\" id \"\"\"\n    id: ID!\n    \"\"\" field \"\"\"\n    field: String!\n}"
  // GqlgenEnums ...
  GqlgenEnums = "\"\"\" DummyEnum \"\"\"\nenum DummyEnum {\n    \"\"\" ENUM_VALUE \"\"\"\n    ENUM_VALUE\n}"
  // GqlgenScalars ...
  GqlgenScalars = "\"\"\" Int32 \"\"\"\nscalar Int32\n\n\"\"\" Int64 \"\"\"\nscalar Int64\n\n\"\"\" UUID \"\"\"\nscalar UUID\n"
)

// Grpc Generator compiled templates
const (
  // GrpcService const for compiled Boiler build with grpc service implementation
  GrpcService = "// Code generated by Boiler; YOU MUST CHANGE THIS.\n\npackage {{toSnakeCase .ServiceName}}\n\nimport (\n  {{- range .ServicePackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\ntype Implementation struct {\n  desc.Unimplemented{{.ServiceName}}Server\n}\n\nfunc New{{.ServiceName}}Service() *Implementation {\n  return &Implementation{}\n}\n\n// RegisterService Code generated by Boiler; DO NOT EDIT.\nfunc (s *Implementation) RegisterService(params *app.RegisterParams) error {\n  params.SetServiceType(app.GrpcServiceTyp)\n  grpcParams := params.Grpc()\n  desc.Register{{.ServiceName}}Server(grpcParams.GrpcServiceRegistrar(), s)\n\n  return desc.Register{{.ServiceName}}HandlerFromEndpoint(\n    params.Context(),\n    grpcParams.GrpcHttpProxyServeMux(),\n    grpcParams.GrpcServerEndpoint(),\n    grpcParams.GrpcClientOptions(),\n  )\n}\n"
  // GrpcRegisterService const for compiled Boiler build with grpc service register method
  GrpcRegisterService = "// RegisterService Code generated by Boiler; DO NOT EDIT.\nfunc (s *Implementation) RegisterService(params *app.RegisterParams) error {\n  params.SetServiceType(app.GrpcServiceTyp)\n  grpcParams := params.Grpc()\n  desc.Register{{.ServiceName}}Server(grpcParams.GrpcServiceRegistrar(), s)\n\n  return desc.Register{{.ServiceName}}HandlerFromEndpoint(\n    params.Context(),\n    grpcParams.GrpcHttpProxyServeMux(),\n    grpcParams.GrpcServerEndpoint(),\n    grpcParams.GrpcClientOptions(),\n  )\n}\n"
  // GrpcStub const for compiled Boiler build with grpc stub implementation
  GrpcStub = "// Code generated by Boiler; YOU MUST CHANGE THIS.\n\npackage {{toSnakeCase .ServiceName}}\n\nimport (\n  {{- range .CallStubPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\n// {{.CallName}} implementation stub. Change this.\nfunc (s *Implementation) {{.CallName}}(ctx context.Context, req *desc.{{.CallInputProto}}) (*desc.{{.CallOutputProto}}, error) {\n  return nil, status.Error(codes.Unimplemented, \"{{.CallName}} not implemented\")\n}\n"
  // GrpcProto const for compiled Boiler build with proto template for grpc service
  GrpcProto = "// Code generated by Boiler; YOU MUST CHANGE THIS.\n\nsyntax = \"proto3\";\n\npackage {{.serviceName}};\n\noption go_package = \"{{.goPackage}}\";\n\n// Use proto imports:\n//   Execute:   boiler proto-deps init & boiler proto-deps gen --github-token=<GITHUB_TOKEN>;\n//   Uncomment: imports declarations.\n\n// import \"proto/ushakovn-org/protobuf/validate/validate.proto\";\n// import \"proto/ushakovn-org/protobuf/protobuf/timestamp.proto\";\n// import \"proto/ushakovn-org/protobuf/protobuf/duration.proto\";\n// import \"proto/ushakovn-org/protobuf/api/annotations.proto\";\n\nservice DummyService {\n  rpc GetDummy(GetDummyRequest) returns (GetDummyResponse);\n}\n\nmessage GetDummyRequest {\n  string id = 1;\n}\n\nmessage GetDummyResponse {\n  repeated Dummy dummy = 1;\n}\n\nmessage Dummy {\n  int64 id = 1;\n  string field = 2;\n}\n"
)

// Grpc Generator Makefile Targets templates
const (
  // GrpcMakeMkBinDeps ...
  GrpcMakeMkBinDeps = "# Target generated by Boiler; DO NOT EDIT.\n.PHONY: bin-deps-protoc\nbin-deps-protoc:\n\tmkdir -p bin\n\tGOBIN=${PWD}/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@latest\n\tGOBIN=${PWD}/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest\n\tGOBIN=${PWD}/bin go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest\n"
  // GrpcMakeMkGenerate ...
  GrpcMakeMkGenerate = "# Target generated by Boiler. DO NOT EDIT.\n.PHONY: generate-protoc\ngenerate-protoc:\n\tprotoc \\\n\t--plugin=protoc-gen-go=${PWD}/bin/protoc-gen-go \\\n\t--plugin=protoc-gen-go-gprc=${PWD}/bin/protoc-gen-go-grpc \\\n\t--plugin=protoc-gen-grpc-gateway=${PWD}/bin/protoc-gen-grpc-gateway \\\n\t--plugin=protoc-gen-openapiv2=${PWD}/bin/protoc-gen-openapiv2 \\\n\t\\\n\t--go_out=. --go_opt=paths=import --go_opt=module={{.goPackageTrim}} \\\n\t--go-grpc_out=. --go-grpc_opt=paths=import --go-grpc_opt=module={{.goPackageTrim}} \\\n\t--grpc-gateway_out=. --grpc-gateway_opt=paths=import --grpc-gateway_opt=module={{.goPackageTrim}} \\\n\t--grpc-gateway_opt generate_unbound_methods=true \\\n\t--validate_out=\"lang=go,module={{.goPackageTrim}},paths=import:.\" \\\n\t\\\n\t./api/**/*.proto"
)

// Grpc Generator Makefile Targets names
const (
  // GrpcMakeMkBinDepsName ...
  GrpcMakeMkBinDepsName = "bin-deps-protoc"
  // GrpcMakeMkGenerateName ...
  GrpcMakeMkGenerateName = "generate-protoc"
)

// Proto Dependencies Generator compiled templates
const (
  // ProtoDepsConfig const for compiled Boiler build with proto deps config file
  ProtoDepsConfig = "# Proto dependencies config generated by Boiler; YOU MAY CHANGE THIS.\n\n# App proto dependencies section; DO NOT EDIT.\napp_deps:\n    - import: github.com/ushakovn-org/protobuf/protoc-gen-validate/validate/validate.proto@main\n    - import: github.com/ushakovn-org/protobuf/google/protobuf/timestamp.proto@main\n    - import: github.com/ushakovn-org/protobuf/google/protobuf/duration.proto@main\n    - import: github.com/ushakovn-org/protobuf/google/api/annotations.proto@main\n\n# Local proto dependencies section\nlocal_deps:\n  # Example path:\n  # - path: .boiler/vendor/<owner>/<repo>/<path>.proto\n\n# External proto dependencies section\nexternal_deps:\n  # Example import:\n  # - import: github.com/<owner>/<repo>/<package>/<path>.proto\n"
  // ProtoDepsDump const for compiled Boiler build with proto deps dump file
  ProtoDepsDump = "# Proto dependencies dump generated by Boiler; DO NOT EDIT.\n\n# App proto dependencies\napp_deps:\n  {{- range .AppDeps}}\n  - import: {{.Import}}\n  {{- end}}\n\n# Local proto dependencies\nlocal_deps:\n  {{- range .LocalDeps}}\n  - path: {{.Path}}\n  {{- end}}\n\n# External proto dependencies\nexternal_deps:\n  {{- range .ExternalDeps}}\n  - import: {{.Import}}\n  {{- end}}\n"
  // ProtoDepsMakeMk const for compiled Boiler build with target for including make.mk file
  ProtoDepsMakeMk = "# Target generated by Boiler; DO NOT EDIT.\n.PHONY: bin-deps-proto\nbin-deps-proto:\n\tmkdir -p bin\n\tGOBIN=${PWD}/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@latest\n\tGOBIN=${PWD}/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest\n"
)

// Proto Dependencies Generator Makefile Targets names
const (
  // ProtoDepsMakeMkBinDepsName ...
  ProtoDepsMakeMkBinDepsName = "bin-deps-proto"
)

// Config Generator compiled templates
const (
  // GenConfigEmpty const for compiled Boiler build with empty config for generation
  GenConfigEmpty = "# Config generated by Boiler; YOU MUST CHANGE THIS.\nversion: \"1\"\n\n# Boiler app section\napp:\n  name: \"app\"\n  version: \"v0.0.1\"\n  description: \"test app description\"\n\n# Custom keys section\ncustom:\n  first_group_key:\n    group: \"first_group\"\n    type: \"int\"\n    value: \"1\"\n    description: \"test first group key\"\n\n  second_group_key:\n    group: \"second_group\"\n    type: \"duration\"\n    value: \"10s\"\n    description: \"test second group key\"\n"
  // GenConfigGroups const for compiled Boiler build with generated config groups
  GenConfigGroups = "// Code generated by Boiler; DO NOT EDIT.\n\npackage config\n\nimport (\n  {{- range .GroupsPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\nvar (\n  {{- range .ConfigGroups}}\n  {{toUpperCamelCase .GroupName}} = {{toLowerCamelCase .GroupName}}ConfigGroup{}\n  {{- end}}\n)\n\nvar (\n  {{- range .ConfigGroups}}\n  _ = {{toUpperCamelCase .GroupName}}\n  {{- end}}\n)\n\ntype (\n  {{- range .ConfigGroups}}\n  {{toLowerCamelCase .GroupName}}ConfigGroup struct{}\n  {{- end}}\n)\n\ntype (\n  {{- range .ConfigGroups}}\n  {{toLowerCamelCase .GroupName}}ConfigKey string\n  {{- end}}\n)\n\nfunc configValue(ctx context.Context, configKey string) config.Value {\n  return config.ClientConfig(ctx).GetValue(configKey)\n}\n"
  // GenConfigConfig const for compiled Boiler build with generated config
  GenConfigConfig = "// Code generated by Boiler; DO NOT EDIT.\n\npackage config\n\nimport (\n  {{- range .ConfigPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\n{{- range $cg := .ConfigGroups}}\nconst (\n  {{- range .GroupKeys}}\n  // {{.KeyComment}}\n  {{toCapitalizeCase .KeyName}} {{toLowerCamelCase $cg.GroupName}}ConfigKey = \"{{toSnakeCase .KeyName}}\"\n  {{- end}}\n)\n{{end}}\n\n{{- range $cg := .ConfigGroups}}\n{{- range .GroupKeys}}\nfunc (c {{toLowerCamelCase $cg.GroupName}}ConfigGroup) {{toUpperCamelCase .KeyNameTrim}}(ctx context.Context) {{.ValueType}} {\n  return configValue(ctx, string({{toCapitalizeCase .KeyName}})).{{.ValueCall}}()\n}\n{{end}}\n{{- end}}\n\n// Suppress unused imports\nvar (\n  _ = time.Time{}\n)\n"
)
