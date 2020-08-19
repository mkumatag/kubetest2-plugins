package deployer

var (
	user_data_template = `#cloud-config
# This is a cloud-config file for bringing the k8s master node up and running

write_files:
  - content: |
      {
        "exec-opts": ["native.cgroupdriver=systemd"],
        "log-driver": "json-file",
        "log-opts": {
          "max-size": "100m"
        },
        "storage-driver": "overlay2",
        "storage-opts": [
          "overlay2.override_kernel_check=true"
        ],
        "mtu": 9000
      }
    path: /etc/docker/daemon.json
  - content: |
      [Unit]
      Description=kubelet: The Kubernetes Node Agent
      Documentation=https://kubernetes.io/docs/
      Wants=network-online.target
      After=network-online.target

      [Service]
      Environment="KUBELET_CONFIG_ARGS=--config=/var/lib/kubelet/config.yaml"
      Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
      Environment="KUBELET_KUBEADM_ARGS=--cgroup-driver=systemd --network-plugin=cni --pod-infra-container-image=k8s.gcr.io/pause:3.2"
      ExecStart=/usr/local/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_CONFIG_ARGS $KUBELET_KUBEADM_ARGS $KUBELET_EXTRA_ARGS
      Restart=always
      StartLimitInterval=0
      RestartSec=10

      [Install]
      WantedBy=multi-user.target
    path: /usr/lib/systemd/system/kubelet.service
  - content: |
      #!/usr/bin/env bash
      # https://storage.googleapis.com/kubernetes-release-dev/ci/latest.txt for build numbers
      CI_VERSION=${CI_VERSION:-"v1.20.0-alpha.0.515+382107e6c84374"}
      ARCH=${ARCH:-"ppc64le"}
      RELEASE_MARKER=${RELEASE_MARKER:-"ci/latest"}

      echo "Downloading kubectl..."
      curl -sSL https://dl.k8s.io/ci/${CI_VERSION}/bin/linux/${ARCH}/kubectl > /usr/local/bin/kubectl
      echo "Downloading kubelet..."
      curl -sSL https://dl.k8s.io/ci/${CI_VERSION}/bin/linux/${ARCH}/kubelet > /usr/local/bin/kubelet
      echo "Downloading kubeadm..."
      curl -sSL https://dl.k8s.io/ci/${CI_VERSION}/bin/linux/${ARCH}/kubeadm > /usr/local/bin/kubeadm

      chmod +x /usr/local/bin/kube*
      DOCKER_TEMP=$(mktemp -d)
      wget https://oplab9.parqtec.unicamp.br/pub/ppc64el/docker/version-19.03.8/centos/docker-ce-19.03.8-3.el7.ppc64le.rpm -P ${DOCKER_TEMP}
      wget https://oplab9.parqtec.unicamp.br/pub/ppc64el/docker/version-19.03.8/centos/docker-ce-cli-19.03.8-3.el7.ppc64le.rpm -P ${DOCKER_TEMP}
      wget https://dl.fedoraproject.org/pub/epel/7/ppc64le/Packages/c/containerd-1.2.4-1.el7.ppc64le.rpm -P ${DOCKER_TEMP}

      yum install -y conntrack-tools socat tc
      yum install -y ${DOCKER_TEMP}/*.rpm
      yum remove -y podman
      systemctl enable docker
      systemctl start docker

      iptables -P INPUT ACCEPT
      iptables -P FORWARD ACCEPT
      iptables -P OUTPUT ACCEPT

      # Flush All Iptables Chains/Firewall rules #
      iptables -F

      # Delete all Iptables Chains #
      iptables -X

      # Flush all counters too #
      iptables -Z
      # Flush and delete all nat and  mangle #
      iptables -t nat -F
      iptables -t nat -X
      iptables -t mangle -F
      iptables -t mangle -X
      iptables -t raw -F
      iptables -t raw -X

      systemctl restart docker
      swapoff -a

      systemctl daemon-reload
      systemctl enable kubelet
      systemctl start kubelet
      kubeadm reset -f
      rm -rf /etc/cni/net.d
      kubeadm init --apiserver-bind-port=%d --apiserver-cert-extra-sans %s --pod-network-cidr=172.20.0.0/16 --kubernetes-version ${RELEASE_MARKER}

      mkdir -p $HOME/.kube
      sudo cp /etc/kubernetes/admin.conf $HOME/.kube/config
      sudo chown $(id -u):$(id -g) $HOME/.kube/config
      curl https://docs.projectcalico.org/manifests/calico.yaml -O
      sed -i 's/veth_mtu\:.*/veth_mtu: \"8940\"/' calico.yaml

      kubectl create -f calico.yaml
    path: /usr/local/bin/bootstrap-k8s.sh

runcmd:
  - [sh, /usr/local/bin/bootstrap-k8s.sh]`
)
