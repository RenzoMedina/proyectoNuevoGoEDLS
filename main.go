package main

import (
	"flag"
	"fmt"
)

func main() {
	//pquete para linea de comandos
	flagPattern := flag.String("p", "", "filter by pattern")
	flagAll := flag.Bool("a", false, "all file includin hide files")
	flaNumberRecords := flag.Int("n", 0, "number of records")

	//bandera por tiempo
	hasOrderByTime := flag.Bool("t", false, "sort by time, oldset first")

	//bandera tamano
	hasOrderBySize := flag.Bool("s", false, "sort buy file size, smallset first")

	//Bandera organizador
	hasOrderReverse := flag.Bool("r", false, "reverse order while sorting")
	//esto siempre se debe hacer para mape
	flag.Parse()
	fmt.Println("pattern:", *flagPattern)
	fmt.Println("all:", *flagAll)
	fmt.Println("records:", *flaNumberRecords)
	fmt.Println("time:", *hasOrderByTime)
	fmt.Println("size:", *hasOrderBySize)
	fmt.Println("reverse:", *hasOrderReverse)
}
