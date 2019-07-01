# terraform-provider-ksyun
Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.11.x
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-ksyun`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-ksyun
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-ksyun
$ go build
```

Using the provider
----------------------

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-ksyun
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## 快速安装


  若用户只是利用插件，不对插件进行二次开发，且不想安装go环境和编译代码，可直接安装编译好的插件。不同操作系统的插件（terraform-provider-ksyun）位于目录bin-os/下，将其直接拷贝至terraform的plugins默认目录即可。
不同操作系统terraform的plugins默认目录：(https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)。
>注意：首次拷贝，需手动创建plugins文件夹。

### Terraform-provider-ksyun 文件配置：

 #### 1、provider配置

 	provider "ksyun" {
       access_key = "你的ak"
       secret_key = "你的sk"
       region = "cn-beijing-6"
     }
 	
  可以在配置文件中指定，也可以在环境变量中配置，若两处都配置，以配置文件为主。
  在环境变量中配置：
  
```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-ksyun
$ export KSYUN_ACCESS_KEY=xxx
$ export KSYUN_SECRET_KEY=xxx
$ export KSYUN_REGION=xxx
$ export TF_LOG=DEBUG
$ export TF_LOG_PATH=/var/log/terraform.log
$ export OS_DEBUG=1
```
#### 2、data-source使用
  以data+产品线命名的文件夹，代表不同资源，可根据配置条件导出全部符合的资源。若配置条件为空，导出全部。
  以dataEips下的main.tf为例：

 	data "ksyun_eips" "default" {
        //导出的资源会输出到output_result文件中
      output_file = “output_result” 
        //只导出eipId包含在ids中的eip信息
      ids = []
       //只导出project_id=1的eip默认只导出project_id=0的eip
      project_id = [“1"]
       //只导出instance_type=“Ipfwd”的eip
      instance_type = ["Ipfwd"]
      network_interface_id = []
      internet_gateway_id = []
      band_width_share_id = []
      line_id = []
      public_ip = []
    }
  在该目录下执行：
```sh
$ terraform init
$ terraform plan
```
  会将符合条件的eip导出到output_result中
#### 3、resource使用
  以产品线命名的文件夹代表对应产品线的单资源配置，以产品线+Service命名的文件夹代表对应该资源的服务编排。
##### 单资源配置：
  以eips下的main.tf为例：

 	resource "ksyun_eip" "default1" {
      line_id ="cf8b7b95-4651-b96c-db67-b38336f2fe70"
      band_width =1
      charge_type = "PostPaidByDay"
      purchase_time =1
      project_id=0
    }

  在该目录下执行：
```sh
$ terraform init
$ terraform plan //获取操作类型（新建，修改或重建）
$ terraform apply //执行操作
$ terraform destroy //删除eip
$ terraform import ksyun_eip.default1 eipId //导入该eipId的eip信息，一般用于对已有实例的修改
```

##### 服务编排：
  以terraform-provider-ksyun/example/instanceService为例：

###### 1、variables.tf定义变量

 	variable "instance_name" {
      default = "ksyun_instance_tf"
    }
    variable "subnet_name" {
      default = "ksyun_subnet_tf"
    }

  定义变量instance_name，其默认值为“ksyun_instance_tf”
  若用户想自定义变量的值，可在执行terraform 命令时指定:
```sh
$ terraform plan -var ‘instance_name=kec’ -var ‘subnet_name=sub’
```
###### 2、outputs.tf控制台输出

 	output "eip_id" {
      value = "${ksyun_eip.default.id}"
    }


  执行terraform apply 后，控制台会输出：
```sh
$ eip_id=‘e9587b84-0da7-4fd7-a26d-bc56df63b01e’
```
###### 3、main.tf 资源编排

  1、定义了provider为ksyun

  2、定义了三个资源，即镜像、线路和可用区

  3、定义了多个资源配置，包括创建vpc、子网、安全组、安全组规则、主机及绑定eip

### Terraform-provider-ksyun 版本升级：
  terraform v0.12 版本的配置文件与v0.11 版本的配置文件格式不同。该example下的配置文件是基于v0.11.13开发的。若想使用v0.12版本的terraform需对配置文件进行修改，该修改不需手动修改，terraform支持自动修改。可在配置文件目录下直接执行：
```sh
$ terraform 0.12upgrade
```
  terraform会询问是否确认修改，输入yes即可。
