package eren

import (
	"fmt"
	"net/url"
	"path"
	"testing"
)

func TestPath(t *testing.T){

	join := path.Join("/", "fuck/")
	fmt.Printf(join)
}

//func TestIp(t *testing.T){
//
//	server := NewServer()
//	server.GET("/a", func(context *Context) {
//		context.Status(200)
//	})
//
//	http.ListenAndServe(":8080",server)
//
//}

func TestUrl(t *testing.T){
	escape:= url.QueryEscape("/a?id=10")
	fmt.Printf(escape)
}