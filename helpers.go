package senml

import (
	"slices"
	"strings"
	"time"
)

type RecordFinder = func(r Record) bool

func FindByName(n string) RecordFinder {
	return func(r Record) bool {
		return strings.EqualFold(r.Name, n)
	}
}

func FindByNormalizedName(bn, n string) RecordFinder {
	return func(r Record) bool {
		rName := r.BaseName + r.Name
		bnn := bn + n
		return strings.EqualFold(rName, bnn)
	}
}

func FindByUnit(u string) RecordFinder {
	return func(r Record) bool {
		return r.Unit == u
	}
}

func (p Pack) GetRecord(fn RecordFinder) (Record, bool) {
	idx := slices.IndexFunc(p, fn)

	if idx == -1 {
		return Record{}, false
	}

	clone := p.Clone()
	clone.Normalize()

	return clone[idx], true
}

func (p Pack) GetValue(fn RecordFinder) (float64, bool) {
	r, ok := p.GetRecord(fn)

	if !ok {
		return 0.0, false
	}

	return r.GetValue()	
}

func (r Record) GetValue() (float64, bool) {
	if r.Value == nil {
		return 0.0, false
	}

	return *r.Value, true
}

func (p Pack) GetValueWithUnit(fn RecordFinder) (float64, string, bool) {
	r, ok := p.GetRecord(fn)
	if !ok || r.Value == nil {
		return 0.0, "", false
	}

	return *r.Value, r.Unit, true
}

func (p Pack) GetBoolValue(fn RecordFinder) (bool, bool) {
	r, ok := p.GetRecord(fn)
	if !ok || r.BoolValue == nil {
		return false, false
	}

	return *r.BoolValue, true
}

func (p Pack) GetStringValue(fn RecordFinder) (string, bool) {
	r, ok := p.GetRecord(fn)
	if !ok {
		return "", false
	}

	return r.StringValue, true
}

func (p Pack) GetTime(fn RecordFinder) (time.Time, bool) {
	r, ok := p.GetRecord(fn)
	if !ok {
		return time.Time{}, false
	}

	return r.GetTime()
}

func (r Record) GetTime() (time.Time, bool) {
	return time.Unix(int64(r.Time), 0), true
}

func (p Pack) GetLatLon() (float64, float64, bool) {
	r, ok := p.GetRecord(FindByUnit(UnitLat))
	if !ok || r.Value == nil {
		return 0.0, 0.0, false
	}

	lat := *r.Value

	r, ok = p.GetRecord(FindByUnit(UnitLon))
	if !ok || r.Value == nil {
		return 0.0, 0.0, false
	}

	lon := *r.Value

	return lat, lon, true
}

func (p Pack) GetSum(fn RecordFinder) (float64, bool) {
	r, ok := p.GetRecord(fn)
	if !ok || r.Sum == nil {
		return 0.0, false
	}

	return *r.Sum, true
}
