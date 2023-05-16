package catu

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/go-catupiry/catu/helpers"
	"github.com/go-catupiry/catu/pagination"
	"github.com/go-catupiry/query_parser_to_db"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RequestContextOpts struct {
	EchoContext echo.Context
}

func NewRequestContext(opts *RequestContextOpts) *RequestContext {
	app := GetApp()
	cfg := app.GetConfiguration()
	ctx := RequestContext{
		App:         app,
		EchoContext: opts.EchoContext,
		// Title:               "",
		Theme:  app.GetTheme(),
		Layout: app.GetLayout(),
		ENV:    cfg.GetF("GO_ENV", "development"),
		Query:  query_parser_to_db.NewQuery(50),
		Pager:  pagination.NewPager(),
	}

	// Is a context used on CLIs, not in HTTP request / echo then skip it
	if opts.EchoContext == nil || ctx.Request().URL == nil {
		return &ctx
	}

	ctx.Pager.CurrentUrl = ctx.Request().URL.Path
	ctx.Pager.Limit, _ = strconv.ParseInt(cfg.GetF("PAGER_LIMIT", "20"), 10, 64)

	if opts.EchoContext.Request().Method != "GET" {
		return &ctx
	}

	limitMax, _ := strconv.ParseInt(app.GetConfiguration().GetF("PAGER_LIMIT_MAX", "50"), 10, 64)

	rawParams := opts.EchoContext.QueryParams()

	filteredParamArray := []string{}

	for key, param := range rawParams {
		// get limit with max value for security:
		if key == "limit" && len(param) == 1 {
			queryLimit, err := strconv.ParseInt(param[0], 10, 64)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"key":   key,
					"param": param,
				}).Error("NewRequestContext invalid query param limit")
				continue
			}
			if queryLimit > 0 && queryLimit < limitMax {
				ctx.Pager.Limit = queryLimit
			}
		}

		if key == "page" && len(param) == 1 {
			page, _ := strconv.ParseInt(param[0], 10, 64)
			ctx.Pager.Page = page
			continue
		}

		ctx.Query.AddQueryParamFromRaw(key, param)
	}

	if len(filteredParamArray) > 0 {
		strings.Join(filteredParamArray[:], ",")
	}

	return &ctx
}

type RequestContext struct {
	EchoContext echo.Context
	App         App

	Title string

	IsAuthenticated   bool
	AuthenticatedUser UserInterface
	// authenticated user role name list
	Roles []string

	Session SessionData

	Widgets   map[string]map[string]string
	Theme     string
	Layout    string
	BodyClass []string
	Content   template.HTML
	Query     query_parser_to_db.QueryInterface
	Pager     *pagination.Pager

	ENV string
}

/// --- Start echo.Context overrides

// Request returns `*http.Request`.
func (c *RequestContext) Request() *http.Request {
	return c.EchoContext.Request()
}

// SetRequest sets `*http.Request`.
func (c *RequestContext) SetRequest(r *http.Request) {
	c.EchoContext.SetRequest(r)
}

// SetResponse sets `*Response`.
func (c *RequestContext) SetResponse(r *echo.Response) {
	c.EchoContext.SetResponse(r)
}

// Response returns `*Response`.
func (c *RequestContext) Response() *echo.Response {
	return c.EchoContext.Response()
}

// IsTLS returns true if HTTP connection is TLS otherwise false.
func (c *RequestContext) IsTLS() bool {
	return c.EchoContext.IsTLS()
}

// IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
func (c *RequestContext) IsWebSocket() bool {
	return c.EchoContext.IsWebSocket()
}

// Scheme returns the HTTP protocol scheme, `http` or `https`.
func (c *RequestContext) Scheme() string {
	return c.EchoContext.Scheme()
}

// RealIP returns the client's network address based on `X-Forwarded-For`
// or `X-Real-IP` request header.
// The behavior can be configured using `Echo#IPExtractor`.
func (c *RequestContext) RealIP() string {
	return c.EchoContext.RealIP()
}

// Path returns the registered path for the handler.
func (c *RequestContext) Path() string {
	return c.EchoContext.Path()
}

