package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/danielpaulus/dtx_codec/dtx"
)

func main() {
	dat, err := ioutil.ReadFile("dtx/fixtures/conn-in6.dump")
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
	var fragmentDecoder *dtx.FragmentDecoder = nil
	for remaining != 0 {
		msg, remainingBytes, err := dtx.Decode(dat)
		if err != nil {
			panic(err)
		}
		if fragmentDecoder != nil {
			fragmentDecoder.AddFragment(msg)
			if fragmentDecoder.HasFinished() {
				bytes := fragmentDecoder.Extract()
				msg1, _, err1 := dtx.Decode(bytes)
				if err1 != nil {
					panic(err1)
				}
				print(msg1.StringDebug())
				f.Write([]byte(msg1.String()))
				f.Write([]byte("\n"))
				f.Write([]byte(msg1.StringDebug()))
				f.Write([]byte("\n\n"))
				fragmentDecoder = nil
			}
		}
		if msg.IsFirstFragment() {
			fragmentDecoder = dtx.NewFragmentDecoder(msg)
		}

		remaining = len(remainingBytes)
		dat = remainingBytes
		f.Write([]byte(msg.String()))
		f.Write([]byte("\n"))
		f.Write([]byte(msg.StringDebug()))
		f.Write([]byte("\n\n"))

	}
	payloadDumpFile.Write([]byte("]"))
}
