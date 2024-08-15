package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	db "github.com/arya2004/xyfin/db/sqlc"
	"github.com/arya2004/xyfin/mail"
	"github.com/arya2004/xyfin/utils"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)


const TaskSendVerifyEmail = "task:send_verify_email"


type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}


func (distributor *RedisTaskDistributor) SendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error{

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cannot marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("cannot enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil

}

func (processor *RedisTaskProcessor)ProcessVerifyEmail(ctx context.Context, task *asynq.Task) error{
	
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("cannot unmarshal payload: %w", err)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows{
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("cannot get user: %w", err)
	}

	// send email to user.Email
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username: user.Username,
		Email: user.Email,
		SecretCode: utils.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("cannot create verify email: %w", err)
	}

	config, err := utils.LoadConfig("..")

	if err != nil {
		fmt.Println("Failed to load config:", err)
		
	}

	sender := mail.NewGmailSender("Xphyrus", config.EmailSenderAddress, config.EmailSenderPassword)
	to := []string{user.Email}
	
	data := map[string]string{
		"Code": verifyEmail.SecretCode,
	}

	err = sender.SendEmail("Welcome to Our Service!", "./welcome_template.html", to, nil, nil, nil, data)
	if err != nil {
		fmt.Println("Failed to send email:", err)
	}
	
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")
	return nil
}