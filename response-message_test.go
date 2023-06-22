package catu

import (
	"html/template"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetResponseMessage(t *testing.T) {
	app := GetAppInstance()
	e := app.GetRouter()

	type args struct {
		key     string
		message *ResponseMessage
	}
	tests := []struct {
		name          string
		savedMessages map[string]*ResponseMessage
		args          args
		wantErr       bool
		expectedLen   int
	}{
		{
			name: "Should set a valid data",
			args: args{
				key: "4",
				message: &ResponseMessage{
					Type:    "success",
					Message: "Test message",
				},
			},
			expectedLen: 1,
		},
		{
			name: "Should add a seccond message",
			savedMessages: map[string]*ResponseMessage{
				"1": {
					Type:    "success",
					Message: "Test message",
				},
			},
			args: args{
				key: "2",
				message: &ResponseMessage{
					Type:    "error",
					Message: "Test message 2",
				},
			},
			expectedLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			if tt.savedMessages != nil {
				c.Set(responseMessageKey, tt.savedMessages)
			}

			err := SetResponseMessage(c, tt.args.key, tt.args.message)
			if err != nil {
				assert.Equal(t, tt.wantErr, true)
			}

			setMessages := c.Get(responseMessageKey).(map[string]*ResponseMessage)
			assert.Equal(t, tt.expectedLen, len(setMessages))

			for k, m := range tt.savedMessages {
				if setMessages[k] != nil {
					assert.Equal(t, setMessages[k].Type, m.Type)
					assert.Equal(t, setMessages[k].Message, m.Message)
				} else {
					assert.Equal(t, setMessages[k], tt.args.message)
				}
			}
		})
	}
}

func TestGetResponseMessages(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*ResponseMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResponseMessages(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResponseMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResponseMessages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseMessagesRender(t *testing.T) {
	type args struct {
		c   echo.Context
		tpl string
	}
	tests := []struct {
		name string
		args args
		want template.HTML
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResponseMessagesRender(tt.args.c, tt.args.tpl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResponseMessagesRender() = %v, want %v", got, tt.want)
			}
		})
	}
}
