package catu

import (
	"encoding/json"

	"github.com/pkg/errors"

	"net/http"

	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
)

type BaseListReponse struct {
	Meta BaseMetaResponse `json:"meta"`
}

type BaseMetaResponse struct {
	Count int64 `json:"count"`
}

type BaseErrorResponse struct {
	Messages []BaseErrorResponseMessage `json:"messages"`
}

type BaseErrorResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func SendManyJSONApi(c echo.Context, count int64, pointerArr interface{}) error {
	p, err := jsonapi.Marshal(pointerArr)
	if err != nil {
		return errors.Wrap(err, "SendManyJSONApi error on jsonapi.Marshal")
	}

	payload, ok := p.(*jsonapi.ManyPayload)
	if !ok {
		return errors.New("SendManyJSONApi was not a many payloader")
	}

	payload.Meta = &jsonapi.Meta{
		"count": count,
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)

	return json.NewEncoder(c.Response()).Encode(payload)
}

func SendOneJSONApi(c echo.Context, record interface{}) error {
	p, err := jsonapi.Marshal(record)
	if err != nil {
		return errors.Wrap(err, "SendOneJSONApi error on jsonapi.Marshal")
	}

	payload, ok := p.(*jsonapi.OnePayload)
	if !ok {
		return errors.New("SendOneJSONApi was not a many payloader")
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)

	return json.NewEncoder(c.Response()).Encode(payload)
}
