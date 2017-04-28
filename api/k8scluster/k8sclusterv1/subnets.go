package k8sclusterv1

import (
	"fmt"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
)

//Subnet ...
type Subnet struct {
	ID          string           `json:"id"`
	Type        string           `json:"type"`
	VlanID      string           `json:"vlan_id"`
	IPAddresses []string         `json:"ip_addresses"`
	Properties  SubnetProperties `json:"properties"`
}

//SubnetProperties ...
type SubnetProperties struct {
	CIDR              string `json:"cidr"`
	NetworkIdentifier string `json:"network_identifier"`
	Note              string `json:"note"`
	SubnetType        string `json:"subnet_type"`
	DisplayLabel      string `json:"display_label"`
	Gateway           string `json:"gateway"`
}

//Subnets interface
type Subnets interface {
	AddSubnet(clusterName string, subnetID string, target *ClusterTargetHeader) error
	List(target *ClusterTargetHeader) ([]Subnet, error)
}

type subnet struct {
	client *clusterClient
	config *bluemix.Config
}

func newSubnetAPI(c *clusterClient) Subnets {
	return &subnet{
		client: c,
		config: c.config,
	}
}

//GetSubnets ...
func (r *subnet) List(target *ClusterTargetHeader) ([]Subnet, error) {
	subnets := []Subnet{}
	_, err := r.client.get("/v1/subnets", &subnets, target)
	if err != nil {
		return nil, err
	}

	return subnets, err
}

//AddSubnetToCluster ...
func (r *subnet) AddSubnet(name string, subnetID string, target *ClusterTargetHeader) error {
	rawURL := fmt.Sprintf("/v1/clusters/%s/subnets/%s", name, subnetID)
	_, err := r.client.put(rawURL, nil, nil, target)
	return err
}
