package util

import (
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

func SetLogger() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
	})

	log.AddHook(filename.NewHook())

	dirPath := "./log"

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, 0755); err != nil {
			log.Fatalln("Error:", err)
		}
	}

	filePath := dirPath + "/pygo_server." +
		time.Now().Format("02_Jan_2006") +
		".log"

	logFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	setRequestLogger()

}

var requestLogger = log.New()

func setRequestLogger() {

	requestLogger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
	})

	dirPath := "./log"

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, 0755); err != nil {
			log.Fatalln("Error:", err)
		}
	}

	filePath := dirPath + "/pygo_http_request." +
		time.Now().Format("02_Jan_2006") +
		".log"

	logFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	mw := io.MultiWriter(os.Stdout, logFile)
	requestLogger.SetOutput(mw)

}

func LogRequest(handler http.Handler) http.HandlerFunc {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			ip, err := getIp(r)
			if err != nil {
				SendErr(w, http.StatusBadRequest, err)
				return
			}

			requestLogger.WithFields(log.Fields{
				"IP":     ip,
				"Method": r.Method,
				"URL":    r.URL,
			}).Infoln("HTTP Request")

			handler.ServeHTTP(w, r)

		})

}

func getIp(req *http.Request) (net.IP, error) {

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, err
	}

	userIp := net.ParseIP(ip)
	if userIp == nil {
		return nil, err
	}

	return userIp, nil

}
