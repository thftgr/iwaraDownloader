package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/pterm/pterm"
	"github.com/thftgr/iwaraDownloader/config"
	"github.com/thftgr/iwaraDownloader/iwaraApi"
	"github.com/thftgr/iwaraDownloader/src"
	"io"
	"log"
	"net/http"
	"os"
)

//TODO 용어 정리
//TODO 	OgVZVfZ5jys65Yk91 = 파일명
//TODO 	sz拓海             = 유저네임
//TODO
//TODO

//TODO 디렉토리 구조  BASE/{username}/{username}_{filename}.mp4
//TODO 입력 가능성 있는 데이터 => 파일명 유저네임
//TODO

//TODO 1. 무작위 파일명 선택
//TODO 2. 유저네임 추춯
//TODO 3. 1번항목의 파일명과 유저네임이 일치하는지 확인
//TODO 4. 3번 false인 경우 기존 유저네임으로 파일명들 추출
//TODO 5. 기존 파일명을 신규 파일명으로 일괄 업데이트 후 캐싱 다시 실행

var urlOnly bool

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile | log.Lmsgprefix)

	//	indent, _ := json.MarshalIndent(src.FileIndex, "", "  ")
	//	fmt.Println(string(indent))
	//	fmt.Println("///////////////////////////////////////////////////////////////////////////")
	//indent, _ = json.MarshalIndent(src.FileByUploader, "", "  ")
	//fmt.Println(string(indent))
	//indent, _ = json.MarshalIndent(src.FileList, "", "  ")
	//fmt.Println(string(indent))

	//indent, _ := json.MarshalIndent(src.Uploaders, "", "  ")
	//fmt.Println(string(indent))

}
func main() {
	//rename Y:/private/iwara/呆音/MonakaWagasino_ljjzvtxvpruwjnbao.mp4 to Y:/private/iwara/呆音/呆音_ljjzvtxvpruwjnbao.mp4
	//indent, _ := json.MarshalIndent(src.FileIndex, "", "  ")
	//fmt.Println(string(indent))
	//fmt.Println("///////////////////////////////////////////////////////////////////////////")
	endedFilename := map[string]bool{}
	ii := 0
	for filename, usernmae := range src.FileIndex.Filename {
		ii++
		if endedFilename[filename] {
			//log.Printf("pass %s_%s.mp4", usernmae, filename)
			continue
		}
		fmt.Println("\n\nvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
		fmt.Println(ii, "/", len(src.FileIndex.Filename))
		newUsername := iwaraApi.FindUsername(&filename)
		filenames := *iwaraApi.GetAllFilenameByUsername(newUsername)
		log.Println(filename, newUsername)
		if newUsername == "" {
			continue
		}

		if err := os.MkdirAll(config.RoorDir+newUsername, 775); err != nil {
			log.Println(err)
		}

		for _, filename := range filenames {
			if src.FileIndex.DirName[filename] == "" {
				if urlOnly {
					log.Println("X\t", filename)
					err := saveUrl(filename)
					if err != nil {
						log.Println(err)
						continue
					}
				} else {
					log.Println(pterm.Red("X start Download ", filename))
					err := download(newUsername, filename)
					if err != nil {
						log.Println(err)
						continue
					}
					log.Println(pterm.Green("O finish Download ", filename))
				}

				continue
			} else {
				log.Println(pterm.Green("O\t", filename))
			}
			oldFilename := fmt.Sprintf("%s%s/%s_%s.mp4", config.RoorDir, src.FileIndex.DirName[filename], usernmae, filename)
			newFilename := fmt.Sprintf("%s%s/%s_%s.mp4", config.RoorDir, newUsername, newUsername, filename)
			endedFilename[filename] = true

			if oldFilename == newFilename {
				continue
			}
			log.Printf("rename %s to %s | %v\n", oldFilename, newFilename, os.Rename(oldFilename, newFilename))
		}
		de, _ := os.ReadDir(config.RoorDir + src.FileIndex.DirName[filename])
		if len(de) < 1 {
			_ = os.RemoveAll(config.RoorDir + src.FileIndex.DirName[filename])
		}
	}
}
func saveUrl(filename string) (err error) {
	f, err := os.OpenFile("iwaraUrls.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	_, err = f.WriteString("https://ecchi.iwara.tv/videos/" + filename + "\r\n")
	return

}
func download(usernmae, filename string) (err error) {
	downloadUrl, err := iwaraApi.GetDownloadUrl(filename)
	if err != nil {
		log.Println(err)
	}
	res, err := http.Get(downloadUrl)
	if err != nil {
		log.Println(err)
		return

	}
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}
	defer res.Body.Close()

	bar := pb.Full.Start64(res.ContentLength)

	// create proxy reader
	barReader := bar.NewProxyReader(res.Body)

	// copy from proxy reader
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, barReader)

	// finish bar
	bar.Finish()
	dir := fmt.Sprintf("%s%s", config.RoorDir, usernmae)
	newFilename := fmt.Sprintf("%s_%s.mp4", usernmae, filename)
	return saveLocal(&buf, dir, newFilename)

}

