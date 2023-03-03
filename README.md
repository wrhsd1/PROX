# ChatGPT Proxy Server

### Note: The configuration file is not shared in this repository.

To run the proxy server using Docker, you can use the following command:

`make run-docker`

Alternatively, you can run the following command to start the Docker container manually:

```bash
docker run -itd --name chatgpt-proxy-v2 -p 10101:10101 acheong08/chatgpt-proxy-v2
```
This will start the container in detached mode with the name `chatgpt-proxy-v2` and map port `10101` on the host to port `10101` in the container.

Once the container is running, you can send requests to the proxy server at `http://localhost:10101`.
