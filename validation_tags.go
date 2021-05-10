package main

import (
	"regexp"
	"strings"
)

func validationHrefTag(tagHref string) []bool {
	linkHref := regexp.MustCompile(`[href]{3,4}="[\s\S]*?"`).FindString(tagHref) // href="link.url"

	// валидация тега a href
	var existABegin, existHref, existUrl, existText, existAEnd bool

	// проверки тега A
	existABegin = strings.HasPrefix(tagHref, "<a ")                              // наличие открывающегося тега a
	existHref = strings.Contains(tagHref, "href")                                // наличие атрибута href
	existUrl = len(linkHref) > 7                                                 // наличие url, href="" (7)
	existText = len(regexp.MustCompile(`>[\s\S]*?</a>`).FindString(tagHref)) > 5 // наличие текста ссылки, ></a> (5)
	existAEnd = strings.HasSuffix(tagHref, "</a>")                               // наличие закрывающего тега a

	// валидация url в тэге a href
	tagUrl := strings.Replace(linkHref, "href=", "", 1)
	tagUrl = strings.Trim(tagUrl, "\"")

	correctnessTagHref := []bool{existABegin, existHref, existUrl, existText, existAEnd}

	// объединение результатов проверки тега и урл []bool + []bool
	correctUrlHref := validationURL(tagUrl)
	for _, boolResult := range correctUrlHref {
		correctnessTagHref = append(correctnessTagHref, boolResult)
	}

	return correctnessTagHref
}

func validationImgTag(tagImg string) []bool {
	linkImg := regexp.MustCompile(`[src]{3}="[\s\S]*?"`).FindString(tagImg) // src="link.url"

	// валидация тега img src
	var existImg, existSrc, existUrl, existEnding bool

	// проверки тега IMG
	existImg = strings.HasPrefix(tagImg, "<img ") // наличие открывающего тега <img
	existSrc = strings.Contains(tagImg, "src")    // наличие атрибута src
	existUrl = len(linkImg) > 6                   // наличие ссылки (src="")(6)
	existEnding = strings.HasSuffix(tagImg, ">")  // закрытие тега

	// валидация url в тэге img
	tagUrl := strings.Replace(linkImg, "src=", "", 1)
	tagUrl = strings.Trim(tagUrl, "\"")

	correctnessTagImg := []bool{existImg, existSrc, existUrl, existEnding}

	// объединение результатов проверки тега и урл []bool + []bool
	correctUrlImg := validationURL(tagUrl)
	for _, boolResult := range correctUrlImg {
		correctnessTagImg = append(correctnessTagImg, boolResult)
	}

	return correctnessTagImg
}

func validationURL(tagURL string) []bool {
	// валидация url относительного(<protocol>://<domain>) и абсолютного типа(/<path>/<path>)
	var existProtocol, existDomain, existSeparators, existSybolsAllowed, existDot bool

	existProtocol = true                                                           // наличие протокола
	existSeparators = true                                                         // наличие разделителей
	existDot = true                                                                // наличие точки
	existDomain = len(tagURL) > 0                                                  // длина ссылки, например vk.ru(5)
	existSybolsAllowed = !regexp.MustCompile(`[^\w:/.#-=?]|'`).MatchString(tagURL) // проверка запрещенных символов

	// если ссылка абсолютная, изменяем проверки
	isAbsolute := regexp.MustCompile(`:/|/{2}|[htps]{4,5}`).MatchString(tagURL)
	if isAbsolute {
		existProtocol = strings.Contains(tagURL, "http")
		existSeparators = strings.Count(tagURL, "/") >= 2
		existDot = strings.Contains(tagURL, ".") // наличие точки

	}

	return []bool{existProtocol, existSeparators, existDomain, existSybolsAllowed, existDot}
}
