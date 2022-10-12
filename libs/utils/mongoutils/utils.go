package mongoutils

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	tUUID       = reflect.TypeOf(uuid.UUID{})
	uuidSubtype = byte(0x04)

	mongoRegistry = bson.NewRegistryBuilder().
			RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(uuidEncodeValue)).
			RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(uuidDecodeValue)).
			Build()
)

func uuidEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tUUID {
		return bsoncodec.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}
	b := val.Interface().(uuid.UUID)
	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}

func uuidDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tUUID {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}

	var data []byte
	var subtype byte
	var err error
	switch vrType := vr.Type(); vrType {
	case bsontype.Binary:
		data, subtype, err = vr.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if err != nil {
		return err
	}
	uuid2, err := uuid.FromBytes(data)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(uuid2))
	return nil
}

// SetRegistryForUuids sets the registry for mongodb such that
// it can encode and decode primitive.Binary UUIDs into google/uuid.UUID type
func SetRegistryForUuids(
	opts *options.ClientOptions,
) *options.ClientOptions {
	return opts.SetRegistry(mongoRegistry)
}

func BinaryToUuid(id interface{}) uuid.UUID {
	bid := id.(primitive.Binary).Data
	return *(*uuid.UUID)(bid)
}

type Config struct {
	Uri string
	DB  string
}

func ProvideMongoClient(
	lc fx.Lifecycle,
	logger *zap.Logger,
	cfg *Config,
) (*mongo.Client, error) {
	client, err := mongo.NewClient(
		SetRegistryForUuids(
			options.Client().ApplyURI(cfg.Uri),
		),
	)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Sugar().Info("Connecting to mongo")
			err := client.Connect(ctx)
			if err != nil {
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Sugar().Info("Disconnecting from database")

			return client.Disconnect(ctx)
		},
	})

	return client, nil
}

func ProvideMongoDatabase(
	client *mongo.Client,
	cfg *Config,
) *mongo.Database {
	return client.Database(cfg.DB)
}
