package main

import (
	"github.com/mrminglang/go-rigger/boot"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		logrus.Println("System shutdown......")
		boot.Destroy()
	}()

	logrus.Println("System start......")

	boot.Boot()
}
