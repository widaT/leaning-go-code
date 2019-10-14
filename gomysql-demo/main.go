package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

/*
CREATE TABLE `tb_test` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` int(20) NOT NULL DEFAULT '0' COMMENT '会员id',
  `img` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
 */
var db *sql.DB
func init()  {
	var err error
	db, err = sql.Open("mysql", "root:wida@tcp(localhost:3306)/test")
	if err != nil {
		panic(err.Error())
	}

	db.SetConnMaxLifetime(10*time.Minute)
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
}

func add()  {
	stmtIns, err := db.Prepare("INSERT INTO tb_test VALUES(?,?,?)") // ? = placeholder
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i+1, i * i,fmt.Sprintf("http://abc.com/%d.jpg",i))
		if err != nil {
			panic(err.Error())
		}
	}
}

func query()  {
	//findone
	stmtOut, err := db.Prepare("SELECT uid FROM tb_test WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()
	var squareNum int

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)  //查询超时时间
	defer cancel()
	err = stmtOut.QueryRowContext(ctx,13).Scan(&squareNum)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("The square number of 13 is: %d \n", squareNum)
	err = stmtOut.QueryRow(5).Scan(&squareNum)  //不带超时
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("The square number of 5 is: %d \n", squareNum)

	//findmany
	stmtOut, err = db.Prepare("SELECT * FROM tb_test")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	type Entry struct{
		Id int32
		Uid int
		Img string}
	var entrys []Entry

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)  //查询超时时间，对耗时的查询使用超时处理对程序的健壮性有很大帮助
	defer cancel()
	rows ,err := stmtOut.QueryContext(ctx)
	if err!=nil {
		println(err)
	}
	defer rows.Close()
	for rows.Next() {
		entry := Entry{}
		rows.Scan(&entry.Id,&entry.Uid,&entry.Img) //这边需要和数据库的字段顺序保持一致,另外一种方法是select中指定字段，这边scan的顺序和指定的字段顺序一致
		entrys = append(entrys,entry)
	}
	fmt.Println(entrys)
}

func update()  {
	stm,_ := db.Prepare("update tb_test set uid=? where id=? ")
	ret ,err := stm.Exec(999,1)
	if err !=nil {
		panic(err)
	}

	fmt.Println(ret.RowsAffected()) //影响条数
	stm.Close()
}

func delete()  {
	stm,_ := db.Prepare("DELETE from tb_test  where id=? ")
	ret ,err := stm.Exec(25)
	if err !=nil {
		panic(err)
	}

	fmt.Println(ret.RowsAffected()) //影响条数
	stm.Close()
}

func transaction()  {
	tx ,_:=db.Begin()
	stmt,err := tx.Prepare("SELECT uid FROM tb_test WHERE id = ?")
	if err !=nil {
		tx.Rollback()
		return
	}
	var uid  int32
	stmt.QueryRow(1).Scan(&uid)

	ret,err := tx.Exec("UPDATE tb_test set img='http://abc.com/3.jpg' where id=?",1)
	if err !=nil {
		tx.Rollback()
		return
	}
	fmt.Println(ret.RowsAffected())

	tx.Commit()
}


func main() {
	//query()
	//update()

	//delete()
	transaction()
}