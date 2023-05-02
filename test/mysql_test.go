package test

import (
	"context"
	"fmt"
	"github.com/weiyouwozuiku/Gateway/handler"
	"github.com/weiyouwozuiku/Gateway/public"
	"testing"
	"time"
)

const (
	beginSQL       = "start transaction;"
	commitSQL      = "commit;"
	rollbackSQL    = "rollback;"
	createTableSQL = "CREATE TABLE `test1` (`id` int(12) unsigned NOT NULL AUTO_INCREMENT" +
		" COMMENT '自增id',`name` varchar(255) NOT NULL DEFAULT '' COMMENT '姓名'," +
		"`created_at` datetime NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB " +
		"DEFAULT CHARSET=utf8"
	insertSQL    = "INSERT INTO `test1` (`id`, `name`, `created_at`) VALUES (NULL, '111', '2022-09-22 00:18:08');"
	dropTableSQL = "DROP TABLE `test1`"
	selectSQL    = "SELECT id,created_at FROM test1 WHERE id>? order by id asc"
)

type Test1 struct {
	Id        int64     `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func Test_DBPool(t *testing.T) {
	SetUp()
	// 获取连接池
	dbpool, err := handler.GetDBPool("default")
	if err != nil {
		t.Fatal(err)
	}
	// 执行事务
	trace := public.NewTrace()
	if _, err := handler.DBPoolQuery(trace, dbpool, beginSQL); err != nil {
		t.Fatal(err)
	}
	// 创建表
	if _, err := handler.DBPoolQuery(trace, dbpool, createTableSQL); err != nil {
		handler.DBPoolQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}
	// 插入数据
	if _, err := handler.DBPoolQuery(trace, dbpool, insertSQL); err != nil {
		handler.DBPoolQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}
	// // 插入数据
	// if _, err := middleware.DBPoolQuery(trace, dbpool, insertSQL); err != nil {
	// 	middleware.DBPoolQuery(trace, dbpool, rollbackSQL)
	// 	t.Fatal(err)
	// }
	current_id := 0
	table_name := "test1"
	fmt.Println("begin read table ", table_name, "")
	fmt.Println("------------------------------------------------------------------------")
	fmt.Printf("%6s | %6s\n", "id", "created_at")
	for {
		rows, err := handler.DBPoolQuery(trace, dbpool, selectSQL, current_id)
		if err != nil {
			handler.DBPoolQuery(trace, dbpool, rollbackSQL)
			t.Fatal(err)
		}
		defer rows.Close()
		row_len := 0
		for rows.Next() {
			create_time := ""
			if err := rows.Scan(&current_id, &create_time); err != nil {
				handler.DBPoolQuery(trace, dbpool, rollbackSQL)
				t.Fatal(err)
			}
			fmt.Printf("%6d | %6s\n", current_id, create_time)
			row_len++
		}
		if row_len == 0 {
			break
		}
	}
	fmt.Println("------------------------------------------------------------------------")
	fmt.Println("finish read table ", table_name, "")
	// 删除表
	if _, err := handler.DBPoolQuery(trace, dbpool, dropTableSQL); err != nil {
		handler.DBPoolQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}
	// 提交事务
	handler.DBPoolQuery(trace, dbpool, commitSQL)
	TearDown()
}

func Test_GormPool(t *testing.T) {
	SetUp()
	dbpool, err := handler.GetGORMPool(handler.DBDefault)
	if err != nil {
		t.Fatal(err)
	}
	db := dbpool.Begin()
	trace := public.NewTrace()
	ctx := context.Background()
	ctx = public.SetTraceContext(ctx, trace)
	db = db.WithContext(ctx)
	if err := db.Exec(createTableSQL).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	t1 := &Test1{Name: "test_name", CreatedAt: time.Now()}
	if err := db.Save(t1).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	list := []Test1{}
	if err := db.Where("name=?", "test_name").Find(&list).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	fmt.Println(list)
	if err := db.Exec(dropTableSQL).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	db.Commit()
	TearDown()
}
