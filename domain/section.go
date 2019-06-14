package domain

import (
	"../util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetSections(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	if resultJson, found := util.Cache("course/" + idString + "/section"); found {
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

	var section []Section
	db.Where("course_id = ?", id).Find(&section)

	var sections = struct {
		Sections []Section `json:"sections"`
	}{section}

	outputJson, err := json.Marshal(sections)
	if err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	fmt.Fprint(w, string(outputJson))

	util.SetCache("course/"+idString+"/section", string(outputJson))

}

func CreateSection(w http.ResponseWriter, r *http.Request) {

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var section Section
	if err := json.NewDecoder(r.Body).Decode(&section); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	if err := db.Create(&section); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

}

func UpdateSection(w http.ResponseWriter, r *http.Request) {

	id, err := util.CheckId(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var section Section
	if db.First(&section, id).RowsAffected == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.NotFound))
		return
	}

	var updSection Section
	if err := json.NewDecoder(r.Body).Decode(&updSection); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}
	updSection.Id = id

	if err := db.Save(&updSection); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

	util.DeleteCache("course/" + mux.Vars(r)["id"] + "/section")

}

func DeleteSection(w http.ResponseWriter, r *http.Request) {

	id, err := util.CheckId(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var lessons []Lesson
	db.Where("section_id = ?", id).Find(&lessons)

	var pot []ParagraphsOrTest
	for _, l := range lessons {
		var p []ParagraphsOrTest
		db.Where("lesson_id = ?", l.Id).Find(&p)
		pot = append(pot, p...)
	}

	var ta []TestsAnswer
	for _, p := range pot {
		var t []TestsAnswer
		db.Where("test_id = ?", p.Id).Find(&t)
		ta = append(ta, t...)
	}

	var ut []UsersTest
	for _, p := range pot {
		var u []UsersTest
		db.Where("test_id = ?", p.Id).Find(&u)
		ut = append(ut, u...)
	}

	for i := range ut {
		db.Delete(&ut[i])
	}
	for i := range ta {
		db.Delete(&ta[i])
	}
	for i := range pot {
		db.Delete(&pot[i])
	}
	for i := range lessons {
		db.Delete(&lessons[i])
	}

	var section Section
	if rows := db.Delete(&section, id).RowsAffected; rows == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.NotFound))
		return
	}

	util.DeleteCache("course/" + mux.Vars(r)["id"] + "/section")

}
