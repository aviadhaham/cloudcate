package search

type SearchResult struct {
	Account string `json:"account"`
	Profile string `json:"profile"`
}

type SearchResultNonGlobal struct {
	SearchResult
	Region string `json:"region"`
}

type VpcSearchResult struct {
	SearchResultNonGlobal
	VpcName   string `json:"vpc_name"`
	VpcId     string `json:"vpc_id"`
	CidrBlock string `json:"cidr_block"`
}

type Ec2SearchResult struct {
	SearchResultNonGlobal
	InstanceId       string `json:"instance_id"`
	InstanceName     string `json:"instance_name"`
	PrivateIpAddress string `json:"private_ip_address"`
	PrivateDnsName   string `json:"private_dns_name"`
	PublicDnsName    string `json:"public_dns_name"`
	PublicIpAddress  string `json:"public_ip_address"`
}

type LoadBalancerSearchResult struct {
	SearchResultNonGlobal
	LoadBalancerName    string `json:"load_balancer_name"`
	LoadBalancerDnsName string `json:"load_balancer_dns_name"`
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

type IamUserSearchResult struct {
	SearchResult
	UserName string `json:"user_name"`
}

type IamUserKeySearchResult struct {
	IamUserSearchResult
	AccessKey string `json:"access_key"`
}

type ElasticIpSearchResult struct {
	SearchResultNonGlobal
	PublicIp   string `json:"public_ip"`
	InstanceId string `json:"instance_id"`
}

type CloudfrontSearchResult struct {
	SearchResult
	DistributionArn string `json:"distribution_arn"`
	DistributionId  string `json:"distribution_id"`
	DomainName      string `json:"domain_name"`
}
