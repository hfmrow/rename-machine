// treeView.go

/*
	Source file auto-generated on Tue, 12 Nov 2019 22:14:11 using Gotk3ObjHandler v1.5 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019 H.F.M - TreeView library
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	This library allow you to facilitate Treeview operations.
	Can manage ListView and TreeView, only one of them at a time.
	This lib handle every kind of column type as defined in the gtk3 development manual.
	Some conveignant functions are available to manage entries, columns, values, rows ...

	i.e:
		func exampleTreeViewStructure() {
			var err error
			var tvs *gi.TreeViewStructure
			var storeSlice [][]interface{}
			var parentIter *gtk.TreeIter

			if tw, err := gtk.TreeViewNew(); err == nil { // Create TreeView. You can use existing one.
				if tvs, err = gi.TreeViewStructureNew(tw, false, false); err == nil { // Create Structure
					tvs.AddColumn("", "active", true, false, false, false, false) // With his columns
					tvs.AddColumn("Category", "markup", true, false, false, false, true)
					tvs.StoreSetup(new(gtk.TreeStore)) // Setup structure with desired TreeModel

					tvs.StoreDetach()        // Free TreeStore from TreeView while fill it. (useful for very large entries)
					for j := 0; j < 3; j++ { // Fill with parent nodes
						parentIter, _ = tvs.AddRow(nil, tvs.ColValuesIfaceToIfaceSlice(false, fmt.Sprintf("Parent %d", j)))

						for i := 0; i < 3; i++ { // Fill parents with childs nodes
							tvs.AddRow(parentIter, tvs.ColValuesIfaceToIfaceSlice(false, fmt.Sprintf("entry %d", i)))
						}
					}
					tvs.StoreAttach() // Say to TreeView that it get his StoreModel right now
				}
			}
			// Retrieve raw values with paths [][]interface{}. Can be done as [][]string too, and [][]interface{} without path.
			if err == nil {
				if storeSlice, err = tvs.StoreToIfaceSliceWithPaths(); err == nil {
					fmt.Println(storeSlice)
				}
			}
			if err != nil {
				log.Fatal(err)
			}
		}
*/

package gtk3Import

import (
	"errors"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// This library allow you to facilitate Treeview operations.
// Can manage ListView and TreeView, only one of them at a time.
// This lib handle every kind of column type as defined in the gtk3
// development manual. Some conveignant functions are available to
// manage entries, columns, values, rows ...
// Notice: All options, functions, if they're needed, must be set
// before starting "StoreSetup" function. Otherwise, you can modify
// all of them at runtime using Gtk3 objects (TreeView, ListStore,
// TreeStore, Columns, and so on). You can access it using the main
// structure.
type TreeViewStructure struct {

	// Actual TreeModel. Used in some functions to avoid use of
	// (switch ... case) type selection.
	Model *gtk.TreeModel

	// Direct acces to implicated objects
	TreeView  *gtk.TreeView
	ListStore *gtk.ListStore
	TreeStore *gtk.TreeStore
	Selection *gtk.TreeSelection

	// Basic options
	MultiSelection      bool
	ActivateSingleClick bool

	// The model has been modified?
	Modified bool

	// All columns option are available throught this structure
	Columns         []column
	ColumnsMinWidth int // Set a minimum width to avoid a warning gtk

	// TODO MAY CAUSE WHOLE FREEZE SOMETIME (undefined)
	// When "HasTooltip" is true, this function is launched,
	// Case of use: display tooltip according to rows currently hovered.
	// returned "bool" means display or not the tooltip.
	CallbackTooltipFunc func(iter *gtk.TreeIter, path *gtk.TreePath, column *gtk.TreeViewColumn, tooltip *gtk.Tooltip) bool
	HasTooltip,
	// Signify that we will use 'query-tooltip' signal callback
	UseQueryTooltip bool

	// Function to call when the selection has (possibly) changed.
	SelectionChangedFunc func()

	// This function is called (if not nil) each time a column value is changed.
	CallbackOnSetColValue func(iter *gtk.TreeIter, col int, value interface{})

	// Used for gtk.Model.ForEach functions
	ModelForEachFunc func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) bool

	// Used for gtk.TreeSelection.SetSelectFunction
	SelectFunction func(selection *gtk.TreeSelection, model *gtk.TreeModel, path *gtk.TreePath, selected bool) bool

	// Used for gtk.TreeSelection.ForEachFunc
	SelectionForEachFunc func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter)

	// Used to substract from Y coordinates when using tooltip
	headerHeight int
	// Used to determine wich TreeModel we work with.
	StoreType gtk.ITreeModel

	colTypeSl []glib.Type
}

type column struct {
	Name      string
	Editable  bool
	ReadOnly  bool
	Sortable  bool
	Resizable bool
	Expand    bool
	Visible   bool

	// Attributes with layout: "text", "markup", "pixbuf", "combo", "progress", "spinner", "active" (toggle button)
	// Attributes without layout: "pointer", "integer", "uint64", "int64"
	Attribute string

	// Direct access to the GtkTreeViewColumn and his CellRenderer
	Column       *gtk.TreeViewColumn
	ColType      glib.Type
	CellRenderer gtk.ICellRenderer

	// There is some default function defined for cell edition, normally, you don't have to define it yourself.
	// But in the case where you need specific operations, you can build you own edition function.
	EditTextFunc   func(cellRendererText *gtk.CellRendererText, path, text string) // "text"
	EditActiveFunc func(cellRendererToggle *gtk.CellRendererToggle, path string)   // "active" (toggle button)

	// Functions below require a type assertion that the CellRenderer type comes from,
	// e.g: cellRenderer.(*Gtk.CellRendererText).
	// Generic conditional edition function. The 'values' can be nonexistant, unique or multiple and sometimes,
	// require a type assertion in callback, that depend on which CellRenderer type the,'signal' come from.
	EditConditionFunc func(cellRenderer interface{}, path string, col int, values ...interface{}) bool

	// Callback function and flag that permit to write cell on cancel 'text'
	// both must be defined before 'store' initialisation.
	WriteOnCancel         bool
	WriteOnCancelCallback func(text string) bool // Usually used for confirmation dialog
}

// GetHeaderButton: Retrieve the button that assigned to the column header.
// May be used after "StoreSetup()" method.
func (col *column) GetHeaderButton() (button *gtk.Button, err error) {
	// var iWdgt *gtk.IWidget
	button = new(gtk.Button)
	if iWdgt, err := col.Column.GetButton(); err == nil {
		// Wrap button
		wdgt := iWdgt.ToWidget()
		actionable := &gtk.Actionable{wdgt.Object}
		button = &gtk.Button{gtk.Bin{gtk.Container{*wdgt}}, actionable}
	}
	return
}

/*********************\
*    Setup Functions   *
* Func that applied to *
* buils and handle     *
* treeview & models    *
***********************/

// Create a new treeview structure.
func TreeViewStructureNew(treeView *gtk.TreeView, multiselection, activateSingleClick bool) (tvs *TreeViewStructure, err error) {
	tvs = new(TreeViewStructure)
	// Store data
	tvs.ActivateSingleClick = activateSingleClick
	tvs.MultiSelection = multiselection
	tvs.TreeView = treeView
	tvs.ClearAll()
	tvs.ColumnsMinWidth = 13 // Set a minimum width to avoid a warning gtk

	tvs.HasTooltip = true

	tvs.CallbackTooltipFunc = nil /*func(iter *gtk.TreeIter, path *gtk.TreePath, column *gtk.TreeViewColumn, tooltip *gtk.Tooltip) bool {
		return false
	}*/

	if tvs.Selection, err = tvs.TreeView.GetSelection(); err != nil {
		return nil, fmt.Errorf("Unable to get gtk.TreeSelection: %s\n", err.Error())
	}
	return
}

