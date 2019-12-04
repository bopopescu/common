package mail

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SendMail(t *testing.T) {
	mc := MailConf{
		User:     "xxxx@163.com",
		Passwd:   "xxxx",
		Server:   "smtp.163.com",
		RecvList: []string{"xxx@163.com"},
		SendUser: "xxx",
	}

	err := SendMail("Test", "This is a test mail", &mc)
	if !assert.Nil(t, err) {
		log.Printf("收到意外错误: %s", err.Error())
	}
}
