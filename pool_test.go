package pool

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"testing"
	"time"
)

type DemoCloser struct {
	Name     string
	activeAt time.Time
}

func (p *DemoCloser) Close() error {
	fmt.Println(p.Name, "closed")
	return nil
}

func (p *DemoCloser) GetActiveTime() time.Time {
	return p.activeAt
}

const addr string = "127.0.0.1:65432"

var serverStart chan int
var connectionCount chan int
var closePool chan int

func server() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
	defer l.Close()
	// fmt.Println("Listening on ", addr)
	serverStart <- 1
	for {
		_, err := l.Accept()
		if err != nil {
			log.Fatal("Error accepting: ", err)
		}
		connectionCount <- 1
	}
}

func TestNewPool(t *testing.T) {
	poolConfig := &Config{
		InitialCap: 1,
		MaxCap:     5,
		Factory:    func() (interface{}, error) { return net.Dial("tcp", addr) },
		Close:      func(v interface{}) error { return v.(net.Conn).Close() },
	}
	_, err := NewPool(poolConfig)
	if err != nil {
		t.Error(err)
		return
	}

	count := 0
	timeCount := 0
	for count < poolConfig.InitialCap && timeCount < 10 {
		select {
		case _, ok := <-connectionCount:
			if ok {
				count++
			} else {
				t.Error("Channel closed!") //Channel 被close.
			}
		case <-time.After(time.Second * 1):
			t.Log("tick..")
			timeCount++
		}
	}
}

func TestGetPub(t *testing.T) {
	poolConfig := &Config{
		InitialCap: 1,
		MaxCap:     5,
		Factory:    func() (interface{}, error) { return net.Dial("tcp", addr) },
		Close:      func(v interface{}) error { return v.(net.Conn).Close() },
	}
	p, err := NewPool(poolConfig)
	if err != nil {
		t.Error(err)
		return
	}

	count := 0
	timeCount := 0
	for count < poolConfig.InitialCap && timeCount < 10 {
		select {
		case _, ok := <-connectionCount:
			if ok {
				count++
			} else {
				t.Error("Channel closed!") //Channel 被close.
			}
		case <-time.After(time.Second * 1):
			t.Log("tick..")
			timeCount++
		}
	}

	var connArray []net.Conn

	for i := 0; i < poolConfig.MaxCap; i++ {
		s, err := p.Get()
		if err != nil {
			t.Fatal(err)
			return
		}
		connArray = append(connArray, s.(net.Conn))
	}

	for count < poolConfig.MaxCap-poolConfig.InitialCap && timeCount < 10 {
		select {
		case _, ok := <-connectionCount:
			if ok {
				count++
			} else {
				t.Error("Channel closed!") //Channel 被close.
			}
		case <-time.After(time.Second * 1):
			t.Log("tick..")
			timeCount++
		}
	}

	for s := range connArray {
		err = p.Put(s)
		if err != nil {
			t.Fatal(err)
			return
		}
	}
}

// func TestPool_Get(t *testing.T) {
// 	pool, err := NewGenericPool(0, 5, time.Minute*10, func() (Poolable, error) {
// 		time.Sleep(time.Second)
// 		name := strconv.FormatInt(time.Now().Unix(), 10)
// 		log.Printf("%s created", name)
// 		// TODO: FIXME &DemoCloser{Name: name}后，pool.Acquire陷入死循环
// 		return &DemoCloser{Name: name, activeAt: time.Now()}, nil
// 	})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	for i := 0; i < 10; i++ {
// 		s, err := pool.Acquire()
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		pool.Release(s)
// 	}
// }

// func TestGenericPool_Shutdown(t *testing.T) {
// 	pool, err := NewGenericPool(0, 10, time.Minute*10, func() (Poolable, error) {
// 		time.Sleep(time.Second)
// 		return &DemoCloser{Name: "test"}, nil
// 	})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if err := pool.Shutdown(); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if _, err := pool.Acquire(); err != ErrPoolClosed {
// 		t.Error(err)
// 	}
// }

func TestMain(m *testing.M) {
	fmt.Println("begin")
	runtime.GOMAXPROCS(8)
	serverStart = make(chan int)
	connectionCount = make(chan int)
	closePool = make(chan int)
	go server()
	<-serverStart
	m.Run()
	fmt.Println("end")
}
