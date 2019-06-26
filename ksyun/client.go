package ksyun

import (
	"github.com/ksc/ksc-sdk-go/service/eip"
	"github.com/ksc/ksc-sdk-go/service/kec"
	"github.com/ksc/ksc-sdk-go/service/slb"
	"github.com/ksc/ksc-sdk-go/service/vpc"
)

type KsyunClient struct {
	region  string
	eipconn *eip.Eip
	slbconn *slb.Slb
	vpcconn *vpc.Vpc
	kecconn *kec.Kec
}
