module github.com/loafoe/terraform-backend-hsdp

go 1.24.1

toolchain go1.24.2

require (
	github.com/bhoriuchi/go-crypto v0.0.0-20190614232206-6aed78a5c061
	github.com/cloudfoundry-community/gautocloud v1.6.0
	github.com/dip-software/gautocloud-connectors v0.9.0
	github.com/dip-software/go-dip-api v0.91.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/minio/minio-go/v7 v7.0.90
	github.com/spf13/viper v1.20.1
)

require (
	github.com/aws/aws-sdk-go v1.55.6 // indirect
	github.com/azer/snakecase v1.0.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudfoundry-community/go-cfenv v1.18.0 // indirect
	github.com/coder/websocket v1.8.12 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-jose/go-jose/v4 v4.0.5 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.24.0 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/go-viper/encoding/hcl v0.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.6 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/vault/api v1.16.0 // indirect
	github.com/hasura/go-graphql-client v0.13.1 // indirect
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kevinburke/go-types v0.0.0-20180804230501-141e48b46866 // indirect
	github.com/kevinburke/go.uuid v1.2.0 // indirect
	github.com/kevinburke/rest v0.0.0-20180424040019-54a8b1109595 // indirect
	github.com/kevinburke/twilio-go v0.0.0-20180929160738-8c4607665b91 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/minio/crc64nvme v1.0.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/segmentio/kafka-go v0.4.47 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/ttacon/builder v0.0.0-20170518171403-c099f663e1c2 // indirect
	github.com/ttacon/libphonenumber v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/oauth2 v0.28.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/cloudfoundry-community/gautocloud v1.1.6 => github.com/loafoe/gautocloud v0.0.0-20201207124432-b51ec5b81955
