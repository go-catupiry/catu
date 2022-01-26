package http_client

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

func GetPageHTML(url string, headers http.Header) (string, error) {
	resp, err := Get(url, headers)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":     url,
			"headers": headers,
			"error":   err,
		}).Error("GetPageHTML error")
		return "", err
	}

	defer resp.Body.Close()

	rdrBody := io.Reader(resp.Body)
	bodyBytes, err := ioutil.ReadAll(rdrBody)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(bodyBytes), nil
}
