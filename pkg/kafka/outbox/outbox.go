package outbox

import (
  "context"
  "fmt"
  "os"
  "reflect"
  "time"

  "github.com/IBM/sarama"
  validation "github.com/go-ozzo/ozzo-validation"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/kafka/models"
  "google.golang.org/protobuf/encoding/protojson"
  "google.golang.org/protobuf/reflect/protoreflect"
  "gopkg.in/yaml.v3"
)

type Outbox struct {
  config   Config
  storage  Storage
  producer Producer
}

type Config struct {
  LockTime     time.Duration `yaml:"lock_time"`
  WorkerIdle   time.Duration `yaml:"worker_idle"`
  WorkerRetry  int           `yaml:"worker_retry"`
  WorkersCount int           `yaml:"workers_count"`
}

type Producer interface {
  SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type Storage interface {
  BatchRecords(ctx context.Context, tableName string) ([]*models.Record, error)
  DeleteRecord(ctx context.Context, tableName, recordID string) error
}

func NewConfig(path string) (*Config, error) {
  buf, err := os.ReadFile(path)
  if err != nil {
    return nil, fmt.Errorf("file reading failed: %w", err)
  }
  config := &Config{}

  if err = yaml.Unmarshal(buf, config); err != nil {
    return nil, fmt.Errorf("yaml unmarshalling failed: %w", err)
  }
  return config, nil
}

func New(config Config, storage Storage, producer Producer) *Outbox {
  return &Outbox{
    config:   config,
    storage:  storage,
    producer: producer,
  }
}

func (o *Outbox) Run(ctx context.Context) error {
  if err := o.Validate(); err != nil {
    return fmt.Errorf("outbox validation failed: %w", err)
  }
  for _, tableName := range tableNames {
    o.sendWithWorkers(ctx, tableName)

    log.Infof("kafka-outbox: sending started for: %s table", tableName)
  }
  log.Infof("kafka-outbox: sending in progress")

  return nil
}

func (o *Outbox) sendWithWorkers(ctx context.Context, tableName string) {
  for worker := 0; worker < o.config.WorkersCount; worker++ {
    go func() { o.send(ctx, tableName) }()
  }
}

func (o *Outbox) send(ctx context.Context, tableName string) {
  ticker := time.NewTicker(o.config.WorkerIdle)
  for {
    select {
    case <-ticker.C:
      if err := o.sendBatch(ctx, tableName); err != nil {
        log.Errorf("outbox.send: table: %s: error: %v", tableName, err)
      }
    case <-ctx.Done():
      log.Errorf("outbox.send: table: %s: context cancelled", tableName)
      return
    }
  }
}

func (o *Outbox) sendBatch(ctx context.Context, tableName string) error {
  records, err := o.storage.BatchRecords(ctx, tableName)
  if err != nil {
    return fmt.Errorf("storage.BatchRecords: %w", err)
  }
  var msgBuf []byte

  for _, record := range records {
    msgBuf, err = marshalRecord(tableName, record)
    if err != nil {
      return fmt.Errorf("marshalRecord: %w", err)
    }
    topicName, ok := tableTopics[tableName]
    if !ok {
      return fmt.Errorf("topic not found: %s table", tableName)
    }
    msgKey := record.ID

    msgHeaders := []sarama.RecordHeader{
      {
        Key:   []byte("action_type"),
        Value: []byte(record.ActionType.String()),
      },
    }
    _, _, err = o.producer.SendMessage(&sarama.ProducerMessage{
      Topic:     topicName,
      Key:       sarama.StringEncoder(msgKey),
      Value:     sarama.StringEncoder(msgBuf),
      Headers:   msgHeaders,
      Timestamp: time.Now().UTC(),
    })
    if err != nil {
      return fmt.Errorf("producer.SendMessage: %w", err)
    }
    if err = o.storage.DeleteRecord(ctx, tableName, record.ID); err != nil {
      return fmt.Errorf("storage.DeleteRecord: %w", err)
    }
  }
  return nil
}

func (o *Outbox) Validate() error {
  return validation.ValidateStruct(o,
    validation.Field(&o.config),
    validation.Field(&o.storage, validation.Required),
    validation.Field(&o.producer, validation.Required),
  )
}

func (c *Config) Validate() error {
  return validation.ValidateStruct(c,
    validation.Field(&c.WorkersCount,
      validation.Required,
      validation.Min(1),
      validation.Max(5),
    ),
    validation.Field(&c.WorkerIdle,
      validation.Required,
      validation.Min(100*time.Millisecond),
      validation.Max(1*time.Second),
    ),
    validation.Field(&c.WorkerRetry,
      validation.Min(0*time.Millisecond),
      validation.Max(5*time.Second),
    ),
    validation.Field(&c.LockTime,
      validation.Min(100*time.Millisecond),
      validation.Max(5*time.Second),
    ),
  )
}

var tableTopics = map[string]string{
  clientsTableName: "kafka-outbox-clients-topic-name",
}

var tableTypes = map[string]any{
  clientsTableName: struct{}{}, // TODO: protobuf struct
}

const (
  clientsTableName = "clients"
)

var tableNames = []string{
  clientsTableName,
}

func marshalRecord(tableName string, record *models.Record) ([]byte, error) {
  typ, ok := tableTypes[tableName]
  if !ok {
    return nil, fmt.Errorf("type not found: %s table", tableName)
  }
  refTyp := reflect.TypeOf(typ)
  pb := reflect.New(refTyp).Interface().(protoreflect.ProtoMessage)

  if err := protojson.Unmarshal(record.JSONOut, pb); err != nil {
    return nil, fmt.Errorf("protojson.Unmarshal: %s table: %w", tableName, err)
  }
  buf, err := protojson.Marshal(pb)
  if err != nil {
    return nil, fmt.Errorf("protojson.Marshal: %s table: %w", tableName, err)
  }
  return buf, nil
}
