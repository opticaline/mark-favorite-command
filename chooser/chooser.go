package chooser

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"strconv"
	"strings"
	"time"
)

const (
	KEY_C = 3
)

type RecordChooser struct {
	currentSelected int
	selected        map[int]bool
	records         []string
	isInput         bool
}

func Construct(records []string) *RecordChooser {
	return &RecordChooser{
		currentSelected: len(records) - 1,
		selected:        make(map[int]bool),
		records:         records,
		isInput:         false,
	}
}

func (rc *RecordChooser) WaitForAnswer() ([]string, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
	termbox.SetOutputMode(termbox.Output256 | termbox.OutputGrayscale)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	rc.draw()

	for {
		select {
		case ev := <-eventQueue:
			// Ctrl+c exit
			if ev.Type == termbox.EventKey && ev.Key == KEY_C {
				return nil, nil
			} else if rc.isInput {
				return rc.getSelected(), nil
			} else {
				switch {
				case ev.Ch == 'j':
					rc.addSelected(1)
				case ev.Ch == 'k':
					rc.addSelected(-1)
				case ev.Ch == 's':
					if rc.selected[rc.currentSelected] {
						delete(rc.selected, rc.currentSelected)
					} else {
						rc.selected[rc.currentSelected] = true
					}
				case ev.Key == 13:
					rc.isInput = true
				case ev.Type == termbox.EventKey:
					_, height := termbox.Size()
					switch {
					case ev.Key == 2:
						// page back
						rc.addSelected(-height)
					case ev.Key == 6:
						// page forward
						rc.addSelected(height)
					}
				default:
					fmt.Println(ev.Key, ev.Type, ev.Ch, ev.Mod, ev.N)
				}
			}
		default:
			rc.draw()
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (rc *RecordChooser) addSelected(num int) {
	foo := rc.currentSelected + num
	if foo < 0 {
		foo = 0
	}
	if foo >= len(rc.records) {
		foo = len(rc.records) - 1
	}
	rc.currentSelected = foo
}

func (rc *RecordChooser) draw() {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		fmt.Println(err)
	}
	if rc.isInput {

	} else {
		rc.drawList()
	}
	if err := termbox.Flush(); err != nil {
		fmt.Println(err)
	}
}

func (rc *RecordChooser) drawList() {
	width, height := termbox.Size()
	middleLine := height / 2
	start, end := 0, 0

	if rc.currentSelected <= middleLine {
		start = 0
		if len(rc.records) < height {
			end = len(rc.records) - 1
		} else {
			end = height
		}
	} else if rc.currentSelected+middleLine > len(rc.records) {
		end = len(rc.records)
		if end > height {
			start = end - height
		} else {
			end = 0
		}
	} else {
		start = rc.currentSelected - middleLine
		end = rc.currentSelected + middleLine
	}

	display := rc.records[start:end]
	length := len(fmt.Sprintf("%d", len(rc.records)+1)) + 1
	for i := 0; i < height; i++ {
		number := fmt.Sprintf("%"+strconv.Itoa(length)+"d ", i+start+1)
		color := termbox.ColorDefault
		if i+start == rc.currentSelected {
			color = termbox.ColorCyan
			number = fmt.Sprintf("%d", i+start+1)
			number = number + strings.Repeat(" ", length+1-len(number))
		}
		if rc.selected[i+start] {
			if i+start != rc.currentSelected {
				color = color | termbox.AttrBold
			}
			drawSentence(0, i, "*", termbox.ColorRed, termbox.ColorDefault)
		}
		drawSentence(1, i, number, termbox.ColorDefault, color)
		padding := width - len(display[i])
		if padding < 0 {
			padding = 0
		}
		drawSentence(length+3, i, fmt.Sprintf("%s%s", display[i], strings.Repeat(" ", padding)),
			termbox.ColorDefault, color)
	}
}

func (rc *RecordChooser) getSelected() []string {
	result := make([]string, len(rc.selected))
	index := 0
	for k := range rc.selected {
		result[index] = rc.records[k]
		index++
	}
	return result
}

func drawSentence(x, y int, sentence string, fg, bg termbox.Attribute) {
	for j := 0; j < len(sentence); j++ {
		termbox.SetCell(j+x, y, rune(sentence[j]), fg, bg)
	}
}
