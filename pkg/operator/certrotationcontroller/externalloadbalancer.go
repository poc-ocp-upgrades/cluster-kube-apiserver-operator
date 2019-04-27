package certrotationcontroller

import (
	"fmt"
	"strings"
	"k8s.io/klog"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
)

func (c *CertRotationController) syncExternalLoadBalancerHostnames() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	infrastructureConfig, err := c.infrastructureLister.Get("cluster")
	if err != nil {
		return err
	}
	hostname := infrastructureConfig.Status.APIServerURL
	hostname = strings.Replace(hostname, "https://", "", 1)
	hostname = hostname[0:strings.LastIndex(hostname, ":")]
	hostname = strings.Replace(hostname, "api-int.", "api.", 1)
	klog.V(2).Infof("syncing external loadbalancer hostnames: %v", hostname)
	c.externalLoadBalancer.setHostnames([]string{hostname})
	return nil
}
func (c *CertRotationController) runExternalLoadBalancerHostnames() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for c.processExternalLoadBalancerHostnames() {
	}
}
func (c *CertRotationController) processExternalLoadBalancerHostnames() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dsKey, quit := c.externalLoadBalancerHostnamesQueue.Get()
	if quit {
		return false
	}
	defer c.externalLoadBalancerHostnamesQueue.Done(dsKey)
	err := c.syncExternalLoadBalancerHostnames()
	if err == nil {
		c.externalLoadBalancerHostnamesQueue.Forget(dsKey)
		return true
	}
	utilruntime.HandleError(fmt.Errorf("%v failed with : %v", dsKey, err))
	c.externalLoadBalancerHostnamesQueue.AddRateLimited(dsKey)
	return true
}
func (c *CertRotationController) externalLoadBalancerHostnameEventHandler() cache.ResourceEventHandler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.ResourceEventHandlerFuncs{AddFunc: func(obj interface{}) {
		c.externalLoadBalancerHostnamesQueue.Add(workQueueKey)
	}, UpdateFunc: func(old, new interface{}) {
		c.externalLoadBalancerHostnamesQueue.Add(workQueueKey)
	}, DeleteFunc: func(obj interface{}) {
		c.externalLoadBalancerHostnamesQueue.Add(workQueueKey)
	}}
}
