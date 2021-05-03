package cmds

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ivanwang123/calendar/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add [reminder] [date]",
	Short:   "Adds a reminder",
	Example: "cali add \"Turn in essay\" 4/29\ncali add \"Turn in essay\" 4/29/2022\ncali add \"Turn in essay\" Apr 29\ncali add \"Turn in essay\" in 2d\ncali add \"Turn in essay\" in 1d 2w 3m 1y",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			printError("Oops, not enough arguments")
			return
		}

		reminder := args[0]
		if len(reminder) == 0 {
			printError("Oops, no reminder was given")
			return
		}
		deadlineDate, err := parseAnyDate(args[1:])
		if err != nil {
			printError(err.Error())
			return
		}

		cmp := compareDateToToday(deadlineDate)
		if cmp < 0 {
			printError("Oops, that deadline has already past")
			return
		}
		dateKey := fmt.Sprint(convertDateToKey(deadlineDate))
		err = db.CreateReminder(dateKey, reminder)
		if err != nil {
			printError("Oops, something went wrong")
			return
		}
		printSuccess(fmt.Sprintf("Added reminder: %s", formatReminder(reminder, deadlineDate)))
	},
}

// 20210429 = dateKey
// 4/29 = numDate
// Apr 29 = strDate
// 1d 1w 1m 1y = durDate

func convertDateToKey(date time.Time) int {
	y, m, d := date.Date()
	keyStr := fmt.Sprintf("%04d%02d%02d", y, m, d)
	key, _ := strconv.Atoi(keyStr)
	return key
}

func parseAnyDate(args []string) (time.Time, error) {
	var date time.Time
	var err error
	if strings.Contains(args[0], "/") {
		date, err = parseNumDate(args[0])
	} else if args[0] == "in" {
		durDate := strings.Join(args[1:], " ")
		date, err = parseDurDate(durDate)
	} else {
		strDate := strings.Join(args[0:], " ")
		date, err = parseStrDate(strDate)
	}
	if err != nil {
		return time.Now(), err
	}
	return date, nil
}

func parseDurDate(durDate string) (time.Time, error) {
	parts := strings.Split(durDate, " ")
	addDays := 0
	addWeeks := 0
	addMonths := 0
	addYears := 0
	for _, p := range parts {
		durStr := p[:len(p)-1]
		durType := p[len(p)-1:]
		dur, err := strconv.Atoi(durStr)
		if err != nil {
			return time.Now(), fmt.Errorf("Oops, invalid duration. Expected number, received '%s'", durStr)
		}
		switch durType {
		case "d":
			addDays = dur
		case "w":
			addWeeks = dur
		case "m":
			addMonths = dur
		case "y":
			addYears = dur
		default:
			return time.Now(), fmt.Errorf("Oops, invalid time notation. Expected 'd', 'w', 'm', or 'y', received '%s'", durType)
		}
	}
	return time.Date(time.Now().Year()+addYears, time.Now().Month()+time.Month(addMonths), time.Now().Day()+addDays+7*addWeeks, 0, 0, 0, 0, time.UTC), nil
}

func parseNumDate(numDate string) (time.Time, error) {
	parts := strings.Split(numDate, "/")
	formattedNumDate := numDate
	if len(parts) == 2 {
		formattedNumDate += fmt.Sprintf("/%d", time.Now().Year())
	} else if len(parts) != 3 {
		return time.Now(), fmt.Errorf("Oops, invalid date. Expected mm/dd or mm/dd/yy format, received '%s'", numDate)
	}
	date, err := time.Parse("1/2/2006", formattedNumDate)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("1/2/06", formattedNumDate)
	if err == nil {
		return date, nil
	}
	return time.Now(), fmt.Errorf("Oops, invalid date. Expected mm/dd or mm/dd/yy format, received '%s'", numDate)
}

func parseStrDate(strDate string) (time.Time, error) {
	parts := strings.Split(strDate, " ")
	formattedStrDate := strDate
	if len(parts) == 2 {
		formattedStrDate += fmt.Sprintf(" %d", time.Now().Year())
	} else if len(parts) != 3 {
		return time.Now(), fmt.Errorf("Oops, invalid date. Expected 'Month Day' or 'Month Day Year', received '%s'", strDate)
	}
	date, err := time.Parse("January _2 2006", formattedStrDate)
	if err == nil {
		return date, nil
	}
	date, err = time.Parse("Jan _2 2006", formattedStrDate)
	if err == nil {
		return date, nil
	}
	return time.Now(), fmt.Errorf("Oops, invalid date. Expected 'Month Day' or 'Month Day Year', received '%s'", strDate)
}

func compareDateToToday(date time.Time) int {
	ty, tm, td := time.Now().Date()
	dy, dm, dd := date.Date()
	if ty != dy {
		return dy - ty
	}
	if tm != dm {
		return int(dm) - int(tm)
	}
	if td != dd {
		return dd - td
	}
	return 0
}

func init() {
	RootCmd.AddCommand(addCmd)
}
