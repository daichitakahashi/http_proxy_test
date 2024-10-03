package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/elazarl/goproxy"
)

func main() {
	addr := os.Args[1]
	os.Unsetenv("HTTPS_PROXY")

	s := goproxy.NewProxyHttpServer()
	s.Verbose = true

	go http.ListenAndServe(":8888", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		log.Println("fake: received request:", string(dump))
		w.Write([]byte("Hello from fake server"))
	}))

	s.OnRequest(goproxy.DstHostIs("google.com:443")).HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		log.Println("proxy: HandleConnect", host)
		return goproxy.MitmConnect, host
	})
	s.OnRequest(goproxy.DstHostIs("google.com:443")).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		r := req.Clone(req.Context())
		r.URL, _ = url.Parse(fmt.Sprintf("http://localhost:8888%s?%s", r.URL.RawPath, r.URL.RawQuery))
		return r, nil
	})

	http.ListenAndServe(addr, s)
}
