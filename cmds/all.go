package cmds

import (
	"fmt"
	"math"
	"time"

	"github.com/ivanwang123/calendar/db"
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:     "all",
	Short:   "Displays all reminders",
	Example: "cali all",
	Run: func(cmd *cobra.Command, args []string) {
		today, upcoming, err := getReminders()
		if err != nil {
			printError("Oops, unable to get reminders")
			return
		}
		printSubtitle(fmt.Sprintf("\nToday is %s (%d/%d)", formatMonthDay(time.Now()), time.Now().Month(), time.Now().Day()))
		if len(today) > 0 {
			for _, r := range today {
				fmt.Println(r)
			}
		} else {
			fmt.Println("No reminders for today")
		}

		printSubtitle("\nUpcoming")
		if len(upcoming) > 0 {
			for _, r := range upcoming {
				fmt.Println(r)
			}
		} else {
			fmt.Println("No upcoming deadlines")
		}
	},
}

func formatMonthDay(date time.Time) string {
	ordinal := "th"
	ones := date.Day() % 10
	tens := date.Day() % 100
	if ones == 1 && tens != 11 {
		ordinal = "st"
	} else if ones == 2 && tens != 12 {
		ordinal = "nd"
	} else if ones == 3 && tens != 13 {
		ordinal = "rd"
	}

	return fmt.Sprintf("%s %d%s", date.Month().String(), date.Day(), ordinal)
}

func getReminders() ([]string, []string, error) {
	reminders, err := db.AllReminders()
	if err != nil {
		return nil, nil, err
	}
	today := make([]string, 0)
	upcoming := make([]string, 0)
	todayKey := convertDateToKey(time.Now())
	idx := 1
	for _, r := range reminders {
		if r.Date == todayKey {
			today = append(today, fmt.Sprintf("%d. %s", idx, r.Reminder))
			idx++
		} else if r.Date > todayKey {
			upcoming = append(upcoming, fmt.Sprintf("%d. %s", idx, formatReminder(r.Reminder, convertKeyToDate(r.Date))))
			idx++
		} else {
			db.DeleteReminder(r.Key)
		}
	}
	return today, upcoming, nil
}

func formatReminder(reminder string, date time.Time) string {
	days := daysFromToday(date)
	if days == 1 {
		return fmt.Sprintf("%s (tomorrow)", reminder)
	}

	year, month, day := date.Date()
	dateStr := fmt.Sprintf("%d/%d", month, day)
	if year != time.Now().Year() {
		dateStr += fmt.Sprintf("/%d", year)
	}
	return fmt.Sprintf("%s - in %d days (%s)", reminder, days, dateStr)
}

func convertKeyToDate(dateKey int) time.Time {
	year, month, day := decodeDateKey(dateKey)
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func decodeDateKey(dateKey int) (year, month, day int) {
	year = dateKey / 10000
	month = (dateKey - year*10000) / 100
	day = (dateKey - year*10000 - month*100)
	return
}

func daysFromToday(date time.Time) int {
	hours := date.Sub(time.Now()).Hours()
	return int(math.Ceil(hours / 24))
}

func init() {
	RootCmd.AddCommand(allCmd)
}
