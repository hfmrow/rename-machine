// progress.go

/*
	Source file auto-generated on Thu, 17 Oct 2019 02:33:04 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	This structure implement a progressbar.
*/

package gtk3Import

import (
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ProgressBarStruct struct {
	RefreshMs    uint
	GifImageName string

	gifImage        *gtk.Image
	box             *gtk.Box
	boxPosition     int
	fMain, fEnd     func() (err error)
	progressBar     *gtk.ProgressBar
	TimeOutContinue bool
	ticker          *time.Ticker
}

func ProgressBarNew(progressBar *gtk.ProgressBar, fMain, fEnd func() (err error)) (pbs *ProgressBarStruct) {
	pbs = new(ProgressBarStruct)
	pbs.Init(progressBar, fMain, fEnd)
	return
}

func ProgressGifNew(gifImage *gtk.Image, box *gtk.Box, position int, fMain, fEnd func() (err error)) (pbs *ProgressBarStruct) {
	pbs = new(ProgressBarStruct)
	pbs.Init(nil, fMain, fEnd)
	pbs.box = box
	pbs.boxPosition = position
	pbs.gifImage = gifImage
	pbs.box.Add(pbs.gifImage)
	pbs.gifImage.SetHAlign(gtk.ALIGN_FILL)
	pbs.gifImage.SetHExpand(true)
	pbs.box.ReorderChild(pbs.gifImage, position)
	// pbs.gifImage.Show()
	return
}

func (pbs *ProgressBarStruct) Init(progressBar *gtk.ProgressBar, fMain, fEnd func() (err error)) {
	pbs.RefreshMs = 100
	pbs.progressBar = progressBar
	pbs.fMain, pbs.fEnd = fMain, fEnd
}

func (pbs *ProgressBarStruct) StartGif() (err error) {
	pbs.gifImage.Show()
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err = pbs.fMain()
	}()

	waitGroup.Wait()
	if err != nil {
		return
	}
	glib.IdleAdd(func() {
		pbs.removeGifFromTopBox()
		err = pbs.fEnd()
	})
	return
}

func (pbs *ProgressBarStruct) removeGifFromTopBox() {
	pbs.box.Remove(pbs.gifImage)
}

func (pbs *ProgressBarStruct) StartTicker() (err error) {

	pbs.ticker = time.NewTicker(time.Millisecond * time.Duration(pbs.RefreshMs))
	go func() {
		for _ = range pbs.ticker.C {
			glib.IdleAdd(func() {
				pbs.progressBar.Pulse()
			})
		}
	}()

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err = pbs.fMain()
	}()

	waitGroup.Wait()
	if err != nil {
		return
	}
	glib.IdleAdd(func() {
		pbs.progressBar.SetFraction(0)
		err = pbs.fEnd()
		pbs.ticker.Stop()
	})

	return
}

func (pbs *ProgressBarStruct) StartTimeOut() {

	pbs.TimeOutContinue = true
	glib.TimeoutAdd(pbs.RefreshMs, func() bool {
		glib.IdleAdd(func() {
			pbs.progressBar.Pulse()
		})
		if pbs.TimeOutContinue {
			return true
		} else {
			glib.IdleAdd(func() {
				pbs.progressBar.SetFraction(0)
				pbs.fEnd()
			})
			return false
		}
	})

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		pbs.fMain()
	}()
	waitGroup.Wait()
	pbs.TimeOutContinue = false
	return
}
