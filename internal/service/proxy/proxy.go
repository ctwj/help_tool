package proxy

import (
	"fmt"
	proxyBiz "helptools/internal/biz/proxy"
	"log"
)

type Proxy struct {
}

var proxy *proxyBiz.TransparentProxy

func NewProxy() *Proxy {
	return &Proxy{}
}

var proxyInitialized bool

func (p *Proxy) StartProxy() {
	fmt.Println("Start Proxy")
	if !proxyInitialized {
		proxy = proxyBiz.NewTransparentProxy("192.168.2.78:7890")
		err := proxy.Start(":80", ":443")
		if err != nil {
			log.Fatalf("Failed to start proxy server: %v", err)
		}
		proxyInitialized = true
	} else {
		fmt.Println("Proxy already initialized")
	}
}

func (p *Proxy) StopProxy() {
	proxy.Stop()
}
