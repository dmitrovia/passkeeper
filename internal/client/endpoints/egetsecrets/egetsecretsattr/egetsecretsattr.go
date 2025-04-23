package egetsecretsattr

import "net/http"

type GetSecretsAttr struct {
	URL    string
	Client *http.Client
	Token  string
}

func (p *GetSecretsAttr) Init(
	eURL string,
	client *http.Client,
	token string,
) {
	p.URL = eURL
	p.Client = client
	p.Token = token
}
