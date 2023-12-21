variable "foo_in" {
  type = string
}

variable "foo_in_list" {
  type = list(string)
}

variable "foo_in_map" {
  type = map(string)
}

output "foo_out" {
  value = var.foo_in
}

output "foo_out_list" {
  value = var.foo_in_list
}

output "foo_out_map" {
  value = var.foo_in_map
}
