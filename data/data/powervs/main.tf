data "ibm_resource_group" "rg" {
  name = var.powervs_resource_group
}

data "ibm_resource_instance" "test-pdns-instance" {
  name = var.powervs_dns
  resource_group_id = data.ibm_resource_group.rg.id
}

data "ibm_dns_zones" "ds_pdnszone" {
  instance_id = data.ibm_resource_instance.test-pdns-instance.guid
}

module "master" {
  source = "./instance"

  ibmcloud_api_key = var.powervs_api_key
  image_name = var.powervs_image_name
  memory = var.powervs_memory
  networks = [var.powervs_network_name]
  powervs_service_instance_id = var.powervs_service_id
  processors = var.powervs_processors
  ssh_key_name = var.powervs_ssh_key
  vm_name = "${var.cluster_name}-master"
  #user_data = base64encode(file("${path.module}/user_data.txt"))
  user_data = base64encode(templatefile("${path.module}/user_data.tmpl",{port=var.apiserver_port, extra_domain="${var.cluster_name}-master.${var.powervs_dns_zone}", release_marker=var.release_marker, build_version=var.build_version}))
  ibmcloud_region = var.powervs_region
  ibmcloud_zone = var.powervs_zone
}

resource "null_resource" "wait-for-cloud-init-completes" {
  connection {
    type = "ssh"
    user = "root"
    host = module.master.addresses[0].external_ip
  }
  provisioner "remote-exec" {
    inline = [
    "cloud-init status -w"
    ]
  }
}

resource "ibm_dns_resource_record" "test-pdns-resource-record-a" {
  instance_id = data.ibm_resource_instance.test-pdns-instance.guid
  zone_id     = local.zoneids[var.powervs_dns_zone]
  type        = "A"
  name        = "${var.cluster_name}-master"
  rdata       = module.master.addresses[0].external_ip
  ttl         = 3600
}

locals {
  zoneids = zipmap(data.ibm_dns_zones.ds_pdnszone.dns_zones[*]["name"], data.ibm_dns_zones.ds_pdnszone.dns_zones[*]["zone_id"])
}
