# NSQ Simple Publish and Subscribe
This repository contains Golang source code for demonstrating simple publish and subscribe message functionality using NSQ. See more details about [NSQ](https://nsq.io/overview/design.html).

This simple application is built using the [Cobra](https://github.com/spf13/cobra) library as the commander CLI.

## Demo

### Clone this repository
```bash
git clone https://github.com/fikri240794/nsq-simple-pubsub.git
```

### Set Up NSQ Server
You can set up the NSQ Server according to your preferences. For more details, you can visit [NSQ Docs Installation](https://nsq.io/deployment/installing.html).

I recommend you to set up the [NSQ Server using Docker](https://nsq.io/deployment/docker.html), especially for this demo. You can use the docker-compose.yml file that I've already created, you just need to run:
```bash
docker-compose up -d
```

### Publish
Before starting the publishing process, ensure that the NSQD (NSQ Daemon) is running. For this demo, NSQD has the following default configuration:
```
# NSQD Configuration
HOST: localhost
PORT: 4150
```
You can begin by thinking about the topic name you will use to publish messages. This is how to publish a message:
```bash
go run main.go publish --topic=your-topic-name
```
See more publish options:
```bash
go run main.go publish -h
```

### Subscribe
Before starting the subscribe process, ensure that the NSQLOOKUPD (NSQ LOOKUP Daemon) is running. For this demo, NSQLOOKUPD has the following default configuration:
```
# NSQLOOKUPD Configuration
HOST: localhost
PORT: 4161
```
You can subscribe to the topic that you have published messages to by:
```bash
go run main.go subscribe --channel=your-subscriber-group --topic=your-topic-name
```
See more subscribe options:
```bash
go run main.go subscribe -h
```

## Extra
![Alt Text](https://github.com/fikri240794/nsq-simple-pubsub/blob/main/nsq-simple-pubsub.gif)