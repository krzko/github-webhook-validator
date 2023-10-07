# GitHub Webhook Validator

A simple Go utility to validate GitHub webhook deliveries against a provided secret.

## Overview

This utility sets up a web server that listens for `POST` requests on the `/webhook` endpoint and validates the signature of the incoming GitHub webhook payload against a predefined secret.


## Usage

Clone the repository:

```sh
git clone https://github.com/your-username/github-webhook-validator.git
cd github-webhook-validator
go build
```

### Service

Run:

```sh
 export WEBHOOK_SECRET="It's a Secret to Everybody"

 ./github-webhook-validator
```

### Client

To test the GitHub webhook validator utility using `curl`, you'll need to mimic a webhook request from GitHub. This means you'll need to create a `POST` request with a `X-Hub-Signature-256 header and a request body. The `X-Hub-Signature-256 header value should be the HMAC hex digest of the request body using your webhook secret.

Here's a simplified example of how you might construct such a curl command:

1. First, you'll need to compute the `X-Hub-Signature-256` header value. You can use a command-line tool like `openssl` to do this:

```sh
echo -n "your-payload-here" | openssl dgst -sha256 -hmac $WEBHOOK_SECRET
```

Replace `"your-payload-here"` with the content of your webhook payload, and `$WEBHOOK_SECRET` with your webhook secret. This will output a hash value.

2. Now, you can use curl to send a POST request to your utility, including the `X-Hub-Signature-256` header value you computed in the previous step:

```sh
curl -X POST \
     -H "Content-Type: application/json" \
     -H "X-Hub-Signature-256: sha256=$WEBHOOK_SECRET \
     --data '{"payload": "your-payload-here"}' \
     http://localhost:8080/webhook
```
