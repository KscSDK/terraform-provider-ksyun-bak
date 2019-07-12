package ksyun

var lineKeys = map[string]bool{
	"LineName": true,
	"LineId":   true,
	"LineType": true,
}
var eipKeys = map[string]bool{
	"CreateTime":         true,
	"ProjectId":          true,
	"PublicIp":           true,
	"AllocationId":       true,
	"State":              true,
	"LineId":             true,
	"BandWidth":          true,
	"InstanceType":       true,
	"InstanceId":         true,
	"NetworkInterfaceId": true,
	"InternetGatewayId":  true,
	"BandWidthShareId":   true,
	"IsBandWidthShare":   true,
}
var slbKeys = map[string]bool{
	"CreateTime":        true,
	"LoadBalancerName":  true,
	"VpcId":             true,
	"LoadBalancerId":    true,
	"Type":              true,
	"SubnetId":          true,
	"PublicIp":          true,
	"State":             true,
	"LoadBalancerState": true,
}
var listenerKeys = map[string]bool{
	"CreateTime":       true,
	"LoadBalancerId":   true,
	"ListenerName":     true,
	"ListenerId":       true,
	"ListenerState":    true,
	"CertificateId":    true,
	"ListenerProtocol": true,
	"ListenerPort":     true,
	"Method":           true,
	"HealthCheck":      true,
	"Session":          true,
	"RealServer":       true,
}
var healthCheckKeys = map[string]bool{
	"HealthCheckId":      true,
	"HealthCheckState":   true,
	"HealthyThreshold":   true,
	"Interval":           true,
	"Timeout":            true,
	"UnhealthyThreshold": true,
}
var sessionKeys = map[string]bool{
	"SessionPersistencePeriod": true,
	"SessionState":             true,
	"CookieType":               true,
	"CookieName":               true,
}

var serverKeys = map[string]bool{
	"RegisterId":      true,
	"RealServerIp":    true,
	"RealServerPort":  true,
	"RealServerType":  true,
	"InstanceId":      true,
	"RealServerState": true,
	"Weight":          true,
}
var lbAclKeys = map[string]bool{
	"CreateTime":              true,
	"LoadBalancerAclName":     true,
	"LoadBalancerAclId":       true,
	"LoadBalancerAclEntrySet": true,
}

var lbAclEntryKeys = map[string]bool{
	"LoadBalancerAclId":      true,
	"LoadBalancerAclEntryId": true,
	"CidrBlock":              true,
	"RuleNumber":             true,
	"RuleAction":             true,
	"Protocol":               true,
}
var availabilityZoneKeys = map[string]bool{
	"AvailabilityZoneName":  true,
	"AvailabilityZoneState": true,
}
var vpcNetworkInterfaceKeys = map[string]bool{
	"NetworkInterfaceId":   true,
	"NetworkInterfaceType": true,
	"MacAddress":           true,
	"SecurityGroupSet":     true,
	"InstanceId":           true,
	"InstanceType":         true,
	"PrivateIpAddress":     true,
	"DNS1":                 true,
	"DNS2":                 true,
}
var subnetAvailableAddresseKeys = map[string]bool{
	"AvailableIpAddress": true,
}
var subnetAllocatedIpAddressesKeys = map[string]bool{
	"AvailableIpAddress": true,
}

var vpcKeys = map[string]bool{
	"CidrBlock":  true,
	"CreateTime": true,
	"IsDefault":  true,
	"VpcName":    true,
	"VpcId":      true,
}

