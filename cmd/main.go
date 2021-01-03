package main

import (
	"github.com/GrafeasGroup/tor-tools-go/shadowban"
)

func main() {
	done := make(chan bool)
	go shadowban.Checker(done)
	<-done
}
