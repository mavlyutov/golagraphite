package golagraphite

import (
	"testing"
)

var metricNameTests = []struct {
	in  string
	out string
}{
	{"simple_metric", "simple_metric"},
	{"this.is.hostname", "this.is.hostname"},
	{" metric_with_spaces	is_ugly ", "metric_with_spaces_is_ugly"},
	{"..metric_with.dots.", "metric_with.dots"},
	{"metric_with_some.±!@#[]strange_symbols!@#$%^&*()_+", "metric_with_some.strange_symbols"},
}

func TestNormalizeMetricName(t *testing.T) {
	for _, test := range metricNameTests {
		normalized := NormalizeMetricName(test.in)
		if normalized != test.out {
			t.Errorf("expected %q from %q but got %q, ", test.out, test.in, normalized)
		}
	}
}
