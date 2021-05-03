package cmds

import (
	"errors"
	"fmt"
	"time"

	"github.com/ivanwang123/calendar/db"
	"github.com/spf13/cobra"
	"gopkg.in/toast.v1"
)

var notifyCmd = &cobra.Command{
	Use:     "notify",
	Short:   "Sends a system notification of your reminders",
	Example: "cali notify",
	Run: func(cmd *cobra.Command, args []string) {
		err := SendNotification()
		if err != nil {
			printError(err.Error())
		}
	},
}

func SendNotification() error {
	today, upcoming, err := getReminders()
	if err != nil {
		return errors.New("Oops, unable to get reminders")
	}

	var title string
	var message string

	if len(today) == 0 {
		title = "Upcoming deadlines"
		if len(upcoming) > 0 {
			for _, r := range upcoming {
				message += r + "\n"
			}
		} else {
			message = "No upcoming deadlines"
		}
	} else {
		title = fmt.Sprintf("Reminders for today, %s", formatMonthDay(time.Now()))
		for _, r := range today {
			message += r + "\n"
		}
	}

	configs := db.GetConfigurations()
	audioStr := configs["audio"]
	durationStr := configs["duration"]

	// TODO: figure out app icon
	notification := toast.Notification{
		AppID:   "Cali",
		Title:   title,
		Message: message,
		// Actions: []toast.Action{
		// 	{Type: "protocol", Label: "I'm a button", Arguments: "bingmaps:?q=sushi"},
		// 	{Type: "protocol", Label: "Me too!", Arguments: ""},
		// },
	}
	withDuration(durationStr, &notification)
	withAudio(audioStr, &notification)

	err = notification.Push()
	if err != nil {
		return errors.New("Oops, unable to send notification")
	}
	return nil
}

func withDuration(durationStr string, notification *toast.Notification) {
	switch durationStr {
	case "short":
		notification.Duration = toast.Short
	case "long":
		notification.Duration = toast.Long
	default:
		notification.Duration = toast.Short
	}
}

func withAudio(audioStr string, notification *toast.Notification) {
	switch audioStr {
	case "default":
		notification.Audio = toast.Default
	case "none":
		notification.Audio = toast.Silent
	case "fancy":
		notification.Audio = toast.Reminder
	default:
		notification.Audio = toast.Default
	}
}

func init() {
	RootCmd.AddCommand(notifyCmd)
}
