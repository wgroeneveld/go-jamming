package common

import "time"

// I know it's public. Not sure how to handle this in tests, package-independent
var Now = time.Now
