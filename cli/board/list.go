/*
 * This file is part of arduino-cli.
 *
 * Copyright 2018 ARDUINO SA (http://www.arduino.cc/)
 *
 * This software is released under the GNU General Public License version 3,
 * which covers the main part of arduino-cli.
 * The terms of this license can be found at:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 * You can be released from the requirements of the above licenses by purchasing
 * a commercial license. Buying such a license is mandatory if you want to modify or
 * otherwise use the software for commercial activities involving the Arduino
 * software without disclosing the source code of your own applications. To purchase
 * a commercial license, send an email to license@arduino.cc.
 */

package board

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/arduino/arduino-cli/cli"
	"github.com/arduino/arduino-cli/commands/board"
	"github.com/arduino/arduino-cli/common/formatter"
	"github.com/arduino/arduino-cli/output"
	"github.com/arduino/arduino-cli/rpc"
	"github.com/spf13/cobra"
)

func initListCommand() *cobra.Command {
	listCommand := &cobra.Command{
		Use:     "list",
		Short:   "List connected boards.",
		Long:    "Detects and displays a list of connected boards to the current computer.",
		Example: "  " + cli.AppName + " board list --timeout 10s",
		Args:    cobra.NoArgs,
		Run:     runListCommand,
	}

	listCommand.Flags().StringVar(&listFlags.timeout, "timeout", "1s",
		"The timeout of the search of connected devices, try to increase it if your board is not found (e.g. to 10s).")
	return listCommand
}

var listFlags struct {
	timeout string // Expressed in a parsable duration, is the timeout for the list and attach commands.
}

// runListCommand detects and lists the connected arduino boards
func runListCommand(cmd *cobra.Command, args []string) {
	instance := cli.CreateInstance()

	// timeout, err := time.ParseDuration(listFlags.timeout)
	// if err != nil {
	// 	formatter.PrintError(err, "Invalid timeout.")
	// 	os.Exit(cli.ErrBadArgument)
	// }

	resp, err := board.BoardList(context.Background(), &rpc.BoardListReq{Instance: instance})
	if err != nil {
		formatter.PrintError(err, "Error detecting boards")
		os.Exit(cli.ErrNetwork)
	}

	if cli.OutputJSONOrElse(resp) {
		outputListResp(resp)
	}
}

func outputListResp(resp *rpc.BoardListResp) {
	sort.Slice(resp.Ports, func(i, j int) bool {
		x, y := resp.Ports[i], resp.Ports[j]
		return x.Protocol < y.Protocol || (x.Protocol == y.Protocol && x.Address < y.Address)
	})
	table := output.NewTable()
	table.SetHeader("Port", "Type", "Board Name", "FQBN")
	for _, port := range resp.Ports {
		address := port.Protocol + "://" + port.Address
		if port.Protocol == "serial" {
			address = port.Address
		}
		protocol := port.ProtocolLabel
		if len(port.Boards) > 0 {
			sort.Slice(port.Boards, func(i, j int) bool {
				x, y := port.Boards[i], port.Boards[j]
				return x.Name < y.Name || (x.Name == y.Name && x.FQBN < y.FQBN)
			})
			for _, b := range port.Boards {
				board := b.Name
				fqbn := b.FQBN
				table.AddRow(address, protocol, board, fqbn)
				// show address and protocol only on the first row
				address = ""
				protocol = ""
			}
		} else {
			board := "Unknown"
			fqbn := ""
			table.AddRow(address, protocol, board, fqbn)
		}
	}
	fmt.Print(table.Render())
}
