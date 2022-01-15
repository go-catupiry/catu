package catu

import (
	"html/template"

	"github.com/go-catupiry/catu/helpers"
	"github.com/go-catupiry/catu/pagination"
	"github.com/sirupsen/logrus"
)

// func renderCachedBlock(templateKey string) template.HTML {
// 	// TODO! add a small time cache in this redis getter ...

// 	// const key = 'HTML-quotes-teaser-b-'+(symbol.replace('BVMF:', ''));
// 	// return siteCache.getAsync(key);

// 	// const key = 'HTML-quote-graph-'+(symbol.replace('BVMF:', ''));
// 	// return siteCache.getAsync(key);

// 	html, err := cache.GetItem(templateKey)
// 	if err != nil {
// 		logrus.WithFields(logrus.Fields{
// 			"templateKey": templateKey,
// 			"error":       err,
// 		}).Warn("renderCachedBlock error on get item from cache")
// 	}
// 	return template.HTML(html)
// }

func noEscapeHTML(str string) template.HTML {
	return template.HTML(str)
}

func paginate(pager *pagination.Pager, queryString string) template.HTML {
	return renderPager(pager, queryString)
}

// func imageTPLHelper(image models.Image, style, class, width string) template.HTML {
// 	html := ""

// 	url := image.GetUrl(style)

// 	if url != "" {
// 		html += `<img`

// 		if image.Description != nil {
// 			html += ` alt="` + *image.Description + `"`
// 		}

// 		html += ` src="` + url + `"`

// 		if class != "" {
// 			html += ` class="` + class + `"`
// 		}

// 		if width != "" {
// 			html += ` width="` + width + `"`
// 		}

// 		html += `>`
// 	}

// 	return template.HTML(html)
// }

// func imagesTPLHelper(images *[]models.Image, style string) template.HTML {

// 	return template.HTML("")
// }

type ContentDates interface {
	GetTeaserDatesHTML(separator string) template.HTML
}

func contentDates(record ContentDates, separator string) template.HTML {
	return record.GetTeaserDatesHTML(separator)
}

func truncate(text string, length int, ellipsis string) template.HTML {
	html, err := helpers.Truncate(text, length, ellipsis)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"text":     text,
			"length":   length,
			"ellipsis": ellipsis,
		}).Error("truncate error on truncate text")
	}
	return html
}

func shareMenu(url, title string, symbol *string, btnClass, btnGroupClass string) template.HTML {
	app := GetApp()

	html := ""
	siteName := app.Configuration.Get("SITE_NAME")
	twitterUsername := app.Configuration.Get("TWITTER_USER_NAME")
	symbolStr := ""

	if symbol != nil {
		symbolStr = *symbol
	}

	if btnClass == "" {
		btnClass = "btn btn-light btn-sm blog-post-share dropdown-toggle"
	}

	if btnGroupClass == "" {
		btnGroupClass = "share-dropdown-menu btn-group pull-right"
	}

	if symbolStr != "" {
		title += ` (` + symbolStr + `) `
	}

	html += `<div class="` + btnGroupClass + `">
		<button type="button" class="` + btnClass + `" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
			<i class="zmdi zmdi-share"></i>
		</button>
		<div class="dropdown-menu">
			<a href="https://twitter.com/intent/tweet?url=` + url + `&original_referer=` + url + `&tw_p=tweetbutton&via=` + twitterUsername + `&text=` + title + ` - " class="dropdown-item s-twitter" target="_blank">
				<img src="/public/theme/material-blog/icons/twitter.png" alt="Compartilhar no Twitter">
			</a>
			<a href="https://www.facebook.com/sharer/sharer.php?u=` + url + `" class="dropdown-item s-facebook" target="_blank">
				<img src="/public/theme/material-blog/icons/facebook.png" alt="Compartilhar no Facebook">
			</a>
			<a href="https://www.linkedin.com/shareArticle?mini=true&url=` + url + `&title=` + title + `&summary=${desc}" class="dropdown-item s-linkedin" target="_blank">
				<img src="/public/theme/material-blog/icons/linkedin.png" alt="Compartilhar no Linkedin">
			</a>
			<a href="whatsapp://send?text=` + title + ` - ` + siteName + `: ` + url + `" data-action="share/whatsapp/share" class="dropdown-item s-whatsapp-mobile">
				<img src="/public/theme/material-blog/icons/whatsapp.png" alt="Compartilhar no Whatsapp">
			</a>
			<a href="https://web.whatsapp.com/send?text=` + title + ` - ` + siteName + `: ` + url + `" data-action="share/whatsapp/share" class="dropdown-item s-whatsapp-site" target="_blank">
				<img src="/public/theme/material-blog/icons/whatsapp.png" alt="Compartilhar no Whatsapp">
			</a>
			<a href="mailto:?subject=` + title + ` - ` + siteName + `&amp;body=` + title + ` no link: ` + url + `" title="Compartilhar por Email" class="dropdown-item">
				<img width="50px" src="/public/theme/material-blog/icons/email.png" alt="Compartilhar por e-mail">
			</a>
		</div>
	</div>`

	return template.HTML(html)
}
