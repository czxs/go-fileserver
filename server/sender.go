package server

import (
	"crypto/md5"
	"fmt"
	"g"
	"io"
	"net"
	"os"
	"time"
)

type ArithmeticError struct {
	error
}

//重写Error()方法
func (this *ArithmeticError) Error() string {
	return "conn is error"
}

func getFileInfo(nfspath string, filepath string, filename string) *g.SysFileInfo {
	fmt.Println("start======getFileInfo ====read server ...")
	filep := nfspath + "/" + filepath + "/" + filename
	fi, err := os.Lstat(filep)
	if err != nil {
		fmt.Println("info ERROR", err)
		return nil
	}
	fileHandle, err := os.Open(filep)
	defer fileHandle.Close()
	if err != nil {
		fmt.Println("open ERROR", err)
		return nil
	}

	h := md5.New()
	_, err = io.Copy(h, fileHandle)
	fileInfo := &g.SysFileInfo{
		FName:  fi.Name(),
		FSize:  fi.Size(),
		FPerm:  fi.Mode().Perm(),
		FMtime: fi.ModTime(),
		FType:  fi.IsDir(),
		FMd5:   fmt.Sprintf("%x", h.Sum(nil)),
		FPath:  filepath,
	}
	fmt.Println(fileInfo)
	fmt.Println("start======getFileInfo ====read server ...")
	return fileInfo
}

func SendFile(nfspath string, filepath string, filename string, conn net.Conn) bool {
	/*
		death, err := checkConnAlive(conn)
		if err != nil || death {
			fmt.Println("the conn is death")
			return &ArithmeticError{}
		}*/
	fmt.Println("+++++getFileInfo ++++++")
	filep := nfspath + "/" + filepath + "/" + filename
	fInfo := getFileInfo(nfspath, filepath, filename)
	newName := fmt.Sprintf("%s", fInfo.FName)
	cmdLine := fmt.Sprintf("upload %s %d %d %d %s %s %s", newName, fInfo.FMtime.Unix(), fInfo.FPerm, fInfo.FSize, fInfo.FMd5, fInfo.FPath, "\r\n")
	conn.Write([]byte(cmdLine))
	fileHandle, err := os.Open(filep)
	defer fileHandle.Close()
	if err != nil {
		fmt.Println("open file error:", err)
		return false
	}
	n, err := io.Copy(conn, fileHandle)
	if err != nil {
		fmt.Println("io copy的复制出问题了", n)
		fmt.Println("io copy的复制出问题了", err)
	}
	for {
		buffer := make([]byte, 1024)
		num, err := conn.Read(buffer)
		if err == nil && num > 0 && (string(buffer[:num]) == "sync success" || string(buffer[:num]) == "syncing success") {
			fmt.Println(string(buffer[:num]))

			return true
		} else if err != nil || string(buffer[:num]) == "sync failed" {
			fmt.Println(string(buffer[:num]))
			fmt.Println(err)
			time.Sleep(time.Second * 10)
			return false
		}

		fmt.Println(string(buffer[:num]))
	}
	fmt.Println("-----getFileInfo-------")
	return true
}
