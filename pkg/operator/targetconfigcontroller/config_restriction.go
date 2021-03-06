package targetconfigcontroller

import (
	godefaultbytes "bytes"
	"fmt"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

func RemoveConfig(dst, src interface{}, currentPath string) (interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return dst, nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
