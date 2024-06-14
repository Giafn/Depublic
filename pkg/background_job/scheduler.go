package background_job

import (
	"fmt"
	"log"

	"github.com/Giafn/Depublic/configs"
	"github.com/gocraft/work"
)

func ScheduleEmails(
	emailAddress string,
	subject string,
	body string,
) {
	cfg, err := configs.NewConfig()
	checkError(err)
	fmt.Println("app_depublic")

	var redisPool = Pool(cfg.Redis.Host, cfg.Redis.Port)

	var enqueuer = work.NewEnqueuer("app_depublic", redisPool)

	_, err = enqueuer.Enqueue("send_email", work.Q{"email_address": emailAddress, "subject": subject, "body": body})
	if err != nil {
		log.Fatal(err)
	}
}
