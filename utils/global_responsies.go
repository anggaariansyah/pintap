package utils
import "net/http"

type Meta struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ModelResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}
type EmptyStruct struct{}
type EmptyResponse struct {
	Content interface{} `json:"content"`
}

func Response(code int, msg string, data interface{}) *ModelResponse {
	res := new(ModelResponse)
	meta := new(Meta)
	meta.Code = code
	switch code {
	case http.StatusOK:
		meta.Status = true
	case http.StatusInternalServerError:
		meta.Status = false
	case http.StatusBadRequest:
		meta.Status = false
	case http.StatusCreated:
		meta.Status = true
	case http.StatusNotFound:
		meta.Status = false
	case http.StatusUnauthorized:
		meta.Status = false
	}
	meta.Message = msg
	res.Meta = *meta
	res.Data = data

	return res
}


