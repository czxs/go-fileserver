package storageCache

import(
	"g"
	"fmt"
	"time"
	"github.com/streadway/amqp"
)


type RabbitConn struct{
	// Dial is an application supplied function for creating and configuring a
    // connection
    Dial func() (interface{}, error)
    //Dial func() (*autorc.Conn, error)
    // Maximum number of idle connections in the pool.
    MaxIdle int
    // Maximum number of connections allocated by the pool at a given time.
    // When zero, there is no limit on the number of connections in the pool.
    MaxActive int
    closed bool
    active int
    idle chan interface{}
}

type IdleConn struct{
	c interface{}
	t time.Time
}

//初始化生成批量的rabbitmqconn
func (this *RabbitConn) InitPool() error {
	this.idle=make(chan interface{},this.MaxActive)
	for x:=0;x<this.MaxActive;x++{
		conn,err:=this.Dial()
		if err!=nil{
			return err 
		}
		this.idle <- IdleConn{c:conn,t:time.Now()}
	}
	return nil
}

func (this *RabbitConn)Get() interface{} {
    // 如果空闲连接为空，初始化连接池
    if this.idle == nil {
        this.InitPool()
    }
    // 赋值一下好给下面回收和返回
    //conn := <-this.idle
    //idleConn
    ic := <-this.idle
    // 这里要用 (idleConn) 把interface{} 类型转化为 idleConn 类型的，否则拿不到里面的属性t、c
    conn := ic.(IdleConn).c
    //fmt.Println(conn.(*DB).conn)
    //fmt.Println(" --- reflect --- ", reflect.TypeOf(conn))
    // 使用完把连接回收到连接池里
    defer this.Release(conn)
    // 因为channel是有锁的，所以就没必要借助sync.Mutex来进行读写锁定
    // container/list就需要锁住，不然并发就互抢出问题了
    return conn
}

// 回收连接到连接池
func (this *RabbitConn)Release(conn interface{}) {
    //this.idle <-conn
    this.idle <-IdleConn{t: time.Now(), c: conn}
}

func newRmqPool() *RabbitConn {
    poolNum:= 500
    fmt.Printf("初始化 rabbitmq 连接池，连接数：%d \n", poolNum)
    return &RabbitConn{
        MaxActive: poolNum,
        //Dial: func() (*autorc.Conn, error) {
        Dial: func() (interface{}, error) {
        	conn,err:=amqp.Dial("amqp://"+g.Config().Rabbitmq.User+":"+g.Config().Rabbitmq.User+"@"+g.Config().Rabbitmq.S_addr)
            return conn, err
        },
    } 
}

func (rmqcon *RabbitConn) GetChannel()(*amqp.Channel,error){
	
	conn:=rmqcon.Get().(*amqp.Connection)

	channel, err := conn.Channel()

	if err != nil {
		fmt.Println("channel err:",err)
		return channel,err
	}
	return channel,nil
}


var Rmqconn *RabbitConn =newRmqPool()
