locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

# source blocks configure your builder plugins; your source is then used inside
# build blocks to create resources. A build block runs provisioners and
# post-processors on an instance created by the source.
source "amazon-ebs" "aws-ami" {
  ami_name      = "packer example ${local.timestamp}"
  instance_type = "t2.micro"
  region        = "us-west-1"
  source_ami_filter {
    filters = {
      name                = "ubuntu-eks/k8s_1.19/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
}

source "vagrant" "local-vm" {
  communicator = "ssh"
  source_path = "ubuntu/xenial64"
  provider = "virtualbox"
  add_force = true
  output_dir = "local-vm"
}


build {
  sources = ["vagrant.local-vm"]

  provisioner "shell" {
    inline = [
      "sudo apt-get update",
      "sudo echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections",
      "DEBIAN_FRONTEND=noninteractive sudo apt-get install -yq openvpn",
    ]
  }

  provisioner "file"{
    source = "./vpn_conf"
    destination = "/tmp/vpn_conf"
  }

  provisioner "shell" {
    inline = [
      "sudo cp /tmp/vpn_conf/* /etc/openvpn/",
      "sudo service openvpn start",
    ]
  }
}
