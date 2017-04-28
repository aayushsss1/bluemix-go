package k8sclusterv1

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/trace"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/file_helpers"
	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/helpers"
)

//ClusterInfo ...
type ClusterInfo struct {
	GUID              string
	CreatedDate       string
	DataCenter        string
	ID                string
	IngressHostname   string
	IngressSecretName string
	Location          string
	MasterKubeVersion string
	ModifiedDate      string
	Name              string
	Region            string
	ServerURL         string
	State             string
	IsPaid            bool
	WorkerCount       int
}

//ClusterCreateResponse ...
type ClusterCreateResponse struct {
	ID string
}

//ClusterTargetHeader ...
type ClusterTargetHeader struct {
	OrgID     string
	SpaceID   string
	AccountID string
}

//ClusterSoftlayerHeader ...
type ClusterSoftlayerHeader struct {
	SoftLayerUsername string
	SoftLayerAPIKey   string
}

//ClusterCreateRequest ...
type ClusterCreateRequest struct {
	Billing     string
	Datacenter  string
	Isolation   string
	MachineType string
	Name        string
	PrivateVlan string
	PublicVlan  string
	WorkerNum   int
	NoSubnet    bool
}

// ServiceBindRequest ...
type ServiceBindRequest struct {
	ClusterNameOrID         string
	SpaceGUID               string `json:"spaceGUID" binding:"required"`
	ServiceInstanceNameOrID string `json:"serviceInstanceGUID" binding:"required"`
	NamespaceID             string `json:"namespaceID" binding:"required"`
}

// ServiceBindResponse ...
type ServiceBindResponse struct {
	ServiceInstanceGUID string `json:"serviceInstanceGUID" binding:"required"`
	NamespaceID         string `json:"namespaceID" binding:"required"`
	SecretName          string `json:"secretName"`
	Binding             string `json:"binding"`
}

//Clusters interface
type Clusters interface {
	Create(params *ClusterCreateRequest, target *ClusterTargetHeader) (ClusterCreateResponse, error)
	List(target *ClusterTargetHeader) ([]ClusterInfo, error)
	Delete(name string, target *ClusterTargetHeader) error
	Find(name string, target *ClusterTargetHeader) (ClusterInfo, error)
	GetClusterConfig(name, homeDir string, target *ClusterTargetHeader) (string, error)
	UnsetCredentials(target *ClusterTargetHeader) error
	SetCredentials(slUsername, slAPIKey string, target *ClusterTargetHeader) error
	BindService(params *ServiceBindRequest, target *ClusterTargetHeader) (ServiceBindResponse, error)
	UnBindService(clusterNameOrID, namespaceID, serviceInstanceGUID string, target *ClusterTargetHeader) error
}

type clusters struct {
	client *clusterClient
	config *bluemix.Config
}

func newClusterAPI(c *clusterClient) Clusters {
	return &clusters{
		client: c,
		config: c.config,
	}
}

//Create ...
func (r *clusters) Create(params *ClusterCreateRequest, target *ClusterTargetHeader) (ClusterCreateResponse, error) {
	var cluster ClusterCreateResponse
	_, err := r.client.post("/v1/clusters", params, &cluster, target)
	return cluster, err
}

//Delete ...
func (r *clusters) Delete(name string, target *ClusterTargetHeader) error {
	rawURL := fmt.Sprintf("/v1/clusters/%s", name)
	_, err := r.client.delete(rawURL, target)
	return err
}

func (r *clusters) List(target *ClusterTargetHeader) ([]ClusterInfo, error) {
	clusters := []ClusterInfo{}
	_, err := r.client.get("/v1/clusters", &clusters, target)
	if err != nil {
		return nil, err
	}

	return clusters, err
}

func (r *clusters) Find(name string, target *ClusterTargetHeader) (ClusterInfo, error) {
	rawURL := fmt.Sprintf("/v1/clusters/%s", name)
	cluster := ClusterInfo{}
	_, err := r.client.get(rawURL, &cluster, target)
	if err != nil {
		return cluster, err
	}

	return cluster, err
}

