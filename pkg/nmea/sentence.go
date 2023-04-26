package nmea

import (
	"fmt"
	"time"
)

type Sentence []string

// BaseYear is the amount that is added to 2 digit years.
const BaseYear = 2000

func (se Sentence) Type() (string, error) {
	if len(se[0]) <= 2 {
		return "", fmt.Errorf("no header")
	}
	return string([]byte(se[0])[2:]), nil
}

//func (se Sentence) TimeAt(index int) (time.Time, error) {
//	var h, m, s int
//	_, err := fmt.Sscanf(se[index], "%02d%02d%02d", &h, &m, &s)
//	if err != nil {
//		return time.Time{}, err
//	}
//	return time.Date(0, 0, 0, h, m, s, 0, time.UTC), nil
//}

//func (se Sentence) DateAt(index int) (time.Time, error) {
//	var y, m, d int
//	_, err := fmt.Sscanf(se[index], "%02d%02d%02d", &d, &m, &y)
//	if err != nil {
//		return time.Time{}, err
//	}
//	return time.Date(BaseYear+y, time.Month(m), d, 0, 0, 0, 0, time.UTC), nil
//}

func (se Sentence) DateTimeAt(dIndex, tIndex int) (time.Time, error) {
	var ye, mo, da int
	_, err := fmt.Sscanf(se[dIndex], "%02d%02d%02d", &da, &mo, &ye)
	if err != nil {
		return time.Time{}, err
	}
	var h, m, s int
	_, err = fmt.Sscanf(se[tIndex], "%02d%02d%02d", &h, &m, &s)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(BaseYear+ye, time.Month(mo), da, h, m, s, 0, time.UTC), nil
}
