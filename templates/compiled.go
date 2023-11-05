package templates

// RPC Generator compiled templates
const (
  // Contracts const for compiled Boiler build with contracts template
  Contracts = "// Code generated by Boiler; DO NOT EDIT.\n\npackage handler\n\n{{range .Requests}}\n// {{.Name}}Request ...\ntype {{.Name}}Request struct {\n  {{range .Fields}}\n  {{.Name}} {{.Type}} `json:\"{{.Tag}}\"`\n  {{- end}}\n}\n{{- end}}\n\n{{range .Responses}}\n// {{.Name}}Response ...\ntype {{.Name}}Response struct {\n  {{range .Fields}}\n  {{.Name}} {{.Type}} `json:\"{{.Tag}}\"`\n  {{- end}}\n}\n{{- end}}\n\n{{range .TypeDefs}}\n// {{.Name}} ...\ntype {{.Name}} struct{\n  {{range .Fields}}\n  {{.Name}} {{.Type}} `json:\"{{.Tag}}\"`\n  {{- end}}\n}\n{{- end}}"
  // Handle const for compiled Boiler build with handle template
  Handle = "// Code generated by Boiler; YOU MAY CHANGE THIS.\n\npackage handler\n\nimport (\n\t\"net/http\"\n\n\t\"github.com/gin-gonic/gin\"\n)\n\n// Handle{{.Name}} stub ...\nfunc (h *Handler) Handle{{.Name}}(ctx *gin.Context) {\n  req := &{{.Name}}Request{}\n\n  if err := ctx.BindJSON(req); err != nil {\n    ctx.JSON(http.StatusBadRequest, err)\n    return\n  }\n\n  resp := &{{.Name}}Response{}\n  ctx.JSON(http.StatusOK, resp)\n}\n"
  // Handler const for compiled Boiler build with handler template
  Handler = "// Code generated by Boiler; DO NOT EDIT.\n\npackage handler\n\nimport (\n  \"github.com/gin-gonic/gin\"\n  log \"github.com/sirupsen/logrus\"\n)\n\n// Handler ...\ntype Handler struct {\n  addr string\n  g    *gin.Engine\n}\n\n// NewHandler ...\nfunc NewHandler() *Handler {\n  return &Handler{}\n}\n\ntype (\n  route  string\n  Routes map[route]gin.IRoutes\n)\n\nconst (\n    {{- range .Handles}}\n\t{{.Name}}Route route = \"{{.Route}}\"\n\t{{- end}}\n)\n\n// RegisterRoutes ...\nfunc (h *Handler) RegisterRoutes() Routes {\n  h.g = gin.New()\n  // Registered routes\n  return Routes{\n    {{- range .Handles}}\n    {{.Name}}Route: h.g.POST(string({{.Name}}Route), h.Handle{{.Name}}),\n    {{- end}}\n  }\n}\n\n// ServeHTTP ...\nfunc (h *Handler) ServeHTTP() {\n  go func() {\n    if err := h.g.Run(h.addr); err != nil {\n      log.Fatalf(\"ServeHTTP error: gin.Run: %v\", err)\n    }\n  }()\n}\n"
)

// Project Generator compiled templates
const (
  // Main const for compiled Boiler build with main template
  Main = "// Code generated by Boiler. YOU MAY CHANGE THIS\npackage main\n\nfunc main() {\n    println(`\n        _           _ _\n       | |         (_) |\n       | |__   ___  _| | ___ _ __\n       | '_ \\ / _ \\| | |/ _ \\ '__|\n       | |_) | (_) | | |  __/ |\n       |_.__/ \\___/|_|_|\\___|_|\n    `)\n}"
  // Gomod const for compiled Boiler build with go mod template
  Gomod = "module main\n\ngo 1.19"
  // Makefile const for compiled Boiler build with makefile template
  Makefile = "# Including generated by Boiler. DO NOT EDIT.\ninclude make.mk\n\n\n# Target section. YOU MAY CHANGE THIS.\nhelp:\n\tboiler --help\n"
)

