package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockStdin(t *testing.T, dummyInput string) (funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin

	tmpfile, err := os.CreateTemp(t.TempDir(), t.Name())
	if err != nil {
		return nil, err
	}

	content := []byte(dummyInput)

	if _, err := tmpfile.Write(content); err != nil {
		return nil, err
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return nil, err
	}

	// Set stdin to the temp file
	os.Stdin = tmpfile

	return func() {
		// clean up
		os.Stdin = oldOsStdin
		os.Remove(tmpfile.Name())
	}, nil
}

func TestPipe(t *testing.T) {
	o := initOls()
	o.Opt.Quiet = true
	userInput := " some ip:  8.8.8.8 \n some url  hTTps://www.google.com/something \n"
	funcDefer, err := mockStdin(t, userInput)
	if err != nil {
		t.Fatal(err)
	}

	defer funcDefer()

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	o.readFromPipe()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	fmt.Println(out)

	l := strings.Split(out, "\n")
	fmt.Println(l[2])
	fmt.Println(l[4])
	fmt.Println(l[8])
	assert.Equal(t, " some url  hTTps://www.google.com/something ", l[2], "test url")
	assert.Equal(t, "URL: https://www.google.com", l[4], "test url")
	assert.Equal(t, "  creation: 1997-09-15T04:00:00Z", l[12], "google creation")
}
