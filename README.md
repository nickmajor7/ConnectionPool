# pool

Golang 实现的连接池


## 功能：

- 连接池中连接类型为`interface{}`，使得更加通用
- 支持用户设定 ping 方法，检查连接的连通性，无效的连接将丢弃

## 基本用法

```go

//factory 创建连接的方法
factory := func() (interface{}, error) { return net.Dial("tcp", "127.0.0.1:4000") }

//close 关闭连接的方法
close := func(v interface{}) error { return v.(net.Conn).Close() }

//ping 检测连接的方法
//ping := func(v interface{}) error { return nil }

//创建一个连接池： 初始化1，最大连接2
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

//从连接池中取得一个连接
v, err := p.Get()

//do something
//conn=v.(net.Conn)

//将连接放回连接池中
p.Put(v)

//释放连接池中的所有连接
p.Release()


```


## License

The MIT License (MIT) - see LICENSE for more details