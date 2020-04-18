package render

import (
	"encoding/json"
	"io"
)

type JSONRendererOptions struct {
	Writer      io.Writer
	Columns     []string
	OnlyEnabled bool
}

type JSONRenderer struct {
	w       io.Writer
	Columns []string
	data    *data
}

type data struct {
	lines []line
}

func NewJSONRenderer(opts *JSONRendererOptions) JSONRenderer {
	if len(opts.Columns) == 0 {
		opts.Columns = DefaultColumns
	}

	return JSONRenderer{
		w:       opts.Writer,
		Columns: opts.Columns,
		data:    &data{},
	}
}

func (j JSONRenderer) AddSeparator() {
	// not used
}

type line struct {
	Profile string
	Status  string
	IP      string
	Host    string
}

func (j JSONRenderer) AppendRow(row *Row) {
	if row.Comment != "" {
		return
	}

	l := line{
		Profile: row.Profile,
		Status:  row.Status,
		IP:      row.IP,
		Host:    row.Host,
	}
	j.data.lines = append(j.data.lines, l)
}

func (j JSONRenderer) Render() error {
	enc := json.NewEncoder(j.w)

	return enc.Encode(j.data.lines)
}