// StoreSetup: Configure the TreeView columns and build the gtk.ListStore or
// gtk.TreeStore object. All parameters, personals or callback functions
// must be defined before calling StoreSetup in case you don't want to use
// the predefined ones and prefer to use your own (functions / callback).
// The "store" argument must be *gtk.ListStore or
// *gtk.TreeStore to indicate which kind of TreeModel we are working with ...
// i.e:
//   StoreSetup(new(gtk.TreeStore)), configure struct to work with a TreeStore.
//   StoreSetup(new(gtk.ListStore)), configure struct to work with a ListStore.
func (tvs *TreeViewStructure) StoreSetup(store gtk.ITreeModel) (err error) {
	var tmpColType glib.Type

	tvs.StoreType = store
	tvs.headerHeight = -1

	// Removing existing columns if there is ...
	for idx := int(tvs.TreeView.GetNColumns()) - 1; idx > -1; idx-- {
		tvs.TreeView.RemoveColumn(tvs.TreeView.GetColumn(idx))
	}

	// Tooltip setup
	if tvs.HasTooltip {

		tvs.TreeView.SetProperty("has-tooltip", true)

		if tvs.UseQueryTooltip {
			// TODO Disabled to show if issue (freeze) with gdpf stop occuring
			tvs.TreeView.Connect("query-tooltip", tvs.treeViewQueryTooltip)
		}

	} else {

		tvs.TreeView.SetProperty("has-tooltip", false)
	}

	// Set options
	tvs.TreeView.SetActivateOnSingleClick(tvs.ActivateSingleClick)
	if tvs.MultiSelection {
		tvs.Selection.SetMode(gtk.SELECTION_MULTIPLE)
	}

	// Build columns and his (default) edit function according to his type.
	for colIdx, _ := range tvs.Columns {
		if tmpColType, err = tvs.insertColumn(colIdx); err != nil {
			return fmt.Errorf("Unable to insert column nb %d: %s\n", colIdx, err.Error())
		} else {
			tvs.Columns[colIdx].ColType = tmpColType
			tvs.colTypeSl = append(tvs.colTypeSl, tmpColType)
		}
	}

	return tvs.buildStore()
}

// buildStore: Build ListStore or TreeStore object. Depending on provided
// object type in "StoreType" variable.
func (tvs *TreeViewStructure) buildStore() (err error) {

	switch tvs.StoreType.(type) {
	case *gtk.ListStore: // Create the ListStore.
		if tvs.ListStore, err = gtk.ListStoreNew(tvs.colTypeSl...); err != nil {
			return fmt.Errorf("Unable to create ListStore: %v", err)
		}
		tvs.Model = &tvs.ListStore.TreeModel
		tvs.TreeView.SetModel(tvs.ListStore)

	case *gtk.TreeStore: // Create the TreeStore.
		if tvs.TreeStore, err = gtk.TreeStoreNew(tvs.colTypeSl...); err != nil {
			return fmt.Errorf("Unable to create TreeStore: %v", err)
		}
		tvs.Model = &tvs.TreeStore.TreeModel
		tvs.TreeView.SetModel(tvs.TreeStore)

	}
	// Emitted whenever the selection has (possibly, RTFM) changed.
	if tvs.SelectionChangedFunc != nil { // link to callback function if exists.
		tvs.Selection.Connect("changed", tvs.SelectionChangedFunc)
	}
	return err
}

/*******************************\
*    Struct columns Functions    *
* Funct that applied to columns  *
* handled by the main structure  *
* this is the step before cols   *
* integration to the treeview    *
*********************************/

// RemoveColumn: Remove column from MainStructure and TreeView.
func (tvs *TreeViewStructure) RemoveColumn(col int) (columnCount int) {
	columnCount = tvs.TreeView.RemoveColumn(tvs.Columns[col].Column)
	tvs.Columns = append(tvs.Columns[:col], tvs.Columns[col+1:]...)
	tvs.Modified = true
	return
}

// InsertColumn: Insert new column to MainStructure, StoreSetup method must be called after.
func (tvs *TreeViewStructure) InsertColumn(name, attribute string, pos int, editable, readOnly,
	sortable, resizable, expand, visible bool) {
	newCol := []column{{Name: name, Attribute: attribute, Editable: editable, ReadOnly: readOnly,
		Sortable: sortable, Resizable: resizable, Expand: expand, Visible: visible}}

	tvs.Columns = append(tvs.Columns[:pos], append(newCol, tvs.Columns[pos:]...)...)
	tvs.Modified = true
}

// AddColumn: Adds a single new column to MainStructure.
// attribute may be: text, markup, pixbuf, combo, progress, spinner, active ...
// see above for complete list.
func (tvs *TreeViewStructure) AddColumn(name, attribute string, editable, readOnly,
	sortable, resizable, expand, visible bool) {

	// col := column{Name: name, Attribute: attrMap[attribute], Editable: editable, ReadOnly: readOnly,
	col := column{Name: name, Attribute: attribute, Editable: editable, ReadOnly: readOnly,
		Sortable: sortable, Resizable: resizable, Expand: expand, Visible: visible}
	tvs.Columns = append(tvs.Columns, col)
	tvs.Modified = true
}

// AddColumns: Adds several new columns to MainStructure.
// attribute may be: text, markup, pixbuf, combo, progress, spinner, active
func (tvs *TreeViewStructure) AddColumns(nameAndAttribute [][]string, editable, readOnly,
	sortable, resizable, expand, visible bool) {
	for _, inCol := range nameAndAttribute {
		tvs.AddColumn(inCol[0], inCol[1], editable, readOnly, sortable, resizable, expand, visible)
	}
	tvs.Modified = true
}

/**************************\
*    Building Functions     *
* Make to create treemodel  *
* and handle the differant  *
* kind of columns with      *
* predefined edit functions *
****************************/

