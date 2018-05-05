package main

import "fmt"

func main() {
    maxGoroutines := 10
	guard := make(chan bool, maxGoroutines)
	for j:= 0;j<10;j++{
		guard <- true
	}

    for i := 0; i < 30; i++ {   
        if (<-guard) {
           go worker(i)
        }else{
			guard <- false
		}
	}
	var a int
	fmt.Scan(a)
	fmt.Println(a)
}

func worker(i int) { 
	fmt.Println("doing work on", i) 
}