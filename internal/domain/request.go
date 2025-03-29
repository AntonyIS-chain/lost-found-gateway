package domain


type ProxyRequest struct {
	Method string
	Path string
	Headers map[string]string
	Body []byte
}


type ProxyResponse struct {
	StatusCode int
	Headers map[string]string
	Body []byte
}