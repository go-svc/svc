package consul

import (
	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/consul/api"
)

func NewClient(conf *api.Config) *api.Client {
	client, err := api.NewClient(conf)
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Error occurred while creating the Consul api client.")
	}

	return client
}
