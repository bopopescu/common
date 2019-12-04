package mail

import (
	"crypto/tls"
	"net/smtp"
	"strings"
)

type MailConf struct {
	User     string   // 邮箱账号
	Passwd   string   // 邮箱密码
	Server   string   // 服务器 x.x.x.x:465
	RecvList []string // 接受人列表
	SendUser string   // 发件人
}

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func send_mail_ssl(msg []byte, mailConf *MailConf) error {
	//到邮件服务器验证账户
	auth := unencryptedAuth{
		smtp.PlainAuth(
			"",
			mailConf.User,
			mailConf.Passwd,
			mailConf.Server,
		),
	}

	//SSL/TLS配置，服务器返回时验证主机名
	tlsCfg := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         mailConf.Server,
	}

	//连接，运行在465端口上的SSL/TLS连接需要
	conn, err := tls.Dial("tcp", mailConf.Server+":465", tlsCfg)
	if nil != err {
		return err
	}

	//获取邮件客户端
	c, err := smtp.NewClient(conn, mailConf.Server)
	if nil != err {
		return err
	}
	defer c.Quit()

	//鉴权
	if err := c.Auth(auth); nil != err {
		return err
	}

	//发送者
	if err := c.Mail(mailConf.User); nil != err {
		return err
	}

	//设置邮件接受者
	for _, sender := range mailConf.RecvList {
		if err := c.Rcpt(sender); nil != err {
			return err
		}
	}

	//告诉服务器，准备发送邮件头部和数据
	w, err := c.Data()
	if nil != err {
		return err
	}

	//发送
	if _, err := w.Write([]byte(msg)); nil != err {
		return err
	}

	//关闭写
	if err := w.Close(); nil != err {
		return err
	}

	//关闭连接

	return nil
}

// 新接口，建议使用
func SendMail(subject string, data string, mailConf *MailConf) error {
	content_type := "Content-Type: text/html; charset=UTF-8\r\n"
	msg := []byte(content_type + "From: " + mailConf.SendUser + "<" + mailConf.User + ">\r\n" +
		"To: " + strings.Join(mailConf.RecvList, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + data)

	return send_mail_ssl(msg, mailConf)
}
