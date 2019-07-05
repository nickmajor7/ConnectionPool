package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pool "github.com/AZsoftAlanZheng/ConnectionPool"
)

const addr string = "127.0.0.1:8910"

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	go server()
	//等待tcp server启动
	time.Sleep(2 * time.Second)
	fmt.Println("使用: ctrl+c 退出服务")
	p := initPool()
	for i := 0; i < 10; i++ {
		go client(i, p)
	}
	<-c
	fmt.Println("關閉連接池…")
	p.Release()
	fmt.Println("服务退出")
}

func initPool() pool.Pool {
	//factory 创建连接的方法
	factory := func() (interface{}, error) { return net.Dial("tcp", addr) }

	//close 关闭连接的方法
	close := func(v interface{}) error { return v.(net.Conn).Close() }

	//创建一个连接池
	poolConfig := &pool.Config{
		InitialCap: 1,
		MaxCap:     2,
		Factory:    factory,
		Close:      close,
	}
	p, err := pool.NewPool(poolConfig)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func client(num int, p pool.Pool) {
	//从连接池中取得一个连接
	fmt.Println("num:", num, " Get()")
	v, err := p.Get()
	if err != nil {
		log.Fatal(err)
	}

	//do something
	cn := v.(net.Conn)
	fmt.Println("num:", num, " Got the connction:", cn.LocalAddr())

	//将连接放回连接池中
	fmt.Println("num:", num, " Put()")
	err = p.Put(v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("End")
}

func server() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
	defer l.Close()
	fmt.Println("Listening on ", addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		}
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
	}
}
