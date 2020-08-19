data "ibm_pi_network" "power_networks" {
    count                = length(var.networks)
    pi_network_name      = var.networks[count.index]
    pi_cloud_instance_id = var.powervs_service_instance_id
}

data "ibm_pi_image" "power_images" {
    pi_image_name        = var.image_name
    pi_cloud_instance_id = var.powervs_service_instance_id
}

resource "ibm_pi_instance" "pvminstance" {
    pi_memory             = var.memory
    pi_processors         = var.processors
    pi_instance_name      = var.vm_name
    pi_proc_type          = var.proc_type
    pi_image_id           = data.ibm_pi_image.power_images.id
    pi_network_ids        = data.ibm_pi_network.power_networks.*.id
    pi_key_pair_name      = var.ssh_key_name
    pi_sys_type           = var.system_type
    pi_cloud_instance_id = var.powervs_service_instance_id
    pi_user_data          = var.user_data
}

