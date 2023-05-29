package media

import (
	"context"
	"fmt"
	"snsDownload/internal/pkg/server/dbconn"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getLatestHistory(uuid string) (map[string]interface{}, error) {
	var dataMap map[string]interface{}
	collection := dbconn.GetCollection(dbconn.DATABASE, "tb_content")
	filter := bson.M{
		"uuid": uuid,
	}

	if err := collection.FindOne(context.Background(), filter).Decode(&dataMap); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("not found - document")
		}
		return nil, err
	}
	return dataMap, nil
}

func GetNotCompleteHistory() ([]map[string]interface{}, error) {
	var dataMapList []map[string]interface{}
	collection := dbconn.GetCollection(dbconn.DATABASE, "tb_content")
	filter := bson.M{
		"completeFlag": 0,
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var foundData bson.M
		if err := cursor.Decode(&foundData); err != nil {
			return nil, err
		}
		dataMapList = append(dataMapList, foundData)
	}

	return dataMapList, nil
}

func saveHistory(uuid string, dataMap map[string]interface{}) {
	collection := dbconn.GetCollection(dbconn.DATABASE, "tb_progress")
	// fmt.Printf("[mongo]\t%+v\n", dataMap)
	for {
		now := time.Now()
		dataMap["uuid"] = uuid
		dataMap["progressTm"] = time.Now().UnixMilli()
		dataMap["writeDate"] = now.Format(time.RFC3339)
		if _, err := collection.InsertOne(context.Background(), dataMap); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				time.Sleep(250 * time.Millisecond)
				continue
			}
		} else {
			break
		}
	}
}

func commitProgress(uuid string) error {
	collection := dbconn.GetCollection(dbconn.DATABASE, "tb_content")
	filter := bson.M{
		"uuid": uuid,
	}

	update := bson.M{
		"$set": map[string]interface{}{
			"completeFlag": 1,
			"modifyDate":   time.Now().Format(time.RFC3339),
		},
	}

	if _, err := collection.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}

	return nil
}

// 해당 uuid의 튜플이 있으면 업데이트 없으면 생성.
func saveContent(uuid string, dataMap map[string]interface{}) error {
	collection := dbconn.GetCollection(dbconn.DATABASE, "tb_content")
	var foundDocument map[string]interface{}
	filter := bson.M{
		"uuid": uuid,
	}

	now := time.Now()
	if err := collection.FindOne(context.Background(), filter).Decode(&foundDocument); err != nil {
		if err == mongo.ErrNoDocuments {
			dataMap["writeDate"] = now.Format(time.RFC3339)
			dataMap["completeFlag"] = 0
			dataMap["uuid"] = uuid
			if _, err := collection.InsertOne(context.Background(), dataMap); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	dataMap["modifyDate"] = now.Format(time.RFC3339)
	for k, v := range dataMap {
		foundDocument[k] = v
	}
	update := bson.M{
		"$set": foundDocument,
	}

	if _, err := collection.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}

	return nil
}
