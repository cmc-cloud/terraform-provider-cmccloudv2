resource "cmccloudv2_elb_l7policy_rule" "rule1" {
	l7policy_id  	= "5f59d19b-0e7b-42e0-a277-f7e1d13d6aad"
	type 			= "COOKIE"
	compare_type	= "CONTAINS"
	key				= "region-id"
	value			= "hn-1"
	invert			= false 
}
resource "cmccloudv2_elb_l7policy_rule" "rule2" {
	l7policy_id  	= "5f59d19b-0e7b-42e0-a277-f7e1d13d6aad"
	type 			= "FILE_TYPE"
	compare_type	= "EQUAL_TO" 
	value			= "png"
	invert			= false  
	depends_on 		= [cmccloudv2_elb_l7policy_rule.rule1]

}
resource "cmccloudv2_elb_l7policy_rule" "rule3" {
	l7policy_id  	= "5f59d19b-0e7b-42e0-a277-f7e1d13d6aad"
	type 			= "HEADER"
	compare_type	= "STARTS_WITH" 
	key				= "key"
	value			= "png"
	invert			= false
	depends_on 		= [cmccloudv2_elb_l7policy_rule.rule2]
}