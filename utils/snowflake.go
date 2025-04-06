package utils

import (
	"sync"
	"time"
)

type IdWorkerInterface interface {
	GenID() uint64
}

type IdWorker struct {
	workCount     uint64
	sequence      int64
	sequenceMask  int64
	lastTimestamp int64
	sync.Mutex
}

var _ IdWorkerInterface = (*IdWorker)(nil)

const (
	workerBits         int64 = 10
	sequenceBits       int64 = 12
	timestampLeftShift int64 = 41
	twEpochMagic       int64 = 1288834974657
)

func NewIdWorker(workCount uint64) *IdWorker {
	return &IdWorker{
		workCount:     workCount,
		sequenceMask:  -1 ^ (-1 << sequenceBits),
		lastTimestamp: -1,
	}
}

func (iw *IdWorker) GenID() uint64 {
	var timestamp = time.Now().UnixMilli()
	iw.Lock()
	defer iw.Unlock()
	//TODO:時鐘回播問題
	if timestamp < iw.lastTimestamp {
		timestamp = time.Now().UnixMilli()
	}
	if iw.lastTimestamp == timestamp {
		iw.sequence = (iw.sequence + 1) & iw.sequenceMask
		if iw.sequence == 0 {
			timestamp = tilNextMillis(iw.lastTimestamp)
		}
	} else {
		iw.sequence = 0
	}
	iw.lastTimestamp = timestamp
	return uint64((timestamp-twEpochMagic)<<timestampLeftShift |
		(int64(iw.workCount) << workerBits) |
		iw.sequence)
}

func tilNextMillis(lastTimestamp int64) int64 {

	timeStamp := time.Now().UnixMilli()
	for timeStamp <= lastTimestamp {
		timeStamp = time.Now().UnixMilli()
	}
	return timeStamp
}
