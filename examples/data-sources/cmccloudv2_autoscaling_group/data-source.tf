# get as group by id
data "autoscaling_group" "example" {
  autoscaling_group_id    = "your_vpc_id" 
}

# get as group by name
data "autoscaling_group" "example" {
  name    = "as-group-xk8a" 
}