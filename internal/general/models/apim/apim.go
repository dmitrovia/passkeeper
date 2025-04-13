package apim

type CfgServer struct {
	ServerAddr       string `json:"address"`
	DBDSN            string `json:"dbdsn"`
	FilesStoragePath string `json:"filesStoragePath"`
	CryptoKeyPath    string `json:"cryptoKeyPath"`
}

type CfgClient struct {
	ServerAddr           string `json:"address"`
	FilesSynchronizePath string `json:"filesSynchronizePath"`
	MetaPath             string `json:"metaPath"`
	TokenPath            string `json:"tokenPath"`
	CryptoKeyPath        string `json:"cryptoKeyPath"`
	GzipFormats          string `json:"gzipFormats"`
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
