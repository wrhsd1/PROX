# ChatGPT Proxy Server

## Building
- Install [Go](https://go.dev/)
- `git clone https://github.com/acheong08/ChatGPT-Proxy-V1`
- `cd ChatGPT-Proxy-V1`
- `go build`

The output should be a binary `ChatGPT-Proxy-V1`

## Running

To run the binary, simply execute it with `./ChatGPT-Proxy-V1`

## Deployment

There are multiple ways to deploy this proxy: `fly.io`, `docker`, or simply running it on a server.

### fly.io

- Install [flyctl](https://fly.io/docs/getting-started/installing-flyctl/)
- Login to flyctl with `flyctl auth login`
- `flyctl apps create` and enter the appropriate details
- `flyctl deploy --remote-only` to deploy the app to fly.io

### Docker

- Install [Docker](https://docs.docker.com/get-docker/)
- Build the docker image with `docker build -t chatgpt-proxy-v1 .`
- Remember to build the binary first

### Running on a server

- Just run the binary on a server of the same architecture as the machine you built it on

### Configuration

The only configuration is the port to run the proxy on. This can be set with the `PORT` environment variable. The default port is `8080`.

You can use `nginx` to reverse proxy the port to a domain and/or Cloudflare to proxy the port to a domain.