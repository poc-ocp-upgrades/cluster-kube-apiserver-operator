package certrotationcontroller

import (
	"sync"
	"k8s.io/apimachinery/pkg/util/sets"
)

type DynamicServingRotation struct {
	lock			sync.RWMutex
	hostnames		[]string
	hostnamesChanged	chan struct{}
}

func (r *DynamicServingRotation) setHostnames(newHostnames []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.isSame(newHostnames) {
		return
	}
	r.lock.Lock()
	r.hostnames = newHostnames
	r.lock.Unlock()
	r.hostnamesChanged <- struct{}{}
}
func (r *DynamicServingRotation) isSame(newHostnames []string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.lock.RLock()
	defer r.lock.RUnlock()
	existingSet := sets.NewString(r.hostnames...)
	newSet := sets.NewString(newHostnames...)
	return existingSet.Equal(newSet)
}
func (r *DynamicServingRotation) GetHostnames() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.hostnames
}
