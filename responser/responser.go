package responser

import (
	"fmt"
	"net/http"
)

type ResponseStruct struct {
	Code    int
	Message string
	W       http.ResponseWriter
}

func (output *ResponseStruct) SendError() {
	output.W.WriteHeader(400)
	fmt.Fprintf(output.W, output.Message)
	panic("requesttheend")
}

func (output *ResponseStruct) AddMessage(str string) {
	output.Message += str + "\n"
}

func (output *ResponseStruct) Finish() {
	output.W.WriteHeader(200)
	fmt.Fprintf(output.W, output.Message)
}
