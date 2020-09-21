module cadence-example

go 1.14

require (
	github.com/pborman/uuid v0.0.0-20160209185913-a97ce2ca70fa
	github.com/spf13/viper v1.7.1
	github.com/uber-go/tally v3.3.17+incompatible
	go.uber.org/cadence v0.13.4
	go.uber.org/yarpc v1.47.2
	go.uber.org/zap v1.16.0
	gopkg.in/neurosnap/sentences.v1 v1.0.6
)

replace github.com/apache/thrift => github.com/apache/thrift v0.0.0-20190309152529-a9b748bb0e02
