// Tasker - A pluggable task server for keeping track of all those To-Do's
// Copyright (C) 2019 Ryan Murphy <ryan@oatmealrais.in>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package util

import (
	"regexp"
	"time"

	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
)

var YYYY_MM_DD *regexp.Regexp = regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}")
var YYYY_MM_DD_HH_MM *regexp.Regexp = regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}")

// StringToTimestamp takes a string either formatted in YYYY-MM-DD or in RFC3339
// and returns a protobuf Timestamp. If there is any issue, it will return nil
func StringToTimestamp(s string) *tspb.Timestamp {
	var t time.Time
	var err error

	// For ease of writing out, we support YYYY-MM-DD defaulting to midnight
	// local time, and YYYY-MM-DD HH:MM, defaulting to local time
	if YYYY_MM_DD_HH_MM.MatchString(s) {
		t, err = time.ParseInLocation("2006-01-02 15:04", s, time.Now().Location())
	} else if YYYY_MM_DD.MatchString(s) {
		t, err = time.ParseInLocation("2006-01-02", s, time.Now().Location())
	} else {
		t, err = time.ParseInLocation(time.RFC3339, s, time.Now().Location())
	}

	if err != nil {
		return nil
	}

	result, err := ptypes.TimestampProto(t)
	if err != nil {
		return nil
	}

	return result
}

// TimestampToString takes a protobuf Timestamp and returns a string formatted
// as YYYY-MM-DD
func TimestampToString(t *tspb.Timestamp) string {
	result, err := ptypes.Timestamp(t)
	if err != nil {
		return ""
	}

	return result.In(time.Now().Location()).Format("2006-01-02")
}