// Project Generator compiled templates names
const (
  // NameMain name for compiled Main file template
  NameMain = "main"
  // NameGomod name for compiled Gomod file template
  NameGomod = "gomod"
)

// Storage Generator compiled templates
const (
  // FilterIfStmtWithPtr const for compiled Boiler build with filter if statement for zero typed filters
  FilterIfStmtWithPtr = "filters.%s.Ptr() != nil"
  // FilterIfStmtWithLen const for compiled Boiler build with filter if statement for slice typed filters
  FilterIfStmtWithLen = "len(filters.%s) > 0"

  InputIfStmtWithPtr = "input.%s.Ptr() != nil"
  InputIfStmtWithLen = "len(input.%s) > 0"

  Builders       = "// Code generated by Boiler; DO NOT EDIT.\n\npackage storage\n\nimport (\n  {{- range .BuildersPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\nfunc NewSelectBuilder() sq.SelectBuilder {\n  return sq.SelectBuilder{}.PlaceholderFormat(sq.Dollar)\n}\n\nfunc NewInsertBuilder() sq.InsertBuilder {\n  return sq.InsertBuilder{}.PlaceholderFormat(sq.Dollar)\n}\n\nfunc NewUpdateBuilder() sq.UpdateBuilder {\n  return sq.UpdateBuilder{}.PlaceholderFormat(sq.Dollar)\n}\n\nfunc NewDeleteBuilder() sq.DeleteBuilder {\n  return sq.DeleteBuilder{}.PlaceholderFormat(sq.Dollar)\n}\n"
  Client         = "// Code generated by Boiler; DO NOT EDIT.\n\npackage storage\n\nimport (\n  {{- range .ClientPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\nvar ErrZeroRowsRetrieved = errors.New(\"zero rows retrieved\")\n\ntype Client interface {\n  Pinger\n  Execer\n  Querier\n  Txer\n}\n\ntype Querier interface {\n  QueryContext(context.Context, string, ...any) (*sql.Rows, error)\n  QueryRowContext(context.Context, string, ...any) *sql.Row\n}\n\ntype Execer interface {\n  ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)\n}\n\ntype Txer interface {\n  Begin() (*sql.Tx, error)\n  BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)\n}\n\ntype Pinger interface {\n  PingContext(ctx context.Context) error\n}\n\ntype Builder interface {\n  ToSql() (statement string, args []any, err error)\n}\n\nfunc DoQueryContext[Model any](ctx context.Context, querier Querier, builder Builder) ([]Model, error) {\n  statement, args, err := builder.ToSql()\n  if err != nil {\n    return nil, fmt.Errorf(\"builder.ToSql: %w\", err)\n  }\n  var models []Model\n  if err = sc.Select(ctx, querier, &models, statement, args...); err != nil {\n    return nil, fmt.Errorf(\"sqlscan.Select: %w\", err)\n  }\n  if len(models) == 0 {\n    return nil, ErrZeroRowsRetrieved\n  }\n  return models, nil\n}\n\nfunc DoExecContext(ctx context.Context, execer Execer, builder Builder) error  {\n  statement, args, err := builder.ToSql()\n  if err != nil {\n    return fmt.Errorf(\"builder.ToSql: %w\", err)\n  }\n  if _, err = execer.ExecContext(ctx, statement, args...); err != nil {\n    return fmt.Errorf(\"execer.ExecContext: %w\", err)\n  }\n  return nil\n}\n\n\n\n\n"
  Consts         = "// Code generated by Boiler; DO NOT EDIT.\n\npackage storage\n\nimport (\n  {{- range .ConstsPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\ntype tableName string\n\nconst (\n  {{- range .Models}}\n  {{.ModelName}}_TableName = \"{{.SqlTableName}}\"\n  {{- end}}\n)\n\ntype (\n  {{- range .Models}}\n  {{toLowerCamelCase .ModelName}}_Field string\n  {{- end}}\n)\n\n{{- range $m := .Models}}\nconst (\n  {{- range .ModelFields}}\n  {{.FieldName}}_{{$m.ModelName}}_Field {{toLowerCamelCase $m.ModelName}}_Field = \"{{.SqlTableFieldName}}\"\n  {{- end}}\n)\n{{end}}\n"
  Implementation = "// Code generated by Boiler; DO NOT EDIT. {{$modelName := .ModelName}}\n\npackage storage\n\nimport (\n  {{- range .ImplementationPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\ntype {{toLowerCamelCase $modelName}}Storage struct {\n  client client.Client\n}\n\nfunc New{{$modelName}}Storage(client client.Client) {{$modelName}}Storage {\n  return &{{toLowerCamelCase $modelName}}Storage{\n    client: client,\n  }\n}\n\nfunc (s *{{toLowerCamelCase $modelName}}Storage) List{{$modelName}}(ctx context.Context, input List{{$modelName}}Input) ([]*models.{{$modelName}}, error) {\n  builder := client.NewSelectBuilder().\n    Columns(\n      {{- range .ModelFields}}\n      string({{.FieldName}}_{{$modelName}}_Field),\n      {{- end}}\n    ).\n    From({{$modelName}}_TableName)\n\n  if err := input.Pagination.validate(); err != nil {\n    return nil, fmt.Errorf(\"pagination.Validate: %w\", err)\n  }\n\n  offset, limit := input.Pagination.orDefault().toOffsetLimit()\n  builder = builder.Offset(offset).Limit(limit)\n\n  if input.Sort != nil {\n    builder = builder.OrderBy(input.Sort.{{toLowerCamelCase $modelName}}Sort())\n  }\n\n  if filters := input.Filters; filters != nil {\n    {{- range $modelField := .ModelFields}}\n    {{- range .ModelsFieldFilters}}\n    if {{.FilterIfStmt}} {\n      builder = builder.Where(sq.{{.FilterSqOperator}}{string({{$modelField.FieldName}}_{{$modelName}}_Field): filters.{{.FilterName}}{{.FilterTypeSuffix}}})\n    }\n    {{- end}}\n    {{- end}}\n  }\n\n  return client.DoQueryContext[*models.{{$modelName}}](ctx, s.client, builder)\n}\n\nfunc (s *{{toLowerCamelCase $modelName}}Storage) {{$modelName}}(ctx context.Context, input {{$modelName}}Input) (*models.{{$modelName}}, error) {\n  builder := client.NewSelectBuilder().\n    Columns(\n      {{- range .ModelFields}}\n      string({{.FieldName}}_{{$modelName}}_Field),\n      {{- end}}\n    ).\n    From({{$modelName}}_TableName)\n\n  if filters := input.Filters; filters != nil {\n    {{- range $modelField := .ModelFields}}\n    {{- range .ModelFieldFilters}}\n    if {{.FilterIfStmt}} {\n      builder = builder.Where(sq.Eq{string({{$modelField.FieldName}}_{{$modelName}}_Field): filters.{{.FilterName}}{{.FilterTypeSuffix}}})\n    }\n    {{- end}}\n    {{- end}}\n  }\n\n  model, err := client.DoQueryContext[*models.{{$modelName}}](ctx, s.client, builder)\n  if err != nil {\n    return nil, err\n  }\n\n  return model[0], nil\n}\n\nfunc (s *{{toLowerCamelCase $modelName}}Storage) Put{{$modelName}}(ctx context.Context, input Put{{$modelName}}Input) (*models.{{$modelName}}, error) {\n  builder := client.NewInsertBuilder().\n    Into({{$modelName}}_TableName)\n\n  fields := map[string]any{\n    {{- range .ModelFields}}\n    {{- if eq .NotNullField true}}\n    string({{.FieldName}}_{{$modelName}}_Field): input.{{.FieldName}},\n    {{- end}}\n    {{- end}}\n  }\n\n  {{- range $modelField := .ModelFields}}\n  {{- if eq .NotNullField false}}\n  if {{.FieldIfStmt}} {\n    fields[string({{$modelField.FieldName}}_{{$modelName}}_Field)] = input.{{.FieldName}}{{.FieldTypeSuffix}}\n  }\n  {{- end}}\n  {{- end}}\n\n  builder = builder.SetMap(fields)\n\n  if err := client.DoExecContext(ctx, s.client, builder); err != nil {\n    return nil, err\n  }\n\n  model := &models.{{$modelName}}{\n    {{- range .ModelFields}}\n    {{.FieldName}}: input.{{.FieldName}},\n    {{- end}}\n  }\n\n  return model, nil\n}\n\nfunc (s *{{toLowerCamelCase $modelName}}Storage) Update{{$modelName}}(ctx context.Context, input Update{{$modelName}}Input) (*models.{{$modelName}}, error) {\n  builder := client.NewUpdateBuilder().\n    Table({{$modelName}}_TableName)\n\n  fields := map[string]any{\n    {{- range .ModelFields}}\n    {{- if eq .NotNullField true}}\n    {{- if ne .FieldBadge \"pk\"}}\n    string({{.FieldName}}_{{$modelName}}_Field): input.{{.FieldName}},\n    {{- end}}\n    {{- end}}\n    {{- end}}\n  }\n\n  {{- range .ModelFields}}\n  {{- if eq .NotNullField false}}\n  {{- if ne .FieldBadge \"pk\"}}\n  if {{.FieldIfStmt}} {\n    fields[string({{.FieldName}}_{{$modelName}}_Field)] = input.{{.FieldName}}{{.FieldTypeSuffix}}\n  }\n  {{- end}}\n  {{- end}}\n  {{- end}}\n\n  builder = builder.\n    SetMap(fields).\n    {{- range .ModelFields}}\n    {{- if eq .FieldBadge \"pk\"}}\n    Where(sq.Eq{string({{.FieldName}}_{{$modelName}}_Field): input.{{.FieldName}}}).\n    Suffix(suffixReturning)\n    {{- end}}\n    {{- end}}\n\n  model, err := client.DoQueryContext[*models.{{$modelName}}](ctx, s.client, builder)\n  if err != nil {\n    return nil, err\n  }\n  return model[0], nil\n}\n\nfunc (s *{{toLowerCamelCase $modelName}}Storage) Delete{{$modelName}}(ctx context.Context, input Delete{{$modelName}}Input) (*models.{{.ModelName}}, error) {\n  builder := client.NewDeleteBuilder().\n    From({{$modelName}}_TableName).\n    {{- range .ModelFields}}\n    {{- if eq .FieldBadge \"pk\"}}\n    Where(sq.Eq{string({{.FieldName}}_{{$modelName}}_Field): input.{{.FieldName}}}).\n    Suffix(suffixReturning)\n    {{- end}}\n    {{- end}}\n\n  model, err := client.DoQueryContext[*models.{{$modelName}}](ctx, s.client, builder)\n  if err != nil {\n    return nil, err\n  }\n  return model[0], nil\n}\n"
  Interface      = "// Code generated by Boiler; DO NOT EDIT.\n\npackage storage\n\nimport (\n  {{- range .InterfacePackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\ntype {{.ModelName}}Storage interface {\n  List{{.ModelName}}(ctx context.Context, input List{{.ModelName}}Input) ([]*models.{{.ModelName}}, error)\n  {{.ModelName}}(ctx context.Context, input {{.ModelName}}Input) (*models.{{.ModelName}}, error)\n  Put{{.ModelName}}(ctx context.Context, input Put{{.ModelName}}Input) (*models.{{.ModelName}}, error)\n  Update{{.ModelName}}(ctx context.Context, input Update{{.ModelName}}Input) (*models.{{.ModelName}}, error)\n  Delete{{.ModelName}}(ctx context.Context, input Delete{{.ModelName}}Input) (*models.{{.ModelName}}, error)\n}\n\ntype {{.ModelName}}Input struct {\n  Filters *{{.ModelName}}Filter\n}\n\ntype {{.ModelName}}Filter struct {\n  {{- range .ModelFields}}\n  {{- range .ModelFieldFilters}}\n  {{.FilterName}} {{.FilterType}}\n  {{- end}}\n  {{- end}}\n}\n\ntype Put{{.ModelName}}Input struct {\n  {{- range .ModelFields}}\n  {{.FieldName}} {{.FieldType}} {{- if eq .FieldBadge \"pk\"}} // PRIMARY KEY{{- end}}\n  {{- end}}\n}\n\ntype Delete{{.ModelName}}Input struct {\n  {{- range .ModelFields}}\n  {{- if eq .FieldBadge \"pk\"}}\n  {{.FieldName}} {{.FieldType}} // PRIMARY KEY\n  {{- end}}\n  {{- end}}\n}\n\ntype Update{{.ModelName}}Input struct {\n  {{- range .ModelFields}}\n  {{.FieldName}} {{.FieldType}} {{- if eq .FieldBadge \"pk\"}} // PRIMARY KEY{{- end}}\n  {{- end}}\n}\n\ntype List{{.ModelName}}Input struct {\n  Filters    *List{{.ModelName}}Filters\n  Sort       *List{{.ModelName}}Sort\n  Pagination *Pagination\n}\n\ntype List{{.ModelName}}Sort struct {\n  Field {{toLowerCamelCase .ModelName}}_Field\n  Order sortOrder\n}\n\nfunc (p *List{{.ModelName}}Sort) {{toLowerCamelCase .ModelName}}Sort() string {\n  return fmt.Sprintf(\"%s %s\", p.Field, p.Order)\n}\n\ntype List{{.ModelName}}Filters struct {\n  {{- range .ModelFields}}\n  {{- range .ModelsFieldFilters}}\n  {{.FilterName}} {{.FilterType}}\n  {{- end}}\n  {{- end}}\n}\n"
  Models         = "// Code generated by Boiler; DO NOT EDIT.\n\npackage models\n\nimport (\n  {{- range .ModelsPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\n{{range .Models}}\ntype {{.ModelName}} struct {\n  {{- range .ModelFields}}\n  {{.FieldName}} {{.FieldType}} `db:\"{{.SqlTableFieldName}}\"` {{- if eq .FieldBadge \"pk\"}} // PRIMARY KEY{{- end}}\n  {{- end}}\n}\n{{end}}"
  Options        = "// Code generated by Boiler; DO NOT EDIT.\n\npackage storage\n\nimport (\n  {{- range .OptionsPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\nconst suffixReturning = \"RETURNING *\"\n\ntype sortOrder string\n\nconst (\n  SortOrderAsc  sortOrder = \"ASC\"\n  SortOrderDesc sortOrder = \"DESC\"\n  SortOrderRand sortOrder = \"RAND()\"\n)\n\ntype Pagination struct {\n  Page    uint64\n  PerPage uint64\n}\n\nfunc (p *Pagination) orDefault() *Pagination {\n  if p != nil {\n    return p\n  }\n  return &Pagination{\n    Page:    0,\n    PerPage: 100,\n  }\n}\n\nfunc (p *Pagination) validate() error {\n  if p == nil {\n    return nil\n  }\n  if p.Page < 0 {\n    return fmt.Errorf(\"pagination.Page=%d must be non-negative\", p.Page)\n  }\n  if p.PerPage < 0 {\n    return fmt.Errorf(\"pagination.PerPage=%d must be positive\", p.PerPage)\n  }\n  return nil\n}\n\nfunc (p *Pagination) toOffsetLimit() (offset uint64, limit uint64) {\n  offset = p.Page * p.PerPage\n  limit = p.PerPage\n  return\n}\n\n"
)

