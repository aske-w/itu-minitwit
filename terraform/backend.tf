terraform {
  backend "s3" {
    region = "us-west-1"
    endpoint = "https://fra1.digitaloceanspaces.com"
    skip_credentials_validation = true
    skip_metadata_api_check = true
  }
}