func (r *clusters) GetClusterConfig(name, dir string, target *ClusterTargetHeader) (string, error) {
	rawURL := fmt.Sprintf("/v1/clusters/%s/config", name)
	if !file_helpers.FileExists(dir) {
		return "", fmt.Errorf("Path: %q, to download the config doesn't exist", dir)
	}

	now := time.Now()
	zipName := fmt.Sprintf("%s_kubeconfig-%d", name, now.UnixNano())
	downloadPath := fmt.Sprintf("%s/%s.zip", dir, zipName)

	trace.Logger.Println("Will download the kubeconfig at", downloadPath)

	var out *os.File
	var err error
	if out, err = os.Create(downloadPath); err != nil {
		return "", err
	}
	defer out.Close()
	defer file_helpers.RemoveFile(downloadPath)
	_, err = r.client.get(rawURL, out, target)
	if err != nil {
		return "", err
	}

	trace.Logger.Println("Downloaded the kubeconfig at", downloadPath)

	if err = helpers.Unzip(downloadPath, dir); err != nil {
		return "", err
	}

	var unzippedFolderPath string
	homeDirFiles, _ := ioutil.ReadDir(dir)
	for _, homeDirFile := range homeDirFiles {
		if homeDirFile.IsDir() && strings.HasPrefix(homeDirFile.Name(), "kubeConfig") {
			unzippedFolderPath = fmt.Sprintf("%s/%s", dir, homeDirFile.Name())
			break
		}
	}

	if unzippedFolderPath == "" {
		return "", errors.New("There is no directory with prefix kubeConfig in the unzipped file")
	}

	//Rename the config folder to prefix with the cluster name for better identification in the directory
	targetDirPath := filepath.Join(filepath.Dir(unzippedFolderPath), zipName)
	if err = os.Rename(unzippedFolderPath, targetDirPath); err != nil {
		return "", err
	}

	homeDirFiles, err = ioutil.ReadDir(targetDirPath)
	if err != nil {
		return "", fmt.Errorf("Couldn't read %q. Error occured: %v", targetDirPath, err)
	}
	for _, homeDirFile := range homeDirFiles {
		if strings.HasSuffix(homeDirFile.Name(), ".yml") {
			return filepath.Join(targetDirPath, homeDirFile.Name()), nil
		}
	}
	return "", errors.New("Unable to locate kube config in zip archive")

}

func (r *clusters) UnsetCredentials(target *ClusterTargetHeader) error {
	rawURL := fmt.Sprintf("/v1/credentials")
	_, err := r.client.delete(rawURL, target)
	return err
}

func (r *clusters) SetCredentials(slUsername, slAPIKey string, target *ClusterTargetHeader) error {
	_, err := r.client.post("/v1/credentials", nil, nil, target, &ClusterSoftlayerHeader{
		SoftLayerAPIKey:   slAPIKey,
		SoftLayerUsername: slUsername,
	})
	return err
}

//BindService ...
func (r *clusters) BindService(params *ServiceBindRequest, target *ClusterTargetHeader) (ServiceBindResponse, error) {
	rawURL := fmt.Sprintf("/v1/clusters/%s/services", params.ClusterNameOrID)
	payLoad := struct {
		SpaceGUID               string `json:"spaceGUID" binding:"required"`
		ServiceInstanceNameOrID string `json:"serviceInstanceGUID" binding:"required"`
		NamespaceID             string `json:"namespaceID" binding:"required"`
	}{
		SpaceGUID:               params.SpaceGUID,
		ServiceInstanceNameOrID: params.ServiceInstanceNameOrID,
		NamespaceID:             params.NamespaceID,
	}
	var cluster ServiceBindResponse
	fmt.Println(params)
	_, err := r.client.post(rawURL, payLoad, &cluster, target)
	return cluster, err
}

//UnBindService ...
func (r *clusters) UnBindService(clusterNameOrID, namespaceID, serviceInstanceGUID string, target *ClusterTargetHeader) error {
	rawURL := fmt.Sprintf("/v1/clusters/%s/services/%s/%s", clusterNameOrID, namespaceID, serviceInstanceGUID)
	_, err := r.client.delete(rawURL, target)
	return err
}
