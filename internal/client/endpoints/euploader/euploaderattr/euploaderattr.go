package euploaderattr

import "net/http"

type UploaderAttr struct {
	URL    string
	Client *http.Client
	Token  string
	Data   *[]byte
}

func (p *UploaderAttr) Init(
	eURL string,
	client *http.Client,
	token string,
) {
	p.URL = eURL
	p.Client = client
	p.Token = token
}
