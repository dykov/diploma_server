package auth

import (
	"../util"
	"errors"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func Middleware(handler http.Handler, accessGrant uint64, middleware ...func(http.Handler, uint64) http.Handler) http.Handler {

	for _, mw := range middleware {
		handler = mw(handler, accessGrant)
	}

	return handler

}

func AuthMiddleware(handler http.Handler, accessGrant uint64) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			tokenHeader := r.Header.Get("Authorization") // get token in format `Bearer {token-body}`
			if tokenHeader == "" {
				log.Errorln("Missing auth token")
				log.WithField("Request Header", r.Header).Debugln("Missing auth token")
				util.SendErr(w, http.StatusUnauthorized, errors.New("Missing auth token "))
				return
			}

			splitted := strings.Split(tokenHeader, " ") // splitted token must have 2 parts: `Bearer` and `{token-body}`
			if len(splitted) != 2 {
				log.Errorln("Invalid auth token")
				log.WithFields(log.Fields{
					"Splitted token": splitted,
					"Request Header": r.Header,
				}).Debugln("Invalid auth token")
				util.SendErr(w, http.StatusUnauthorized, errors.New("Invalid auth token "))
				return
			}
			tokenPart := splitted[1] // get token body

			token, tokenErr := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
				return []byte(util.SecretKeyForToken), nil
			})
			// if token is invalid or expired
			if tokenErr != nil || token == nil {
				log.Debugln("Token error", tokenErr)
				util.SendErr(w, http.StatusUnauthorized, errors.New("Unauthorized: "+tokenErr.Error()))
				return
			}

			if (token.Claims.(jwt.MapClaims)["type"]).(string) != "access_token" {
				log.Errorln("Incorrect token type")
				log.WithFields(log.Fields{
					"Expected": "access_token",
					"Got":      (token.Claims.(jwt.MapClaims)["type"]).(string),
				}).Debugln("Incorrect token type")
				util.SendErr(w, http.StatusUnauthorized, errors.New("Unauthorized: Incorrect token type "))
				return
			}

			db := util.GetDB(w)
			if db.DB() == nil {
				return
			}

			userGrant := uint64((token.Claims.(jwt.MapClaims)["grant"]).(float64))

			if userGrant < accessGrant {
				log.Errorln("Not enough grant")
				log.WithFields(log.Fields{
					"Request":      r,
					"User grant":   userGrant,
					"Access grant": accessGrant,
				}).Debugln("Not enough grant")
				util.SendErr(w, http.StatusForbidden, errors.New("Not enough grant "))
				return
			}

			handler.ServeHTTP(w, r)

		},
	)

}
