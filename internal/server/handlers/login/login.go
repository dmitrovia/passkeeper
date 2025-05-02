package login

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
	"github.com/dmitrovia/passkeeper/internal/general/validate"
	"github.com/dmitrovia/passkeeper/internal/server/handlers/login/loginattr"
	"github.com/dmitrovia/passkeeper/internal/server/service"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var errEmptyData = errors.New("data is empty")

type Login struct {
	serv service.AuthService
	attr *loginattr.LoginAttr
}

func NewHandler(
	s service.AuthService,
	inAttr *loginattr.LoginAttr,
) *Login {
	return &Login{serv: s, attr: inAttr}
}

func (h *Login) LoginHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	reqAttr := &apim.InLoginUser{}

	err := getReqData(req, reqAttr, h.attr)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.LogE("login->getReqData", err, h.attr.ZapLogger)

		return
	}

	isValid := isValid(reqAttr)
	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(
		req.Context(), h.attr.Dbtimeout)
	defer cancel()

	exist, user, err := h.serv.UserIsExist(ctx, reqAttr.Login)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.LogE("login->UserIsExist", err, h.attr.ZapLogger)

		return
	}

	if !exist {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = checkPass(*user.Password, reqAttr.Password)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		logger.LogE("login->checkPass", err, h.attr.ZapLogger)

		return
	}

	token, err := generateToken(reqAttr.Login, h.attr)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.LogE("login->GT", err, h.attr.ZapLogger)

		return
	}

	writer.Header().Set("Authorization", token)
	writer.WriteHeader(http.StatusOK)
}

func generateToken(
	id string,
	attr *loginattr.LoginAttr,
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

func checkPass(hash string, pass string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash), []byte(pass))
	if err != nil {
		return fmt.Errorf("checkPass->bcrypt.CHAP: %w", err)
	}

	return nil
}

func isValid(reqAttr *apim.InLoginUser) bool {
	if reqAttr.Login == "" || reqAttr.Password == "" {
		return false
	}

	res := validate.IsValidLogin(reqAttr.Login)
	if !res {
		return false
	}

	res = validate.IsValidPass(reqAttr.Password)

	return res
}

func getReqData(
	req *http.Request,
	reqAttr *apim.InLoginUser,
	attrHandler *loginattr.LoginAttr,
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
