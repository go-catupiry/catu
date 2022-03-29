package http_client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
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
		logrus.WithFields(logrus.Fields{
			"err": fmt.Sprintf("%+v\n", err),
		}).Debug("catu.GetPageHTML error")
		return "", errors.Wrap(err, "GetPageHTML error")
	}

	return string(bodyBytes), nil
}
