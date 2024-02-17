---

# CloudCate

Welcome to CloudCate! 🚀 This handy tool tackles the headache of searching for AWS resources across multiple accounts. Ever found yourself juggling between AWS accounts to track down resources? CloudCate is here to make your life easier by consolidating searches for various AWS resources into one simple, user-friendly interface.

## What You Can Search

Right out of the box, CloudCate lets you search across these AWS resource types in multiple accounts:
- S3 Buckets
- DNS (Hosted Zones or Records)
- Load Balancers
- EC2 Instances
- IAM Access Keys
- Elastic IPs

## Quick Start

### Prerequisites

Before diving in, make sure you have:
- Docker (for Docker users) or Go (for local runners)
- Your AWS CLI configured with `.aws/credentials` containing the profiles you want to search

### Important Note on AWS Access

It's crucial to ensure that the AWS access keys used with CloudCate have the necessary permissions to search the resources you're interested in. You're responsible for creating and managing these access keys safely. Make sure they're properly secured and have the right permissions set up across all accounts you plan to search.

### Run It Locally

To get CloudCate up and running on your machine, execute:

```bash
go run cmd/main.go
```

### Run It with Docker

For a containerized experience, build and run CloudCate using the following commands:

1. **Build the Docker Image**

```bash
docker build -t cloudcate .
```

2. **Run the Docker Container**

Remember to replace `/path/to/credentials` with the actual path to your `.aws/credentials` file.

```bash
docker run --rm -d -p 8080:80 -v /path/to/credentials:/root/.aws/credentials cloudcate:latest
```

CloudCate will now be accessible at `http://localhost:8080`.

## How to Use It

Visit `http://localhost:8080` in your browser. The UI is straightforward: select the AWS service you're searching for (e.g., S3, EC2) and input your search terms. CloudCate will search through the specified AWS profiles and regions, showing you the resources that match your query.


## License

CloudCate is released under the [MIT license](https://choosealicense.com/licenses/mit/).
