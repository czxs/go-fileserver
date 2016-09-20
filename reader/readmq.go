package reader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"g"
	"log"
	"server"
	"storageCache"
	"time"
	"reflect"
	"github.com/streadway/amqp"
)

type Proxyfile struct {
	Nfspath   string
	Filepath  string
	Filename  string
	TimeStamp time.Time
}

func ServerInit() {
	//start to read rabbitmq and sendfile
	fmt.Println("start======ServerInit ====read server ...")
	var rmqcon *storageCache.RabbitConn = storageCache.Rmqconn

	fmt.Println("start======ServerInit ====read server ...")
	fmt.Println("sd")
	go readRabbitmq(rmqcon)
	


}

func readRabbitmq(rmqcon *storageCache.RabbitConn) {
	fmt.Println("start======readRabbitmq ====read server ...")
	fmt.Println(g.Config().Rabbitmq)
	channel, err := rmqcon.GetChannel()
	if err != nil {
		fmt.Println("wrabbitmq channel get err:", err)
	}
	defer channel.Close()

	msg, err := channel.Consume(
		g.Config().Rabbitmq.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("%s:%s", "readRabbitmq consume err:", err)
	}
	forever := make(chan bool)
	go func(){
		for d := range msg {
			fmt.Println(len(msg))
			fmt.Println(reflect.TypeOf(d))
			var file *Proxyfile
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(time.Second * t)
			err := json.Unmarshal(d.Body, &file)
			fmt.Println("读取消息%s",d.Body)
			if err == nil {
				go func(locald amqp.Delivery,f *Proxyfile){
					alen:=make(chan string,len(g.Config().Reciver.Client.Agent))
					clen:=make(chan bool,len(g.Config().Reciver.Client.Agent))
					fmt.Println("read message from rabbitmq success!!!")
					fmt.Println(time.Now())
					log.Printf("[x]%s", f)
					//stat := 0
					for _, destination := range g.Config().Reciver.Client.Agent {
						go SendFileToDest(destination, f,alen,clen)
					}
					for {
						if len(alen) == len(g.Config().Reciver.Client.Agent) {//表示所有机器都传完了
							if len(clen) == len(g.Config().Reciver.Client.Agent){ //表示所有机器都传正确了
								fmt.Println("消息确认")
								fmt.Println(len(alen))
								fmt.Println(len(clen))
								locald.Ack(false)
							}else{
								fmt.Println(len(alen))
								fmt.Println(len(clen))
								fmt.Println("消息没确认")
								locald.Nack(false, true)
							}
							break
						}
					}
					log.Printf("Done")
				}(d,file)
			} else {
				fmt.Println("message from rabbitmq error:", err)
			}
			fmt.Println("消息处理完%s",d.Body)
		}
	}()
	time.Sleep(time.Second * time.Duration(len(msg)) / 10)

	fmt.Println("end======readRabbitmq ====read server ...")
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
    <-forever
}

func SendFileToDest(destination string, file *Proxyfile,alen chan string,clen chan bool) bool {
	fmt.Println("start======SendFileToDest ====read server ...")

	conn_destination, err := server.TcpConn.GetConn(destination)
	fmt.Println(file.Filepath, file.Filename, destination, conn_destination)
	if err != nil {
		fmt.Println("创建链接有问题", err)
		fmt.Println("fults",destination)
		alen <- destination
		return false
	}

	death,err:=server.CheckConnAlive(conn_destination)
	if err != nil || death{
		alen <- destination
		fmt.Println("dssds",destination)
		return false 
	}

	ok := server.SendFile(file.Nfspath, file.Filepath, file.Filename, conn_destination)
	
	if ok {
		fmt.Println(destination)
		fmt.Println("%s 文件发送成功",file.Filename,ok)
		clen <- ok
	}
	alen <- destination
	fmt.Println("HERE")
	server.TcpConn.ReleaseConn(destination, conn_destination)
	fmt.Println(len(server.TcpConn.Idle[destination]))
	fmt.Println("文件发送失败")
	fmt.Println("--------------SendFileToDest-------------------")
	return ok
}
