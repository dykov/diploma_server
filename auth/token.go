package auth

import (
	"../util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {

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

	if (token.Claims.(jwt.MapClaims)["type"]).(string) != "refresh_token" {
		log.Errorln("Incorrect token type")
		log.WithFields(log.Fields{
			"Expected": "refresh_token",
			"Got":      (token.Claims.(jwt.MapClaims)["type"]).(string),
		}).Debugln("Incorrect token type")
		util.SendErr(w, http.StatusUnauthorized, errors.New("Unauthorized: Incorrect token type "))
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	userId, _ := strconv.ParseUint((token.Claims.(jwt.MapClaims)["user_id"]).(string), 10, 64)
	userGrant := uint64((token.Claims.(jwt.MapClaims)["grant"]).(float64))

	var jsonToken = createToken(userGrant, userId)
	outputJson, _ := json.Marshal(jsonToken)
	fmt.Fprint(w, string(outputJson))
}

func createToken(grant uint64, userId uint64) (jsonToken Token) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = userId
	claims["type"] = "access_token"
	claims["exp"] = time.Now().Add(util.AccessTokenExpiration).Unix()
	claims["grant"] = grant
	accessToken, _ := token.SignedString(util.SecretKeyForToken)

	claims["type"] = "refresh_token"
	claims["exp"] = time.Now().Add(util.RefreshTokenExpiration).Unix()
	refreshToken, _ := token.SignedString(util.SecretKeyForToken)

	jsonToken = Token{
		accessToken,
		refreshToken,
	}

	return jsonToken

}
