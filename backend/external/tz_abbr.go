package external

import "fmt"

var timezoneAbbrMap = map[string]string{
	"SST":   "Pacific/Midway",
	"HST":   "Pacific/Honolulu",
	"AKST":  "America/Anchorage",
	"PST":   "America/Los_Angeles",
	"PDT":   "America/Los_Angeles",
	"MST":   "America/Denver",
	"CST":   "America/Chicago",
	"EST":   "America/New_York",
	"-04":   "America/Caracas",
	"-03":   "America/Santiago",
	"-02":   "America/Noronha",
	"-01":   "Atlantic/Cape_Verde",
	"UTC":   "UTC",
	"GMT":   "Europe/London",
	"CET":   "Europe/Brussels",
	"EET":   "Africa/Cairo",
	"SAST":  "Africa/Johannesburg",
	"MSK":   "Europe/Moscow",
	"+0330": "Asia/Tehran",
	"+04":   "Asia/Baku",
	"+0430": "Asia/Kabul",
	"+05":   "Asia/Tashkent",
	"IST":   "Asia/Kolkata",
	"+0530": "Asia/Colombo",
	"+0545": "Asia/Kathmandu",
	"+0630": "Indian/Cocos",
	"+07":   "Asia/Bangkok",
	"+08":   "Asia/Singapore",
	"KST":   "Asia/Seoul",
	"JST":   "Asia/Tokyo",
	"ACST":  "Australia/Darwin",
	"ChST":  "Pacific/Guam",
	"AEST":  "Australia/Brisbane",
	"AEDT":  "Australia/Melbourne",
	"+11":   "Australia/Lord_Howe",
	"+12":   "Pacific/Fiji",
	"NZDT":  "Pacific/Auckland",
	"+1345": "Pacific/Chatham",
	"+13":   "Pacific/Tongatapu",
}

func GetTimezoneForStatusioTzAbbr(timezoneAbbr string) (string, error) {
	offset, ok := timezoneAbbrMap[timezoneAbbr]
	if !ok {
		return "", fmt.Errorf("can find timezone name for the given abbr %v", timezoneAbbr)
	}
	return offset, nil
}
