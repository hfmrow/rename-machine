// spinButton.go

/*
	©2019 H.F.M. MIT license
*/

package gtk3Import

import (
	"fmt"
	"log"
	"math/big"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	gltses "github.com/hfmrow/genLib/tools/errors"
)

// SpinScaleSet: is a structure which contains parameters for the simple
// configuration of GtkSpinButton / GtkScaleButton. The options available
// are the most commonly used. If a callback function is present and
// 'AutoUpdateOnSignalCall' flag is set to 'true' (default value = false),
// the object's parameters will be updated on each 'value-updated' signal.
// Note: GtkSpinButton, GtkScaleButton are accepted as arguments, a
// GtkAdjustement object is created and linked to the given widget, it
// remain accessible using the main structure as all others parameters.
type SpinScaleSet struct {
	Widget interface{}
	Value,
	Min,
	Max,
	IncStep,
	PageIncrement,
	PageSize,
	SpinSetAccell float64
	SetDispDigits uint
	SpinSetNumeric,
	SpinSetSnapToTicks,
	SpinSetWrap,
	ScaleSetDrawValue,
	ScaleSetHasOrig,
	// On each "value-changed" signal call, values from structure are updated
	AutoUpdateOnSignalCall bool
	ScaleSetValuePos gtk.PositionType
	CallbackOnChange interface{}
	Adjustment       *gtk.Adjustment
}

// SpinScaleSetNew: create a new structure to which contains parameters
// for the simple configuration of GtkSpinButton / GtkScaleButton.
func SpinScaleSetNew(widget interface{}, min, max, value, step float64, callbackChangeOnChange ...interface{}) (*SpinScaleSet, error) {
	sss := new(SpinScaleSet)
	sss.Widget = widget

	sss.Value = value
	sss.Min = min
	sss.Max = max
	sss.IncStep = step
	sss.PageIncrement = 0
	sss.PageSize = 0
	// Defines the number of decimal to be displayed by the 'widget',
	// it is based by default on the value given in the increment step.
	// i.e: '1' or '1.0' will define 0 digits (default) '0.01' or '.01'
	// will set to 2 digits, '0.001' set it to 3 digits...
	// This method is not perfect but work in most common cases.
	f := strings.Split(fmt.Sprint(big.NewFloat(step).String()), ".")
	if len(f) > 1 {
		sss.SetDispDigits = uint(len(f[1]))
	}

	sss.SpinSetAccell = 0
	sss.SpinSetNumeric = false
	sss.SpinSetSnapToTicks = false
	sss.SpinSetWrap = false

	sss.ScaleSetDrawValue = true
	sss.ScaleSetHasOrig = true
	sss.ScaleSetValuePos = gtk.POS_BOTTOM

	if len(callbackChangeOnChange) > 0 {
		sss.CallbackOnChange = callbackChangeOnChange[0]
	} else {
		sss.CallbackOnChange = nil
	}
	err := sss.configure()
	return sss, err
}

