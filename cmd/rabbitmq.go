/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var connection string
var producers int
var consumers int
var payloadSize int
var queueName string

// rabbitmqCmd represents the rabbitmq command
var rabbitmqCmd = &cobra.Command{
	Use:   "rabbitmq",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Got the following flags", connection, producers, consumers, payloadSize, queueName)
	},
}

func init() {
	rootCmd.AddCommand(rabbitmqCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rabbitmqCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rabbitmqCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rabbitmqCmd.Flags().StringVarP(&connection, "connection", "c", "", "RabbitMQ connection string")
	rabbitmqCmd.Flags().IntVarP(&producers, "producers", "x", 1, "Number of RabbitMQ producers")
	rabbitmqCmd.Flags().IntVarP(&consumers, "consumers", "y", 1, "Number of RabbitMQ consumers")
	rabbitmqCmd.Flags().IntVarP(&payloadSize, "payload", "s", 1000, "Message payload size")
	rabbitmqCmd.Flags().StringVarP(&queueName, "queue", "q", "", "Queue name")
}
