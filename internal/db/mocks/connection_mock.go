package mocks

import (
	"github.com/derickit/go-rest-api/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	PingFunc func() error
)

type MockMongoMgr struct{}

func (m *MockMongoMgr) Ping() error {
	return PingFunc()
}

func (m *MockMongoMgr) Database() db.MongoDatabase {
	return &MockMongoDataBase{}
}

func (m *MockMongoMgr) Disconnect() error {
	return nil
}

type MockMongoDataBase struct{}

func (m *MockMongoDataBase) Collection(_ string, _ ...*options.CollectionOptions) *mongo.Collection {
	return nil
}
