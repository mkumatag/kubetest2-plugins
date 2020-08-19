package deployer

import (
	"bytes"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/providers"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/wait"
	"log"
	"net"
	gohttp "net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/spf13/pflag"
	//"github.com/tmc/scp"
	"github.com/pkg/sftp"
	"sigs.k8s.io/kubetest2/pkg/types"
)

const (
	powervsSessionTimeout = time.Duration(2) * time.Minute
	// Name is the name of the deployer
	Name = "powervs"
	pollInterval = 2 * time.Minute
	pollTimeout = 30 * time.Minute
)

// New implements deployer.New for kind
func New(opts types.Options) (types.Deployer, *pflag.FlagSet) {
	d := &deployer{
		commonOptions: opts,
		logsDir:       filepath.Join(opts.ArtifactsDir(), "logs"),
	}
	return d, bindFlags(d)
}

func (d *deployer) setIBMPISession() (err error) {
	sess, err := session.New(&bluemix.Config{BluemixAPIKey: d.apikey})
	if err != nil{
		log.Printf("Unable to Create a new Bluemix session: %+v", err)
		return
	}
	iamAuthRepository, err := authentication.NewIAMAuthRepository(sess.Config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		log.Printf("Unable to NewIAMAuthRepository: %+v", err)
		return
	}
	err = iamAuthRepository.AuthenticateAPIKey(sess.Config.BluemixAPIKey)
	if err != nil {
		log.Printf("Unable to AuthenticateAPIKey: %+v", err)
		return
	}
	d.IBMPISession, err = ibmpisession.New(sess.Config.IAMAccessToken,
		d.region,
		d.debug,
		powervsSessionTimeout,
		d.userAccount,
		d.zone)
	return
}

func bindFlags(d *deployer) *pflag.FlagSet {
	os.Setenv("KUBECONFIG", "/tmp/kubeconfig")
	flags := pflag.NewFlagSet(Name, pflag.ContinueOnError)
	flags.StringVar(
		&d.clusterName, "cluster-name", "k8s-cluster-ppc64le", "Cluster name",
	)
	flags.StringVar(
		&d.userAccount, "user-account", "", "IBM Cloud User Account ID(command: ibmcloud account list)",
	)
	flags.BoolVar(
		&d.debug, "debug", false, "Enable debug flag for APIs(Caution: This will print the headers including the tokens)",
	)
	flags.StringVar(
		&d.apikey, "api-key", "", "IBM Cloud API Key used for accessing the APIs",
	)
	flags.StringVar(
		&d.region, "region", "", "IBM Cloud PowerVS region name",
	)
	flags.StringVar(
		&d.zone, "zone", "", "IBM Cloud PowerVS zone name",
	)
	flags.StringVar(
		&d.instanceID, "instance-id", "", "IBM Cloud PowerVS service instance ID(get GUID from command: ibmcloud resource service-instances --long)",
	)
	flags.StringVar(
		&d.networkID, "network-id", "", "Network ID(command: ibmcloud pi nets)",
	)
	flags.StringVar(
		&d.imageID, "image-id", "", "Image ID(command: ibmcloud pi imgs)",
	)
	flags.Float64Var(
		&d.memory, "memory", 8, "Memory in GBs",
	)
	flags.Float64Var(
		&d.processor, "processor", 0.5, "Processor Units",
	)
	flags.StringVar(
		&d.processorType, "processor-type", "shared", "Processor Types(dedicated, shared)",
	)
	flags.IntVar(
		&d.apiServerPort, "api-server-port", 992, "API Server Port Address",
		)
	return flags
}

// assert that New implements types.NewDeployer
var _ types.NewDeployer = New

type deployer struct {
	// generic parts
	commonOptions      types.Options
	logsDir            string // dir to export logs to
	doInit             sync.Once
	IBMPISession       *ibmpisession.IBMPISession
	debug              bool
	apikey             string
	clusterName        string
	userAccount        string
	region             string
	zone               string
	instanceID         string
	networkID          string
	ipaddresstable     map[string]bool
	privateIPAddresses []string
	publicIPAddresses  []string
	imageID            string
	memory             float64
	processor          float64
	processorType      string
	masterPublicIP     string
	masterPrivateIP    string
	apiServerPort      int
}

func (d *deployer) init() error {
	var err error
	d.doInit.Do(func() { err = d.initialize() })
	return err
}

