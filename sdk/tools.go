package sdk

func createChunk[T any](newElements *[]T, chunckSize int) *[][]T {
	var chuncks [][]T

	if newElements == nil || len(*newElements) == 0 {
		return nil
	}

	for i := 0; i < len(*newElements); i += chunckSize {
		chunk := (*newElements)[i:min(i+chunckSize, len(*newElements))]
		chuncks = append(chuncks, chunk)
	}

	return &chuncks
}
