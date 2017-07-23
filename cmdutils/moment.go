package cmdutils

import (
	"fmt"
	"github.com/go-shadow/moment"
	"strconv"
	"strings"
	"time"
)

const (
	buildFormat = "MMMM Do, YYYY"
)

// Lapsed retrieves the time lapsed since timestamp date
func Lapsed(timestamp string) string {

	n := moment.New()
	d := toMoment(timestamp)

	return lapsed(d, n)

}

// BuildDate returns a string representing build date
func BuildDate(timestamp string) string {

	return toMoment(timestamp).Format(buildFormat)

}

/*

	Convert an ISO looking string to a moment date

*/
func toMoment(timestamp string) *moment.Moment {

	d := moment.New()

	sp := strings.Split(timestamp, "T")
	ds := strings.Split(sp[0], "-")
	ts := strings.Split(sp[1], ".")

	var m, dy int64
	y, _ := strconv.ParseInt(ds[0], 10, 32)
	m, _ = strconv.ParseInt(ds[1], 10, 32)
	dy, _ = strconv.ParseInt(ds[2], 10, 32)

	d.Strtotime(ts[0])
	d.SetYear(int(y))
	d.SetMonth(time.Month(int(m)))
	d.SetDay(int(dy))

	return d

}

/*

	Calculate the time lapsed between start and end moment dates

*/
func lapsed(s *moment.Moment, e *moment.Moment) string {

	diff := e.GetDiff(s)
	res := diff.Humanize()

	// humanize only works for on same day for now
	if res != "diff is in days" {
		return res
	}

	// otherwise do a hacky fix
	res = ""

	plural := "s"
	if diff.InYears() >= 1 {
		if diff.InYears() == 1 {
			plural = ""
		}
		res = res + fmt.Sprintf("%d year%s ", diff.InYears(), plural)
	} else if diff.InMonths() >= 1 {
		if diff.InMonths() == 1 {
			plural = ""
		}
		res = res + fmt.Sprintf("%d month%s ", diff.InMonths(), plural)
	} else {
		if diff.InDays() == 1 {
			plural = ""
		}
		res = res + fmt.Sprintf("%d day%s ", diff.InDays(), plural)
	}
	res = res + "ago"

	return res

}
