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
  region  = "us-east-1"
  alias   = "global"
  profile = "kodmain"
}

provider "aws" {
  region  = "eu-west-3"
  profile = "kodmain"
}

resource "aws_acm_certificate" "cert" {
  provider          = aws.global
  domain_name       = "kodmain.run"
  validation_method = "DNS"
  subject_alternative_names = ["staging.kodmain.run"] 
  tags = {
    Environment = "production"
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_s3_bucket_policy" "bucket_policy" {
  bucket = aws_s3_bucket.app.id  // Assurez-vous que cela correspond au nom de votre bucket S3

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action    = ["s3:GetObject"]
        Effect    = "Allow"
        Principal = "*"
        Resource  = ["arn:aws:s3:::${aws_s3_bucket.app.bucket}/*"]
      }
    ]
  })
}


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

resource "aws_cloudfront_function" "path_rewrite_function" {
  name    = "PathRewriteFunction"
  runtime = "cloudfront-js-1.0"

  publish = true

  code = <<-EOT
    function handler(event) {
    var request = event.request;
    var uri = request.uri;
    var headers = request.headers;
    var host = headers.host.value;

    // Rediriger en fonction du domaine
    if (host === 'staging.kodmain.run') {
        uri = '/staging' + uri;
    } else if (host === 'kodmain.run') {
        uri = '/production' + uri;
    }

    // Rediriger vers un document par défaut si nécessaire
    if (uri.endsWith('/') || !/\.[^/]+$/.test(uri)) {
        uri = uri.endsWith('/') ? uri + 'index.html' : uri + '/index.html';
    }

    request.uri = uri;
    return request;
}

  EOT
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

  aliases = ["kodmain.run","staging.kodmain.run"]

  default_cache_behavior {
    function_association {
      event_type = "viewer-request"
      function_arn = aws_cloudfront_function.path_rewrite_function.arn
    }
    
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = "S3-${aws_s3_bucket.app.id}"
    compress         = true

    forwarded_values {
      query_string = true
      headers      = ["Origin"]
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
  zone_id = "Z10052173VRSYMBUSS942"
  name    = "kodmain.run"
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.s3_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.s3_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}


resource "aws_route53_record" "staging_kodmain_cloudfront" {
  zone_id = "Z10052173VRSYMBUSS942"
  name    = "staging.kodmain.run"
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.s3_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.s3_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}