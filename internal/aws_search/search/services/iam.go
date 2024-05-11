package aws_search

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func FindIam(config aws.Config, region string, searchValue string) (string, string, error) {

	config.Region = region

	iamClient := iam.NewFromConfig(config)

	accessKey := searchValue
	users, err := iamClient.ListUsers(context.TODO(), &iam.ListUsersInput{})
	if err != nil {
		return "", "", fmt.Errorf("failed to list IAM users: %v", err)
	}
	for _, user := range users.Users {
		accessKeys, err := iamClient.ListAccessKeys(context.TODO(), &iam.ListAccessKeysInput{
			UserName: user.UserName,
		})
		if err != nil {
			return "", "", fmt.Errorf("failed to list access keys for user %s: %v", *user.UserName, err)
		}
		for _, accessKeyMetadata := range accessKeys.AccessKeyMetadata {
			if strings.Contains(*accessKeyMetadata.AccessKeyId, accessKey) {
				return accessKey, *user.UserName, nil
			}
		}
	}
	return "", "", nil
}
