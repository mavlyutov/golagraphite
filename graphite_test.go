package golagraphite

import (
	"testing"
)

var metricNameTests = []struct {
	in  string
	out string
}{
	{"metric_a", "metric_a"},
	{" metric_a ", "metric_a"},
	{"\\Processor(_Total)\\% Processor Time", "processor_total._processor_time"},
}

func TestFlagParser(t *testing.T) {
	for _, test := range metricNameTests {
		normalized := NormalizeMetricName(test.in)
		if normalized != test.out {
			t.Errorf("expected %q from %q but got %q, ", test.out, test.in, normalized)
		}
	}
}
