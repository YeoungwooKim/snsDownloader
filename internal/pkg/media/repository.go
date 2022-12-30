package media

import (
	"context"
	"fmt"
	"headless/internal/pkg/server/dbconn"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func saveHistory(uuid string, dataMap map[string]interface{}) {
	collection := dbconn.GetCollection(dbconn.DATABASE, "tb_progress")
	fmt.Printf("[mongo]\t%+v\n", dataMap)
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
