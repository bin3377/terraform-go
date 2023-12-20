data "aws_vpcs" "foo" {
}

output "foo" {
  value = data.aws_vpcs.foo.ids
}
