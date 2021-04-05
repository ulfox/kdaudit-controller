package main

import (
	"github.com/sirupsen/logrus"
	"github.com/ulfox/kdaudit/controllers"
	"github.com/ulfox/kdaudit/utils"
)

var (
	logger *logrus.Logger

	// SRVER is set on build time
	SRVER string = "<not set>"
)

func main() {
	sigs := utils.NewOSSignal()

	logger = logrus.New()
	logger.WithField("KDAudit", SRVER).Info("Initiated")

	// sig := utils.NewSignal()
	// go sig.Wait()
	// sig.Stop()

	kdAuditFlags := utils.ParseFlags()

	client := utils.NewClientAuth(
		kdAuditFlags.ClientAuthType,
		kdAuditFlags.LocalKubeCFG,
		logger,
	)

	kdAuditCFG := utils.NewKDAuditCFG(logger)
	err := kdAuditCFG.ReadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	if kdAuditCFG.Namespaces.Enabled {
		nsController := controllers.NewNSController(
			client,
			kdAuditCFG.Namespaces.Config,
			kdAuditCFG.SlackWebhook,
			logger,
		)

		informerTerm := make(chan struct{})
		defer close(informerTerm)

		go nsController.Run(informerTerm)
	}

	sigs.Wait()
}
