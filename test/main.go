package main

import "fmt"

func f(values ...int){
	fmt.Print(values[1])
}

func main(){
	f(1,2,3)
}