package informers

import (
	"context"

	log "github.com/sirupsen/logrus"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func NewNSInformer(client kubernetes.Interface) (cache.SharedIndexInformer, []*workqueue.RateLimitingInterface) {
	return multiQueueGenericInformer(
		client,
		cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().Namespaces().List(context.TODO(), meta_v1.ListOptions{})
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Namespaces().Watch(context.TODO(), meta_v1.ListOptions{})
			},
		},
		api_v1.Namespace{},
	)
}

func multiQueueGenericInformer(
	client kubernetes.Interface,
	i cache.ListWatch,
	ns api_v1.Namespace,
) (
	cache.SharedIndexInformer,
	[]*workqueue.RateLimitingInterface,
) {

	nsInformer := cache.NewSharedIndexInformer(&i, &ns, 0, cache.Indexers{})

	queueAdd := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	queueUpdate := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	queueDelete := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	nsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			log.Infof("multiQueueGenericInformer.AddFunc: %s", key)
			if err == nil {
				queueAdd.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			log.Infof("multiQueueGenericInformer.UpdateFunc: %s", key)
			if err == nil {
				queueUpdate.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			log.Infof("multiQueueGenericInformer.DeleteFunc: %s", key)
			if err == nil {
				queueDelete.Add(key)
			}
		},
	})

	queues := []*workqueue.RateLimitingInterface{&queueAdd, &queueUpdate, &queueDelete}
	return nsInformer, queues
}
