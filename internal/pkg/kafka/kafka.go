package kafka

// import (
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// )

// func init() {

// }

// /*
// 	여기 지금있는건 무한루프로 돌아가는데
// 	원래 설계는 크론으로 돌고 있을때 몇개 정도만 체킹해서 작업하는거라
// 	요기서 추가해야할 사항으론 머 무한루프 돌고 있으니.. 플래그 처리를 하던
// 	메세지 오프셋 체킹을 하던 해야할듯?
// */
// func ConsumeMessage() {
// 	c, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers": "localhost:9092",
// 		"group.id":          "myGroup",
// 		// "auto.offset.reset": "207",
// 		"auto.offset.reset": "earliest",
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	c.SubscribeTopics([]string{"exam-topic"}, nil)

// 	for {
// 		msg, err := c.ReadMessage(time.Second)
// 		if err == nil {
// 			dataMap := make(map[string]interface{})
// 			// fmt.Printf("[offset:%v]: %s\n", msg.TopicPartition.Offset, string(msg.Value))
// 			json.Unmarshal(msg.Value, dataMap)
// 			fmt.Printf("msg %#v", dataMap)
// 		} else if err.(kafka.Error).Code() != kafka.ErrTimedOut {
// 			fmt.Printf("Consumer error: %v, msg:%v\n", err, msg)
// 		}
// 		time.Sleep(time.Second * 1)
// 	}

// 	c.Close()
// }
