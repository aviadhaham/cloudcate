package aws_search

type SearchResult struct {
	Account string `json:"account"`
	Profile string `json:"profile"`
	Region  string `json:"region"`
}

type LoadBalancerSearchResult struct {
	SearchResult
	LoadBalancerArn string `json:"load_balancer_arn"`
}

type S3SearchResult struct {
	SearchResult
	BucketName string `json:"bucket_name"`
}

type DNSSearchResult struct {
	SearchResult
	HostedZoneName string `json:"hosted_zone_name"`
	DnsRecordName  string `json:"dns_record_name"`
	DnsRecordType  string `json:"dns_record_type"`
}

type IamSearchResult struct {
	SearchResult
	UserName string `json:"user_name"`
}

type ElasticIpSearchResult struct {
	SearchResult
	PublicIp   string `json:"public_ip"`
	InstanceId string `json:"instance_id"`
}
