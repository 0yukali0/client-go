/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type tenant struct {
	User            string
	Num             uint64
	ResourcesLabels map[string]string
	Resources       map[string]string
}

var ResourceTypes []string
var appsID map[string]string //appID: user
/*
func main() {
	var wg sync.WaitGroup
	num := uint64(50)
	tenants := []tenant{
                tenant{
                        "user1",
                        num,
                        map[string]string{"cpu": "1000", "memory": "8000000000", "duration": "10"},
                        map[string]string{"cpu": "1", "memory": "8G", "duration": "10"},
                },
                tenant{
                        "user2",
                        num * 2,
                        map[string]string{"cpu": "1000", "memory": "4000000000", "duration": "10"},
                        map[string]string{"cpu": "1", "memory": "4G", "duration": "10"},
                },
                tenant{
                        "user3",
                        num * 2,
                        map[string]string{"cpu": "2000", "memory": "2000000000", "duration": "10"},
                        map[string]string{"cpu": "2", "memory": "2G", "duration": "10"},
                },
                tenant{
                        "user4",
                        num,
                        map[string]string{"cpu": "4000", "memory": "2000000000", "duration": "10"},
                        map[string]string{"cpu": "4", "memory": "2G", "duration": "10"},
                },
        }

	for _, user := range tenants {
		go create(user)
		wg.Add(1)
	}
	wg.Wait()
}
*/
func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
        num := uint64(50)
	ResourceTypes = []string{"cpu", "memory"}
	appsID = make(map[string]string, 0)
	podsClient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)
	tenants := []tenant{
		tenant{
			"user1",
			num,
			map[string]string{"cpu": "1000", "memory": "8000000000", "duration": "10"},
			map[string]string{"cpu": "1", "memory": "8G", "duration": "10"},
		},
		tenant{
			"user2",
			num * 2,
			map[string]string{"cpu": "1000", "memory": "4000000000", "duration": "10"},
			map[string]string{"cpu": "1", "memory": "4G", "duration": "10"},
		},
		tenant{
			"user3",
			num * 2,
			map[string]string{"cpu": "2000", "memory": "2000000000", "duration": "10"},
			map[string]string{"cpu": "2", "memory": "2G", "duration": "10"},
		},
		tenant{
			"user4",
			num,
			map[string]string{"cpu": "4000", "memory": "2000000000", "duration": "10"},
			map[string]string{"cpu": "4", "memory": "2G", "duration": "10"},
		},
	}
	for _, user := range tenants {
		for index := uint64(1); index <= user.Num; index++ {
			c1Resources := make(map[apiv1.ResourceName]resource.Quantity)
			c1Resources[apiv1.ResourceMemory] = resource.MustParse(user.Resources["memory"])
			c1Resources[apiv1.ResourceCPU] = resource.MustParse(user.Resources["cpu"])
			appID := fmt.Sprintf("app-%s-%06d", user.User, index)
			appsID[appID] = user.User
			pod := &apiv1.Pod{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Pod",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("pod-%s-%06d", user.User, index),
					Namespace: "default",
					Labels: map[string]string{
						"applicationId":                appID,
						"queue":                        "root.sandbox",
						"yunikorn.apache.org/username": user.User,
						"vcore":                        user.ResourcesLabels["cpu"],
						"memory":                       user.ResourcesLabels["memory"],
						"duration":                     user.ResourcesLabels["duration"],
					},
				},
				Spec: apiv1.PodSpec{
					SchedulerName: "yunikorn",
					RestartPolicy: "Never",
					Containers: []apiv1.Container{
						{
							Name:    "sleep",
							Image:   "alpine:latest",
							Command: []string{"sleep", user.Resources["duration"]},
							Resources: apiv1.ResourceRequirements{
								Requests: c1Resources,
								Limits:   c1Resources,
							},
						},
					},
				},
			}
			if _, err :=  podsClient.Create(context.TODO(), pod, metav1.CreateOptions{}); err != nil {
				panic(err)
			}
		}
		fmt.Printf("%s has %d application\n", user.User, user.Num)
	}
}

