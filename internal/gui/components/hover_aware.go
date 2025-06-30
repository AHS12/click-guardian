package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// HoverAware is a container that shows a tooltip when the mouse is over it.

type HoverAware struct {
	widget.BaseWidget
	content      fyne.CanvasObject
	tooltip      *widget.Label
	tooltipText  func() string
}

// NewHoverAware creates a new HoverAware widget.
func NewHoverAware(content fyne.CanvasObject, tooltipText func() string) *HoverAware {
	h := &HoverAware{
		content:     content,
		tooltip:     widget.NewLabel(""),
		tooltipText: tooltipText,
	}
	h.tooltip.Hide()
	h.ExtendBaseWidget(h)
	return h
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (h *HoverAware) CreateRenderer() fyne.WidgetRenderer {
	return &hoverAwareRenderer{
		hoverAware: h,
	}
}

// MouseIn is called when the mouse enters the widget.
func (h *HoverAware) MouseIn(*desktop.MouseEvent) {
	h.tooltip.SetText(h.tooltipText())
	h.tooltip.Show()
}

// MouseOut is called when the mouse leaves the widget.
func (h *HoverAware) MouseOut() {
	h.tooltip.Hide()
}

// MouseMoved is called when the mouse moves over the widget.
func (h *HoverAware) MouseMoved(*desktop.MouseEvent) {}

type hoverAwareRenderer struct {
	hoverAware *HoverAware
}

func (r *hoverAwareRenderer) Layout(size fyne.Size) {
	r.hoverAware.content.Resize(size)
	r.hoverAware.tooltip.Resize(r.hoverAware.tooltip.MinSize())
	r.hoverAware.tooltip.Move(fyne.NewPos(size.Width, size.Height/2-r.hoverAware.tooltip.MinSize().Height/2))
}

func (r *hoverAwareRenderer) MinSize() fyne.Size {
	return r.hoverAware.content.MinSize()
}

func (r *hoverAwareRenderer) Refresh() {
	r.hoverAware.content.Refresh()
	r.hoverAware.tooltip.Refresh()
}

func (r *hoverAwareRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.hoverAware.content, r.hoverAware.tooltip}
}

func (r *hoverAwareRenderer) Destroy() {}
