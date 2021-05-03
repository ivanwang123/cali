package db

import (
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

var reminderBucket = []byte("reminders")
var settingsBucket = []byte("settings")

var db *bolt.DB

type Reminder struct {
	Key      uint64
	Date     int
	Reminder string
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(reminderBucket)
		return err
	})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(settingsBucket)
		return err
	})
}

func GetConfigurations() map[string]string {
	configs := make(map[string]string)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(settingsBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			configs[string(k)] = string(v)
		}
		return nil
	})
	return configs
}

func SetConfigurations(configs map[string]string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(settingsBucket)
		for k, v := range configs {
			err := b.Put([]byte(k), []byte(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func GetLastNotification() (time.Time, error) {
	var notificationDate time.Time
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(settingsBucket)
		timeStr := b.Get([]byte("notification"))
		if timeStr == nil {
			timeStr = []byte(time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05.999999999 -0700 MST"))
		}
		var err error
		notificationDate, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(timeStr))
		return err
	})
	if err != nil {
		return time.Now(), err
	}
	return notificationDate, nil
}

func SetLastNotification(date time.Time) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(settingsBucket)
		return b.Put([]byte("notification"), []byte(date.Format("2006-01-02 15:04:05.999999999 -0700 MST")))
	})
}

func CreateReminder(dateKey, reminder string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(reminderBucket)
		id64, _ := b.NextSequence()
		key := itob(id64)
		value := fmt.Sprintf("%s %s", dateKey, reminder)
		return b.Put(key, []byte(value))
	})
}

func AllReminders() ([]Reminder, error) {
	var reminders []Reminder
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(reminderBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			values := strings.SplitN(string(v), " ", 2)
			dateStr := values[0]
			reminder := values[1]
			date, _ := strconv.Atoi(dateStr)
			reminders = append(reminders, Reminder{
				Key:      btoi(k),
				Date:     date,
				Reminder: reminder,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(reminders, func(a, b int) bool {
		return reminders[a].Date < reminders[b].Date
	})
	return reminders, nil
}

func DeleteReminder(key uint64) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(reminderBucket)
		return b.Delete(itob(key))
	})
}

func DeleteAllReminders() error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(reminderBucket)
	})
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
