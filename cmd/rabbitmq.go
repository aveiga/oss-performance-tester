/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"sync"

	"github.com/aveiga/oss-performance-tester/pkg/amqp"
	"github.com/spf13/cobra"
	streadwayAmqp "github.com/streadway/amqp"
)

var connection string
var queueName string
var producers int
var consumers int
var payloadSize int
var wg sync.WaitGroup

// rabbitmqCmd represents the rabbitmq command
var rabbitmqCmd = &cobra.Command{
	Use:   "rabbitmq",
	Short: "RabbitMQ Load Tester",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("Got the following flags", connection, producers, consumers, payloadSize, queueName)
		client := amqp.NewMessagingClient(connection)
		serializedPayload := make([]byte, payloadSize)

		for i := 1; i <= consumers; i++ {
			wg.Add(1)
			SetupConsumer(client)
		}

		for i := 1; i <= producers; i++ {
			// fmt.Println("Setting up producer")
			wg.Add(1)
			go SetupProducer(client, serializedPayload)
		}

		wg.Wait()
	},
}

func SetupProducer(client *amqp.MessagingClient, payload []byte) {
	for {
		// fmt.Println("publishing")
		err := client.PublishOnQueue(payload, queueName)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SetupConsumer(client *amqp.MessagingClient) {
	err := client.SubscribeToQueue(queueName, "Load Tester", func(d streadwayAmqp.Delivery) {})
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
}
