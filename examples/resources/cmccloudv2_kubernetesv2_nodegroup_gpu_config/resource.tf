resource "cmccloudv2_kubernetesv2_nodegroup_gpu_config" "nodegroup_gpu_1" {
    cluster_id 		= "93cf35f9-f570-4a6c-aeb1-115a0de256c9"
	nodegroup_id 	= "db478032-0d65-4f12-8d43-bdfeeab879f8"
	gpu_model 		= "GPU-Accelerated-1080Ti"
	driver 			= "580.65.06"
	strategy 		= "single" 
	mig_profile 	= "all-1g.5gb" 
	gpu_profiles { 
		name 		= "nvidia.com/gpu"
		replicas 	= 2
	}
	gpu_profiles {
		name 		= "nvidia.com/mig-1g.5gb"
		replicas 	= 2
	} 	
}