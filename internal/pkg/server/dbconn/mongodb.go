package dbconn

import (
	"context"
	"snsDownloader/internal/pkg/colorLog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var client *mongo.Client

var parentContext = context.Background()
var parentContextCancelFunc context.CancelFunc

var collection *mongo.Collection

const DATABASE = "TEST_DATA_BASE"

func Create(dbUri string) {
	colorLog.Info("init mongo")

	var cancel context.CancelFunc
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// server context
	parentContext, parentContextCancelFunc = context.WithCancel(context.Background())

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if err != nil {
		colorLog.Error("connection fail. error=%v", err.Error())
		panic(err)
	}

	// cluster일 경우는 모든 곳에 다 연결되므로 localhost시에는 문제가 된다.
	if 1 != 1 {
		if err = client.Ping(ctx, nil); err != nil {
			colorLog.Error("connection ping fail. error=%v", err.Error())
			panic(err)
		}
	}

	// cluster 연결이라 로컬에는 안됨
	if 1 != 1 {
		if serverStatus, err := client.Database("admin").RunCommand(ctx, bsonx.Doc{{"serverStatus", bsonx.Int32(1)}}).DecodeBytes(); err != nil {
			colorLog.Error("fail server status. error=%v", err.Error())
		} else if version, err := serverStatus.LookupErr("version"); err != nil {
			colorLog.Error("fail server version. error=%v", err.Error())
		} else {
			colorLog.Info("server version : %v", version.String())
		}
	}
	initSchema()
}

func initSchema() {
	collection := GetCollection(DATABASE, "tb_progress")
	mongod := mongo.IndexModel{
		Keys: bson.M{
			"progressTm": 1,
		}, Options: options.Index().SetUnique(true).SetName("idx_tb_progress_tm_01"),
	}

	if indexName, err := collection.Indexes().CreateOne(context.Background(), mongod); err != nil {
		colorLog.Error("	failed create index. index=%+v,err=%+v", indexName, err)
	} else {
		colorLog.Debug("	create unique index. index=%+v", indexName)
	}

	collection = GetCollection(DATABASE, "tb_content")
	mongod = mongo.IndexModel{
		Keys: bson.M{
			"uuid": 1,
		}, Options: options.Index().SetUnique(true).SetName("idx_tb_content_uuid_01"),
	}
	if indexName, err := collection.Indexes().CreateOne(context.Background(), mongod); err != nil {
		colorLog.Error("	failed create index. index=%+v,err=%+v", indexName, err)
	} else {
		colorLog.Debug("	create unique index. index=%+v", indexName)
	}

}

// GetCollection. 특정 Collection을 반환
func GetCollection(database string, collection string) *mongo.Collection {
	db := client.Database(database)
	return db.Collection(collection)
}

func GetContext() context.Context {
	return parentContext
}

// Close. 연결된 Client를 종료한다.
func Close() {
	colorLog.Info("close mongodb")
	if client == nil {
		return
	}

	// context 관련 작업 종료 (context canceled 처리. <-Done()
	parentContextCancelFunc()

	if err := client.Disconnect(parentContext); err != nil {
		colorLog.Error("disconnect fail. error=%v", err.Error())
		//panic(err)
	}
}
