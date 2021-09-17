package theend

import (
	"fmt"
	"net/http"
)

type ResponseTheEnd struct {
	Code    int
	Message string
}

func (output *ResponseTheEnd) TheEnd(w http.ResponseWriter) {
	if output.Code == 0 {
		output.Code = 200
	}

	w.WriteHeader(output.Code)
	fmt.Fprintf(w, output.Message)
}

func (output *ResponseTheEnd) SetMessage(str string) {
	output.Message += str
}
