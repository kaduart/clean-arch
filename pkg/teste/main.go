package main

import "fmt"

func main() {
	evento := []string{"teste", "teste 2", "teste 3", "teste 4", "teste 5"}

	evento = append(evento[:0], evento[1:]...)
	fmt.Println(evento)
}
