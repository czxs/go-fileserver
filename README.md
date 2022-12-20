# go-fileserver
##目的：target
####解决运维的大量小文件实时分发的问题
      
##原理：schematic 
#####1.提供http接口，前端上传文件的时候给http接口提供需要传递我文件名，文件相对路径，文件存储的基本路径
#####2.proxy接收到http请求的信息之后，传递给配置的rabbitmq队列
#####3.server端连接rabbitmq的server，读取消息，按照消息提供的路径和名称，将文件读取。
#####4.server初始化一批长连接和client进行通讯。
#####5.server将文件读取之后通过socket长连接，将文件传送给client。
#####6.client接收文件

    
##下一个版本的功能：next version
#####1.client端检查自己是不是叶子节点，如果不是，将传递的文件字节流写本地的同时也传给下一级叶子节点。
#####2.client也初始化并维护一个与各个叶子节点的连接池
    
##安装：install 
#####go build 
#####./src 
    
## 配置：configure

client 配置：
    {
		"debug":true,     //打开调试模式<br>
		"rabbitmq":{     //mq的配置<br>
	    	"s_addr":"127.0.0.1:5672",  //地址和端口<br>
       	    "user":"dev",  //用户<br>
       		"pass":"dev",  //密码<br>
	        "enable":false,  //是否连接mq<br>
       		"exchange":"test", //exchange的名字<br>
       		"queue":"hello",  //队列的名字<br>
       		"routing_key":"hello_test" //route_key<br>
   		},
   		"reciver":{<br>
       		"client":{  //接收文件端<br>
           		"agent":["127.0.0.1:55556"],  //接收文件的地址，可以多个<br>
	                "enable":false,                //是否开启发送文件模块<br>
        	   	"agent_process":50            //初始化的长连接个数<br>
       		},  <br>
	       	"server":{                      //server 发送文件的地址<br>
        		"enable":true,              //是否开启接收文件模块<br>
           		"listen":":55556",           //接收文件模块监听的端口<br>
	           	"dirpath":"/data/static1",   //文件保存的基本路径<br>
        	   	"proxy":false              <br> //是否开启代理模式（功能没有开发完全，设想如果开启，就将文件从此节点转发给下一级子节点）<br>
	       	}   <br>
   		},<br>
		"http":{<br>
 	        "listen":"127.0.0.1:6082", //http代理监听的地址<br>
       		"enable":false              //是否开启http模块<br>
   		}<br>
    }<br>
server 配置
    {
		"debug":true,     ##打开调试模式<br>
		"rabbitmq":{     ##mq的配置<br>
			"s_addr":"127.0.0.1:5672",  ##地址和端口<br>
			"user":"dev",  ##用户<br>
			"pass":"dev",  ##密码<br>
			"enable":true,  ##是否连接mq<br>
			"exchange":"test", ##exchange的名字<br>
			"queue":"hello",  ##队列的名字<br>
			"routing_key":"hello_test" ##route_key<br>
		},
		"reciver":{<br>
			"client":{  ##接收文件端<br>
			"agent":["127.0.0.1:55556"],  ##接收文件的地址，可以多个<br>
			"enable":true,                ##是否开启发送文件模块<br>
			"agent_process":50            ##初始化的长连接个数<br>
			},<br>
			"server":{                      ##server 发送文件的地址<br>
				"enable":false,              ##是否开启接收文件模块<br>
				"listen":":55556",           ##接收文件模块监听的端口<br>
				"dirpath":"/data/static1",   ##文件保存的基本路径<br>
				"proxy":false    ##是否开启代理模式（功能没有开发完全，设想如果开启，就将文件从此节点转发给下一级子节点）<br>
			}<br>
		},<br>
		"http":{<br>
			"listen":"127.0.0.1:6082", ##http代理监听的地址<br>
			"enable":true              ##是否开启http模块<br>
		}<br>
    }<br>
## server 端需要安装rabbitmq
#####mq 新建队列和exchange 
#####创建routingkey  
##运行：run
#####client ./src
 
#####server ./src
    
    
