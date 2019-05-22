package kube

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sort"
)

func GetNodeAddresses(clientset *kubernetes.Clientset) ([]string, error) {
	ipAddresses := []string{}

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return ipAddresses, err
	}

	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == v1.NodeExternalIP {
				ipAddresses = append(ipAddresses, address.Address)
			}
		}
	}

	sort.Strings(ipAddresses)

	return ipAddresses, nil
}
