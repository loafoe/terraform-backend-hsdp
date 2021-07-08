module github.com/philips-labs/terraform-backend-hsdp

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.28 // indirect
	github.com/bhoriuchi/go-crypto v0.0.0-20190614232206-6aed78a5c061
	github.com/cloudfoundry-community/gautocloud v1.1.7
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/minio/md5-simd v1.1.1 // indirect
	github.com/minio/minio-go/v7 v7.0.12
	github.com/philips-software/gautocloud-connectors v0.4.0
	github.com/philips-software/go-hsdp-api v0.26.1-0.20201208155559-7375257898d2
	github.com/spf13/viper v1.8.1
)

replace github.com/cloudfoundry-community/gautocloud v1.1.6 => github.com/loafoe/gautocloud v0.0.0-20201207124432-b51ec5b81955
