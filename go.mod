module github.com/ppc64le-cloud/kubetest2-plugins

go 1.15

require (
	cloud.google.com/go v0.57.0 // indirect
	cloud.google.com/go/storage v1.7.0 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20191115210519-2b2be6cc8ed4 // indirect
	github.com/ChrisTrenkamp/goxpath v0.0.0-20190607011252-c5096ec8773d // indirect
	github.com/IBM-Cloud/bluemix-go v0.0.0-20200724061452-7e266e248891
	github.com/IBM-Cloud/power-go-client v1.0.46
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/antchfx/xpath v1.1.2 // indirect
	github.com/aws/aws-sdk-go v1.32.3 // indirect
	github.com/bmatcuk/doublestar v1.2.1 // indirect
	github.com/golang/protobuf v1.4.1 // indirect
	github.com/hashicorp/go-hclog v0.13.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-plugin v1.2.2
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/hashicorp/hcl/v2 v2.5.0 // indirect
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93 // indirect
	github.com/hashicorp/logutils v1.0.0
	github.com/hashicorp/terraform v0.12.21
	github.com/hashicorp/terraform-plugin-sdk v1.14.0
	github.com/hashicorp/terraform-svchost v0.0.0-20191119180714-d2e4933b9136 // indirect
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/masterzen/simplexml v0.0.0-20190410153822-31eea3082786 // indirect
	github.com/masterzen/winrm v0.0.0-20190308153735-1d17eaf15943 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mitchellh/cli v1.1.1
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.3.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/openshift/installer v0.16.1
	github.com/packer-community/winrmcp v0.0.0-20180921211025-c76d91c1e7db // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.10.1
	github.com/posener/complete v1.2.3 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/spf13/pflag v1.0.5
	//github.com/terraform-providers/terraform-provider-aws v0.0.0
	github.com/ulikunitz/xz v0.5.7 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/xlab/treeprint v0.0.0-20181112141820-a009c3971eca // indirect
	github.com/zclconf/go-cty v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79
	golang.org/x/net v0.0.0-20200506145744-7e3656a0809f // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.0.0-20200616195046-dc31b401abb5 // indirect
	google.golang.org/api v0.25.0 // indirect
	google.golang.org/genproto v0.0.0-20200507105951-43844f6eee31 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	//k8s.io/api v0.18.3 // indirect
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
	sigs.k8s.io/kubetest2 v0.0.0-20200815004617-8027cdc89d56
)

replace (
	github.com/Azure/go-autorest => github.com/tombuildsstuff/go-autorest v14.0.1-0.20200416184303-d4e299a3c04a+incompatible
	github.com/Azure/go-autorest/autorest => github.com/tombuildsstuff/go-autorest/autorest v0.10.1-0.20200416184303-d4e299a3c04a
	github.com/Azure/go-autorest/autorest/azure/auth => github.com/tombuildsstuff/go-autorest/autorest/azure/auth v0.4.3-0.20200416184303-d4e299a3c04a
	//github.com/go-log/log => github.com/go-log/log v0.1.1-0.20181211034820-a514cf01a3eb // Pinned by MCO
	github.com/hashicorp/terraform => github.com/openshift/terraform v0.12.20-openshift-4 // Pin to fork with deduplicated rpc types v0.12.20-openshift-4
	github.com/hashicorp/terraform-plugin-sdk => github.com/openshift/hashicorp-terraform-plugin-sdk v1.14.0-openshift // Pin to fork with public rpc types
	//github.com/terraform-providers/terraform-provider-aws => github.com/openshift/terraform-provider-aws v1.60.1-0.20200630224953-76d1fb4e5699 // Pin to openshift fork with tag v2.67.0-openshift
	//github.com/terraform-providers/terraform-provider-azurerm => github.com/openshift/terraform-provider-azurerm v1.40.1-0.20200707062554-97ea089cc12a // release-2.17.0 branch
	//github.com/terraform-providers/terraform-provider-ignition/v2 => github.com/community-terraform-providers/terraform-provider-ignition/v2 v2.1.0
	//github.com/terraform-providers/terraform-provider-vsphere => github.com/openshift/terraform-provider-vsphere v1.18.1-openshift-1
	//github.com/vmware/govmomi => github.com/vmware/govmomi v0.22.2-0.20200420222347-5fceac570f29
	//k8s.io/api => k8s.io/api v0.17.1 // Replaced by MCO/CRI-O
	//k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.1 // Replaced by MCO/CRI-O
	//k8s.io/apiserver => k8s.io/apiserver v0.17.1 // Replaced by MCO/CRI-O
	//k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/client-go => k8s.io/client-go v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/code-generator => k8s.io/code-generator v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/component-base => k8s.io/component-base v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/cri-api => k8s.io/cri-api v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/kubectl => k8s.io/kubectl v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/kubelet => k8s.io/kubelet v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/kubernetes => k8s.io/kubernetes v1.17.1 // Replaced by MCO/CRI-O
//k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/metrics => k8s.io/metrics v0.17.1 // Replaced by MCO/CRI-O
//k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.1 // Replaced by MCO/CRI-O
)
