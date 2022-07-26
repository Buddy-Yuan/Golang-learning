package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
)

type mysql_db struct {
	db *sql.DB
}

type user struct {
	id    int
	name  string
	email string
	age   int
}

func (f *mysql_db) mysqlOpen() {
	db, err := sql.Open("mysql", "tidb:123456@tcp(xx.xx.xx.xx:4000)/test")
	if err != nil {
		fmt.Println("connection database failed")
	} else {
		fmt.Println("connection database successd")
	}
	f.db = db
}

func (f *mysql_db) mysqlClose() {
	defer f.db.Close()
}

//插入数据
func (f *mysql_db) mysqlInsert(ch1 chan struct{}) {
	stmt, err := f.db.Prepare("INSERT INTO user(id,name,email,age) VALUES(?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("insert failed")
		return
	}
	stmt.Exec("1", "name1", "name1@pingcap.com", 10)
	stmt.Exec("2", "name2", "name2@pingcap.com", 20)
	stmt.Exec("3", "name3", "name3@pingcap.com", 30)
	stmt.Exec("4", "name4", "name4@pingcap.com", 40)
	stmt.Exec("5", "name5", "name5@pingcap.com", 50)
	fmt.Println("insert success")
	//数据插入成功，通知下一个goroutine
	//有时候使用 channel 不需要发送任何的数据，只用来通知子协程(goroutine)执行任务，或只用来控制协程并发度。
	//这种情况下，使用空结构体作为占位符就非常合适了。
	fmt.Println("notification goroutine mysqlSelect")
	ch1 <- struct{}{}
}

//查询数据
func (f *mysql_db) mysqlSelect(ch1 chan struct{}, ch2 chan string) {
	<-ch1
	fmt.Println("now,start goroutine mysqlSelect")
	//现在可以关闭关闭通道1
	close(ch1)

	sqlStr := "select id, name,email, age from user"
	rows, err2 := f.db.Query(sqlStr)
	if err2 != nil {
		fmt.Printf("query failed, err:%v\n", err2)
		return
	}
	defer rows.Close()
	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.email, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		s1 := strconv.Itoa(u.id) + u.name + u.email + strconv.Itoa(u.age)
		fmt.Println("notification goroutine WriteFile")
		ch2 <- s1
	}
	close(ch2)
}

func WriteFile(ch2 chan string, ch3 chan struct{}) {
	filePath := "/tmp/WriteFile.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file failed, err:%v\n", err)
		return
	}
	//及时关闭file句柄
	defer file.Close()

	//循环接收通道的数据
	for {
		ch2strings, ok := <-ch2
		if !ok {
			//全部写入完成，通知结束
			fmt.Println("notification goroutine done")
			break
		} else {
			write := bufio.NewWriter(file)
			write.WriteString(ch2strings + "\n")
			write.Flush()
		}
	}
	ch3 <- struct{}{}
}

func main() {
	db := &mysql_db{}
	ch1 := make(chan struct{})
	ch2 := make(chan string)
	ch3 := make(chan struct{})
	db.mysqlOpen()
	//
	go db.mysqlInsert(ch1)
	go db.mysqlSelect(ch1, ch2)
	go WriteFile(ch2, ch3)

	<-ch3
	close(ch3)
	db.mysqlClose()

}
