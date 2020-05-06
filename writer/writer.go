package writer

import (
	"bufio"
	"os"
	"strings"
)

// Writer object is used to hold details of the markdown being written
type Writer struct {
	fileName string
	buffer   string
}

// New get a new writer object
func New(file string) *Writer {
	return &Writer{fileName: file}
}

// WriteFile writes markdown to disk
func (writer *Writer) WriteFile() error {
	f, err := os.Create(writer.fileName)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	_, werr := w.WriteString(writer.buffer)
	if werr != nil {
		return werr
	}
	w.Flush()
	f.Close()
	return nil
}

// MakeLink makes links to include in markdown files
func MakeLink(title string, url string) string {
	return "[" + title + "](" + url + ")"
}

// InlineCode returns code in backticks
func InlineCode(code string) string {
	return "`" + code + "`"
}

// GetBuf returns the buffer from memory
func (writer *Writer) GetBuf() string {
	return writer.buffer
}

func (writer *Writer) writeLine(line string) {
	writer.buffer = writer.buffer + line + "\n"
}

// H1 adds a h1 to the file
func (writer *Writer) H1(line string) {
	writer.writeLine("# " + line)
	writer.writeLine("")
}

// H1Underline adds a h1 with underline to the file
func (writer *Writer) H1Underline(line string) {
	writer.writeLine(line)
	writer.writeLine("======")
	writer.writeLine("")
}

// H2 adds a h2 to the file
func (writer *Writer) H2(line string) {
	writer.writeLine("## " + line)
	writer.writeLine("")
}

// H2Underline adds a h2 with underline to the file
func (writer *Writer) H2Underline(line string) {
	writer.writeLine(line)
	writer.writeLine("------")
	writer.writeLine("")
}

// H3 adds a h3 to the file
func (writer *Writer) H3(line string) {
	writer.writeLine("### " + line)
	writer.writeLine("")
}

// P adds a block of text
func (writer *Writer) P(line string) {
	writer.writeLine("")
	writer.writeLine(line)
	writer.writeLine("")
}

// Bullet creates a bullet point
func (writer *Writer) Bullet(line string) {
	writer.writeLine("* " + line)
}

// Table creates a table
func (writer *Writer) Table(headers []string, rows [][]string) {
	var headerLine string
	headerLine = "|" + strings.Join(headers, " | ") + "|"
	writer.writeLine(headerLine)
	hyphens := make([]string, len(headers))
	var hyphensLine string
	for i := 0; i < len(headers); i++ {
		hyphens[i] = "---"
	}
	hyphensLine = strings.Join(hyphens, " | ")
	writer.writeLine(hyphensLine)
	for _, row := range rows {
		rowLine := strings.Join(row, " | ")
		writer.writeLine(rowLine)
	}
	writer.writeLine("")
}
