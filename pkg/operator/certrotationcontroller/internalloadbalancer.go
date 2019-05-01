package certrotationcontroller

import (
	"fmt"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
	"strings"
)

func (c *CertRotationController) syncInternalLoadBalancerHostnames() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	infrastructureConfig, err := c.infrastructureLister.Get("cluster")
	if err != nil {
		return err
	}
	hostname := infrastructureConfig.Status.APIServerURL
	hostname = strings.Replace(hostname, "https://", "", 1)
	hostname = hostname[0:strings.LastIndex(hostname, ":")]
	klog.V(2).Infof("syncing internal loadbalancer hostnames: %v", hostname)
	c.internalLoadBalancer.setHostnames([]string{hostname})
	return nil
}
func (c *CertRotationController) runInternalLoadBalancerHostnames() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for c.processExternalLoadBalancerHostnames() {
	}
}
func (c *CertRotationController) processInternalLoadBalancerHostnames() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dsKey, quit := c.internalLoadBalancerHostnamesQueue.Get()
	if quit {
		return false
	}
	defer c.internalLoadBalancerHostnamesQueue.Done(dsKey)
	err := c.syncInternalLoadBalancerHostnames()
	if err == nil {
		c.internalLoadBalancerHostnamesQueue.Forget(dsKey)
		return true
	}
	utilruntime.HandleError(fmt.Errorf("%v failed with : %v", dsKey, err))
	c.internalLoadBalancerHostnamesQueue.AddRateLimited(dsKey)
	return true
}
func (c *CertRotationController) internalLoadBalancerHostnameEventHandler() cache.ResourceEventHandler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.ResourceEventHandlerFuncs{AddFunc: func(obj interface{}) {
		c.internalLoadBalancerHostnamesQueue.Add(workQueueKey)
	}, UpdateFunc: func(old, new interface{}) {
		c.internalLoadBalancerHostnamesQueue.Add(workQueueKey)
	}, DeleteFunc: func(obj interface{}) {
		c.internalLoadBalancerHostnamesQueue.Add(workQueueKey)
	}}
}
