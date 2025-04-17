package eloadattr

import "net/http"

type LoadAttr struct {
	URL    string
	Client *http.Client
	Token  string
	Data   *[]byte
}

func (p *LoadAttr) Init(
	eURL string,
	client *http.Client,
	token string,
	data *[]byte,
) {
	p.URL = eURL
	p.Client = client
	p.Token = token
	p.Data = data
}
