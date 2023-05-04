package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

func main() {
	//pquete para linea de comandos
	flagPattern := flag.String("p", "", "filter by pattern")
	flagAll := flag.Bool("a", false, "all file including hide files")
	flaNumberRecords := flag.Int("n", 0, "number of records")

	//bandera por tiempo
	hasOrderByTime := flag.Bool("t", false, "sort by time, oldset first")

	//bandera tamano
	hasOrderBySize := flag.Bool("s", false, "sort buy file size, smallset first")

	//Bandera organizador en reversa
	hasOrderReverse := flag.Bool("r", false, "reverse order while sorting")

	//esto siempre se debe hacer para mape
	flag.Parse()
	path := flag.Arg(0)

	if path == "" {
		path = "."
	}
	dirs, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	fs := []file{}

	for _, dir := range dirs {
		isHidden := isHidden(dir.Name(), path)

		if isHidden && !*flagAll {
			continue
		}

		f, err := getFile(dir, isHidden)
		if err != nil {
			panic(err)
		}

		//aqui validamos si el flag tiene valor o no
		if *flagPattern != "" {
			//aqui validamos si es necesario hacer macth, con la expression regular se aplica que sea insetive
			isMatched, err := regexp.MatchString("(?i)"+*flagPattern, f.name)
			if err != nil {
				panic(err)
			}

			//si es false no llena el slice
			if !isMatched {
				continue
			}
		}
		fs = append(fs, f)
	}

	//aqui validaremos si no nos entrega ningun flag de ordenamiento
	if !*hasOrderBySize || !*hasOrderByTime {
		orderByName(fs, *hasOrderReverse)
	}

	//valida si no le envia la validacion por tamaño
	if *hasOrderBySize && !*hasOrderByTime {
		orderBySize(fs, *hasOrderReverse)
	}

	//validamos orden por tiempo
	if *hasOrderByTime {
		orderByTime(fs, *hasOrderReverse)
	}
	//esto para validar la cantidad de registro que se necesita mostrar
	if *flaNumberRecords == 0 || *flaNumberRecords > len(fs) {
		*flaNumberRecords = len(fs)
	}
	printList(fs, *flaNumberRecords)
}

// esto poblar nuestra ruta
func getFile(dir fs.DirEntry, isHidden bool) (file, error) {
	info, err := dir.Info()
	if err != nil {
		return file{}, fmt.Errorf("dir.Info(): %v", err)
	}
	f := file{
		name:             dir.Name(),
		isDir:            dir.IsDir(),
		userName:         "Renzo",
		groupName:        "Edteam",
		size:             info.Size(),
		modificationTime: info.ModTime(),
		mode:             info.Mode().String(),
	}
	setFile(&f)
	return f, nil
}

// esta funcion es para setear nuestro formato
func setFile(f *file) {

	//hacemos un switch para poder agregar el archivo según como lo encuentre
	switch {
	case isLink(*f):
		f.filetype = fileLink
	case f.isDir:
		f.filetype = fileDirectory
	case isExec(*f):
		f.filetype = fileExecutable
	case isCompress(*f):
		f.filetype = fileCompress
	case isImage(*f):
		f.filetype = fileImage
	default:
		f.filetype = fileRegular
	}

}

// funcion generica
func mySort[T constraints.Ordered](i, j T, isReverse bool) bool {

	if isReverse {
		return i > j
	}
	return i < j
}

// funcion de ordenamiento por nombre
func orderByName(files []file, isReverse bool) {

	sort.SliceStable(files, func(i, j int) bool {

		return mySort(
			strings.ToLower(files[i].name),
			strings.ToLower(files[j].name),
			isReverse,
		)

	})
}

// funcion ordena por tamaño
func orderBySize(files []file, isReverse bool) {
	sort.SliceStable(files, func(i, j int) bool {
		return mySort(
			files[i].size,
			files[j].size,
			isReverse,
		)
	})
}

// function ordenamiento por tiempo
func orderByTime(files []file, isReverse bool) {
	sort.SliceStable(files, func(i, j int) bool {
		return mySort(
			files[i].modificationTime.Unix(),
			files[j].modificationTime.Unix(),
			isReverse,
		)
	})
}

// function para validar si es un archivo link
func isLink(f file) bool {
	return strings.HasPrefix(strings.ToUpper(f.mode), "L")

}

// function para validar si es un archivo ejecutable
func isExec(f file) bool {
	if runtime.GOOS == Windows {
		return strings.HasSuffix(f.name, exe)
	}
	return strings.Contains(f.mode, "x")
}

// function para validar si un archivo comprimido
func isCompress(f file) bool {
	return strings.HasSuffix(f.name, zip) || strings.HasSuffix(f.name, gz) || strings.HasSuffix(f.name, tar) || strings.HasSuffix(f.name, rar) || strings.HasSuffix(f.name, deb)
}

// function para validar si es imagen
func isImage(f file) bool {
	return strings.HasSuffix(f.name, png) || strings.HasSuffix(f.name, jpg) || strings.HasSuffix(f.name, gif)
}

// funcion para pintar de manera visible nuestra cmd
func printList(fs []file, nRecords int) {
	for _, file := range fs[:nRecords] {

		style := mapStyleByFileType[file.filetype]

		fmt.Printf("%s %s %s %10d %s %s %s %s\n", file.mode, file.userName, file.groupName, file.size, file.modificationTime.Format(time.DateTime), style.icon, file.name, style.symbol)
	}
}

// esta función es para validar los archivos ocultos
func isHidden(fileName, basePath string) bool {
	return strings.HasPrefix(fileName, ".")
}
