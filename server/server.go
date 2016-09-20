package server

import (
	"fmt"
	"g"
	"net"
	"time"
)

type ServerError struct {
	error
}

//重写Error()方法
func (this *ServerError) Error() string {
	return "conn is error"
}

type TcpConnPool struct {
	Dial      func(addr string) (interface{}, error)
	MaxIdle   int
	MaxActive int
	Close     bool
	Active    int
	//idle chan map[string]interface{}
	Idle map[string]chan interface{}
}

type IdleConn struct {
	c interface{}
	t time.Time
}



var TcpConn TcpConnPool = newTcpPool()

func (this *TcpConnPool) InitTcpConnPool() error {
	fmt.Println("=====InitTcpConnPool=========-")
	agents := g.Config().Reciver.Client.Agent
	num := 0
	this.Idle = make(map[string]chan interface{}, len(agents))
	for _, ag := range agents {
		fmt.Println(ag)
		n,err:=this.InitTcpSigleConnPool(ag)
			fmt.Println("==============-")
		if err != nil || n  == g.Config().Reciver.Client.Agent_Process {
			fmt.Println("%s no success connection!!!")
			num++
		}
	}

	if num == len(agents) {
		return &ServerError{}
	}
	fmt.Println("=====InitTcpConnPool=========-")
	return nil
}




func (this *TcpConnPool) InitTcpSigleConnPool(agent string) (int,error) {
	fmt.Println("=====InitTcpSigleConnPool=========-")
	t := time.Now()
	n:=0
	agchan := make(chan interface{}, g.Config().Reciver.Client.Agent_Process)
	for i := 0; i < g.Config().Reciver.Client.Agent_Process; i++ {
		conn, err := this.Dial(agent)
		if err != nil {
			fmt.Println("create tcp connection err to %s:", agent, err)
			continue
			n++
		}

		agchan <- IdleConn{c: conn, t: t}
	}
	if len(agchan) > 0 {
		this.Idle[agent] = agchan	
	}
	fmt.Println(this.Idle[agent])
	fmt.Println("=====InitTcpSigleConnPool=========-")
	return n,nil
}




func (this *TcpConnPool) GetConn(addr string) (net.Conn, error) {
	//最简单的版本
	fmt.Println("=====GetConn=========-")
	if this.Idle == nil {
		fmt.Println("连接池为空，重新初始化")
		err := this.InitTcpConnPool()
		if err != nil {
			fmt.Println("some tcp conn is err,suggestion to exit the app or change the config")
		}
		fmt.Println(this)
	}
	
	
	if len(this.Idle[addr]) == 0 {
		fmt.Println("因为没有连接了，所以重新初始化连接")
		this.InitTcpSigleConnPool(addr)
	}
	fmt.Println("%s 这个机器的连接数有多少个：",len(this.Idle[addr]))
	ic := <-this.Idle[addr]
	conn := ic.(IdleConn).c.(net.Conn)	
	fmt.Println("连接分别是：",conn)	
	fmt.Println("=====GetConn=========-")
	return conn, nil
}


func newTcpPool() TcpConnPool {
	fmt.Println("=====TcpConnPool=========-")
	return TcpConnPool{
		Dial: func(addr string) (interface{}, error) {
			conn, err := net.Dial("tcp", addr)
			return conn, err
		},
	}

}

func (this *TcpConnPool) ReleaseConn(addr string, conn net.Conn) {
	fmt.Println("+++++++++++++ReleaseConn++++++++++++++")
	t := time.Now()
	this.Idle[addr] <- IdleConn{t: t, c: conn}
	fmt.Println(len(this.Idle[addr]))
	fmt.Println("------------------ReleaseConn-------------")
}






func CheckConnAlive(conn net.Conn) (bool, error) {
	fmt.Println("=====checkConnAlive=========-")
	n,e:=conn.Write([]byte("ping ... I am alive"))
	fmt.Println(n)
	if e != nil {
		fmt.Println("连接已死")
		fmt.Println(n)
		return true,e
	}
	fmt.Println("---checkConnAlive---")
	fmt.Println(conn)
	fmt.Println("+++checkConnAlive+++")
	var num int64
	var timeout chan bool = make(chan bool,1)
	var bfer chan int = make(chan int,1)
	go func(){
		var ts int64 
		to:=time.Now().Unix()
		for{
			if num > 0 {
				timeout <- true
				break
			}
		    ts=time.Now().Unix()
		    if (ts - to) > 2{
		    	timeout <- false
		    	break
		    }
		}
	}()
	go func(){
		for {
			bf := make([]byte, 100)
			num, err:= conn.Read(bf)
			if err != nil {
				bfer <- 2
				break 
			}

			if string(bf[:num]) == "ping ... I am alive"{
				bfer <- 1
				break 
			}		
		}
	}()
	
	select{

		case timeo:= <-timeout:
			if timeo {
				fmt.Println("---timeo == 1 ---")
				return false,nil
				

			}else{
				fmt.Println("---timeo == 0 ---")
				return true ,&ServerError{}
			}

		case br := <-bfer:
			if br == 2 {
				fmt.Println("---br == 2---")
				return true ,&ServerError{}
			}else{
				fmt.Println("---br == 1---")
				return false,nil
			}
	}


}
