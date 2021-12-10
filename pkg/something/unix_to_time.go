package something

import "time"

func ToTime(ts float64) time.Time {
	sec := int64(ts)
	nano := (ts - float64(sec)) * 1000000000

	return time.Unix(sec, int64(nano))
}
