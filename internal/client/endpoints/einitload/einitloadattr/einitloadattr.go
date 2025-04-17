package einitloadattr

import "net/http"

type InitLoadAttr struct {
	URL    string
	Client *http.Client
	Token  string
}

func (p *InitLoadAttr) Init(
	eURL string,
	client *http.Client,
	token string,
) {
	p.URL = eURL
	p.Client = client
	p.Token = token
}