/*
func create(user tenant) {
 var kubeconfig *string
        if home := homedir.HomeDir(); home != "" {
                kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
        } else {
                kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
        }
        flag.Parse()

        config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
        if err != nil {
                panic(err)
        }
        clientset, err := kubernetes.NewForConfig(config)
        if err != nil {
                panic(err)
        }
        ResourceTypes = []string{"cpu", "memory"}
        appsID = make(map[string]string, 0)
        podsClient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)

        for index := uint64(1); index <= user.Num; index++ {
		c1Resources := make(map[apiv1.ResourceName]resource.Quantity)
                c1Resources[apiv1.ResourceMemory] = resource.MustParse(user.Resources["memory"])
                c1Resources[apiv1.ResourceCPU] = resource.MustParse(user.Resources["cpu"])
                appID := fmt.Sprintf("app-%s-%06d", user.User, index)
                appsID[appID] = user.User
                pod := &apiv1.Pod{
			TypeMeta: metav1.TypeMeta{
                        Kind:       "Pod",
                        APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
                                        Name:      fmt.Sprintf("pod-%s-%06d", user.User, index),
                                        Namespace: "default",
                                        Labels: map[string]string{
                                                "applicationId":                appID,
                                                "queue":                        "root.sandbox",
                                                "yunikorn.apache.org/username": user.User,
                                                "vcore":                        user.ResourcesLabels["cpu"],
                                                "memory":                       user.ResourcesLabels["memory"],
                                                "duration":                     user.ResourcesLabels["duration"],
                                        },
                                },
                                Spec: apiv1.PodSpec{
                                        SchedulerName: "yunikorn",
                                        RestartPolicy: "Never",
                                        Containers: []apiv1.Container{
                                                {
                                                        Name:    "sleep",
                                                        Image:   "alpine:latest",
                                                        Command: []string{"sleep", user.Resources["duration"]},
                                                        Resources: apiv1.ResourceRequirements{
                                                                Requests: c1Resources,
                                                                Limits:   c1Resources,
                                                        },
                                                },
                                        },
                                },
                }
                        _, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
                        if err != nil {
                                panic(err)
                        }
        }
                fmt.Printf("%s has %d application\n", user.User, user.Num)
}
*/
/*
	nodesClient := clientset.CoreV1().Nodes()
	wg := new(sync.WaitGroup)
	for timeStamp := uint64(0); timeStamp < 20000; timeStamp++ {
		wg.Add(1)
		nodes, err := nodesClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		//pods := podsClient.List(context.TODO(), metav1.ListOptions{})
		go MIGAndBias(nodes, wg)
		// go Res(pods)
		time.Sleep(time.Second)
		wg.Wait()
	}
}

func MIGAndBias(nodes *apiv1.NodeList, wg *sync.WaitGroup) {
	defer wg.Done()
	utilizations := make(map[string][]float64, 0)
	averageUtilizations := make(map[string]float64, 0)
	for _, node := range nodes.Items {
		var mig float64
		for index, resrouceType := range ResourceTypes {
			cap := node.Status.Capacity[apiv1.ResourceName(resrouceType)]
			avail := node.Status.Allocatable[apiv1.ResourceName(resrouceType)]
			capValue, ok1 := cap.AsInt64()
			availValue, ok2 := avail.AsInt64()
			if !ok1 || !ok2 {
				fmt.Printf("%s fail\n", resrouceType)
				return
			}
			fmt.Printf("res %s, %d/%d\n", resrouceType, availValue, capValue)

			utilization := float64(1) - (float64(capValue) / float64(availValue))
			utilizations[resrouceType] = append(utilizations[resrouceType], utilization)
			averageUtilizations[resrouceType] += utilization
			if index == 0 {
				mig = utilization
			} else {
				fmt.Printf("MIG:%s %f\n", node.Spec.ProviderID, math.Abs(mig-utilization))
			}
		}
	}

	for _, resrouceType := range ResourceTypes {
		average := averageUtilizations[resrouceType] / float64(len(nodes.Items))

		bias := float64(0)
		for _, utilization := range utilizations[resrouceType] {
			bias += math.Pow(utilization-average, 2)
		}
		fmt.Printf("%s bias: %f\nc", resrouceType, math.Sqrt(bias))
	}
	return
}

func Res(pods *apiv1.PodList, wg *sync.WaitGroup) {
	defer wg.Done()

	return
}
*/
