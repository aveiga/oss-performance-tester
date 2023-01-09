/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/aveiga/oss-performance-tester/pkg/amqp"
	"github.com/aveiga/oss-performance-tester/pkg/serialization"
	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
)

type DeclarationJson struct {
	ExchangeName string
	ExchangeType string
	Durable      bool
	AutoDelete   bool
	Internal     bool
	NoWait       bool
	Queue        struct {
		QueueName  string
		Durable    bool
		AutoDelete bool
		Exclusive  bool
		NoWait     bool
		Arguments  []struct {
			ArgumentKey   string
			ArgumentValue string
		}
	}
}

var connection string
var queueName string
var producers int
var consumers int
var payloadSize int
var declarationJson string
var wg sync.WaitGroup

// rabbitmqCmd represents the rabbitmq command
var rabbitmqCmd = &cobra.Command{
	Use:   "rabbitmq",
	Short: "RabbitMQ Load Tester",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("Got the following flags", connection, producers, consumers, payloadSize, queueName)
		client := amqp.NewMessagingClient(connection)
		serializedPayload := make([]byte, payloadSize)

		if declarationJson != "" {
			content, err := ioutil.ReadFile(declarationJson)
			if err != nil {
				log.Fatal("Error when opening file: ", err)
			}

			declaration, err := serialization.Deserialize[DeclarationJson](content)
			if err != nil {
				log.Fatal("Error deserializing file: ", err)
			}

			for i := 1; i <= consumers; i++ {
				wg.Add(1)
				SetupComplexConsumer(client, declaration)
			}

			for i := 1; i <= producers; i++ {
				// fmt.Println("Setting up producer")
				wg.Add(1)
				go SetupComplexProducer(client, serializedPayload, declaration)
			}

		} else {
			for i := 1; i <= consumers; i++ {
				wg.Add(1)
				SetupSimpleConsumer(client)
			}

			for i := 1; i <= producers; i++ {
				// fmt.Println("Setting up producer")
				wg.Add(1)
				go SetupSimpleProducer(client, serializedPayload)
			}
		}

		wg.Wait()
	},
}

func SetupSimpleProducer(client *amqp.MessagingClient, payload []byte) {
	for {
		// fmt.Println("publishing")
		err := client.PublishOnQueue(payload, queueName)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SetupComplexProducer(client *amqp.MessagingClient, payload []byte, declaration DeclarationJson) {
	for {
		// fmt.Println("publishing")
		exchange := amqp.Exchange{
			ExchangeName: declaration.ExchangeName,
			ExchangeType: declaration.ExchangeType,
			Durable:      declaration.Durable,
			AutoDelete:   declaration.AutoDelete,
			Internal:     declaration.Internal,
			NoWait:       declaration.NoWait,
		}
		queue := amqp.Queue{
			QueueName:  declaration.Queue.QueueName,
			Durable:    declaration.Queue.Durable,
			AutoDelete: declaration.Queue.AutoDelete,
			Exclusive:  declaration.Queue.Exclusive,
			NoWait:     declaration.Queue.NoWait,
		}
		err := client.Publish(payload, &exchange, &queue)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SetupSimpleConsumer(client *amqp.MessagingClient) {
	err := client.SubscribeToQueue(queueName, "Load Tester", func(d amqp091.Delivery) {})
	if err != nil {
		fmt.Println(err)
	}
}

func SetupComplexConsumer(client *amqp.MessagingClient, declaration DeclarationJson) {
	exchange := amqp.Exchange{
		ExchangeName: declaration.ExchangeName,
		ExchangeType: declaration.ExchangeType,
		Durable:      declaration.Durable,
		AutoDelete:   declaration.AutoDelete,
		Internal:     declaration.Internal,
		NoWait:       declaration.NoWait,
	}
	queue := amqp.Queue{
		QueueName:  declaration.Queue.QueueName,
		Durable:    declaration.Queue.Durable,
		AutoDelete: declaration.Queue.AutoDelete,
		Exclusive:  declaration.Queue.Exclusive,
		NoWait:     declaration.Queue.NoWait,
	}
	err := client.Subscribe(&exchange, &queue, "Load Tester", func(d amqp091.Delivery) {})
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(rabbitmqCmd)
	rabbitmqCmd.Flags().StringVarP(&connection, "connection", "c", "", "RabbitMQ connection string")
	rabbitmqCmd.Flags().StringVarP(&queueName, "queue", "q", "", "RabbitMQ queue name")
	rabbitmqCmd.Flags().IntVarP(&producers, "producers", "x", 1, "Number of queue producers")
	rabbitmqCmd.Flags().IntVarP(&consumers, "consumers", "y", 1, "Number of queue consumers")
	rabbitmqCmd.Flags().IntVarP(&payloadSize, "payload", "s", 1000, "Message payload size")
	rabbitmqCmd.Flags().StringVarP(&declarationJson, "file", "f", "", "RabbitMQ Exchange and Queue spec file path")
}
