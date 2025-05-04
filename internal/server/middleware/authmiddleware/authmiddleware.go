package authmiddleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dmitrovia/passkeeper/internal/general/logger"
	"github.com/dmitrovia/passkeeper/internal/server/middleware/authmiddleware/authmiddlewareattr"
	"github.com/dmitrovia/passkeeper/internal/server/models/ctxm"
	"github.com/dmitrovia/passkeeper/internal/server/models/userm"
	"github.com/golang-jwt/jwt/v4"
)

// const tokenLen = 2

var errUnexpectedMethod = errors.New("errUnexpectedMethod")

var errUserNotExist = errors.New("user is not exist")

func AuthMiddleware(
	attr *authmiddlewareattr.AuthMiddlewareAttr,
) func(http.Handler) http.Handler {
	handler := func(hand http.Handler) http.Handler {
		authFn := func(writer http.ResponseWriter,
			req *http.Request,
		) {
			authHeader := req.Header.Get("Authorization")

			if authHeader == "" {
				setErrStr(writer, attr, "header Authorization is empty")
				return
			}

			authToken := strings.Split(authHeader, " ")
			/*isBearer := authToken[0] == "Bearer"
			isLenValid := len(authToken) == tokenLen

			if !isLenValid || !isBearer {
				setErrStr(writer, attr, "Invalid token format")

				return
			}*/

			token, err := parseToken(authToken[0], attr)
			if err != nil {
				setErr(writer, attr, err)
				return
			}

			ctx, cancel := context.WithTimeout(
				req.Context(), attr.Dbtimeout)

			defer cancel()

			user, isValid, err := isValidToken(ctx, token, attr)
			if err != nil {
				setErr(writer, attr, err)
				return
			}

			if !isValid {
				setErrStr(writer, attr, "token is invalid")
				return
			}

			req = req.WithContext(
				context.WithValue(req.Context(), ctxm.UserKey, user))

			hand.ServeHTTP(writer, req)
		}

		return http.HandlerFunc(authFn)
	}

	return handler
}

func isValidToken(ctx context.Context,
	token *jwt.Token,
	attr *authmiddlewareattr.AuthMiddlewareAttr,
) (*userm.User, bool, error) {
	if !token.Valid {
		return nil, false, nil
	}

	claims, oka := token.Claims.(jwt.MapClaims)
	if !oka {
		return nil, false, nil
	}

	timeNow := float64(time.Now().Unix())
	claimsExp, oka := claims["exp"].(float64)

	if !oka {
		return nil, false, nil
	}

	if timeNow > claimsExp {
		return nil, false, nil
	}

	login, ok := claims["id"].(string)
	if !ok {
		return nil, false, nil
	}

	exist, user, err := attr.AuthService.UserIsExist(ctx,
		login)
	if err != nil {
		return nil, false, fmt.Errorf("IVT->UIE: %w", err)
	}

	if !exist {
		return nil, false, errUserNotExist
	}

	return user, true, nil
}

func parseToken(inToken string,
	attr *authmiddlewareattr.AuthMiddlewareAttr,
) (*jwt.Token, error) {
	token, err := jwt.Parse(inToken,
		func(token *jwt.Token) (interface{}, error) {
			_, isHMAC := token.Method.(*jwt.SigningMethodHMAC)

			if !isHMAC {
				headerAlg, oka := token.Header["alg"].(string)

				if !oka {
					return nil, errUnexpectedMethod
				}

				msg := "Unexpected signing method " + headerAlg
				logger.Log("AuthMiddleware", msg, attr.ZapLogger)

				return nil, errUnexpectedMethod
			}

			return []byte(attr.Secret), nil
		})
	if err != nil {
		return nil, fmt.Errorf("parseToken>jwt.Parse: %w", err)
	}

	return token, nil
}

func setErrStr(writer http.ResponseWriter,
	attr *authmiddlewareattr.AuthMiddlewareAttr,
	txt string,
) {
	writer.WriteHeader(http.StatusUnauthorized)
	logger.Log("AuthMiddleware",
		txt, attr.ZapLogger)
}

func setErr(writer http.ResponseWriter,
	attr *authmiddlewareattr.AuthMiddlewareAttr,
	err error,
) {
	writer.WriteHeader(http.StatusUnauthorized)
	logger.LogE("AuthMiddleware", err, attr.ZapLogger)
}
