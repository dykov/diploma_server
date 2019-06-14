package domain

import (
	"../util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetAnswers(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	if resultJson, found := util.Cache("paragraph/" + idString + "/answer"); found {
		fmt.Fprintln(w, resultJson.(string))
		return
	}

	id, err := util.CheckId(w, idString)
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var answer []TestsAnswer
	db.Where("test_id = ?", id).Find(&answer)

	var answers = struct {
		TestsAnswers []TestsAnswer `json:"tests_answers"`
	}{answer}

	outputJson, err := json.Marshal(answers)
	if err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	fmt.Fprint(w, string(outputJson))

	util.SetCache("paragraph/"+idString+"/answer", string(outputJson))

}

func CreateAnswer(w http.ResponseWriter, r *http.Request) {

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var answer TestsAnswer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	if err := db.Create(&answer); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

}

func UpdateAnswer(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	id, err := util.CheckId(w, idString)
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var answer TestsAnswer
	if db.First(&answer, id).RowsAffected == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.NotFound))
		return
	}

	var upd TestsAnswer
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}
	upd.Id = id

	if err := db.Save(&upd); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

	util.DeleteCache("paragraph/" + idString + "/answer")

}

func DeleteAnswer(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	id, err := util.CheckId(w, idString)
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	if rows := db.Delete(&TestsAnswer{}, id).RowsAffected; rows == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.NotFound))
		return
	}

	util.DeleteCache("paragraph/" + idString + "/answer")

}
