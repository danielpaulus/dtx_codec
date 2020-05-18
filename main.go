package main

import (
	"fmt"
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

	payloadDumpFile, err2 := os.Create("payload_dump.json")
	if err2 != nil {
		log.Fatal("couldnt create file")
	}
	defer payloadDumpFile.Close()

	remaining := 1
	payloadDumpFile.Write([]byte("["))
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

		if msg.HasPayload() {
			payloadDumpFile.Write([]byte(fmt.Sprintf("\"%x\",", msg.GetPayloadBytes())))
		}

	}
	payloadDumpFile.Write([]byte("]"))
}
