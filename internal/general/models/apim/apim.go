package apim

type CfgServer struct {
	ServerAddr           string `json:"address"`
	DBDSN                string `json:"dbdsn"`
	FilesStoragePath     string `json:"filesStoragePath"`
	CryptoKeyPathPrivate string `json:"cryptoKeyPathPrivate"`
	CryptoKeyPathPublic  string `json:"cryptoKeyPathPublic"`
	SecretAuth           string `json:"secretAuth"`
}

type CfgClient struct {
	ServerAddr           string `json:"address"`
	FilesUploadPath      string `json:"filesUploadPath"`
	MetaPath             string `json:"metaPath"`
	TokenPath            string `json:"tokenPath"`
	CryptoKeyPathPublic  string `json:"cryptoKeyPathPublic"`
	CryptoKeyPathPrivate string `json:"cryptoKeyPathPrivate"`
	GzipFormats          string `json:"gzipFormats"`
	Aes256key            string `json:"aes256key"`
}

type InInitSingleLoad struct {
	FileName string `json:"fileName"`
}

type InGetSecretByID struct {
	Identifier string `json:"identifier"`
}

type InRegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type InLoginUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CfgToken struct {
	Token string `json:"token"`
}
