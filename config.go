package main

import (
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	AccessToken       string `envconfig:"OAUTH_ACCESS_TOKEN" required:"true"`
	UserToken         string `envconfig:"USER_TOKEN" required:"true"`
	BotAccessToken    string `envconfig:"BOT_ACCESS_TOKEN" required:"true"`
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`
	BotID             string `envconfig:"BOT_ID" required:"true"`
	ChannelID         string `envconfig:"CHANNEL_ID" required:"true"`
	AdminChannelID    string `envconfig:"ADMIN_CHANNEL_ID" required:"false"`
}

// Env ... Variable for environment loading
var Env envConfig

func init() {
	loadEnvironment()
	Env.setDecrypt()
}

func loadEnvironment() {
	if err := godotenv.Load(); err != nil {
		log.Printf("error: Failed to load env file")
	}
	loadEnvconfig()
}

func loadEnvconfig() {
	if err := envconfig.Process("", &Env); err != nil {
		log.Printf("error: Failed to process env var: %s", err)
	}
}

// decryptEnv ... 暗号化されている環境変数を複合する関数
func decryptEnv(v string) string {
	svc := kms.New(session.New(), aws.NewConfig().WithRegion("ap-northeast-1"))
	blob64, err := base64.StdEncoding.DecodeString(v)
	result, err := svc.Decrypt(&kms.DecryptInput{CiphertextBlob: blob64})

	if err != nil {
		log.Printf("error: Got error decrypting data: %s", err)
		return v
	}

	blobString := string(result.Plaintext)

	return blobString
}

// setDecryptEnv ... tokenなどをDecryptして格納し直す関数
func (env *envConfig) setDecrypt() {
	env.AccessToken = decryptEnv(env.AccessToken)
	env.UserToken = decryptEnv(env.UserToken)
	env.BotAccessToken = decryptEnv(env.BotAccessToken)
	env.VerificationToken = decryptEnv(env.VerificationToken)
	env.BotID = decryptEnv(env.BotID)
	env.ChannelID = decryptEnv(env.ChannelID)
	env.AdminChannelID = decryptEnv(env.AdminChannelID)
}
