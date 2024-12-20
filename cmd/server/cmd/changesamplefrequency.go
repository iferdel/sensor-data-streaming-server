package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/iferdel/sensor-data-streaming-server/internal/pubsub"
	"github.com/iferdel/sensor-data-streaming-server/internal/routing"
	"github.com/spf13/cobra"
)

var changeSampleFrequencyCmd = &cobra.Command{
	Use:   "changeSampleFrequency",
	Short: "Change the sample frequency [Hz] of a sensor",
	Run: func(cmd *cobra.Command, args []string) {

		if conn != nil {
			defer conn.Close()
		}

		sensorSerialNumber, err := cmd.Flags().GetString("sensor")
		if err != nil {
			log.Printf("error retrieving sensorid flag: %v", err)
			return
		}
		newSampleFrequency, err := cmd.Flags().GetString("changeSampleFrequency")
		if err != nil {
			log.Printf("error retrieving new sample frequency flag: %v", err)
			return
		}
		value, err := strconv.Atoi(newSampleFrequency)
		if err != nil {
			fmt.Println("sample frequency must be an integer greater than 0")
			return
		}

		err = pubsub.PublishGob(
			publishCh,                // amqp.Channel
			routing.ExchangeTopicIoT, // exchange
			fmt.Sprintf(routing.BindKeySensorCommandFormat, sensorSerialNumber), // routing key
			routing.CommandMessage{
				SensorName: sensorSerialNumber,
				Timestamp:  time.Now(),
				Command:    "changeSampleFrequency",
				Params: map[string]interface{}{
					"sampleFrequency": value,
				},
			}, // value
		)
		if err != nil {
			log.Printf("could not publish change sample frequency command: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(changeSampleFrequencyCmd)
	changeSampleFrequencyCmd.Flags().StringP("sensor", "s", "", "sensorid")
	changeSampleFrequencyCmd.Flags().StringP("changeSampleFrequency", "f", "", "sample frequency")
}
