module github.com/MeilleursAgents/terraform-provider-ansiblevault/v2

go 1.12

require (
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/terraform v0.13.4
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/sosedoff/ansible-vault-go v0.0.0-20181205202858-ab5632c40bf5
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/client-go v11.0.0+incompatible // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
