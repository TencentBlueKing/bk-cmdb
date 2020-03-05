package metadata

type Vpc struct {
	VpcId   string `json:"bk_vpc_id", bson:"bk_vpc_id"`
	VpcName string `json:"bk_vpc_name", bson:"bk_vpc_name"`
}

type Instance struct {
	InstanceId    string `json:"bk_instance_id", bson:"bk_instance_id"`
	InstanceName  string `json:"bk_instance_name", bson:"bk_instance_name"`
	PrivateIp     string `json:"bk_host_innerip", bson:"bk_host_innerip"`
	PublicIp      string `json:"bk_host_outerip", bson:"bk_host_outerip"`
	InstanceState string `json:"bk_instance_state", bson:"bk_instance_state"`
	VpcId         string `json:"bk_vpc_id", bson:"bk_vpc_id"`
}
