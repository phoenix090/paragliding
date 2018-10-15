package model

import (
	"strconv"
	"time"
)

/*
** getUptime updates uptime and formates it in ISO 8601 standard
 */
func GetUptime(t time.Time) (uptime string) {
	now := time.Now()
	newTime := now.Sub(t)
	hours := int(newTime.Hours())
	sek := strconv.Itoa(int(newTime.Seconds()) % 36000 % 60)
	var min, hour, y, m, d string

	// checking and setting when min gets to 1 or more
	if int(newTime.Seconds())%36000 >= 60 {
		min = strconv.Itoa(int(newTime.Minutes()) % 60)
	}

	// checking and setting when hour gets to 1 or more
	if hours >= 1 {
		hour = strconv.Itoa(hours)
	}

	// Setting the days correct
	if hours > 23 {
		d = strconv.Itoa(hours / 24)
		hour = strconv.Itoa(hours % 24)
	}
	days, _ := strconv.Atoi(d)
	// Setting the month correct
	if days > 31 {
		m = strconv.Itoa(days / 31)
		d = strconv.Itoa(days % 31)

	}
	months, _ := strconv.Atoi(m)
	// Setting the year correct
	if months > 12 {
		y = strconv.Itoa(months / 12)
		m = strconv.Itoa(months % 12)
	}

	uptime = "P" + y + "Y" + m + "M" + d + "DT" + hour + "H" + min + "M" + sek + "S"

	return uptime
}
