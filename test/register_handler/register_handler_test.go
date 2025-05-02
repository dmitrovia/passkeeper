package loginhandler_test

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

	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/register"
	"github.com/dmitrovia/passkeeper/internal/server/migrator"
	"github.com/dmitrovia/passkeeper/internal/server/proc/serverproc/serverpa"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	tn     string
	login  string
	pass   string
	expcod int
	exbody string
	meth   string
	data   *[]byte
}

const stok int = http.StatusOK

const url = "https://localhost:8443"

const path = "../../internal/client/crypto/keys/public.pem"

const (
	statusISE = http.StatusInternalServerError
	statusBR  = http.StatusBadRequest
)

//nolint:lll,funlen
func getTestData(encKey *[]byte) *[]testData {
	tmp := make([]byte, 0)

	incd := GetIncorrectData()

	incd1 := GetIncorrectDataWithCrypto(encKey)

	return &[]testData{
		{
			tn:     "1",
			login:  "test" + randomString(),
			pass:   "temppass",
			expcod: stok,
			exbody: "",
			data:   nil,
		},
		{
			tn:     "2",
			login:  "test" + randomString(),
			pass:   "temppass",
			expcod: statusISE,
			exbody: "",
			data:   &tmp,
		},
		{
			tn:     "3",
			login:  "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest" + randomString(),
			pass:   "temppass",
			expcod: statusBR,
			exbody: "",
			data:   nil,
		},
		{
			tn:     "4",
			login:  "test" + randomString(),
			pass:   "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest",
			expcod: statusBR,
			exbody: "",
			data:   nil,
		},
		{
			tn:     "5",
			login:  "test" + randomString(),
			pass:   "temppass",
			expcod: statusISE,
			exbody: "",
			data:   incd,
		},
		{
			tn:     "6",
			login:  "test" + randomString(),
			pass:   "temppass",
			expcod: statusBR,
			exbody: "",
			data:   incd1,
		},
		{
			tn:     "7",
			login:  "",
			pass:   "",
			expcod: statusBR,
			exbody: "",
			data:   nil,
		},
		{
			tn:     "8",
			login:  "upload_test",
			pass:   "test",
			expcod: stok,
			exbody: "",
			data:   nil,
		},
	}
}

//nolint:funlen
func TestRegisterHandler(t *testing.T) {
	t.Helper()
	t.Parallel()

	time.Sleep(20 * time.Second)

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

	testCases := getTestData(&encKey)

	register := register.NewHandler(
		attr.AuthService, attr.RigsterAttr).RegisterHandler

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
				test.meth,
				url+"/api/user/register", bytes.NewReader(bodyReq))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			newr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/api/user/register",
				register)
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
	data := apim.InLoginUser{}
	data.Login = testd.login
	data.Password = testd.pass

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Input->Marshal: %w", err)
	}

	encrypt, err := rsa.Encrypt(&marshal, encKey)
	if err != nil {
		return nil, fmt.Errorf("Input->Encrypt: %w", err)
	}

	return encrypt, nil
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
