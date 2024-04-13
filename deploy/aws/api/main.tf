terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "5.41.0"
    }
    random = {
      source = "hashicorp/random"
      version = "3.6.0"
    }
  }
}

provider "aws" {
  region = "eu-west-3"
  profile = "kodmain"
}

variable "github_token" {
  description = "GitHub token"
}

resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_instance" "free_tier_arm_instance" {
  # X86 ami-089c89a80285075f7 t2.micro  Amazon Linux 2 # WORKING
  # ARM ami-09e82d7942ffb02d3 t4g.micro Amazon Linux 2 # NOT WORKING (fixed by client.cpu_total_compute in nomad-server.hcl)
  ami           = "ami-09e82d7942ffb02d3"
  instance_type = "t4g.micro"

  tags = {
    Name = "NomadServer"
  }
  
  user_data = <<-EOF
    #!/bin/bash
    yum install -y yum-utils
    yum-config-manager --add-repo https://rpm.releases.hashicorp.com/AmazonLinux/hashicorp.repo
    yum-config-manager --add-repo https://cli.github.com/packages/rpm/gh-cli.repo
    yum -y install nomad docker gh cni-plugins httpd-tools htop
    git clone https://github.com/kodmain/thetiptop /home/ec2-user/thetiptop
    cp /home/ec2-user/thetiptop/deploy/aws/api/nomad.service /etc/systemd/system/nomad.service
    sleep 1
    systemctl enable nomad 
    systemctl enable docker
    systemctl start docker
    systemctl start nomad
    sleep 1
    nomad acl bootstrap > /home/ec2-user/bootstrap.token
    export NOMAD_TOKEN=$(cat /home/ec2-user/bootstrap.token | grep "Secret" |awk '{print $4}')
    export GH_TOKEN="${var.github_token}"
    export GF_ADMIN_PASSWORD="${random_password.password.result}"
    echo "export NOMAD_TOKEN=$NOMAD_TOKEN" >> /home/ec2-user/.bashrc
    echo "export GH_TOKEN='$GH_TOKEN'" >> /home/ec2-user/.bashrc
    echo "export GF_ADMIN_PASSWORD='$GF_ADMIN_PASSWORD'" >> /home/ec2-user/.bashrc
    nomad acl policy apply -description "Deployment" deploy /home/ec2-user/thetiptop/deploy/aws/api/nomad-policy.hcl
    nomad acl token create -name="github" -policy="deploy" > /home/ec2-user/github.token
    export GITHUB_NOMAD_TOKEN=$(cat /home/ec2-user/github.token | grep "Secret" |awk '{print $4}')
    gh secret set NOMAD_TOKEN -b"$GITHUB_NOMAD_TOKEN" --repo kodmain/thetiptop
    sed -i 's/NOMADTOKEN/'"$NOMAD_TOKEN"'/g' /home/ec2-user/thetiptop/deploy/aws/api/jobs/server.hcl
    sleep 1
    nomad job run -token=$NOMAD_TOKEN /home/ec2-user/thetiptop/deploy/aws/api/jobs/server.hcl
    nomad job run -token=$NOMAD_TOKEN -var="grafana_admin_password=$GF_ADMIN_PASSWORD" /home/ec2-user/thetiptop/deploy/aws/api/jobs/middlewares.hcl
  EOF

  iam_instance_profile = aws_iam_instance_profile.traefik_instance_profile.name

  root_block_device {
    volume_size = 10
    volume_type = "gp3"
    encrypted = true
  }
  
  security_groups = [aws_security_group.nomad.name]
  key_name = aws_key_pair.remote.key_name 
}

resource "aws_key_pair" "remote" {
  key_name   = "kodmain"
  public_key = file("~/.ssh/kodmain.pub")
}

resource "aws_security_group" "nomad" {
  name        = "nomad"
  description = "Security Group for Nomad Server"

  /* Disable use nomad.kodmain.run
  ingress {
    from_port   = 4646
    to_port     = 4646
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
   */

  /* Disable use traefik.kodmain.run
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  */

  /* Disable SSH 
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  */

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_iam_role" "traefik_route53_role" {
  name = "TraefikRoute53Role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "traefik_route53_policy" {
  name = "TraefikRoute53Policy"
  role = aws_iam_role.traefik_route53_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "s3:GetObject",
          "s3:ListBucket"
        ],
        Resource = [
          "${aws_s3_bucket.config_bucket.arn}",
          "${aws_s3_bucket.config_bucket.arn}/*"
        ]
      },
      {
        Action = [
          "ssm:UpdateInstanceInformation",
          "route53:GetChange",
          "route53:ChangeResourceRecordSets",
          "route53:ListResourceRecordSets",
          "route53:ListHostedZones",
          "route53:ListHostedZonesByName"
        ],
        Effect = "Allow",
        Resource = "*"
      },
      {
        Action: [
          "cloudwatch:ListMetrics",
          "cloudwatch:GetMetricData",
          "cloudwatch:GetMetricStatistics",
          "cloudwatch:DescribeAlarms",
          "cloudwatch:DescribeAlarmHistory",
          "cloudwatch:DescribeAlarmsForMetric"
        ],
        Effect: "Allow",
        Resource: "*"
      },
      {
        Action = [
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams",
          "logs:GetLogEvents",
          "logs:FilterLogEvents"
        ],
        Effect = "Allow",
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_instance_profile" "traefik_instance_profile" {
  name = "TraefikInstanceProfile"
  role = aws_iam_role.traefik_route53_role.name
}

resource "aws_route53_record" "kodmain_wildcard" {
  zone_id = "Z10052173VRSYMBUSS942"

  name    = "*.kodmain.run"  # Enregistrement wildcard pour tous les sous-domaines
  type    = "A"
  ttl     = 10
  records = [aws_instance.free_tier_arm_instance.public_ip]

  allow_overwrite = true
}

resource "aws_route53_record" "kodmain_internal" {
  zone_id = "Z10052173VRSYMBUSS942"

  name    = "internal.kodmain.run"  # Enregistrement wildcard pour tous les sous-domaines
  type    = "A"
  ttl     = 10
  records = [aws_instance.free_tier_arm_instance.private_ip]

  allow_overwrite = true
}

resource "aws_s3_bucket" "config_bucket" {
  bucket = "config.kodmain"

  tags = {
    Name        = "ConfigBucket"
    Environment = "Production"
  }
}

resource "aws_s3_bucket_policy" "config_bucket_policy" {
  bucket = aws_s3_bucket.config_bucket.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Deny",
        Principal = "*",
        Action    = "s3:*",
        Resource  = [
          "${aws_s3_bucket.config_bucket.arn}",
          "${aws_s3_bucket.config_bucket.arn}/*"
        ],
        Condition = {
          Bool = {
            "aws:SecureTransport": "false"
          }
        }
      }
    ]
  })
}

resource "aws_s3_bucket_versioning" "config_bucket_versioning" {
  bucket = aws_s3_bucket.config_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_object" "config_file" {
  bucket = aws_s3_bucket.config_bucket.bucket
  key    = "config.yml"
  source = "../../../api/config.yml"
}

