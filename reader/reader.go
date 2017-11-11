package reader

import (
	"bufio"
	"io"
)

// Run start buffered reader
func Run(reader io.Reader, writerCloser io.WriteCloser) (err error) {

	w := bufio.NewWriter(writerCloser)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		buf := scanner.Bytes()

		if _, err = w.WriteString("-> "); err != nil {
			return err
		}

		if _, err = w.Write(buf); err != nil {
			return err
		}

		if _, err = w.WriteRune('\n'); err != nil {
			return err
		}

		if err = w.Flush(); err != nil {
			return err
		}
	}

	if err = writerCloser.Close(); err != nil {
		return err
	}

	return scanner.Err()
}
