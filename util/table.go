package util

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type Table interface {
	Header(...string)
	Row(...string)
	End()
}

type table struct {
	tw *tabwriter.Writer
}

func (t table) Header(headers ...string) {
	sep := make([]string, len(headers))
	for i, v := range headers {
		sep[i] = strings.Repeat("-", len(v))
	}
	t.Row(headers...)
	t.Row(sep...)
}

func (t table) Row(fields ...string) {
	fmt.Fprintln(t.tw, strings.Join(fields, "\t"))
}

func (t table) End() {
	t.tw.Flush()
}

func NewTable() Table {
	return &table{tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)}
}
