locals {
  availability_zones = ["${var.region}a", "${var.region}c"]
}

resource "aws_vpc" "horahora_vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "horahora-vpc-${var.environment}"
  }
}

resource "aws_internet_gateway" "internet_gateway" {
  vpc_id = aws_vpc.horahora_vpc.id

  tags = {
    Name = "horahora-internet-gateway-${var.environment}"
  }
}

module "vpc_subnets" {
  source = "git::https://github.com/cloudposse/terraform-aws-dynamic-subnets.git?ref=master"

  name        = "horahora"
  environment = var.environment

  vpc_id             = aws_vpc.horahora_vpc.id
  igw_id             = aws_internet_gateway.internet_gateway.id
  cidr_block         = aws_vpc.horahora_vpc.cidr_block
  max_subnet_count   = length(local.availability_zones) * 4 // 4 per AZ, good enough
  availability_zones = local.availability_zones
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
  source                 = "git::https://github.com/cloudposse/terraform-aws-s3-bucket?ref=0.17.1"
  enabled                = true
  versioning_enabled     = true
  allowed_bucket_actions = ["s3:GetObject", "s3:ListBucket", "s3:GetBucketLocation", "s3:PutObject"]
  name                   = "otomads"
  namespace              = "horahora"
  stage                  = var.environment
  region                 = var.region
  block_public_acls      = false
  block_public_policy    = false
  ignore_public_acls     = false

  cors_rule_inputs = list({
    allowed_headers = ["*"]
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
    expose_headers  = []
    max_age_seconds = 300
  })
}

module "eks_cluster" {
  source = "git::https://github.com/cloudposse/terraform-aws-eks-cluster.git?ref=master"
  name   = "horahora-eks-cluster-${var.environment}"

  workers_security_group_ids = [module.eks_workers.security_group_id]
  workers_role_arns          = [module.eks_workers.workers_role_arn]

  vpc_id     = aws_vpc.horahora_vpc.id
  subnet_ids = module.vpc_subnets.private_subnet_ids
  region     = var.region
}

module "eks_workers" {
  source     = "git::https://github.com/cloudposse/terraform-aws-eks-workers.git?ref=master"
  name       = "horahora-eks-workers"
  stage      = var.environment
  subnet_ids = module.vpc_subnets.private_subnet_ids
  vpc_id     = aws_vpc.horahora_vpc.id

  instance_type                          = "t2.micro"
  min_size                               = 1
  max_size                               = 2
  cpu_utilization_high_threshold_percent = 60
  cpu_utilization_low_threshold_percent  = 20

  cluster_endpoint                   = module.eks_cluster.eks_cluster_endpoint
  cluster_security_group_id          = module.eks_cluster.security_group_id
  cluster_certificate_authority_data = module.eks_cluster.eks_cluster_certificate_authority_data
  cluster_name                       = "horahora-eks-cluster-${var.environment}"
}


