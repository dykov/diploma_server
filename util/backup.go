package util

import (
	"encoding/json"
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type backup struct {
	Time     time.Time `json:"time"`
	Filepath string    `json:"filepath"`
	Name     string    `json:"name"`
}

func Backup(w http.ResponseWriter, r *http.Request) {

	var file backup
	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
		SendErr(w, http.StatusBadRequest, err)
		return
	}
	file.Name += ".dump"
	file.Time = time.Now()

	db := GetDB(w)
	if db.DB() == nil {
		return
	}
	var existFile backup
	if db.Limit(1).Order("time desc").Find(&existFile).RowsAffected != 0 &&
		existFile.Time.Add(12*time.Hour).After(time.Now()) {
		SendErr(w, http.StatusBadRequest, errors.New("Backup was created less than 12 hours ago."))
		return
	}

	if _, err := os.Stat(file.Filepath); os.IsNotExist(err) {
		if err = os.MkdirAll(file.Filepath, 0755); err != nil {
			SendErr(w, http.StatusBadRequest, err)
			log.Errorln("Error:", err)
			return
		}
	}

	cmd := exec.Command(
		"pg_dump",
		"-h", host,
		"-p", strconv.Itoa(port),
		"-d", dbname,
		"-U", user,
		"-f", file.Filepath+string(filepath.Separator)+file.Name,
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+password)

	if err := cmd.Run(); err != nil {
		SendErr(w, http.StatusBadRequest, err)
		log.WithField("Error", err).Errorln("Unable to create backup file")
	}

	if err := db.Create(&file); err.Error != nil {
		SendErr(w, http.StatusBadRequest, err.Error)
		return
	}

}

func Restore(w http.ResponseWriter, r *http.Request) {

	var file backup
	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
		SendErr(w, http.StatusBadRequest, err)
		return
	}
	file.Name += ".dump"

	if _, err := os.Stat(file.Filepath); os.IsNotExist(err) {
		log.Fatalln("Error:", err)
		SendErr(w, http.StatusBadRequest, err)
		return
	}

	cmd := exec.Command(
		"pg_restore",
		"-h", host,
		"-p", strconv.Itoa(port),
		"-d", dbname,
		"-U", user,
		"-v", file.Filepath+string(filepath.Separator)+file.Name,
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+password)

	if err := cmd.Run(); err != nil {
		SendErr(w, http.StatusBadRequest, err)
		log.WithField("Error", err).Errorln("Unable to restore backup file")
	}

}
