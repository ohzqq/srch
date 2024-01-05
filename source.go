package srch

type Src func(args ...any) []any

func SliceSrc(data ...any) Src {
	return func(...any) []any {
		return data
	}
}

func FileSrc(file ...string) Src {
	return func(...any) []any {
		data, err := NewDataFromFiles(file...)
		if err != nil {
			return []any{}
		}
		return data
	}
}
