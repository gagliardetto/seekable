package seekable

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func _WriteUint32LE(writer io.Writer, i uint32) (err error) {
	return _WriteUint32(writer, i, binary.LittleEndian)
}

func _WriteUint32(writer io.Writer, i uint32, order binary.ByteOrder) (err error) {
	buf := make([]byte, 4)
	order.PutUint32(buf, i)
	_, err = writer.Write(buf)
	return err
}

func _ReadUint32LE(reader io.Reader) (out uint32, err error) {
	return _ReadUint32(reader, binary.LittleEndian)
}

func _ReadUint32(reader io.Reader, order binary.ByteOrder) (out uint32, err error) {
	buf := make([]byte, 4)
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, fmt.Errorf("expected 4 bytes, got %v", n)
	}
	out = order.Uint32(buf)
	return
}

func TestRegisterByLen(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	tracker := NewOffsetTracker()
	items := [][]byte{
		[]byte("foo"),
		[]byte("bar"),
		[]byte("hello"),
		[]byte("world"),
	}
	{
		for itemIndex, item := range items {
			ln := len(item)

			err := _WriteUint32LE(buf, uint32(ln))
			require.NoError(t, err)

			gotLn, err := buf.Write(item)
			require.NoError(t, err)
			require.Equal(t, ln, gotLn)

			lineNum, err := tracker.RegisterByLen(4 + ln)
			require.NoError(t, err)
			require.Equal(t, itemIndex+1, lineNum)
		}
		spew.Dump(tracker)
	}
	reader := bytes.NewReader(buf.Bytes())
	{
		require.Equal(t, []int{0, 7, 14, 23, 32}, tracker.offsets)
	}
	{
		for i := 0; i < len(items)-1; i++ {
			lineNum := i + 1
			got, err := tracker.GetLine(reader, lineNum)
			require.NoError(t, err)
			lineReader := bytes.NewReader(got)
			gotContentLen, err := _ReadUint32LE(lineReader)
			require.NoError(t, err)
			content, err := io.ReadAll(lineReader)
			require.NoError(t, err)
			require.Equal(t, len(content), int(gotContentLen))
			require.Equal(t, items[lineNum-1], content)
			require.Equal(t, items[i], content)
		}
	}
}

