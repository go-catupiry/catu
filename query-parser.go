package catu

import (
	"strings"
)

type QueryAttr struct {
	Operator   string
	Values     []string
	IsMultiple bool
	ParamName  string
}

type Query struct {
	Fields      []QueryAttr
	QueryString string
}

func (r *Query) AddQueryParamFromRaw(paramName string, values []string) error {
	if len(values) == 0 {
		return nil
	}

	if paramName == "page" {
		return nil
	}

	r.AddQueryString(paramName, values)

	var qAttr QueryAttr

	if len(values) > 1 {
		qAttr.IsMultiple = true
	}

	if !strings.Contains(paramName, "_") {
		qAttr.Values = values
		qAttr.ParamName = paramName
		qAttr.Operator = "="
		r.Fields = append(r.Fields, qAttr)
		return nil
	}

	return nil
}

func (r *Query) AddQueryString(paramName string, values []string) {
	if r.QueryString != "" {
		r.QueryString += "&"
	}

	if len(values) > 1 {
		for i := range values {
			r.QueryString += paramName + "[]=" + values[i]
		}
	} else {
		r.QueryString += paramName + "=" + values[0]
	}
}

func (r *Query) GetQueryString(paramName string) string {
	for i := range r.Fields {
		if r.Fields[i].ParamName == paramName {
			if len(r.Fields[i].Values) == 0 {
				return ""
			} else if len(r.Fields[i].Values) == 1 {
				return paramName + `=` + r.Fields[i].Values[0]
			} else {
				var results []string
				for vi := range r.Fields[i].Values {
					result := paramName + "[]=" + r.Fields[i].Values[vi]
					results = append(results, result)
				}

				return strings.Join(results, "&")
			}
		}
	}

	return ""
}

func (r *Query) GetParamValue(paramName string) string {
	for i := range r.Fields {
		if r.Fields[i].ParamName == paramName {
			if len(r.Fields[i].Values) != 0 {
				return r.Fields[i].Values[0]
			}
		}
	}

	return ""
}
