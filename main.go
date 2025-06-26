package main

import "github.com/joho/godotenv"

func main() {
	println("Hello, World!")

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}
