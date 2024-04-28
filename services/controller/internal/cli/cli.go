package cli

import (
	"fmt"
	"hack/services/controller/internal/pinger"
	. "hack/services/controller/internal/status"
	cfg "hack/services/controller/pkg/handleconfig"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type Ports struct {
	OldPort string
	NewPort string
}

func RunCLI() {
	var action, service, port string

	var rootCmd = &cobra.Command{
		Use:   "controller",
		Short: "Controller is a CLI for managing microservices",
		Long:  `Controller is a command line interface for managing microservices through start, stop, and add commands.`,
	}

	var cmdSearchBC = &cobra.Command{
		Use:   "searchbc",
		Short: "Searches for the microservice with keyword 'bc'",
		Long:  `Searches for the microservice identified by the keyword 'bc' and displays information about it.`,
		Run: func(cmd *cobra.Command, args []string) {
			action = "searchbc"
			fmt.Printf("%s...:\n", action)

		},
	}

	var cmdStatus = &cobra.Command{
		Use:   "status",
		Short: "Displays the status of the microservices",
		Long:  `Displays current status information about all managed microservices, including their operational state and health.`,
		Run: func(cmd *cobra.Command, args []string) {
			HandleStatus()
		},
	}

	var cmdStart = &cobra.Command{
		Use:   "start [model] [port]",
		Short: "Starts the model",
		Long:  `Starts the specified model on the given port.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action = "start"
			service = args[0]
			if service == "bc" {
				cfg.ChangeJson(action, service, "")
			}
			if len(args) > 1 {
				port = args[1]
				if port == "oldport" || port == "" {
					cfg.ChangeJson(action, service, "oldport")
				} else {
					if _, err := strconv.Atoi(port); err != nil {
						fmt.Fprintf(os.Stderr, "Error: Port must be a number\n")
						os.Exit(1)
					}
				}
			} else {
				port = "not specified"
			}
			fmt.Printf("%sing %s on port %s...\n", action, service, port)
		},
	}

	var cmdStop = &cobra.Command{
		Use:   "stop [model]",
		Short: "Stops the model",
		Long:  `Stops the specified model.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action = "stop"
			service = args[0]
			if service == "bc" {
				return
			}
			port = "oldport"
			fmt.Printf("%sping %s on %s...\n", action, service, port)
		},
	}

	var cmdChPort = &cobra.Command{
		Use:   "chport [model] [new_port]",
		Short: "Changes the port of the model",
		Long:  `Changes the port of the specified model from old_port to new_port.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			action = "chport"
			service = fmt.Sprintf("m%s", args[0])
			port = args[1]

			fmt.Printf("changing port for %s service...\n", service[1:])

			if _, err := strconv.Atoi(port); err != nil {
				fmt.Fprintf(os.Stderr, "Error: New port must be a number\n")
				os.Exit(1)
			}

			var ps Ports

			ps.OldPort = cfg.ChangeJson(action, service, port)
			ps.NewPort = port
			pinger.DoAction(ps.OldPort, ps.NewPort)

			fmt.Printf("port is changed to %s\n", port)
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
			cfg.ChangeJson(action, service, port)
			fmt.Printf("%sing new service %s on port: %s\n", action, service, port)
		},
	}

	var cmdDel = &cobra.Command{
		Use:   "del [model]",
		Short: "Deletes the model",
		Long:  `Deletes the specified model from the system.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			action = "del"
			service = args[0]
			fmt.Printf("Action: %s\nModel: %s deleted\n", action, service)
		},
	}

	rootCmd.AddCommand(cmdStart, cmdStop, cmdAdd, cmdChPort, cmdDel, cmdSearchBC, cmdStatus)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
