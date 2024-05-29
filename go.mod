module github.com/cmc-cloud/terraform-provider-cmccloudv2

go 1.22

toolchain go1.22.2

require (
	// run: go get -d github.com/cmc-cloud/gocmcapiv2@7c5f385 to get correct lastest version, eb9c186 = hash github commit
	github.com/cmc-cloud/gocmcapiv2 v0.0.0-20240529175245-7c5f385170dd
	github.com/hashicorp/terraform-plugin-sdk v1.17.2
)

require (
	cloud.google.com/go v0.65.0 // indirect
	cloud.google.com/go/bigquery v1.8.0 // indirect
	cloud.google.com/go/datastore v1.1.0 // indirect
	cloud.google.com/go/pubsub v1.3.1 // indirect
	cloud.google.com/go/storage v1.10.0 // indirect
	dmitri.shuralyov.com/gpu/mtl v0.0.0-20190408044501-666a987793e9 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802 // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/alcortesm/tgz v0.0.0-20161220082320-9c5fe88206d7 // indirect
	github.com/andybalholm/crlf v0.0.0-20171020200849-670099aa064f // indirect
	github.com/anmitsu/go-shlex v0.0.0-20161002113705-648efa622239 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/apparentlymart/go-dump v0.0.0-20190214190832-042adf3cf4a0 // indirect
	github.com/apparentlymart/go-textseg v1.0.0 // indirect
	github.com/apparentlymart/go-textseg/v12 v12.0.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5 // indirect
	github.com/aws/aws-sdk-go v1.37.0 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.2.1 // indirect
	github.com/cheggaaa/pb v1.0.27 // indirect
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/client9/misspell v0.3.4 // indirect
	github.com/cncf/udpa/go v0.0.0-20191209042840-269d4d468f6f // indirect
	github.com/creack/pty v1.1.9 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/envoyproxy/go-control-plane v0.9.4 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.1.0 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/gliderlabs/ssh v0.2.2 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.1.0 // indirect
	github.com/go-git/go-git-fixtures/v4 v4.0.2-0.20200613231340-f56387b50c12 // indirect
	github.com/go-git/go-git/v5 v5.3.0 // indirect
	github.com/go-gl/glfw v0.0.0-20190409004039-e6da0acd62b1 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20200222043503-6f7a984d4dc4 // indirect
	github.com/go-resty/resty/v2 v2.11.0 // indirect
	github.com/go-test/deep v1.0.3 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/mock v1.4.4 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/martian v2.1.0+incompatible // indirect
	github.com/google/martian/v3 v3.0.0 // indirect
	github.com/google/pprof v0.0.0-20200708004538-1a94d8640e99 // indirect
	github.com/google/renameio v0.1.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-checkpoint v0.5.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.5.3 // indirect
	github.com/hashicorp/go-hclog v0.9.2 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/hashicorp/go-plugin v1.3.0 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/hashicorp/go-version v1.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/hashicorp/hcl v0.0.0-20170504190234-a4b07c25de5f // indirect
	github.com/hashicorp/hcl/v2 v2.8.2 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-config-inspect v0.0.0-20191212124732-c6ae6269b9d7 // indirect
	github.com/hashicorp/terraform-exec v0.13.3 // indirect
	github.com/hashicorp/terraform-json v0.10.0 // indirect
	github.com/hashicorp/terraform-plugin-test/v2 v2.2.1 // indirect
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20181102032728-5e5cf60278f6 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jessevdk/go-flags v1.5.0 // indirect
	github.com/jhump/protoreflect v1.6.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jmespath/go-jmespath/internal/testify v1.5.1 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/keybase/go-crypto v0.0.0-20161004153544-93f5b35093ba // indirect
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/klauspost/compress v1.11.2 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/pty v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mattn/go-isatty v0.0.5 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mitchellh/cli v1.1.2 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.0.4 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/posener/complete v1.2.1 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/rogpeppe/go-internal v1.3.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.4.1 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/ulikunitz/xz v0.5.8 // indirect
	github.com/vmihailenco/msgpack v3.3.3+incompatible // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.1 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/yuin/goldmark v1.4.13 // indirect
	github.com/zclconf/go-cty v1.8.2 // indirect
	github.com/zclconf/go-cty-debug v0.0.0-20191215020915-b22d67c1ba0b // indirect
	github.com/zclconf/go-cty-yaml v1.0.2 // indirect
	go.opencensus.io v0.22.4 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mobile v0.0.0-20190719004257-d2bd2a29d028 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/term v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.34.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d // indirect
	google.golang.org/grpc v1.32.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.27 // indirect
	gopkg.in/errgo.v2 v2.1.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
	rsc.io/binaryregexp v0.2.0 // indirect
	rsc.io/quote/v3 v3.1.0 // indirect
	rsc.io/sampler v1.3.0 // indirect
)

