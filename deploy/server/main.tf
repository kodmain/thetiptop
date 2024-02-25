provider "aws" {
  region = "eu-west-3"
  profile = "kodmain"
}

variable "github_token" {
  description = "GitHub token"
}

locals {
  nomad_server              = base64encode(file("${path.module}/nomad-server.hcl"))
  nomad_policy              = base64encode(file("${path.module}/nomad-policy.hcl"))
  nomad_service             = base64encode(file("${path.module}/nomad.service"))
  service_admin             = base64encode(file("${path.module}/../jobs/admin/admin.hcl"))
  service_project           = base64encode(file("${path.module}/../jobs/project/project.hcl"))
  service_monitoring        = base64encode(file("${path.module}/../jobs/monitoring/monitoring.hcl"))
}

resource "aws_instance" "free_tier_arm_instance" {
  # X86 ami-089c89a80285075f7 t2.micro  Amazon Linux 2 # WORKING
  # ARM ami-0ddd50b03e7b395c4 t4g.micro Amazon Linux 2 # NOT WORKING (fixed by client.cpu_total_compute in nomad-server.hcl)
  ami           = "ami-0ddd50b03e7b395c4"
  instance_type = "t4g.micro"

  tags = {
    Name = "NomadServer"
  }
  
  user_data = <<-EOF
    #!/bin/bash
    sudo yum install -y yum-utils
    sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/AmazonLinux/hashicorp.repo
    sudo yum-config-manager --add-repo https://cli.github.com/packages/rpm/gh-cli.repo
    sudo yum -y install nomad docker gh cni-plugins httpd-tools
    mkdir -p /home/ec2-user/services
    echo "${local.nomad_server}"       | base64 --decode > /home/ec2-user/nomad-server.hcl
    echo "${local.nomad_policy}"       | base64 --decode > /home/ec2-user/nomad-policy.hcl
    echo "${local.nomad_service}"      | base64 --decode > /etc/systemd/system/nomad.service
    echo "${local.service_monitoring}" | base64 --decode > /home/ec2-user/services/monitoring.hcl
    echo "${local.service_admin}"      | base64 --decode > /home/ec2-user/services/admin.hcl
    echo "${local.service_project}"    | base64 --decode > /home/ec2-user/services/project.hcl
    systemctl enable nomad 
    systemctl enable docker
    systemctl start docker
    systemctl start nomad
    nomad acl bootstrap > /home/ec2-user/bootstrap.token
    export NOMAD_TOKEN=$(cat /home/ec2-user/bootstrap.token | grep "Secret" |awk '{print $4}')
    echo "export NOMAD_TOKEN=$NOMAD_TOKEN" >> /home/ec2-user/.bashrc
    echo "export GH_TOKEN=${var.github_token}" >> /home/ec2-user/.bashrc
    nomad acl policy apply -description "Deployment" deploy /home/ec2-user/nomad-policy.hcl
    nomad acl token create -name="github" -policy="deploy" > /home/ec2-user/github.token
    export GITHUB_NOMAD_TOKEN=$(cat /home/ec2-user/github.token | grep "Secret" |awk '{print $4}')
    gh secret set GITHUB_NOMAD_TOKEN -b"$GITHUB_NOMAD_TOKEN" --repo kodmain/thetiptop
    sed -i 's/NOMADTOKEN/'"$NOMAD_TOKEN"'/g' /home/ec2-user/services/admin.hcl
    nomad job run -token=$NOMAD_TOKEN /home/ec2-user/services/admin.hcl
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

  ingress {
    from_port   = 22
    to_port     = 22
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

resource "aws_route53_record" "kodmain" {
  zone_id = "Z10052173VRSYMBUSS942"

  name    = "kodmain.run"  # Nom de domaine à rediriger
  type    = "A"
  ttl     = 300
  records = [aws_instance.free_tier_arm_instance.public_ip]

  allow_overwrite = true
}

resource "aws_route53_record" "kodmain_wildcard" {
  zone_id = "Z10052173VRSYMBUSS942"

  name    = "*.kodmain.run"  # Enregistrement wildcard pour tous les sous-domaines
  type    = "A"
  ttl     = 300
  records = [aws_instance.free_tier_arm_instance.public_ip]

  allow_overwrite = true
}

