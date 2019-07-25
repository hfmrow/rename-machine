// treeView.go

/*
*	©2019 H.F.M. MIT license
*	Handle ListStore gotk3 object.
 */

package gtk3Import

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

/**********************
*	TreeView section
 *********************/
// TreeViewGetRowsNb: Get selected rows numbers
func TreeViewGetRowsNb(tw *gtk.TreeView, ls *gtk.ListStore) (rowsNb []int) {
	selected, err := tw.GetSelection()
	if err != nil {
		fmt.Println(err)
	}
	rows := selected.GetSelectedRows(ls)
	for l := rows; l != nil; l = l.Next() {
		path := l.Data().(*gtk.TreePath)
		rowsNb = append(rowsNb, path.GetIndices()[0])
	}
	return rowsNb
}

// TreeViewListStoreSetup: Setup tree view columns and the list store that holds its data.
// options set combine multiselection and editable property.
func TreeViewListStoreSetup(treeView *gtk.TreeView, multiSelect bool, columns [][]string, editable ...bool) (listStore *gtk.ListStore) {

	var err error
	// Set selection mode
	if multiSelect {
		selMode, err := treeView.GetSelection()
		if err != nil {
			log.Fatal("Unable to set selection mode:", err)
		}
		selMode.SetMode(gtk.SELECTION_MULTIPLE)
	}
	// Build columns
	for colIdx, colDescr := range columns {
		var attributes []string
		colName := colDescr[0]
		for idx := 1; idx < len(colDescr); idx++ {
			attributes = append(attributes, colDescr[idx])
		}
		treeView.AppendColumn(TreeViewCreateTextColumn(colName, colIdx, attributes, editable...))
		column := treeView.GetColumn(colIdx)
		column.SetResizable(true)      // Set column resizable
		column.SetSortColumnID(colIdx) // Set column sortable
	}

	// Set glibType to string for the list store.
	colType := make([]glib.Type, len(columns))
	for idx, _ := range colType {
		colType[idx] = glib.TYPE_STRING
	}
	// Creating a list store. This is what holds the data that will be shown on our tree view.
	listStore, err = gtk.ListStoreNew(colType...)
	if err != nil {
		log.Fatal("Unable to create list store:", err)
	}
	treeView.SetModel(listStore)
	return listStore
}

// TreeViewCreateTextColumn: Add a column to the tree view
func TreeViewCreateTextColumn(title string, id int, cellAttributes []string, editable ...bool) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if len(editable) != 0 {
		cellRenderer.SetProperty("editable", editable[0])
	}

	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, cellAttributes[0], id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}

	for idx := 1; idx < len(cellAttributes); idx++ {
		column.AddAttribute(cellRenderer, cellAttributes[idx], id)
	}

	return column
}

// ListStoreAddRow: Append a row to the listStore
func ListStoreAddRow(listStore *gtk.ListStore, row []string) (iter *gtk.TreeIter) {
	// Get an iterator for a new row at the end of the list store
	iter = listStore.Append()
	// Set columns index and interface of strings slice
	colIdx := make([]int, len(row))
	var interfaceSlice []interface{} = make([]interface{}, len(row))
	for idx, col := range row {
		interfaceSlice[idx] = col
		colIdx[idx] = idx
	}
	// Set the contents of the list store row that the iterator represents
	err := listStore.Set(iter, colIdx, interfaceSlice)
	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
	return iter
}

/**********************
*	TreeStore section
 *********************/
// TreeViewTreeStoreSetup: Creates a TreeView and the TreeStore that holds its data
func TreeViewTreeStoreSetup(treeView *gtk.TreeView, multiSelect bool, columns [][]string) (treeStore *gtk.TreeStore) {
	var err error
	// Set selection mode
	if multiSelect {
		selMode, err := treeView.GetSelection()
		if err != nil {
			log.Fatal("Unable to set selection mode:", err)
		}
		selMode.SetMode(gtk.SELECTION_MULTIPLE)
	}
	// Build columns
	for colIdx, colDescr := range columns {
		var attributes []string
		colName := colDescr[0]
		for idx := 1; idx < len(colDescr); idx++ {
			attributes = append(attributes, colDescr[idx])
		}
		treeView.AppendColumn(TreeViewCreateTextColumn(colName, colIdx, attributes))
		column := treeView.GetColumn(colIdx)
		column.SetResizable(true)      // Set column resizable
		column.SetSortColumnID(colIdx) // Set column sortable
	}
	// Set glibType to string for the TreeStore.
	colType := make([]glib.Type, len(columns))
	for idx, _ := range colType {
		colType[idx] = glib.TYPE_STRING
	}
	// Creating a list store. This is what holds the data that will be shown on our tree view.
	treeStore, err = gtk.TreeStoreNew(colType...)
	if err != nil {
		log.Fatal("Unable to create TreeStore:", err)
	}
	treeView.SetModel(treeStore)

	//	TEST_TreeStore(treeStore) /*TEST*/

	return treeStore
}

// TreeStoreAddRow: Append a toplevel row to the tree store for the tree view
func TreeStoreAddRow(treeStore *gtk.TreeStore, row string) *gtk.TreeIter {
	return TreeStoreAddSubRow(treeStore, nil, []string{row}...)
}

// TreeStoreAddSubRow: Append a sub row to the tree store for the tree view
func TreeStoreAddSubRow(treeStore *gtk.TreeStore, iter *gtk.TreeIter, row ...string) *gtk.TreeIter {
	// Get an iterator for a new row at the end of the list store
	subIter := treeStore.Append(iter)
	// Set columns index and interface of strings slice
	//	var interfaceSlice []interface{} = make([]interface{}, len(row))
	for idx, line := range row {
		//		interfaceSlice[idx] = col
		err := treeStore.SetValue(subIter, idx, line)
		//	err := treeStore.Set(subIter, colIdx..., interfaceSlice)
		if err != nil {
			log.Fatal("Unable to add subRow:", err)
		}
	}
	return subIter
}

/*	TEST purpose with pango markup style */
func TEST_TreeStore(treeStore *gtk.TreeStore) {
	filename := "/media/oiuytreza/storage/Documents/dev/go/src/github.com/gfdbvckjh/sandr/files-tst.tmp"
	Values := []string{"<b>5</b> this is the fifth line"}
	Values1 := []string{"<b>8</b> this is the eighth line"}
	Values2 := []string{markup("  12  ", "bgc", "#E4DDDD") + "this is the twelfth line"}

	iter := TreeStoreAddRow(treeStore, filename)
	TreeStoreAddSubRow(treeStore, iter, Values...)
	TreeStoreAddSubRow(treeStore, iter, Values1...)
	TreeStoreAddSubRow(treeStore, iter, Values2...)

}
