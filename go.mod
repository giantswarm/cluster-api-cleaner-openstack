module github.com/giantswarm/cluster-api-cleaner-openstack

go 1.13

require (
	github.com/giantswarm/microerror v0.4.0
	github.com/go-logr/logr v1.2.0
	github.com/gophercloud/gophercloud v0.16.0
	github.com/prometheus/client_golang v1.12.0 // indirect
	go.uber.org/zap v1.19.1
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	sigs.k8s.io/cluster-api v1.1.3
	sigs.k8s.io/cluster-api-provider-openstack v0.6.3
	sigs.k8s.io/controller-runtime v0.11.1
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.5.8
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.0 => github.com/gorilla/websocket v1.4.2
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.0.1-0.20211028151834-d72fd59c8483
)
