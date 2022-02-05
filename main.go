package main

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/thftgr/iwaraDownloader/iwaraApi"
	"github.com/thftgr/iwaraDownloader/pool"
	"github.com/thftgr/iwaraDownloader/src"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

//var rootDownloadPath = `./iwara`

var rootDownloadPath = `Y:/private/iwara/`

func init() {
	if rootDownloadPath[len(rootDownloadPath)-1:] != "/" {
		rootDownloadPath += "/"
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile | log.Lmsgprefix)
	src.ReadDir(rootDownloadPath)
	fmt.Println(len(src.Uploaders), "uploaders.")
	fmt.Println(len(src.FileList), "files.")

	//indent, _ := json.MarshalIndent(src.FileList, "", "  ")
	//fmt.Println(string(indent))
	//fmt.Println("///////////////////////////////////////////////////////////////////////////")
	//indent, _ = json.MarshalIndent(src.FileList, "", "  ")
	//fmt.Println(string(indent))

	//indent, _ := json.MarshalIndent(src.Uploaders, "", "  ")
	//fmt.Println(string(indent))

}

//TODO

func main() {
	//hash := `yxzkpimvl2tobal8j`
	//fmt.Println(iwaraApi.FindUsername(&hash))

	//腿 玩 年
	//syncs("xinhai999")
	syncs("腿 玩 年")
	//for _, uploader := range src.Uploaders {
	//	syncs(uploader)
	//}
}

func syncs(username string) {
	//return
	st := time.Now()
	USERNAME := url.QueryEscape(username)
	hashs := *iwaraApi.GetAllHashByUsername(USERNAME)

	hashSize := len(hashs)

	//hashSize = 1 //테스트용

	jobs := pool.Jobs{}
	for i := 0; i < hashSize; i++ {
		dirName := USERNAME
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
			b, _ := iwaraApi.Fetch(&downloadUrl)
			err := saveLocal(&b, rootDownloadPath+dirName+"/", fileName)
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
