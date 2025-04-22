package MongoDB

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/event"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongo(t *testing.T) {
	// 控制初始化超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			// 每个命令查询之前
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, startedEvent *event.CommandSucceededEvent) {
			// 执行成功
		},
		Failed: func(ctx context.Context, startedEvent *event.CommandFailedEvent) {
			// 执行失败
		},
	}
	opts := options.Client().ApplyURI("mongodb://root:example@localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)
	mdb := client.Database("test")
	col := mdb.Collection("table")
	res, err := col.InsertOne(ctx, &Table{Id: 1, Title: "test", Content: "test"})
	assert.NoError(t, err)
	// 这是MongoDB的ID，即_id字段
	fmt.Printf("ID:%s", res.InsertedID)

}

type Table struct {
	Id      int
	Title   string
	Content string
}
