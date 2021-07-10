package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("-", 100))

	//body, err := getBodyFromFile("testdata/urlsTest2.txt")
	//if err != nil {
	//	log.Fatalf("get body page: %v", err)
	//}

	bodyPage, err := getBodyFromUrl("https://pkg.go.dev/golang.org/x/net/html")
	if err != nil {
		log.Fatalf("get body page: %v", err)
	}

	// извлечение невалидных тегов href
	hrefTags, err := extractingTags("a", bodyPage)
	if err != nil {
		log.Fatalf("extracting tags <a: %v", err)
	}

	// извлечение невалидных тегов img
	imgTags, err := extractingTags("img", bodyPage)
	if err != nil {
		log.Fatalf("extracting tags <img: %v", err)
	}

	// отображение информации о ошибках
	displayingErrors("a", hrefTags)
	displayingErrors("img", imgTags)

	fmt.Println(strings.Repeat("-", 100))
}

// html по url
func getBodyFromUrl(url string) (string, error) {
	req, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("get the page body: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("close request - FAIL")
		}
	}(req.Body)

	bodyHtml, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", fmt.Errorf("read body(html): %v", err)
	}

	return string(bodyHtml), nil
}

// html из файла
func getBodyFromFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("open file - %s: %v", fileName, err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("close file - %s: %v", fileName, err)
		}
	}()

	body, err := ioutil.ReadAll(file)
	strBody := strings.ToLower(string(body))

	return strBody, nil
}

// добыча тегов
func extractingTags(typeTag string, bodyHtml string) (map[string][]bool, error) {
	// для добавления новых регулярок, для случаев ломаных тегов
	const (
		aType1   = `<a ([\s\S]*?)a>`         // <a "https://url.ru">text a>
		aType2   = `<[\s\S]href([\s\S]*?)a>` // < href="https://url.ru">text a>
		imgType1 = `<img([\s\S]*?)>`
	)

	var tagFromTypes []string
	switch typeTag {
	case "a":
		tagFromTypes = []string{aType1, aType2}
	case "img":
		tagFromTypes = []string{imgType1}
	default:
		return nil, fmt.Errorf("unknown tag - %s: ", typeTag)
	}

	resultTags := map[string][]bool{}

	// поиск тегов по регуляркам из const
	countTags := 0
	for _, tagFromType := range tagFromTypes {
		listTags := regexp.MustCompile(tagFromType).FindAllString(bodyHtml, -1)
		countTags += len(listTags)

		// отправка найденых тегов на валидацию
		var checksTag []bool
		for _, tag := range listTags {
			switch typeTag {
			case "a":
				checksTag = validationATag(tag)
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
	fmt.Printf("%s tags found - %d, %d - invalid\n", strings.ToUpper(typeTag), countTags, len(resultTags))

	return resultTags, nil
}

// обработка найденых ошибок в тегах
func errorHandling(tag string, tagErrors []bool, bodyHtml string) string {
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

					// формирование ключа карты(map)
					tagInform := fmt.Sprintf("Row - %d, Column - %d, Tag - %s\n", i+1, positionTag, tag)
					tagInform = strings.Join(strings.Fields(tagInform), " ")

					return tagInform
				}
			}
		}
	}
	return "ok"
}

// отображение(расшифровка) ошибок
func displayingErrors(tagType string, tags map[string][]bool) {
	for tag, tagErrors := range tags {
		switch tagType {
		case "a":
			fmt.Println("\n", tag, "\n", "Errors ->", tagErrors)
			if !tagErrors[0] {
				fmt.Println(">> отсутствует открывающий тег <a ")
			}
			if !tagErrors[1] {
				fmt.Println(">> отсутствует атрибут href")
			}
			if !tagErrors[2] {
				fmt.Println(">> не найден url в атрибуте href или href не соответсвует требованиям")
			}
			if !tagErrors[3] {
				fmt.Println(">> отсутствует текст ссылки(невидимая ссылка) или неверно указан href")
			}
			if !tagErrors[4] {
				fmt.Println(">> отсутствует закрывающий тег </a>")
			}

			if !tagErrors[5] {
				fmt.Println(">> отсутствует указание протокола http или https")
			}
			if !tagErrors[6] {
				fmt.Println(">> отустствуют разделители или меньше чем требуется")
			}
			if !tagErrors[7] {
				fmt.Println(">> короткая длина ссылки, возможно ссылка неверна")
			}
			if !tagErrors[8] {
				fmt.Println(">> в ссылке присутствуют запрещённые символы или неверно указан href")
			}
			if !tagErrors[9] {
				fmt.Println(">> в ссылке не найдена точка, возможно так не должно быть")
			}
		case "img":
			fmt.Println("проверки ->", tagErrors)
			if !tagErrors[0] {
				fmt.Println(">> отсутствует открывающий тег <img ")
			}
			if !tagErrors[1] {
				fmt.Println(">> отсутствует атрибут src")
			}
			if !tagErrors[2] {
				fmt.Println(">> не найден url в атрибуте src или src не соответсвует требованиям")
			}
			if !tagErrors[3] {
				fmt.Println(">> отсутствует закрывающий тег >")
			}

			if !tagErrors[4] {
				fmt.Println(">> отсутствует указание протокола http или https")
			}
			if !tagErrors[5] {
				fmt.Println(">> отустствуют разделители или меньше чем требуется")
			}
			if !tagErrors[6] {
				fmt.Println(">> короткая длина ссылки, возможно ссылка неверна")
			}
			if !tagErrors[7] {
				fmt.Println(">> в ссылке присутствуют запрещённые символы")
			}
			if !tagErrors[8] {
				fmt.Println(">> в ссылке не найдена точка, возможно так не должно быть")
			}
		}
	}
}
