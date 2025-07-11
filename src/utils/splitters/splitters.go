package splitters

func StringSlice(data []string, chunkSize int) [][]string {
	var result [][]string
	dataLength := len(data)
	for index := 0; index < dataLength; index += chunkSize {
		end := index + chunkSize
		if end > dataLength {
			end = dataLength
		}
		result = append(result, data[index:end])
	}

	return result
}