func (d *deployer) initialize() error {
	if d.commonOptions.ShouldBuild() {
		if err := d.verifyBuildFlags(); err != nil {
			return fmt.Errorf("init failed to check build flags: %s", err)
		}
	}

	if d.commonOptions.ShouldUp() {
		if err := d.verifyUpFlags(); err != nil {
			return fmt.Errorf("init failed to verify flags for up: %s", err)
		}
	}
	if d.commonOptions.ShouldDown() {
		if err := d.verifyDownFlags(); err != nil {
			return fmt.Errorf("init failed to verify flags for down: %s", err)
		}
	}
	err := d.setIBMPISession()
	if err != nil {
		log.Fatalf("Failed to setIBMPISession with err: %+v\n", err)
	}
	return nil
}

func (d *deployer) verifyBuildFlags() error{
	//TODO: yet to implement
	return nil
}

func (d *deployer) verifyUpFlags() error{
	//TODO: yet to implement
	return nil
}

func (d *deployer) verifyDownFlags() error{
	//TODO: yet to implement
	return nil
}

func cidrRange(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, nil
	default:
		return ips[1 : len(ips)-1], nil
	}
}

// getIPAddresses IP address range for the mentioned start and end address
func getIPAddresses(start, end string) (ips []string) {
	for ip := net.ParseIP(start); bytes.Compare(ip, net.ParseIP(end)) <=0; inc(ip) {
		ips = append(ips, ip.String())
	}
	return
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func Find(a []string, x string) (int, error){
	for i, n := range a {
		if x == n {
			return i, nil
		}
	}
	return 0, errors.New("element not found")
}

func (d *deployer) reserveIPAddress() (string, string, error){

	for ip, used := range d.ipaddresstable {
		if !used {
			if i, err := Find(d.privateIPAddresses, ip); err == nil{
				d.ipaddresstable[ip] = true
				return d.privateIPAddresses[i], d.publicIPAddresses[i], nil
			}
		}
	}
	return "", "", errors.New("Unable to reserve the IP address, no free IP left in the network")
}

func getKeyFile(privateKeyFile string) (key ssh.Signer, err error) {
	buf, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}

func getKubeconfigFile(RemoteMachineIP string) error {
	key, err := getKeyFile("/Users/manjunath/.ssh/id_rsa")
	if err != nil {
		panic(err)
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", RemoteMachineIP), config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftp.Close()

	srcFile, err := sftp.Open("/etc/kubernetes/admin.conf")
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create("/tmp/kubeconfig")
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile.WriteTo(dstFile)
	return nil
}

func (d *deployer) Up() error {
	if err := d.init(); err != nil {
		return fmt.Errorf("up failed to init: %s", err)
	}

	networkClient := instance.NewIBMPINetworkClient(d.IBMPISession, d.instanceID)
	network, err := networkClient.Get(d.networkID, d.instanceID, powervsSessionTimeout)
	if err != nil{
		return err
	}

	ipaddrs, err := cidrRange(*network.Cidr)
	if err != nil{
		return err
	}

	for _, ipaddressRange := range network.IPAddressRanges{
		d.privateIPAddresses = append(d.privateIPAddresses, getIPAddresses(*ipaddressRange.StartingIPAddress, *ipaddressRange.EndingIPAddress)...)
	}
	for _, publicipaddressRange := range network.PublicIPAddressRanges{
		d.publicIPAddresses = append(d.publicIPAddresses, getIPAddresses(*publicipaddressRange.StartingIPAddress, *publicipaddressRange.EndingIPAddress)...)
	}

	fmt.Printf("privateIPAddresses: %+v, publicipaddressRange: %+v", d.privateIPAddresses, d.publicIPAddresses)

	d.ipaddresstable = make(map[string]bool)
	for i := 0; i < len(ipaddrs); i +=1 {
		d.ipaddresstable[ipaddrs[i]] = false
	}

	//spew.Dump(network)

	d.ipaddresstable[network.Gateway] = true

	networkPorts, err := networkClient.GetAllPort(d.networkID, d.instanceID, powervsSessionTimeout)
	if err != nil{
		return err
	}
	//spew.Dump(networkPorts)
	for i := 0; i < len(networkPorts.Ports); i +=1 {
		d.ipaddresstable[*networkPorts.Ports[i].IPAddress] = true
	}

	//spew.Dump(d.ipaddresstable)

	instanceClient := instance.NewIBMPIInstanceClient(d.IBMPISession, d.instanceID)

	d.masterPrivateIP, d.masterPublicIP, err = d.reserveIPAddress()
	if err != nil {
		return err
	}
	user_data := fmt.Sprintf(user_data_template, d.apiServerPort, d.masterPublicIP)
	fmt.Printf("user_data: %s\n", user_data)
	params := p_cloud_p_vm_instances.PcloudPvminstancesPostParams{
		Body: &models.PVMInstanceCreate{
			ImageID: &d.imageID,
			KeyPairName: "mkumatag-pub-key",
			Memory: &d.memory,
			Processors: &d.processor,
			Networks: []*models.PVMInstanceAddNetwork{{NetworkID: &d.networkID, IPAddress: d.masterPrivateIP}},
			ServerName: &d.clusterName,
			ProcType: &d.processorType,
			UserData: b64.StdEncoding.EncodeToString([]byte(user_data)),
		},
	}
	instances, err := instanceClient.Create(&params,d.instanceID, powervsSessionTimeout)
	if err != nil{
		log.Printf("error while creating instance: %+v\n", err)
		return err
	}
	id := (*instances)[0].PvmInstanceID
	//id := "2b6008eb-47fa-4298-a5ae-22d524193079"
	pollErr := wait.PollImmediate(pollInterval, pollTimeout, func() (bool, error) {
		masterInstance, err := instanceClient.Get(*id, d.instanceID, powervsSessionTimeout)
		//spew.Dump(masterInstance)
		if err != nil{
			return false, fmt.Errorf("failed to get the powervm: %v", err)
		}
		if *masterInstance.Status == "ERROR" {
			return false, fmt.Errorf("instance went into error state, fault message: %s", masterInstance.Fault.Message)
		}
		if masterInstance.Health.Status != "OK" {
			return false, nil
		}
		return true, nil
	})

	if pollErr != nil {
		return fmt.Errorf("failed to bring lpar online: %v", pollErr)
	}
	err = getKubeconfigFile(d.masterPublicIP)
	if err != nil{
		return err
	}
	//if _, err := waitPoll(20 * time.Minute, func()(bool, error){
	//	masterInstance, err := instanceClient.Get(id, d.instanceID, powervsSessionTimeout)
	//	spew.Dump(masterInstance)
	//	if err != nil{
	//		log.Printf("Failed to get the powervm %+v", err)
	//		return false, err
	//	}
	//	if *masterInstance.Status == "ERROR" {
	//		return false, errors.New(masterInstance.Fault.Message)
	//	}
	//	if masterInstance.Health.Status != "ACTIVE" {
	//		return false, nil
	//	}
	//	return true, nil
	//}); err != nil{
	//	return errors.New("failed to bring instance to active state")
	//}
	//imgclient := instance.NewIBMPIImageClient(d.IBMPISession, "ddc04489-b76e-40ba-b17a-b76ae25087f5")
	//
	//images, err := imgclient.GetAll("ddc04489-b76e-40ba-b17a-b76ae25087f5")
	//if err != nil{
	//	log.Printf("error while getting the images! err: %+v", err)
	//}
	//fmt.Println(images)
	//spew.Dump(images)
	//panic("implement me")
	return nil
}

func waitPoll(timeout time.Duration, do func()(bool, error))(bool, error){
	t := time.After(timeout)
	tick := time.Tick(1 * time.Minute)
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-t:
			return false, errors.New("timed out")
		// Got a tick, we should check on doSomething()
		case <-tick:
			ok, err := do()
			// Error from doSomething(), we should bail
			if err != nil {
				return false, err
				// doSomething() worked! let's finish up
			} else if ok {
				return true, nil
			}
			// doSomething() didn't work yet, but it didn't fail, so let's try again
			// this will exit up to the for loop
		}
	}
}

func (d *deployer) kubectl() string{
	return fmt.Sprintf("kubectl -s https://%s:%s/", d.masterPublicIP, d.apiServerPort)
}

func (d *deployer) Down() error {
	if err := d.init(); err != nil {
		return fmt.Errorf("up failed to init: %s", err)
	}
	panic("implement me")
}

func (d *deployer) IsUp() (up bool, err error) {
	panic("implement me")
}

func (d *deployer) DumpClusterLogs() error {
	if err := d.init(); err != nil {
		return fmt.Errorf("dumpClusterLogs failed to init: %s", err)
	}
	panic("implement me")
}

func (d *deployer) Build() error {
	if err := d.init(); err != nil {
		return fmt.Errorf("build failed to init: %s", err)
	}
	panic("implement me")
}

// assert that deployer implements types.Deployer
var _ types.Deployer = &deployer{}