// insertColumn: Insert column at defined position
func (tvs *TreeViewStructure) insertColumn(colIdx int) (colType glib.Type, err error) {
	// renderCell: Set cellRenderer type and column options
	var renderCell = func(cellRenderer gtk.ICellRenderer, colIdx int) (err error) {
		var column *gtk.TreeViewColumn

		if column, err = gtk.TreeViewColumnNewWithAttribute(tvs.Columns[colIdx].Name,
			cellRenderer,
			tvs.Columns[colIdx].Attribute,
			colIdx); err == nil {
			tvs.Columns[colIdx].CellRenderer = cellRenderer    // Store de CellRenderer used by this column
			tvs.Columns[colIdx].Column = column                // store column object to main struct.
			column.SetMinWidth(tvs.ColumnsMinWidth)            // Set a minimum width (13) to avoid a warning gtk
			column.SetExpand(tvs.Columns[colIdx].Expand)       // Set expand option
			column.SetResizable(tvs.Columns[colIdx].Resizable) // Set resizable option
			column.SetVisible(tvs.Columns[colIdx].Visible)     // Set visible option
			if tvs.Columns[colIdx].Sortable {                  // Set sortable option
				column.SetSortColumnID(colIdx)
			}
			tvs.TreeView.InsertColumn(column, colIdx)
		}
		return err
	}
	attribute := tvs.Columns[colIdx].Attribute
	switch {
	case attribute == "active": // "toggle"

		var cellRendererToggle *gtk.CellRendererToggle
		if cellRendererToggle, err = gtk.CellRendererToggleNew(); err == nil {

			// Define DEFAULT *conditional* function for checkbox edition if is not user defined
			if tvs.Columns[colIdx].EditConditionFunc == nil {
				tvs.Columns[colIdx].EditConditionFunc = func(cellRenderer interface{}, path string, col int, values ...interface{}) bool {
					return true
				}
			}

			// An edit function may be user defined before structure initialisation,
			// if not, this one, is set as the default edit function for checkboxes.
			if tvs.Columns[colIdx].EditActiveFunc == nil {
				tvs.Columns[colIdx].EditActiveFunc = func(cellRenderer *gtk.CellRendererToggle, path string) {
					if !tvs.Columns[colIdx].ReadOnly {

						// Conditional function: the "EditActiveFunc" is launched if condition is respected
						if tvs.Columns[colIdx].EditConditionFunc(cellRendererToggle, path, colIdx) {

							var (
								err        error
								currState  bool
								currIter   = new(gtk.TreeIter)
								parentIter = new(gtk.TreeIter)
							)

							if currIter, err = tvs.Model.GetIterFromString(path); err == nil {
								currState = tvs.GetColValue(currIter, colIdx).(bool)
								if err = tvs.SetColValue(currIter, colIdx, !currState); err == nil {
									// Change the state of the children if it exists
									if tvs.Model.IterHasChild(currIter) {
										tvs.changeChildrenTreeState(path, currIter, colIdx, !currState)
									}

									// Get parent if exits
									for tvs.Model.IterParent(parentIter, currIter) {
										// if tvs.Model.IterParent(parentIter, currIter) {
										// Change parent state if all child are checked, otherwise, uncheck it
										if err = tvs.SetColValue(parentIter, colIdx, tvs.AllChildsCheckedState(parentIter, colIdx, false)); err == nil {
											currIter = parentIter
											parentIter = new(gtk.TreeIter)
										}
									}
								}
							}
							if err != nil {
								log.Fatalf("Unable to edit (toggle) cell col %d, path %s: %s\n", colIdx, path, err.Error())
							}
						}
					}
				}
			}

			if err == nil {
				cellRendererToggle.Connect("toggled", tvs.Columns[colIdx].EditActiveFunc)
				if err = renderCell(cellRendererToggle, colIdx); err == nil {
					colType = glib.TYPE_BOOLEAN
				}

			}
		}
	case attribute == "spinner":
		var cellRendererSpinner *gtk.CellRendererSpinner
		if cellRendererSpinner, err = gtk.CellRendererSpinnerNew(); err == nil {
			cellRendererSpinner.SetProperty("editable", tvs.Columns[colIdx].Editable)
			if err = renderCell(cellRendererSpinner, colIdx); err == nil {
				colType = glib.TYPE_FLOAT
			}
		}
	case attribute == "progress":
		var cellRendererProgress *gtk.CellRendererProgress
		if cellRendererProgress, err = gtk.CellRendererProgressNew(); err == nil {
			cellRendererProgress.SetProperty("editable", tvs.Columns[colIdx].Editable)
			if err = renderCell(cellRendererProgress, colIdx); err == nil {
				colType = glib.TYPE_OBJECT
			}
		}
	case attribute == "pixbuf":
		var cellRendererPixbuf *gtk.CellRendererPixbuf
		if cellRendererPixbuf, err = gtk.CellRendererPixbufNew(); err == nil {
			if err = renderCell(cellRendererPixbuf, colIdx); err == nil {
				colType = glib.TYPE_OBJECT
			}
		}
	case attribute == "combo":
		var cellRendererCombo *gtk.CellRendererCombo
		if cellRendererCombo, err = gtk.CellRendererComboNew(); err == nil {
			cellRendererCombo.SetProperty("editable", tvs.Columns[colIdx].Editable)
			if err = renderCell(cellRendererCombo, colIdx); err == nil {
				colType = glib.TYPE_OBJECT
			}
		}
	case attribute == "text" || attribute == "markup":
		var cellRendererText *gtk.CellRendererText
		if cellRendererText, err = gtk.CellRendererTextNew(); err == nil {
			cellRendererText.SetProperty("editable", tvs.Columns[colIdx].Editable)
			/*
			 *
			 */

			if tvs.Columns[colIdx].WriteOnCancel {
				// Ask to record entry if the 'editing-canceled' signal appear (like focus out treeview)
				cellRendererText.Connect("editing-started",
					func(cellRendererText *gtk.CellRendererText, editable gtk.ICellEditable, path string) {
						var PreviousText string
						// Transtype ICellEditable interface to GtkCellEditable object
						cEditable, _ := editable.(*gtk.CellEditable)
						// Retrieve a GtkEntry from GtkCellEditable
						entry := cEditable.ToEntry()
						PreviousText, _ = entry.GetText()
						entry.Connect("editing-done", func() {
							text, _ := entry.GetText()
							entry.SetText(text)
							// If there is no change in the cell, simply exit
							if PreviousText == text {
								return
							}
							// On 'editing-canceled' signal, handling it to get Gtkentry content and write it to GtkListStore
							if value, _ := cEditable.GetProperty("editing-canceled"); value.(bool) {
								if tvs.Columns[colIdx].WriteOnCancelCallback != nil {
									if !tvs.Columns[colIdx].WriteOnCancelCallback(text) {
										return
									}
								}

								iter, _ := tvs.Model.GetIterFromString(path)
								tvs.SetColValue(iter, colIdx, text)
							}
						})
					})
			}

			/*
			 *
			 */
			// Define DEFAULT *conditional* function for text edition if is not user defined
			if tvs.Columns[colIdx].EditConditionFunc == nil {
				tvs.Columns[colIdx].EditConditionFunc = func(cellRenderer interface{}, path string, col int, values ...interface{}) bool {
					return true
				}
			}

			// An edit function may be user-defined before structure initialisation,
			// if not, this one, is set as the DEFAULT edit function for text cells.
			if tvs.Columns[colIdx].EditTextFunc == nil {

				tvs.Columns[colIdx].EditTextFunc = func(cellRenderer *gtk.CellRendererText, path string, text string) {
					if !tvs.Columns[colIdx].ReadOnly {

						// Conditional function: the "EditTextFunc" is launched if condition is respected
						if tvs.Columns[colIdx].EditConditionFunc(cellRendererText, path, colIdx, text) {

							var iter *gtk.TreeIter
							if iter, err = tvs.Model.GetIterFromString(path); err == nil {
								switch tvs.StoreType.(type) {
								case *gtk.ListStore:
									if err = tvs.ListStore.SetValue(iter, colIdx, text); err == nil {
										tvs.Modified = true
									}
								case *gtk.TreeStore:
									if err = tvs.TreeStore.SetValue(iter, colIdx, text); err == nil {
										tvs.Modified = true
									}
								}
							}
							if err != nil {
								log.Fatalf("Unable to edit (text) cell col %v, path %v, text %v: %v\n", colIdx, path, text, err)
							}
						}
					}
				}
			}

			if err == nil {
				cellRendererText.Connect("edited", tvs.Columns[colIdx].EditTextFunc)
				if err = renderCell(cellRendererText, colIdx); err == nil {
					colType = glib.TYPE_STRING
				}
			}
		}

	// For these Type, there is no layout
	case attribute == "pointer": // Pointer

		colType = glib.TYPE_POINTER

	case attribute == "integer": // INT

		colType = glib.TYPE_INT

	case attribute == "uint64": // UINT64

		colType = glib.TYPE_UINT64

	case attribute == "int64": // INT64

		colType = glib.TYPE_INT64

	default:
		err = fmt.Errorf("Error on setting attribute: %s is not implemented or inexistent.\n", tvs.Columns[colIdx].Attribute)
	}
	if err != nil {
		err = fmt.Errorf("Unable to make Renderer Cell: %s\n", err.Error())
	} else {
		// Add type to columns structure.
		tvs.Columns[colIdx].ColType = colType
	}
	return colType, err
}

/**********************\
*    Iters Functions    *
* Func that            *
* applied to Iters      *
************************/

// GetSelectedIters: retrieve list of selected iters,
// return nil whether nothing selected.
func (tvs *TreeViewStructure) GetSelectedIters() (iters []*gtk.TreeIter) {
	iters = make([]*gtk.TreeIter, tvs.Selection.CountSelectedRows())
	var count int
	// tvs.Selection.SelectedForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{}) {
	tvs.Selection.SelectedForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) {
		iters[count] = iter
		count++
	})
	if len(iters) == 0 {
		iters = nil
	}
	return
}

// GetSelectedPaths: retrieve list of selected paths,
// return nil whether nothing selected.
func (tvs *TreeViewStructure) GetSelectedPaths() (paths []*gtk.TreePath) {
	paths = make([]*gtk.TreePath, tvs.Selection.CountSelectedRows())
	var count int
	tvs.Selection.SelectedForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) {
		paths[count] = path
		count++
	})
	if len(paths) == 0 {
		paths = nil
	}
	return
}

