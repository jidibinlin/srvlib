module github.com/gzjjyz/srvlib

go 1.20

require (
	github.com/995933447/confloader v0.0.0-20230314141707-e7b191386ae2
	github.com/995933447/log-go v0.0.0-20230420123341-5d684963433b
	github.com/995933447/redisgroup v0.0.0-20230510085956-718f047520a1
	github.com/gorilla/websocket v1.5.0
	github.com/gzjjyz/micro v0.0.2
	github.com/huandu/go-clone v1.6.0
	github.com/json-iterator/go v1.1.12
	github.com/nats-io/nats.go v1.25.0
	google.golang.org/grpc v1.33.1
	gorm.io/driver/mysql v1.5.0
	gorm.io/gorm v1.25.1
)

require (
	github.com/995933447/std-go v0.0.0-20220806175833-ab3496c0b696 // indirect
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/etcd v3.3.27+incompatible // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20230327231512-ba87abf18a23 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/howeyc/fsnotify v0.9.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nats-server/v2 v2.9.16 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.43.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace (
	github.com/coreos/bbolt v1.3.7 => go.etcd.io/bbolt v1.3.7
	github.com/derekparker/delve => github.com/go-delve/delve v1.20.1
	github.com/go-delve/delve => github.com/derekparker/delve v1.4.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
