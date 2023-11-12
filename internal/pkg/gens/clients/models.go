package clients

import (
  "fmt"
  "path/filepath"
  "strings"
)

type protoClients struct {
  ProtoClients []*protoClient `json:"proto_clients" yaml:"proto_clients"`
}

type protoClient struct {
  // Proto file import path
  ProtoPath string `json:"import" yaml:"import"`
  // Is local proto file import
  Local bool `json:"local" yaml:"local"`
}

func validateProtoClient(client *protoClient) error {
  partsImport := strings.Split(client.ProtoPath, "@")

  if isLocal := client.Local; isLocal {
    if len(partsImport) != 1 {
      return fmt.Errorf("invalid local proto client import: %s", client.ProtoPath)
    }
    protoPath := partsImport[0]
    filePrefix := filepath.Join(".boiler", "vendor")

    if !strings.HasPrefix(protoPath, filePrefix) {
      return fmt.Errorf("local proto client must be placed in: %s", filePrefix)
    }
  } else {
    if len(partsImport) != 2 {
      return fmt.Errorf("invalid github proto client import: %s", client.ProtoPath)
    }
    protoPath := partsImport[0]
    filePrefix := "github.com"

    if !strings.HasPrefix(protoPath, filePrefix) {
      return fmt.Errorf("github proto client must be placed in: %s", filePrefix)
    }
  }
  return nil
}
