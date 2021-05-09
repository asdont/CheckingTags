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
	file, err := os.Open(`urls1.txt`)
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

//type checkTagsHref struct {
//	numberType, numberHref int
//	url string
//	validTag, validUrl []bool
//}
//
//func (tag *checkTagsHref) result()  {
//	res := map[string]string{}
//	fmt.Println(tag.numberType, tag.numberHref, tag.url, tag.validTag, tag.validUrl)
//}

func extractingTags(bodyString string) {
	// добыча тегов
	const (
		hrefType1 = `<a ([\s\S]*?)a>`         // <a "https://url.ru">text a>
		hrefType2 = `<[\s\S]href([\s\S]*?)a>` // < href="https://url.ru">text</a>
		imgType1  = `<img([\s\S]*?)>`
	)

	// для добавления новых регулярок, для случаев ломаных тегов
	typesHref := []string{hrefType1, hrefType2}
	typesImg := []string{imgType1}

	resHref := map[string][]bool{}
	for _, typeHref := range typesHref {
		listHref := regexp.MustCompile(typeHref).FindAllString(bodyString, -1)
		for _, href := range listHref {
			tagHref, checksTag := validationHref(cleanedTag(href))
			resHref[tagHref] = checksTag
		}
	}

	resImg := map[string][]bool{}
	for _, typeImg := range typesImg {
		listImg := regexp.MustCompile(typeImg).FindAllString(bodyString, -1)
		for _, img := range listImg {
			tagImg, checksTag := validationImg(cleanedTag(img))
			resImg[tagImg] = checksTag
		}
	}
}

func cleanedTag(tag string) string {
	// очистка тегов от избыточных символов
	tag = strings.Replace(tag, "\n", "", -1)
	tag = strings.Join(strings.Fields(tag), " ")
	return tag
}

func validationHref(tagHref string) (string, []bool) {
	tagHref = strings.Replace(tagHref, "'", "\"", -1)

	link := regexp.MustCompile(`[href]{3,4}="[\s\S]*?"`).FindString(tagHref)

	// валидация тега
	var existABegin, existHref, existUrl, existText, existAEnd bool

	existABegin = strings.HasPrefix(tagHref, "<a ")                              // наличие открывающегося тега a
	existHref = strings.Contains(tagHref, "href")                                // наличие атрибута href
	existUrl = len(link) > 7                                                           // url указан, href="" (7)
	existText = len(regexp.MustCompile(`>[\s\S]*?</a>`).FindString(tagHref)) > 5 // наличие текста ссылки, ></a> (5)
	existAEnd = strings.HasSuffix(tagHref, "</a>")                               // наличие закрывающего тега a

	// валидация url в тэге
	tagUrl := strings.Replace(link, "href=", "", 1)
	tagUrl = strings.Trim(tagUrl, "\"")

	correctTag := []bool{existABegin, existHref, existUrl, existText, existAEnd}

	correctUrl := validationURL(tagUrl)
	for _, boolResult := range correctUrl {
		correctTag = append(correctTag, boolResult)
	}

	return tagHref, correctTag
}

func validationImg(tagImg string) (string, []bool) {
	tagImg = strings.Replace(tagImg, "'", "\"", -1)
	fmt.Println(tagImg)
	link := regexp.MustCompile(`[src]{2,4}="[\s\S]*?"`).FindString(tagImg)
	fmt.Println(link)

	tagSrc := strings.Replace(link, "src=", "", 1)
	tagSrc = strings.Trim(tagSrc, "\"")
	fmt.Println(tagSrc)

	return "", []bool{}
}

func validationURL(tagURL string) []bool {
	// валидация url относительного(<protocol>://<domain>) и абсолютного типа(/<path>/<path>)
	isAbsolute := regexp.MustCompile(`:/|/{2}|[htps]{4,5}|[.]`).MatchString(tagURL)

	var existProtocol, existDomain, existSeparators, existSybolsAllowed bool
	if isAbsolute {
		existProtocol = strings.Contains(tagURL, "http")                             // наличие протокола
		existDomain = len(regexp.MustCompile(`://[\s\S]*?"`).FindString(tagURL)) > 4 // ://url.ru"(>4)
		existSeparators = strings.Count(tagURL, "/") > 2                             // наличие разделителей
		existSybolsAllowed = !regexp.MustCompile(`[^\w\d:/.#-]`).MatchString(tagURL) // проверка запрещенных символов
	} else {
		existProtocol = true
		existDomain = true
		existSeparators = true
		existSybolsAllowed = !regexp.MustCompile(`[^\w\d:/.#-]`).MatchString(tagURL) // проверка запрещенных символов
	}

	return []bool{existProtocol, existDomain, existSeparators, existSybolsAllowed}
}


