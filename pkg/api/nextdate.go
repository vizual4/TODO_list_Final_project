package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const formatTime = "20060102"

func NextDate(now time.Time, dstart, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("incorrect repeat format")
	}

	nowDate, err := time.Parse(formatTime, now.Format(formatTime))
	if err != nil {
		return "", err
	}

	rule := strings.Split(repeat, " ")

	date, err := time.Parse(formatTime, dstart)
	if err != nil {
		return "", err
	}

	switch rule[0] {
	case "d":
		if len(rule) < 2 {
			return "", fmt.Errorf("missing days count for repeat rule %s", rule[0])
		}
		repeatNum, err := strconv.Atoi(rule[1])

		if err != nil {
			return "", err
		}

		if repeatNum <= 0 || repeatNum > 400 {
			return "", fmt.Errorf("repeat number bigger than 400")
		}

		for {
			date = date.AddDate(0, 0, repeatNum)
			if date.After(nowDate) {
				break
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(nowDate) {
				break
			}
		}
	case "w":
		if len(rule) < 2 {
			return "", fmt.Errorf("missing days count for repeat rule %s", rule[0])
		}

		daysMap, err := parseWeekToMap(rule[1])
		if err != nil {
			return "", err
		}

		for {
			date = date.AddDate(0, 0, 1)
			weekday := date.Weekday()
			if weekday == time.Sunday {
				weekday = 7
			}
			if daysMap[int(weekday)] && date.After(nowDate) {
				break
			}
		}
	case "m":
		if len(rule) < 2 {
			return "", fmt.Errorf("incorrect format for repeat rule %s", rule[0])
		}

		daysMap, err := parseDaysToMap(rule[1])
		if err != nil {
			return "", err
		}

		var monthsMap map[int]bool
		if len(rule) > 2 {
			monthsMap, err = parseMonthsToMap(rule[2])
			if err != nil {
				return "", err
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			if matchMonthDay(date, daysMap, monthsMap) && date.After(nowDate) {
				break
			}
		}

	default:
		return "", fmt.Errorf("incorrect format")
	}

	return date.Format(formatTime), nil
}

func parseWeekToMap(days string) (map[int]bool, error) {
	if days == "" {
		return nil, fmt.Errorf("empty days rule")
	}

	daysMap := make(map[int]bool)
	parts := strings.Split(days, ",")

	for _, v := range parts {
		day, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		if day < 1 || day > 7 {
			return nil, fmt.Errorf("incorrect number: %s", v)
		}

		if daysMap[day] {
			return nil, fmt.Errorf("duplicate day number: %d", day)
		}

		daysMap[day] = true
	}
	return daysMap, nil
}

func parseDaysToMap(days string) (map[int]bool, error) {
	if days == "" {
		return nil, fmt.Errorf("empty days rule")
	}
	daysMap := make(map[int]bool)
	parts := strings.Split(days, ",")

	for _, v := range parts {
		day, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		if day < -2 || day > 31 || day == 0 {
			return nil, fmt.Errorf("incorrect number: %s", v)
		}

		if daysMap[day] {
			return nil, fmt.Errorf("duplicate day number: %d", day)
		}

		daysMap[day] = true
	}

	return daysMap, nil
}

func parseMonthsToMap(months string) (map[int]bool, error) {
	if months == "" {
		return nil, fmt.Errorf("empty month rule")
	}

	monthsMap := make(map[int]bool)
	parts := strings.Split(months, ",")

	for _, v := range parts {
		month, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		if month < 1 || month > 12 {
			return nil, fmt.Errorf("incorrect number: %s", v)
		}

		if monthsMap[month] {
			return nil, fmt.Errorf("duplicate day number: %d", month)
		}

		monthsMap[month] = true
	}

	return monthsMap, nil
}

func matchMonthDay(t time.Time, daysMap, monthsMap map[int]bool) bool {
	if len(monthsMap) > 0 && !monthsMap[int(t.Month())] {
		return false
	}

	firstOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	daysInMonth := firstOfNextMonth.AddDate(0, 0, -1).Day() // 28 ф

	currentDay := t.Day() // 1 ф

	for targetDay := range daysMap {
		actualDay := targetDay
		if targetDay < 0 {
			actualDay = daysInMonth + targetDay + 1
		}

		if currentDay == actualDay {
			return true
		}
	}

	return false
}
