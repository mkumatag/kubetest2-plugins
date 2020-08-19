terraform {
  required_version = ">= 0.12"
}

variable "cluster_name" {
  description = "K8s cluster name"
}

variable "release_marker" {
  description = "Kubernetes release marker"
  default = "ci/latest"
}

variable "build_version" {
  description = "Kubernetes Build Number"
}
