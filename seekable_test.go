package seekable

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
