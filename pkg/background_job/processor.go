package background_job

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Giafn/Depublic/configs"
	"github.com/gocraft/work"
	"gopkg.in/gomail.v2"
)

func Init(redisHost, redisPort string) {

	var redisPool = Pool(redisHost, redisPort)

	pool := work.NewWorkerPool(Context{}, 10, "app_depublic", redisPool)

	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).CheckParam)

	pool.Job("send_email", (*Context).SendEmail)

	pool.Start()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	pool.Stop()
}

type Context struct {
	email   string
	Subject string
	Body    string
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting Job: ", job.Name)
	return next()
}

func (c *Context) CheckParam(job *work.Job, next work.NextMiddlewareFunc) error {
	if _, ok := job.Args["email_address"]; !ok {
		c.email = job.ArgString("email_address")
		c.Subject = job.ArgString("subject")
		c.Body = job.ArgString("body")
		if err := job.ArgError(); err != nil {
			return fmt.Errorf("arg error %v", err.Error())
		}
	}
	return next()
}

func (c *Context) SendEmail(job *work.Job) error {
	return SendEmail(
		job.ArgString("email_address"),
		job.ArgString("subject"),
		job.ArgString("body"),
	)
}

func SendEmail(to, subject, body string) error {
	cfg, err := configs.NewConfig()
	checkError(err)

	from := cfg.SMTP.From
	password := cfg.SMTP.Pass
	smtpHost := cfg.SMTP.Host
	smtpPort := 587

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(smtpHost, smtpPort, cfg.SMTP.User, password)

	err = dialer.DialAndSend(mailer)
	checkError(err)

	fmt.Println("Email sent successfully to:", to)
	return nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
