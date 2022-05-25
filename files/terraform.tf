resource "aws_s3_bucket" "collector-data-bucket" {
  bucket         = "collector-data-bucket"
  force_destroy  = false
  acl            = "private"
}

resource "aws_s3_bucket" "athena-querydata-bucket" {
  bucket         = "athena-querydata-bucket"
  force_destroy  = false
  acl            = "private"
}