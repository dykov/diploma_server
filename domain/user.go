package domain

import (
	"../util"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"
)

func GetUser(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	if resultJson, found := util.Cache("user/" + idString); found {
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

	var user User
	if rows := db.First(&user, id).RowsAffected; rows == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.UserNotFound))
		return
	}
	user.Password = ""

	outputJson, err := json.Marshal(user)
	if err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	fmt.Fprint(w, string(outputJson))

	util.SetCache("user/"+idString, string(outputJson))

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	var res string
	db.Raw("select password_validation( ? )", user.Password).Row().Scan(&res)
	if res != "" {
		util.SendErr(w, http.StatusBadRequest, errors.New(res))
		return
	}
	db.Raw("select login_validation( ? )", user.Login).Row().Scan(&res)
	if res != "" {
		util.SendErr(w, http.StatusBadRequest, errors.New(res))
		return
	}

	user.Password = util.HashPassword(user.Password, util.Salt)
	user.VerificationCode = randomString()

	if err := db.Create(&user).Scan(&user); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

	link := "http://localhost:80/user/" + strconv.FormatUint(user.Id, 10) + "/verification?code=" + user.VerificationCode

	// если ошибка - ссылка должна отобразиться на странице клиента
	if sendEmail(user.Login, user.Email, link) != nil {
		util.SendErr(w, http.StatusBadRequest, errors.New(link))
	}

}

func randomString() string {

	var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	stringLength := 30
	str := make([]rune, stringLength)
	for i, _ := range str {
		str[i] = runes[rand.Intn(len(runes))]
	}

	return string(str)

}

func sendEmail(login, email, link string) error {

	auth1 := smtp.PlainAuth("", "onaft1pygo@gmail.com", "onaft_pygoonaft_pygo", "smtp.gmail.com")

	to := []string{email}
	fmt.Println(to)
	msg := []byte("To: " + email + "\r\n" +
		"Subject: PyGo verification\r\n" +
		"\r\n" +
		"Dear, " + login + " !" +
		"\nPlease, verify your account with the link:" +
		"\n" + link + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth1, "onaft1pygo@gmail.com", to, msg)
	fmt.Println(msg)
	if err != nil {
		logrus.WithField("Error", err).Errorln("Email sending error")
		return err
	}

	return nil

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	id, err := util.CheckId(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var user User
	if rows := db.First(&user, id).RowsAffected; rows == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.UserNotFound))
		return
	}

	var updUser User
	if err := json.NewDecoder(r.Body).Decode(&updUser); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	var res string
	db.Raw("select password_validation( ? )", updUser.Password).Row().Scan(&res)
	if res != "" {
		util.SendErr(w, http.StatusBadRequest, errors.New(res))
		return
	}
	db.Raw("select login_validation( ? )", updUser.Login).Row().Scan(&res)
	if res != "" {
		util.SendErr(w, http.StatusBadRequest, errors.New(res))
		return
	}

	updUser.Id = id
	updUser.Password = util.HashPassword(updUser.Password, util.Salt)

	if err := db.Save(&updUser); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

	util.DeleteCache("user/" + mux.Vars(r)["id"])
	util.DeleteCache("all_users")

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	id, err := util.CheckId(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	db.Where("user_id = ?", id).Delete(&UsersTest{})

	if db.Delete(&User{}, id).RowsAffected == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.UserNotFound))
		return
	}

	util.DeleteCache("user/" + mux.Vars(r)["id"])
	util.DeleteCache("all_users")

}

func GetUsersTests(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	if resultJson, found := util.Cache("user/" + idString + "/test"); found {
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

	if rows := db.First(&User{}, id).RowsAffected; rows == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.UserNotFound))
		return
	}

	var usersTest []UsersTest
	db.Where("user_id = ?", id).Find(&usersTest)

	var userTests = struct {
		TestId []uint64 `json:"test_id"`
	}{}

	for i := range usersTest {
		userTests.TestId = append(userTests.TestId, usersTest[i].TestId)
	}

	outputJson, err := json.Marshal(userTests)
	if err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	fmt.Fprint(w, string(outputJson))

	util.SetCache("user/"+idString+"/test", string(outputJson))

}

func AddUsersTest(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]
	id, err := util.CheckId(w, idString)
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var user User
	if db.First(&user, id).RowsAffected == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.UserNotFound))
		return
	}

	// Get JSON like: {"id":1 , "points":10}
	var test ParagraphsOrTest
	if err := json.NewDecoder(r.Body).Decode(&test); err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}
	user.Rating += test.Points
	// Save user changes
	if err := db.Save(&user); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

	//Add record to UsersTest
	var userTest = UsersTest{id, test.Id}
	if err := db.Create(&userTest); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

	util.DeleteCache("user/" + idString + "/test")
	util.DeleteCache("all_users")

}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	if resultJson, found := util.Cache("all_users"); found {
		fmt.Fprintln(w, resultJson.(string))
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var users []User
	db.Find(&users)
	for i := range users {
		users[i].Password = ""
	}

	var allUsers = struct {
		Users []User `json:"users"`
	}{users}

	outputJson, err := json.Marshal(allUsers)
	if err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	fmt.Fprint(w, string(outputJson))

	util.SetCache("all_users", string(outputJson))

}

func OnaftReview(w http.ResponseWriter, r *http.Request) {

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var users []User
	if rows := db.Where("is_onaft_student = ?", 1).Find(&users).RowsAffected; rows == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.UserNotFound))
		return
	}

	var u = struct {
		Users []User `json:"users"`
	}{users}

	outputJson, err := json.Marshal(u)
	if err != nil {
		util.SendErr(w, http.StatusBadRequest, err)
		return
	}

	fmt.Fprint(w, string(outputJson))

}

func UserVerification(w http.ResponseWriter, r *http.Request) {

	id, err := util.CheckId(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}

	db := util.GetDB(w)
	if db.DB() == nil {
		return
	}

	var user User
	if rows := db.First(&user, id).RowsAffected; rows == 0 {
		util.SendErr(w, http.StatusBadRequest, errors.New(util.UserNotFound))
		return
	}

	if user.VerificationCode != r.URL.Query().Get("code") {
		util.SendErr(w, http.StatusBadRequest, errors.New("Wrong verification code"))
		return
	}

	user.Role = util.GrantUser
	user.VerificationCode = ""

	if err := db.Save(&user); err.Error != nil {
		util.SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

	util.DeleteCache("user/" + mux.Vars(r)["id"])
	util.DeleteCache("all_users")

}
