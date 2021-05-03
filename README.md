# Database

Reminders bucket

- store sequence id as key
- store deadline as value in format [01-12][01-31][00-99] (ex. 042521)
- store reminder as value after deadline with space inbetween

# CLI

Add reminder: cali "Turn in homework" 4/20/2021 or cali "Turn in homework" in 2d

- year is optional, default is current year
- order of date and reminder is optional?
- quotes are optional if last param
- use 'in xd yw zm' to set reminder some time from today

List reminders: cali or cali all

- output today's date and reminders
- output upcoming: 4/20 (Today/Tomorrow/In x days/In x months) Turn in homework

Today is April 20th (4/20)

1. Turn in homework
2. Finish project

Upcoming

3. Finish project (in 1 day on 4/21)
4. Turn in project (in >1 week on 4/28)

List reminders on a day: cali today, cali 4/20

- output all reminders on a specific date

Send notification: cali notify

- output all today
- output 3 or 5 upcoming

\*\*List past reminders: cali past

- output: 4/20 (Yesterday/x days ago) Turn in homework
- save up to 1 week

Delete reminder: cali delete [reminder]

- reminder is the index of the reminder when printed
- reminder can be "today" or 4/20 date to delete all on that day

Set configuration: cali set --color=blue --audio=reminder

# Notification

-Title: Reminders for today, April 26th
-Message: 1. Finish project\n2. Do homework

-Title: Upcoming deadlines
-Message: 1. Turn in project (in 1 week on 4/28)

# Extra

-Improve documentation w/ examples, longs, etc.
-Mac support
-Add app icon

# TODO
