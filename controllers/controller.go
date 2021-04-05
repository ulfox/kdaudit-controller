package controllers

type Controller interface {
	Run(informerTerm <-chan struct{})
	HasSynced() bool
	runWorker()
}
