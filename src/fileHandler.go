package src

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
)

type FileData struct {
	ServerKey string
	Uploader  string
	FullPath  string
	File      fs.FileInfo
}

var (
	//FileList    = map[string]fs.FileInfo{}
	FileList    = map[string]FileData{}
	Uploaders   = []string{}
	mutex       = sync.Mutex{}
	regFilename *regexp.Regexp
)

func init() {
	regFilename, _ = regexp.Compile(`(.+)_([A-Za-z0-9]{15,20})[.]mp4`)
}

func ReadDir(path string) {
	defer func() {
		err, _ := recover().(error)
		if err != nil {
			fmt.Println("", err)
		}
	}()
	if path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}

	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		fileName := file.Name()

		if file.IsDir() {
			if path[len(path)-1:] != "/" {
				path += "/"
			}
			fmt.Println("find Dir:", path+fileName)
			ReadDir(path + fileName)
			Uploaders = append(Uploaders, fileName)
		} else {
			if !regFilename.Match([]byte(fileName)) {
				continue
			}
			r := regFilename.FindAllStringSubmatch(fileName, -1)[0]
			uploader := r[1]
			serverKey := r[2]
			mutex.Lock()

			FileList[strings.ToUpper(serverKey)] = FileData{
				ServerKey: serverKey,
				Uploader:  uploader,
				FullPath:  path,
				File:      file,
			}
			mutex.Unlock()

		}
	}
	Uploaders = removeDuplicate(Uploaders)

}
func removeDuplicate(sa []string) (sr []string) {

	allKeys := make(map[string]bool)

	for _, item := range sa {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			sr = append(sr, item)
		}
	}
	return
}
