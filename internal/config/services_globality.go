package config

var ServicesGlobality = map[string]bool{
	"loadbalancer": false,
	"s3":           true,
	"dns":          true,
	"iam":          true,
	"elastic_ip":   false,
	"cloudfront":   true,
}
