package e2e

import (
	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	test "github.com/openshift/cluster-kube-apiserver-operator/test/library"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/cert"
	"strings"
	"testing"
)

func TestUserClientCABundle(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	kubeConfig, err := test.NewClientConfigForTest()
	require.NoError(t, err)
	kubeClient, err := clientcorev1.NewForConfig(kubeConfig)
	require.NoError(t, err)
	configClient, err := configclient.NewForConfig(kubeConfig)
	require.NoError(t, err)
	clientCA := test.NewCertificateAuthorityCertificate(t, nil)
	configMapName := strings.ToLower(test.GenerateNameForTest(t, "UserCA"))
	_, err = kubeClient.ConfigMaps("openshift-config").Create(&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: configMapName}, Data: map[string]string{"ca-bundle.crt": string(cert.EncodeCertPEM(clientCA.Certificate))}})
	require.NoError(t, err)
	defer func() {
		_, err := updateAPIServerClusterConfigSpec(configClient, func(apiServer *configv1.APIServer) {
			apiServer.Spec.ClientCA.Name = ""
		})
		assert.NoError(t, err)
	}()
	_, err = updateAPIServerClusterConfigSpec(configClient, func(apiServer *configv1.APIServer) {
		apiServer.Spec.ClientCA.Name = configMapName
	})
	require.NoError(t, err)
	var lastResourceVersion string
	err = wait.Poll(test.WaitPollInterval, test.WaitPollTimeout, func() (bool, error) {
		caBundle, err := kubeClient.ConfigMaps("openshift-kube-apiserver").Get("client-ca", metav1.GetOptions{})
		if err != nil || caBundle.ResourceVersion == lastResourceVersion {
			return false, nil
		}
		certificates, err := cert.ParseCertsPEM([]byte(caBundle.Data["ca-bundle.crt"]))
		if err != nil {
			return false, err
		}
		for _, certificate := range certificates {
			if certificate.SerialNumber.String() == clientCA.Certificate.SerialNumber.String() {
				return true, nil
			}
		}
		return false, nil
	})
	require.NoError(t, err, "user client-ca not found in combined client-ca bundle")
}
