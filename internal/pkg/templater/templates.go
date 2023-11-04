package templater

import (
  "bytes"
  "fmt"
  "go/format"
  "os"
  "text/template"
)

func CopyTemplate(templateCompiled string, filePath string) error {
  if err := os.WriteFile(filePath, []byte(templateCompiled), os.ModePerm); err != nil {
    return fmt.Errorf("os.WriteFile: %w", err)
  }
  return nil
}

func ExecTemplate(templateCompiled string, dataPtr any, funcMap template.FuncMap) ([]byte, error) {
  t := template.New("")

  if len(funcMap) > 0 {
    t = t.Funcs(funcMap)
  }
  t, err := t.Parse(templateCompiled)
  if err != nil {
    return nil, fmt.Errorf("template.New().Parse: %w", err)
  }
  var buffer bytes.Buffer

  if err = t.Execute(&buffer, dataPtr); err != nil {
    return nil, fmt.Errorf("t.Execute: %w", err)
  }
  return buffer.Bytes(), nil
}

func ExecTemplateCopy(templateCompiled, filePath string, dataPtr any, funcMap template.FuncMap) error {
  var (
    buf []byte
    err error
  )
  if buf, err = ExecTemplate(templateCompiled, dataPtr, funcMap); err != nil {
    return err
  }
  if err = os.WriteFile(filePath, buf, os.ModePerm); err != nil {
    return fmt.Errorf("os.WriteFile: %w", err)
  }
  return nil
}

func ExecTemplateCopyWithGoFmt(templateCompiled, filePath string, dataPtr any, funcMap template.FuncMap) error {
  var (
    buf []byte
    err error
  )
  if buf, err = ExecTemplate(templateCompiled, dataPtr, funcMap); err != nil {
    return err
  }
  if buf, err = format.Source(buf); err != nil {
    return fmt.Errorf("format.Source: %w", err)
  }
  if err = os.WriteFile(filePath, buf, os.ModePerm); err != nil {
    return fmt.Errorf("os.WriteFile: %w", err)
  }
  return nil
}