// SetPath sets the registered path for the handler.
func (c *RequestContext) SetPath(p string) {
	c.EchoContext.SetPath(p)
}

// Param returns path parameter by name.
func (c *RequestContext) Param(name string) string {
	return c.EchoContext.Param(name)
}

// ParamNames returns path parameter names.
func (c *RequestContext) ParamNames() []string {
	return c.EchoContext.ParamNames()
}

// SetParamNames sets path parameter names.
func (c *RequestContext) SetParamNames(names ...string) {
	c.EchoContext.SetParamNames(names...)
}

// ParamValues returns path parameter values.
func (c *RequestContext) ParamValues() []string {
	return c.EchoContext.ParamValues()
}

// SetParamValues sets path parameter values.
func (c *RequestContext) SetParamValues(values ...string) {
	c.EchoContext.SetParamValues(values...)
}

// QueryParam returns the query param for the provided name.
func (c *RequestContext) QueryParam(name string) string {
	return c.EchoContext.QueryParam(name)
}

// QueryParams returns the query parameters as `url.Values`.
func (c *RequestContext) QueryParams() url.Values {
	return c.EchoContext.QueryParams()
}

// QueryString returns the URL query string.
func (c *RequestContext) QueryString() string {
	return c.EchoContext.QueryString()
}

// FormValue returns the form field value for the provided name.
func (c *RequestContext) FormValue(name string) string {
	return c.EchoContext.FormValue(name)
}

// FormParams returns the form parameters as `url.Values`.
func (c *RequestContext) FormParams() (url.Values, error) {
	return c.EchoContext.FormParams()
}

// FormFile returns the multipart form file for the provided name.
func (c *RequestContext) FormFile(name string) (*multipart.FileHeader, error) {
	return c.EchoContext.FormFile(name)
}

// MultipartForm returns the multipart form.
func (c *RequestContext) MultipartForm() (*multipart.Form, error) {
	return c.EchoContext.MultipartForm()
}

// Cookie returns the named cookie provided in the request.
func (c *RequestContext) Cookie(name string) (*http.Cookie, error) {
	return c.EchoContext.Cookie(name)
}

// SetCookie adds a `Set-Cookie` header in HTTP response.
func (c *RequestContext) SetCookie(cookie *http.Cookie) {
	c.EchoContext.SetCookie(cookie)
}

// Cookies returns the HTTP cookies sent with the request.
func (c *RequestContext) Cookies() []*http.Cookie {
	return c.EchoContext.Cookies()
}

// Bind binds the request body into provided type `i`. The default binder
// does it based on Content-Type header.
func (c *RequestContext) Bind(i interface{}) error {
	return c.EchoContext.Bind(i)
}

// Validate validates provided `i`. It is usually called after `Context#Bind()`.
// Validator must be registered using `Echo#Validator`.
func (c *RequestContext) Validate(i interface{}) error {
	return c.EchoContext.Validate(i)
}

// Render renders a template with data and sends a text/html response with status
// code. Renderer must be registered using `Echo.Renderer`.
func (c *RequestContext) Render(code int, name string, data interface{}) error {
	return c.EchoContext.Render(code, name, data)
}

// HTML sends an HTTP response with status code.
func (c *RequestContext) HTML(code int, html string) error {
	return c.EchoContext.HTML(code, html)
}

// HTMLBlob sends an HTTP blob response with status code.
func (c *RequestContext) HTMLBlob(code int, b []byte) error {
	return c.EchoContext.HTMLBlob(code, b)
}

// String sends a string response with status code.
func (c *RequestContext) String(code int, s string) error {
	return c.EchoContext.String(code, s)
}

// JSON sends a JSON response with status code.
func (c *RequestContext) JSON(code int, i interface{}) error {
	return c.EchoContext.JSON(code, i)
}

// JSONPretty sends a pretty-print JSON with status code.
func (c *RequestContext) JSONPretty(code int, i interface{}, indent string) error {
	return c.EchoContext.JSONPretty(code, i, indent)
}

