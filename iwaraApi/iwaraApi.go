package iwaraApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Username string
	Keys     []string
}

func FindUsername(filename *string) (username string) {
	defer func() {
		_, _ = recover().(error)
	}()
	//class="username">清炽</a>
	url := `https://ecchi.iwara.tv/videos/` + *filename
	r, err := Fetch(&url)
	if err != nil {
		log.Println(err)
	}
	reg, _ := regexp.Compile(`class="username">(.+?)</a>`)
	username = reg.FindAllStringSubmatch(string(r), -1)[0][1]
	time.Sleep(time.Second)

	return
}

func GetAllFilenameByUsername(username string) (filename *[]string) {

	URL := "https://ecchi.iwara.tv/users/" + username + "/videos"
	log.Println(URL)
	res, err := http.Get(URL)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("=========================================")
		fmt.Println(res.StatusCode, res.Status)
		fmt.Println("client error: ", err)
		fmt.Println("=========================================")
	}
	page := GetMaxPage(&body) + 1

	var hashs []string
	for i := 0; i < page; i++ {
		ress, errr := http.Get(URL + `?page=` + strconv.Itoa(i))
		if errr != nil {
			log.Println(errr)
			continue
		}
		if ress.StatusCode != http.StatusOK {
			fmt.Println("=========================================")
			fmt.Println(username)
			fmt.Println(ress.StatusCode, ress.Status)
			fmt.Println("=========================================")
		}
		defer ress.Body.Close()
		bodyy, _ := ioutil.ReadAll(ress.Body)
		reg, _ := regexp.Compile(`<a href="/videos/(.+?)(?:[?].+?|["])>`)
		urls := reg.FindAllStringSubmatch(string(bodyy), -1)
		hashs = append(hashs, GetSubMatchData(urls, 1)...)
	}

	fmt.Println("=========================================")
	fmt.Println(fmt.Sprintf("%s =>%s", URL, username))
	fmt.Println(res.StatusCode, res.Status)
	fmt.Println(fmt.Sprintf("find %d keys from %d page", len(hashs), page))
	fmt.Println("=========================================")
	return &hashs

}

func GetUsername(s *string) (uname string) {
	defer func() {
		_, _ = recover().(error)
	}()
	reg, _ := regexp.Compile(`https://ecchi.iwara.tv/users/(.+?)(?:(/videos)|[?]|/|$)`)
	uname = reg.FindAllStringSubmatch(*s, -1)[0][1]
	return
}

func GetBaseUrl(s *string) (err error) {
	defer func() {
		err, _ = recover().(error)
	}()
	reg, _ := regexp.Compile(`(https://ecchi.iwara.tv/users/.+?)(?:(/videos)|[?]|/|$)`)
	*s = reg.FindAllStringSubmatch(*s, -1)[0][1] + "/videos"
	return
}

func GetSubMatchData(sa [][]string, index int) (sr []string) {

	allKeys := make(map[string]bool)

	for _, item := range sa {
		if len(item) < index+1 {
			continue
		}
		if _, value := allKeys[item[index]]; !value {
			allKeys[item[index]] = true
			sr = append(sr, item[index])
		}
	}
	return
}

type downloadUrlStruct struct {
	Resolution string `json:"resolution"`
	Uri        string `json:"uri"`
	Mime       string `json:"mime"`
}

//mikoto 호스트가 가장 빠름 //mikoto.iwara.tv
func GetDownloadUrl(hashs string) (urls string, err error) {
	defer func() { // 함수 빠져나가기 직전 무조건 실행된다
		err, _ = recover().(error) // 프로그램이 죽는경우 살린다
		if err != nil {            // 죽이고 살린 후 처리
			log.Println(err)
		}
	}()
	for {

		res, _ := http.Get("https://ecchi.iwara.tv/api/video/" + hashs)
		body, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()

		var ress []downloadUrlStruct
		_ = json.Unmarshal(body, &ress)

		for i := 0; i < len(ress); i++ {
			if ress[i].Resolution == `Source` && strings.Contains(ress[i].Uri, "mikoto.iwara.tv") {
				urls = `https:` + ress[i].Uri
				return
			}
		}
		time.Sleep(time.Millisecond * 500)
	}

	return
}

func GetMaxPage(body *[]byte) (page int) {
	defer func() {
		err, _ := recover().(error)
		if err != nil {
			page = 0
		}
	}()

	reg, _ := regexp.Compile(`<li class="pager-last last"><a title=".+?" href="/users/.+?/videos\?.*?page=([0-9]{1,3})">`)
	urls := reg.FindAllStringSubmatch(string(*body), -1)
	page, _ = strconv.Atoi(urls[0][1])
	return
}

func Fetch(url *string) (data []byte, err error) {
	defer func() {
		err, _ = recover().(error)
	}()
	res, _ := http.Get(*url)
	if res.StatusCode == http.StatusOK {
		data, _ = ioutil.ReadAll(res.Body)
	}
	defer res.Body.Close()
	return
}
