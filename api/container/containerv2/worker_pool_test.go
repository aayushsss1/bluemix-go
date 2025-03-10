package containerv2

import (
	"log"
	"net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("workerpools", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//Create
	Describe("Create", func() {
		Context("When creating workerpool is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/vpc/createWorkerPool"),
						ghttp.VerifyJSON(`{"cluster":"bm64u3ed02o93vv36hb0","flavor":"b2.4x16", "hostPool":"hostpoolid1", "name":"mywork211","vpcID":"6015365a-9d93-4bb4-8248-79ae0db2dc26","workerCount":1,"zones":[], "entitlement":""}`),
						ghttp.RespondWith(http.StatusCreated, `{
							"workerPoolID":"string"
						}`),
					),
				)
			})

			It("should create Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}
				params := WorkerPoolRequest{
					Cluster:    "bm64u3ed02o93vv36hb0",
					HostPoolID: "hostpoolid1",
					CommonWorkerPoolConfig: CommonWorkerPoolConfig{
						Flavor:      "b2.4x16",
						Name:        "mywork211",
						VpcID:       "6015365a-9d93-4bb4-8248-79ae0db2dc26",
						WorkerCount: 1,
						Zones:       []Zone{},
						Entitlement: "",
					},
				}
				_, err := newWorkerPool(server.URL()).CreateWorkerPool(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When creating workerpool is successful with OS", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/vpc/createWorkerPool"),
						ghttp.VerifyJSON(`{"cluster":"bm64u3ed02o93vv36hb0","flavor":"b2.4x16", "hostPool":"hostpoolid1", "name":"mywork211","vpcID":"6015365a-9d93-4bb4-8248-79ae0db2dc26","workerCount":1,"zones":[], "entitlement":"", "operatingSystem":"REDHAT_7_64"}`),
						ghttp.RespondWith(http.StatusCreated, `{
							"workerPoolID":"string"
						}`),
					),
				)
			})

			It("should create Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}
				params := WorkerPoolRequest{
					Cluster:    "bm64u3ed02o93vv36hb0",
					HostPoolID: "hostpoolid1",
					CommonWorkerPoolConfig: CommonWorkerPoolConfig{
						Flavor:          "b2.4x16",
						Name:            "mywork211",
						VpcID:           "6015365a-9d93-4bb4-8248-79ae0db2dc26",
						WorkerCount:     1,
						Zones:           []Zone{},
						Entitlement:     "",
						OperatingSystem: "REDHAT_7_64",
					},
				}
				_, err := newWorkerPool(server.URL()).CreateWorkerPool(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When creating workerpool is successful with kms enabled and provided by different account", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/vpc/createWorkerPool"),
						ghttp.VerifyJSON(`{"cluster":"bm64u3ed02o93vv36hb0","flavor":"b2.4x16", "hostPool":"hostpoolid1", "name":"mywork211","vpcID":"6015365a-9d93-4bb4-8248-79ae0db2dc26","workerCount":1,"zones":[], "entitlement":"", "workerVolumeEncryption": {"kmsInstanceID": "kmsid", "workerVolumeCRKID": "rootkeyid", "kmsAccountID":"OtherAccountID"}}`),
						ghttp.RespondWith(http.StatusCreated, `{
							"workerPoolID":"string"
						}`),
					),
				)
			})

			It("should create Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}
				params := WorkerPoolRequest{
					Cluster:    "bm64u3ed02o93vv36hb0",
					HostPoolID: "hostpoolid1",
					CommonWorkerPoolConfig: CommonWorkerPoolConfig{
						Flavor:      "b2.4x16",
						Name:        "mywork211",
						VpcID:       "6015365a-9d93-4bb4-8248-79ae0db2dc26",
						WorkerCount: 1,
						Zones:       []Zone{},
						Entitlement: "",
						WorkerVolumeEncryption: &WorkerVolumeEncryption{
							KmsInstanceID:     "kmsid",
							WorkerVolumeCRKID: "rootkeyid",
							KMSAccountID:      "OtherAccountID",
						},
					},
				}
				_, err := newWorkerPool(server.URL()).CreateWorkerPool(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When creating workerpool is successful with kms enabled and provided by different account", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/vpc/createWorkerPool"),
						ghttp.VerifyJSON(`{"cluster":"bm64u3ed02o93vv36hb0","flavor":"b2.4x16", "hostPool":"hostpoolid1", "name":"mywork211","vpcID":"6015365a-9d93-4bb4-8248-79ae0db2dc26","workerCount":1,"zones":[], "entitlement":"", "workerVolumeEncryption": {"kmsInstanceID": "kmsid", "workerVolumeCRKID": "rootkeyid", "kmsAccountID":"OtherAccountID"}}`),
						ghttp.RespondWith(http.StatusCreated, `{
							"workerPoolID":"string"
						}`),
					),
				)
			})

			It("should create Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}
				params := WorkerPoolRequest{
					Cluster:    "bm64u3ed02o93vv36hb0",
					HostPoolID: "hostpoolid1",
					CommonWorkerPoolConfig: CommonWorkerPoolConfig{
						Flavor:      "b2.4x16",
						Name:        "mywork211",
						VpcID:       "6015365a-9d93-4bb4-8248-79ae0db2dc26",
						WorkerCount: 1,
						Zones:       []Zone{},
						Entitlement: "",
						WorkerVolumeEncryption: &WorkerVolumeEncryption{
							KmsInstanceID:     "kmsid",
							WorkerVolumeCRKID: "rootkeyid",
							KMSAccountID:      "OtherAccountID",
						},
					},
				}
				_, err := newWorkerPool(server.URL()).CreateWorkerPool(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When creating workerpool is successful with SecondaryStorageOption", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/vpc/createWorkerPool"),
						ghttp.VerifyJSON(`{"cluster":"bm64u3ed02o93vv36hb0","flavor":"b2.4x16", "hostPool":"hostpoolid1", "name":"mywork211","vpcID":"6015365a-9d93-4bb4-8248-79ae0db2dc26","workerCount":1,"zones":[], "entitlement":"", "secondaryStorageOption": "secondarystoragename1"}`),
						ghttp.RespondWith(http.StatusCreated, `{
							"workerPoolID":"string"
						}`),
					),
				)
			})

			It("should create Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}
				params := WorkerPoolRequest{
					Cluster:    "bm64u3ed02o93vv36hb0",
					HostPoolID: "hostpoolid1",
					CommonWorkerPoolConfig: CommonWorkerPoolConfig{
						Flavor:                 "b2.4x16",
						Name:                   "mywork211",
						VpcID:                  "6015365a-9d93-4bb4-8248-79ae0db2dc26",
						WorkerCount:            1,
						Zones:                  []Zone{},
						Entitlement:            "",
						SecondaryStorageOption: "secondarystoragename1",
					},
				}
				_, err := newWorkerPool(server.URL()).CreateWorkerPool(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When creating workerpool is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/vpc/createWorkerPool"),
						ghttp.VerifyJSON(`{"cluster":"bm64u3ed02o93vv36hb0","flavor":"b2.4x16","name":"mywork211","vpcID":"6015365a-9d93-4bb4-8248-79ae0db2dc26","workerCount":1,"zones":[],"entitlement":""}
`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to create workerpool`),
					),
				)
			})

			It("should return error during creating workerpool", func() {
				params := WorkerPoolRequest{
					Cluster: "bm64u3ed02o93vv36hb0",
					CommonWorkerPoolConfig: CommonWorkerPoolConfig{
						Flavor:      "b2.4x16",
						Name:        "mywork211",
						VpcID:       "6015365a-9d93-4bb4-8248-79ae0db2dc26",
						WorkerCount: 1,
						Zones:       []Zone{},
						Entitlement: "",
					},
				}
				target := ClusterTargetHeader{}
				_, err := newWorkerPool(server.URL()).CreateWorkerPool(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//getworkerpools
	Describe("Get", func() {
		Context("When Get workerpool is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorkerPool"),
						ghttp.RespondWith(http.StatusCreated, `{
							"dedicatedHostPoolId": "dedicatedhostpoolid1",
							"flavor": "string",
							"id": "string",
							"isolation": "string",
							"lifecycle": {
							  "actualState": "string",
							  "desiredState": "string"
							},
							"operatingSystem": "REDHAT_7_64",
							"poolName": "string",
							"provider": "string",
							"vpcID": "string",
							"workerCount": 0,
							"zones": [
							  {
								"id": "string",
								"subnets": [
								  {
									"id": "string",
									"primary": true
								  }
								],
								"workerCount": 0
							  }
							]
						  }`),
					),
				)
			})

			It("should get Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}

				wp, err := newWorkerPool(server.URL()).GetWorkerPool("aaa", "bbb", target)
				Expect(err).NotTo(HaveOccurred())
				Expect(wp.HostPoolID).To(BeIdenticalTo("dedicatedhostpoolid1"))
				Expect(wp.OperatingSystem).To(BeIdenticalTo("REDHAT_7_64"))
			})
		})
		Context("When get workerpool is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorkerPool"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get workerpool`),
					),
				)
			})

			It("should return error during get workerpool", func() {
				target := ClusterTargetHeader{}
				_, err := newWorkerPool(server.URL()).GetWorkerPool("aaa", "bbb", target)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When Get workerpool is successful and worker volume encyiption is enabled", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorkerPool"),
						ghttp.RespondWith(http.StatusCreated, `{
							"flavor": "string",
							"id": "string",
							"isolation": "string",
							"lifecycle": {
							  "actualState": "string",
							  "desiredState": "string"
							},
							"poolName": "string",
							"provider": "string",
							"vpcID": "string",
							"workerCount": 0,
							"zones": [
							  {
								"id": "string",
								"subnets": [
								  {
									"id": "string",
									"primary": true
								  }
								],
								"workerCount": 0
							  }
							],
							"workerVolumeEncryption": {
								"workerVolumeCRKID": "crk",
								"kmsInstanceID": "kmsid"
							}
						  }`),
					),
				)
			})

			It("should get Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}

				wpresp, err := newWorkerPool(server.URL()).GetWorkerPool("aaa", "bbb", target)
				Expect(err).NotTo(HaveOccurred())
				Expect(wpresp.WorkerVolumeEncryption.KmsInstanceID).Should(Equal("kmsid"))
				Expect(wpresp.WorkerVolumeEncryption.WorkerVolumeCRKID).Should(Equal("crk"))
			})
		})
		Context("When Get workerpool is successful and worker volume encyiption is enabled and provided by another account", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorkerPool"),
						ghttp.RespondWith(http.StatusCreated, `{
							"flavor": "string",
							"id": "string",
							"isolation": "string",
							"lifecycle": {
							  "actualState": "string",
							  "desiredState": "string"
							},
							"poolName": "string",
							"provider": "string",
							"vpcID": "string",
							"workerCount": 0,
							"zones": [
							  {
								"id": "string",
								"subnets": [
								  {
									"id": "string",
									"primary": true
								  }
								],
								"workerCount": 0
							  }
							],
							"workerVolumeEncryption": {
								"workerVolumeCRKID": "crk",
								"kmsInstanceID": "kmsid",
								"kmsAccountID":"OtherAccountID"
							}
						  }`),
					),
				)
			})

			It("should get Workerpool in a cluster with KMSAccountID present in WorkerVolumeEncryption", func() {
				target := ClusterTargetHeader{}

				wpresp, err := newWorkerPool(server.URL()).GetWorkerPool("aaa", "bbb", target)
				Expect(err).NotTo(HaveOccurred())
				Expect(wpresp.WorkerVolumeEncryption.KmsInstanceID).Should(Equal("kmsid"))
				Expect(wpresp.WorkerVolumeEncryption.WorkerVolumeCRKID).Should(Equal("crk"))
				Expect(wpresp.WorkerVolumeEncryption.KMSAccountID).Should(Equal("OtherAccountID"))
			})
		})
		Context("When Get workerpool is successful with UserDefinedSecondaryDisk", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorkerPool"),
						ghttp.RespondWith(http.StatusCreated, `{
							"flavor": "flavor1",
							"id": "id1",
							"lifecycle": {
							  "actualState": "actual",
							  "desiredState": "desired"
							},
							"operatingSystem": "REDHAT_7_64",
							"poolName": "name1",
							"provider": "provider1",
							"vpcID": "vpcid1",
							"workerCount": 0,
							"secondaryStorageOption": {
								"Count": 1,
								"DeviceType": "devtype1",
								"RAIDConfiguration": "config1",
								"Size": 2,
								"name": "name3",
								"profile": "profile3"
							  }
						  }`),
					),
				)
			})

			It("should get Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}

				wp, err := newWorkerPool(server.URL()).GetWorkerPool("aaa", "bbb", target)
				Expect(err).NotTo(HaveOccurred())
				storage := DiskConfigResp{
					Count:             1,
					DeviceType:        "devtype1",
					RAIDConfiguration: "config1",
					Size:              2,
					Name:              "name3",
					Profile:           "profile3",
				}
				Expect(*wp.SecondaryStorageOption).To(BeIdenticalTo(storage))
			})
		})
	})

	//List
	//getworkerpools
	Describe("List", func() {
		Context("When list workerpool is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorkerPools"),
						ghttp.RespondWith(http.StatusCreated, `[{
							"flavor": "string",
							"id": "string",
							"isolation": "string",
							"lifecycle": {
							  "actualState": "string",
							  "desiredState": "string"
							},
							"poolName": "string",
							"provider": "string",
							"vpcID": "string",
							"workerCount": 0,
							"zones": [
							  {
								"id": "string",
								"subnets": [
								  {
									"id": "string",
									"primary": true
								  }
								],
								"workerCount": 0
							  }
							]
						  }]`),
					),
				)
			})

			It("should list Workerpools in a cluster", func() {
				target := ClusterTargetHeader{}

				_, err := newWorkerPool(server.URL()).ListWorkerPools("aaa", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When list workerpool is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorkerPools"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list workerpool`),
					),
				)
			})

			It("should return error during get workerpools", func() {
				target := ClusterTargetHeader{}
				_, err := newWorkerPool(server.URL()).ListWorkerPools("aaa", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//Delete
	Describe("Delete", func() {
		Context("When delete of worker is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v1/clusters/test/workerpools/abc-123-def-ghi"),
						ghttp.RespondWith(http.StatusOK, `{							
						}`),
					),
				)
			})

			It("should delete workerpool", func() {
				target := ClusterTargetHeader{}
				err := newWorkerPool(server.URL()).DeleteWorkerPool("test", "abc-123-def-ghi", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When cluster delete is failed", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v1/clusters/test/workerpools/abc-123-def-ghi"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to delete worker`),
					),
				)
			})

			It("should return error service key delete", func() {
				target := ClusterTargetHeader{}
				err := newWorkerPool(server.URL()).DeleteWorkerPool("test", "abc-123-def-ghi", target)
				Expect(err).To(HaveOccurred())
			})
		})

		//Resize
		Describe("Resize", func() {
			Context("When resizing workerpool is successful", func() {
				BeforeEach(func() {
					server = ghttp.NewServer()
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest(http.MethodPost, "/v2/resizeWorkerPool"),
							ghttp.VerifyJSON(`{"cluster":"bm64u3ed02o93vv36hb0","workerpool":"mywork211","size":5}`),
						),
					)
				})
				It("should resize Workerpool in a cluster", func() {
					target := ClusterTargetHeader{}
					params := ResizeWorkerPoolReq{
						Cluster:    "bm64u3ed02o93vv36hb0",
						Workerpool: "mywork211",
						Size:       5,
					}
					err := newWorkerPool(server.URL()).ResizeWorkerPool(params, target)
					Expect(err).NotTo(HaveOccurred())
				})
			})
			Context("When resizing workerpool is unsuccessful", func() {
				BeforeEach(func() {
					server = ghttp.NewServer()
					server.SetAllowUnhandledRequests(true)
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest(http.MethodPost, "/v2/resizeWorkerPool"),
							ghttp.VerifyJSON(`{"cluster":"bm64u3ed02o93vv36hb0","workerpool":"mywork211","size":5}`),
							ghttp.RespondWith(http.StatusInternalServerError, `Failed to resize workerpool`),
						),
					)
				})

				It("should return error during resizing workerpool", func() {
					params := ResizeWorkerPoolReq{
						Cluster:    "bm64u3ed02o93vv36hb0",
						Workerpool: "mywork211",
						Size:       5,
					}
					target := ClusterTargetHeader{}
					err := newWorkerPool(server.URL()).ResizeWorkerPool(params, target)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})

func newWorkerPool(url string) WorkerPool {

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.VpcContainerService,
	}
	return newWorkerPoolAPI(&client)
}
