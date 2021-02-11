module github.com/tinkerbell/infrastructure

go 1.14

replace github.com/altoros/gosigma => github.com/juju/gosigma v0.0.0-20200420012028-063911838a9e

replace github.com/hashicorp/raft => github.com/juju/raft v2.0.0-20200420012049-88ad3b3f0a54+incompatible

replace gopkg.in/mgo.v2 => github.com/juju/mgo v2.0.0-20201106044211-d9585b0b0d28+incompatible

require (
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/juju/juju v0.0.0-20210211140050-bee2b099303a
	github.com/juju/packaging v0.0.0-20200421095529-970596d2622a
	github.com/pulumi/pulumi-equinix-metal/sdk v1.1.0
	github.com/pulumi/pulumi-ns1/sdk v1.4.0
	github.com/pulumi/pulumi/sdk/v2 v2.20.0
)
