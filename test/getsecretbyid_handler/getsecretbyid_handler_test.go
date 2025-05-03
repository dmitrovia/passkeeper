package getsecretbyidhandler_test

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

	"github.com/dmitrovia/passkeeper/internal/client/auth/authcfg"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/getsecretbyid"
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
}

const url = "https://localhost:8443"

const path = "../../internal/client/crypto/keys/public.pem"

const (
	statusISE = http.StatusInternalServerError
	statusBR  = http.StatusBadRequest
)

//nolint:lll
func getTestData(encKey *[]byte) *[]testData {
	tmp := make([]byte, 0)

	incd := GetIncorrectData()

	incd1 := GetIncorrectDataWithCrypto(encKey)

	return &[]testData{
		{
			tn:           "1",
			inIdentifier: "test" + randomString(),
			expcod:       statusISE,
			exbody:       "",
			data:         &tmp,
		},
		{
			tn:           "2",
			inIdentifier: "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest" + randomString(),
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
		},
		{
			tn:           "3",
			inIdentifier: "test" + randomString(),
			expcod:       statusISE,
			exbody:       "",
			data:         incd,
		},
		{
			tn:           "4",
			inIdentifier: "test" + randomString(),
			expcod:       statusBR,
			exbody:       "",
			data:         incd1,
		},
		{
			tn:           "5",
			inIdentifier: "",
			expcod:       statusBR,
			exbody:       "",
			data:         nil,
		},
	}
}

//nolint:funlen
func TestLoginHandler(t *testing.T) {
	t.Helper()
	t.Parallel()

	// time.Sleep(20 * time.Second)

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
	fmt.Println(Token)

	testCases := getTestData(&encKey)

	getSecretByIDH := getsecretbyid.NewHandler(
		attr.SecretService,
		attr.GetSecretByIDAttr).GetSecretByIDHadnler

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
				http.MethodGet,
				url+"/api/user/getsecretbyid", bytes.NewReader(bodyReq))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Token)

			newr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Use(authmiddleware.AuthMiddleware(
				attr.AuthMidAttr))
			router.HandleFunc("/api/user/getsecretbyid",
				getSecretByIDH)
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
	outAttr := &apim.InGetSecretByID{}
	outAttr.Identifier = testd.inIdentifier

	marshal, err := json.Marshal(outAttr)
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

	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
