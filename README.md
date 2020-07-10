# slack-payload-handler

[![Docker build](https://img.shields.io/docker/automated/bitsofinfo/slack-payload-handler)](https://hub.docker.com/repository/docker/bitsofinfo/slack-payload-handler)

Simple utility you can use as a custom [Tekton triggers webhook interceptor](https://github.com/tektoncd/triggers/blob/master/docs/eventlisteners.md#Webhook-Interceptors) when receiving [Slack interactive message payloads](https://api.slack.com/interactivity/handling#payloads) in response to user interaction (i.e. clicking on buttons etc) to trigger things in your Tekton CICD system.

For convience it lightly mutates the original JSON and adds an `action_values` property which is just an array of the selected `values` under `actions`. Completely ignorable if you don't want to use it.

## Usage

```
Usage of ./slack-payload-handler:
  -debug-request
        Optional, print requests to STDOUT, default false
  -debug-response
        Optional, print responses to STDOUT, default false
  -listen-port int
        Optional, port to listen on, default 8080 (default 8080)
```

# Docker example
```
docker run -p 8080:8080 -it bitsofinfo/slack-payload-handler \
  slack-payload-handler --debug-request true --debug-response true --listen-port 8080
```

# Example

```
$> ./slack-payload-handler --debug-request true --debug-response true

$> curl -k -X POST localhost:8080 --data-urlencode 'payload={"x":"b","actions":[{"value":"1"}]}'

{"action_values":["1"],"actions":[{"value":"1"}],"x":"b"}

```

# Test slack dummy payload

```
$> ./slack-payload-handler --debug-request true --debug-response true

$> curl -k -X POST localhost:8080 --data-urlencode "payload=$(cat slack.test.json)"

```


## notes

https://github.com/nlopes/slack/pull/638

https://github.com/gorilla/mux/issues/531
