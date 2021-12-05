package src

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"regexp"
	"sync"
)

var (
	FileList    = map[string]map[string]fs.FileInfo{}
	mutex       = sync.Mutex{}
	regFilename *regexp.Regexp
)

func init() {
	regFilename, _ = regexp.Compile(`(.+)_([A-Za-z0-9]{15,20})[.]mp4`)
}

func ReadDir(path string) {

	defer func() { // 함수 빠져나가기 직전 무조건 실행된다
		err, _ := recover().(error) // 프로그램이 죽는경우 살린다
		if err != nil {
			fmt.Println(err)
		}
	}()
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			fmt.Println("find Dir:", path+"/"+file.Name())
			ReadDir(path + "/" + file.Name())
		} else {
			if regFilename.Match([]byte(file.Name())) {
				names := regFilename.FindAllStringSubmatch(file.Name(), -1)
				mutex.Lock()
				if FileList[names[0][1]] == nil {
					FileList[names[0][1]] = map[string]fs.FileInfo{}
					FileList[names[0][1]][names[0][2]] = file
				} else {
					FileList[names[0][1]][names[0][2]] = file
				}
				mutex.Unlock()
			}
		}
	}
}
