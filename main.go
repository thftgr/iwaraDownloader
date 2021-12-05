package main

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/thftgr/iwaraDownloader/iwaraApi"
	"github.com/thftgr/iwaraDownloader/pool"
	"github.com/thftgr/iwaraDownloader/src"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var rootDownloadPath = `./iwara`

//const rootDownloadPath = `Y:/private/iwara/`

func init() {
	if rootDownloadPath[len(rootDownloadPath)-1:] != "/" {
		rootDownloadPath += "/"
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile | log.Lmsgprefix)
	src.ReadDir(rootDownloadPath)

	//indent, _ := json.MarshalIndent(src.FileList, "", "  ")
	//fmt.Println(string(indent))
	//fmt.Println("///////////////////////////////////////////////////////////////////////////")
	//indent, _ = json.MarshalIndent(src.FileList, "", "  ")
	//fmt.Println(string(indent))

	//indent, _ := json.MarshalIndent(src.Uploaders, "", "  ")
	//fmt.Println(string(indent))

}

func main() {
	for _, uploader := range src.Uploaders {
		syncs(uploader)
	}
}

func syncs(username string) {
	//return
	var (
		USERNAME = url.QueryEscape(username)
		URL      string
	)

	st := time.Now()

	URL = "https://ecchi.iwara.tv/users/" + USERNAME
	err := iwaraApi.GetBaseUrl(&URL)
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
	page := iwaraApi.GetMaxPage(&body) + 1

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
			hashs = append(hashs, iwaraApi.GetSubMatchData(urls, 1)...)
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

	hashSize := len(hashs)

	//hashSize = 1 //테스트용

	jobs := pool.Jobs{}
	for i := 0; i < hashSize; i++ {
		dirName, _ := url.QueryUnescape(USERNAME)
		i := i

		if src.FileList[strings.ToUpper(hashs[i])].File != nil {
			fmt.Println(pterm.Green("O\t", hashs[i]))
			continue
		} else {
			fmt.Println(pterm.Red("X\t", hashs[i]))
		}

		jobs = append(jobs, func() interface{} {

			downloadUrl, _ := iwaraApi.GetDownloadUrl(hashs[i])
			fileName := fmt.Sprintf("%s_%s.mp4", dirName, hashs[i])
			fmt.Println(downloadUrl)
			fmt.Println("filename: ", fileName)
			fmt.Println("==========================================")
			fmt.Println("started download.")
			fmt.Println("path:", rootDownloadPath+dirName+"/"+fileName)
			fmt.Println("filename:", fileName)
			fmt.Println("==========================================")
			b, _ := iwaraApi.DownloadFile(&downloadUrl)
			err = saveLocal(&b, rootDownloadPath+dirName+"/", fileName)
			fmt.Println("==========================================")
			fmt.Println("download Finished.")
			fmt.Println("path:", rootDownloadPath+dirName+"/"+fileName)
			fmt.Println("filename:", fileName)
			fmt.Println("==========================================")
			return err
		})
	}
	if len(jobs) > 1 {
		pool.StartPool(jobs, 4)
	}

	et := time.Now()
	fmt.Println("Total Time:", et.UnixMilli()-st.UnixMilli(), "ms")
}
func saveLocal(data *[]byte, dir, name string) (err error) {
	defer func() {
		err, _ = recover().(error)
	}()
	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}
	fullPath := dir + name
	_ = os.MkdirAll(dir, 775)
	file, _ := os.Create(fullPath + ".idownload")
	if file == nil {
		return
	}
	_, _ = file.Write(*data)
	file.Close()

	if _, err = os.Stat(fullPath); !os.IsNotExist(err) {
		_ = os.Remove(fullPath)
	}
	_ = os.Rename(fullPath+".idownload", fullPath)

	return
}
