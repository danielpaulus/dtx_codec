package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/danielpaulus/dtx_codec/dtx"
)

func main() {
	dat, err := ioutil.ReadFile("dtx/fixtures/conn-out6.dump")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("dump.txt")
	if err != nil {
		log.Fatal("couldnt create file")
	}
	defer f.Close()

	remaining := 1
	for remaining != 0 {
		msg, remainingBytes, err := dtx.Decode(dat)
		if err != nil {
			log.Fatal(err)
		}
		remaining = len(remainingBytes)
		dat = remainingBytes
		f.Write([]byte(msg.String()))
		f.Write([]byte("\n"))
		f.Write([]byte(msg.StringDebug()))
		f.Write([]byte("\n"))

	}
}
