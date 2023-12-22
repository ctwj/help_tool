package proxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
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

func getTLSConfig(certPEM, keyPEM string) (*tls.Config, error) {
	// 创建一个空的 tls.Config
	config := &tls.Config{}

	// 解析证书和私钥
	tlsCert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return nil, err
	}

	// 将证书和私钥设置到配置中
	config.Certificates = []tls.Certificate{tlsCert}

	return config, nil
}

type TransparentProxy struct {
	listenerHTTP  net.Listener
	listenerHTTPS net.Listener
	destination   string
	requestHook   RequestInterceptor // 添加 requestHook 字段
	responseHook  ResponseInterceptor
	proxyIPs      []string
	tlsConfig     *tls.Config
}

type RequestInterceptor func(req *http.Request) *http.Request
type ResponseInterceptor func(req *http.Response) *http.Response

func NewTransparentProxy(destination string) *TransparentProxy {
	ips, _ := getAllLocalIPs()
	tlsConfig, _ := getTLSConfig(cert, key)
	return &TransparentProxy{
		destination: destination,
		proxyIPs:    ips,
		tlsConfig:   tlsConfig,
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

	// 打印请求内容
	fmt.Println("Received request from client:")
	fmt.Println(request.Method, request.URL.String())
	request.Header.Write(os.Stdout)

	// 获取实际目标地址
	destination := request.Host
	if !strings.Contains(request.Host, ":") {
		destination = fmt.Sprintf("%s:%s", request.Host, "443") // 默认使用 HTTPS 的端口 443
	}
	targetURL := fmt.Sprintf("https://%s%s", destination, request.URL.String())
	fmt.Println("Target URL:", targetURL)

	// 发起连接到实际目标地址
	remoteConn, err := tls.Dial("tcp", destination, p.tlsConfig)
	if err != nil {
		log.Printf("Failed to connect to the destination server: %v", err)
		return
	}
	defer remoteConn.Close()

	// 发送客户端的原始请求
	err = request.Write(remoteConn)
	if err != nil {
		log.Printf("Failed to send client request to the destination server: %v", err)
		return
	}

	// 读取目标服务器的响应
	remoteResp, err := http.ReadResponse(bufio.NewReader(remoteConn), request)
	if err != nil {
		log.Printf("Failed to read response from the destination server: %v", err)
		return
	}
	defer remoteResp.Body.Close()

	// 打印响应内容
	fmt.Println("Received response from server:")
	fmt.Println(remoteResp.Status)
	remoteResp.Header.Write(os.Stdout)

	// 读取实际目标地址返回的内容
	responseBody, err := io.ReadAll(remoteResp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return
	}

	// 修改返回内容
	modifiedResponse := strings.ReplaceAll(string(responseBody), "code", "CXXX")
	fmt.Println("Modified response:", modifiedResponse)

	// 将修改后的内容写回到 clientConn 返回给浏览器
	_, err = clientConn.Write([]byte(modifiedResponse))
	if err != nil {
		log.Printf("Failed to write response to client: %v", err)
		return
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
	// 从 clientConn 获取请求内容
	clientRequest := request
	destination := request.Host
	if !strings.Contains(request.Host, ":") {
		destination = fmt.Sprintf("%s:%s", request.Host, "80")
	}

	// 打印请求内容
	fmt.Println("Received request from client:")
	fmt.Println(clientRequest.Method, clientRequest.URL.String())
	clientRequest.Header.Write(os.Stdout)

	// 获取实际目标地址
	targetURL := fmt.Sprintf("http://%s%s", destination, clientRequest.URL.String())
	fmt.Println("Target URL:", targetURL)

	// 复制原始请求的主体
	var requestBody bytes.Buffer
	_, err := io.Copy(&requestBody, clientRequest.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		return
	}
	clientRequest.Body.Close()

	// 创建新的请求对象，包含原始请求的主体
	newRequest, err := http.NewRequest(clientRequest.Method, targetURL, &requestBody)
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		return
	}
	newRequest.Header = clientRequest.Header

	// 发起请求到实际目标地址
	client := &http.Client{}
	remoteResp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Failed to send request to target URL: %v", err)
		return
	}
	defer remoteResp.Body.Close()

	// 读取实际目标地址返回的内容
	responseBody, err := io.ReadAll(remoteResp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return
	}

	// 修改返回内容
	modifiedResponse := strings.ReplaceAll(string(responseBody), "code", "CXXX")
	fmt.Println("Modified response:", modifiedResponse)

	// 将修改后的内容写回到 clientConn 返回给浏览器
	_, err = clientConn.Write([]byte(modifiedResponse))
	if err != nil {
		log.Printf("Failed to write response to client: %v", err)
		return
	}

}

func (p *TransparentProxy) handleHTTP1(clientConn net.Conn, clientHello []byte, request *http.Request, index uint64) {
	// 从 clientConn 获取请求内容
	clientRequest := request
	destination := request.Host
	if !strings.Contains(request.Host, ":") {
		destination = fmt.Sprintf("%s:%s", request.Host, "80")
	}

	// 打印请求内容
	fmt.Println("Received request from client:")
	fmt.Println(clientRequest.Method, clientRequest.URL.String())
	clientRequest.Header.Write(os.Stdout)

	// 获取实际目标地址
	targetURL := clientRequest.URL.String()
	fmt.Println("Target URL:", targetURL)

	// 发起请求到实际目标地址
	fmt.Println("Connecting to target URL...")
	remoteConn, err := net.Dial("tcp", destination)
	if err != nil {
		log.Printf("Failed to connect to target URL: %v", err)
		return
	}
	defer remoteConn.Close()

	// 将请求发送到实际目标地址
	fmt.Println("Sending request to target URL...")
	err = clientRequest.Write(remoteConn)
	if err != nil {
		log.Printf("Failed to send request to target URL: %v", err)
		return
	}

	// 从 remoteConn 中获取实际目标地址返回的内容
	fmt.Println("Reading response from target URL...")
	remoteResp, err := http.ReadResponse(bufio.NewReader(remoteConn), clientRequest)
	if err != nil {
		log.Printf("Failed to read response from target URL: %v", err)
		return
	}
	defer remoteResp.Body.Close()

	// 读取实际目标地址返回的内容
	responseBody, err := io.ReadAll(remoteResp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return
	}

	// 修改返回内容
	modifiedResponse := strings.ReplaceAll(string(responseBody), "code", "XXXX")
	fmt.Println("Modified response:", modifiedResponse)

	// 将修改后的内容写回到 clientConn 返回给浏览器
	fmt.Println("Writing response to client...")
	_, err = clientConn.Write([]byte(modifiedResponse))
	if err != nil {
		log.Printf("Failed to write response to client: %v", err)
		return
	}
}
