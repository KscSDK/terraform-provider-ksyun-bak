# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "eu-east-1"
}

# Create Load Balancer Listener with tcp protocol
resource "ksyun_lb_listener" "default" {
  listener_name = "tf-xun",
  listener_port = "8000",
  listener_protocol = "HTTP",
  listener_state = "stop",
  load_balancer_id = "7fae85e4-ab1a-415c-aef9-03a402c79d97",
  method = "RoundRobin"
  certificate_id = ""
  session {
    session_state = "stop"
    session_persistence_period = 100
    cookie_type = "ImplantCookie"
    cookie_name = "cookiexunqq"
  }
}
