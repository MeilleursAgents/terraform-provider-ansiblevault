module github.com/MeilleursAgents/terraform-provider-ansiblevault/v2

go 1.12

require (
	github.com/hashicorp/terraform v0.14.9
	github.com/sosedoff/ansible-vault-go v0.0.0-20181205202858-ab5632c40bf5
	gopkg.in/yaml.v2 v2.4.0
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
