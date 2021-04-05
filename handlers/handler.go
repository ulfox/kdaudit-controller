package handlers

type Handler interface {
	ObjectHandler(t string, obj interface{}, labelsWatcher map[string]interface{}, slackWebhook string)
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(obj interface{})
}
