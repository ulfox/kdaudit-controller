package controllers

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ulfox/kdaudit/handlers"
	"github.com/ulfox/kdaudit/informers"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type NSController struct {
	logger           *logrus.Logger
	client           kubernetes.Interface
	queueAdd         *workqueue.RateLimitingInterface
	queueUpdate      *workqueue.RateLimitingInterface
	queueDelete      *workqueue.RateLimitingInterface
	informer         cache.SharedIndexInformer
	handler          handlers.Handler
	cfgLabelsWatcher map[string]interface{}
	slackWebhook     string
	exit             bool
}

func NewNSController(client *kubernetes.Clientset, kda map[string]interface{}, slackWebhook string, logger *logrus.Logger) *NSController {
	informer, queues := informers.NewNSInformer(client)

	return &NSController{
		logger:           logger,
		client:           client,
		informer:         informer,
		queueAdd:         queues[0],
		queueUpdate:      queues[1],
		queueDelete:      queues[2],
		handler:          &handlers.NSHandler{},
		cfgLabelsWatcher: kda,
		slackWebhook:     slackWebhook,
		exit:             false,
	}
}

func (c *NSController) Run(informerTerm <-chan struct{}) {
	defer runtime.HandleCrash()

	defer (*c.queueAdd).ShutDown()
	defer (*c.queueUpdate).ShutDown()
	defer (*c.queueDelete).ShutDown()

	go c.informer.Run(informerTerm)

	if !cache.WaitForCacheSync(informerTerm, c.HasSynced) {
		runtime.HandleError(errors.New("NSController: Could not sync cache"))
		return
	}
	wait.Until(c.runWorker, time.Second, informerTerm)
}

func (c *NSController) HasSynced() bool {
	return c.informer.HasSynced()
}

func (c *NSController) runWorker() {
	for !c.exit {
		if (*c.queueUpdate).Len() > 0 {
			c.processQueue("update", c.queueUpdate)
		} else if (*c.queueAdd).Len() > 0 {
			c.processQueue("add", c.queueAdd)
		} else if (*c.queueDelete).Len() > 0 {
			c.processQueue("delete", c.queueDelete)
		}
	}
}

func (c *NSController) processQueue(
	t string,
	q *workqueue.RateLimitingInterface,
) {

	key, quit := (*q).Get()
	if quit {
		return
	}

	defer (*q).Done(key)

	item, _, err := c.informer.GetIndexer().GetByKey(key.(string))
	if err != nil && item != nil {
		if (*q).NumRequeues(key) < 5 {
			c.logger.Errorf("NSController.processQueue: Failed getting key %s with error %v, retrying", key, err)
			(*q).AddRateLimited(key)
		} else {
			c.logger.Errorf("NSController.processQueue: Timedount trying to get %s", key)
			(*q).Forget(key)
			runtime.HandleError(err)
		}
	} else if item == nil {
		c.logger.Debugf("NSController.processQueue: Got empty item using key %s. Skipping", key)
		(*q).Forget(key)
		return
	}
	c.handler.ObjectHandler(t, item, c.cfgLabelsWatcher, c.slackWebhook)
	(*q).Forget(key)
}
