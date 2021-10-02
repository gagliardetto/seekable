package seekable

import (
	"fmt"
	"io"
	"sync"
)

type OffsetTracker struct {
	offsets []int
	len     int
	mu      *sync.RWMutex
}

func NewOffsetTracker() *OffsetTracker {
	return &OffsetTracker{
		mu: &sync.RWMutex{},
	}
}

func getOffsetsOfNewlines(r io.Reader, limit int) ([]int, error) {
	indexes := make([]int, 0)
	indexes = append(indexes, 0)
	buf := make([]byte, 32*1024)
	newline := '\n'

	fromStart := 0

	for {
		c, err := r.Read(buf)

		switch {
		case err == io.EOF:
			this := fromStart + c
			if indexes[len(indexes)-1] != this {
				indexes = append(indexes, this)
			}
			return indexes, nil

		case err != nil:
			return indexes, err
		}
		for i := 0; i < c; i++ {
			if rune(buf[i]) == newline {
				// If +1, then the newline is included in the offset.
				indexes = append(indexes, fromStart+i+1)

				if limit > 0 && len(indexes) > limit {
					return indexes, nil
				}
			}
		}
		fromStart += c
	}
}

//
func (tracker *OffsetTracker) CompileIndex(r io.Reader, limit int) error {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	var err error
	tracker.offsets, err = getOffsetsOfNewlines(r, limit)
	tracker.len = len(tracker.offsets)
	return err
}

func (tracker *OffsetTracker) validateLineNumber(lineNum int) error {
	if lineNum < 1 {
		return fmt.Errorf("lineNum out of bounds: got %v, but min is 1", lineNum)
	}
	maxLineNum := tracker.len - 1
	if lineNum > maxLineNum {
		return fmt.Errorf("lineNum out of bounds: got %v, but max is %v", lineNum, maxLineNum)
	}
	return nil
}

func (tracker *OffsetTracker) GetLine(r io.ReaderAt, lineNum int) ([]byte, error) {
	tracker.mu.RLock()
	defer tracker.mu.RUnlock()

	if err := tracker.validateLineNumber(lineNum); err != nil {
		return nil, err
	}

	i := lineNum - 1
	o := tracker.offsets[i]

	untilIndex := i
	if i < tracker.len-1 {
		untilIndex++
	}

	// The -1 is to exclude the final newline:
	l := tracker.offsets[untilIndex] - o

	if l <= 0 {
		// TODO
		return nil, nil
	}
	lineAt := make([]byte, l)
	if _, err := r.ReadAt(lineAt, int64(o)); err != nil {
		return nil, err
	}
	return lineAt, nil
}

func (tracker *OffsetTracker) GetLineReader(r io.ReaderAt, lineNum int) (io.Reader, error) {
	tracker.mu.RLock()
	defer tracker.mu.RUnlock()

	if err := tracker.validateLineNumber(lineNum); err != nil {
		return nil, err
	}

	i := lineNum - 1
	o := tracker.offsets[i]

	untilIndex := i
	if i < tracker.len-1 {
		untilIndex++
	}

	// The -1 is to exclude the final newline:
	l := tracker.offsets[untilIndex] - o

	if l <= 0 {
		// TODO:
		return nil, nil
	}
	return io.NewSectionReader(r, int64(o), int64(l)), nil
}

// func (tracker *OffsetTracker) AddOffset(offset int) error {
// 	tracker.mu.Lock()
// 	defer tracker.mu.Unlock()

// 	lastOffset := tracker.offsets[len(tracker.offsets)-1]
// 	if offset <= lastOffset {
// 		return errors.New("provided offset not valid; must be greater than the last registered offset")
// 	}

// 	tracker.offsets = append(tracker.offsets, offset)
// 	tracker.len++
// 	return nil
// }

// func (tracker *OffsetTracker) AddOffsetByLen(length int) error {
// 	tracker.mu.Lock()
// 	defer tracker.mu.Unlock()

// 	lastOffset := tracker.offsets[len(tracker.offsets)-1]

// 	offset := lastOffset + length

// 	tracker.offsets = append(tracker.offsets, offset)
// 	tracker.len++
// 	return nil
// }
