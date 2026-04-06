package handler

import "sync"

var (
	lastParsedData []ParsedFlatQuestion
	mu             sync.RWMutex

	// Decision & Quiz store
	lastParsedDQData []ParsedFlatDQQuestion
	muDQ             sync.RWMutex
)

// SetLastParsedData saves the latest parsed output to memory safely
func SetLastParsedData(data []ParsedFlatQuestion) {
	mu.Lock()
	defer mu.Unlock()
	lastParsedData = data
}

// GetLastParsedData retrieves the latest parsed output from memory safely
func GetLastParsedData() []ParsedFlatQuestion {
	mu.RLock()
	defer mu.RUnlock()
	return lastParsedData
}

// SetLastParsedDQData saves the latest Decision & Quiz parsed output to memory safely
func SetLastParsedDQData(data []ParsedFlatDQQuestion) {
	muDQ.Lock()
	defer muDQ.Unlock()
	lastParsedDQData = data
}

// GetLastParsedDQData retrieves the latest Decision & Quiz parsed output from memory safely
func GetLastParsedDQData() []ParsedFlatDQQuestion {
	muDQ.RLock()
	defer muDQ.RUnlock()
	return lastParsedDQData
}