// ItersSelect: Select provided Iters.
func (tvs *TreeViewStructure) ItersSelect(iters ...*gtk.TreeIter) {
	for _, iter := range iters {
		if !tvs.Selection.IterIsSelected(iter) {
			tvs.Selection.SelectIter(iter)
		}
	}
}

// ItersUnselectAll: Unselect all selected iters.
func (tvs *TreeViewStructure) ItersUnselectAll() {
	tvs.Selection.UnselectAll()
}

// ItersUnselect: Unselect provided Iters.
func (tvs *TreeViewStructure) ItersUnselect(iters ...*gtk.TreeIter) {
	for _, iter := range iters {
		if tvs.Selection.IterIsSelected(iter) {
			tvs.Selection.UnselectIter(iter)
		}
	}
}

// ItersSelectRange: Select range between start and end iters.
func (tvs *TreeViewStructure) ItersSelectRange(startIter, endIter *gtk.TreeIter) (err error) {
	var startPath, endPath *gtk.TreePath

	if startPath, err = tvs.Model.GetPath(startIter); err == nil {
		if endPath, err = tvs.Model.GetPath(endIter); err == nil {
			tvs.Selection.SelectRange(startPath, endPath)
		}
	}
	return err
}

// ScrollToIter: scroll to iter, pointing to the column if it has been specified.
func (tvs *TreeViewStructure) IterScrollTo(iter *gtk.TreeIter, column ...int) (err error) {
	var path *gtk.TreePath
	var colNb int
	if len(column) > 0 {
		colNb = column[0]
	}
	if col := tvs.TreeView.GetColumn(colNb); col != nil {
		if path, err = tvs.Model.GetPath(iter); err == nil {
			if path != nil {
				tvs.TreeView.ScrollToCell(path, col, true, 0.5, 0.5)
			} else {
				err = fmt.Errorf("IterScrollTo: Unable to get path from iter\n")
			}
		}
	} else {
		err = fmt.Errorf("IterScrollTo: Unable to get column %d\n", colNb)
	}
	return
}

/*******************************\
*    Cols Functions              *
* Funct that applied to Columns *
*********************************/

// GetColValue: Get value from iter of specific column as interface type.
func (tvs *TreeViewStructure) GetColValue(iter *gtk.TreeIter, col int) (value interface{}) {
	var gValue *glib.Value
	var err error
	if gValue, err = tvs.Model.GetValue(iter, col); err == nil {
		if value, err = gValue.GoValue(); err == nil {
			return
		}
	}
	log.Fatalf("GetColValue: %s\n", err.Error())
	return
}

// SetColValueWithCallback: set the value to iter for a specific column as an interface type.
// A callback function is called giving the possibility to perform interactive operations.
// func (tvs *TreeViewStructure) SetColValueWithCallback (iter *gtk.TreeIter, col int, value interface{}) (err error) {

// 	var (
// 		oldVal   *glib.Value
// 		oldValue interface{}
// 	)

// 	switch tvs.StoreType.(type) {
// 	case *gtk.ListStore:
// 		if oldVal, err = tvs.ListStore.GetValue(iter, col); err == nil {
// 			if oldValue, err = oldVal.GoValue(); err == nil {
// 				err = tvs.ListStore.SetValue(iter, col, value)
// 			}
// 		}
// 	case *gtk.TreeStore:
// 		if oldVal, err = tvs.TreeStore.GetValue(iter, col); err == nil {
// 			if oldValue, err = oldVal.GoValue(); err == nil {
// 				err = tvs.TreeStore.SetValue(iter, col, value)
// 			}
// 		}
// 	}
// 	if err == nil {
// 		tvs.Modified = true
// 		if tvs.CallbackOnSetColValue != nil {
// 			tvs.CallbackOnSetColValue(iter, col, value, oldValue)
// 			return
// 		}
// 	}
// 	return
// }

// SetColValue: Set value to iter for a specific column as interface type.
func (tvs *TreeViewStructure) SetColValue(iter *gtk.TreeIter, col int, value interface{}) (err error) {

	switch tvs.StoreType.(type) {
	case *gtk.ListStore:

		err = tvs.ListStore.SetValue(iter, col, value)
	case *gtk.TreeStore:

		err = tvs.TreeStore.SetValue(iter, col, value)
	}
	if err == nil {
		tvs.Modified = true
		if tvs.CallbackOnSetColValue != nil {
			tvs.CallbackOnSetColValue(iter, col, value)
		}
	}
	return
}

// GetColValueFromPath: Get value from path of specific column as interface type.
// Note: should be used only if there is no other choice, prefer using iter to get values.
func (tvs *TreeViewStructure) GetColValuePath(path *gtk.TreePath, col int) (value interface{}) {
	var err error
	var iter *gtk.TreeIter

	switch tvs.StoreType.(type) {
	case *gtk.ListStore:
		iter, err = tvs.ListStore.GetIter(path)
	case *gtk.TreeStore:
		iter, err = tvs.TreeStore.GetIter(path)
	}
	if err != nil {
		log.Fatalf("GetColValuePath: unable to get iter from path: %s\n", err.Error())
		return
	}
	return tvs.GetColValue(iter, col)
}

// SetColValue: Set value to path for a specific column as interface type.
// Note: should be used only if there is no other choice, prefer using iter to set values.
func (tvs *TreeViewStructure) SetColValuePath(path *gtk.TreePath, col int, goValue interface{}) (err error) {

	var iter *gtk.TreeIter

	switch tvs.StoreType.(type) {

	case *gtk.ListStore:
		if iter, err = tvs.ListStore.GetIter(path); err == nil {
			err = tvs.ListStore.SetValue(iter, col, goValue)
		}
	case *gtk.TreeStore:
		if iter, err = tvs.TreeStore.GetIter(path); err == nil {
			err = tvs.TreeStore.SetValue(iter, col, goValue)
		}
	}
	if err != nil {
		return
	}
	tvs.Modified = true
	return
}

/***************************\
*    Rows Functions          *
* Func that applied to rows *
*****************************/

// CountRows: Return the number of rows in treeview.
func (tvs *TreeViewStructure) CountRows() (count int) {
	switch tvs.StoreType.(type) {
	case *gtk.ListStore:
		count = tvs.Model.IterNChildren(nil)
	case *gtk.TreeStore:
		// two different way because TreeStore return the number of toplevel nodes intead of ListStore

		// TODO May cause 'fatal error: invalid pointer found on stack' when used intensivly (querytooltip)
		// Find another way to doing the same thing
		tvs.Model.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) bool {
			count++
			return false
		})
	}
	return
}

// GetRowNbIter: Return the row number handled by the given iter,
// without any depth consideration.
func (tvs *TreeViewStructure) GetRowNbIter(iter *gtk.TreeIter) int {
	path, err := tvs.Model.GetPath(iter)
	if err != nil {
		fmt.Printf("Unable to get row number: %s\n", err.Error())
		return -1
	}
	ind := path.GetIndices()
	return ind[len(ind)-1:][0]
}

// AddRow: Append a row to the Store (defined by type of "StoreType" variable).
// "parent" is useless for ListStore, if its set to nil on TreeStore,
// it will create a new parent
func (tvs *TreeViewStructure) AddRow(parent *gtk.TreeIter, row ...interface{}) (iter *gtk.TreeIter, err error) {
	return tvs.InsertRow(parent, -1, row...)
}

// InsertRow: Insert a row to the Store (defined by type of "StoreType" variable).
// "parent" is useless for ListStore, if its set to nil on TreeStore,
// it will create a new parent otherwise the new iter will be a child of it.
// "insertPos" indicate row number for insertion, set to -1 means append at the end.
func (tvs *TreeViewStructure) InsertRow(parent *gtk.TreeIter, insertPos int, row ...interface{}) (iter *gtk.TreeIter, err error) {

	iter = new(gtk.TreeIter)

	var colIdx = make([]int, len(row))
	for idx := 0; idx < len(row); idx++ {
		colIdx[idx] = idx
	}

	switch tvs.StoreType.(type) {
	case *gtk.ListStore:
		err = tvs.ListStore.InsertWithValues(iter, insertPos, colIdx, row)
	case *gtk.TreeStore:
		err = tvs.TreeStore.InsertWithValues(iter, parent, insertPos, colIdx, row)
	}

	if err != nil {
		return nil, fmt.Errorf("Unable to add row %d: %s\n", insertPos, err.Error())
	}
	tvs.Modified = true
	return
}