// JSONBlob sends a JSON blob response with status code.
func (c *RequestContext) JSONBlob(code int, b []byte) error {
	return c.EchoContext.JSONBlob(code, b)
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (c *RequestContext) JSONP(code int, callback string, i interface{}) error {
	return c.EchoContext.JSONP(code, callback, i)
}

// JSONPBlob sends a JSONP blob response with status code. It uses `callback`
// to construct the JSONP payload.
func (c *RequestContext) JSONPBlob(code int, callback string, b []byte) error {
	return c.EchoContext.JSONPBlob(code, callback, b)
}

// XML sends an XML response with status code.
func (c *RequestContext) XML(code int, i interface{}) error {
	return c.EchoContext.XML(code, i)
}

// XMLPretty sends a pretty-print XML with status code.
func (c *RequestContext) XMLPretty(code int, i interface{}, indent string) error {
	return c.EchoContext.XMLPretty(code, i, indent)
}

// XMLBlob sends an XML blob response with status code.
func (c *RequestContext) XMLBlob(code int, b []byte) error {
	return c.EchoContext.XMLBlob(code, b)
}

// Blob sends a blob response with status code and content type.
func (c *RequestContext) Blob(code int, contentType string, b []byte) error {
	return c.EchoContext.Blob(code, contentType, b)
}

// Stream sends a streaming response with status code and content type.
func (c *RequestContext) Stream(code int, contentType string, r io.Reader) error {
	return c.EchoContext.Stream(code, contentType, r)
}

// File sends a response with the content of the file.
func (c *RequestContext) File(file string) error {
	return c.EchoContext.File(file)
}

// Attachment sends a response as attachment, prompting client to save the
// file.
func (c *RequestContext) Attachment(file string, name string) error {
	return c.EchoContext.Attachment(file, name)
}

// Inline sends a response as inline, opening the file in the browser.
func (c *RequestContext) Inline(file string, name string) error {
	return c.EchoContext.Inline(file, name)
}

// NoContent sends a response with no body and a status code.
func (c *RequestContext) NoContent(code int) error {
	return c.EchoContext.NoContent(code)
}

// Redirect redirects the request to a provided URL with status code.
func (c *RequestContext) Redirect(code int, url string) error {
	return c.EchoContext.Redirect(code, url)
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *RequestContext) Error(err error) {
	c.EchoContext.Error(err)
}

// Handler returns the matched handler by router.
func (c *RequestContext) Handler() echo.HandlerFunc {
	return c.EchoContext.Handler()
}

// SetHandler sets the matched handler by router.
func (c *RequestContext) SetHandler(h echo.HandlerFunc) {
	c.EchoContext.SetHandler(h)
}

// Logger returns the `Logger` instance.
func (c *RequestContext) Logger() echo.Logger {
	return c.EchoContext.Logger()
}

// Set the logger
func (c *RequestContext) SetLogger(l echo.Logger) {
	c.EchoContext.SetLogger(l)
}

// Echo returns the `Echo` instance.
func (c *RequestContext) Echo() *echo.Echo {
	return c.EchoContext.Echo()
}

/// --- END echo.Context overrides

// Reset resets the context after request completes. It must be called along
// with `Echo#AcquireContext()` and `Echo#ReleaseContext()`.
// See `Echo#ServeHTTP()`
func (c *RequestContext) Reset(r *http.Request, w http.ResponseWriter) {
	c.EchoContext.Reset(r, w)
}

type SessionData struct {
	UserID string
}

func (r *RequestContext) Set(name string, value interface{}) {
	r.EchoContext.Set(name, value)
}

func (r *RequestContext) Get(name string) interface{} {
	return r.EchoContext.Get(name)
}

func (r *RequestContext) GetString(name string) string {
	v := r.EchoContext.Get(name)
	if v == nil {
		return ""
	}
	return v.(string)
}

// Get value from echo context data in boolean format
func (r *RequestContext) GetBool(name string) bool {
	v := r.Get(name)
	if v == nil {
		return false
	}

	return v.(bool)
}

// Get data in string map format ([]string) from echo context data
func (r *RequestContext) GetStringMap(name string) []string {
	v := r.Get(name)
	if v == nil {
		return []string{}
	}

	return v.([]string)
}

func (r *RequestContext) GetTemplateHTML(name string) template.HTML {
	v := r.Get(name)
	if v == nil {
		return template.HTML("")
	}

	return v.(template.HTML)
}

func (r *RequestContext) RenderPagination(name string) string {
	html := ""
	return html
}

// Render one template, with support for themes
func (r *RequestContext) RenderTemplate(wr io.Writer, name string, data interface{}) error {
	return r.App.GetTemplates().ExecuteTemplate(wr, path.Join(r.Theme, name), data)
}

// Partial - Include and render one template inside other
func (r *RequestContext) Partial(name string, data interface{}) template.HTML {
	var htmlBuffer bytes.Buffer
	err := r.RenderTemplate(&htmlBuffer, name, data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"partialName": name,
			"error":       fmt.Sprintf("%+v\n", err),
		}).Error("catu.Partial error on render partial template")
		return template.HTML("")
	}

	return template.HTML(htmlBuffer.String())
}

