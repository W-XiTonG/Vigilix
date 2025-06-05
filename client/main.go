package main

import (
	"client/interactive"
	"client/util"
	"log"
)

func main() {
	err, ip, port, user, passwd := util.Parameter()
	if err != nil {
		log.Fatalln(err)
	}
	interactive.Connect(ip, user, passwd, port)
}
