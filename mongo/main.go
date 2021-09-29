package mongo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var ctx = context.TODO()
var files *mongo.Collection

var FileNotFoundError = errors.New("mongo: no documents in result")

type FileRecord struct {
	Id uuid.UUID `bson:"_id"`

	Type string   `bson:"type"`
	Payload     []byte   `bson:"payload"`

	UploadedAt time.Time `bson:"uploaded_at"`
	IP         string    `bson:"ip"`
}

func Init() {
	var client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("MongoDB init failed:", err)
	}

	err = client.Connect(ctx)
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("MongoDB connect-ping failed:", err)
	}

	var db *mongo.Database

	db = client.Database("pictrick")
	files = db.Collection("files")
}


func SaveFile(payload []byte, contentType, ip string) (uuid.UUID, error) {
	id := uuid.New()

	var fileRecord = FileRecord{
		Id: id,
		Type: contentType,
		Payload:    payload,
		UploadedAt: time.Now(),
		IP:         ip,
	}

	_, err := files.InsertOne(ctx, fileRecord)
	if err != nil {
		log.Println(err)
		return uuid.UUID{}, err
	}

	return id, nil
}

func GetFilePayload(id uuid.UUID) ([]byte, string, error) {
	var filter = struct {
		_id uuid.UUID `bson:"_id"`
	}{
		_id: id,
	}

	var data struct {
		Type string              `bson:"type"`
		Payload     []byte       `bson:"payload"`
	}

	err := files.FindOne(ctx,
		filter,
	).Decode(&data)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, "", FileNotFoundError
		}
	}

	return data.Payload, data.Type, err
}