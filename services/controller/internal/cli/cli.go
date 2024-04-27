package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func RunCLI() (string, string, string) {
	var action, service, port string

	var rootCmd = &cobra.Command{
		Use:   "controller",
		Short: "Controller is a CLI for managing models",
		Long:  `Controller is a command line interface for managing models through start, stop, and add commands.`,
	}

	var cmdStart = &cobra.Command{
		Use:   "start [model] [port]",
		Short: "Starts the model",
		Long:  `Starts the specified model on the given port.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			action = "start"
			service = args[0]
			port = args[1]

			if port != "oldport" {
				if _, err := strconv.Atoi(port); err != nil {
					fmt.Fprintf(os.Stderr, "Error: Port must be a number\n")
					os.Exit(1)
				}
			}

			fmt.Printf("%sing %s on %s...\n", action, service, port)
		},
	}

	var cmdStop = &cobra.Command{
		Use:   "stop [model]",
		Short: "Stops the model",
		Long:  `Stops the specified model.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action = "stop"
			service = args[0]
			port = "oldport"
			fmt.Printf("%sping %s on %s...\n", action, service, port)
		},
	}

	var cmdChPort = &cobra.Command{
		Use:   "chport [model] [old_port] [new_port]",
		Short: "Changes the port of the model",
		Long:  `Changes the port of the specified model from old_port to new_port.`,
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			action = "chport"
			service = args[0]
			port = args[1]

			if _, err := strconv.Atoi(port); err != nil {
				fmt.Fprintf(os.Stderr, "Error: New port must be a number\n")
				os.Exit(1)
			}

			fmt.Printf("change port for %s service, port changed to %s\n", service, port)
		},
	}

	var cmdAdd = &cobra.Command{
		Use:   "add [model] [port]",
		Short: "Adds a new microservice",
		Long:  `Adds a new microservice to the system.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			action = "add"
			service = args[0]
			port = args[1]
			fmt.Printf("%sing new service %s on port: %s\n", action, service, port)
		},
	}

	rootCmd.AddCommand(cmdStart, cmdStop, cmdAdd, cmdChPort)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return action, service, port
}