// InsertRowAtIter: Insert a row after/before iter to "StoreType": ListStore/Treestore.
// Parent may be nil for Liststore.
func (tvs *TreeViewStructure) InsertRowAtIterN(inIter, parent *gtk.TreeIter, row []interface{}, before ...bool) (iter *gtk.TreeIter, err error) {
	var tmpBefore bool
	var colIdx []int
	// var path *gtk.TreePath
	if len(before) != 0 {
		tmpBefore = before[0]
	}
	for idx, _ := range row {
		colIdx = append(colIdx, idx)
	}
	switch tvs.StoreType.(type) {
	case *gtk.ListStore:
		if tmpBefore { // Get the insertion iter
			iter = tvs.ListStore.InsertBefore(inIter)
		} else {
			iter = tvs.ListStore.InsertAfter(inIter)
		}
		err = tvs.ListStore.Set(iter, colIdx, row)
	case *gtk.TreeStore:
		if tmpBefore { // Get the insertion iter
			iter = tvs.TreeStore.InsertBefore(parent, inIter)
		} else {
			iter = tvs.TreeStore.InsertAfter(parent, inIter)
		}
		err = tvs.TreeStore.SetValue(iter, colIdx[0], row[0])
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to insert row: %s\n", err.Error()))
	}
	tvs.Modified = true
	return iter, err
}

// InsertRowAtIter: Insert a row after/before iter to "StoreType": ListStore/Treestore.
// Parent may be nil for Liststore.
func (tvs *TreeViewStructure) InsertRowAtIter(inIter, parent *gtk.TreeIter, row ...interface{}) (iter *gtk.TreeIter, err error) {
	if iter, err = tvs.InsertRow(parent, tvs.GetRowNbIter(inIter)+1, row...); err == nil {
		tvs.Modified = true
	}
	return iter, err
}

// TODO rewrite to be much faster
// DuplicateRow: Copy a row after iter to the listStore
func (tvs *TreeViewStructure) DuplicateRow(inIter, parent *gtk.TreeIter) (iter *gtk.TreeIter, err error) {
	var glibValue *glib.Value
	var goValue interface{}
	switch tvs.StoreType.(type) {
	case *gtk.ListStore:
		iter = tvs.ListStore.InsertAfter(inIter)
		for colIdx, _ := range tvs.Columns {
			if glibValue, err = tvs.ListStore.GetValue(inIter, colIdx); err == nil {
				if goValue, err = glibValue.GoValue(); err == nil {
					err = tvs.ListStore.SetValue(iter, colIdx, goValue)
				}
			}
		}
	case *gtk.TreeStore:
		iter = tvs.TreeStore.InsertAfter(parent, inIter)
		for colIdx, _ := range tvs.Columns {
			if glibValue, err = tvs.TreeStore.GetValue(inIter, colIdx); err == nil {
				if goValue, err = glibValue.GoValue(); err == nil {
					err = tvs.TreeStore.SetValue(iter, colIdx, goValue)
				}
			}
		}
	}
	if err != nil {
		return nil, fmt.Errorf("Unable to duplicating row: %s\n", err.Error())
	}
	tvs.Modified = true
	tvs.ItersUnselect(inIter)
	tvs.ItersSelect(iter)
	return iter, err
}

// RemoveSelectedRows: Delete entries from selected iters or from given iters.
func (tvs *TreeViewStructure) RemoveSelectedRows(iters ...*gtk.TreeIter) (count int) {

	if len(iters) == 0 {
		iters = tvs.GetSelectedIters()
	}

	tmpIters := make([]*gtk.TreeIter, len(iters))
	// Reverse iters position to avoid error on remove.
	// (when iter is removed, the next come to be invalid)
	countIters := len(iters) - 1
	for idx := countIters; idx >= 0; idx-- {
		tmpIters[countIters-idx] = iters[idx]
	}

	return tvs.RemoveRows(tmpIters...)
}

// RemoveRows: Unified remove iter(s) function
func (tvs *TreeViewStructure) RemoveRows(iters ...*gtk.TreeIter) (count int) {

	switch tvs.StoreType.(type) {
	case *gtk.ListStore:
		for _, iter := range iters {
			if tvs.ListStore.Remove(iter) {
				count++
			}
		}
	case *gtk.TreeStore:
		for _, iter := range iters {
			if tvs.TreeStore.Remove(iter) {
				count++
			}
		}
	}
	if count > 0 {
		tvs.Modified = true
	}
	return
}

// RemoveRowsPath: Unified remove path(s) function. accept '*gtk.TreePath' & 'TreePath' as string
func (tvs *TreeViewStructure) RemoveRowsPath(paths ...interface{}) (count int, err error) {

	var iter *gtk.TreeIter

	for _, p := range paths {

		switch path := p.(type) {

		case *gtk.TreePath:
			iter, err = tvs.Model.GetIter(path)
		case string:
			iter, err = tvs.Model.GetIterFromString(path)
		}

		if err != nil {
			break
		}

		switch tvs.StoreType.(type) {

		case *gtk.ListStore:
			if tvs.ListStore.Remove(iter) {
				count++
			}
		case *gtk.TreeStore:
			if tvs.TreeStore.Remove(iter) {
				count++
			}
		}
	}

	if count > 0 {
		tvs.Modified = true
	}
	return
}

// GetSelectedRows: Get entries from selected iters as [][]string.
func (tvs *TreeViewStructure) GetSelectedRows() (outSlice [][]string, err error) {
	outSlice = make([][]string, tvs.Selection.CountSelectedRows())
	var count int
	tvs.Selection.SelectedForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) {
		if outSlice[count], err = tvs.GetRow(iter); err != nil {
			err = fmt.Errorf("Unable to get selected row: %d\n", count)
		}
		count++
	})
	return
}

// getRow: Get row from iter as []string
func (tvs *TreeViewStructure) GetRow(iter *gtk.TreeIter) (outSlice []string, err error) {
	var glibValue *glib.Value
	var valueString string

	for colIdx := 0; colIdx < len(tvs.Columns); colIdx++ {
		if glibValue, err = tvs.Model.GetValue(iter, colIdx); err == nil {
			if valueString, err = tvs.getStringCellValueByType(glibValue); err == nil {
				outSlice = append(outSlice, valueString)
			}
		}
		if err != nil {
			break
		}
	}
	return outSlice, err
}

// GetRowIface: Get row from iter as []interface{}
func (tvs *TreeViewStructure) GetRowIface(iter *gtk.TreeIter) (outIface []interface{}, err error) {
	var glibValue *glib.Value
	var value interface{}

	for colIdx := 0; colIdx < len(tvs.Columns); colIdx++ {
		if glibValue, err = tvs.Model.GetValue(iter, colIdx); err == nil {
			if value, err = glibValue.GoValue(); err == nil {
				outIface = append(outIface, value)
			}
		}
		if err != nil {
			break
		}
	}
	return outIface, err
}

/*******************************\
*    Convenient Functions       *
* Designed to make life easier *
********************************/

// GetColumns: Retieve columns available in the current TreeView.
func (tvs *TreeViewStructure) GetColumns() (out []*gtk.TreeViewColumn) {
	glibList := tvs.TreeView.GetColumns()
	glibList.Foreach(func(item interface{}) {
		out = append(out, item.(*gtk.TreeViewColumn))
	})
	return
}

// StoreDetach: Unlink "TreeModel" from TreeView. Useful when lot of rows need
// be inserted. After insertion, StoreAttach() must be used to restore the link
// with the treeview. tips: must be used before ListStore/TreeStore.Clear().
func (tvs *TreeViewStructure) StoreDetach() {
	if tvs.StoreType != nil {
		tvs.Model.Ref()
		tvs.TreeView.SetModel(nil)
	}
}

