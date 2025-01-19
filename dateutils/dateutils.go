package dateutils

import (
	"sort"
	"time"
)

// HashWithTimestamp represents a hash and its corresponding timestamp.
type HashWithTimestamp struct {
	Hash      string
	Timestamp int64 // Assuming timestamp is in Unix epoch format
}

// GroupHashesByDate groups hashes by their date, starting from today to older dates.
func GroupHashesByDate(data []HashWithTimestamp) map[string][]string {
	// Step 1: Parse timestamps and filter data starting from today
	today := time.Now().Truncate(24 * time.Hour) // Get the start of today
	var filteredData []HashWithTimestamp
	for _, item := range data {
		itemTime := time.Unix(item.Timestamp, 0)
		if itemTime.After(today) || itemTime.Equal(today) {
			filteredData = append(filteredData, item)
		}
	}

	// Step 2: Group data by date
	groupedData := make(map[string][]string)
	for _, item := range filteredData {
		itemTime := time.Unix(item.Timestamp, 0)
		dateKey := itemTime.Format("2006-01-02") // Format as YYYY-MM-DD
		groupedData[dateKey] = append(groupedData[dateKey], item.Hash)
	}

	// Step 3: Sort the groups in descending order (from today to older dates)
	var sortedDates []string
	for date := range groupedData {
		sortedDates = append(sortedDates, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(sortedDates)))

	// Step 4: Create a sorted map for the grouped data
	sortedGroupedData := make(map[string][]string)
	for _, date := range sortedDates {
		sortedGroupedData[date] = groupedData[date]
	}

	return sortedGroupedData
}
