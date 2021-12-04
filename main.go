package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const videoHashRegex = `<a href="/videos/(.+?)(?:[?].+?|["])>`

func main() {

	st := time.Now()
	url := "https://ecchi.iwara.tv/users/%E8%BF%99%E8%85%BF%E5%80%9Fwo%E7%8E%A9%E4%B8%80%E5%A4%A9?language=ja"
	//url := "https://ecchi.iwara.tv/users/%E4%B8%89%E4%BB%81%E6%9C%88%E9%A5%BC"
	err := getBaseUrl(&url)
	if err != nil {
		fmt.Println("cannot parse iwara user url")
		panic(err)
	}
	fmt.Println(url)

	res, err := http.Get(url)
	if err != nil {
		logErr(err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("=========================================")
		fmt.Println(res.StatusCode, res.Status)
		fmt.Println("client error: ", err)
		fmt.Println("=========================================")
	}
	page := getMaxPage(&body) + 1
	{
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
						fmt.Println(err.Error())
					}
				}()

				res, _ := http.Get(url + `?page=` + strconv.Itoa(i))
				if res.StatusCode != http.StatusOK {
					fmt.Println("=========================================")
					fmt.Println(res.StatusCode, res.Status)
					fmt.Println("=========================================")
				}
				body, _ := ioutil.ReadAll(res.Body)
				reg, err := regexp.Compile(videoHashRegex)
				if err != nil {
					logErr(err)
					return
				}
				urls := reg.FindAllStringSubmatch(string(body), -1)
				mutex.Lock()
				defer mutex.Unlock()
				hashs = append(hashs, getSubMatchData(urls, 1)...)
			}()
		}
		for i := 0; i < page; i++ {
			<-ch
		}
		fmt.Println("=========================================")
		fmt.Println(strings.Join(hashs, "\n"))
		fmt.Println(res.StatusCode, res.Status)
		fmt.Println(fmt.Sprintf("find %d keys from %d page", len(hashs), page))
		fmt.Println("=========================================")
		//_, _ = getDownloadUrl(hashs)
	}

	et := time.Now()
	fmt.Println("Total Time:", et.UnixMilli()-st.UnixMilli(), "ms")

}

func getBaseUrl(s *string) (err error) {
	defer func() {
		err, _ = recover().(error)
	}()
	reg, _ := regexp.Compile(`(https://ecchi.iwara.tv/users/.+?)(?:(/videos)|[?]|/|$)`)
	*s = reg.FindAllStringSubmatch(*s, -1)[0][1] + "/videos"
	return
}

func logErr(err error) {
	fmt.Println("=========================================")
	fmt.Println("client error: ", err)
	fmt.Println("=========================================")
}

func getSubMatchData(sa [][]string, index int) (sr []string) {

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

const downloadBaseUrl = `https://ecchi.iwara.tv/api/video/`

func getDownloadUrl(hashs []string) (urls []string, err error) {
	for i := 0; i < len(hashs); i++ {
		res, e := http.Get(downloadBaseUrl + hashs[i])
		if e != nil {
			fmt.Println(e.Error())
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("=========================================")
			fmt.Println(res.StatusCode, res.Status)
			fmt.Println("client error: ", err)
			fmt.Println("=========================================")
		}
		var ress []downloadUrlStruct

		if err := json.Unmarshal(body, &ress); err != nil {
			fmt.Println("=========================================")
			fmt.Println("body parse error: ", err)
			fmt.Println("=========================================")
			continue
		}
		for i := 0; i < len(ress); i++ {
			if ress[i].Resolution == `Source` {
				fmt.Println("successful get url")
				fmt.Println(ress[i].Uri)
				urls = append(urls, `https:`+ress[i].Uri)
			}
		}
	}
	return
}

func getMaxPage(body *[]byte) (page int) {
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
