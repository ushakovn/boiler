package kafkaoutbox

import (
  "fmt"
  "os"

  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/yoheimuta/go-protoparser"
  "github.com/yoheimuta/go-protoparser/parser"
)

type parsedProtoMessage struct {
  messageName string
  tableName   string
}

type parsedProto struct {
  messages []*parsedProtoMessage
}

func parseProto(filePath, optionName string) (*parsedProto, error) {
  file, err := os.Open(filePath)
  if err != nil {
    return nil, fmt.Errorf("file open failed: %w", err)
  }
  parsed, err := protoparser.Parse(file)
  if err != nil {
    return nil, fmt.Errorf("protoparser.Parse: %w", err)
  }
  var (
    messages []*parsedProtoMessage
    token    parser.Visitee
  )
  for _, token = range parsed.ProtoBody {
    message, ok := token.(*parser.Message)
    if !ok {
      continue
    }
    for _, token = range message.MessageBody {
      option, ok2 := token.(*parser.Option)
      if !ok2 {
        continue
      }
      if option.OptionName != optionName {
        continue
      }
      tableName := stringer.UnquoteString(option.Constant)

      messages = append(messages, &parsedProtoMessage{
        messageName: message.MessageName,
        tableName:   tableName,
      })
    }
  }
  return &parsedProto{
    messages: messages,
  }, nil
}
