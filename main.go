package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

func main() {
	httpsProxy, err := url.Parse(os.Getenv("HTTPS_PROXY"))
	if err != nil {
		panic(err)
	}
	defer launchProxy(":" + httpsProxy.Port())()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	log.Println("GET https://example.com")
	resp1, err := http.Get("https://example.com")
	if err != nil {
		panic(err)
	}
	body1, _ := io.ReadAll(resp1.Body)
	log.Println(string(body1))

	fmt.Println()

	log.Println("GET https://google.com/foo?search=bar")
	resp2, err := http.Get("https://google.com/foo?search=bar")
	if err != nil {
		panic(err)
	}
	body2, _ := io.ReadAll(resp2.Body)
	log.Println(string(body2))
}

func launchProxy(addr string) func() {
	cmd := exec.Command("./_proxy", addr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 500)
	return func() { cmd.Process.Kill() }
}
