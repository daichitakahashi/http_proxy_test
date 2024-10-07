package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
)

func main() {
	addr := os.Args[1]
	os.Unsetenv("HTTPS_PROXY")

	go http.ListenAndServe(":8888", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		log.Println("fake: received request:", string(dump))
		w.Write([]byte("Hello from fake server"))
	}))

	go func() {
		backend := s3mem.New()
		fake := gofakes3.New(backend) // , gofakes3.WithHostBucket(true))
		http.ListenAndServe(":7777", fake.Server())
	}()

	h := goProxy()

	http.ListenAndServe(addr, h)
}

func goProxy() http.Handler {
	s := goproxy.NewProxyHttpServer()
	s.Verbose = true

	google := s.OnRequest(goproxy.UrlMatches(regexp.MustCompile("google.com:443")))
	google.HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return goproxy.MitmConnect, host
	})
	google.DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		r := req.Clone(req.Context())
		r.URL, _ = url.Parse(fmt.Sprintf("http://localhost:8888%s?%s", r.URL.Path, r.URL.Query().Encode()))
		r.RequestURI = ""
		ctx.Logf("original: %s", req.URL.String())
		ctx.Logf("proxied: %s", r.URL.String())
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			log.Println("Proxy error:", err)
		}
		return r, resp
	})

	s3Rx := regexp.MustCompile(`([^\.]*)\.?s3\.ap-northeast-1\.amazonaws\.com:443`)
	s3 := s.OnRequest(goproxy.UrlMatches(s3Rx))
	s3.HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return goproxy.MitmConnect, host
	})
	s3.DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		r := req.Clone(req.Context())
		if matches := s3Rx.FindStringSubmatch(req.URL.Host); matches[1] != "" {
			r.URL, _ = url.Parse(fmt.Sprintf("http://localhost:7777%s?%s", path.Join("/", matches[1], r.URL.Path), r.URL.Query().Encode()))
		} else {
			r.URL, _ = url.Parse(fmt.Sprintf("http://localhost:7777%s?%s", r.URL.Path, r.URL.Query().Encode()))
		}
		r.RequestURI = ""
		ctx.Logf("original: %s", req.URL.String())
		ctx.Logf("proxied: %s", r.URL.String())
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			log.Println("Proxy error:", err)
		}
		return r, resp
	})

	return s
}
