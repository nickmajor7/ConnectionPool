package pool

import "errors"

var (
	ErrInvalidCapacity    = errors.New("invalid capacity settings")
	ErrInvalidFactoryFun  = errors.New("invalid factory func settings")
	ErrInvalidCloseFun    = errors.New("invalid close func settings")
	ErrOpenNumber         = errors.New("numOpen > maxOpen")
	ErrConnIsNil          = errors.New("connection is nil. rejecting")
	ErrPoolClosed         = errors.New("pool is closed")
	ErrPoolClosedAndClose = errors.New("connction pool is closed. close connection")
)

// Config 连接池相关配置
type Config struct {
	//连接池中初始化的连接数
	InitialCap int
	//连接池中拥有的最大的连接数(默认值0，无限制)
	MaxCap int
	//生成连接的方法
	Factory func() (interface{}, error)
	//关闭连接的方法
	Close func(interface{}) error
	//检查连接是否有效的方法
	Ping func(interface{}) error
}

// Pool 基本方法
type Pool interface {
	Get() (interface{}, error)

	Put(interface{}) error

	Ping(interface{}) error

	Close(interface{}) error

	Release()
}
