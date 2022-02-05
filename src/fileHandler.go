package src

import (
	"github.com/thftgr/iwaraDownloader/config"
	"os"
	"regexp"
)

type fileIndex struct {
	Username map[string][]string // filename array
	Filename map[string]string   // username
	DirName  map[string]string   // username
}

var (
	FileIndex = fileIndex{
		Username: map[string][]string{},
		Filename: map[string]string{},
		DirName:  map[string]string{},
	}
	regFilename, _ = regexp.Compile(`(.*)_([A-Za-z0-9]{15,20})[.]mp4`)
)

func init() {
	ReadAllFiles()
}

func ReadAllFiles() {
	root := config.RoorDir
	dirs, _ := os.ReadDir(root) // ROOT/{usernmae}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		files, _ := os.ReadDir(root + dir.Name()) // ROOT/{usernmae}/{usernmae}_{filename}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if !regFilename.Match([]byte(file.Name())) {
				continue
			}
			r := regFilename.FindAllStringSubmatch(file.Name(), -1)[0]
			FileIndex.Filename[r[2]] = r[1]
			FileIndex.DirName[r[2]] = dir.Name()
			FileIndex.Username[r[1]] = append(FileIndex.Username[r[1]], r[2])
		}
	}
}
