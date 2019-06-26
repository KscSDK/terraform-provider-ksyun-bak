package ksyun

import (
	"github.com/ksc/ksc-sdk-go/ksc"
	"github.com/ksc/ksc-sdk-go/ksc/utils"
	"github.com/ksc/ksc-sdk-go/service/eip"
	"github.com/ksc/ksc-sdk-go/service/kec"
	"github.com/ksc/ksc-sdk-go/service/slb"
	"github.com/ksc/ksc-sdk-go/service/vpc"
)

// Config is the configuration of ksyun meta data
type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Insecure  bool
}

// Client will returns a client with connections for all product
func (c *Config) Client() (*KsyunClient, error) {
	var client KsyunClient
	//init ksc client info
	client.region = c.Region
	cli := ksc.NewClient(c.AccessKey, c.SecretKey)
	cfg := &ksc.Config{
		Region: &c.Region,
	}
	url := &utils.UrlInfo{
		UseSSL: c.Insecure,
		Locate: false,
	}
	client.vpcconn = vpc.SdkNew(cli, cfg, url)
	client.eipconn = eip.SdkNew(cli, cfg, url)
	client.slbconn = slb.SdkNew(cli, cfg, url)
	client.kecconn = kec.SdkNew(cli, cfg, url)
	return &client, nil
}
