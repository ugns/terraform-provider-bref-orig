terraform {
  required_providers {
    bref = {
      version = "0.1"
      source  = "bref.sh/bref/bref"
    }
  }
}

provider "bref" {}

provider "bref" {
  alias        = "us-west-2"
  region       = "us-west-2"
}

provider "bref" {
  alias        = "eu-central-1"
  region       = "eu-central-1"
  bref_version = "1.0.2"
}

data "bref_lambda_layer" "php7" {
  provider = bref.us-west-2
  layer_name = "php-74"
}

data "bref_lambda_layer" "php8" {
  layer_name = "php-80-fpm"
}

data "bref_lambda_layer" "console" {
  provider   = bref.eu-central-1
  layer_name = "console"
}

output "php7_layer" {
  value = data.bref_lambda_layer.php7.layer_arn
}

output "php8_layer" {
  value = data.bref_lambda_layer.php8.arn
}

output "console_layer" {
  value = data.bref_lambda_layer.console.arn
}
