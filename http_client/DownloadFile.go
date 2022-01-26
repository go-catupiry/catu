package http_client

import (
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// DownloadFile - Download one file
func DownloadFile(url string, dest *os.File, headers http.Header) (bool, error) {
	logrus.WithFields(logrus.Fields{
		"url":  url,
		"dest": dest.Name(),
	}).Debug("DownloadFile will download")

	var err error

	defer dest.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header = headers

	res, err := HttpClient.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":     url,
			"headers": headers,
			"error":   err,
		}).Error("DownloadFile error")
		return false, err
	}

	defer res.Body.Close()

	_, err = io.Copy(dest, res.Body)
	if err != nil {
		return false, err
	}

	logrus.WithFields(logrus.Fields{
		"url":  url,
		"dest": dest.Name(),
	}).Debug("DownloadFile done download")

	return true, err
}
