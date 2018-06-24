package core

import (
	"fmt"
	"github.com/zaypen/jbt/screen"
	"github.com/zaypen/jbt/strong"
	"github.com/zaypen/jbt/unit"
	"strings"
)

type ProgressWriter struct {
	screenWidth    int
	current, total uint64
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	inc := len(p)
	pw.current += uint64(inc)
	pw.printProgress()
	return inc, nil
}

func (pw *ProgressWriter) printProgress() {
	const barBorderLength = 2 // len("[]") = 2
	const minBarLength = 1    // len("-") = 1
	const maxStatsLength = 26 // len("1234.0XB / 1234.0XB 100.0%") = 26
	pw.Clear()
	percentage := float32(pw.current) / float32(pw.total)
	downloaded := unit.Bytes(pw.current)
	remaining := unit.Bytes(pw.total)
	stats := fmt.Sprintf("%s / %s, %.1f%%", downloaded, remaining, percentage*100)
	if barLength := pw.screenWidth - maxStatsLength - barBorderLength; barLength >= minBarLength {
		filling := int(float32(barLength) * percentage)
		bar := fmt.Sprintf("[%s%s]", strings.Repeat("=", filling), strings.Repeat("-", barLength-filling))
		fmt.Printf("%s%s", bar, strong.Pad(stats, maxStatsLength))
	} else {
		fmt.Printf("%s", stats)
	}
}
func (pw *ProgressWriter) Clear() {
	fmt.Printf("\r%s\r", strings.Repeat(" ", pw.screenWidth))
}

func newProgressWriter(total uint64) *ProgressWriter {
	width, _ := screen.Size()
	return &ProgressWriter{width, 0, total}
}
