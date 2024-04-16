export type SearchResult = {
  Account: string;
  Profile: string;
  Region: string;
};

export type Ec2SearchResult = SearchResult &{
  instanceId: string;
  InstanceName: string;
  PrivateIpAddress: string;
  PrivateDnsName: string;
  PublicIpAddress: string;
  PublicDnsName: string;
};

export type S3SearchResult = SearchResult &{
  BucketName: string;
};

export type LoadBalancerSearchResult = SearchResult &{
  LoadBalancerName: string;
  LoadBalancerDnsName: string;
};

export type DNSSearchResult = SearchResult &{
  HostedZoneName: string;
  DnsRecordName: string;
  DnsRecordType: string;
};

export type IamSearchResult = SearchResult &{
  UserName: string;
};

export type ElasticIpSearchResult = SearchResult &{
  PublicIp: string;
  InstanceId: string;
};

export type CloudfrontSearchResult = SearchResult &{
  DistributionArn: string;
  DistributionId: string;
  DomainName: string;
};

export type AllSearchResults = Ec2SearchResult | S3SearchResult | LoadBalancerSearchResult | DNSSearchResult | IamSearchResult | ElasticIpSearchResult | CloudfrontSearchResult;
