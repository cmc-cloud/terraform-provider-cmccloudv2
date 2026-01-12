# get security_group by id
data "cmccloudv2_security_group" "security_group_1" {
    security_group_id = "3ac75ff6-9744-4cc1-a21c-3bd4d6154bcb"
}

# get security_group by name
data "cmccloudv2_security_group" "security_group_2" {
    name = "sg-9wjd"
}