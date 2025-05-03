package uploadsecret_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/models/secret"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/uploadsecret"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware"
	"github.com/dmitrovia/passkeeper/internal/server/migrator"
	"github.com/dmitrovia/passkeeper/internal/server/proc/serverproc/serverpa"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	tn           string
	inIdentifier string
	expcod       int
	exbody       string
	data         *[]byte
	token        *string
	secret       string
}

const url = "https://localhost:8443"

const path = "../../internal/client/crypto/keys/public.pem"

const (
	statusISE = http.StatusInternalServerError
	statusBR  = http.StatusBadRequest
	statusU   = http.StatusUnauthorized
)

const strLen = 5

//nolint:lll,funlen
func getTestData(encKey *[]byte) *[]testData {
	tmp := make([]byte, 0)

	incd := GetIncorrectData()

	incd1 := GetIncorrectDataWithCrypto(encKey)

	tok := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDYzNTI1MTEsImlkIjoicXFxIn0.6O4koFaQNd4vdWWb2WsftHL0Ewb45_4-tng1qcPn980"

	tok1 := ""

	return &[]testData{
		{
			tn:           "1",
			inIdentifier: "test" + randomString(),
			expcod:       statusISE,
			exbody:       "",
			data:         &tmp,
			token:        nil,
			secret:       "test",
		},
		{
			tn:           "2",
			inIdentifier: "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest" + randomString(),
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
			secret:       "test",
		},
		{
			tn:           "3",
			inIdentifier: "test" + randomString(),
			expcod:       statusISE,
			exbody:       "",
			data:         incd,
			token:        nil,
			secret:       "test",
		},
		{
			tn:           "4",
			inIdentifier: "test" + randomString(),
			expcod:       statusBR,
			exbody:       "",
			data:         incd1,
			token:        nil,
			secret:       "test",
		},
		{
			tn:           "5",
			inIdentifier: "",
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
			secret:       "test",
		},
		{
			tn:           "6",
			inIdentifier: "",
			expcod:       statusU,
			exbody:       "",
			data:         nil,
			token:        &tok,
			secret:       "test",
		},
		{
			tn:           "7",
			inIdentifier: "",
			expcod:       statusU,
			exbody:       "",
			data:         nil,
			token:        &tok1,
			secret:       "test",
		},
		{
			tn:           "8",
			inIdentifier: "test2233",
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
			token:        nil,
			secret:       "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest",
		},
	}
}

//nolint:funlen,cyclop
func TestUploadSecretHandler(t *testing.T) {
	t.Helper()
	t.Parallel()

	time.Sleep(30 * time.Second)

	attr := &serverpa.ServerProcAttr{}

	err := attr.Init()
	if err != nil {
		t.Errorf("Init: %v", err)

		return
	}

	encKey, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("ReadFile: %v", err)

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

	testCases := getTestData(&encKey)

	uploadSecretH := uploadsecret.NewHandler(
		attr.SecretService,
		attr.UploadSecretAttr).UploadSecretHandler

	for _, test := range *testCases {
		t.Run(http.MethodPost, func(t *testing.T) {
			t.Parallel()

			reqData, err := formReqBody(&test, &encKey)
			if err != nil {
				fmt.Println(err)

				return
			}

			var bodyReq []byte
			if test.data != nil {
				bodyReq = *test.data
			} else {
				bodyReq = *reqData
			}

			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodPost,
				url+"/api/user/uploadsecret", bytes.NewReader(bodyReq))
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
			router.HandleFunc("/api/user/uploadsecret",
				uploadSecretH)
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

func formReqBody(
	testd *testData,
	encKey *[]byte,
) (*[]byte, error) {
	secret := &secret.Secret{}
	secret.Identifier = &testd.inIdentifier
	secret.Value = &testd.secret

	marshal, err := json.Marshal(secret)
	if err != nil {
		return nil, fmt.Errorf("formReqBody->Marshal: %w", err)
	}

	encrypt, err := rsa.Encrypt(&marshal, encKey)
	if err != nil {
		return nil, fmt.Errorf("formReqBody->Encrypt: %w", err)
	}

	return encrypt, nil
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

func randomString() string {
	letters := []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, strLen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
