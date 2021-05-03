package cmds

import (
	"fmt"
	"strconv"

	"github.com/ivanwang123/calendar/db"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "del [index]...",
	Aliases: []string{"delete"},
	Short:   "Deletes one or more reminders",
	Example: "cali del 1\ncali del 1 2 3\ncali del all",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printError("Oops, expected at least one argument")
			return
		}
		reminders, err := db.AllReminders()
		if err != nil {
			printError("Oops, unable to delete reminder")
			return
		}
		if args[0] == "all" {
			err := db.DeleteAllReminders()
			if err != nil {
				printError("Oops, unable to delete reminders")
				return
			}
			plural := "s"
			if len(reminders) == 1 {
				plural = ""
			}
			printSuccess(fmt.Sprintf("Deleted %d reminder%s", len(reminders), plural))
			return
		}
		for _, r := range args {
			rIdx, err := strconv.Atoi(r)
			if err != nil {
				printError(fmt.Sprintf("Oops, index must be a number, received %s", r))
				return
			}
			if rIdx < 1 || rIdx > len(reminders) {
				printError(fmt.Sprintf("Oops, there is no reminder with index %d", rIdx))
				return
			}
			reminder := reminders[rIdx-1]
			err = db.DeleteReminder(reminder.Key)
			if err != nil {
				printError("Oops, unable to delete reminder")
				return
			}
			printSuccess(fmt.Sprintf("Deleted reminder: %s", reminder.Reminder))
		}
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