// StoreAttach: To use after data insertion to restore the link with TreeView.
func (tvs *TreeViewStructure) StoreAttach() {
	if tvs.StoreType != nil {
		tvs.TreeView.SetModel(tvs.Model)
		tvs.Model.Unref()
	}
}

// Clear: Clear the current used Model:
// unified version of gtk.TreeStore.Clear() or gtk.ListStore.Clear()
func (tvs *TreeViewStructure) Clear() {
	switch tvs.StoreType.(type) {
	case *gtk.ListStore:
		tvs.ListStore.Clear()
	case *gtk.TreeStore:
		tvs.TreeStore.Clear()
	}
}

// ClearAll: Clear TreeView's columns, ListStore / TreeStore object.
// Depending on provided object type into the "StoreType" variable.
// To reuse structure, you must execute StoreSetup() again after
// added new columns.
func (tvs *TreeViewStructure) ClearAll() (err error) {
	if tvs.TreeView != nil {
		// Removing existing columns if exists ...
		for idx := int(tvs.TreeView.GetNColumns()) - 1; idx > -1; idx-- {
			tvs.TreeView.RemoveColumn(tvs.TreeView.GetColumn(idx))
		}
		tvs.Columns = tvs.Columns[:0]
		tvs.TreeView.SetModel(nil)
		switch tvs.StoreType.(type) {
		case *gtk.ListStore:
			if tvs.ListStore != nil {
				tvs.ListStore.Clear()
				tvs.ListStore.Unref()
			}
		case *gtk.TreeStore:
			if tvs.TreeStore != nil {
				tvs.TreeStore.Clear()
				tvs.TreeStore.Unref()
			}
		}
		tvs.Modified = false
	}
	return
}

// StoreToStringSlice: Retrieve all the rows values from a 'StoreType'
// as [][]string
func (tvs *TreeViewStructure) StoreToStringSlice() (out [][]string, err error) {
	var row []string
	tvs.Model.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) bool {
		if row, err = tvs.GetRow(iter); err == nil {
			out = append(out, row)
		} else {
			return true
		}
		return false
	})
	return
}

// GetTreeViewIface: retrieve whole content of a TreeView as
// [][]interface
func (tvs *TreeViewStructure) StoreToIfaceSlice() (out [][]interface{}, err error) {
	var row []interface{}
	tvs.Model.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) bool {
		if row, err = tvs.GetRowIface(iter); err == nil {
			out = append(out, row)
		} else {
			return true
		}
		return false
	})
	return
}

// ColNamesToStringSlice:  Retrieve column names as string slice
func (tvs *TreeViewStructure) ColNamesToStringSlice() (outSlice []string) {
	for _, col := range tvs.Columns {
		outSlice = append(outSlice, col.Column.GetTitle())
	}
	return
}

// GetTreeCol: This method retrieve data from a [single column] of the current
// 'GtkTreeStore' as []string. Use GetTreeFullIface to retrieve multiple columns
// at once.
func (tvs *TreeViewStructure) GetTreeCol(toggleCol, dataCols int) (checked, unChecked []string, err error) {
	var row []interface{}

	tvs.TreeStore.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) bool {

		if row, err = tvs.GetRowIface(iter); err == nil {

			if row[toggleCol].(bool) {
				checked = append(checked, row[dataCols].(string))
			} else {
				unChecked = append(unChecked, row[dataCols].(string))
			}
		} else {

			err = fmt.Errorf("GetTreeCol: %v", err)
			return true
		}

		return false
	})
	return
}

// GetTreeFullIface: Retrieve the whole content of the current 'GtkTreeStore'.
func (tvs *TreeViewStructure) GetTreeFullIface(toggleCol int) (checked, unChecked [][]interface{}, err error) {
	var row []interface{}

	tvs.TreeStore.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) bool {

		if row, err = tvs.GetRowIface(iter); err == nil {

			if row[toggleCol].(bool) {
				checked = append(checked, row)
			} else {
				unChecked = append(unChecked, row)
			}
		} else {

			err = fmt.Errorf("GetTreeFullIface: %v", err)
			return true
		}

		return false
	})
	return
}

// GetTreeParentIface: return the 'row' holded by the parent of 'iter' or from itself
// whether having children. After call with success, provided 'iter' argument points
// to the one the line came from.
func (tvs *TreeViewStructure) GetTreeParentIface(iter *gtk.TreeIter) ([]interface{}, error) {

	if !tvs.Model.IterHasChild(iter) {

		parent := new(gtk.TreeIter)
		if !tvs.Model.IterParent(parent, iter) {
			return nil, fmt.Errorf("Unable to get parent iter from child")
		}
		iter = parent
	}
	return tvs.GetRowIface(iter)
}

// AddTree: This function, add a full tree to a TreeStore, childs will be added to
// parent's tree if exists or to a new created parent.
// Each entry handle a checkbox and a name only.
//
// e.g: in 3 calls to 'AddTree'
// 'iFace' = []interface{"github.com", "hfmrow", "gtk3Import", "pango.go"}
// 'iFace' = []interface{"github.com", "hfmrow", "gtk3Import", "main.go"}
// 'iFace' = []interface{"github.com", "hfmrow", "gtk3Import", "view.go"}
// add all nodes for this tree to column 'filepathCol' and checkbox at 'toggleCol'
//
// []  github.com
// []  ╰── hfmrow
// []      ╰── gtk3Import
// []          ├── pango.go
// []          ├── main.go
// []          ╰── view.go
//
// - 'stateDefault' default value for 'toggleCol'.
// - 'callbackParentCreation' can be 'nil', otherwise, this callback corresponds
//   to the creation of the iter during the construction of the tree, allow to control
//   (allows to manually add some usefull data to the row).
// - The returned 'outIter', target the iter used by the the last entry, useful for
//   manually adding more entries to columns if needed.
func (tvs *TreeViewStructure) AddTree(
	toggleCol,
	filepathCol int,
	stateDefault bool,
	callbackParentCreation func(store *gtk.TreeStore, iter *gtk.TreeIter, currentAddIdx int, iRow *[]interface{}) bool,
	iFace ...interface{}) (outIter *gtk.TreeIter, err error) {

	var (
		iterate,
		childsCount int
		tmpIter      *gtk.TreeIter
		ok           bool
		value        string
		pathSplitted []interface{}
	)

	// addItem: It only does what the name says.
	var addItem = func(toAdd string, iter *gtk.TreeIter) (tmpIter *gtk.TreeIter) {

		var err error
		tmpIter = tvs.TreeStore.Append(iter)

		if err = tvs.TreeStore.SetValue(tmpIter, toggleCol, stateDefault); err == nil {
			if err = tvs.TreeStore.SetValue(tmpIter, filepathCol, toAdd); callbackParentCreation != nil && err == nil {

				if !callbackParentCreation(tvs.TreeStore, tmpIter, iterate, &pathSplitted) {
					// remove created iter and quit if callback is not true
					tvs.TreeStore.Remove(tmpIter)
					err = fmt.Errorf("callbackParentCreation: has rejected the 'iter' creation")
				}
				// if iterate > len(pathSplitted)-1 {
				// 	tmpIter = nil
				// }
			}
		}
		if err != nil {
			log.Printf("AddTree: Unable to addItem. %s\n", err.Error())
			tmpIter = nil
		}
		return
	}

	// findCreatFirstParent: Add or find first parent that match "name" and return it's iter.
	var findCreatFirstParent = func(name string) (tmpIter *gtk.TreeIter, ok bool, err error) {

		if len(name) > 0 {
			tmpIter, ok = tvs.TreeStore.GetIterFirst()

			for ok {
				value := tvs.GetColValue(tmpIter, filepathCol).(string)
				if value == name {
					return
				}
				ok = tvs.TreeStore.IterNext(tmpIter)
			}

			// Nothing found, then create it
			tmpIter = addItem(pathSplitted[iterate].(string), outIter)
			ok = tmpIter != nil

		} else {
			err = errors.New("findCreatFirstParent: Could not proceed with empty parent.")
		}
		return
	}

	// searchMatch: Walk trought iter to retrieve the one matching "toMatch".
	var searchMatch = func(toMatch string, outIter *gtk.TreeIter) (childIter *gtk.TreeIter, ok bool) {

		childIter = new(gtk.TreeIter)
		ok = tvs.TreeStore.IterChildren(outIter, childIter)
		for ok {
			value = tvs.GetColValue(childIter, filepathCol).(string)
			if value == toMatch {
				return
			}
			ok = tvs.Model.IterNext(childIter)
		}
		return
	}

	pathSplitted = iFace

	// parse "dirPath" entry into treestore.
	if outIter, ok, err = findCreatFirstParent(pathSplitted[0].(string)); err == nil {

		for ok {
			value = tvs.GetColValue(outIter, filepathCol).(string)
			if value == pathSplitted[iterate] {
				childsCount = tvs.Model.IterNChildren(outIter)
				if childsCount > 0 {
					iterate++
					if iterate >= len(pathSplitted) {
						break
					}
					if tmpIter, ok = searchMatch(pathSplitted[iterate].(string), outIter); !ok {
						outIter = addItem(pathSplitted[iterate].(string), outIter)
						ok = outIter != nil

						// last := pathSplitted[len(pathSplitted)-1].(string)
						// fmt.Printf("item: %v, last: %v\n", pathSplitted[iterate], pathSplitted[len(pathSplitted)-1])

					} else {
						outIter = tmpIter
					}
					continue
				} else {
					iterate++
					if iterate >= len(pathSplitted) {
						break
					}
					outIter = addItem(pathSplitted[iterate].(string), outIter)
					ok = outIter != nil

					// fmt.Println(pathSplitted[iterate].(string))

				}
			} else {
				if tmpIter, ok = searchMatch(pathSplitted[iterate].(string), outIter); ok {
					outIter = addItem(pathSplitted[iterate].(string), tmpIter)
					ok = outIter != nil
					// fmt.Println(pathSplitted[iterate].(string))
				}
			}
		}
	}
	return
}

