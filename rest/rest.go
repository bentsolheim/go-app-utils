package rest

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

type Response struct {
	Message string
	Items   interface{}
}

func JsonResponse(w http.ResponseWriter, f func() (interface{}, error)) {
	encoder := json.NewEncoder(w)
	data, err := f()
	var response Response
	if err != nil {
		response = Response{Message: err.Error()}
	} else {
		response = Response{Message: "OK", Items: data}
	}
	if err := encoder.Encode(conventionalMarshaller{response}); err != nil {
		println(err.Error())
	}
}

var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)

type conventionalMarshaller struct {
	Value interface{}
}

func (m conventionalMarshaller) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(m.Value)

	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			args := string(match)
			return []byte((args[:1] + strings.ToLower(string(args[1])) + args[2:]))
		},
	)

	return converted, err
}
