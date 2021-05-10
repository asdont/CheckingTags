//func main() {
//	url := "https://mnemag.ru/articles/12-samyh-neozhidannyh-allergii.html"
//	getBody(url)
//}
//
//func getBody(url string) {
//	req, err := http.Get(url)
//
//	if err != nil {
//		panic("не удалось получить тело страницы")
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			panic("не удалось закрыть соединение")
//		}
//	}(req.Body)
//
//	bodyHtml, err := ioutil.ReadAll(req.Body)
//	if err != nil {
//		panic("не удалось прочитать тело страницы")
//	}
//
//	extractingTags(bodyHtml)
//}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	file, err := os.Open(`urls2.txt`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(file)
	str := fmt.Sprintf("%s", b)

	extractingTags(str)
}

func extractingTags(bodyHtml string) {
	// добыча тегов
	const (
		hrefType1 = `<a ([\s\S]*?)a>`         // <a "https://url.ru">text a>
		hrefType2 = `<[\s\S]href([\s\S]*?)a>` // < href="https://url.ru">text a>
		imgType1  = `<img([\s\S]*?)>`
	)

	// для добавления новых регулярок, для случаев ломаных тегов
	typesHref := []string{hrefType1, hrefType2}
	typesImg := []string{imgType1}

	//bodyHtmlClean := strings.Join(strings.Fields(bodyHtml), " ")

	countHref := 0
	for _, typeHref := range typesHref {
		listHref := regexp.MustCompile(typeHref).FindAllString(bodyHtml, -1)
		countHref = countHref + len(listHref)
		for _, hrefTag := range listHref {
			link, checksTag := validationHrefTag(hrefTag)
			errorHandling("href", hrefTag, link, checksTag, bodyHtml)
		}
	}
	fmt.Printf("Найдено %d <a href>\n", countHref)

	countImg := 0
	for _, typeImg := range typesImg {
		listImg := regexp.MustCompile(typeImg).FindAllString(bodyHtml, -1)
		countImg = countImg + len(listImg)
		for _, imgTag := range listImg {
			link, checksTag := validationImgTag(imgTag)
			errorHandling("img", imgTag, link, checksTag, bodyHtml)
		}
	}
	fmt.Printf("Найдено %d <img src>\n", countImg)
}

func errorHandling(tagType string, tag string, link string, tagErrors []bool, bodyHtml string) {
	// обработка найденых ошибок в тегах
	bodyHtmlLines := strings.Split(bodyHtml, "\n")
	//bodyHtml = strings.Trim(bodyHtml, "\n")

	for i := range tagErrors {
		if !tagErrors[i] {
			positionTag := strings.Index(bodyHtml, tag)

			for i, stroke := range bodyHtmlLines {
				positionTag = positionTag - len(stroke) - 1
				if positionTag <= 0 {
					positionTag = positionTag + len(stroke) + 1
					fmt.Printf("строка - %d столбец - %d, тег - %s\n", i+1, positionTag, tag)

					switch tagType {
					case "href":
						fmt.Print("href\n")
					case "img":
						fmt.Print("img")
					}

					break
				}
			}
		}
	}

}

func validationHrefTag(tagHref string) (string, []bool) {
	tagHref = strings.ToLower(tagHref)

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

	correctUrlHref := validationURL(tagUrl)
	for _, boolResult := range correctUrlHref {
		correctnessTagHref = append(correctnessTagHref, boolResult)
	}

	return linkHref, correctnessTagHref
}

func validationImgTag(tagImg string) (string, []bool) {
	tagImg = strings.ToLower(tagImg)

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

	correctUrlImg := validationURL(tagUrl)
	for _, boolResult := range correctUrlImg {
		correctnessTagImg = append(correctnessTagImg, boolResult)
	}

	return linkImg, correctnessTagImg
}

func validationURL(tagURL string) []bool {
	// валидация url относительного(<protocol>://<domain>) и абсолютного типа(/<path>/<path>)
	var existProtocol, existDomain, existSeparators, existSybolsAllowed, existDot bool

	existProtocol = true                                                           // наличие протокола
	existSeparators = true                                                         // наличие разделителей
	existDomain = len(tagURL) > 0                                                  // длина ссылки, например vk.ru(5)
	existSybolsAllowed = !regexp.MustCompile(`[^\w:/.#-=?]|'`).MatchString(tagURL) // проверка запрещенных символов
	existDot = true                                                                // наличие точки

	isAbsolute := regexp.MustCompile(`:/|/{2}|[htps]{4,5}`).MatchString(tagURL)
	if isAbsolute {
		existProtocol = strings.Contains(tagURL, "http")
		existSeparators = strings.Count(tagURL, "/") >= 2
		existDot = strings.Contains(tagURL, ".") // наличие точки

	}

	return []bool{existProtocol, existSeparators, existDomain, existSybolsAllowed, existDot}
}
