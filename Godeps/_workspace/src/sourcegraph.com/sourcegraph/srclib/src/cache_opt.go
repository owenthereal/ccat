package src

type BuildCacheOpt struct {
	NoCacheRead  bool `long:"no-cache-read" description:"do not read from build cache"`
	NoCacheWrite bool `long:"no-cache-write" description:"do not write results to build cache"`
}
