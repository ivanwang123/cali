package cmds

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/ivanwang123/calendar/db"
	"github.com/spf13/cobra"
)

var subtitleColor = color.New(color.FgCyan)
var successColor = color.New(color.FgGreen)
var errorColor = color.New(color.FgRed)

var RootCmd = &cobra.Command{
	Use:   "cali",
	Short: "Cali is your deadline reminder buddy",
}

func setConfig() {
	configs := db.GetConfigurations()
	tc := configs["subtitle"]
	sc := configs["success"]
	ec := configs["error"]
	subtitleColor, _ = strToColor(tc, subtitleColor)
	successColor, _ = strToColor(sc, successColor)
	errorColor, _ = strToColor(ec, errorColor)
}

func printSubtitle(message string) {
	setConfig()
	if subtitleColor == color.New() {
		fmt.Println(message)
	} else {
		subtitleColor.Println(message)
	}
}

func printSuccess(message string) {
	setConfig()
	if successColor == color.New() {
		fmt.Println(message)
	} else {
		successColor.Println(message)
	}
}

func printError(message string) {
	setConfig()
	if errorColor == color.New() {
		fmt.Println(message)
	} else {
		errorColor.Println(message)
	}
}
