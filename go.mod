module github.com/philips-labs/terraform-backend-http

go 1.15

require (
	github.com/bhoriuchi/go-crypto v0.0.0-20190614232206-6aed78a5c061
	github.com/bhoriuchi/terraform-backend-http v0.0.0-20190615070304-ad22a976cbe3
	github.com/cloudfoundry-community/gautocloud v1.1.6
	github.com/cloudfoundry-community/go-cfenv v1.18.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/mock v1.4.4 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/minio/md5-simd v1.1.1 // indirect
	github.com/minio/minio-go/v7 v7.0.6
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/philips-software/gautocloud-connectors v0.4.0
	github.com/spf13/afero v1.4.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	go.mongodb.org/mongo-driver v1.4.4
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/cloudfoundry-community/gautocloud v1.1.6 => github.com/loafoe/gautocloud v0.0.0-20201207124432-b51ec5b81955
