package domain

import (
	"../util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetLessons(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	if resultJson, found := util.Cache("section/" + idString + "/lesson"); found {
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

	var lesson []Lesson
	db.Where("section_id = ?", id).Find(&lesson)

	var lessons = struct {
		Lessons []Lesson `json:"lessons"`
	}{lesson}

	outputJson, err := json.Marshal(lessons)
	if err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	fmt.Fprint(w, string(outputJson))

	util.SetCache("section/"+idString+"/lesson", string(outputJson))

}

func CreateLesson(w http.ResponseWriter, r *http.Request) {

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var lesson Lesson
	if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	if err := db.Create(&lesson); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

}

func UpdateLesson(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	id, err := util.CheckId(w, idString)
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var lesson Lesson
	if db.First(&lesson, id).RowsAffected == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.NotFound))
		return
	}

	var upd Lesson
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}
	upd.Id = id

	if err := db.Save(&upd); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

	util.DeleteCache("section/" + idString + "/lesson")

}

func DeleteLesson(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	id, err := util.CheckId(w, idString)
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var pot []ParagraphsOrTest
	db.Where("lesson_id = ?", id).Find(&pot)

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

	if rows := db.Delete(&Lesson{}, id).RowsAffected; rows == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.NotFound))
		return
	}

	util.DeleteCache("section/" + idString + "/lesson")

}
