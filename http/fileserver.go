package http

import(
	"fmt"
	"net/http"
	"proxy"
)

func configProxyRoutes(){
	http.HandleFunc("/fileserver",func(w http.ResponseWriter,r *http.Request){
		//调用mqproxy去处理
		if r.Method != "POST"{
			w.Write([]byte(string(http.StatusMethodNotAllowed)))
		}


		r.ParseForm()
		var nfspath  string = r.PostFormValue("nfspath")
		var filename string = r.PostFormValue("filename")
		var filepath string = r.PostFormValue("filepath")
		if nfspath == "" || filename == "" || filepath == "" {
			w.Write([]byte("three args : nfspath,filepath,filename ,any one can not be null"))
			fmt.Println(nfspath)
			fmt.Println(filename)
			fmt.Println(filepath)
			return 
		}

		go proxy.ProxyAddMessage(nfspath,filepath,filename)
		w.Write([]byte("/fileserver?status=ok"))

	})
}