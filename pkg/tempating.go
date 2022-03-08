package pkg 

const SINGLE_TEST = 
`---
test_name: test load balancer api
includes:
  - !include ../env_test.yaml
#######
# Post: create a full LB
#######
stages:
  - name: get lb
    delay_before: 5
    max_retries: 5
    request:
      url: "{url}/lb/v1/zones/{tavern.env_vars.ZONE}/lbs/{tavern.env_vars.LB_ID}"
      method: GET
      headers:
        X-Auth-Token: "{tavern.env_vars.API_CANARY_TOKEN}"
    response:
      status_code: 200
      save:
        json:
          lb_ip_id: ip[0].id`

func TemplateTest() {

}