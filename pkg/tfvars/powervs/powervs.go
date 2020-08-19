package powervs

type TFVars struct {
	ResourceGroup string	`json:"powervs_resource_group"`
	DNSName 	  string	`json:"powervs_dns"`
	DNSZone		  string	`json:"powervs_dns_zone"`
	Apikey        string	`json:"powervs_api_key"`
//	UserAccount   string	`json:"powervs_user_account"`
	Region        string	`json:"powervs_region"`
	Zone          string	`json:"powervs_zone"`
	ServiceID    string	`json:"powervs_service_id"`
	NetworkName     string	`json:"powervs_network_name"`
	ImageName       string	`json:"powervs_image_name"`
	Memory        float64	`json:"powervs_memory"`
	Processors     float64	`json:"powervs_processors"`
	//ApiServerPort int		`json:"apiserver_port"`
	//Debug         bool		`json:"debug"`
	//ClusterName   string	`json:"powervs_cluster_name"`
	SSHKey        string	`json:"powervs_ssh_key"`
}