package cvm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/griesbacher/check_azure/azureHttp"
	"github.com/griesbacher/check_x"
)

func ShowMachines(ac azureHttp.AzureConnector, resourceGroup string) error {
	body, typ, err := ac.RequestWithSub("2014-04-01", fmt.Sprintf("resourceGroups/%s/providers/Microsoft.ClassicCompute/virtualMachines", resourceGroup), "")
	if err != nil {
		return err
	}
	if typ != azureHttp.ContentTypeJSON {
		return azureHttp.ContentError(azureHttp.ContentTypeJSON, typ)
	}
	var vms virtualMachines
	err = json.Unmarshal(body, &vms)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	for i, vm := range vms.Value {
		buffer.WriteString(fmt.Sprintf("%d - %s\n", i, vm.Name))
		buffer.WriteString(fmt.Sprintf("\tType: %s\n", vm.Type))
		buffer.WriteString(fmt.Sprintf("\tLocation: %s\n", vm.Location))
		buffer.WriteString(fmt.Sprintf("\tProvisioningState: %s\n", vm.Properties.ProvisioningState))
		buffer.WriteString(fmt.Sprintf("\tPowerState: %s\n", vm.Properties.InstanceView.PowerState))
		buffer.WriteString(fmt.Sprintf("\tPrivateIpAddress: %s\n", vm.Properties.InstanceView.PrivateIPAddress))
		buffer.WriteString(fmt.Sprintf("\tPublicIpAddress: %s\n", vm.Properties.InstanceView.PublicIPAddresses))
		buffer.WriteString(fmt.Sprintf("\tFQDN: %s\n", vm.Properties.InstanceView.FullyQualifiedDomainName))
		buffer.WriteString(fmt.Sprintf("\tGuestAgentStatus: %s\n", vm.Properties.InstanceView.GuestAgentStatus.Status))
		buffer.WriteString(fmt.Sprintf("\tGuestAgentVersion: %s\n", vm.Properties.InstanceView.GuestAgentStatus.GuestAgentVersion))

		buffer.WriteString(fmt.Sprintf("\tDomainName: %s\n", vm.Properties.DomainName.Name))

		buffer.WriteString(fmt.Sprintf("\tHardwareSize: %s\n", vm.Properties.HardwareProfile.Size))

		buffer.WriteString("\tnetwork\n")
		for j, e := range vm.Properties.NetworkProfile.InputEndpoints {
			buffer.WriteString(fmt.Sprintf("\t%d:\n", j))
			buffer.WriteString(fmt.Sprintf("\t\tEndpointName: %s\n", e.EndpointName))
			buffer.WriteString(fmt.Sprintf("\t\tPublicIPAddress: %s\n", e.PublicIPAddress))
			buffer.WriteString(fmt.Sprintf("\t\tPublicPort: %d\n", e.PublicPort))
			buffer.WriteString(fmt.Sprintf("\t\tPrivatePort: %d\n", e.PrivatePort))
			buffer.WriteString(fmt.Sprintf("\t\tProtocol: %s\n", e.Protocol))
			buffer.WriteString(fmt.Sprintf("\t\tEnableDirectServerReturn: %t\n", e.EnableDirectServerReturn))
		}
		buffer.WriteString(fmt.Sprintf("\tVirtualNetwork: %s\n", vm.Properties.NetworkProfile.VirtualNetwork.Name))
		buffer.WriteString(fmt.Sprintf("\tVirtualNetworkSubnetNames: %s\n", vm.Properties.NetworkProfile.VirtualNetwork.SubnetNames))

		buffer.WriteString(fmt.Sprintf("\tOperatingSystemDisk-Names: %s\n", vm.Properties.StorageProfile.OperatingSystemDisk.DiskName))
		buffer.WriteString(fmt.Sprintf("\tOperatingSystemDisk-OperatingSystem: %s\n", vm.Properties.StorageProfile.OperatingSystemDisk.OperatingSystem))
		buffer.WriteString(fmt.Sprintf("\tOperatingSystemDisk-IoType: %s\n", vm.Properties.StorageProfile.OperatingSystemDisk.IoType))
		buffer.WriteString(fmt.Sprintf("\tOperatingSystemDisk-SourceImageName: %s\n", vm.Properties.StorageProfile.OperatingSystemDisk.SourceImageName))
		buffer.WriteString(fmt.Sprintf("\tOperatingSystemDisk-VhdURI: %s\n", vm.Properties.StorageProfile.OperatingSystemDisk.VhdURI))
		buffer.WriteString(fmt.Sprintf("\tStorageAccount-Name: %s\n", vm.Properties.StorageProfile.OperatingSystemDisk.StorageAccount.Name))
		buffer.WriteString(fmt.Sprintf("\tStorageAccount-Type: %s\n", vm.Properties.StorageProfile.OperatingSystemDisk.StorageAccount.Type))
	}
	check_x.LongExit(check_x.OK, "see below", buffer.String())
	return nil
}