//
//TODO
//526ad5a9d46b4330e1d105232106f948fef75f49   OgVZVfZ5jys65Yk91
//func main() {
//	//hash := `yxzkpimvl2tobal8j`
//	//fmt.Println(iwaraApi.FindUsername(&hash))
//
//	//腿 玩 年
//	//syncs("xinhai999")
//	//syncs("腿 玩 年")
//	//for _, uploader := range src.Uploaders {
//	//	syncs(uploader)
//	//}
//
//	us := []string{
//		"3dhgames",
//		"sz拓海",
//		"toyhentai",
//		"cramoisi",
//		"vtol-neko",
//		"ilixiya",
//		"xraymmd",
//		"ay",
//		"orion-0",
//		"akeginu",
//		"niziiro-ageha",
//		"%E6%9D%A5%E4%B8%80%E5%8F%91%E5%92%95%E5%99%9C%E7%81%B5%E6%B3%A2",
//		"%E6%B7%B1%E7%94%B0%E3%83%A1%E3%82%A4",
//		"%E5%AE%87%E8%BD%A9%E5%91%80",
//		"smixix",
//		"cdream",
//	}
//	for _, uploader := range us {
//		syncs(uploader)
//	}
//
//}
//
//func syncs(username string) {
//	//return
//	st := time.Now()
//	USERNAME := url.QueryEscape(username)
//	hashs := *iwaraApi.GetAllFilenameByUsername(USERNAME)
//
//	hashSize := len(hashs)
//
//	//hashSize = 1 //테스트용
//
//	jobs := pool.Jobs{}
//	for i := 0; i < hashSize; i++ {
//		dirName := USERNAME
//		i := i
//
//		if src.FileList[strings.ToUpper(hashs[i])].File != nil {
//			fmt.Println(pterm.Green("O\t", hashs[i]))
//			continue
//		} else {
//			fmt.Println(pterm.Red("X\t", hashs[i]))
//		}
//
//		jobs = append(jobs, func() interface{} {
//
//			downloadUrl, _ := iwaraApi.GetDownloadUrl(hashs[i])
//			fileName := fmt.Sprintf("%s_%s.mp4", dirName, hashs[i])
//			fmt.Println(downloadUrl)
//			fmt.Println("filename: ", fileName)
//			fmt.Println("==========================================")
//			fmt.Println("started download.")
//			fmt.Println("path:", rootDownloadPath+dirName+"/"+fileName)
//			fmt.Println("filename:", fileName)
//			fmt.Println("==========================================")
//			b, _ := iwaraApi.Fetch(&downloadUrl)
//			err := saveLocal(&b, rootDownloadPath+dirName+"/", fileName)
//			fmt.Println("==========================================")
//			fmt.Println("download Finished.")
//			fmt.Println("path:", rootDownloadPath+dirName+"/"+fileName)
//			fmt.Println("filename:", fileName)
//			fmt.Println("==========================================")
//			return err
//		})
//	}
//	if len(jobs) > 1 {
//		pool.StartPool(jobs, 4)
//	}
//
//	et := time.Now()
//	fmt.Println("Total Time:", et.UnixMilli()-st.UnixMilli(), "ms")
//}

func saveLocal(data *bytes.Buffer, dir, name string) (err error) {
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

	_, _ = file.Write(data.Bytes())
	file.Close()

	if _, err = os.Stat(fullPath); !os.IsNotExist(err) {
		_ = os.Remove(fullPath)
	}
	_ = os.Rename(fullPath+".idownload", fullPath)

	return
}
