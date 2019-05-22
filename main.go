package main

import (
	"strings"
	"time"

	"github.com/chickenzord/kube-node-publish/pkg/config"
	"github.com/chickenzord/kube-node-publish/pkg/dns"
	"github.com/chickenzord/kube-node-publish/pkg/kube"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	dns.InitConfig()
	cfg := config.NewConfigFromEnv()

	// logging config
	logLevel, _ := log.ParseLevel(cfg.LogLevel)
	log.SetLevel(logLevel)

	// check
	log.Printf("version %s (%s)", config.Version, config.GitSha)
	log.Println(cfg)
	if cfg.DomainName == "" {
		log.Fatal("domain name required")
	}

	// prepare client
	clientset, _ := kube.NewClientSet(cfg.InCluster, cfg.KubeConfig)
	if cfg.InCluster {
		log.Info("using in-cluster kubeconfig")
	} else {
		log.Infof("using kubeconfig: %s", cfg.KubeConfig)
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
			if err := dns.EnsureRecord(cfg.DomainName, "A", newValue); err != nil {
				log.Errorf("cannot update dns record: %v", err)
			} else {
				log.Infof("updated A %s record: %s", cfg.DomainName, newValue)
				curValue = newValue
			}
		}

		// delay
		duration, _ := time.ParseDuration(cfg.CheckPeriod)
		log.Debugf("waiting %s", cfg.CheckPeriod)
		time.Sleep(duration)
	}

}
