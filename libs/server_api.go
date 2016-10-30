package libs

import (
	"net/http"
	"time"
	"strings"
	"encoding/json"
	"fmt"
)

type customHandler struct{

}

func(cb *customHandler) ServeHTTP( w http.ResponseWriter, r *http.Request ) {

	switch {
	case strings.HasPrefix(r.URL.Path, "/all") :
		b , e := json.Marshal(GlobalAllocates)

		if e == nil {
			w.Write(b);
		}else {
			Log.Infof("server api error : %v",e)
		}
	case strings.HasPrefix(r.URL.Path, "/token/"):
		fmt.Println("customHandler!!");
		w.Write([]byte("customHandler!!"));
	default:
		fmt.Println("customHandler!!");
		w.Write([]byte("customHandler!!"));

	}

}

func RunServerApi() {
	Log.Infof("run server api")
	var server *http.Server = &http.Server{
	Addr:           ":3838",
	Handler:        &customHandler{},
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 <<20,
	}
	go server.ListenAndServe();
}