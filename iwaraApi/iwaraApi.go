package iwaraApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

const downloadBaseUrl = `https://ecchi.iwara.tv/api/video/`

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

func GetDownloadUrl(hashs string) (urls string, err error) {
	defer func() { // 함수 빠져나가기 직전 무조건 실행된다
		err, _ = recover().(error) // 프로그램이 죽는경우 살린다
		if err != nil {            // 죽이고 살린 후 처리
			fmt.Println(err)
		}
	}()
	res, _ := http.Get(downloadBaseUrl + hashs)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var ress []downloadUrlStruct
	_ = json.Unmarshal(body, &ress)

	for i := 0; i < len(ress); i++ {
		if ress[i].Resolution == `Source` {
			urls = `https:` + ress[i].Uri
			return
		}
	}
	return
}

func GetMaxPage(body *[]byte) (page int) {
	defer func() { // 함수 빠져나가기 직전 무조건 실행된다
		err, _ := recover().(error) // 프로그램이 죽는경우 살린다
		if err != nil {             // 죽이고 살린 후 처리
			fmt.Println(err)
			page = 0
		}
	}()

	reg, _ := regexp.Compile(`<li class="pager-last last"><a title=".+?" href="/users/.+?/videos\?.*?page=([0-9]{1,3})">`)
	urls := reg.FindAllStringSubmatch(string(*body), -1)
	page, _ = strconv.Atoi(urls[0][1]) // 널 포인터 에러가 날것임 => 별다른 처리가 없다면 프로그램이 죽는다
	return
}

func DownloadFile(url *string) (data []byte, err error) {
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
