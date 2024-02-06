package radio

import (
	"fmt"
	"net/http"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

type Widget struct {
	view.ScrollableWidget

	jobs  []Job
	err   error
	index int
}

func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := &Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.Common),

		index: -1,
	}

	widget.SetRenderFunction(widget.Render)
	widget.initializeKeyboardControls()

	return widget
}

func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	resp, err := http.Get("https://radio.serveo.net/api/list")
	if err != nil {
		widget.err = err
		widget.jobs = nil
		widget.SetItemCount(0)
	} else {
		var data struct {
			Jobs struct {
				Jobs []Job `json:"jobs"`
			} `json:"jobs"`
		}

		if err := utils.ParseJSON(resp.Body, &data); err != nil {
			widget.err = err
			widget.jobs = nil
			widget.SetItemCount(0)
		} else {
			widget.jobs = data.Jobs.Jobs
			widget.SetItemCount(len(widget.jobs))
		}
	}

	widget.Render()
}

func (widget *Widget) Render() {
	widget.Redraw(widget.content)
}

func (widget *Widget) content() (string, string, bool) {
	title := "Radio Jobs"

	if widget.err != nil {
		return title, widget.err.Error(), true
	}

	if len(widget.jobs) == 0 {
		return title, "No jobs to display", false
	}

	var str string
	for idx, job := range widget.jobs {
		row := fmt.Sprintf(
			`[%s]%2d. %s`,
			widget.RowColor(idx),
			idx+1,
			job.Name,
		)

		str += utils.HighlightableHelper(widget.View, row, idx, len(job.Name))
	}

	return title, str, false
}

func (widget *Widget) startJob() {
	job := widget.selectedJob()
	if job != nil {
		utils.OpenURL(job.StartURL)
	}
}

func (widget *Widget) burnJob() {
	job := widget.selectedJob()
	if job != nil {
		utils.OpenURL(job.BurnURL)
	}
}

func (widget *Widget) selectedJob() *Job {
	var job *Job

	sel := widget.GetSelected()
	if sel >= 0 && widget.jobs != nil && sel < len(widget.jobs) {
		job = &widget.jobs[sel]
	}

	return job
}