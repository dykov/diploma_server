package domain

import (
	"../util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetParagraphs(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	if resultJson, found := util.Cache("lesson/" + idString + "/paragraph"); found {
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

	var paragraph []ParagraphsOrTest
	db.Where("lesson_id = ?", id).Find(&paragraph)

	var paragraphs = struct {
		ParagraphsOrTests []ParagraphsOrTest `json:"paragraphs_or_tests"`
	}{paragraph}

	outputJson, err := json.Marshal(paragraphs)
	if err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	fmt.Fprint(w, string(outputJson))

	util.SetCache("lesson/"+idString+"/paragraph", string(outputJson))

}

func CreateParagraph(w http.ResponseWriter, r *http.Request) {

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var paragraph ParagraphsOrTest
	if err := json.NewDecoder(r.Body).Decode(&paragraph); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	if err := db.Create(&paragraph); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

}

func UpdateParagraph(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	id, err := util.CheckId(w, idString)
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var paragraph ParagraphsOrTest
	if db.First(&paragraph, id).RowsAffected == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.NotFound))
		return
	}

	var upd ParagraphsOrTest
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}
	upd.Id = id

	if err := db.Save(&upd); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

	util.DeleteCache("lesson/" + idString + "/paragraph")

}

func DeleteParagraph(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	id, err := util.CheckId(w, idString)
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var ta []TestsAnswer
	db.Where("test_id = ?", id).Find(&ta)

	var ut []UsersTest
	db.Where("test_id = ?", id).Find(&ut)

	for i := range ut {
		db.Delete(&ut[i])
	}
	for i := range ta {
		db.Delete(&ta[i])
	}

	if rows := db.Delete(&ParagraphsOrTest{}, id).RowsAffected; rows == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.NotFound))
		return
	}

	util.DeleteCache("lesson/" + idString + "/paragraph")

}
