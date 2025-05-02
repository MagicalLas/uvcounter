package core

type Counter interface {
	AddDelta(delta int64)
	AddDeltaAndGet(delta int64) int64
	Get() int64
}

type UVCounter interface {
	Count(viewerID string)
	CountAndGet(viewerID string) int64
	Get() int64
}
