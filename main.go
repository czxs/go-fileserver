package main
import (
	"g"
	"fmt"
	"flag"
	"http"
	"client"
	"proxy"
	"reader"
)


var (
	cfg = flag.String("c","cfg.json","configuration file")
	syncSer = flag.Bool( "server",false,"server mode")
	syncCli = flag.Bool( "client",false,"server mode")
	syncPro = flag.Bool( "proxy",false,"server mode")
)





func main(){


	flag.Parse()

	g.ParseConfig(*cfg)

	forever := make(chan bool)
	if *syncSer || g.Config().Reciver.Client.Enable {
		fmt.Println("start http server and fileserver server ...")
		go http.Start()	
		go proxy.ProxyStorage()
		go reader.ServerInit()
	}

	if *syncCli ||g.Config().Reciver.Server.Enable {
		fmt.Println("start fileserver client ...")
		if *syncPro || g.Config().Reciver.Server.Proxy{
			client.Proxyserver = true
			fmt.Println("start proxyserver ...")
		}
		go client.StartClient()
	}
	
	<-forever


}