package tfvars

type TFVars struct {
	ReleaseMarker 	string	`json:"release_marker"`
	BuildVersion	string	`json:"build_version"`
	ClusterName   string	`json:"cluster_name"`
	ApiServerPort int		`json:"apiserver_port"`
}