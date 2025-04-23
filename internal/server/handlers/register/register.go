package register

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/general/models/apim"
	"github.com/dmitrovia/passkeeper/internal/general/rsa"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/register/registerattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/dmitrovia/passkeeper/internal/server/service"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	authService service.AuthService
	attr        *registerattr.RegisterAttr
}

var errEmptyData = errors.New("data is empty")

const (
	statusISE = http.StatusInternalServerError
)

func NewHandler(
	authS service.AuthService,
	inAttr *registerattr.RegisterAttr,
) *Register {
	return &Register{
		authService: authS,
		attr:        inAttr,
	}
}

func (h *Register) RegisterHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	reqAttr := &apim.InRegisterUser{}

	err := getReqData(req, reqAttr, h.attr)
	if err != nil {
		setErr(writer, h.attr, err, "getReqData")

		return
	}

	isValid := validate(reqAttr)
	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(
		req.Context(), h.attr.Dbtimeout)
	defer cancel()

	exist, _, err := h.authService.UserIsExist(ctx,
		reqAttr.Login)
	if err != nil {
		setErr(writer, h.attr, err, "UserIsExist")

		return
	}

	if exist {
		writer.WriteHeader(http.StatusConflict)

		return
	}

	err = createUser(ctx, h, reqAttr)
	if err != nil {
		setErr(writer, h.attr, err, "CreateUser")

		return
	}

	token, err := generateToken(reqAttr.Login, h.attr)
	if err != nil {
		setErr(writer, h.attr, err, "generateToken")

		return
	}

	writer.Header().Set("Authorization", token)
	writer.WriteHeader(http.StatusOK)
}

func createUser(ctx context.Context,
	handler *Register,
	reqAttr *apim.InRegisterUser,
) error {
	passwHash, err := cryptPass(reqAttr.Password)
	if err != nil {
		return fmt.Errorf("CreateUser->cryptPass: %w", err)
	}

	user := &userm.User{}
	user.Login = &reqAttr.Login
	user.Password = &passwHash

	err = handler.authService.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("CreateUser->authService.CU: %w", err)
	}

	return nil
}

func setErr(writer http.ResponseWriter,
	inAttr *registerattr.RegisterAttr,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.LogE("register->"+method, err, inAttr.ZapLogger)
}

func validate(reqAttr *apim.InRegisterUser) bool {
	if reqAttr.Login == "" || reqAttr.Password == "" {
		return false
	}

	return true
}

func cryptPass(pass string) (string, error) {
	passwHash, err := bcrypt.GenerateFromPassword(
		[]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("cryptPass->GFP: %w", err)
	}

	return string(passwHash), nil
}

func getReqData(
	req *http.Request,
	reqAttr *apim.InRegisterUser,
	attrHandler *registerattr.RegisterAttr,
) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("getReqData->io.ReadAll: %w", err)
	}

	if len(body) == 0 {
		return fmt.Errorf("getReqData: %w", errEmptyData)
	}

	dec, err := rsa.Decrypt(&body, attrHandler.DecKey)
	if err != nil {
		return fmt.Errorf("getReqData->Decrypt: %w", err)
	}

	err = json.Unmarshal(*dec, &reqAttr)
	if err != nil {
		return fmt.Errorf("getReqData->json.Unmarshal: %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return fmt.Errorf("getReqData->req.Body.Close: %w", err)
	}

	return nil
}

func generateToken(
	id string,
	attr *registerattr.RegisterAttr,
) (string, error) {
	generateToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"id": id,
			"exp": time.Now().Add(
				time.Hour * time.Duration(
					attr.TokenExpHour)).Unix(),
		})

	token, err := generateToken.SignedString(
		[]byte(attr.Secret))
	if err != nil {
		return token, fmt.Errorf("generateToken->SS: %w", err)
	}

	return token, nil
}