// Add a body class string checking if is unique
func (r *RequestContext) AddBodyClass(class string) {
	if helpers.SliceContains(r.BodyClass, class) {
		return
	}

	r.BodyClass = append(r.BodyClass, class)
}

// Remove a body class string checking if is unique
func (r *RequestContext) RemoveBodyClass(class string) {
	if !helpers.SliceContains(r.BodyClass, class) {
		return
	}

	r.BodyClass = append(r.BodyClass, class)
}

// Get body class as string,
func (r *RequestContext) GetBodyClassText() string {
	return strings.Join(r.BodyClass, " ")
}

// Get selected response type
func (r *RequestContext) GetResponseContentType() string {
	v := r.GetString("responseContentType")
	if v == "" {
		return r.Request().Header.Get(echo.HeaderContentType) // default ...
	}

	return v
}

// Set response type, returns error if the type is invalid
func (r *RequestContext) SetResponseContentType(v string) error {
	if v == "" {
		return errors.New("RequestContext.SetResponseContentType value should not be empty")
	}

	r.Set("responseContentType", v)
	return nil
}

func (r *RequestContext) GetLimit() int {
	return int(r.Pager.Limit)
}

func (r *RequestContext) GetOffset() int {
	page := int(r.Pager.Page)

	if page < 2 {
		return 0
	}

	limit := int(r.Pager.Limit)
	return limit * (page - 1)
}

func (r *RequestContext) ParseQueryFromReq(c echo.Context) error {
	return nil
}

func (r *RequestContext) GetAuthenticatedRoles() *[]string {
	if r.IsAuthenticated {
		return &r.Roles
	}

	return &[]string{"unAuthenticated"}
}

func (r *RequestContext) SetAuthenticatedUser(user UserInterface) {
	r.AuthenticatedUser = user
	r.IsAuthenticated = true
}

func (r *RequestContext) SetAuthenticatedUserAndFillRoles(user UserInterface) {
	r.SetAuthenticatedUser(user)
	r.Roles = user.GetRoles()
	r.Roles = append(r.Roles, "authenticated")
}

func (r *RequestContext) Can(permission string) bool {
	roles := r.GetAuthenticatedRoles()
	return r.App.Can(permission, *roles)
}

func GetQueryIntFromReq(param string, c echo.Context) int {
	var err error
	var valueInt int
	page := c.QueryParam(param)
	if page != "" {
		valueInt, err = strconv.Atoi(page)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path":  c.Path(),
				"param": param,
				"page":  page,
			}).Warn("NewRequestRequestContext invalid page query param")
		}
	}

	return valueInt
}

func GetQueryInt64FromReq(param string, c echo.Context) int64 {
	var err error
	var valueInt int64
	value := c.QueryParam(param)
	if value != "" {
		valueInt, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path":  c.Path(),
				"param": param,
				"value": value,
			}).Warn("GetQueryInt64FromReq invalid page query param")
		}
	}

	return valueInt
}
