package test_report

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PostmanReport struct {
	Run        PostmanRun        `json:"run"`
	Collection PostmanCollection `json:"collection"`
}

type PostmanCollection struct {
	Item []PostmanItem `json:"item"`
}

type PostmanRun struct {
	Executions []PostmanExecution `json:"executions"`
}

type PostmanExecution struct {
	Item       PostmanItem     `json:"item,omitempty"`
	Request    PostmanRequest  `json:"request,omitempty"`
	Response   PostmanResponse `json:"response,omitempty"`
	Id         string          `json:"id,omitempty"`
	FileOrigin string
	JsonPath   string
}

type PostmanItem struct {
	Name string        `json:"name,omitempty"`
	Id   string        `json:"id"`
	Item []PostmanItem `json:"item"`
}

type PostmanURL struct {
	Protocol string              `json:"protocol,omitempty"`
	Port     string              `json:"port,omitempty"`
	Path     []string            `json:"path,omitempty"`
	Host     []string            `json:"host,omitempty"`
	Query    []PostmanQueryParam `json:"query,omitempty"`
	Variable []any               `json:"variable,omitempty"`
}

func (u PostmanURL) GetUrl() string {
	var sb strings.Builder

	if u.Protocol != "" {
		sb.WriteString(u.Protocol + "://")
	}
	sb.WriteString(u.Host[0])

	if u.Port != "" {
		sb.WriteString(":" + u.Port)
	}
	sb.WriteString("/")
	sb.WriteString(strings.Join(u.Path, "/"))

	if len(u.Query) > 0 {
		start := true
		for _, elem := range u.Query {
			if !elem.Disabled {
				if start {
					sb.WriteString("?")
					start = false
				} else {
					sb.WriteString("&")
				}
				sb.WriteString(elem.Key)
				sb.WriteString("=")
				sb.WriteString(elem.Value)
			}
		}
	}

	return sb.String()
}

type PostmanQueryParam struct {
	Disabled bool   `json:"disabled"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type PostmanHeader struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type PostmanRequest struct {
	URL    PostmanURL      `json:"url,omitempty"`
	Header []PostmanHeader `json:"header,omitempty"`
	Body   PostmanBody     `json:"body,omitempty"`
	Method string          `json:"method,omitempty"`
}

type PostmanBody struct {
	Raw string `json:"raw"`
}

type Header struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type PostmanStream string

func (p *PostmanStream) UnmarshalJSON(data []byte) error {
	var jsonStruct postmanStreamJson
	err := json.Unmarshal(data, &jsonStruct)
	if err != nil {
		return err
	}
	if jsonStruct.Type != "Buffer" {
		return fmt.Errorf("%s stream is not supported", jsonStruct.Type)
	}
	*p = PostmanStream(jsonStruct.Data)
	return nil
}

type postmanStreamJson struct {
	Type string `json:"type,omitempty"`
	Data []byte `json:"data,omitempty"`
}

type PostmanResponse struct {
	ID           string          `json:"id,omitempty"`
	Code         int             `json:"code,omitempty"`
	Header       []PostmanHeader `json:"header,omitempty"`
	Stream       PostmanStream   `json:"stream,omitempty"`
	ResponseSize int             `json:"responseSize,omitempty"`
}
