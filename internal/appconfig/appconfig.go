package appconfig

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/care-giver-app/care-giver-api/internal/log"
	"go.uber.org/zap"
)

const (
	LocalEnv = "local"
)

type AppConfig struct {
	Env               string
	AWSConfig         aws.Config
	Logger            *zap.Logger
	UserTableName     string
	ReceiverTableName string
}

func NewAppConfig() *AppConfig {
	appCfg := &AppConfig{}
	appCfg.ReadEnvVars()
	appCfg.Logger = log.GetLoggerWithEnv(log.InfoLevel, appCfg.Env)
	return appCfg
}

func (a *AppConfig) ReadEnvVars() {
	a.Env = getEnvVarStringOrDefault("ENV", LocalEnv)
	a.UserTableName = getEnvVarStringOrDefault("USER_TABLE_NAME", fmt.Sprintf("%s-%s", "user-table", LocalEnv))
	a.ReceiverTableName = getEnvVarStringOrDefault("RECEIVER_TABLE_NAME", fmt.Sprintf("%s-%s", "receiver-table", LocalEnv))
}

func getEnvVarStringOrDefault(envVar string, defaultValue string) string {
	env, present := os.LookupEnv(envVar)
	if present {
		return env
	}
	return defaultValue
}
