package handlers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/ulfox/kdaudit/alerts"
	core_v1 "k8s.io/api/core/v1"
)

var (
	slackCLient alerts.KDAuditSlackAlert
)

type NSHandler struct {
	labelsWatcher map[string]interface{}
}

func getLabels(obj interface{}) map[string]string {
	return obj.(*core_v1.Namespace).Labels
}

func checkLabels(obj interface{}, labelsWatcher map[string]interface{}) {
	labels := getLabels(obj)
	for v, k := range labelsWatcher {
		val, ok := labels[v]
		if !ok {
			msg := fmt.Sprintf("NSHandler.checkLabels: %s is missing %v label",
				obj.(*core_v1.Namespace).Name, v)
			log.Error(msg)
			err := slackCLient.SendSlackNotification(msg)
			if err != nil {
				log.Error(err)
			}
		} else {
			if val != k {
				msg := fmt.Sprintf("NSHandler.checkLabels: %s has wrong %s label value",
					obj.(*core_v1.Namespace).Name, v)
				log.Errorf(msg)
				err := slackCLient.SendSlackNotification(msg)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}

func (t *NSHandler) ObjectHandler(s string, obj interface{}, labelsWatcher map[string]interface{}, slackWebhook string) {
	t.labelsWatcher = labelsWatcher
	if slackWebhook != "" {
		slackCLient = alerts.NewSlackAlert(slackWebhook)
	}

	if s == "add" {
		t.ObjectCreated(obj)
	} else if s == "update" {
		t.ObjectUpdated(obj)
	} else if s == "delete" {
		t.ObjectDeleted(obj)
	}
}

func (t *NSHandler) ObjectCreated(obj interface{}) {
	checkLabels(obj, t.labelsWatcher)
}

func (t *NSHandler) ObjectDeleted(obj interface{}) {
	log.Debug("NSHandler.NamespaceDeleted")
}

func (t *NSHandler) ObjectUpdated(obj interface{}) {
	checkLabels(obj, t.labelsWatcher)
}
