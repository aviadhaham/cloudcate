---

# CloudCate

Hey there! 👋 CloudCate is a neat little tool that helps you find AWS resources like EC2 instances, S3 buckets, and more across different profiles and regions. It's built in Go and runs both as a simple local server or in a Docker container. Just plug in your AWS credentials, and you're good to go!

## Quick Start

### What You Need

- Docker (if you're going the container route)
- Go (for running locally)
- Your AWS credentials set up in `.aws/credentials`

### Run It Locally

Got Go installed? Run this in the project's root:

```bash
go run cmd/main.go
```

### Run It with Docker

Prefer Docker? Build and run CloudCate with these commands:

1. Build the image:

```bash
docker build -t cloudcate .
```

2. Run the container (don't forget to replace `/path/to/credentials` with your actual credentials path):

```bash
docker run --rm -d -p 8080:80 -v /path/to/credentials:/root/.aws/credentials cloudcate:latest
```

Now, CloudCate will be up at `http://localhost:8080`.

## How to Use It

Once CloudCate is running, head over to `http://localhost:8080`. You'll see a simple UI where you can pick the AWS service you're interested in (like EC2 or S3) and type in what you're looking for. Hit search, and it'll comb through your profiles and regions to find the resources matching your query.

## License

This project is under the [MIT license](https://choosealicense.com/licenses/mit/).
