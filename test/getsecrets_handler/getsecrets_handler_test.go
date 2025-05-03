package getsecretshandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/getsecrets"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/migrator"
	"github.com/dmitrovia/passkeeper/internal/server/proc/serverproc/serverpa"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	tn     string
	expcod int
	exbody string
	token  *string
}

const url = "https://localhost:8443"

const (
	statusISE = http.StatusInternalServerError
	statusBR  = http.StatusBadRequest
	statusU   = http.StatusUnauthorized
)

//nolint:lll
func getTestData() *[]testData {
	tok := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDYzNTI1MTEsImlkIjoicXFxIn0.6O4koFaQNd4vdWWb2WsftHL0Ewb45_4-tng1qcPn980"

	tok1 := ""

	return &[]testData{
		{
			tn:     "1",
			expcod: statusU,
			exbody: "",
			token:  &tok,
		},
		{
			tn:     "2",
			expcod: statusU,
			exbody: "",
			token:  &tok1,
		},
	}
}

//nolint:funlen
func TestGetSecretsHandler(t *testing.T) {
	t.Helper()
	t.Parallel()

	time.Sleep(40 * time.Second)

	attr := &serverpa.ServerProcAttr{}

	err := attr.Init()
	if err != nil {
		t.Errorf("Init: %v", err)

		return
	}

	err = migrator.UseMigrations(attr)
	if err != nil {
		t.Errorf("Init: %v", err)

		return
	}

	path := "../../internal/client/auth/token.json"

	tok, err := authcfg.GetToken(path)
	if err != nil {
		t.Errorf("GetToken: %v", err)
	}

	Token := tok

	testCases := getTestData()

	getSecretsH := getsecrets.NewHandler(attr.SecretService,
		attr.GetSecretsAttr).GetSecretsHandler

	for _, test := range *testCases {
		t.Run(http.MethodPost, func(t *testing.T) {
			t.Parallel()

			tmp := make([]byte, 0)

			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				url+"/api/user/getsecrets", bytes.NewReader(tmp))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			if test.token != nil {
				req.Header.Set("Authorization", *test.token)
			} else {
				req.Header.Set("Authorization", Token)
			}

			newr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Use(authmiddleware.AuthMiddleware(
				attr.AuthMidAttr))
			router.HandleFunc("/api/user/getsecrets",
				getSecretsH)
			router.ServeHTTP(newr, req)
			status := newr.Code
			body, _ := io.ReadAll(newr.Body)

			assert.Equal(t,
				test.expcod,
				status, test.tn+": Response code didn't match expected")

			if test.exbody != "" {
				assert.JSONEq(t, test.exbody, string(body))
			}
		})
	}
}

//nolint:errchkjson
func GetIncorrectData() *[]byte {
	incd := &apim.IncorrectData{}
	incd.IncorrectData = "IncorrectData"
	marshal, _ := json.Marshal(incd)

	return &marshal
}

//nolint:errchkjson
func GetIncorrectDataWithCrypto(encKey *[]byte) *[]byte {
	incd := &apim.IncorrectData{}
	incd.IncorrectData = "IncorrectData"
	marshal, _ := json.Marshal(incd)

	encrypt, _ := rsa.Encrypt(&marshal, encKey)

	return encrypt
}
