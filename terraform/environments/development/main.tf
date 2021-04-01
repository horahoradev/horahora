locals {
  availability_zones = ["${var.region}a", "${var.region}c"]
}

resource "aws_vpc" "horahora_vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "horahora-vpc-${var.environment}"
  }

  enable_dns_hostnames = true
  enable_dns_support   = true
}

//resource "aws_vpc" "horahora_vpc_alt" {
//  cidr_block = "172.31.0.0/16"
//
//  tags = {
//    Name = "horahora-vpc-alt-${var.environment}"
//  }
//
//  enable_dns_hostnames = true
//  enable_dns_support   = true
//}

resource "aws_internet_gateway" "internet_gateway" {
  vpc_id = aws_vpc.horahora_vpc.id

  tags = {
    Name = "horahora-internet-gateway-${var.environment}"
  }
}

//resource "aws_internet_gateway" "internet_gateway_alt" {
//  vpc_id = aws_vpc.horahora_vpc_alt.id
//
//  tags = {
//    Name = "horahora-internet-gateway-alt-${var.environment}"
//  }
//}

module "vpc_subnets" {
  source = "git::https://github.com/cloudposse/terraform-aws-dynamic-subnets.git?ref=master"

  name        = "horahora"
  environment = var.environment

  vpc_id                  = aws_vpc.horahora_vpc.id
  igw_id                  = aws_internet_gateway.internet_gateway.id
  cidr_block              = aws_vpc.horahora_vpc.cidr_block
  max_subnet_count        = length(local.availability_zones) * 4 // 4 per AZ, good enough
  availability_zones      = local.availability_zones
  map_public_ip_on_launch = true
  tags = {
    "kubernetes.io/cluster/horahora-dev-cluster" = "shared"
  }
}


// Origin bucket for video files
// TODO: https://docs.aws.amazon.com/AmazonS3/latest/dev/example-bucket-policies.html#example-bucket-policies-use-case-3
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
  source = "git::https://github.com/cloudposse/terraform-aws-eks-cluster.git?ref=0.29.0"
  name   = "horahora-${var.environment}"

  kubernetes_version = "1.17"

  workers_security_group_ids = [module.eks_workers.security_group_id]
  workers_role_arns          = [module.eks_workers.workers_role_arn]

  vpc_id     = aws_vpc.horahora_vpc.id
  subnet_ids = module.vpc_subnets.public_subnet_ids
  region     = var.region
}

module "eks_workers" {
  source     = "git::https://github.com/cloudposse/terraform-aws-eks-workers.git?ref=master"
  name       = "horahora-eks-workers"
  stage      = var.environment
  subnet_ids = module.vpc_subnets.public_subnet_ids
  vpc_id     = aws_vpc.horahora_vpc.id

  instance_type                          = "c5.large"
  min_size                               = 1
  max_size                               = 1
  cpu_utilization_high_threshold_percent = 60
  cpu_utilization_low_threshold_percent  = 20

  //  allowed_cidr_blocks = ["0.0.0.0/0"]

  //  use_custom_image_id = true
  //  image_id            = "ami-0a9a4d789dc726512"

  cluster_endpoint                   = module.eks_cluster.eks_cluster_endpoint
  cluster_security_group_id          = module.eks_cluster.security_group_id
  cluster_certificate_authority_data = module.eks_cluster.eks_cluster_certificate_authority_data
  cluster_name                       = "horahora-${var.environment}-cluster"
  associate_public_ip_address        = true

  mixed_instances_policy = {
    instances_distribution = {
      on_demand_base_capacity                  = 0
      on_demand_allocation_strategy            = "prioritized"
      spot_allocation_strategy                 = "lowest-price"
      spot_instance_pools                      = 2
      on_demand_percentage_above_base_capacity = 0
      spot_max_price                           = ""
    }
    override = [{
      instance_type     = "c5.large"
      weighted_capacity = 1
    }]
  }

  key_name = "bastion"
}

resource "aws_security_group" "rds_whitelist" {
  name   = "rds_whitelist-${var.environment}"
  vpc_id = aws_vpc.horahora_vpc.id

  ingress {
    description = "TLS from VPC"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/8",
      "172.16.0.0/12",
    "217.138.252.0/24"] // lol bruh moment FIXME
  }
}

module "scheduledb" {
  source         = "git::https://github.com/cloudposse/terraform-aws-rds.git?ref=master"
  name           = "scheduledb-${var.environment}"
  engine         = "postgres"
  engine_version = "12.5"
  instance_class = "db.t3.micro"

  vpc_id              = aws_vpc.horahora_vpc.id
  subnet_ids          = module.vpc_subnets.public_subnet_ids
  publicly_accessible = true // a bit lazy, but makes it easier to apply migrations for now
  db_parameter_group  = "postgres12"
  database_port       = 5432
  database_name       = "scheduler"

  database_user     = "scheduler"
  database_password = var.schedulerdb_password

  backup_retention_period = 5 // in days
  allocated_storage       = "20"

  associate_security_group_ids = [aws_security_group.rds_whitelist.id]

  apply_immediately = true
}

module "userdb" {
  source         = "git::https://github.com/cloudposse/terraform-aws-rds.git?ref=master"
  name           = "userdb-${var.environment}"
  engine         = "postgres"
  engine_version = "12.5"
  instance_class = "db.t3.micro"

  vpc_id              = aws_vpc.horahora_vpc.id
  subnet_ids          = module.vpc_subnets.public_subnet_ids
  publicly_accessible = true // a bit lazy, but makes it easier to apply migrations for now
  db_parameter_group  = "postgres12"
  database_port       = 5432
  database_name       = "userservice"

  database_user     = "userservice"
  database_password = var.userdb_password

  backup_retention_period = 5 // in days
  allocated_storage       = "20"

  apply_immediately            = true
  associate_security_group_ids = [aws_security_group.rds_whitelist.id]
}

module "videodb" {
  source         = "git::https://github.com/cloudposse/terraform-aws-rds.git?ref=master"
  name           = "videodb-${var.environment}"
  engine         = "postgres"
  engine_version = "12.5"
  instance_class = "db.t3.micro"

  vpc_id              = aws_vpc.horahora_vpc.id
  subnet_ids          = module.vpc_subnets.public_subnet_ids
  publicly_accessible = true // a bit lazy, but makes it easier to apply migrations for now
  db_parameter_group  = "postgres12"
  database_port       = 5432
  database_name       = "videoservice"

  database_user     = "videoservice"
  database_password = var.videodb_password

  backup_retention_period = 5 // in days
  allocated_storage       = "20"

  apply_immediately            = true
  associate_security_group_ids = [aws_security_group.rds_whitelist.id]
}

// Shared by scheduler and videoservice for now
// TODO: cluster mode, HA
module "video_redis" {
  source         = "git::https://github.com/cloudposse/terraform-aws-elasticache-redis.git?ref=master"
  name           = "video_redis-${var.environment}"
  engine_version = "5.0.5"
  family         = "redis5.0"

  availability_zones                   = local.availability_zones
  subnets                              = module.vpc_subnets.private_subnet_ids
  vpc_id                               = aws_vpc.horahora_vpc.id
  instance_type                        = "cache.t2.micro"
  cluster_mode_enabled                 = false
  cluster_size                         = 1
  cluster_mode_num_node_groups         = 1
  cluster_mode_replicas_per_node_group = 0
  automatic_failover_enabled           = false
  allowed_cidr_blocks                  = ["10.0.0.0/8", "172.16.0.0/12", "217.138.252.0/24"]
  transit_encryption_enabled           = false
}