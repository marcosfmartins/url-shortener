package storage

import (
	"context"
	"errors"
	"github.com/marcosfmartins/url_shortener/internal/entity"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

func NewClient(URL string) (*mongo.Client, *mongo.Database, error) {
	dbName, err := getDBName(URL)
	if err != nil {
		return nil, nil, err
	}

	opt := options.Client().ApplyURI(URL)

	conn, err := mongo.Connect(opt)
	if err != nil {
		return nil, nil, err
	}

	err = conn.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, nil, err
	}

	return conn, conn.Database(dbName), nil
}

func getDBName(URL string) (string, error) {
	obj, err := connstring.Parse(URL)
	if err != nil {
		return "", err
	}
	return obj.Database, nil
}

func parserError(err error) error {
	if err == nil {
		return nil
	}

	if mongo.IsDuplicateKeyError(err) {
		return entity.DuplicateKeyError.WithError(err)
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return entity.NotFoundError.WithError(err)
	}

	return entity.InternalServerError.WithError(err)
}
