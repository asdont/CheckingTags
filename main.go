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
	// тело страницы
	body := getBody("urls1.txt")

	// извлечение требуемых тегов
	hrefTags := extractingTags("href", body)
	imgTags := extractingTags("img", body)

	fmt.Println(hrefTags)
	fmt.Println(imgTags)
}

func getBody(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(file)
	strBody := fmt.Sprintf("%s", b)
	strBody = strings.ToLower(strBody)

	return strBody
}

func extractingTags(typeTag string, bodyHtml string) map[string][]bool {
	// добыча тегов

	// для добавления новых регулярок, для случаев ломаных тегов
	const (
		hrefType1 = `<a ([\s\S]*?)a>`         // <a "https://url.ru">text a>
		hrefType2 = `<[\s\S]href([\s\S]*?)a>` // < href="https://url.ru">text a>
		imgType1  = `<img([\s\S]*?)>`
	)

	var tagFromTypes []string
	switch typeTag {
	case "href":
		tagFromTypes = []string{hrefType1, hrefType2}
	case "img":
		tagFromTypes = []string{imgType1}
	}

	countTags := 0
	resultTags := map[string][]bool{}
	// поиск тегов по регуляркам из const
	for _, tagFromType := range tagFromTypes {
		listTags := regexp.MustCompile(tagFromType).FindAllString(bodyHtml, -1)
		countTags += len(listTags)

		// отправка найденых тегов на валидацию
		var checksTag []bool
		for _, tag := range listTags {
			switch typeTag {
			case "href":
				checksTag = validationHrefTag(tag)
			case "img":
				checksTag = validationImgTag(tag)
			}

			// получить позицию тега, если есть ошибки
			tagErrorPosition := errorHandling(tag, checksTag, bodyHtml)
			if tagErrorPosition != "ok" {
				resultTags[tagErrorPosition] = checksTag
			}
		}
	}
	fmt.Printf("найдено %d %s-тега, %d неправильных\n", countTags, strings.ToUpper(typeTag), len(resultTags))

	return resultTags
}

func errorHandling(tag string, tagErrors []bool, bodyHtml string) string {
	// обработка найденых ошибок в тегах
	bodyHtmlLines := strings.Split(bodyHtml, "\n")

	for i := range tagErrors {
		// если есть ошибки в теге
		if !tagErrors[i] {
			positionTag := strings.Index(bodyHtml, tag)

			// поиск позиции тега
			for i, stroke := range bodyHtmlLines {
				positionTag = positionTag - len(stroke) - 1
				if positionTag <= 0 {
					positionTag = positionTag + len(stroke) + 1

					// формирование ключа map
					tagInform := fmt.Sprintf("строка - %d столбец - %d, тег - %s\n", i+1, positionTag, tag)
					tagInform = strings.Join(strings.Fields(tagInform), " ")

					return tagInform
				}
			}
		}
	}
	return "ok"
}

func displayingErrors(tagType string, tagErrors []bool) {
	// отображение(расшифровка) ошибок
	switch tagType {
	case "href":
		fmt.Println("проверки ->", tagErrors)
		if !tagErrors[0] {
			fmt.Println("отсутствует открывающий тег <a ")
		}
		if !tagErrors[1] {
			fmt.Println("отсутствует атрибут href")
		}
		if !tagErrors[2] {
			fmt.Println("не найден url в атрибуте href или href не соответсвует требованиям")
		}
		if !tagErrors[3] {
			fmt.Println("отсутствует текст ссылки(невидимая ссылка) или неверно указан href")
		}
		if !tagErrors[4] {
			fmt.Println("отсутствует закрывающий тег </a>")
		}

		if !tagErrors[5] {
			fmt.Println("отсутствует указание протокола http или https")
		}
		if !tagErrors[6] {
			fmt.Println("отустствуют разделители или меньше чем требуется")
		}
		if !tagErrors[7] {
			fmt.Println("короткая длина ссылки, возможно ссылка неверна")
		}
		if !tagErrors[8] {
			fmt.Println("в ссылке присутствуют запрещённые символы или неверно указан href")
		}
		if !tagErrors[9] {
			fmt.Println("в ссылке не найдена точка, возможно так не должно быть")
		}
	case "img":
		fmt.Println("проверки ->", tagErrors)
		if !tagErrors[0] {
			fmt.Println("отсутствует открывающий тег <img ")
		}
		if !tagErrors[1] {
			fmt.Println("отсутствует атрибут src")
		}
		if !tagErrors[2] {
			fmt.Println("не найден url в атрибуте src или src не соответсвует требованиям")
		}
		if !tagErrors[3] {
			fmt.Println("отсутствует закрывающий тег >")
		}

		if !tagErrors[4] {
			fmt.Println("отсутствует указание протокола http или https")
		}
		if !tagErrors[5] {
			fmt.Println("отустствуют разделители или меньше чем требуется")
		}
		if !tagErrors[6] {
			fmt.Println("короткая длина ссылки, возможно ссылка неверна")
		}
		if !tagErrors[7] {
			fmt.Println("в ссылке присутствуют запрещённые символы")
		}
		if !tagErrors[8] {
			fmt.Println("в ссылке не найдена точка, возможно так не должно быть")
		}
	}
	fmt.Print("\n")
}
