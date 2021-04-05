locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

source "amazon-ebs" "aws-ami" {
  ami_name      = "packer example ${local.timestamp}"
  instance_type = "t2.micro"
  region        = "us-west-1"
  source_ami_filter {
    filters = {
      name                = "ubuntu-eks/k8s_1.17/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"
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
  source_path = "ubuntu/bionic64"
  provider = "virtualbox"
  add_force = true
  output_dir = "local-vm"
}


build {
  sources = [
    "vagrant.local-vm",
    "amazon-ebs.aws-ami"
  ]

  provisioner "shell" {
    inline = [
      "sudo apt-get update",
      "sudo echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections",
      "sudo apt-get install -yq software-properties-common && sudo add-apt-repository 'deb http://us.archive.ubuntu.com/ubuntu/ bionic main restricted'",
      "sudo apt-get update",
      "sudo apt-get install -yq openvpn",
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
