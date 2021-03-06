package certrotationcontroller

import (
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
	"net"
)

const workQueueKey = "key"

func (c *CertRotationController) syncServiceHostnames() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hostnames := sets.NewString("kubernetes", "kubernetes.default", "kubernetes.default.svc")
	hostnames.Insert("kubernetes.default.svc." + "cluster.local")
	networkConfig, err := c.networkLister.Get("cluster")
	if err != nil {
		return err
	}
	for _, cidrString := range networkConfig.Status.ServiceNetwork {
		_, serviceCIDR, err := net.ParseCIDR(cidrString)
		if err != nil {
			return err
		}
		ip, err := cidr.Host(serviceCIDR, 1)
		if err != nil {
			return err
		}
		hostnames.Insert(ip.String())
	}
	klog.V(2).Infof("syncing servicenetwork hostnames: %v", hostnames.List())
	c.serviceNetwork.setHostnames(hostnames.List())
	return nil
}
func (c *CertRotationController) runServiceHostnames() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for c.processServiceHostnames() {
	}
}
func (c *CertRotationController) processServiceHostnames() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dsKey, quit := c.serviceHostnamesQueue.Get()
	if quit {
		return false
	}
	defer c.serviceHostnamesQueue.Done(dsKey)
	err := c.syncServiceHostnames()
	if err == nil {
		c.serviceHostnamesQueue.Forget(dsKey)
		return true
	}
	utilruntime.HandleError(fmt.Errorf("%v failed with : %v", dsKey, err))
	c.serviceHostnamesQueue.AddRateLimited(dsKey)
	return true
}
func (c *CertRotationController) serviceHostnameEventHandler() cache.ResourceEventHandler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.ResourceEventHandlerFuncs{AddFunc: func(obj interface{}) {
		c.serviceHostnamesQueue.Add(workQueueKey)
	}, UpdateFunc: func(old, new interface{}) {
		c.serviceHostnamesQueue.Add(workQueueKey)
	}, DeleteFunc: func(obj interface{}) {
		c.serviceHostnamesQueue.Add(workQueueKey)
	}}
}
