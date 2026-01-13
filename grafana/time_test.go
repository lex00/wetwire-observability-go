package grafana

import (
	"testing"
)

func TestNewTimeRange(t *testing.T) {
	tr := NewTimeRange("now-1h", "now")
	if tr == nil {
		t.Fatal("NewTimeRange() returned nil")
	}
	if tr.From != "now-1h" {
		t.Errorf("From = %v, want now-1h", tr.From)
	}
	if tr.To != "now" {
		t.Errorf("To = %v, want now", tr.To)
	}
}

func TestTimeRange_LastHour(t *testing.T) {
	tr := LastHour()
	if tr.From != "now-1h" {
		t.Errorf("From = %v, want now-1h", tr.From)
	}
	if tr.To != "now" {
		t.Errorf("To = %v, want now", tr.To)
	}
}

func TestTimeRange_Last6Hours(t *testing.T) {
	tr := Last6Hours()
	if tr.From != "now-6h" {
		t.Errorf("From = %v, want now-6h", tr.From)
	}
	if tr.To != "now" {
		t.Errorf("To = %v, want now", tr.To)
	}
}

func TestTimeRange_Last24Hours(t *testing.T) {
	tr := Last24Hours()
	if tr.From != "now-24h" {
		t.Errorf("From = %v, want now-24h", tr.From)
	}
	if tr.To != "now" {
		t.Errorf("To = %v, want now", tr.To)
	}
}

func TestTimeRange_Last7Days(t *testing.T) {
	tr := Last7Days()
	if tr.From != "now-7d" {
		t.Errorf("From = %v, want now-7d", tr.From)
	}
	if tr.To != "now" {
		t.Errorf("To = %v, want now", tr.To)
	}
}

func TestTimeRange_Last30Days(t *testing.T) {
	tr := Last30Days()
	if tr.From != "now-30d" {
		t.Errorf("From = %v, want now-30d", tr.From)
	}
	if tr.To != "now" {
		t.Errorf("To = %v, want now", tr.To)
	}
}

func TestTimeRange_Today(t *testing.T) {
	tr := Today()
	if tr.From != "now/d" {
		t.Errorf("From = %v, want now/d", tr.From)
	}
	if tr.To != "now/d" {
		t.Errorf("To = %v, want now/d", tr.To)
	}
}

func TestTimeRange_ThisWeek(t *testing.T) {
	tr := ThisWeek()
	if tr.From != "now/w" {
		t.Errorf("From = %v, want now/w", tr.From)
	}
	if tr.To != "now/w" {
		t.Errorf("To = %v, want now/w", tr.To)
	}
}

func TestTimeRange_ThisMonth(t *testing.T) {
	tr := ThisMonth()
	if tr.From != "now/M" {
		t.Errorf("From = %v, want now/M", tr.From)
	}
	if tr.To != "now/M" {
		t.Errorf("To = %v, want now/M", tr.To)
	}
}

func TestDefaultTimeRange(t *testing.T) {
	tr := DefaultTimeRange()
	if tr.From != "now-6h" {
		t.Errorf("From = %v, want now-6h", tr.From)
	}
	if tr.To != "now" {
		t.Errorf("To = %v, want now", tr.To)
	}
}
