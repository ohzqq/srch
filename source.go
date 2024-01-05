package srch

import (
	"os"
)

type Src func(args ...any) []any

func SliceSrc(data ...any) Src {
	return func(...any) []any {
		return data
	}
}

func FileSrc(file string) Src {
	return func(...any) []any {
		f, err := os.Open(file)
		if err != nil {
			return []any{}
		}
		defer f.Close()

		data, err := DecodeData(f)
		if err != nil {
			return []any{}
		}
		return data
	}
}
