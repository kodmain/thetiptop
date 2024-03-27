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
  region = "us-east-1"
  alias  = "global"
  profile = "kodmain"
}

provider "aws" {
  region = "eu-west-3"
  profile = "kodmain"
}

resource "aws_acm_certificate" "cert" {
  provider          = aws.global
  domain_name       = "kodmain.run"
  validation_method = "DNS"

  tags = {
    Environment = "production"
  }
  lifecycle {
    create_before_destroy = true
  }
}

# Attendez la validation du certificat avant de créer la distribution CloudFront
resource "aws_acm_certificate_validation" "cert_validation" {
  provider        = aws.global
  certificate_arn = aws_acm_certificate.cert.arn
  validation_record_fqdns = [for record in aws_route53_record.cert_validation : record.fqdn]
}

resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      type   = dvo.resource_record_type
      record = dvo.resource_record_value
    }
  }

  zone_id = "Z10052173VRSYMBUSS942"
  name    = each.value.name
  type    = each.value.type
  records = [each.value.record]
  ttl     = 60
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
  associate_public_ip_address = true

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
    cp /home/ec2-user/thetiptop/deploy/server/nomad.service /etc/systemd/system/nomad.service
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
    nomad acl policy apply -description "Deployment" deploy /home/ec2-user/thetiptop/deploy/server/nomad-policy.hcl
    nomad acl token create -name="github" -policy="deploy" > /home/ec2-user/github.token
    export GITHUB_NOMAD_TOKEN=$(cat /home/ec2-user/github.token | grep "Secret" |awk '{print $4}')
    gh secret set NOMAD_TOKEN -b"$GITHUB_NOMAD_TOKEN" --repo kodmain/thetiptop
    sed -i 's/NOMADTOKEN/'"$NOMAD_TOKEN"'/g' /home/ec2-user/thetiptop/deploy/jobs/server.hcl
    sleep 1
    nomad job run -token=$NOMAD_TOKEN /home/ec2-user/thetiptop/deploy/jobs/server.hcl
    nomad job run -token=$NOMAD_TOKEN -var="grafana_admin_password=$GF_ADMIN_PASSWORD" /home/ec2-user/thetiptop/deploy/jobs/middlewares.hcl
  EOF

  iam_instance_profile = aws_iam_instance_profile.traefik_instance_profile.name

  root_block_device {
    volume_size = 10
    volume_type = "gp3"
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
  */
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

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


resource "aws_s3_bucket" "app" {
  bucket = "kodmain"
  force_destroy = true
}

resource "aws_s3_bucket_public_access_block" "app_public_access_block" {
  bucket = aws_s3_bucket.app.id

  block_public_acls       = false
  ignore_public_acls      = false
  block_public_policy     = false
  restrict_public_buckets = false
}


resource "aws_s3_bucket_public_access_block" "logs_public_access_block" {
  bucket = aws_s3_bucket.logs.id

  block_public_acls       = false
  ignore_public_acls      = true
  block_public_policy     = true
  restrict_public_buckets = true
}


resource "aws_s3_access_point" "kodmain_access_point" {
  name         = "kodmain"
  bucket       = aws_s3_bucket.app.id

  public_access_block_configuration {
    block_public_acls       = false
    block_public_policy     = false
    ignore_public_acls      = false
    restrict_public_buckets = false
  }
}

resource "aws_s3_bucket_policy" "bucket_log_policy" {
  bucket = aws_s3_bucket.logs.id

   policy = jsonencode({
    Version = "2012-10-17"
    Id      = "logs_policy"
    Statement = [
      {
        Sid       = "HTTPSOnly"
        Effect    = "Deny"
        Principal = "*"
        Action    = "s3:*"
        Resource = [
          aws_s3_bucket.logs.arn,
          "${aws_s3_bucket.logs.arn}/*",
        ]
        Condition = {
          Bool = {
            "aws:SecureTransport" = "false"
          }
        }
      },
    ]
  })
}

resource "aws_s3_bucket" "logs" {
  bucket = "logs.kodmain"
  force_destroy = true
}


resource "aws_cloudfront_distribution" "s3_distribution" {
  origin {
    domain_name = aws_s3_bucket.app.bucket_domain_name
    origin_id   = "S3-${aws_s3_bucket.app.id}"
  }

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"

  tags = {
    Environment = "production"
  }

  aliases = ["kodmain.run"]

  default_cache_behavior {
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = "S3-${aws_s3_bucket.app.id}"
    compress         = true

    forwarded_values {
      query_string = true
      cookies {
        forward = "all"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  viewer_certificate {
    acm_certificate_arn            = aws_acm_certificate.cert.arn
    ssl_support_method             = "sni-only"
    minimum_protocol_version       = "TLSv1.2_2019"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }
  
  logging_config {
    include_cookies = false
    bucket          = "${aws_s3_bucket.logs.bucket}.s3.amazonaws.com"
    prefix          = "thetiptop"
  }
  
}

resource "aws_s3_bucket_policy" "bucket_app_policy" {
  bucket = aws_s3_bucket.app.id
  
  policy = jsonencode({
    Version = "2012-10-17",
    Id      = "AllowGetObjects",
    Statement = [
      {
        Sid       = "HTTPSOnly"
        Effect    = "Deny"
        Principal = "*"
        Action    = "s3:*"
        Resource = [
          aws_s3_bucket.app.arn,
          "${aws_s3_bucket.app.arn}/*",
        ]
        Condition = {
          Bool = {
            "aws:SecureTransport" = "false"
          }
        }
      }
    ]
  })
}

resource "aws_s3_bucket_website_configuration" "project_website" {
  bucket = aws_s3_bucket.app.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}

resource "aws_route53_record" "kodmain_cloudfront" {
  zone_id = "Z10052173VRSYMBUSS942"  # Remplacez par l'ID de votre zone Route 53
  name    = "kodmain.run"
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.s3_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.s3_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}