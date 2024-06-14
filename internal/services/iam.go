package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func FindIamUser(config aws.Config, region string, searchValue string) ([]string, error) {

	config.Region = region
	iamClient := iam.NewFromConfig(config)

	users, err := iamClient.ListUsers(context.TODO(), &iam.ListUsersInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list IAM users: %v", err)
	}
	filteredUsers := []string{}
	for _, user := range users.Users {
		if user.UserName != nil && strings.Contains(strings.ToLower(*user.UserName), searchValue) {
			filteredUsers = append(filteredUsers, *user.UserName)
		}
	}

	return filteredUsers, nil
}

func FindIamUserKey(config aws.Config, region string, searchValue string) (map[string]string, error) {

	config.Region = region
	iamClient := iam.NewFromConfig(config)

	users, err := iamClient.ListUsers(context.TODO(), &iam.ListUsersInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list IAM users: %v", err)
	}

	filteredAccessKeys := make(map[string]string)

	for _, user := range users.Users {
		accessKeys, err := iamClient.ListAccessKeys(context.TODO(), &iam.ListAccessKeysInput{
			UserName: user.UserName,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list access keys for user %s: %v", *user.UserName, err)
		}

		for _, accessKeyMetadata := range accessKeys.AccessKeyMetadata {
			if strings.Contains(*accessKeyMetadata.AccessKeyId, searchValue) {
				filteredAccessKeys[*user.UserName] = *accessKeyMetadata.AccessKeyId
			}
		}
	}

	return filteredAccessKeys, nil
}
