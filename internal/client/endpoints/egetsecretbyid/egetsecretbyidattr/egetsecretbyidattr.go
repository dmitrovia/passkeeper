package egetsecretbyidattr

import "net/http"

type GetSecretByIDAttr struct {
	URL    string
	Client *http.Client
	Token  string
	Data   *[]byte
}

func (p *GetSecretByIDAttr) Init(
	eURL string,
	client *http.Client,
	token string,
) {
	p.URL = eURL
	p.Client = client
	p.Token = token
}
