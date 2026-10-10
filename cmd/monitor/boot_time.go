package main

import "time"

// bootTimeFromSample estimates the last system boot from the Agent's own
// timestamp and the Linux uptime counter. No additional Agent fields are
// required and existing node credentials remain unchanged.
func bootTimeFromSample(sample Sample) string {
 if sample.Timestamp.IsZero() || sample.Uptime==0 { return "" }
 // Guard against arbitrary out-of-range values in corrupted samples.
 if sample.Uptime > uint64((365*150)*24*3600) { return "" }
 return sample.Timestamp.Add(-time.Duration(sample.Uptime)*time.Second).UTC().Format(time.RFC3339)
}