// ColValuesStringSliceToIfaceSlice: Convert string list to []interface, for simplify adding text rows
func (tvs *TreeViewStructure) ColValuesStringSliceToIfaceSlice(inSlice ...string) (outIface []interface{}) {
	outIface = make([]interface{}, len(inSlice))
	for idx, data := range inSlice {
		outIface[idx] = data
	}
	return
}

// glibType:  glib value type List structure.
var glibType = map[glib.Type]string{
	0:  "glib.TYPE_INVALID",
	4:  "glib.TYPE_NONE",
	8:  "glib.TYPE_INTERFACE",
	12: "glib.TYPE_CHAR",
	16: "glib.TYPE_UCHAR",
	20: "glib.TYPE_BOOLEAN",
	24: "glib.TYPE_INT",
	28: "glib.TYPE_UINT",
	32: "glib.TYPE_LONG",
	36: "glib.TYPE_ULONG",
	40: "glib.TYPE_INT64",
	44: "glib.TYPE_UINT64",
	48: "glib.TYPE_ENUM",
	52: "glib.TYPE_FLAGS",
	56: "glib.TYPE_FLOAT",
	60: "glib.TYPE_DOUBLE",
	64: "glib.TYPE_STRING",
	68: "glib.TYPE_POINTER",
	72: "glib.TYPE_BOXED",
	76: "glib.TYPE_PARAM",
	80: "glib.TYPE_OBJECT",
	84: "glib.TYPE_VARIANT",
}

/************************\
*    Helpers Functions    *
* Made to simplify some  *
* more complex functions  *
**************************/

// TODO MAY CAUSE WHOLE FREEZE SOMETIME (undefined)
// treeViewQueryTooltip: function to display tooltip according to rows currently hovered
// Note: return TRUE if tooltip should be shown right now, FALSE otherwise
func (tvs *TreeViewStructure) treeViewQueryTooltip(tw *gtk.TreeView, x, y int, KeyboardMode bool, tooltip *gtk.Tooltip) bool {

	if tvs.CallbackTooltipFunc != nil && tvs.HasTooltip {

		if tvs.CountRows() > 0 {
			// we need to substract header height to "y" position to get the correct path.
			if path, column, _, _, isBlank := tvs.TreeView.IsBlankAtPos(x, y-tvs.getHeaderHeight()); !isBlank {
				if iter, err := tvs.Model.GetIter(path); err == nil {

					return tvs.CallbackTooltipFunc(iter, path, column, tooltip)
				} else {

					log.Printf("treeViewQueryTooltip:GetIter: %s\n", err.Error())
				}
			}
		}
	}
	return false
}

// getHeaderHeight: Used to get height of header to use with [TreeView.IsBlankAtPos],
// It is needed to decrease y pos by height of cells to get a correct path value.
func (tvs *TreeViewStructure) getHeaderHeight() (height int) {
	if tvs.headerHeight < 0 { // That means that is launched only at the first call.
		for gtk.EventsPending() {
			gtk.MainIteration() // Wait for pending events (until the widget is redrawn)
		}
		// Getting header height
		backupVisibleHeader := tvs.TreeView.GetHeadersVisible()
		tvs.TreeView.SetHeadersVisible(true)
		withHeader, _ := tvs.TreeView.GetPreferredHeight()
		tvs.TreeView.SetHeadersVisible(false)
		withoutHeader, _ := tvs.TreeView.GetPreferredHeight()
		tvs.TreeView.SetHeadersVisible(backupVisibleHeader)
		tvs.headerHeight = withHeader - withoutHeader
	}
	return tvs.headerHeight
}

// getCellValueByType: Retrieve cell value and convert it to string based on
// his type. Used by GetRow func.
func (tvs *TreeViewStructure) getStringCellValueByType(glibValue *glib.Value) (valueString string, err error) {
	var actualType glib.Type
	var valueIface interface{}

	if actualType, _, err = glibValue.Type(); err == nil {
		switch actualType {

		// String
		case glib.TYPE_STRING:
			valueString, _ = glibValue.GetString()

		// Numeric values: int, uint64, int64
		case glib.TYPE_INT64, glib.TYPE_UINT64, glib.TYPE_INT:
			if valueIface, err = glibValue.GoValue(); err == nil {
				switch val := valueIface.(type) {
				case int, uint64, int64:
					valueString = fmt.Sprintf("%d", val)
				}
			}

		// Pointer, just say that it's what's it ...
		case glib.TYPE_POINTER:
			valueString = fmt.Sprintf("%v", glibValue.GetPointer())

		// Boolean
		case glib.TYPE_BOOLEAN:
			if valueIface, err = glibValue.GoValue(); err == nil {
				if valueIface.(bool) {
					valueString = "true"
				} else {
					valueString = "false"
				}
			}

		// Need to be implemented
		default:
			err = fmt.Errorf("getStringCellValueByType: Type %s, not yet implemented\n", glibType[actualType])
		}
	}
	return
}

// changeTreeState: Modify the state of the entire tree starting at parent.
func (tvs *TreeViewStructure) changeChildrenTreeState(path string, parent *gtk.TreeIter, col int, goValue interface{}) {
	var err error
	var ok bool
	childIter := new(gtk.TreeIter)
	ok = parent != nil
	for ok {
		if err = tvs.SetColValue(parent, col, goValue); err == nil {
			if ok = tvs.Model.IterHasChild(parent); ok {
				if ok = tvs.Model.IterChildren(parent, childIter); ok {
					for ok {
						tvs.changeChildrenTreeState(path, childIter, col, goValue)
						ok = tvs.Model.IterNext(childIter)
					}
				}
			}
		}
	}
	if err != nil {
		fmt.Printf("Unable to changeChildrenTreeState (toggle) cell col: %d, path: %s, %s\n", col, path, err.Error())
	}
	return
}

// IsNotEmpty: Returns 'true' if the TreeView is not empty.
func (tvs *TreeViewStructure) IsNotEmpty() bool {
	_, ok := tvs.Model.GetIterFirst()
	return ok
}

