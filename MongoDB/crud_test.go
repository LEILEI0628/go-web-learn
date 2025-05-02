package MongoDB

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
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
	// bson
	filter := bson.D{bson.E{Key: "id", Value: 1}} // id=1
	var table Table
	err = col.FindOne(ctx, filter).Decode(&table)
	assert.NoError(t, err)
	fmt.Printf("%#v \n", table)
	err = col.FindOne(ctx, Table{Id: 1}).Decode(&table)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// 没有数据
		t.Log("没有数据")
	} else {
		assert.NoError(t, err)
		fmt.Printf("%#v \n", table)
	}
	sets := bson.D{bson.E{Key: "$set", Value: bson.E{Key: "title", Value: "New Title 新标题"}}}
	updateRes, err := col.UpdateMany(ctx, filter, sets)
	assert.NoError(t, err)
	fmt.Println(updateRes.MatchedCount, updateRes.ModifiedCount)

	updateRes, err = col.UpdateMany(ctx, filter, bson.D{
		bson.E{Key: "$set", Value: Table{Title: "新标题2", Content: "新内容"}}})
	assert.NoError(t, err)
	fmt.Println("affected", updateRes.MatchedCount, updateRes.ModifiedCount)

	////or := bson.A{bson.D{bson.E{"id", 123}},
	////	bson.D{bson.E{"id", 456}}}
	//or := bson.A{bson.M{"id": 123}, bson.M{"id": 456}}
	//orRes, err := col.Find(ctx, bson.D{bson.E{"$or", or}})
	//assert.NoError(t, err)
	//var ars []Table
	//err = orRes.All(ctx, &ars)
	//assert.NoError(t, err)
	//
	//and := bson.A{bson.D{bson.E{"id", 123}},
	//	bson.D{bson.E{"title", "我的标题2"}}}
	//andRes, err := col.Find(ctx, bson.D{bson.E{"$and", and}})
	//assert.NoError(t, err)
	//ars = []Table{}
	//err = andRes.All(ctx, &ars)
	//assert.NoError(t, err)
	//
	////in := bson.D{bson.E{"id", bson.D{bson.E{"$in", []any{123, 456}}}}}
	//in := bson.D{bson.E{"id", bson.M{"$in": []any{123, 456}}}}
	//inRes, err := col.Find(ctx, in)
	//ars = []Table{}
	//err = inRes.All(ctx, &ars)
	//assert.NoError(t, err)
	//
	//inRes, err = col.Find(ctx, in, options.Find().SetProjection(bson.M{
	//	"id":    1,
	//	"title": 1,
	//}))
	//ars = []Table{}
	//err = inRes.All(ctx, &ars)
	//assert.NoError(t, err)
	//
	//idxRes, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
	//	{
	//		Keys:    bson.M{"id": 1},
	//		Options: options.Index().SetUnique(true),
	//	},
	//	{
	//		Keys: bson.M{"author_id": 1},
	//	},
	//})
	//assert.NoError(t, err)
	//fmt.Println(idxRes)

	delRes, err := col.DeleteMany(ctx, filter)
	assert.NoError(t, err)
	fmt.Println("deleted", delRes.DeletedCount)
}

type Table struct {
	Id      int    `bson:"id,omitempty"`    // 如果不加omitempty使用FindOne(ctx, Table{Id: 1})查询时会报错
	Title   string `bson:"title,omitempty"` // MongoDB不会自动忽略零值
	Content string `bson:"content,omitempty"`
}
