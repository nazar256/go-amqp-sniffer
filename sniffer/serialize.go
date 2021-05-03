package sniffer

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"
)

type serializer func(r *OutputRecord) ([]byte, error)

func getSerializer(format Format) serializer {
	switch format {
	case JSON:
		return serializeJSON
	case CSV:
		return serializeCsv
	}

	panic(fmt.Sprintf("Specified format %d is unsupported. It's a bug", format))
}

func serializeJSON(r *OutputRecord) ([]byte, error) {
	jsonRecord, err := json.Marshal(*r)
	return append(jsonRecord, '\n'), err
}

func serializeCsv(r *OutputRecord) ([]byte, error) {
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	timeFormatted := r.Timestamp.Format(time.RFC3339Nano)

	headers, err := json.Marshal(r.Headers)
	if err != nil {
		return buffer.Bytes(), err
	}

	err = writer.Write([]string{
		r.MessageID,
		r.AppID,
		timeFormatted,
		string(headers),
		r.Payload,
		r.ContentType,
		r.ContentEncoding,
	})

	writer.Flush()

	return buffer.Bytes(), err
}
