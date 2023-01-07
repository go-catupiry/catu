package catu

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

type EmptyResponse struct{}
