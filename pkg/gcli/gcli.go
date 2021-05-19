package gcli

// type GDriveCLI struct {
// 	Service         *drive.Service
// 	FilesToDownload map[string]*drive.File
// 	SizeBytes       int
// }

// func (cli GDriveCLI) Display() {

// 	app := tview.NewApplication()
// files, err := cli.GetFiles(folderID)
// if err != nil {
// 	fmt.Printf("ERROR1: %s", err)
// }
// table, err := cli.GenerateTable(app, files)
// if err != nil {
// 	//return err
// }
// if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
// 	//return err
// }
// //return nil

// }

// func (cli GDriveCLI) GenerateTable(app *tview.Application, files *drive.FileList) (*tview.Table, error) {
// 	table := tview.NewTable().SetBorders(true).SetEvaluateAllRows(true)

// 	table.SetCell(0, 0,
// 		tview.NewTableCell("File").
// 			SetTextColor(tcell.ColorYellow).
// 			SetAlign(tview.AlignCenter))

// 	rows := len(files.Files)
// 	for r := 1; r <= rows; r++ {
// 		file := files.Files[r-1]
// 		color := tcell.ColorWhite
// 		_, ok := cli.FilesToDownload[file.Id]
// 		if ok {
// 			color = tcell.ColorRed
// 		}
// 		table.SetCell(r, 0,
// 			tview.NewTableCell(file.Name).
// 				SetTextColor(color).
// 				SetAlign(tview.AlignCenter))
// 	}
// 	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
// 		if key == tcell.KeyEscape || key == tcell.KeyBackspace {
// 			app.Stop()
// 		}
// 		if key == tcell.KeyEnter {
// 			table.SetSelectable(true, true)
// 		}
// 	}).SetSelectedFunc(func(row int, column int) {

// 		//table.SetSelectable(false, false)
// 		file := files.Files[row-1]
// 		if file.MimeType == "application/vnd.google-apps.folder" {
// 			app.Suspend(func() {
// 				cli.Display(file.Id)
// 			})
// 		} else {
// 			cli.FilesToDownload[file.Id] = file
// 			table.GetCell(row, column).SetTextColor(tcell.ColorRed)
// 		}

// 	})
// 	return table, nil
// }
