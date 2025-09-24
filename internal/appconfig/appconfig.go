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
	Env                   string
	AWSConfig             aws.Config
	Logger                *zap.Logger
	UserTableName         string
	ReceiverTableName     string
	EventTableName        string
	RelationshipTableName string
}

func NewAppConfig() *AppConfig {
	appCfg := &AppConfig{}
	appCfg.ReadEnvVars()
	logger, err := log.GetLoggerWithEnv(log.InfoLevel, appCfg.Env)
	if err != nil {
		panic(err)
	}
	appCfg.Logger = logger
	return appCfg
}

func (a *AppConfig) ReadEnvVars() {
	a.Env = getEnvVarStringOrDefault("ENV", LocalEnv)
	a.UserTableName = getEnvVarStringOrDefault("USER_TABLE_NAME", fmt.Sprintf("%s-%s", "user-table", LocalEnv))
	a.ReceiverTableName = getEnvVarStringOrDefault("RECEIVER_TABLE_NAME", fmt.Sprintf("%s-%s", "receiver-table", LocalEnv))
	a.EventTableName = getEnvVarStringOrDefault("EVENT_TABLE_NAME", fmt.Sprintf("%s-%s", "event-table", LocalEnv))
	a.RelationshipTableName = getEnvVarStringOrDefault("RELATIONSHIP_TABLE_NAME", fmt.Sprintf("%s-%s", "relationship-table", LocalEnv))
}

func getEnvVarStringOrDefault(envVar string, defaultValue string) string {
	env, present := os.LookupEnv(envVar)
	if present {
		return env
	}
	return defaultValue
}
