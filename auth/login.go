package auth

import (
	"../domain"
	"../util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	var checkUser domain.User
	if db.Where("login = ?", user.Login).First(&checkUser).RowsAffected == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.UserNotFound))
		return
	}

	if util.HashPassword(user.Password, util.Salt) != checkUser.Password {
		util.SendErr(w, http.StatusBadRequest, errors.New("Wrong password "))
		return
	}

	if checkUser.VerificationCode != "" {
		util.SendErr(w, http.StatusBadRequest, errors.New("Please, pass verification using the link in your email "))
		return
	}

	token := createToken(checkUser.Role, checkUser.Id)
	outputJson, _ := json.Marshal(token)
	fmt.Fprint(w, string(outputJson))

}
