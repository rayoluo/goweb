package email

import (
	"crypto/tls"
	"ginblog/utils"
	"github.com/jordan-wright/email"
	"log"
	"mime"
	"net/smtp"
	"strings"
)

/**
 * @Author: rayoluo
 * @Author: hustly123@gmail.com
 * @Date: 2021/6/2 23:24
 * @Desc: 发送邮件功能模块
 */

// func main() {
// 	sendEmail("测试第三方 email 库", "xuxinhua@studygolang.com")
// }

func SendEmail(subject string, content string, tos ...string) error {
	e := email.NewEmail()

	smtpUsername := utils.MailAddress
	e.From = mime.QEncoding.Encode("UTF-8", "雨落的博客") + "<" + smtpUsername + ">"
	e.To = tos
	e.Subject = subject
	e.Text = []byte(content)

	auth := smtp.PlainAuth("", smtpUsername, utils.MailPass, "smtp.qq.com")
	// err := e.Send("smtp.qq.com:465", auth)
	err := e.SendWithTLS("smtp.qq.com:465", auth, &tls.Config{ServerName: "smtp.qq.com"})
	if err != nil {
		log.Println("Send Mail to", strings.Join(tos, ","), "error:", err)
		return err
	}
	log.Println("Send Mail to", strings.Join(tos, ","), "Successfully")
	return nil
}