// Gqlgen Generator compiled templates
const (
  // GqlgenYaml const for compiled Boiler build with gqlgen yaml config
  GqlgenYaml = "# Code generated by Boiler; YOU MAY CHANGE THIS.\n# Config for 99designs/gqlgen.\n\n# defined schema\nschema:\n  - api/graphql/**/*.graphql\n\n# generated code\nexec:\n  layout: follow-schema\n  dir: internal/app/graph/generated\n  package: generated\n\n# generated models\nmodel:\n  filename: internal/pkg/model/generated.go\n  package: model\n\n# generated resolvers\nresolver:\n  layout: follow-schema\n  dir: internal/app/graph\n  package: graph\n  filename_template: \"{name}.resolvers.go\"\n\n# models binding\nmodels:\n  # ID graphql type resolver\n  ID:\n    model:\n      - github.com/99designs/gqlgen/graphql.ID\n      - github.com/99designs/gqlgen/graphql.Int\n      - github.com/99designs/gqlgen/graphql.Int64\n      - github.com/99designs/gqlgen/graphql.Int32\n\n  # Int graphql type resolver\n  Int:\n    model:\n      - github.com/99designs/gqlgen/graphql.Int\n      - github.com/99designs/gqlgen/graphql.Int64\n      - github.com/99designs/gqlgen/graphql.Int32\n\n  # UUID graphql type resolver\n  UUID:\n    model:\n      - github.com/99designs/gqlgen/graphql.UUID"
)

