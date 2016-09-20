package client

import (
	"fmt"
	"g"
	"io"
	"net"
	"os"
	"time"
)

var (
	Proxyserver bool = false
)

func StartClient() {
	serverport := fmt.Sprintf("%s", g.Config().Reciver.Server.Listen)
	listen, err := getNewListen(serverport)
	if err != nil {
		fmt.Println("获取listener失败！")
	}
	fmt.Println("start client")
	client(listen)
}

func getNewListen(serverport string) (net.Listener, error) {

	listen, err := net.Listen("tcp", serverport)
	if err != nil {
		fmt.Println("net listen failed", err)
		return nil, err
	}

	//检查指定的根目录是否存在
	fi, err := os.Stat(g.Config().Reciver.Server.Dirpath)
	if err != nil {
		fmt.Println("the base dir is not exist!", err)
	} else if !fi.IsDir() {
		fmt.Println("the path is not dir!", err)
	}
	fmt.Println("return listen success")
	return listen, nil
}

func client(listen net.Listener) {
	var Nfspath string = g.Config().Reciver.Server.Dirpath
	for {
		conn, err := listen.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() { //判断错误是否是临时网络有问题
				continue
			}
			fmt.Println("network error", err)
		}
		fmt.Println("go handler")

		go Handler(conn, Nfspath)
	}

}

func Handler(conn net.Conn, Nfspath string) {
	var tempfilename string
	fileheader := true
	var lastsize int
	var cmd g.SysFileInfo
	t := time.Now().Unix()
	fmt.Println(conn)
	fmt.Println(t)
	var fp *os.File
	for {
		
		for {
			buffer := make([]byte, 2048)
			var n int64
			num, err := conn.Read(buffer)
			t = time.Now().Unix()
			if err != nil && err != io.EOF {
				fmt.Println("can not read file", err)
				return 
			}
			if string(buffer[:num]) == "ping ... I am alive" {
				conn.Write([]byte("ping ... I am alive"))
				fmt.Println(string(buffer[:num]))
				fmt.Println("this is a heartbeat check package,not file data")
				time.Sleep(time.Second * 5)
				break
			}

			if num == 0 { //防止进程异常中断
				return
			}

			if fileheader {
				fmt.Println("fileheader is true and start to parse")
				n, cmd = cmdParse(buffer[:num])
				//if n > 0 {
				fii,err:= os.Stat(Nfspath+"/"+cmd.FPath)
				if err != nil {
					//创建目录
					os.MkdirAll(Nfspath+"/"+cmd.FPath,0755)
				}else{
					if ! fii.IsDir() {
						os.MkdirAll(Nfspath+"/"+cmd.FPath,0755)
					}
				}
				tempfilename = fmt.Sprintf("%s/%s/%s.newsync.%s", Nfspath, cmd.FPath, cmd.FName,cmd.FMd5)

				//判断传输的临时文件是否已经存在
				_,err = os.Stat(tempfilename)
				if err == nil || os.IsExist(err) {
					fmt.Println("file is reciving_----------------------------------------=---------+-----------+	")
					sendInfo := fmt.Sprintf("syncing success")
					conn.Write([]byte(sendInfo))
					time.Sleep(time.Second * 10)
				}

				lastsize = int(cmd.FSize)
				err,fp = writeToFile(buffer[int(n):int(num)], tempfilename, cmd.FPerm)
				if err != nil {
					fmt.Println("write file error :", err)
				}
				lastsize = lastsize - num + int(n)
				fileheader = false
				continue

			} else {
				fmt.Println(cmd)
				err,fp = writeToFile(buffer[:int(num)], tempfilename, cmd.FPerm)
				if err != nil {
					fmt.Println("write file error :", err)
				}
				lastsize = lastsize - num
			}
			if lastsize > 0 {
				fileheader = false
				continue
			}

			if lastsize < 0 {
				fmt.Println(string(buffer))
			}

			err = setFileAttr(Nfspath, cmd.FPath, cmd.FName,cmd.FMd5, cmd.FMtime)
			if err != nil {
				fmt.Println("设置文件属性出问题了")
			}

			newfMd5, err := getMd5Info(Nfspath + "/" + cmd.FPath + "/" + cmd.FName)
			if err != nil {
				fmt.Println("获取文件的MD5值过程有问题:", err)
			}

			if newfMd5 == cmd.FMd5 {
				sendInfo := fmt.Sprintf("sync success")
				conn.Write([]byte(sendInfo))
				time.Sleep(time.Second * 10)
			} else {
				sendInfo := fmt.Sprintf("sync failed")
				fmt.Println("___+++")
				conn.Write([]byte(sendInfo))
				time.Sleep(time.Second * 10)
			}
			fmt.Println("文件传输完成")
			fileheader = true
			fp.Close()
			break
		}
	}
}
