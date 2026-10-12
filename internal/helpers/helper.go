package helpers

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GetTokenFromBearer(token string) (string, error) {
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}
	return parts[1], nil
}

func ParseTimestampUTC(timestamp time.Time) (time.Time, string) {
	parseTime := timestamp.Format("2006-01-02 15:04:05.000-07:00")
	timezone := GetTimezone(parseTime[len(parseTime)-6:])
	timestamp = timestamp.UTC()
	return timestamp, timezone
}

func GetTimezone(offset string) string {
	switch offset {
	case "+07:00":
		return "WIB"
	case "+08:00":
		return "WITA"
	case "+09:00":
		return "WIT"
	default:
		return "unknown"
	}
}

func GetTimeLocation(inputTime time.Time) (string, error) {
	_, offset := inputTime.Zone()

	offsetStr := fmt.Sprintf("%+03d:%02d", offset/3600, (offset%3600)/60)

	offsetToTimeZone := map[string]string{
		"+07:00": "Asia/Jakarta",
		"+08:00": "Asia/Makassar",
		"+09:00": "Asia/Jayapura",
	}

	timeZoneName, found := offsetToTimeZone[offsetStr]
	if !found {
		log.Println("Timezone not found, using default timezone Asia/Jakarta")
		return "Asia/Jakarta", nil
	}

	return timeZoneName, nil
}

func StringToPascalCase(input string) string {
	segments := strings.Split(input, ".")
	titleCaser := cases.Title(language.Und)

	for i, segment := range segments {
		segment = strings.ReplaceAll(segment, "_", " ")
		words := strings.Fields(segment)
		for j := range words {
			words[j] = titleCaser.String(words[j])
		}
		segments[i] = strings.Join(words, "")
	}
	return strings.Join(segments, ".")
}

func ConvertToKebabCase(input string) string {
	input = strings.ToLower(input)

	var result strings.Builder
	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			result.WriteRune(char)
		} else if unicode.IsSpace(char) {
			result.WriteRune('-')
		}
	}
	return result.String()
}
