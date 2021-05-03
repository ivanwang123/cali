package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ivanwang123/calendar/cmds"
	"github.com/ivanwang123/calendar/db"
)

func main() {
	rootPath, err := os.UserCacheDir()
	if err != nil {
		log.Fatal("Unable to find local AppData directory")
	}

	dbDir := filepath.Join(rootPath, "Cali")
	err = os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		log.Fatal("Unable to create directory for database")
	}
	dbPath := filepath.Join(dbDir, "cali.db")
	must(db.Init(dbPath))

	if len(os.Args) == 1 {
		sendDailyNotification()
	}

	if len(os.Args) > 0 && (os.Args[0] == "cali" || os.Args[0] == "cali.exe") {
		must(cmds.RootCmd.Execute())
	}
}

func sendDailyNotification() {
	notificationDate, err := db.GetLastNotification()
	if err != nil {
		fmt.Println("Oops, unable to get last notification")
	}
	ty, tm, td := time.Now().Date()
	dy, dm, dd := notificationDate.Date()
	if ty != dy || tm != dm || td != dd {
		err = cmds.SendNotification()
		if err == nil {
			db.SetLastNotification(time.Now())
		}
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