// Configure: configuration of GtkSpinButton / GtkScaleButton with
// given/default parameters.
func (sss *SpinScaleSet) configure() (err error) {
	if sss.Adjustment, err = gtk.AdjustmentNew(sss.Value, sss.Min, sss.Max, sss.IncStep, sss.PageIncrement, sss.PageSize); err == nil {
		if sss.Widget != nil {
			callerInfos := gltses.WarnCallerMess(1) // error callback type handling
			switch obj := sss.Widget.(type) {
			case *gtk.SpinButton:
				obj.Configure(sss.Adjustment, sss.SpinSetAccell, sss.SetDispDigits)
				sss.updateConf()
				if sss.CallbackOnChange != nil {
					switch f := sss.CallbackOnChange.(type) {
					case func(spinBtn *gtk.SpinButton):
						obj.Connect("value-changed", func(spinBtn *gtk.SpinButton) {
							if sss.AutoUpdateOnSignalCall {
								sss.updateConf()
							}
							f(spinBtn)
						})
					default:
						// fmt.Printf(
						// 	"%s SpinScaleSetNew/configure: callback type 'func(spinBtn *gtk.SpinButton)' required!", callerInfos)
						return fmt.Errorf(
							"%s SpinScaleSetNew/configure: callback type 'func(spinBtn *gtk.SpinButton)' required!", callerInfos)
					}
				}
			case *gtk.Scale:
				obj.SetAdjustment(sss.Adjustment)
				sss.updateConf()
				if sss.CallbackOnChange != nil {
					switch f := sss.CallbackOnChange.(type) {
					case func(scaleBtn *gtk.Scale):
						obj.Connect("value-changed", func(scale *gtk.Scale) {
							if sss.AutoUpdateOnSignalCall {
								sss.updateConf()
							}
							f(scale)
						})
					default:
						// fmt.Printf(
						// 	"%s SpinScaleSetNew/configure: callback type 'func(scale *gtk.Scale)' required!", callerInfos)
						return fmt.Errorf(
							"%s SpinScaleSetNew/configure: callback type 'func(scale *gtk.Scale)' required!", callerInfos)
					}
				}
			}
		} else {
			err = fmt.Errorf("SpinScaleSetNew/Configure: There is not object to assign to!")
		}
	} else {
		err = fmt.Errorf("SpinScaleSetNew/Configure: %v", err)
	}
	return
}

// Set values ...
func (sss *SpinScaleSet) SetMin(v float64) {
	sss.Adjustment.SetProperty("lower", v)
	sss.Min = v
}
func (sss *SpinScaleSet) SetMax(v float64) {
	sss.Adjustment.SetProperty("upper", v)
	sss.Max = v
}
func (sss *SpinScaleSet) SetIncStep(v float64) {
	sss.Adjustment.SetProperty("step-increment", v)
	sss.IncStep = v
}
func (sss *SpinScaleSet) SetPageIncrement(v float64) {
	sss.Adjustment.SetProperty("page-increment", v)
	sss.PageIncrement = v
}
func (sss *SpinScaleSet) SetPageSize(v float64) {
	sss.Adjustment.SetProperty("page-size", v)
	sss.PageSize = v
}
func (sss *SpinScaleSet) SetSpinAccell(v float64) {
	switch obj := sss.Widget.(type) {
	case *gtk.SpinButton:
		obj.SetProperty("climb-rate", v)
	default:
		log.Println("SetSpinAccell: unable to configure. Work only with GtkSPinButton")
	}
	sss.SpinSetAccell = v
}
func (sss *SpinScaleSet) SetDigits(v uint) {
	switch obj := sss.Widget.(type) {
	case *gtk.SpinButton:
		obj.SetProperty("digits", v)
	case *gtk.Scale:
		obj.SetProperty("digits", v)
	default:
		log.Println("SetDigits: unable to configure. Work only with GtkSPinButton/GtkScale")
		return
	}
	sss.SetDispDigits = v
}
func (sss *SpinScaleSet) SetSpinNumeric(v bool) {
	switch obj := sss.Widget.(type) {
	case *gtk.SpinButton:
		obj.SetProperty("numeric", v)
	default:
		log.Println("SetSpinNumeric: unable to configure. Work only with GtkSPinButton")
	}
	sss.SpinSetNumeric = v
}
func (sss *SpinScaleSet) SetSpinSnapToTicks(v bool) {
	switch obj := sss.Widget.(type) {
	case *gtk.SpinButton:
		obj.SetProperty("snap-to-ticks", v)
		obj.SetProperty("wrap", v)
	default:
		log.Println("SetSpinSnapToTicks: unable to configure. Work only with GtkSPinButton")
	}
	sss.SpinSetSnapToTicks = v
}
func (sss *SpinScaleSet) SetSpinWrap(v bool) {
	switch obj := sss.Widget.(type) {
	case *gtk.SpinButton:
		obj.SetProperty("wrap", v)
	default:
		log.Println("SetSpinWrap: unable to configure. Work only with GtkSPinButton")
	}
	sss.SpinSetWrap = v
}
func (sss *SpinScaleSet) SetScaleDrawValue(v bool) {
	switch obj := sss.Widget.(type) {
	case *gtk.Scale:
		obj.SetProperty("draw-value", v)
	default:
		log.Println("SetScaleDrawValue: unable to configure. Work only with GtkScale")
	}
	sss.ScaleSetDrawValue = v
}
func (sss *SpinScaleSet) SetScaleHasOrig(v bool) {
	switch obj := sss.Widget.(type) {
	case *gtk.Scale:
		obj.SetProperty("has-origin", sss.ScaleSetHasOrig)
	default:
		log.Println("SetScaleHasOrig: unable to configure. Work only with GtkScale")
	}
	sss.ScaleSetHasOrig = v
}
func (sss *SpinScaleSet) SetScaleValuePos(v gtk.PositionType) {
	switch obj := sss.Widget.(type) {
	case *gtk.Scale:
		obj.SetProperty("value-pos", v)
	default:
		log.Println("SetScaleValuePos: unable to configure. Work only with GtkScale")
	}
	sss.ScaleSetValuePos = v
}

