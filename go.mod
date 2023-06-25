module controller

go 1.18

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	github.com/dustin/go-humanize v1.0.0
	github.com/filecoin-project/go-statemachine v1.0.3
	github.com/filecoin-project/go-statestore v0.1.0
	github.com/fsnotify/fsnotify v1.6.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/ipfs/go-datastore v0.6.0
	github.com/ipfs/go-log/v2 v2.5.1
	github.com/jinzhu/gorm v1.9.16
	github.com/klauspost/compress v1.15.11
	github.com/lestrrat-go/strftime v1.0.6
	github.com/libp2p/go-libp2p v0.23.2
	github.com/minio/minio-go/v7 v7.0.42
	github.com/minio/pkg v1.5.4
	github.com/minio/sio v0.3.0
	github.com/multiformats/go-multiaddr v0.7.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/viper v1.13.0
	github.com/stretchr/testify v1.8.0
	github.com/valyala/tcplisten v1.0.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20200123233031-1cdf64d27158
	go.opencensus.io v0.23.0
	go.uber.org/zap v1.23.0
	go4.org v0.0.0-20180809161055-417644f6feb5
	golang.org/x/sys v0.1.0
	golang.org/x/xerrors v0.0.0-20220517211312-f3a8303e98df
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gotest.tools v2.2.0+incompatible
)

replace github.com/ipfs/go-datastore v0.6.0 => github.com/ipfs/go-datastore v0.1.1

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/containerd/cgroups v1.0.4 // indirect
	github.com/coreos/go-systemd/v22 v22.4.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/davidlazar/go-crypto v0.0.0-20200604182044-b73af7476f6c // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/elastic/gosigar v0.14.2 // indirect
	github.com/filecoin-project/go-cbor-util v0.0.0-20191219014500-08c40a1e63a2 // indirect
	github.com/flynn/noise v1.0.0 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huin/goupnp v1.0.3 // indirect
	github.com/iancoleman/orderedmap v0.2.0 // indirect
	github.com/ipfs/go-block-format v0.0.2 // indirect
	github.com/ipfs/go-cid v0.3.2 // indirect
	github.com/ipfs/go-ipfs-util v0.0.2 // indirect
	github.com/ipfs/go-ipld-cbor v0.0.3 // indirect
	github.com/ipfs/go-ipld-format v0.0.2 // indirect
	github.com/ipfs/go-log v1.0.0 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/jbenet/go-temp-err-catcher v0.1.0 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.1.2 // indirect
	github.com/koron/go-ssdp v0.0.3 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/libp2p/go-cidranger v1.1.0 // indirect
	github.com/libp2p/go-flow-metrics v0.1.0 // indirect
	github.com/libp2p/go-libp2p-asn-util v0.2.0 // indirect
	github.com/libp2p/go-msgio v0.2.0 // indirect
	github.com/libp2p/go-nat v0.1.0 // indirect
	github.com/libp2p/go-netroute v0.2.0 // indirect
	github.com/libp2p/go-openssl v0.1.0 // indirect
	github.com/libp2p/go-reuseport v0.2.0 // indirect
	github.com/libp2p/go-yamux/v4 v4.0.0 // indirect
	github.com/lucas-clemente/quic-go v0.29.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/marten-seemann/qtls-go1-18 v0.1.2 // indirect
	github.com/marten-seemann/qtls-go1-19 v0.1.0 // indirect
	github.com/marten-seemann/tcp v0.0.0-20210406111302-dfbc87cc63fd // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/mikioh/tcpinfo v0.0.0-20190314235526-30a79bb1804b // indirect
	github.com/mikioh/tcpopt v0.0.0-20190314235656-172688c1accc // indirect
	github.com/minio/madmin-go v1.6.6 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multiaddr-dns v0.3.1 // indirect
	github.com/multiformats/go-multiaddr-fmt v0.1.0 // indirect
	github.com/multiformats/go-multibase v0.1.1 // indirect
	github.com/multiformats/go-multicodec v0.6.0 // indirect
	github.com/multiformats/go-multihash v0.2.1 // indirect
	github.com/multiformats/go-multistream v0.3.3 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/opencontainers/runtime-spec v1.0.2 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/polydawn/refmt v0.0.0-20190809202753-05966cbd336a // indirect
	github.com/prometheus/client_golang v1.13.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/raulk/go-watchdog v1.3.0 // indirect
	github.com/spacemonkeygo/spacelog v0.0.0-20180420211403-2296661a0572 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20221012134737-56aed061732a // indirect
	golang.org/x/exp v0.0.0-20220916125017-b168a2c6b86b // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20221017152216-f25eb7ecb193 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
	google.golang.org/api v0.81.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
)
