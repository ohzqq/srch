package data

import (
	"bufio"
	"encoding/json"
	"io"
)

type Table struct {
	rs      io.ReadSeeker
	w       io.Writer
	offsets map[int]int64
}

func NewTable(seeker io.ReadSeeker, w io.Writer) (*Table, error) {
	offsets, err := DecodeHareTable(seeker)
	if err != nil {
		return nil, err
	}

	tableFile := Table{
		rs:      seeker,
		w:       w,
		offsets: offsets,
	}

	return &tableFile, nil
}

func DecodeHareTable(s io.ReadSeeker) (map[int]int64, error) {

	offsets := make(map[int]int64)
	var totalOffset int64
	var recLen int
	var recMap map[string]interface{}
	var currentOffset int64

	r := bufio.NewReader(s)

	_, err := s.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	for {
		rec, err := r.ReadBytes('\n')

		recLen = len(rec)
		totalOffset += int64(recLen)
		currentOffset = totalOffset

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// Skip dummy records.
		if (rec[0] == '\n') || (rec[0] == dummyRune) {
			continue
		}

		//Unmarshal so we can grab the record ID.
		if err := json.Unmarshal(rec, &recMap); err != nil {
			return nil, err
		}
		recMapID := int(recMap["id"].(float64))

		//println(string(rec))
		offsets[recMapID] = currentOffset
	}

	return offsets, nil
}