// AllChildsCheckedState: verify all childs and childs of childs then return false
// if one of them correspond to 'state' parameter. Otherwise, true is returned.
func (tvs *TreeViewStructure) AllChildsCheckedState(parentIter *gtk.TreeIter, col int, state bool) bool {

	nChilds := tvs.Model.IterNChildren(parentIter)
	if nChilds > 0 {
		newChild := new(gtk.TreeIter)
		for idx := 0; idx < nChilds; idx++ {
			if tvs.Model.IterNthChild(newChild, parentIter, idx) {

				// if tvs.Model.IterHasChild(newChild) {
				// 	if !tvs.AllChildsCheckedState(newChild, col, state) {
				// 		return true
				// 	} else {
				// 		return false
				// 	}
				// }

				if state == tvs.GetColValue(newChild, col).(bool) {
					return false
				}
			}
		}
	}
	return true

	// if tvs.Model.IterHasChild(parentIter) {
	// if !tvs.AllChildsCheckedState(parentIter, col, state) {
	// 	return true
	// }
	// }
}

/* Same as above but with recurse capabilities */

// // allChildsChecked: verify all childs and return false if one of them
// // is unchecked. Otherwise, true is returned
// func (tvs *TreeViewStructure) allChildsChecked(parentIter *gtk.TreeIter, col int) bool {
// 	nChilds := tvs.Model.IterNChildren(parentIter)
// 	if nChilds > 0 {
// 		newChild := new(gtk.TreeIter)
// 		for idx := 0; idx < nChilds; idx++ {
// 			if tvs.Model.IterNthChild(newChild, parentIter, idx) {
// 				if !tvs.GetColValue(newChild, col).(bool) {
// 					return false
// 				}
// 			}
// 			if tvs.Model.IterHasChild(newChild) {
// 				tvs.allChildsChecked(newChild, col)
// 			}
// 		}
// 	}
// 	return true
// }

// ChildsPropagateColValue: Add value to parent and to all childs at specific column.
// If 'value' is nil, no modification will be done but all children are returned.
func (tvs *TreeViewStructure) ChildsPropagateColValue(parentIter *gtk.TreeIter, col int, value interface{}) (modifiedIters []*gtk.TreeIter, err error) {

	// Add value to parent if not nil
	if value != nil {
		if err = tvs.SetColValue(parentIter, col, value); err != nil {
			return
		}
	}
	modifiedIters = append(modifiedIters, parentIter)

	nChilds := tvs.Model.IterNChildren(parentIter)
	// Add value to his childs
	if nChilds > 0 {
		for idx := 0; idx < nChilds; idx++ {
			newChild := new(gtk.TreeIter)
			if tvs.Model.IterNthChild(newChild, parentIter, idx) {
				if value != nil {
					if err = tvs.SetColValue(newChild, col, value); err != nil {
						return
					}
				}
				if tvs.Model.IterHasChild(newChild) {
					if mi, err := tvs.ChildsPropagateColValue(newChild, col, value); err == nil {
						modifiedIters = append(modifiedIters, mi...)
					} else {
						return modifiedIters, err
					}
				}
				modifiedIters = append(modifiedIters, newChild)
			}
		}
	}
	return
}

/********************\
*    TEST Functions   *
* Not designed to be *
* used as it !!       *
**********************/

// func (tvs *TreeViewStructure) getDecendants(iter *gtk.TreeIter) (descendants [][]interface{}, err error) {
// 	// var parentPath, path *gtk.TreePath
// 	var iterDesc *gtk.TreeIter
// 	var ok bool = true
// 	var rowIface []interface{}

// 	ok = tvs.TreeStore.IterChildren(iter, iterDesc)
// 	for ok {
// 		rowIface, err = tvs.GetRowIface(iterDesc)
// 		descendants = append(descendants, rowIface)
// 		ok = tvs.TreeStore.IterNext(iterDesc)
// 	}
// 	return
// }

func (tvs *TreeViewStructure) selectRange(start, end *gtk.TreeIter) (err error) {
	var startPath, endPath *gtk.TreePath
	if startPath, err = tvs.ListStore.GetPath(start); err == nil {
		if endPath, err = tvs.ListStore.GetPath(end); err == nil {
			tvs.Selection.SelectRange(startPath, endPath)
		}
	}
	return err
}

func (tvs *TreeViewStructure) pathSelected(start *gtk.TreeIter) (err error) {
	var startPath *gtk.TreePath
	if startPath, err = tvs.ListStore.GetPath(start); err == nil {
		fmt.Println("iter", tvs.Selection.IterIsSelected(start))
		fmt.Println("path", tvs.Selection.PathIsSelected(startPath))
	}
	return err
}

func (tvs *TreeViewStructure) forEach() {
	var err error
	var model gtk.ITreeModel
	var ipath *gtk.TreePath
	var foreachFunc gtk.TreeModelForeachFunc
	foreachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) bool {
		if ipath, err = model.GetPath(iter); err == nil {
			fmt.Printf("path: %s, iter: %s\n", path.String(), ipath.String())
		} else {
			fmt.Println("error occured inside func: " + err.Error())
			return true
		}
		return false
	}
	if model, err = tvs.TreeView.GetModel(); err == nil {
		model.ToTreeModel().ForEach(foreachFunc)
	}
	if err != nil {
		fmt.Println("error occured outside func: " + err.Error())
	}
}

func (tvs *TreeViewStructure) idx() {
	var err error
	// var model *gtk.TreeModel
	var path, cpypath *gtk.TreePath
	if path, err = gtk.TreePathNewFirst(); err == nil {
		fmt.Printf("path: %s\n", path.String())
		path.AppendIndex(3)
		fmt.Printf("depth: %d\n", path.GetDepth())
		path.PrependIndex(6)
		fmt.Printf("depth to copy: %d:%s\n", path.GetDepth(), path.String())
		if cpypath, err = path.Copy(); err == nil {
			fmt.Printf("copied: %d:%s\n", cpypath.GetDepth(), cpypath.String())
			fmt.Printf("compared: %d\n", cpypath.Compare(cpypath))
			cpypath.Next()
			fmt.Printf("next: :%s\n", cpypath.String())
			cpypath.Prev()
			fmt.Printf("prev: :%s\n", cpypath.String())
			cpypath.Up()
			fmt.Printf("up: :%s\n", cpypath.String())
			cpypath.Down()
			fmt.Printf("down: :%s\n", cpypath.String())
			fmt.Printf("IsAncestor: :%v\n", cpypath.IsAncestor(path))
			fmt.Printf("IsDescendant: :%v\n", cpypath.IsDescendant(path))
			if path, err = gtk.TreePathNewFromIndicesv([]int{2, 3, 4, 7, 8}); err == nil {
				fmt.Printf("new indices: %d:%s\n", path.GetDepth(), path.String())
			}
		}
	}
	if err != nil {
		fmt.Println("error occured outside func: " + err.Error())
	}
}

func (tvs *TreeViewStructure) indices() {
	var err error
	var model gtk.ITreeModel
	var ipath, jpath *gtk.TreePath
	var foreachFunc gtk.TreeModelForeachFunc
	foreachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter) bool {
		if ipath, err = model.GetPath(iter); err == nil {
			indices := ipath.GetIndices()
			jpath, _ = gtk.TreePathNewFromIndicesv(indices)
			indices1 := jpath.GetIndices()
			fmt.Printf("indices %v -> pathString: %v -> indices %v\n", indices, jpath.String(), indices1)
		} else {
			fmt.Println("error occured inside func: " + err.Error())
			return true
		}
		return false
	}
	if model, err = tvs.TreeView.GetModel(); err == nil {
		model.ToTreeModel().ForEach(foreachFunc)
	}
	if err != nil {
		fmt.Println("error occured outside func: " + err.Error())
	}
}

// func (tvs *TreeViewStructure) getColsNames() (err error) {
// 	var glist *glib.List
// 	if glist, err = tvs.TreeView.GetColumns(); err == nil {
// 		for l := glist; l != nil; l = l.Next() {
// 			col := l.Data().(*gtk.TreeViewColumn)
// 			fmt.Println(col.GetTitle())
// 		}
// 	}
// 	if err != nil {
// 		err = errors.New("error occured while reading cols names: " + err.Error())
// 	}
// 	return err
// }
