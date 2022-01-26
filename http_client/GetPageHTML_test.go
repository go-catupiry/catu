package http_client

// var mocks_folder = "../../../_mocks"

// func init() {
// 	HttpClient = &MockClient{}
// }

// func TestGetPageHTML(t *testing.T) {
// 	t.Run("GetPageHTML should work with mock and return the mocked data", func(t *testing.T) {
// 		htmlByte, err := ioutil.ReadFile(mocks_folder + "/cvm/dfp-list.html")
// 		assert.Nil(t, err)

// 		r := ioutil.NopCloser(bytes.NewReader([]byte(htmlByte)))

// 		GetDoFunc = func(*http.Request) (*http.Response, error) {
// 			return &http.Response{
// 				StatusCode: 200,
// 				Body:       r,
// 			}, nil
// 		}

// 		url := "http://dados.cvm.gov.br/dados/CIA_ABERTA/DOC/DFP/DADOS/"
// 		var headers http.Header

// 		body, err := GetPageHTML(url, headers)
// 		assert.NotNil(t, body)
// 		assert.Nil(t, err)

// 		t.Log("body length", len(body))
// 		assert.EqualValues(t, string(htmlByte), body)
// 	})
// }
