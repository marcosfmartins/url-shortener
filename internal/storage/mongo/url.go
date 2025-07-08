package storage

import (
	"context"
	"github.com/marcosfmartins/url-shortener/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const URLStorageCollection = "urls"

type URLStorage struct {
	coll *mongo.Collection
}

func NewURLStorage(coll *mongo.Database) *URLStorage {
	return &URLStorage{
		coll: coll.Collection(URLStorageCollection),
	}
}

func (s *URLStorage) Insert(ctx context.Context, url *entity.URL) error {
	_, err := s.coll.InsertOne(ctx, url)
	if err != nil {
		return parserError(err)
	}
	return nil
}

func (s *URLStorage) FindByID(ctx context.Context, id string) (*entity.URL, error) {
	var url entity.URL
	err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&url)
	if err != nil {
		return nil, parserError(err)
	}
	return &url, nil
}

func (s *URLStorage) DeleteByID(ctx context.Context, id string) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return parserError(err)
	}
	return nil
}

func (s *URLStorage) Increment(ctx context.Context, objs []entity.URL) error {
	operations := make([]mongo.WriteModel, 0, len(objs))
	for _, value := range objs {
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.D{{Key: "_id", Value: value.ID}}).
			SetUpdate(bson.D{
				{Key: "$inc", Value: bson.D{{Key: "hits", Value: value.Hits}}},
				{Key: "$set", Value: bson.D{{Key: "lastAccess", Value: value.LastAccess}}},
			}).
			SetUpsert(true)

		operations = append(operations, model)
	}

	_, err := s.coll.BulkWrite(ctx, operations)
	if err != nil {
		return parserError(err)
	}

	return nil
}
