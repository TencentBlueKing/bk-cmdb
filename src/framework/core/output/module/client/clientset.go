package client

import (
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/output/module/client/v3"
)

var _ Interface = &Clientset{}

type Interface interface {
	CCV3() v3.CCV3Interface
}

type Clientset struct {
	ccv3 *v3.Client
}

func (c *Clientset) CCV3() v3.CCV3Interface {
	if c == nil {
		return nil
	}
	return c.ccv3
}

func NewForConfig(c *config.Config) (*Clientset, error) {
	var cs Clientset
	cs.ccv3 = &v3.New(c, nil)
	client = cs
	return &cs, nil
}

var client *Clientset

func GetClient() *Clientset {
	return client
}
