package eloginattr

import "net/http"

type LoginAttr struct {
	URL    string
	Client *http.Client
	Data   *[]byte
}

func (p *LoginAttr) Init(
	eURL string,
	client *http.Client,
) {
	p.URL = eURL
	p.Client = client
}
