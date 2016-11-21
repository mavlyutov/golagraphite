package golagraphite

import (
	"testing"
)

var metricNameTests = []struct {
	in  string
	out string
}{
	{"simple_metric", "simple_metric"},
	{" metric_with_spaces is_ugly ", "metric_with_spaces_is_ugly"},
	{"\\Processor(_Total)\\% Processor Time", "processor_total._processor_time"},
	{"\\PhysicalDisk(*)\\Avg. Disk Write Queue Length", "physicaldisk.avg__disk_write_queue_length"},
	{"\\Network Interface(*)\\Bytes Received/sec", "network_interface.bytes_receivedsec"},
}

func TestFlagParser(t *testing.T) {
	for _, test := range metricNameTests {
		normalized := NormalizeMetricName(test.in)
		if normalized != test.out {
			t.Errorf("expected %q from %q but got %q, ", test.out, test.in, normalized)
		}
	}
}
