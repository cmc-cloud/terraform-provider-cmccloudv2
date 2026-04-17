# get opensearch dashboard flavor by name 
data "cmccloudv2_opensearch_dashboard_flavor" "flavor_1" {
    name = "dash.small"  
}

# get opensearch dashboard flavor by id
data "cmccloudv2_opensearch_dashboard_flavor" "flavor_2" {
    flavor_id = "ecd1971a-e43a-47d6-8c33-5ede45b458fb" 
}

# get opensearch dashboard flavor by CPU & RAM 
data "cmccloudv2_opensearch_dashboard_flavor" "flavor_3" {
    cpu = 1
    ram = 2
}