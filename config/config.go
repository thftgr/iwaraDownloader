package config

import "strings"

//var RoorDir = `./iwara`

var RoorDir = `Y:/private/iwara`

func init() {
	RoorDir = strings.TrimRight(RoorDir, "/") + "/"
}
