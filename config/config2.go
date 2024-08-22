package config

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
	"sync"
)

type AppConfiguration struct {
	ENV              string
	DBDebug          bool
	ApiPrefix        string
	ApiKey           string
	SuffixForTracing string
	Version          string
}

type ServerConfiguration struct {
	Port int
}

type DatabaseConfiguration struct {
	MasterPayment db.Config
	SlavePayment  db.Config
	FlipDBMaster  db.Config
	FlipDBSlave   db.Config
}

// Configuration config
type Configuration struct {
	App               AppConfiguration
	Server            ServerConfiguration
	Database          DatabaseConfiguration
	Redis             RedisConfiguration
	Sentry            SentryConfiguration
	ExternalAPI       ExternalAPIConfiguration
	ApiKey            ApiKeyConfiguration
	GCloud            GCloudConfig
	Consumers         ConsumerConfigs
	Publishers        PublisherConfigs
	Cron              CronConfiguration
	CommService       CommServiceConfiguration
	SquadcastService  SquadcastServiceConfiguration
	HelpCenterService HelpCenterServiceConfiguration
}

type RedisConfiguration struct {
	PayRedis redis.Config
}

type ExternalAPIConfiguration struct {
	Porta CommonExternalAPI
}

type CommonExternalAPI struct {
	BaseURL string
	ApiKey  string
}
type SentryConfiguration struct {
	DSN *string
}

type ApiKeyConfiguration struct {
	FlipServer string
	General    string
}

type GCloudConfig struct {
	ProjectID string
}

type ConsumerConfigs struct {
	Enable              bool
	DecisionSession     ConsumerConfig
	DecisionTransaction ConsumerConfig
}

type PublisherConfigs struct {
	DecisionSession     PublisherConfig
	DecisionTransaction PublisherConfig
}

type ConsumerConfig struct {
	Topic                  string
	Subscription           string
	MaxOutstandingMessages int
	NumGoroutines          int
	Toggle                 bool
}

type PublisherConfig struct {
	Topic  string
	Toggle bool
}

type CronConfiguration struct {
	BlacklistDevice CronJob
}

type CommServiceConfiguration struct {
	Username   string
	Password   string
	SenderMail string
	CSEmail    string
	URL        string
	ProxyURL   string
}

type CronJob struct {
	Toggle   bool
	Interval string
	Limit    int
}

type SquadcastServiceConfiguration struct {
	RefreshToken        string
	PageId              int64
	StateId             map[string]int64
	StatusId            map[string]int64
	ComponentId         map[string]map[string]int64
	GetAccessTokenUrl   string
	CreateIssueUrl      string
	ResolveIssueUrl     string
	CreateIssueMessage  string
	ResolveIssueMessage string
}

type HelpCenterServiceConfiguration struct {
	BaseURL string
	Token   string
}

var (
	configuration *Configuration
	once          sync.Once
)

// All get all config
func All(opts ...Option) *Configuration {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AllowEmptyEnv(true)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		if err := viper.Unmarshal(&configuration); err != nil {
			log.Fatalf("Unable to decode into struct, %v", err)
		}

		for _, opt := range opts {
			opt(configuration)
		}

		if configuration.App.ENV != constants.LocalEnvString {
			ctx := context.Background()
			client, err := secretmanager.NewClient(ctx)
			if err != nil {
				log.Fatalf("failed for setup client secret manager: %v", err)
			}
			defer client.Close()

			loadAndSetConfigFromSecretManager(ctx, client)
		}
	})

	return configuration
}

func loadAndSetConfigFromSecretManager(ctx context.Context, client *secretmanager.Client) error {
	configuration.Database.FlipDBMaster.User = loadFromSecretManager(ctx, client, FlipDBUser)
	configuration.Database.FlipDBMaster.Password = loadFromSecretManager(ctx, client, FlipDBPassword)
	configuration.Database.FlipDBSlave.User = loadFromSecretManager(ctx, client, FlipDBUser)
	configuration.Database.FlipDBSlave.Password = loadFromSecretManager(ctx, client, FlipDBPassword)
	configuration.Database.MasterPayment.User = loadFromSecretManager(ctx, client, PayDBUser)
	configuration.Database.MasterPayment.Password = loadFromSecretManager(ctx, client, PayDBPassword)
	configuration.Database.SlavePayment.User = loadFromSecretManager(ctx, client, PayDBUser)
	configuration.Database.SlavePayment.Password = loadFromSecretManager(ctx, client, PayDBPassword)
	configuration.ExternalAPI.Porta.ApiKey = loadFromSecretManager(ctx, client, PortaAPIKey)
	configuration.CommService.Username = loadFromSecretManager(ctx, client, CommServiceUsername)
	configuration.CommService.Password = loadFromSecretManager(ctx, client, CommServicePassword)
	configuration.SquadcastService.RefreshToken = loadFromSecretManager(ctx, client, SquadcastServiceRefreshToken)
	configuration.HelpCenterService.Token = loadFromSecretManager(ctx, client, HelpCenterToken)
	return nil
}

func loadFromSecretManager(ctx context.Context, client *secretmanager.Client, keyName string) string {
	name := fmt.Sprintf("projects/%v/secrets/%v/versions/latest", configuration.GCloud.ProjectID, keyName)

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		panic(fmt.Errorf("failed to access secret version: %v", err))
	}

	return string(result.Payload.Data)
}

func Get() *Configuration {
	return configuration
}

func GetSuffixForTracing() string {
	if configuration == nil {
		return ""
	}

	return configuration.App.SuffixForTracing
}

func GetEnv() string {
	return configuration.App.ENV
}

// Comm Service Related Config
func GetCommServiceAuth() (string, string) {
	return configuration.CommService.Username, configuration.CommService.Password
}

func GetCommSenderMail() string {
	return configuration.CommService.SenderMail
}

func GetCSEmail() string {
	return configuration.CommService.CSEmail
}

func GetCommServiceUrl() string {
	return configuration.CommService.URL
}

func GetCommServiceProxyUrl() string {
	return configuration.CommService.ProxyURL
}

func GetSquadcastServiceRefreshToken() string {
	return configuration.SquadcastService.RefreshToken
}

func GetSquadcastServicePageId() int64 {
	return configuration.SquadcastService.PageId
}

func GetSquadcastServiceStateId() map[string]int64 {
	return configuration.SquadcastService.StateId
}

func GetSquadcastServiceStatusId() map[string]int64 {
	return configuration.SquadcastService.StatusId
}

func GetSquadcastServiceComponentId() map[string]map[string]int64 {
	return configuration.SquadcastService.ComponentId
}

func GetSquadcastServiceGetAccessTokenUrl() string {
	return configuration.SquadcastService.GetAccessTokenUrl
}

func GetSquadcastServiceCreateIssueUrl() string {
	return configuration.SquadcastService.CreateIssueUrl
}

func GetSquadcastServiceResolveIssueUrl() string {
	return configuration.SquadcastService.ResolveIssueUrl
}

func GetSquadcastServiceCreateIssueMessage() string {
	return configuration.SquadcastService.CreateIssueMessage
}

func GetSquadcastServiceResolveIssueMessage() string {
	return configuration.SquadcastService.ResolveIssueMessage
}

// HC Service Related Config
func GetHelpCenterServiceBaseUrl() string {
	return configuration.HelpCenterService.BaseURL
}

func GetHelpCenterServiceToken() string {
	return configuration.HelpCenterService.Token
}

func GetPortaApiKey() string {
	return configuration.ExternalAPI.Porta.ApiKey
}

// for testing purpose
func Set(conf *Configuration) {
	configuration = conf
}
