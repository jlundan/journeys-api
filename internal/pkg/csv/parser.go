package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"golang.org/x/text/encoding"
	"io"
	"os"
	"regexp"
	"strings"
)

func ParseFile(path string, firstLineAsHeaders bool) (map[string]uint8, [][]string, error) {
	return ParseFileWithDecoderAndDelimiter(path, firstLineAsHeaders, nil, ',')
}

func ParseFileWithDecoderAndDelimiter(path string, firstLineAsHeaders bool, decoder *encoding.Decoder, delimiter rune) (map[string]uint8, [][]string, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	var r *csv.Reader
	if decoder != nil {
		r = csv.NewReader(decoder.Reader(trimReader{csvFile}))
	} else {
		r = csv.NewReader(trimReader{csvFile})
	}

	if delimiter != ',' {
		r.Comma = delimiter
	}

	var headers = map[string]uint8{}
	var data = make([][]string, 0)

	headersRead := false
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("%v: %v", path, err.Error()))
		}

		if firstLineAsHeaders && !headersRead {
			for i, v := range record {
				headers[v] = uint8(i)
			}
			headersRead = true
		} else {
			var row []string
			for _, v := range record {
				row = append(row, strings.TrimSpace(v))
			}
			data = append(data, row)
		}

	}
	return headers, data, nil
}

var trailingWs = regexp.MustCompile(`\s\n`)

type trimReader struct{ io.Reader }

func (tr trimReader) Read(bs []byte) (int, error) {
	n, err := tr.Reader.Read(bs)
	if err != nil {
		return n, err
	}

	lines := string(bs[:n])
	trimmed := []byte(trailingWs.ReplaceAllString(lines, "\n"))
	copy(bs, trimmed)
	return len(trimmed), nil
}
