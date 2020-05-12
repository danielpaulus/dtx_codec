package dtx_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/dtx_codec/dtx"

	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/notifyOfPublishedCapabilites")
	if err != nil {
		log.Fatal(err)
	}
	msg, err := dtx.Decode(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, msg.Fragments, uint16(1))
		assert.Equal(t, msg.FragmentIndex, uint16(0))
		assert.Equal(t, msg.MessageLength, 612)
		assert.Equal(t, dtx.MethodinvocationWithoutExpectedReply, msg.MessageType)
		assert.Equal(t, false, msg.ExpectsReply)
		print(msg.String())
	}

}
