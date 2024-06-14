package service

import "github.com/Giafn/Depublic/pkg/background_job"

func ScheduleEmails(emailAddresses string, subject, body string) {
	background_job.ScheduleEmails(
		emailAddresses,
		subject,
		body,
	)
}
