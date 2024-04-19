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
- CloudFront Distributions

## Quick Start

### Prerequisites

Before diving in, make sure you have:
- Docker (for Docker users) or Go (for local runners)
- Your AWS CLI configured with `.aws/credentials` containing the profiles you want to search

### Important Note on AWS Access

It's crucial to ensure that the AWS access keys used with CloudCate have the necessary permissions to search the resources you're interested in. You're responsible for creating and managing these access keys safely. Make sure they're properly secured and have the right permissions set up across all accounts you plan to search.

#### Example AWS IAM Policy that should support all the operations:
[`aws-policy.json`](aws-policy.json)

### Run It Locally

Using [Taskfile.yml](./Taskfile.yml)

#### Start the UI
```bash
task ui:dev
```
#### Start the backend server
```bash
task server:dev
```

**Server would be availble at `http://localhost:5173`**

### Build for Production

#### Build the UI
```bash
task ui:build
```
#### Build the backend server
```bash
task server:build
```

To start the production server (this will automatically build the UI and server if they haven't been built yet):

```bash
task server
```

**Server would be availble at `http://localhost`**

### Run It with Docker

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

#### Example `docker-compose.yml` to run with Docker Compose
[`docker-compose.yml`](docker-compose.yml)


## How to Use It

Select the AWS service you're searching for (e.g., S3, EC2) and input your search terms. CloudCate will search through the specified AWS profiles and regions, showing you the resources that match your query.


## License

CloudCate is released under the [MIT license](https://choosealicense.com/licenses/mit/).
