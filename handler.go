package main

import (
	"context"
	"strconv"
	"time"
)

// Handler of pasted text
type Handler interface {
	Match(ctx context.Context, text string) MatchResult
	Convert(ctx context.Context, text string, m MatchResult) string
}

type MatchResult struct {
	Match bool
	Value any
}

type timeHandler struct{}

type timeVal struct {
	timeType string
	format   string
}

var knownTimeFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05-0700",
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC850,
	time.Layout,
	time.RubyDate,
	time.ANSIC,
	time.UnixDate,
	time.RFC822,
	time.RFC822Z,
}

const defaultTimeFormat = time.RFC3339

func (h *timeHandler) Match(ctx context.Context, text string) MatchResult {
	// try to convert text to timestamp
	_, err := strconv.ParseInt(text, 10, 64)
	if err == nil {
		return MatchResult{Match: true, Value: &timeVal{timeType: "timestamp", format: defaultTimeFormat}}
	}

	for _, format := range knownTimeFormats {
		// try to convert text to time
		_, err = time.ParseInLocation(format, text, tz)
		if err == nil {
			return MatchResult{Match: true, Value: &timeVal{timeType: "time", format: format}}
		}
	}

	return MatchResult{Match: false}
}

func (h *timeHandler) Convert(ctx context.Context, text string, m MatchResult) string {
	tv := m.Value.(*timeVal)
	if tv.timeType == "timestamp" {
		// try to convert text to timestamp
		timestamp, err := strconv.ParseInt(text, 10, 64)
		if err == nil {
			tt := time.Unix(timestamp, 0)
			if secType == "ms" {
				tt = time.UnixMilli(timestamp)
			}
			return tt.Format(defaultTimeFormat)
		}

		return text
	}

	if tv.timeType == "time" {
		// try to convert text to time
		t, err := time.ParseInLocation(tv.format, text, tz)
		if err == nil {
			ts := t.Unix()
			if secType == "ms" {
				ts = t.UnixMilli()
			}
			return strconv.FormatInt(ts, 10)
		}

		return text
	}

	return text
}
