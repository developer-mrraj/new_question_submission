package handler

import "sync"

var (
	lastParsedData []ParsedFlatQuestion
	mu             sync.RWMutex
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
