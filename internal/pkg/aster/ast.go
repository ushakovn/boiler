package aster

import (
  "fmt"
  "go/ast"
  "go/parser"
  "go/token"

  "github.com/ushakovn/boiler/internal/pkg/filer"
)

func FindMethodDeclaration(filePath string, methodName string) (bool, error) {
  if fileExtension := filer.ExtractFileExtension(filePath); fileExtension != "go" {
    return false, fmt.Errorf("not a .go file specified: extension: %s", fileExtension)
  }
  if methodName == "" {
    return false, fmt.Errorf("method name not specified")
  }
  fileSet := token.NewFileSet()

  astFile, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
  if err != nil {
    return false, fmt.Errorf("parser.ParseFile: %w", err)
  }
  var methodFound bool

  ast.Inspect(astFile, func(node ast.Node) bool {
    funcDecl, ok := node.(*ast.FuncDecl)
    if !ok || funcDecl.Name == nil || funcDecl.Name.Name == "" {
      return true
    }
    if methodFound = funcDecl.Name.Name == methodName; methodFound {
      return false
    }
    return true
  })

  return methodFound, nil
}
