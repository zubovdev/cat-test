package reminder

import (
	"cat-test/internal/domain"
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
	"time"
)

type taskReminder struct {
	usecase domain.TaskUsecase
	ctx     context.Context
	dialer  *gomail.Dialer
}

func (r taskReminder) notifyUsers() {
	tasks, err := r.usecase.GetOverdueTasks(context.Background(), time.Now().Unix())
	if err != nil {
		return
	}

	for _, task := range tasks {
		msg := gomail.NewMessage()
		msg.SetHeader("To", task.User.Email)
		msg.SetHeader("From", "mullydeveloper@gmail.com")
		msg.SetHeader("Subject", "Notification")
		msg.SetBody("text/plain", fmt.Sprintf(`Your task "%s" is overdued!`, task.Title))
		_ = r.dialer.DialAndSend(msg)

		task.MailSent = true
		_ = r.usecase.Update(context.Background(), task)
	}
}

func (r taskReminder) run() {
	ticker := time.NewTicker(time.Second * 3)

	go func() {
		for {
			select {
			case <-ticker.C:
				r.notifyUsers()
			}
		}
	}()

	<-r.ctx.Done()
	ticker.Stop()
}

func RunTaskReminder(ctx context.Context, usecase domain.TaskUsecase, dialer *gomail.Dialer) {
	reminder := taskReminder{usecase: usecase, ctx: ctx, dialer: dialer}
	reminder.run()
}
