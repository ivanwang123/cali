package cmds

import (
	"fmt"
	"time"

	"github.com/ivanwang123/calendar/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [date]",
	Short:   "Displays all reminders on a specific day",
	Example: "cali list 4/27",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			printError(fmt.Sprintf("Oops, not enough arguments, please specify a date"))
			return
		}
		date, err := parseAnyDate(args)
		if err != nil {
			printError(err.Error())
			return
		}
		reminders, err := getRemindersOnDay(date)
		if err != nil {
			printError("Oops, unable to get reminders")
			return
		}
		printSubtitle(fmt.Sprintf("\nReminders for %s (%s)", formatMonthDay(date), fmt.Sprintf("%d/%d", date.Month(), date.Day())))
		if len(reminders) == 0 {
			fmt.Printf("No reminders for %s\n", formatMonthDay(date))
		}
		for _, r := range reminders {
			fmt.Println(r)
		}
	},
}

func getRemindersOnDay(date time.Time) ([]string, error) {
	reminders, err := db.AllReminders()
	if err != nil {
		return nil, err
	}
	dateReminders := make([]string, 0)
	dateKey := convertDateToKey(date)
	for _, r := range reminders {
		if r.Date == dateKey {
			dateReminders = append(dateReminders, r.Reminder)
		}
	}
	return dateReminders, nil
}

func init() {
	RootCmd.AddCommand(listCmd)
}