var groupIdentifierKeys = map[string]bool{
	"SecurityGroupId":   true,
	"SecurityGroupName": true,
}
var subnetKeys = map[string]bool{
	"CreateTime":           true,
	"VpcId":                true,
	"SubnetId":             true,
	"SubnetType":           true,
	"SubnetName":           true,
	"CidrBlock":            true,
	"DhcpIpFrom":           true,
	"DhcpIpTo":             true,
	"GatewayIp":            true,
	"Dns1":                 true,
	"Dns2":                 true,
	"NetworkAclId":         true,
	"NatId":                true,
	"AvailbleIPNumber":     true,
	"AvailabilityZoneName": true,
}
var vpcSecurityGroupKeys = map[string]bool{
	"CreateTime":            true,
	"VpcId":                 true,
	"SecurityGroupName":     true,
	"SecurityGroupId":       true,
	"SecurityGroupType":     true,
	"SecurityGroupEntrySet": true,
}
var vpcSecurityGroupEntrySetKeys = map[string]bool{
	"Description":          true,
	"SecurityGroupEntryId": true,
	"CidrBlock":            true,
	"Direction":            true,
	"Protocol":             true,
	"IcmpType":             true,
	"IcmpCode":             true,
	"PortRangeFrom":        true,
	"PortRangeTo":          true,
}
var instanceKeys = map[string]bool{
	"InstanceId":            true,
	"ProjectId":             true,
	"InstanceName":          true,
	"InstanceType":          true,
	"InstanceConfigure":     true,
	"ImageId":               true,
	"SubnetId":              true,
	"PrivateIpAddress":      true,
	"InstanceState":         true,
	"Monitoring":            true,
	"NetworkInterfaceSet":   true,
	"SriovNetSupport":       true,
	"IsShowSriovNetSupport": true,
	"CreationDate":          true,
	"AvailabilityZone":      true,
	"AvailabilityZoneName":  true,
	"AutoScalingType":       true,
	"ProductWhat":           true,
	"ChargeType":            true,
	"SystemDisk":            true,
}
var instanceConfigureKeys = map[string]bool{
	"VCPU":       true,
	"GPU":        true,
	"MemoryGb":   true,
	"DataDiskGb": true,
	//"RootDiskGb":   true,
	//"DataDiskType": true,
}
var instanceStateKeys = map[string]bool{
	"Name": true,
}
var monitoringKeys = map[string]bool{
	"State": true,
}
var kecNetworkInterfaceKeys = map[string]bool{
	"NetworkInterfaceId":   true,
	"NetworkInterfaceType": true,
	"MacAddress":           true,
	"SecurityGroupSet":     true,
	"PrivateIpAddress":     true,
	"DNS1":                 true,
	"DNS2":                 true,
}
var kecNetworkInterfaceSetKeys = map[string]bool{
	"NetworkInterfaceId":   true,
	"NetworkInterfaceType": true,
	"PrivateIpAddress":     true,
	"GroupSet":             true,
	"SecurityGroupSet":     true,
}
var systemDiskKeys = map[string]bool{
	"DiskType": true,
	"DiskSize": true,
}
var groupSetKeys = map[string]bool{
	"GroupId": true,
}
var kecSecurityGroupSetKeys = map[string]bool{
	"SecurityGroupId": true,
}
var kecSecurityGroupKeys = map[string]bool{
	"CreateTime":            true,
	"VpcId":                 true,
	"SecurityGroupName":     true,
	"SecurityGroupId":       true,
	"SecurityGroupType":     true,
	"SecurityGroupEntrySet": true,
}
var imageKeys = map[string]bool{
	"ImageId":      true,
	"Name":         true,
	"ImageState":   true,
	"CreationDate": true,
	"Platform":     true,
	"IsPublic":     true,
	"InstanceId":   true,
	"IsNpe":        true,
	"UserCategory": true,
	"SysDisk":      true,
	"Progress":     true,
	"ImageSource":  true,
}
var networkInterfaceKeys = map[string]bool{
	"NetworkInterfaceId":   true,
	"NetworkInterfaceType": true,
	"MacAddress":           true,
	"SecurityGroupSet":     true,
	"InstanceId":           true,
	"InstanceType":         true,
	"PrivateIpAddress":     true,
	"SubnetId":             true,
	"ProjectId":            true,
	"DNS1":                 true,
	"DNS2":                 true,
}
