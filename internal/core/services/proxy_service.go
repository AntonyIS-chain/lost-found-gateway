package services

import (
	"io/ioutil"
	"net/http"

	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/ports"
)

type ProxyServiceImpl struct{}

func NewProxyService() ports.ProxyService {
	return &ProxyServiceImpl{}
}

func (p *ProxyServiceImpl) ForwardRequest(req domain.ProxyRequest) (domain.ProxyResponse, error) {
	client := &http.Client{}

	request, err := http.NewRequest(req.Method, req.Path, nil)

	if err != nil {
		return domain.ProxyResponse{}, err
	}

	for key, value := range req.Headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)

	if err != nil {
		return domain.ProxyResponse{}, err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	return domain.ProxyResponse{
		StatusCode: response.StatusCode,
		Headers:    map[string]string{"Content-Type": response.Header.Get("Content-Type")},
		Body:       body,
	}, nil
}
