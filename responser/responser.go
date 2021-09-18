package responser

import (
	"fmt"
	"net/http"
)

type ResponseStruct struct {
	Code    int
	Message string
}

func (output *ResponseStruct) SendError(w http.ResponseWriter) {
	w.WriteHeader(400)
	fmt.Fprintf(w, output.Message)
}

func (output *ResponseStruct) AddMessage(str string) {
	output.Message += str + "\n"
}

func (output *ResponseStruct) Finish(w http.ResponseWriter) {
	w.WriteHeader(200)
	fmt.Fprintf(w, output.Message)
}