type virtualMachines struct {
	Value []struct {
		ID         string `json:"id"`
		Location   string `json:"location"`
		Name       string `json:"name"`
		Properties struct {
			DebugProfile struct {
				BootDiagnosticsEnabled bool `json:"bootDiagnosticsEnabled"`
			} `json:"debugProfile"`
			DomainName struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"domainName"`
			HardwareProfile struct {
				DeploymentID       string `json:"deploymentId"`
				DeploymentLabel    string `json:"deploymentLabel"`
				DeploymentLocked   bool   `json:"deploymentLocked"`
				DeploymentName     string `json:"deploymentName"`
				PlatformGuestAgent bool   `json:"platformGuestAgent"`
				Size               string `json:"size"`
			} `json:"hardwareProfile"`
			InstanceView struct {
				ComputerName             string `json:"computerName"`
				FaultDomain              int    `json:"faultDomain"`
				FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
				GuestAgentStatus         struct {
					FormattedMessage struct {
						Language string `json:"language"`
						Message  string `json:"message"`
					} `json:"formattedMessage"`
					GuestAgentVersion string `json:"guestAgentVersion"`
					ProtocolVersion   string `json:"protocolVersion"`
					Status            string `json:"status"`
					Timestamp         string `json:"timestamp"`
				} `json:"guestAgentStatus"`
				PowerState        string   `json:"powerState"`
				PrivateIPAddress  string   `json:"privateIpAddress"`
				PublicIPAddresses []string `json:"publicIpAddresses"`
				Status            string   `json:"status"`
				StatusMessage     string   `json:"statusMessage"`
				UpdateDomain      int      `json:"updateDomain"`
			} `json:"instanceView"`
			NetworkProfile struct {
				InputEndpoints []struct {
					EnableDirectServerReturn bool   `json:"enableDirectServerReturn"`
					EndpointName             string `json:"endpointName"`
					PrivatePort              int    `json:"privatePort"`
					Protocol                 string `json:"protocol"`
					PublicIPAddress          string `json:"publicIpAddress"`
					PublicPort               int    `json:"publicPort"`
				} `json:"inputEndpoints"`
				VirtualNetwork struct {
					ID          string   `json:"id"`
					Name        string   `json:"name"`
					SubnetNames []string `json:"subnetNames"`
					Type        string   `json:"type"`
				} `json:"virtualNetwork"`
			} `json:"networkProfile"`
			ProvisioningState string `json:"provisioningState"`
			StorageProfile    struct {
				OperatingSystemDisk struct {
					Caching         string `json:"caching"`
					DiskName        string `json:"diskName"`
					IoType          string `json:"ioType"`
					OperatingSystem string `json:"operatingSystem"`
					SourceImageName string `json:"sourceImageName"`
					StorageAccount  struct {
						ID   string `json:"id"`
						Name string `json:"name"`
						Type string `json:"type"`
					} `json:"storageAccount"`
					VhdURI string `json:"vhdUri"`
				} `json:"operatingSystemDisk"`
			} `json:"storageProfile"`
		} `json:"properties"`
		Type string `json:"type"`
	} `json:"value"`
}
