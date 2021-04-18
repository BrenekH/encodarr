package mock

import netHTTP "net/http"

type HTTPClient struct {
	DoResponse  netHTTP.Response
	LastRequest netHTTP.Request
}

func (h *HTTPClient) Do(req *netHTTP.Request) (*netHTTP.Response, error) {
	h.LastRequest = *req
	return &h.DoResponse, nil
}
