package http

import(
	"g"
	"log"
	"time"
	"net/http"
	"encoding/json"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RenderJson(w http.ResponseWriter,v interface{}) {
	bs,err:=json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}


func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}

	RenderDataJson(w, data)
}



func init(){
	configProxyRoutes()
}


func Start(){
	
	if !g.Config().Http.Enable{
		return 
	} 

	addr := g.Config().Http.Listen
	if addr == "" {
		return 
	}

	server := &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err:=server.ListenAndServe()
	if err != nil{
		log.Println("http server error:",err)
	}
	log.Println("listening: ",addr)

}