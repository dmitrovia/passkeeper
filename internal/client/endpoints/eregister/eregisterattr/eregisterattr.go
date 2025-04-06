package eregisterattr

import "net/http"

type RegisterAttr struct {
	URL    string
	Client *http.Client
	Data   *[]byte
}

func (p *RegisterAttr) Init(
	eURL string,
	client *http.Client,
) {
	p.URL = eURL
	p.Client = client
}
