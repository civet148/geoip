package client

import (
	"github.com/civet148/log"
	"testing"
)

func TestGetIpMsg(t *testing.T) {
	msg, err := GetIpMsg(Language_EN, "113.110.225.29")
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Infof("%+v", msg)
}
