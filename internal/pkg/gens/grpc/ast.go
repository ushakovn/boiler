package grpc

import (
  "fmt"
  "go/ast"
  "go/parser"
  "go/printer"
  "go/token"
  "os"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/filer"
)

const (
  unsafePrefix = "Unsafe"
  serverSuffix = "Server"
)

type grpcServerInterface struct {
  ServerName  string
  ServerCalls []*grpcServerInterfaceCall
}

type grpcServerInterfaceCall struct {
  CallName    string
  InputProto  string
  OutputProto string
}

func buildAstFuncParams(callDesc *grpcServiceCallDesc) *ast.FieldList {
  return &ast.FieldList{
    List: []*ast.Field{
      {
        Names: []*ast.Ident{
          {
            Name: "ctx",
          },
        },
        Type: &ast.SelectorExpr{
          X: &ast.Ident{
            Name: "context",
          },
          Sel: &ast.Ident{
            Name: "Context",
          },
        },
      },
      {
        Names: []*ast.Ident{
          {
            Name: "req",
          },
        },
        Type: &ast.StarExpr{
          X: &ast.SelectorExpr{
            X: &ast.Ident{
              Name: "desc",
            },
            Sel: &ast.Ident{
              Name: callDesc.CallInputProto,
            },
          },
        },
      },
    },
  }
}

func buildAstFuncResults(callDesc *grpcServiceCallDesc) *ast.FieldList {
  return &ast.FieldList{
    List: []*ast.Field{
      {
        Type: &ast.StarExpr{
          X: &ast.SelectorExpr{
            X: &ast.Ident{
              Name: "desc",
            },
            Sel: &ast.Ident{
              Name: callDesc.CallOutputProto,
            },
          },
        },
      },
      {
        Type: &ast.Ident{
          Name: "error",
        },
      },
    },
  }
}

func scanGrpcServerInterface(filePath string) (*grpcServerInterface, error) {
  if err := validateGrpcProtocFileName(filePath); err != nil {
    return nil, fmt.Errorf("validateGrpcProtocFileName: %w", err)
  }
  fileSet := token.NewFileSet()

  astFile, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
  if err != nil {
    return nil, fmt.Errorf("parser.ParseFile: %w", err)
  }

  var grpcServer *grpcServerInterface

  ast.Inspect(astFile, func(node ast.Node) bool {
    if grpcServer != nil {
      return false
    }
    typeSpec, ok := node.(*ast.TypeSpec)
    if !ok {
      return true
    }
    if !typeSpec.Name.IsExported() {
      return true
    }
    interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
    if !ok {
      return true
    }
    interfaceTypeIdent := typeSpec.Name

    if interfaceTypeIdent == nil || interfaceTypeIdent.Name == "" {
      return true
    }
    if strings.HasPrefix(interfaceTypeIdent.Name, unsafePrefix) || !strings.HasSuffix(interfaceTypeIdent.Name, serverSuffix) {
      return true
    }
    grpcServerName := strings.TrimSuffix(interfaceTypeIdent.Name, serverSuffix)

    interfaceTypeMethods := interfaceType.Methods

    if interfaceTypeMethods == nil || len(interfaceTypeMethods.List) == 0 {
      return true
    }
    grpcServerCalls := make([]*grpcServerInterfaceCall, 0, len(interfaceTypeMethods.List))

    for _, interfaceMethod := range interfaceTypeMethods.List {
      if names := interfaceMethod.Names; len(names) == 0 || names[0].Name == "" {
        err = fmt.Errorf("encountered invalid grpc server interface method")
        return false
      }
      grpcCallName := interfaceMethod.Names[0].Name

      funcTyp, ok := interfaceMethod.Type.(*ast.FuncType)
      if !ok {
        continue
      }
      if funcTyp.Params == nil || len(funcTyp.Params.List) == 0 {
        continue
      }
      if funcTyp.Results == nil || len(funcTyp.Results.List) == 0 {
        continue
      }
      var grpcCallInputProto string

      for _, methodParam := range funcTyp.Params.List {
        paramStarExpr, ok := methodParam.Type.(*ast.StarExpr)
        if !ok {
          continue
        }
        paramIdent, ok := paramStarExpr.X.(*ast.Ident)
        if !ok {
          continue
        }
        if paramIdent.Name == "" {
          continue
        }
        grpcCallInputProto = paramIdent.Name
        break
      }

      if grpcCallInputProto == "" {
        err = fmt.Errorf("encountered invalid grpc server interface method")
        return false
      }
      var grpcCallOutputProto string

      for _, methodResult := range funcTyp.Results.List {
        resultStarExpr, ok := methodResult.Type.(*ast.StarExpr)
        if !ok {
          continue
        }
        resultIdent, ok := resultStarExpr.X.(*ast.Ident)
        if !ok {
          continue
        }
        if resultIdent.Name == "" {
          continue
        }
        grpcCallOutputProto = resultIdent.Name
        break
      }

      if grpcCallOutputProto == "" {
        err = fmt.Errorf("encountered invalid grpc server interface method")
        return false
      }
      grpcCall := &grpcServerInterfaceCall{
        CallName:    grpcCallName,
        InputProto:  grpcCallInputProto,
        OutputProto: grpcCallOutputProto,
      }
      grpcServerCalls = append(grpcServerCalls, grpcCall)
    }

    grpcServer = &grpcServerInterface{
      ServerName:  grpcServerName,
      ServerCalls: grpcServerCalls,
    }
    return false
  })

  if grpcServer == nil {
    err = fmt.Errorf("encountered invalid grpc server")
  }
  if err != nil {
    return nil, fmt.Errorf("ast.Inspect: %w", err)
  }
  return grpcServer, nil
}

