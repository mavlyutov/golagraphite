package golagraphite

import (
	"testing"
)

var perfCounterMetricNameTests = []struct {
	in  string
	out string
}{
	{"\\Processor(_Total)\\% Processor Time", "processor_total._processor_time"},
	{"\\PhysicalDisk(*)\\Avg. Disk Write Queue Length", "physicaldisk.avg_disk_write_queue_length"},
	{"\\Network Interface(*)\\Bytes Received/sec", "network_interface.bytes_receivedsec"},
	{"\\myServer\\physicaldisk(1 e:)\\avg. disk write queue length", "myserver.physicaldisk1_e.avg_disk_write_queue_length"},
}

func TestNormalizePerfCounterMetricName(t *testing.T) {
	for _, test := range perfCounterMetricNameTests {
		normalized := normalizePerfCounterMetricName(test.in)
		if normalized != test.out {
			t.Errorf("expected %q from %q but got %q, ", test.out, test.in, normalized)
		}
	}
}
