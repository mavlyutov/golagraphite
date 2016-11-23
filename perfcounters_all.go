package golagraphite

import "strings"

func normalizePerfCounterMetricName(rawName string) (normalizedName string) {

	normalizedName = rawName

	// thanks to Microsoft Windows,
	// we have performance counter metric like `\\Processor(_Total)\\% Processor Time`
	// which we need to convert to `processor_total.processor_time` see perfcounter_test.go for more beautiful examples
	r := strings.NewReplacer(
		".", "",
		"\\", ".",
		" ", "_",
	)
	normalizedName = r.Replace(normalizedName)

	normalizedName = NormalizeMetricName(normalizedName)
	return
}
