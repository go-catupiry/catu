package catu

import (
	"html/template"
	"io"
	"strconv"
	"strings"

	"github.com/go-catupiry/catu/helpers"
	"github.com/go-catupiry/catu/pagination"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AppContext struct {
	PathBeforeAlias string
	Protocol        string
	Domain          string
	AppOrigin       string
	Title           string

	IsAuthenticated   bool
	AuthenticatedUser UserInterface
	// authenticated user role name list
	Roles []string

	Session SessionData

	// Widgets     map[string]map[string]string
	Layout              string
	BodyClass           []string
	Content             template.HTML
	ContentData         map[string]interface{}
	Query               Query
	Pager               *pagination.Pager
	MetaTags            HTMLMetaTags
	ResponseContentType string

	ENV string
}

type SessionData struct {
	UserID int64
}

func (r *AppContext) Set(name string, value interface{}) {
	r.ContentData[name] = value
}

func (r *AppContext) Get(name string) interface{} {
	return r.ContentData[name]
}

func (r *AppContext) GetString(name string) string {
	if r.ContentData[name] == nil {
		return ""
	}
	return r.ContentData[name].(string)
}

func (r *AppContext) GetBool(name string) bool {
	return r.ContentData[name].(bool)
}

func (r *AppContext) GetStringMap(name string) []string {
	return r.ContentData[name].([]string)
}

func (r *AppContext) GetTemplateHTML(name string) template.HTML {
	return r.ContentData[name].(template.HTML)
}

func (r *AppContext) RenderPagination(name string) string {
	html := ""

	return html
}

// Render one template, alias to app.templates.ExecuteTemplate()
func (r *AppContext) Render(wr io.Writer, name string, data interface{}) error {
	app := GetApp()

	return app.GetTemplates().ExecuteTemplate(wr, name, data)
}

// Add a body class string checking if is unique
func (r *AppContext) AddBodyClass(class string) {
	if helpers.SliceContains(r.BodyClass, class) {
		return
	}

	r.BodyClass = append(r.BodyClass, class)
}

// Remove a body class string checking if is unique
func (r *AppContext) RemoveBodyClass(class string) {
	if !helpers.SliceContains(r.BodyClass, class) {
		return
	}

	r.BodyClass = append(r.BodyClass, class)
}

// Get body class as string,
func (r *AppContext) GetBodyClassText() string {
	return strings.Join(r.BodyClass, " ")
}

func (r *AppContext) GetLimit() int {
	return int(r.Pager.Limit)
}

func (r *AppContext) GetOffset() int {
	page := int(r.Pager.Page)

	if page < 2 {
		return 0
	}

	limit := int(r.Pager.Limit)
	return limit * (page - 1)
}

func (r *AppContext) ParseQueryFromReq(c echo.Context) error {
	return nil
}

func (r *AppContext) GetAuthenticatedRoles() *[]string {
	if r.IsAuthenticated {
		roles := r.AuthenticatedUser.GetRoles()
		return &roles
	}

	return &[]string{"unAuthenticated"}
}

func (r *AppContext) SetAuthenticatedUser(user UserInterface) {
	r.AuthenticatedUser = user
	r.IsAuthenticated = true
}

func (r *AppContext) SetAuthenticatedUserAndFillRoles(user UserInterface) {
	r.SetAuthenticatedUser(user)
	r.Roles = user.GetRoles()
	r.Roles = append(r.Roles, "authenticated")
}

func (r *AppContext) Can(permission string) bool {
	// roles := r.GetAuthenticatedRoles()
	// log.Println("roles:", roles)
	// return acl.Can(permission, *roles)
	// TODO!
	return true
}

func NewAppContext() AppContext {
	app := GetApp()

	port := app.Configuration.GetF("PORT", "8080")
	protocol := app.Configuration.GetF("PROTOCOL", "http")
	domain := app.Configuration.GetF("DOMAIN", "localhost")

	ctx := AppContext{
		Protocol:  protocol,
		Domain:    domain,
		AppOrigin: app.Configuration.GetF("APP_ORIGIN", protocol+"://"+domain+":"+port),
		// Title:               "",
		// ResponseContentType: "text/html",
		Layout: "site/layouts/default",
		ENV:    app.Configuration.GetF("GO_ENV", "development"),
	}

	ctx.Pager = pagination.NewPager()
	ctx.Pager.Limit, _ = strconv.ParseInt(app.Configuration.GetF("PAGER_LIMIT", "20"), 10, 64)
	ctx.ContentData = map[string]interface{}{}

	ctx.MetaTags.Title = "Monitor do Mercado"
	ctx.MetaTags.Description = " Com tecnologia de ponta, o site Monitor do Mercado reúne as mais importantes informações sobre investimentos."
	ctx.MetaTags.ImageURL = "https://storage.googleapis.com/mm-images/static/favicon.ico"

	return ctx
}

func NewRequestAppContext(c echo.Context) AppContext {
	app := GetApp()
	ctx := NewAppContext()
	ctx.Pager.CurrentUrl = c.Request().URL.Path

	if c.Request().Method != "GET" {
		return ctx
	}

	limitMax, _ := strconv.ParseInt(app.Configuration.GetF("PAGER_LIMIT_MAX", "50"), 10, 64)

	rawParams := c.QueryParams()

	filteredParamArray := []string{}

	for key, param := range rawParams {
		// get limit with max value for security:
		if key == "limit" && len(param) == 1 {
			queryLimit, err := strconv.ParseInt(param[0], 10, 64)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"key":   key,
					"param": param,
				}).Error("NewRequestAppContext invalid query param limit")
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

	return ctx
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
			}).Warn("NewRequestAppContext invalid page query param")
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

func GetPathLimitFromReq() {

}

func (r *AppContext) RenderMetaTags() template.HTML {
	app := GetApp()
	html := ""

	pageUrl := r.AppOrigin + r.PathBeforeAlias

	if pageUrl != "" {
		html += `<meta property="og:url" content="` + pageUrl + `" />`
		html += `<link rel="canonical" href="` + pageUrl + `" />`
	}

	siteName := app.Configuration.Get("SITE_NAME")

	if siteName != "" {
		html += `<meta property="og:site_name" content="` + siteName + `" />`
		html += `<meta content="` + siteName + `" itemprop="name">`
	}

	html += `<meta content="` + siteName + `" name="twitter:site">`
	html += `<meta property="og:type" content="website" />`

	if r.MetaTags.Description != "" {
		html += `<meta property="og:description" content="` + r.MetaTags.Description + `" />`
		html += `<meta content="` + r.MetaTags.Description + `">`
		html += `<meta content="` + r.MetaTags.Description + `" name="twitter:description">`
	}

	if r.MetaTags.Title != "" {
		html += `<meta content="` + r.MetaTags.Title + `" name="twitter:title">`
		html += `<meta property="og:title" content="` + r.MetaTags.Title + `" />`
	}

	if r.MetaTags.ImageURL != "" {
		html += `<meta property="og:image" content="` + r.MetaTags.ImageURL + `" />`
	}

	if r.MetaTags.Keywords != "" {
		html += `<meta name="keywords" content="` + r.MetaTags.Keywords + `" />`
	}

	return template.HTML(html)
}
