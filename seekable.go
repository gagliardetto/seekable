package seekable

import (
	"fmt"
	"io"
	"sync"
)

type OffsetTracker struct {
	offsets  []int
	numLines int
	mu       *sync.RWMutex
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
	tracker.numLines = len(tracker.offsets) - 1
	return err
}

func (tracker *OffsetTracker) validateLineNumber(lineNum int) error {
	if lineNum < 1 {
		return fmt.Errorf("lineNum out of bounds: got %v, but min is 1", lineNum)
	}
	maxLineNum := tracker.numLines
	if lineNum > maxLineNum {
		return fmt.Errorf("lineNum out of bounds: got %v, but max is %v", lineNum, maxLineNum)
	}
	return nil
}

// GetLine returns the the provided line bytes; lineNum is 1-based.
func (tracker *OffsetTracker) GetLine(r io.ReaderAt, lineNum int) ([]byte, error) {
	tracker.mu.RLock()
	defer tracker.mu.RUnlock()

	o, l, err := tracker.getOffset(lineNum)
	if err != nil {
		return nil, err
	}
	lineAt := make([]byte, l)
	if _, err := r.ReadAt(lineAt, int64(o)); err != nil {
		return nil, err
	}
	return lineAt, nil
}

// get offset and length for the provided line; lineNum is 1-based
func (tracker *OffsetTracker) getOffset(lineNum int) (int, int, error) {
	if err := tracker.validateLineNumber(lineNum); err != nil {
		return 0, 0, err
	}

	i := lineNum - 1
	o := tracker.offsets[i]

	untilIndex := i
	if i < tracker.numLines {
		untilIndex++
	}

	l := tracker.offsets[untilIndex] - o

	if l <= 0 {
		// TODO
		return 0, 0, fmt.Errorf("l is %v", l)
	}

	return o, l, nil
}

// GetLineReader returns a reader for the provided line; lineNum is 1-based.
func (tracker *OffsetTracker) GetLineReader(r io.ReaderAt, lineNum int) (io.Reader, error) {
	tracker.mu.RLock()
	defer tracker.mu.RUnlock()

	o, l, err := tracker.getOffset(lineNum)
	if err != nil {
		return nil, err
	}
	return io.NewSectionReader(r, int64(o), int64(l)), nil
}

// RegisterByLen will register a new item of given length.
func (tracker *OffsetTracker) RegisterByLen(length int) (int, error) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	var lastOffset int
	if len(tracker.offsets) > 0 {
		lastOffset = tracker.offsets[len(tracker.offsets)-1]
	} else {
		tracker.offsets = append(tracker.offsets, 0)
	}

	offset := lastOffset + length

	tracker.offsets = append(tracker.offsets, offset)
	tracker.numLines++
	return tracker.numLines, nil
}

func (tracker *OffsetTracker) NumItems() int {
	tracker.mu.RLock()
	defer tracker.mu.RUnlock()

	return tracker.numLines
}
