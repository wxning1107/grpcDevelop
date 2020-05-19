module grpcClient

go 1.14

require (
	github.com/EDDYCJY/go-grpc-example v0.0.0-20181014074047-0f68708edbcb
	github.com/apache/thrift v0.13.0 // indirect
	github.com/docker/docker v1.13.1
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2 // indirect
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120 // indirect
	golang.org/x/sys v0.0.0-20200515095857-1151b9dac4a9 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200515170657-fc4c6c6a6587 // indirect
	google.golang.org/grpc v1.29.1
)

replace github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5 => github.com/openzipkin/zipkin-go-opentracing v0.2.2

replace github.com/openzipkin/zipkin-go-opentracing v0.2.2 => github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
