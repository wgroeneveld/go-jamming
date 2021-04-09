package common

import "time"

// https://labs.yulrizka.com/en/stubbing-time-dot-now-in-golang/
// None of the above are very appealing. For now, just use the lazy way.
var Now = time.Now
