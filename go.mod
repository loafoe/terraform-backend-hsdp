module github.com/philips-labs/terraform-backend-hsdp

go 1.15

require (
	github.com/bhoriuchi/go-crypto v0.0.0-20190614232206-6aed78a5c061
	github.com/cloudfoundry-community/gautocloud v1.1.10
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/minio/md5-simd v1.1.1 // indirect
	github.com/minio/minio-go/v7 v7.0.32
	github.com/philips-software/gautocloud-connectors v0.6.0
	github.com/philips-software/go-hsdp-api v0.75.6
	github.com/spf13/viper v1.12.0
)

replace github.com/cloudfoundry-community/gautocloud v1.1.6 => github.com/loafoe/gautocloud v0.0.0-20201207124432-b51ec5b81955
