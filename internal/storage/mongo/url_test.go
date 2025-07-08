package storage

import (
	"context"
	"github.com/marcosfmartins/url-shortener/internal/entity"
	"github.com/marcosfmartins/url-shortener/pkg/id"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"os"
	"testing"
)

var (
	conn *mongo.Client
	db   *mongo.Database
	ctx  context.Context = context.Background()
)

func init() {
	URL := os.Getenv("MONGO_URL")

	var err error
	conn, db, err = NewClient(URL)
	if err != nil {
		panic(err)
	}
}

func TestNewURLStorage(t *testing.T) {
	obj := NewURLStorage(db)
	assert.NotNil(t, obj)
}

func TestURLStorage_Insert(t *testing.T) {
	cli := NewURLStorage(db)
	assert.NotNil(t, cli)

	id, _ := id.GenerateID()

	obj := &entity.URL{
		ID: id,
	}
	err := cli.Insert(ctx, obj)
	assert.NoError(t, err)

	result := &entity.URL{}
	err = cli.coll.FindOne(ctx, bson.M{"_id": id}).Decode(result)
	assert.NoError(t, err)
	assert.Equal(t, obj, result)
}

func TestURLStorage_FindByID(t *testing.T) {
	cli := NewURLStorage(db)
	assert.NotNil(t, cli)

	id, _ := id.GenerateID()
	obj := &entity.URL{
		ID: id,
	}
	err := cli.Insert(ctx, obj)
	assert.NoError(t, err)

	result, err := cli.FindByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, obj, result)
}

func TestURLStorage_DeleteByID(t *testing.T) {
	cli := NewURLStorage(db)
	assert.NotNil(t, cli)

	id, _ := id.GenerateID()
	obj := &entity.URL{
		ID: id,
	}
	err := cli.Insert(ctx, obj)
	assert.NoError(t, err)

	err = cli.DeleteByID(ctx, id)
	assert.NoError(t, err)

	result, err := cli.FindByID(ctx, id)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestURLStorage_Increment(t *testing.T) {
	cli := NewURLStorage(db)
	assert.NotNil(t, cli)

	id, _ := id.GenerateID()
	obj := entity.URL{ID: id}
	err := cli.Insert(ctx, &obj)
	assert.NoError(t, err)

	obj.Hits = 10

	err = cli.Increment(ctx, []entity.URL{obj})
	assert.NoError(t, err)

	result, err := cli.FindByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, obj, *result)
}
