package main

import "fmt"

// Env ... Variable for environment loading
var Env envConfig

func init() {
	loadEnvironment()
}

func main() {
	fmt.Printf("%s", "Hello world")
}
