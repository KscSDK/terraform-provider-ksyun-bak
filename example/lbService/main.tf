# Specify the provider and access details
provider "ksyun" {
  region = "cn-shanghai-2"
}
data "ksyun_lines" "default" {
  output_file=""
  line_name="BGP"
}
resource "ksyun_vpc" "default" {
  vpc_name        = "tf-acc-vpc"
  cidr_block = "192.168.0.0/16"
}
# Create Load Balancer
resource "ksyun_lb" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  load_balancer_name = "tf-xun-2"
  type = "public"
  subnet_id = ""
  load_balancer_state = "start"
  private_ip_address = ""
}
# Create Load Balancer Listener with tcp protocol
resource "ksyun_lb_listener" "default" {
  listener_name = "tf-xun-2"
  listener_port = "8080"
  listener_protocol = "HTTP"
  listener_state = "start"
  load_balancer_id = "${ksyun_lb.default.id}"
  method = "RoundRobin"
  certificate_id = ""
  session {
    session_state = "start"
    session_persistence_period = 100
    cookie_type = "RewriteCookie"
    cookie_name = "cookiexunqq"
  }
}

# Attach instances to Load Balancer
resource "ksyun_lb_listener_server" "default" {
  listener_id = "${ksyun_lb_listener.default.id}"
  real_server_ip = "10.0.77.20"
  real_server_port = 8000
  real_server_type = "host"
  instance_id = "dc746b15-2afc-481c-926d-eb9538b01f39"
  weight = 10
}
# Create Load Balancer Listener Acl
resource "ksyun_lb_acl" "default" {
  load_balancer_acl_name = "tf-xun-2"
}
# Create Load Balancer Listener Acl Entry
resource "ksyun_lb_acl_entry" "default" {
  load_balancer_acl_id = "${ksyun_lb_acl.default.id}"
  cidr_block = "192.168.11.1/32"
  rule_number = 10
  rule_action = "allow"
  protocol = "ip"
}
resource "ksyun_lb_listener_associate_acl" "default" {
  listener_id = "${ksyun_lb_listener.default.id}"
  load_balancer_acl_id = "${ksyun_lb_acl.default.id}"
}
# Create an eip
resource "ksyun_eip" "default" {
  line_id ="${data.ksyun_lines.default.lines.0.line_id}"
  band_width = 1
  charge_type = "PostPaidByDay"
  purchase_time = 1
  project_id =0
}

# Bind eip to Load Balancer
resource "ksyun_eip_associate" "default" {
  allocation_id = "${ksyun_eip.default.id}"
  instance_type = "Slb"
  instance_id = "${ksyun_lb.default.id}"
  network_interface_id = ""
}
resource "ksyun_healthcheck" "default" {
  listener_id = "${ksyun_lb_listener.default.id}"
  health_check_state = "start"
  healthy_threshold = 2
  interval = 20
  timeout = 200
  unhealthy_threshold = 2
  url_path = "/monitor"
  is_default_host_name = false
  host_name = "www.baidu.com"
}


