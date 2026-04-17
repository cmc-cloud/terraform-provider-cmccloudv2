# get opensearch flavor by name 
data "cmccloudv2_opensearch_flavor" "flavor_1" {
    name = "os.small.1"  
}

# get opensearch flavor by id
data "cmccloudv2_opensearch_flavor" "flavor_2" {
    flavor_id = "88705daf-1119-419d-874c-0c894704d37c" 
}

# get opensearch flavor by CPU & RAM 
data "cmccloudv2_opensearch_flavor" "flavor_3" {
    cpu = 2
    ram = 4
}