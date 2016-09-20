package proxy

import(
	"g"
	"fmt"
	"log"
	"time"
	"storageCache"
	"encoding/json"
)


var(
	ProxyContent chan *ProxyT = make(chan *ProxyT ,5000)	
)


type ProxyT struct {
	Nfspath 	string
	Filepath	string
	Filename	string
	TimeStamp 	time.Time 
}

//var Proxyts *ProxyT = &ProxyT{Nfspath:"",Filepath:"",Filename:""} //进入chan之前的  


func ProxyAddMessage(nfspath string ,filepath string ,filename string ){
	//锁住进程
	var Proxyts *ProxyT = &ProxyT{}
	Proxyts.Nfspath=nfspath
	Proxyts.Filename=filename
	Proxyts.Filepath=filepath
	Proxyts.TimeStamp=time.Now()
	ProxyContent <- Proxyts
}

func ProxyStorage(){

	if g.Config().Rabbitmq.Enable {
		var rmqcon *storageCache.RabbitConn =storageCache.Rmqconn
		for{
			if len(ProxyContent) > 0{
				fmt.Println(">>")
				go writeRabbitMq(ProxyContent,rmqcon)
				log.Println("rabbitmq message in")	
				time.Sleep(time.Second * time.Duration(len(ProxyContent)) /100)
			}
			time.Sleep(time.Second*10)
		}	
		
	}

}


func writeRabbitMq(proxy chan *ProxyT,rmqcon *storageCache.RabbitConn)  error {
	fmt.Println("==writeRabbitMq=")
	channel,err:=rmqcon.GetChannel()
	if err != nil {
		fmt.Println("wrabbitmq channel get err:",err)
	}
	for pro := range proxy{


		b,err:=json.Marshal(pro)
		fmt.Println(string(b))
		if err != nil {
			fmt.Println("json encode error:",err)
			return err
		}
		go func(){

			err=storageCache.SendRabbitMq(b,channel)
			if err != nil{
				fmt.Println("storageCache.rabbitmq err :",err)
			}
		}()
		//这里调用rabbitmq的连接函数并将消息写进去


		fmt.Println("==writeRabbitMq=")

	}
	fmt.Println("==writeRabbitMq=")
	return nil
}
