package main

import "time"

//file type
const (
	fileRegular int = iota
	fileDirectory
	fileExecutable
	fileCompress
	fileImage
	fileLink
)

//file extension
const (
	exe = ".exe"
	deb = ".deb"
	zip = ".zip"
	gz  = ".gz"
	tar = ".tar"
	rar = ".rar"
	png = ".png"
	jpg = ".jpg"
	gif = ".gif"
)

//estructura
type file struct {
	name             string
	filetype         int
	isDir            bool
	isHidden         bool
	userName         string
	groupName        string
	size             int64
	modificationTime time.Time
	mode             string
}

//estrutura para reconcer los colores icons y simbolo
type styleFileType struct {
	icon   string
	color  string
	symbol string
}

//nuestro mapa junto con la estrucutra para que nos muestre en el cmd
var mapStyleByFileType = map[int]styleFileType{
	fileRegular:    {icon: "📄"},
	fileDirectory:  {icon: "📂", color: "BLUE", symbol: "/"},
	fileExecutable: {icon: "🚀", color: "GREEN", symbol: "*"},
	fileCompress:   {icon: "📦", color: "RED"},
	fileImage:      {icon: "📷", color: "MAGENTA"},
	fileLink:       {icon: "🔗", color: "CYAN"},
}
