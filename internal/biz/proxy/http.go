package proxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type Config struct {
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

const cert = `-----BEGIN CERTIFICATE-----
MIIDZzCCAk8CFAlfen+f/q9DLoHSXr8/aMBjgPuFMA0GCSqGSIb3DQEBCwUAMHAx
CzAJBgNVBAYTAmtyMQswCQYDVQQIDAJqcDELMAkGA1UEBwwCanAxDDAKBgNVBAoM
A2x4ZDENMAsGA1UECwwEc2FuZzENMAsGA1UEAwwEbm9uZTEbMBkGCSqGSIb3DQEJ
ARYMZm9yQHRlc3QuY29tMB4XDTIzMTIxNTA2NDYxNVoXDTMzMTIxMjA2NDYxNVow
cDELMAkGA1UEBhMCa3IxCzAJBgNVBAgMAmpwMQswCQYDVQQHDAJqcDEMMAoGA1UE
CgwDbHhkMQ0wCwYDVQQLDARzYW5nMQ0wCwYDVQQDDARub25lMRswGQYJKoZIhvcN
AQkBFgxmb3JAdGVzdC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQDtnxu/ShOOFkp6B9qOp9+79gHk8Ut/wXqQ5U4roNMHYOkOmqkp+fbNe36VNgvG
mUx2l3dDnbRj/j/mMceNxEGOHqTkm3BeCSUakkV//TFJUvJbMQyw5op2Fr3Ueprh
KX31JFDZoYVjbd7Nt36JDEtZjrIKuWKuRkXUpPXUeMrQTYIwiN86Dw/+TPdQhSXM
2uQ76xU+4IwKqaym1+uNoxy6q9NLdr35yDPr5PO0QIu2wBUwj37cJspDLmpZL4Jy
VOF5pLiIKfP7RbY7Aj9ec+QufvtAJxFEVYKiDDeCONyGtMtIh2A5+K1nsP9ivFuY
lUB+v/clgHrLXrJP4n79gvcLAgMBAAEwDQYJKoZIhvcNAQELBQADggEBABN7PDyI
mOsPI1i366f1EuZpF3SMl8rpMKTDoeNWVMSYchtSjPBE5ScgZBfvsnVu82qAUsnN
8oaOl0E0F7/ldC1rmjpbf/8TCQRr23IYrx9qk+QkqNNrqXSAbqNIYmLPcZCqbrpr
MqTCm7+fTZS8dNT/ZsAVoWWwlBnpNPG7St3RBu0Bq6J13PWdc4yj0i90YlH+M2sX
xglcppgrwJenVxtkaCIPCqP0468pTIWqTgq3hfZAqBjHjk9DXa+tTIj/VNlQ6ggP
7k+5pSUkoeuocHyxoRcS8UDKXFHqeWBole6z2n+6Bzaq0F1LnBSLbGgXDCuVS5N6
Q//41X8LTUT2iZk=
-----END CERTIFICATE-----`

const key = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDtnxu/ShOOFkp6
B9qOp9+79gHk8Ut/wXqQ5U4roNMHYOkOmqkp+fbNe36VNgvGmUx2l3dDnbRj/j/m
MceNxEGOHqTkm3BeCSUakkV//TFJUvJbMQyw5op2Fr3UeprhKX31JFDZoYVjbd7N
t36JDEtZjrIKuWKuRkXUpPXUeMrQTYIwiN86Dw/+TPdQhSXM2uQ76xU+4IwKqaym
1+uNoxy6q9NLdr35yDPr5PO0QIu2wBUwj37cJspDLmpZL4JyVOF5pLiIKfP7RbY7
Aj9ec+QufvtAJxFEVYKiDDeCONyGtMtIh2A5+K1nsP9ivFuYlUB+v/clgHrLXrJP
4n79gvcLAgMBAAECggEAQNbowGLbQStN+UyJU+H1CBoR3EIHgv3tjvozQW5qBxpn
neeP/3dI5DstiFOuFOwh1D5yec3fumVXVF4DAOkrpjcwcX0ExXQjXsPJGSqit0pd
/Yo910uhPqXn+MHX31buGuVk9m2/syj44hOPAKCNMwvgA1Mg3UMprOOyfN7VIM9u
8iRxkVCwhaj1qY7z3HbydLvoKB1IlAYsRPMQZ2L4sgsnHXy1PyVlOo7Phgvz6opl
sIoaRKijcjWoyKQlibbZ5SCTIY0iWWSlEfawButzpdHdWr5Jqyp+P70yHt39vEYF
8OyUiB+9pmi4OnYwNqLulXCQJWVmKwwQSTarhnD9qQKBgQD8idCj31+c5jSz2voJ
OMP3lHbRHH6qzywi23nUrLKJfkx+lyiI72VS8dmUVV2+8/L04gYG0cCLOFyIRFTZ
/U3CaEYyXmbLTpUS9uM4oKbp3OrwsfurIonq22NAoj7ItSJcn38ohlKhb9tozfTd
tbjFv/Tblmayp9PeEHTTJ4jAbwKBgQDw4PLXvNsVuZQWA9tmZXerDUJZsVtilK2X
JhfL3DXYdto5p1L2sxx05aiukJTs4JnnIetUsxswlgUhxAugxZauNODyYoEYZ6XE
Db2uwEU3zIKwSfK0X289NttT1l6l2/flVaq0Jg1n/7PS3DTddRqecTeCQeBXQawx
dwA8FtjJJQKBgQCrwJ8hlJ3We8qEN/2tn+nHzDUy6wpK6TO/UT1+oyWZ1Uf5IJz0
5LwouUudUqG7aPZoDgDDSoyFIwPruW1sBJaKDZkQUJvg0cUZbMgEj15110YCBUqA
jbD3BdZu8ul4X5jLHb7BtPklyomSseBDmX/dHjxNy/B0uSei89ZAdbbQCQKBgQCB
2MEPir6O93rcYzfh+tCHZJ5fuzuH6J2q3N33Br3/8hGxAoG2etbcLPDBKS8egfR0
o9Q31FTT3AroKMYb5GdVgSvBfVgZz8WL6dxWV074xUWtwi8TDF7qoKeaifR7dBgt
iAB9HAYeCbjl8c6NkpLG1kEV9mz4nG42O+/kdGxoAQKBgQDoprlF9NYXKPav0OYz
wGZ4WD6c6soW5JfDmEqBEGsEde1yFnDtT2PHJgvNImkwx5U0/YV8qO1a5jIN9c/u
a3JN6WVxLlgO77stKmzJ51kDK7voY45RFSBxZvoBYkTQHHkX6kiThggSGKtWnX+S
RFR9z9eT7VNxHRoZ7YwPIk3h8A==
-----END PRIVATE KEY-----`

type TransparentProxy struct {
	listenerHTTP  net.Listener
	listenerHTTPS net.Listener
	destination   string
	requestHook   RequestInterceptor // 添加 requestHook 字段
	responseHook  ResponseInterceptor
	proxyIPs      []string
}

type RequestInterceptor func(req *http.Request) *http.Request
type ResponseInterceptor func(req *http.Response) *http.Response

func NewTransparentProxy(destination string) *TransparentProxy {
	ips, _ := getAllLocalIPs()
	return &TransparentProxy{
		destination: destination,
		proxyIPs:    ips,
	}
}

func getAllLocalIPs() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			// 排除环回地址和非 IPv4 地址
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}

	return ips, nil
}

func (p *TransparentProxy) Start(httpPort, httpsPort string) error {
	httpListener, err := net.Listen("tcp", httpPort)
	if err != nil {
		return err
	}

	httpsConfig := &tls.Config{
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := tls.X509KeyPair([]byte(cert), []byte(key))
			if err != nil {
				fmt.Println("Failt to load certificate:", err)
				return nil, err
			}
			return &cert, nil
		},
	}
	httpsListener, err := tls.Listen("tcp", httpsPort, httpsConfig)
	if err != nil {
		httpListener.Close()
		return err
	}

	p.listenerHTTP = httpListener
	p.listenerHTTPS = httpsListener

	go p.acceptConnections(p.listenerHTTP, false)
	go p.acceptConnections(p.listenerHTTPS, true)

	log.Printf("Proxy server started. HTTP port: %s", httpPort)
	log.Printf("Proxy server started. HTTPS port: %s", httpsPort)

	return nil
}

func (p *TransparentProxy) Stop() error {
	if p.listenerHTTP != nil {
		err := p.listenerHTTP.Close()
		if err != nil {
			return err
		}
	}

	if p.listenerHTTPS != nil {
		err := p.listenerHTTPS.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *TransparentProxy) acceptConnections(listener net.Listener, isHTTPS bool) {
	for {
		clientConn, err := listener.Accept()
		fmt.Printf(": %v\n", clientConn)
		if err != nil {
			log.Printf("Failed to accept client connection: %v", err)
			continue
		}

		go p.handleClient(clientConn, isHTTPS)
	}
}

func (p *TransparentProxy) handleClient(clientConn net.Conn, isHTTPS bool) {
	defer clientConn.Close()

	index := GetRequestIndex()

	// Read the client's request
	requestBuf := make([]byte, 4096)
	n, err := clientConn.Read(requestBuf)
	if err != nil {
		log.Printf("Failed to read client request: %v", err)
		return
	}

	// Parse the request
	request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(requestBuf[:n])))
	if err != nil {
		log.Printf("Failed to parse client request: %v", err)
		return
	}

	p.PrintRequestData(requestBuf, index)

	if p.localCheck(request.Host) { // 目的ip不能是本地ip，回到导致无限循环
		log.Printf("不能直接访问代理ip %v", err)
		response := "Access to this server is prohibited."
		_, err := clientConn.Write([]byte(response))
		if err != nil {
			log.Printf("Failed to send response: %v", err)
		}
		return
	}
	if isHTTPS || (requestBuf[0] == 0x16 && requestBuf[1] == 0x03 && requestBuf[5] == 0x01) {
		// HTTPS connection, terminate SSL/TLS
		p.handleHTTPS(clientConn, requestBuf, request, index)
	} else {
		// HTTP connection, pass through
		p.handleHTTP(clientConn, requestBuf, request, index)
	}
}

func (p *TransparentProxy) handleHTTPS(clientConn net.Conn, clientHello []byte, request *http.Request, index uint64) {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	destination := request.Host
	if !strings.Contains(request.Host, ":") {
		destination = fmt.Sprintf("%s:%s", request.Host, "443")
	}

	fmt.Println(destination)
	remoteConn, err := tls.Dial("tcp", destination, tlsConfig)
	if err != nil {
		log.Printf("handleHTTPs: %d Failed to connect to the destination server: %v", index, err)
		return
	}
	defer remoteConn.Close()

	// Send client hello to the destination server
	_, err = remoteConn.Write(clientHello)
	if err != nil {
		log.Printf("handleHTTPs: Failed to send client hello to the destination server: %v", err)
		return
	}

	// Forward data between client and server
	go func() {
		_, err := io.Copy(remoteConn, clientConn)
		if err != nil {
			log.Printf("handleHTTPs: Error while copying data from client to server: %v", err)
		}
	}()

	_, err = io.Copy(clientConn, remoteConn)
	if err != nil {
		log.Printf("handleHTTPs: Error while copying data from server to client: %v", err)
	}
}

func (p *TransparentProxy) localCheck(domain string) bool {
	host, _, _ := net.SplitHostPort(domain)
	for _, proxyIP := range p.proxyIPs {
		if proxyIP == host {
			return true
		}
	}
	return false
}

func (p *TransparentProxy) handleHTTP(clientConn net.Conn, clientHello []byte, request *http.Request, index uint64) {

	// 获取实际目标地址，从请求中获取HOST，如果没有附带端口添加默认端口
	destination := request.Host
	if !strings.Contains(request.Host, ":") {
		destination = fmt.Sprintf("%s:%s", request.Host, "80")
	}

	// 开始连接
	remoteConn, err := net.Dial("tcp", destination)
	if err != nil {
		log.Printf("handleHTTP: Failed to connect to the destination server: %v", err)
		return
	}
	defer remoteConn.Close()

	// 向目标服务器发送客户端的原始请求
	_, err = remoteConn.Write(clientHello)
	if err != nil {
		log.Printf("handleHTTP: Failed to send client hello to the destination server: %v", err)
		return
	}

	// 将目标服务器的返回数据直接copy 到客户端请求中
	go func() {
		_, err := io.Copy(remoteConn, clientConn)
		if err != nil {
			log.Printf("handleHTTP: %d Error while copying data from client to server: %v", index, err)
		}
	}()

	// 创建一个 http.Request 对象
	req, err := http.ReadRequest(bufio.NewReader(clientConn))
	if err != nil && !errors.Is(err, io.EOF) {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			log.Printf("handleHTTP: 读取客户端请求失败：%v", err)
		}
		return
	}

	// 读取服务器响应
	response, err := http.ReadResponse(bufio.NewReader(remoteConn), req)
	if err != nil && !errors.Is(err, io.EOF) {
		log.Printf("handleHTTP: 读取服务器响应失败：%v", err)
		return
	}
	defer response.Body.Close()

	// 调用拦截器钩子函数来修改请求
	if p.responseHook != nil {
		response = p.responseHook(response)
	} else {
		log.Print("requestHook is nil")
	}

	responseBody := p.GetResponseData(response)

	p.PrintResponse(response, index)

	_, err = clientConn.Write(responseBody)
	if err != nil {
		log.Printf("handleHTTP: Failed to send response to the client: %v", err)
		return
	}

	// // 将响应从服务器复制到客户端
	// _, err = io.Copy(clientConn, remoteConn)
	// if err != nil {
	// 	log.Printf("handleHTTP: Error while copying data from server to client: %v", err)
	// }
	// 将修改后的响应写回客户端
	// err = response.Write(clientConn)
	// if err != nil {
	// 	log.Printf("handleHTTP: Failed to write modified response to the client: %v", err)
	// 	return
	// }
}
