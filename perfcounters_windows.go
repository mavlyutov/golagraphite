package golagraphite

import (
	"errors"
	"fmt"
	"github.com/lxn/win"
	"github.com/marpaia/graphite-golang"
	"log"
	"time"
	"unsafe"
)

func SendPerfCounters(c Config, metrics_channel chan []graphite.Metric) {

	if c.Performance_counters.Counters == nil {
		log.Println("Cannot find valid performance counters in config, will skip pc collection")
		return
	}

	for _, v := range c.Performance_counters.Counters {
		go func(counterName string) {
			for {	// Retry forever reading from conn counters
				// Handy when perfocunters reading is not yet available (soon after boot) or counter is temporarily missing
				ololosha, err := ReadPerformanceCounter(counterName, c.Performance_counters.Interval)
				if err == nil {
					for metrics := range ololosha {
						for i, v := range metrics {
							metrics[i].Name = fmt.Sprintf("%s.%s", c.Performance_counters.Metric_prefix, v.Name)
						}
						metrics_channel <- metrics
					}
				} else {
					log.Println(err)
				}
				time.Sleep(time.Duration(c.Performance_counters.Interval) * time.Second)
			}
		}(v)
	}
}

func ReadPerformanceCounter(counter string, sleepInterval int) (chan []graphite.Metric, error) {

	var queryHandle win.PDH_HQUERY
	var counterHandle win.PDH_HCOUNTER

	ret := win.PdhOpenQuery(0, 0, &queryHandle)
	if ret != win.ERROR_SUCCESS {
		return nil, errors.New("Unable to open query through DLL call")
	}

	// test path
	ret = win.PdhValidatePath(counter)
	if ret == win.PDH_CSTATUS_BAD_COUNTERNAME {
		return nil, errors.New("Unable to fetch counter (this is unexpected)")
	}

	ret = win.PdhAddEnglishCounter(queryHandle, counter, 0, &counterHandle)
	if ret != win.ERROR_SUCCESS {
		return nil, errors.New(fmt.Sprintf("Unable to add process counter. Error code is %x\n", ret))
	}

	ret = win.PdhCollectQueryData(queryHandle)
	if ret != win.ERROR_SUCCESS {
		return nil, errors.New(fmt.Sprintf("Got an error: 0x%x\n", ret))
	}

	out := make(chan []graphite.Metric)

	go func() {
		for {
			ret = win.PdhCollectQueryData(queryHandle)
			if ret == win.ERROR_SUCCESS {

				var metric []graphite.Metric

				var bufSize uint32
				var bufCount uint32
				var size uint32 = uint32(unsafe.Sizeof(win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE{}))
				var emptyBuf [1]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE // need at least 1 addressable null ptr.

				ret = win.PdhGetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &emptyBuf[0])
				if ret == win.PDH_MORE_DATA {
					filledBuf := make([]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE, bufCount*size)
					ret = win.PdhGetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &filledBuf[0])
					if ret == win.ERROR_SUCCESS {
						for i := 0; i < int(bufCount); i++ {
							c := filledBuf[i]
							s := win.UTF16PtrToString(c.SzName)

							metricName := NormalizeMetricName(counter)
							if len(s) > 0 {
								metricName = fmt.Sprintf("%s.%s", NormalizeMetricName(counter), NormalizeMetricName(s))
							}

							metric = append(metric, graphite.Metric{
								metricName,
								fmt.Sprintf("%v", c.FmtValue.DoubleValue),
								time.Now().Unix()})
						}
					}
				}
				out <- metric
			}

			time.Sleep(time.Duration(sleepInterval) * time.Second)
		}
	}()

	return out, nil

}
