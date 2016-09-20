package client

import (
	"crypto/md5"
	"fmt"
	"g"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func setFileAttr(nfspath string, filepath string, filename string,filemd5 string ,mtime time.Time) error {
	fmt.Println(nfspath)
	fmt.Println(filepath)
	fmt.Println(filename)
	os.Remove(nfspath + "/" + filepath + "/" + filename) //删除当前文件
	tempfilename := fmt.Sprintf("%s.newsync.%s", nfspath+"/"+filepath+"/"+filename,filemd5)
	err := os.Rename(tempfilename, nfspath+"/"+filepath+"/"+filename)
	if err != nil {
		fmt.Println("rename ", tempfilename, " to ", nfspath+"/"+filepath+"/"+filename, " failed", err)
		return err
	}

	fmt.Println("here")
	err = os.Chtimes(nfspath+"/"+filepath+"/"+filename, time.Now(), mtime)
	if err != nil {
		fmt.Println("change the mtime error ", err)
		return err
	}
	return nil
}

func getMd5Info(filename string) (string, error) {
	fileHandle, err := os.Open(filename)
	defer fileHandle.Close()
	if err != nil {
		fmt.Println("open the create file", filename, "failed:", err)
		return "", err
	}

	h := md5.New()
	_, err = io.Copy(h, fileHandle)
	if err != nil {
		fmt.Println("error")
		return "", err
	}
	newfMd5 := fmt.Sprintf("%x", h.Sum(nil))

	return newfMd5, nil
}

func writeToFile(buffer []byte, filename string, perm os.FileMode) (error,*os.File) {
	fmt.Println("start to write file")
	writeFile, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, perm)
	if err != nil {
		fmt.Println("write file error:", err)
		return err,writeFile
	}
	defer writeFile.Close()
	_, err = writeFile.Write(buffer)
	if err != nil {
		fmt.Println("write file error", err)
		return err,writeFile
	}
	fmt.Println("write file")
	return nil ,writeFile
}

func cmdParse(infor []byte) (int64, g.SysFileInfo) {
	var i int64
	var fileInfo g.SysFileInfo
	for i = 0; i < int64(len(infor)); i++ {
		if infor[i] == '\n' && infor[i-1] == '\r' {
			fmt.Println("start to parse")
			cmdLine := strings.Split(string(infor[:i-2]), " ")
			fmt.Println(cmdLine)
			fileName := fmt.Sprintf("%s", cmdLine[1])
			filePerm, _ := strconv.Atoi(cmdLine[3])
			fileMtime, _ := strconv.ParseInt(cmdLine[2], 10, 64)
			fileSize, _ := strconv.ParseInt(cmdLine[4], 10, 64)

			fileInfo = g.SysFileInfo{
				FName:  fileName,
				FMtime: time.Unix(fileMtime, 0),
				FPerm:  os.FileMode(filePerm),
				FSize:  fileSize,
				FMd5:   string(cmdLine[5]),
				FPath:  string(cmdLine[6]),
			}
			fmt.Println("---------------")
			fmt.Println(fileInfo)
			fmt.Println("===============")
			return i + 1, fileInfo

		}
	}
	return 0, fileInfo
}

func checkConnAlive(conn net.Conn) (bool, error) {
	t := time.Now().Unix()
	var n int
	var err error
	for {
		conn.Write([]byte("ping ... I am alive"))
		for {
			bf := make([]byte, 100)
			n, err = conn.Read(bf)
			if err != nil {
				return true, err
			}
			if int(time.Now().Unix())-int(t) > 300 {
				return true, nil
			}
		}
		if n > 0 {
			return false, nil
		}
	}

}
