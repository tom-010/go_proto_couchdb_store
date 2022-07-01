package main

import "log"

func main() {
	p := Person{
		Name: "Tom",
	}
	log.Println(p.String())
}
