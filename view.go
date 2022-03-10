package catu

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-catupiry/catu/pagination"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// type tplWrapper struct {
// 	Ctx *AppContext
// }

type TemplateCTX struct {
	EchoContext echo.Context
	Ctx         interface{}
	Record      interface{}
	Records     interface{}
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	switch v := data.(type) {
	case int:
		// v is an int here, so e.g. v + 1 is possible.
		fmt.Printf("Integer: %v", v)
	case float64:
		// v is a float64 here, so e.g. v + 1.0 is possible.
		fmt.Printf("Float64: %v", v)
	case string:
		// v is a string here, so e.g. v + " Yeah!" is possible.
		fmt.Printf("String: %v", v)
	default:
		htmlContext := data.(*TemplateCTX)
		htmlContext.EchoContext = c

		logrus.WithFields(logrus.Fields{
			"name":          name,
			"htmlContext":   htmlContext,
			"len templates": len(t.templates.Templates()),
		}).Debug("Render")

		var contentBuffer bytes.Buffer
		err := t.templates.ExecuteTemplate(&contentBuffer, name, htmlContext)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"name":  name,
			}).Error("catu.theme.Render error on execute template")
			return err
		}

		ctx := htmlContext.Ctx.(*AppContext)
		ctx.Content = template.HTML(contentBuffer.String())

		var layoutBuffer bytes.Buffer
		err = t.templates.ExecuteTemplate(&layoutBuffer, ctx.Layout, htmlContext)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"name":   name,
				"layout": ctx.Layout,
			}).Error("catu.theme.Render error on execute layout template")
			return err
		}

		ctx.Content = template.HTML(layoutBuffer.String())

		return t.templates.ExecuteTemplate(w, "site/html", htmlContext)
	}

	return nil
}

func findAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if info != nil && !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			name = strings.Replace(name, ".html", "", 1)

			t := root.New(name).Funcs(funcMap)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

func renderPager(r *pagination.Pager, queryString string) template.HTML {
	var htmlBuffer bytes.Buffer
	app := GetApp()
	templates := app.GetTemplates()

	logrus.WithFields(logrus.Fields{
		"count": r.Count,
		"Pager": string(r.ToJSON()),
	}).Debug("paginate params")

	if r.Count == 0 {
		return template.HTML("")
	}

	currentUrl := r.CurrentUrl
	queryParamsStr := ""

	if queryString != "" {
		queryParamsStr += "&" + queryString
	}

	pageCountFloat := float64(r.Count) / float64(r.Limit)
	pageCount := int64(math.Ceil(pageCountFloat))
	totalLinks := (r.MaxLinks * 2) + 1
	startInPage := int64(1)
	endInPage := pageCount

	if pageCount == 0 {
		return template.HTML("")
	}

	// logrus.WithFields(logrus.Fields{
	// 	"pageCount":  pageCount,
	// 	"totalLinks": totalLinks,
	// 	"MaxLinks":   r.MaxLinks,
	// 	"Page":       r.Page,
	// 	"before":     r.MaxLinks+2 < r.Page,
	// 	"after":      r.MaxLinks+r.Page+1 < pageCount,
	// }).Debug("Calculing 1>>>")

	if totalLinks < pageCount {
		if r.MaxLinks+2 < r.Page {
			startInPage = r.Page - r.MaxLinks
			r.FirstPath = currentUrl + "?page=1" + queryParamsStr
			r.FirstNumber = "1"
			r.HasMoreBefore = true
		}

		if (r.MaxLinks + r.Page + 1) < pageCount {
			endInPage = r.MaxLinks + r.Page
			r.LastPath = currentUrl + "?page=" + strconv.FormatInt(pageCount, 10) + queryParamsStr
			r.LastNumber = strconv.FormatInt(pageCount, 10)
			r.HasMoreAfter = true
		}
	}

	// Each link
	for i := startInPage; i <= endInPage; i++ {
		number := strconv.FormatInt(i, 10)
		var link = pagination.Link{
			Path:   currentUrl + "?page=" + number + queryParamsStr,
			Number: number,
		}

		if i == r.Page {
			link.IsActive = true
		}

		// logrus.WithFields(logrus.Fields{
		// 	"i":    i,
		// 	"Page": r.Page,
		// }).Debug("Calculing afterEach")

		r.Links = append(r.Links, link)
	}

	if r.Page > 1 {
		r.HasPrevius = true
		number := strconv.FormatInt(r.Page-1, 10)
		r.PreviusPath = currentUrl + "?page=" + number + queryParamsStr
		r.PreviusNumber = number
	}

	if r.Page < pageCount {
		r.HasNext = true
		number := strconv.FormatInt(r.Page+1, 10)
		r.NextPath = currentUrl + "?page=" + number + queryParamsStr
		r.NextNumber = number
	}

	// logrus.WithFields(logrus.Fields{
	// 	"pagger":      string(r.ToJSON()),
	// 	"startInPage": startInPage,
	// 	"endInPage":   endInPage,
	// }).Debug("Calculing end")

	err := templates.ExecuteTemplate(&htmlBuffer, "site/components/paginate", TemplateCTX{
		Ctx: &r,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pagger": &r,
			"error":  err,
		}).Error("theme.paginate Error on render template")
	}

	return template.HTML(htmlBuffer.String())
}
