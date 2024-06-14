package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

func FindLoadBalancer(config aws.Config, region string, searchValue string) ([][]string, error) {
	config.Region = region

	elbv2Client := elasticloadbalancingv2.NewFromConfig(config)
	v2_output, err := elbv2Client.DescribeLoadBalancers(context.TODO(), &elasticloadbalancingv2.DescribeLoadBalancersInput{})
	if err != nil {
		var accessDeniedErr *http.ResponseError
		if errors.As(err, &accessDeniedErr) && accessDeniedErr.HTTPStatusCode() == 403 {
			return nil, nil
		}
		fmt.Printf("Unable to list load balancers (v2), %v", err)
	}

	elbv1Client := elasticloadbalancing.NewFromConfig(config)
	v1_output, err := elbv1Client.DescribeLoadBalancers(context.TODO(), &elasticloadbalancing.DescribeLoadBalancersInput{})
	if err != nil {
		var accessDeniedErr *http.ResponseError
		if errors.As(err, &accessDeniedErr) && accessDeniedErr.HTTPStatusCode() == 403 {
			return nil, nil
		}
		fmt.Printf("Unable to list load balancers (v1), %v", err)
	}

	filteredLoadBalancers := [][]string{}
	if v2_output != nil {
		for _, lb := range v2_output.LoadBalancers {
			if strings.Contains(*lb.LoadBalancerArn, searchValue) || strings.Contains(*lb.DNSName, searchValue) || strings.Contains(*lb.LoadBalancerArn, searchValue) {
				filteredLoadBalancers = append(filteredLoadBalancers, []string{*lb.LoadBalancerName, *lb.DNSName})
				fmt.Println(*lb.LoadBalancerName)
			}
		}
	}
	if v1_output != nil {
		for _, lb := range v1_output.LoadBalancerDescriptions {
			if strings.Contains(*lb.LoadBalancerName, searchValue) || strings.Contains(*lb.DNSName, searchValue) {
				filteredLoadBalancers = append(filteredLoadBalancers, []string{*lb.LoadBalancerName, *lb.DNSName})
				fmt.Println(*lb.LoadBalancerName)
			}
		}
	}

	return filteredLoadBalancers, nil
}
