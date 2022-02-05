package iwaraApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

const downloadBaseUrl = `https://ecchi.iwara.tv/api/video/`

func GetAllHashByUsername(username string) (hashes *[]string) {
	var (
		URL string
	)

	URL = "https://ecchi.iwara.tv/users/" + username
	err := GetBaseUrl(&URL)
	if err != nil {
		fmt.Println("cannot parse iwara user URL")
		return
	}
	fmt.Println(URL)

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

	ch := make(chan struct{}, page)
	var hashs []string
	mutex := sync.Mutex{}
	for i := 0; i < page; i++ {
		i := i
		go func() {
			defer func() {
				ch <- struct{}{}
				err, _ = recover().(error)
				if err != nil {
					log.Println(err)
				}
			}()

			res, _ := http.Get(URL + `?page=` + strconv.Itoa(i))
			if res.StatusCode != http.StatusOK {
				fmt.Println("=========================================")
				fmt.Println(res.StatusCode, res.Status)
				fmt.Println("=========================================")
			}
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			reg, _ := regexp.Compile(`<a href="/videos/(.+?)(?:[?].+?|["])>`)
			urls := reg.FindAllStringSubmatch(string(body), -1)
			mutex.Lock()
			defer mutex.Unlock()
			hashs = append(hashs, GetSubMatchData(urls, 1)...)
		}()
	}
	for i := 0; i < page; i++ {
		<-ch
	}
	fmt.Println("=========================================")
	//fmt.Println(strings.Join(hashs, "\n"))
	fmt.Println(res.StatusCode, res.Status)
	fmt.Println(fmt.Sprintf("find %d keys from %d page", len(hashs), page))
	fmt.Println("=========================================")
	return &hashs

}

func FindUsername(hash *string) (username string) {
	defer func() {
		_, _ = recover().(error)
	}()
	//class="username">清炽</a>
	url := `https://ecchi.iwara.tv/videos/` + *hash
	r, _ := Fetch(&url)
	reg, _ := regexp.Compile(`class="username">(.+?)</a>`)
	username = reg.FindAllStringSubmatch(string(r), -1)[0][1]
	return
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

func GetDownloadUrl(hashs string) (urls string, err error) {
	defer func() { // 함수 빠져나가기 직전 무조건 실행된다
		err, _ = recover().(error) // 프로그램이 죽는경우 살린다
		if err != nil {            // 죽이고 살린 후 처리
			log.Println(err)
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
