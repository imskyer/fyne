package widget

import (
	"bufio"
	"image/color"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

type labelRenderer struct {
	objects []fyne.CanvasObject

	texts []*canvas.Text

	label *Label
	lines int
}

func (l *labelRenderer) parseText(text string) []string {
	if !strings.Contains(text, "\n") {
		return []string{text}
	}

	var texts []string
	s := bufio.NewScanner(strings.NewReader(text))
	for s.Scan() {
		texts = append(texts, s.Text())
	}
	// this checks if Scan() ended on a blank line
	if string(text[len(text)-1]) == "\n" {
		texts = append(texts, "")
	}

	return texts
}

func (l *labelRenderer) updateTexts(strings []string) {
	l.lines = len(strings)
	count := len(l.texts)
	layout := false

	for i, str := range strings {
		if i >= count {
			text := canvas.NewText("", theme.TextColor())
			l.texts = append(l.texts, text)
			l.objects = append(l.objects, text)

			layout = true
		}
		l.texts[i].Text = str
	}

	for i := l.lines; i < len(l.texts); i++ {
		l.texts[i].Text = ""
		layout = true
	}

	l.Refresh()
	if layout {
		// TODO invalidate container size (to shrink)
		l.Layout(l.label.Size())
	}
}

// MinSize calculates the minimum size of a label.
// This is based on the contained text with a standard amount of padding added.
func (l *labelRenderer) MinSize() fyne.Size {
	height := 0
	width := 0
	for i := 0; i < l.lines; i++ {
		min := l.texts[i].MinSize()
		if l.texts[i].Text == "" {
			min = emptyTextMinSize(l.label.TextStyle)
		}
		height += min.Height
		width = fyne.Max(width, min.Width)
	}

	return fyne.NewSize(width, height).Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
}

func (l *labelRenderer) Layout(size fyne.Size) {
	yPos := theme.Padding()
	lineHeight := emptyTextMinSize(l.label.TextStyle).Height
	lineSize := fyne.NewSize(size.Width-theme.Padding()*2, lineHeight)
	for i := 0; i < l.lines; i++ {
		text := l.texts[i]
		text.Resize(lineSize)
		text.Move(fyne.NewPos(theme.Padding(), yPos))
		yPos += lineHeight
	}
}

func (l *labelRenderer) Objects() []fyne.CanvasObject {
	return l.objects
}

// ApplyTheme is called when the Label may need to update it's look
func (l *labelRenderer) ApplyTheme() {
	for _, text := range l.texts {
		if l.label.color != nil {
			text.Color = l.label.color
		} else {
			text.Color = theme.TextColor()
		}
	}
}

func (l *labelRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (l *labelRenderer) Refresh() {
	l.ApplyTheme()
	for _, text := range l.texts {
		text.Alignment = l.label.Alignment
		text.TextStyle = l.label.TextStyle
	}

	canvas.Refresh(l.label)
}

// Label widget is a basic text component with appropriate padding and layout.
type Label struct {
	baseWidget

	Text      string         // The content of the label
	Alignment fyne.TextAlign // The alignment of the Text
	TextStyle fyne.TextStyle // The style of the label text

	color color.Color
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (l *Label) Resize(size fyne.Size) {
	l.resize(size, l)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (l *Label) Move(pos fyne.Position) {
	l.move(pos, l)
}

// MinSize returns the smallest size this widget can shrink to
func (l *Label) MinSize() fyne.Size {
	return l.minSize(l)
}

// Show this widget, if it was previously hidden
func (l *Label) Show() {
	l.show(l)
}

// Hide this widget, if it was previously visible
func (l *Label) Hide() {
	l.hide(l)
}

// SetText updates the text of the label widget
func (l *Label) SetText(text string) {
	l.Text = text

	render := Renderer(l).(*labelRenderer)
	for _, obj := range render.texts {
		obj.Text = ""
	}

	render.updateTexts(render.parseText(l.Text))
}

// Rows returns the number of text rows in this text entry.
// The entry may be longer than required to show this amount of content.
func (l *Label) Rows() int {
	return Renderer(l).(*labelRenderer).lines
}

// RowLength returns the number of visible characters in the row specified.
// The row parameter should be between 0 and l.Rows()-1.
func (l *Label) RowLength(row int) int {
	return len(Renderer(l).(*labelRenderer).texts[row].Text)
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (l *Label) CreateRenderer() fyne.WidgetRenderer {
	render := &labelRenderer{label: l}

	render.texts = []*canvas.Text{}
	render.objects = []fyne.CanvasObject{}
	render.updateTexts(render.parseText(l.Text))

	return render
}

// NewLabel creates a new layout widget with the set text content
func NewLabel(text string) *Label {
	var style fyne.TextStyle

	l := &Label{
		baseWidget{},
		text,
		fyne.TextAlignLeading,
		style,
		nil,
	}

	Renderer(l).Layout(l.MinSize())
	return l
}