// Grpc Generator compiled templates
const (
  // GrpcService const for compiled Boiler build with grpc service implementation
  GrpcService = "// Code generated by Boiler; YOU MUST CHANGE THIS.\n\npackage {{toSnakeCase .ServiceName}}\n\nimport (\n  {{- range .ServicePackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\ntype Implementation struct {\n  desc.Unimplemented{{.ServiceName}}Server\n}\n\nfunc New{{.ServiceName}}() *Implementation {\n  return &Implementation{}\n}\n\nfunc (s *Implementation) Register(params *app.RegisterParams) {\n  desc.RegisterBarnieServer(params.GrpcServer, &Implementation{})\n}\n"
  // GrpcStub const for compiled Boiler build with grpc stub implementation
  GrpcStub = "// Code generated by Boiler; YOU MUST CHANGE THIS.\n\npackage {{toSnakeCase .ServiceName}}\n\nimport (\n  {{- range .CallStubPackages}}\n  {{.ImportAlias}} \"{{.ImportLine}}\"\n  {{- end}}\n)\n\n// {{.CallName}} implementation stub. Change this.\nfunc (s *Implementation) {{.CallName}}(ctx context.Context, req *desc.{{.CallInputProto}}) (*desc.{{.CallOutputProto}}, error) {\n  return nil, status.Error(codes.Unimplemented, \"{{.CallName}} not implemented\")\n}\n"
  // GrpcProto const for compiled Boiler build with proto template for grpc service
  GrpcProto = "// Code generated by Boiler; YOU MUST CHANGE THIS.\n\nsyntax = \"proto3\";\n\npackage {{.serviceName}};\n\noption go_package = \"{{.goPackage}}\";"
  // GrpcMakeMk const for compiled Boiler build with target for including make.mk file
  GrpcMakeMk = "# Target generated by Boiler. DO NOT EDIT.\ngenerate-protoc:\n\tprotoc \\\n\t--go_out=. --go_opt=paths=import --go_opt=module={{.goPackageTrim}} \\\n\t--go-grpc_out=. --go-grpc_opt=paths=import --go-grpc_opt=module={{.goPackageTrim}} \\\n\t./api/**/*.proto"
)
