package proxy

import (
	"fmt"
	"io"
	"net/http"
)

func (p *TransparentProxy) PrintRequest(request *http.Request) {
	fmt.Println(">-------------------------")
	fmt.Printf("%s %s %s\n", request.Method, request.URL, request.Proto)

	for key, values := range request.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		return
	}
	defer request.Body.Close()

	fmt.Println("")
	fmt.Println(string(body))
	fmt.Println("-------------------------<")
	fmt.Println("")
}

func (p *TransparentProxy) PrintRequestData(body []byte, index uint64) {
	fmt.Printf("--%d-----------------------\n", index)
	fmt.Printf("%s\n", string(body))
	fmt.Printf("--%d-----------------------\n", index)
}

func (p *TransparentProxy) GetResponseData(response *http.Response) []byte {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		return []byte(fmt.Sprintf("Error reading request body: %v", err))
	}

	return body
}

func (p *TransparentProxy) PrintResponse(response *http.Response, index uint64) {
	body, _ := io.ReadAll(response.Body)
	fmt.Printf("--%d-----------------------\n", index)
	fmt.Printf("%s\n", string(body))
	fmt.Printf("--%d-----------------------\n", index)
}
