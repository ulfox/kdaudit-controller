package utils

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

// KDAuditNamespaces stroing kdaudit namespace configuration
type KDAuditNamespaces struct {
	Enabled bool
	Config  map[string]interface{}
}

// KDAuditCFG for storing kdaudit configMap
// The configMap is used to check which kdaudit
// services will be loaded
type KDAuditCFG struct {
	logger       *logrus.Logger
	Service      string
	SlackWebhook string
	Namespaces   KDAuditNamespaces
}

// NewKDAuditCFG for creating a new KDAuditCFG struct
func NewKDAuditCFG(logger *logrus.Logger) KDAuditCFG {
	return KDAuditCFG{
		logger: logger,
	}
}

// ReadConfig for reading kdaudit configMap
func (c *KDAuditCFG) ReadConfig() error {
	service := os.Getenv("KDAUDIT_SERVICE")
	c.Service = service
	c.cfgCheckSlack()
	if c.Service == "namespaceWatcher" {
		return c.cfgCheckNamespaces()
	}
	return nil
}

func (c *KDAuditCFG) cfgCheckSlack() {
	slackWebhook := os.Getenv("SLACK_WEBHOOK")
	if slackWebhook != "" {
		c.SlackWebhook = slackWebhook
	}
}

// cfgCheckNamespaces for enabling the namespace kdaudit service
func (c *KDAuditCFG) cfgCheckNamespaces() error {
	nsEnabled := os.Getenv("NAMESPACE_WATCHER_ENABLED")
	if nsEnabled == "true" {
		c.Namespaces.Enabled = true

		decodedData, err := base64.StdEncoding.DecodeString(os.Getenv("NAMESPACE_WATCHER_CONFIGMAP"))
		if err != nil {
			return err
		}

		jsonMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(decodedData), &jsonMap)
		if err != nil {
			return err
		}

		c.Namespaces.Config = jsonMap
	}
	return nil
}