func regenerateGrpcService(filePath string, serviceDesc *grpcServiceDesc) error {
  if err := validateGrpcFileName(filePath); err != nil {
    return fmt.Errorf("validateGrpcFileName: %w", err)
  }

  // TODO: complete this method

  // TODO: add (s *Implementation) Register(params *app.RegisterParams) regenerate

  return nil
}

func regenerateGrpcServiceStub(filePath string, serviceCallDesc *grpcServiceCallDesc) error {
  if err := validateGrpcFileName(filePath); err != nil {
    return fmt.Errorf("validateGrpcFileName: %w", err)
  }
  fileSet := token.NewFileSet()

  astFile, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
  if err != nil {
    return fmt.Errorf("parser.ParseFile: %w", err)
  }

  ast.Inspect(astFile, func(node ast.Node) bool {
    funcDecl, ok := node.(*ast.FuncDecl)
    if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
      return true
    }
    funcRecv := funcDecl.Recv.List[0]

    if funcRecv == nil || funcRecv.Type == nil {
      return true
    }
    funcRecvStarExpr, ok := funcRecv.Type.(*ast.StarExpr)
    if !ok || funcRecvStarExpr.X == nil {
      return true
    }
    const funcRecvName = "Implementation"

    funcRecvIdent, ok := funcRecvStarExpr.X.(*ast.Ident)
    if !ok || funcRecvIdent.Name != funcRecvName {
      return true
    }

    if funcDecl.Name == nil || !funcDecl.Name.IsExported() || funcDecl.Name.Name != serviceCallDesc.CallName {
      return true
    }

    if funcDecl.Type == nil {
      err = fmt.Errorf("encountered invalid grpc service method")
      return false
    }

    funcDecl.Type.Params = buildAstFuncParams(serviceCallDesc)
    funcDecl.Type.Results = buildAstFuncResults(serviceCallDesc)

    return false

  })

  if err != nil {
    return fmt.Errorf("ast.Inspect: %w", err)
  }

  osFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
  if err != nil {
    return fmt.Errorf("os.OpenFile: %w", err)
  }

  if err = printer.Fprint(osFile, fileSet, astFile); err != nil {
    return fmt.Errorf("printer.Fprint: %w", err)
  }

  return nil
}

func validateGrpcFileName(filePath string) error {
  if ext := filer.ExtractFileExtension(filePath); ext != "go" {
    return fmt.Errorf("not a .go file specified: extension: %s", ext)
  }
  return nil
}

func validateGrpcProtocFileName(filePath string) error {
  if ext := filer.ExtractFileExtension(filePath); ext != "go" {
    return fmt.Errorf("not a .go file specified: extension: %s", ext)
  }
  fileName := filer.ExtractFileName(filePath)

  if fileName = strings.TrimSuffix(fileName, ".pb.go"); !strings.HasSuffix(fileName, "_grpc") {
    return fmt.Errorf("not a grpc generated .go file: name: %s", fileName)
  }
  return nil
}