// require (
// 	cloud.google.com/go/compute v1.24.0 // indirect
// 	cloud.google.com/go/compute/metadata v0.2.3 // indirect
// 	cloud.google.com/go/iam v1.1.6 // indirect
// 	github.com/Masterminds/goutils v1.1.1 // indirect
// 	github.com/Masterminds/semver v1.5.0 // indirect
// 	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
// 	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
// 	github.com/felixge/httpsnoop v1.0.4 // indirect
// 	github.com/go-logr/logr v1.4.1 // indirect
// 	github.com/go-logr/stdr v1.2.2 // indirect
// 	github.com/google/s2a-go v0.1.7 // indirect
// 	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
// 	github.com/huandu/xstrings v1.3.3 // indirect
// 	github.com/imdario/mergo v0.3.15 // indirect
// 	github.com/klauspost/compress v1.11.2 // indirect
// 	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
// 	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
// 	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.47.0 // indirect
// 	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.47.0 // indirect
// 	go.opentelemetry.io/otel v1.22.0 // indirect
// 	go.opentelemetry.io/otel/metric v1.22.0 // indirect
// 	go.opentelemetry.io/otel/trace v1.22.0 // indirect
// 	golang.org/x/sync v0.6.0 // indirect
// 	golang.org/x/time v0.5.0 // indirect
// 	google.golang.org/genproto/googleapis/api v0.0.0-20240227224415-6ceb2ff114de // indirect
// 	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
// )

// require (
// 	cloud.google.com/go v0.112.0 // indirect
// 	cloud.google.com/go/storage v1.36.0 // indirect
// 	github.com/agext/levenshtein v1.2.2 // indirect
// 	github.com/apparentlymart/go-cidr v1.1.0 // indirect
// 	github.com/armon/go-radix v1.0.0 // indirect
// 	github.com/aws/aws-sdk-go v1.37.0 // indirect
// 	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
// 	github.com/bgentry/speakeasy v0.1.0 // indirect
// 	github.com/davecgh/go-spew v1.1.1 // indirect
// 	github.com/fatih/color v1.16.0 // indirect
// 	github.com/go-resty/resty/v2 v2.11.0 // indirect
// 	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
// 	github.com/golang/protobuf v1.5.4 // indirect
// 	github.com/google/go-cmp v0.6.0 // indirect
// 	github.com/google/uuid v1.6.0 // indirect
// 	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
// 	github.com/hashicorp/errwrap v1.1.0 // indirect
// 	github.com/hashicorp/go-checkpoint v0.5.0 // indirect
// 	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
// 	github.com/hashicorp/go-getter v1.5.3 // indirect
// 	github.com/hashicorp/go-hclog v1.6.2 // indirect
// 	github.com/hashicorp/go-multierror v1.1.1 // indirect
// 	github.com/hashicorp/go-plugin v1.6.0 // indirect
// 	github.com/hashicorp/go-safetemp v1.0.0 // indirect
// 	github.com/hashicorp/go-uuid v1.0.3 // indirect
// 	github.com/hashicorp/go-version v1.6.0 // indirect
// 	github.com/hashicorp/hcl v0.0.0-20170504190234-a4b07c25de5f // indirect
// 	github.com/hashicorp/hcl/v2 v2.20.0 // indirect
// 	github.com/hashicorp/logutils v1.0.0 // indirect
// 	github.com/hashicorp/terraform-config-inspect v0.0.0-20191212124732-c6ae6269b9d7 // indirect
// 	github.com/hashicorp/terraform-exec v0.13.3 // indirect
// 	github.com/hashicorp/terraform-json v0.10.0 // indirect
// 	github.com/hashicorp/terraform-plugin-test/v2 v2.2.1 // indirect
// 	github.com/hashicorp/terraform-svchost v0.1.1 // indirect
// 	github.com/hashicorp/yamux v0.1.1 // indirect
// 	github.com/jmespath/go-jmespath v0.4.0 // indirect
// 	github.com/mattn/go-colorable v0.1.13 // indirect
// 	github.com/mattn/go-isatty v0.0.20 // indirect
// 	github.com/mitchellh/cli v1.1.2 // indirect
// 	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
// 	github.com/mitchellh/copystructure v1.2.0 // indirect
// 	github.com/mitchellh/go-homedir v1.1.0 // indirect
// 	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
// 	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
// 	github.com/mitchellh/mapstructure v1.5.0 // indirect
// 	github.com/mitchellh/reflectwalk v1.0.2 // indirect
// 	github.com/oklog/run v1.0.0 // indirect
// 	github.com/posener/complete v1.2.3 // indirect
// 	github.com/spf13/afero v1.2.2 // indirect
// 	github.com/ulikunitz/xz v0.5.8 // indirect
// 	github.com/zclconf/go-cty v1.14.4 // indirect
// 	github.com/zclconf/go-cty-yaml v1.0.2 // indirect
// 	go.opencensus.io v0.24.0 // indirect
// 	golang.org/x/crypto v0.21.0 // indirect
// 	golang.org/x/mod v0.16.0 // indirect
// 	golang.org/x/net v0.23.0 // indirect
// 	golang.org/x/oauth2 v0.17.0 // indirect
// 	golang.org/x/sys v0.18.0 // indirect
// 	golang.org/x/text v0.14.0 // indirect
// 	golang.org/x/tools v0.13.0 // indirect
// 	google.golang.org/api v0.162.0 // indirect
// 	google.golang.org/appengine v1.6.8 // indirect
// 	google.golang.org/genproto v0.0.0-20240227224415-6ceb2ff114de // indirect
// 	google.golang.org/grpc v1.63.2 // indirect
// 	google.golang.org/protobuf v1.34.0 // indirect
// )

// uncomment this line when build from code
// replace github.com/cmc-cloud/gocmcapiv2 => D:\code\CMC\openstack\terraform\gocmcapiv2
