package domain

import (
	"../util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetCourses(w http.ResponseWriter, r *http.Request) {

	if resultJson, found := util.Cache("courses"); found {
		fmt.Fprintln(w, resultJson.(string))
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var course []Course
	db.Find(&course)

	var courses = struct {
		Courses []Course `json:"courses"`
	}{course}

	outputJson, err := json.Marshal(courses)
	if err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	fmt.Fprint(w, string(outputJson))

	util.SetCache("courses", string(outputJson))

}

func CreateCourse(w http.ResponseWriter, r *http.Request) {

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	if err := db.Create(&course); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

}

func UpdateCourse(w http.ResponseWriter, r *http.Request) {

	id, err := util.CheckId(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var course Course
	if db.First(&course, id).RowsAffected == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.NotFound))
		return
	}

	var updCourse Course
	if err := json.NewDecoder(r.Body).Decode(&updCourse); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}
	updCourse.Id = id

	if err := db.Save(&updCourse); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

	util.DeleteCache("courses")

}

func DeleteCourse(w http.ResponseWriter, r *http.Request) {

	id, err := util.CheckId(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var sections []Section
	db.Where("course_id = ?", id).Find(&sections)

	var lessons []Lesson
	for _, s := range sections {
		var l []Lesson
		db.Where("section_id = ?", s.Id).Find(&l)
		lessons = append(lessons, l...)
	}

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
	for i := range sections {
		db.Delete(&sections[i])
	}

	var course Course
	if rows := db.Delete(&course, id).RowsAffected; rows == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.NotFound))
		return
	}

	util.DeleteCache("courses")

}
