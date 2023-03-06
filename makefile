build:
    docker build -t chatgpt-proxy-v2 .

run-docker:
    docker run -itd --name chatgpt-proxy-v2 -p 10101:10101 chatgpt-proxy-v2