// updateConf: internal usage, update widget and GtkAdjustement with current
// values.
func (sss *SpinScaleSet) updateConf() {
	var setAdj = func() {
		sss.Adjustment.SetProperty("lower", sss.Min)
		sss.Adjustment.SetProperty("page-increment", sss.PageIncrement)
		sss.Adjustment.SetProperty("page-size", sss.PageSize)
		sss.Adjustment.SetProperty("step-increment", sss.IncStep)
		sss.Adjustment.SetProperty("upper", sss.Max)
		// sss.Adjustment.SetProperty("value", sss.Value)
	}
	switch obj := sss.Widget.(type) {
	case *gtk.SpinButton:
		setAdj()
		obj.SetProperty("climb-rate", sss.SpinSetAccell)
		obj.SetProperty("digits", sss.SetDispDigits)
		obj.SetProperty("numeric", sss.SpinSetNumeric)
		obj.SetProperty("snap-to-ticks", sss.SpinSetSnapToTicks)
		obj.SetProperty("wrap", sss.SpinSetWrap)
	case *gtk.Scale:
		setAdj()
		obj.SetProperty("digits", sss.SetDispDigits)
		obj.SetProperty("draw-value", sss.ScaleSetDrawValue)
		obj.SetProperty("has-origin", sss.ScaleSetHasOrig)
		obj.SetProperty("value-pos", sss.ScaleSetValuePos)
	}
}

// SpinbuttonSetValues: Configure a GtkSpinButton:
// where 'step' values means in order: stepIncrement, pageIncrement, pageSize
// 'nil' value send as 'sb' will just return a configured *gtk.Adjustment.
// Otherwise, the GtkSpinButton will be configured with the given values.
// Note: Accept GtkSpinButton, GtkScale as arguments.
func SpinbuttonSetValues(widget interface{}, min, max, value float64, step ...float64) (adjustment *gtk.Adjustment, err error) {
	var incStep, pageIncrement, pageSize float64 = 1, 0, 0

	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("./%s:%d Warning!, SpinbuttonSetValues implementation must be changed using new version\n", filepath.Base(file), line)

	switch len(step) {
	case 1:
		incStep = step[0]
	case 2:
		incStep = step[0]
		pageIncrement = step[1]
	case 3:
		incStep = step[0]
		pageIncrement = step[1]
		pageSize = step[2]
	}
	if adjustment, err = gtk.AdjustmentNew(value, min, max, incStep, pageIncrement, pageSize); err == nil {
		if widget != nil {
			switch obj := widget.(type) {
			case *gtk.SpinButton:
				obj.Configure(adjustment, 1, 0)
			case *gtk.Scale:
				obj.SetAdjustment(adjustment)
			}
		}
	} else {
		err = fmt.Errorf("SpinbuttonSetValues: %v", err)
	}
	return
}
