package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// func (o *Observation) Hash() string {
// 	// serialize the objects map
// 	b := new(bytes.Buffer)
// 	e := json.NewEncoder(b)
// 	err := e.Encode(o.objects)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var h maphash.Hash
// 	defer h.Reset()
// 	h.Write(b.Bytes())

// 	return fmt.Sprint(h.Sum64())
// }

func main() {
	// initialize a kubernets client
	localConfig, err := clientcmd.BuildConfigFromFlags("", "/Users/tgoodwin/.kube/config")
	if err != nil {
		panic(err.Error())
	}
	c, err := kubernetes.NewForConfig(localConfig)
	pods, err := c.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Pod list resource version: ", pods.ResourceVersion)
	for _, pod := range pods.Items {
		uid := pod.GetUID()
		fmt.Println("Pod Name:", pod.Name, "UID:", uid, "RV:", pod.ResourceVersion)
	}
}
