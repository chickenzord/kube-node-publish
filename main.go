package main

import (
	"strings"
	"time"

	"github.com/chickenzord/kube-node-publish/pkg/dns"
	"github.com/chickenzord/kube-node-publish/pkg/kube"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	dns.InitConfig()
	config := NewConfigFromEnv()

	// logging config
	logLevel, _ := log.ParseLevel(config.LogLevel)
	log.SetLevel(logLevel)

	// check
	log.Println(config)
	if config.DomainName == "" {
		log.Fatal("domain name required")
	}

	// prepare client
	clientset, _ := kube.NewClientSet(config.InCluster, config.KubeConfig)
	if config.InCluster {
		log.Debug("using in-cluster kubeconfig")
	} else {
		log.Debugf("using kubeconfig: %s", config.KubeConfig)
	}

	// start updating
	var curValue = ""
	for {
		ipAddresses, _ := kube.GetNodeAddresses(clientset)
		newValue := strings.Join(ipAddresses, " ")

		if curValue == newValue {
			log.Debug("no ip address changes")
		} else {
			log.Infof("found %d ip address: %v", len(ipAddresses), ipAddresses)
			if err := dns.EnsureRecord(config.DomainName, "A", newValue); err != nil {
				log.Errorf("cannot update dns record: %v", err)
			} else {
				log.Infof("updated A %s record: %s", config.DomainName, newValue)
				curValue = newValue
			}
		}

		// delay
		duration, _ := time.ParseDuration(config.CheckPeriod)
		log.Debugf("waiting %s", config.CheckPeriod)
		time.Sleep(duration)
	}

}