func TestNewOffsetTracker(t *testing.T) {
	{
		reader := strings.NewReader("\n\n")
		tracker := NewOffsetTracker()
		err := tracker.CompileIndex(reader, 0)
		require.NoError(t, err)
		{
			require.Equal(t, []int{0, 1, 2}, tracker.offsets)
		}
		{
			got, err := tracker.GetLine(reader, 1)
			require.NoError(t, err)
			require.Equal(t, []byte("\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 2)
			require.NoError(t, err)
			require.Equal(t, []byte("\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 12)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, 0)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, -100)
			require.Error(t, err)
			require.Nil(t, got)
		}
	}
	{
		reader := strings.NewReader("\nhello\n\n\nworld\n")
		tracker := NewOffsetTracker()
		err := tracker.CompileIndex(reader, 0)
		require.NoError(t, err)
		{
			require.Equal(t, []int{0, 1, 7, 8, 9, 15}, tracker.offsets)
		}
		{
			got, err := tracker.GetLine(reader, 1)
			require.NoError(t, err)
			require.Equal(t, []byte("\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 2)
			require.NoError(t, err)
			require.Equal(t, []byte("hello\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 3)
			require.NoError(t, err)
			require.Equal(t, []byte("\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 4)
			require.NoError(t, err)
			require.Equal(t, []byte("\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 5)
			require.NoError(t, err)
			require.Equal(t, []byte("world\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 12)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, 0)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, -100)
			require.Error(t, err)
			require.Nil(t, got)
		}
	}
	{
		reader := strings.NewReader("hello\n\n\nworld\n")
		tracker := NewOffsetTracker()
		err := tracker.CompileIndex(reader, 0)
		require.NoError(t, err)
		{
			require.Equal(t, []int{0, 6, 7, 8, 14}, tracker.offsets)
		}
		{
			got, err := tracker.GetLine(reader, 1)
			require.NoError(t, err)
			require.Equal(t, []byte("hello\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 2)
			require.NoError(t, err)
			require.Equal(t, []byte("\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 3)
			require.NoError(t, err)
			require.Equal(t, []byte("\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 4)
			require.NoError(t, err)
			require.Equal(t, []byte("world\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 12)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, 0)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, -100)
			require.Error(t, err)
			require.Nil(t, got)
		}
	}
	{
		reader := strings.NewReader("hello\nworld\n")
		tracker := NewOffsetTracker()
		err := tracker.CompileIndex(reader, 0)
		require.NoError(t, err)
		{
			require.Equal(t, []int{0, 6, 12}, tracker.offsets)
		}
		{
			got, err := tracker.GetLine(reader, 1)
			require.NoError(t, err)
			require.Equal(t, []byte("hello\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 2)
			require.NoError(t, err)
			require.Equal(t, []byte("world\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 3)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, 0)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, -100)
			require.Error(t, err)
			require.Nil(t, got)
		}
	}
	{
		reader := strings.NewReader("hello\nworld")
		tracker := NewOffsetTracker()
		err := tracker.CompileIndex(reader, 0)
		require.NoError(t, err)
		{
			require.Equal(t, []int{0, 6, 11}, tracker.offsets)
		}
		{
			got, err := tracker.GetLine(reader, 1)
			require.NoError(t, err)
			require.Equal(t, []byte("hello\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 2)
			require.NoError(t, err)
			require.Equal(t, []byte("world"), got)
		}
		{
			got, err := tracker.GetLine(reader, 3)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, 0)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, -100)
			require.Error(t, err)
			require.Nil(t, got)
		}
	}
	{
		reader := strings.NewReader("hello\n")
		tracker := NewOffsetTracker()
		err := tracker.CompileIndex(reader, 0)
		require.NoError(t, err)
		{
			require.Equal(t, []int{0, 6}, tracker.offsets)
		}
		{
			got, err := tracker.GetLine(reader, 1)
			require.NoError(t, err)
			require.Equal(t, []byte("hello\n"), got)
		}
		{
			got, err := tracker.GetLine(reader, 3)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, 0)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, -100)
			require.Error(t, err)
			require.Nil(t, got)
		}
	}
	{
		reader := strings.NewReader("hello")
		tracker := NewOffsetTracker()
		err := tracker.CompileIndex(reader, 0)
		require.NoError(t, err)
		{
			require.Equal(t, []int{0, 5}, tracker.offsets)
		}
		{
			got, err := tracker.GetLine(reader, 1)
			require.NoError(t, err)
			require.Equal(t, []byte("hello"), got)
		}
		{
			got, err := tracker.GetLine(reader, 3)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, 0)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, -100)
			require.Error(t, err)
			require.Nil(t, got)
		}
	}
	{
		reader := strings.NewReader("")
		tracker := NewOffsetTracker()
		err := tracker.CompileIndex(reader, 0)
		require.NoError(t, err)
		{
			require.Equal(t, []int{0}, tracker.offsets)
		}
		{
			got, err := tracker.GetLine(reader, 1)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, 3)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, 0)
			require.Error(t, err)
			require.Nil(t, got)
		}
		{
			got, err := tracker.GetLine(reader, -100)
			require.Error(t, err)
			require.Nil(t, got)
		}
	}
}

func TestGetLineReader(t *testing.T) {
	{
		reader := strings.NewReader("\nhello\n\n\nworld\n")
		tracker := NewOffsetTracker()
		err := tracker.CompileIndex(reader, 0)
		require.NoError(t, err)
		{
			require.Equal(t, []int{0, 1, 7, 8, 9, 15}, tracker.offsets)
		}
		{
			gotReader, err := tracker.GetLineReader(reader, 1)
			require.NoError(t, err)
			got, err := ioutil.ReadAll(gotReader)
			require.NoError(t, err)
			require.Equal(t, []byte("\n"), got)
		}
		{
			gotReader, err := tracker.GetLineReader(reader, 2)
			require.NoError(t, err)
			got, err := ioutil.ReadAll(gotReader)
			require.NoError(t, err)
			require.Equal(t, []byte("hello\n"), got)
		}
		{
			gotReader, err := tracker.GetLineReader(reader, 3)
			require.NoError(t, err)
			got, err := ioutil.ReadAll(gotReader)
			require.NoError(t, err)
			require.Equal(t, []byte("\n"), got)
		}
		{
			gotReader, err := tracker.GetLineReader(reader, 4)
			require.NoError(t, err)
			got, err := ioutil.ReadAll(gotReader)
			require.NoError(t, err)
			require.Equal(t, []byte("\n"), got)
		}
		{
			gotReader, err := tracker.GetLineReader(reader, 5)
			require.NoError(t, err)
			got, err := ioutil.ReadAll(gotReader)
			require.NoError(t, err)
			require.Equal(t, []byte("world\n"), got)
		}
		{
			gotReader, err := tracker.GetLineReader(reader, 12)
			require.Error(t, err)
			require.Nil(t, gotReader)
		}
		{
			gotReader, err := tracker.GetLineReader(reader, 0)
			require.Error(t, err)
			require.Nil(t, gotReader)
		}
		{
			gotReader, err := tracker.GetLineReader(reader, -100)
			require.Error(t, err)
			require.Nil(t, gotReader)
		}
	}
}
