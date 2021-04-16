package mock

import netHTTP "net/http"

type HTTPClient struct {
	DoResponse netHTTP.Response
}

func (h *HTTPClient) Do(req *netHTTP.Request) (*netHTTP.Response, error) {
	return &h.DoResponse, nil
}
