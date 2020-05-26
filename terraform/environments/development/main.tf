resource "aws_vpc" "horahora_vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "horahora-vpc-${var.environment}"
  }
}


// FIXME:   Could probably use count to simplify
resource "aws_subnet" "subnet_1" {
  availability_zone = "${var.region}a"
  vpc_id            = aws_vpc.horahora_vpc.id
  cidr_block        = "10.0.1.0/24"
}

resource "aws_subnet" "subnet_2" {
  availability_zone = "${var.region}c"
  vpc_id            = aws_vpc.horahora_vpc.id
  cidr_block        = "10.0.2.0/24"
}

resource "aws_subnet" "subnet_3" {
  availability_zone = "${var.region}a"
  vpc_id            = aws_vpc.horahora_vpc.id
  cidr_block        = "10.0.3.0/24"
}

/*module "eks-cluster" {
  source      = "../../modules/eks_cluster"
  name        = "horahora-eks-cluster"
  environment = var.environment

  subnets = [
    aws_subnet.subnet_1.id,
    aws_subnet.subnet_2.id
  ]
}*/

// Origin bucket for video files
module "video_origin" {
  source                 = "git::https://github.com/cloudposse/terraform-aws-s3-bucket.git?ref=master"
  enabled                = true
  versioning_enabled     = true
  allowed_bucket_actions = ["s3:GetObject", "s3:ListBucket", "s3:GetBucketLocation", "s3:PutObject"]
  name                   = "otomads"
  namespace              = "horahora"
  stage                  = var.environment
  region                 = var.region
}
