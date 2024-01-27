package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/nsqio/go-nsq"
	"github.com/spf13/cobra"
)

const (
	hostFlag    = "host"
	portFlag    = "port"
	topicFlag   = "topic"
	channelFlag = "channel"
)

var (
	publisher     *nsq.Producer
	subscriber    *nsq.Consumer
	publisherCmd  *cobra.Command
	subscriberCmd *cobra.Command
	mainCmd       *cobra.Command
)

type messageHandler struct{}

func buildPublisher(cmd *cobra.Command, args []string) (err error) {
	host, err := cmd.Flags().GetString(hostFlag)
	if err != nil {
		log.Printf("failed to get host with an error: %s\n", err.Error())
		return
	}

	port, err := cmd.Flags().GetInt(portFlag)
	if err != nil {
		log.Printf("failed to get port with an error: %s\n", err.Error())
		return
	}

	config := nsq.NewConfig()
	addr := fmt.Sprintf("%s:%d", host, port)

	publisher, err = nsq.NewProducer(addr, config)
	if err != nil {
		log.Printf("failed to create new publisher with an error: %s\n", err.Error())
		return
	}

	return
}

func publish(topic string) (err error) {
	fmt.Print("message: ")
	reader := bufio.NewReader(os.Stdin)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("failed to read input with an error: %s\n", err.Error())
		return
	}

	err = publisher.Publish(topic, []byte(message))
	if err != nil {
		log.Printf("failed to publish message %s with an error: %s\n", message, err.Error())
		return
	}

	return publish(topic)
}

func runPublisher(cmd *cobra.Command, args []string) (err error) {
	topic, err := cmd.Flags().GetString(topicFlag)
	if err != nil {
		log.Printf("failed to get topic with an error: %s\n", err.Error())
		return
	}

	errChan := make(chan error, 1)
	signalChan := make(chan os.Signal, 1)

	go func() {
		errChan <- publish(topic)
	}()

	signal.Notify(signalChan, os.Interrupt)

	select {
	case err = <-errChan:
		log.Printf("failed to publish message with an error: %s\n", err.Error())
	case <-signalChan:
	}

	publisher.Stop()

	return
}

func (h *messageHandler) HandleMessage(m *nsq.Message) (err error) {
	if len(m.Body) == 0 {
		return
	}

	log.Println(string(m.Body))

	return
}

func buildSubscriber(cmd *cobra.Command, args []string) (err error) {
	topic, err := cmd.Flags().GetString(topicFlag)
	if err != nil {
		log.Printf("failed to get topic with an error: %s\n", err.Error())
		return
	}

	channel, err := cmd.Flags().GetString(channelFlag)
	if err != nil {
		log.Printf("failed to get channel with an error: %s\n", err.Error())
		return
	}

	config := nsq.NewConfig()

	subscriber, err = nsq.NewConsumer(topic, channel, config)
	if err != nil {
		log.Printf("failed to create new subscriber with an error: %s\n", err.Error())
		return
	}

	subscriber.AddHandler(&messageHandler{})

	return
}

func runSubscriber(cmd *cobra.Command, args []string) (err error) {
	nsqHost, err := cmd.Flags().GetString(hostFlag)
	if err != nil {
		log.Printf("failed to get host with an error: %s\n", err.Error())
		return
	}

	nsqPort, err := cmd.Flags().GetInt(portFlag)
	if err != nil {
		log.Printf("failed to get port with an error: %s\n", err.Error())
		return
	}

	addr := fmt.Sprintf("%s:%d", nsqHost, nsqPort)
	err = subscriber.ConnectToNSQLookupd(addr)
	if err != nil {
		log.Printf("failed to connect to %s with an error: %s", addr, err.Error())
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	subscriber.Stop()

	return
}

func initPublisher() {
	publisherCmd = &cobra.Command{
		Use:     "publish",
		Short:   "Publish message(s) to given topic",
		Long:    "Publish message(s) to given topic",
		PreRunE: buildPublisher,
		RunE:    runPublisher,
	}

	publisherCmd.Flags().String(hostFlag, "localhost", "Publisher Host")
	publisherCmd.Flags().Int(portFlag, 4150, "Publisher Port")
	publisherCmd.Flags().String(topicFlag, "", "Topic")
	publisherCmd.MarkFlagRequired(topicFlag)
}

func initSubscriber() {
	subscriberCmd = &cobra.Command{
		Use:     "subscribe",
		Short:   "Subscribe message(s) to given topic",
		Long:    "Subscribe message(s) to given topic",
		PreRunE: buildSubscriber,
		RunE:    runSubscriber,
	}

	subscriberCmd.Flags().String(hostFlag, "localhost", "Subscriber Host")
	subscriberCmd.Flags().Int(portFlag, 4161, "Subscriber Port")
	subscriberCmd.Flags().String(channelFlag, "", "Channel")
	subscriberCmd.MarkFlagRequired(channelFlag)
	subscriberCmd.Flags().String(topicFlag, "", "Topic")
	subscriberCmd.MarkFlagRequired(topicFlag)
}

func init() {
	initPublisher()
	initSubscriber()

	mainCmd = &cobra.Command{
		Long: "NSQ Simple Publish and Subscribe",
	}

	mainCmd.AddCommand(
		publisherCmd,
		subscriberCmd,
	)
}

func main() {
	mainCmd.Execute()
}
