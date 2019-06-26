# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_instance" "default" {
  image_id="d3290df6-3597-4f83-b5ae-48356e91ad46"
  instance_type="N3.2B"
  system_disk{
    disk_type="SSD3.0"
    disk_size=30
  }
  data_disk_gb=100
  data_disk =[
   {
      type="SSD3.0"
      size=20
      delete_with_instance=true
   }
 ]
  max_count=1
  min_count=1
  subnet_id="9a9ac083-cd22-4e75-af56-593a91463972"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id="b8591529-2741-4f09-af3c-49a954e3e4fa"
  private_ip_address=""
  instance_name="xuan-tf-update"
  instance_name_suffix=""
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=[]
}
