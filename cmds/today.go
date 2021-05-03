package cmds

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Displays today's reminders",
	Run: func(cmd *cobra.Command, args []string) {
		today, _, err := getReminders()
		if err != nil {
			printError("Oops, unable to get reminders")
			return
		}
		printSubtitle(fmt.Sprintf("\nReminders for today, %s (%d/%d)", formatMonthDay(time.Now()), time.Now().Month(), time.Now().Day()))
		if len(today) > 0 {
			for _, r := range today {
				fmt.Println(r)
			}
		} else {
			fmt.Println("No reminders for today")
		}
	},
}

func init() {
	RootCmd.AddCommand(todayCmd)
}
