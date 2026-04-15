package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	fynecanvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type eatWhatApp struct {
	app    fyne.App
	window fyne.Window

	manager       *FoodManager
	selectedIndex int
	menuName      string
	menuPath      string

	inputBox     *widget.Entry
	resultTitle  *fynecanvas.Text
	resultStatus *fynecanvas.Text
	resultText   *fynecanvas.Text
	countLabel   *widget.Label
	fileLabel    *widget.Label
	memoryLabel  *widget.Label
	list         *widget.List
}

func newEatWhatApp() *eatWhatApp {
	a := app.NewWithID(appID)
	a.Settings().SetTheme(theme.LightTheme())

	window := a.NewWindow("今天吃什么？")
	window.Resize(fyne.NewSize(980, 640))
	window.CenterOnScreen()

	e := &eatWhatApp{
		app:           a,
		window:        window,
		manager:       NewFoodManager(),
		selectedIndex: -1,
	}

	e.restoreState()
	e.window.SetContent(e.buildUI())
	e.refreshUI()
	return e
}

func (e *eatWhatApp) Run() {
	e.window.ShowAndRun()
}

func (e *eatWhatApp) buildUI() fyne.CanvasObject {
	eyebrow := fynecanvas.NewText("EAT WHAT", color.NRGBA{R: 136, G: 104, B: 82, A: 255})
	eyebrow.Alignment = fyne.TextAlignCenter
	eyebrow.TextSize = 16
	eyebrow.TextStyle = fyne.TextStyle{Bold: true}

	title := widget.NewLabel("今天吃什么？")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	e.countLabel = widget.NewLabel("当前共有 0 个选项")
	e.countLabel.Alignment = fyne.TextAlignCenter
	e.fileLabel = widget.NewLabel("当前菜单：未命名")
	e.fileLabel.Alignment = fyne.TextAlignCenter
	e.memoryLabel = widget.NewLabel("菜单会自动保存到本地 txt")
	e.memoryLabel.Alignment = fyne.TextAlignCenter

	e.resultTitle = fynecanvas.NewText("今日推荐", color.NRGBA{R: 184, G: 92, B: 41, A: 255})
	e.resultTitle.Alignment = fyne.TextAlignCenter
	e.resultTitle.TextSize = 30
	e.resultTitle.TextStyle = fyne.TextStyle{Bold: true}

	e.resultStatus = fynecanvas.NewText("", color.NRGBA{R: 166, G: 96, B: 66, A: 255})
	e.resultStatus.Alignment = fyne.TextAlignCenter
	e.resultStatus.TextSize = 34
	e.resultStatus.TextStyle = fyne.TextStyle{Bold: true}

	e.resultText = fynecanvas.NewText("", color.NRGBA{R: 188, G: 86, B: 50, A: 255})
	e.resultText.Alignment = fyne.TextAlignCenter
	e.resultText.TextSize = 48
	e.resultText.TextStyle = fyne.TextStyle{Bold: true}
	e.resetResultText()

	resultCardBg := fynecanvas.NewRectangle(color.NRGBA{R: 252, G: 244, B: 233, A: 255})
	resultCardBg.CornerRadius = 36

	resultContent := container.NewCenter(
		container.NewVBox(
			container.NewCenter(e.resultTitle),
			widget.NewSeparator(),
			layout.NewSpacer(),
			container.NewCenter(e.resultStatus),
			container.NewCenter(e.resultText),
			layout.NewSpacer(),
		),
	)

	resultCard := container.NewMax(
		resultCardBg,
		container.NewPadded(resultContent),
	)

	e.list = e.newOptionList()
	e.inputBox = widget.NewMultiLineEntry()
	e.inputBox.SetMinRowsVisible(9)
	e.inputBox.SetPlaceHolder("请输入食物选项，一行一个")

	createButton := widget.NewButtonWithIcon("新建菜单", theme.DocumentCreateIcon(), e.showCreateMenuDialog)
	importButton := widget.NewButtonWithIcon("导入 TXT", theme.FolderOpenIcon(), e.handleImportOptions)
	addButton := widget.NewButtonWithIcon("添加选项", theme.ContentAddIcon(), e.showAddOptionsDialog)
	addButton.Importance = widget.HighImportance
	manageButton := widget.NewButtonWithIcon("管理菜单", theme.MenuIcon(), e.showManageOptionsDialog)
	pickButton := widget.NewButton("吃什么", e.handlePickRandom)
	pickButton.Importance = widget.HighImportance

	heroBg := fynecanvas.NewRectangle(color.NRGBA{R: 255, G: 252, B: 247, A: 248})
	heroBg.CornerRadius = 38

	heroCard := container.NewMax(
		heroBg,
		container.NewPadded(
			container.NewVBox(
				container.NewCenter(eyebrow),
				container.NewCenter(title),
				container.NewCenter(e.fileLabel),
				container.NewCenter(e.countLabel),
				container.NewCenter(e.memoryLabel),
				widget.NewSeparator(),
				container.New(
					layout.NewCustomPaddedLayout(0, 0, 0, 0),
					container.NewGridWrap(fyne.NewSize(520, 340), resultCard),
				),
				pickButton,
				container.NewGridWithColumns(2, createButton, importButton),
				container.NewGridWithColumns(2, addButton, manageButton),
			),
		),
	)

	backdrop := fynecanvas.NewRectangle(color.NRGBA{R: 243, G: 239, B: 234, A: 255})
	blobOne := fynecanvas.NewCircle(color.NRGBA{R: 234, G: 210, B: 191, A: 95})
	blobTwo := fynecanvas.NewCircle(color.NRGBA{R: 214, G: 227, B: 221, A: 88})
	blobThree := fynecanvas.NewCircle(color.NRGBA{R: 242, G: 224, B: 202, A: 95})
	background := container.NewWithoutLayout(backdrop, blobOne, blobTwo, blobThree)
	backdrop.Resize(fyne.NewSize(2200, 1600))
	blobOne.Resize(fyne.NewSize(420, 420))
	blobOne.Move(fyne.NewPos(-60, -30))
	blobTwo.Resize(fyne.NewSize(380, 380))
	blobTwo.Move(fyne.NewPos(760, 120))
	blobThree.Resize(fyne.NewSize(320, 320))
	blobThree.Move(fyne.NewPos(110, 480))

	return container.NewMax(
		background,
		container.NewCenter(
			container.New(
				layout.NewCustomPaddedLayout(36, 36, 120, 120),
				heroCard,
			),
		),
	)
}

func (e *eatWhatApp) newOptionList() *widget.List {
	list := widget.NewList(
		func() int {
			return e.manager.Count()
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("template")
			label.Wrapping = fyne.TextWrapOff
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				label,
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			row := item.(*fyne.Container)
			label := row.Objects[1].(*widget.Label)
			label.SetText(e.manager.Get(id))
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		e.selectedIndex = id
	}
	return list
}

func (e *eatWhatApp) refreshUI() {
	e.countLabel.SetText("当前共有 " + strconv.Itoa(e.manager.Count()) + " 个选项")
	e.fileLabel.SetText("当前菜单：" + e.currentMenuName())
	e.list.Refresh()

	if e.manager.Count() == 0 {
		e.selectedIndex = -1
		return
	}
	if e.selectedIndex >= e.manager.Count() {
		e.selectedIndex = e.manager.Count() - 1
	}
}

func (e *eatWhatApp) showCreateMenuDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("例如：工作日午餐 / 周末快乐餐")

	createDialog := dialog.NewForm(
		"新建菜单",
		"创建",
		"取消",
		[]*widget.FormItem{
			widget.NewFormItem("菜单名称", nameEntry),
		},
		func(ok bool) {
			if !ok {
				return
			}
			name := normalizeOption(nameEntry.Text)
			if name == "" {
				dialog.ShowInformation("提示", "菜单名称不能为空。", e.window)
				return
			}
			if err := e.setManagedMenu(name); err != nil {
				dialog.ShowError(err, e.window)
				return
			}
			e.manager.Clear()
			e.selectedIndex = -1
			e.persistState()
			e.refreshUI()
			e.resetResultText()
			e.showAddOptionsDialog()
		},
		e.window,
	)
	createDialog.Resize(fyne.NewSize(480, 220))
	createDialog.Show()
}

func (e *eatWhatApp) showAddOptionsDialog() {
	menuNameEntry := widget.NewEntry()
	menuNameEntry.SetPlaceHolder("例如：工作日午餐 / 周末快乐餐")
	if e.menuName != "" {
		menuNameEntry.SetText(e.menuName)
	}
	hint := widget.NewLabel("一行一个，例如：麻辣烫 / 牛肉面 / 寿司")
	hint.Wrapping = fyne.TextWrapWord

	topContent := []fyne.CanvasObject{
		widget.NewLabel("当前内容会直接写入这份菜单对应的 txt 文件。"),
	}
	if e.menuName == "" {
		topContent = append(topContent, widget.NewLabel("还没有菜单时，可以直接在这里创建并添加选项。"))
		topContent = append(topContent, widget.NewForm(
			widget.NewFormItem("菜单名称", menuNameEntry),
		))
	} else {
		topContent = append(topContent, widget.NewLabel("当前菜单："+e.currentMenuName()))
	}
	topContent = append(topContent, hint)

	content := container.NewBorder(
		container.NewVBox(topContent...),
		container.NewGridWithColumns(
			2,
			widget.NewButtonWithIcon("添加", theme.ConfirmIcon(), func() {
				if e.menuName == "" {
					name := normalizeOption(menuNameEntry.Text)
					if name == "" {
						dialog.ShowInformation("提示", "请先填写菜单名称。", e.window)
						return
					}
					if err := e.setManagedMenu(name); err != nil {
						dialog.ShowError(err, e.window)
						return
					}
				}
				e.handleAddOptions()
			}),
			widget.NewButtonWithIcon("清空输入", theme.ContentClearIcon(), e.handleClearInput),
		),
		nil,
		nil,
		e.inputBox,
	)

	addDialog := dialog.NewCustom("添加选项", "关闭", content, e.window)
	addDialog.Resize(fyne.NewSize(620, 420))
	addDialog.Show()
}

func (e *eatWhatApp) handleAddOptions() {
	lines := strings.Split(e.inputBox.Text, "\n")
	added := e.manager.AddOptions(lines)
	e.persistState()
	e.refreshUI()

	if added == 0 {
		dialog.ShowInformation("提示", "没有新增有效选项。可能都是空行或重复项。", e.window)
		return
	}

	e.inputBox.SetText("")
	dialog.ShowInformation("添加成功", "成功添加 "+strconv.Itoa(added)+" 个选项。", e.window)
}

func (e *eatWhatApp) handleClearInput() {
	e.inputBox.SetText("")
}

func (e *eatWhatApp) handleImportOptions() {
	open := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, e.window)
			return
		}
		if reader == nil {
			return
		}

		path := reader.URI().Path()
		fileName := reader.URI().Name()
		_ = reader.Close()

		lines, readErr := loadOptionsFromTxt(path)
		if readErr != nil {
			dialog.ShowError(readErr, e.window)
			return
		}

		e.showImportModeDialog(lines, fileName)
	}, e.window)

	open.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
	open.Show()
}

func (e *eatWhatApp) showImportModeDialog(lines []string, fileName string) {
	description := widget.NewRichTextFromMarkdown(
		fmt.Sprintf("已读取 `%s`，导入后会以它的文件名作为当前菜单名。", fileName),
	)

	var modeDialog dialog.Dialog
	appendButton := widget.NewButtonWithIcon("在原菜单基础上增加", theme.ContentAddIcon(), func() {
		if err := e.setManagedMenu(defaultMenuNameFromFile(fileName)); err != nil {
			dialog.ShowError(err, e.window)
			return
		}
		added := e.manager.AddOptions(lines)
		e.persistState()
		e.refreshUI()
		modeDialog.Hide()
		dialog.ShowInformation("导入完成", "成功新增 "+strconv.Itoa(added)+" 个选项。", e.window)
	})
	overwriteButton := widget.NewButtonWithIcon("覆盖原菜单", theme.DeleteIcon(), func() {
		if err := e.setManagedMenu(defaultMenuNameFromFile(fileName)); err != nil {
			dialog.ShowError(err, e.window)
			return
		}
		added := e.manager.ReplaceOptions(lines)
		e.selectedIndex = -1
		e.persistState()
		e.refreshUI()
		e.resetResultText()
		modeDialog.Hide()
		dialog.ShowInformation("导入完成", "已用新文件覆盖菜单，共载入 "+strconv.Itoa(added)+" 个选项。", e.window)
	})
	cancelButton := widget.NewButton("取消", func() {
		modeDialog.Hide()
	})

	content := container.NewVBox(
		description,
		widget.NewLabel("追加模式会自动去重；覆盖模式会清空当前菜单后重新导入。"),
		appendButton,
		overwriteButton,
		cancelButton,
	)

	modeDialog = dialog.NewCustom("选择导入方式", "", content, e.window)
	modeDialog.Resize(fyne.NewSize(560, 290))
	modeDialog.Show()
}

func (e *eatWhatApp) showManageOptionsDialog() {
	listTitle := widget.NewLabel("管理当前菜单")
	listTitle.TextStyle = fyne.TextStyle{Bold: true}
	countLabel := widget.NewLabel("当前共有 " + strconv.Itoa(e.manager.Count()) + " 个选项")
	fileLabel := widget.NewLabel("当前菜单：" + e.currentMenuName())

	content := container.NewBorder(
		container.NewVBox(
			listTitle,
			countLabel,
			fileLabel,
		),
		container.NewGridWithColumns(
			3,
			widget.NewButtonWithIcon("重命名", theme.DocumentCreateIcon(), e.showRenameMenuDialog),
			widget.NewButtonWithIcon("删除选中", theme.DeleteIcon(), e.handleRemoveSelected),
			widget.NewButtonWithIcon("清空全部", theme.ContentClearIcon(), e.handleClearAll),
		),
		nil,
		nil,
		container.NewMax(e.list),
	)

	manageDialog := dialog.NewCustom("管理菜单", "关闭", container.NewPadded(content), e.window)
	manageDialog.Resize(fyne.NewSize(760, 620))
	manageDialog.Show()
}

func (e *eatWhatApp) showRenameMenuDialog() {
	if !e.ensureMenuReady() {
		return
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("例如：工作日午餐 / 周末快乐餐")
	nameEntry.SetText(e.currentMenuName())

	renameDialog := dialog.NewForm(
		"重命名菜单",
		"保存",
		"取消",
		[]*widget.FormItem{
			widget.NewFormItem("菜单名称", nameEntry),
		},
		func(ok bool) {
			if !ok {
				return
			}
			name := normalizeOption(nameEntry.Text)
			if name == "" {
				dialog.ShowInformation("提示", "菜单名称不能为空。", e.window)
				return
			}
			if err := e.renameMenu(name); err != nil {
				dialog.ShowError(err, e.window)
				return
			}
			e.persistState()
			e.refreshUI()
		},
		e.window,
	)
	renameDialog.Resize(fyne.NewSize(480, 220))
	renameDialog.Show()
}

func (e *eatWhatApp) handleRemoveSelected() {
	if e.selectedIndex < 0 || e.selectedIndex >= e.manager.Count() {
		dialog.ShowInformation("提示", "请先在列表中选中一个选项。", e.window)
		return
	}

	name := e.manager.Get(e.selectedIndex)
	e.manager.RemoveAt(e.selectedIndex)
	e.selectedIndex = -1
	e.persistState()
	e.refreshUI()
	dialog.ShowInformation("删除成功", "已删除选项："+name, e.window)
}

func (e *eatWhatApp) handleClearAll() {
	if e.manager.Count() == 0 {
		dialog.ShowInformation("提示", "当前没有可清空的选项。", e.window)
		return
	}

	dialog.ShowConfirm("确认清空", "确定要清空当前菜单中的全部选项吗？", func(ok bool) {
		if !ok {
			return
		}
		e.manager.Clear()
		e.selectedIndex = -1
		e.persistState()
		e.refreshUI()
		e.resetResultText()
	}, e.window)
}

func (e *eatWhatApp) handlePickRandom() {
	if e.manager.Count() == 0 {
		dialog.ShowInformation("提示", "请先新建菜单或导入 txt，再添加选项。", e.window)
		return
	}

	options := e.manager.All()
	iterations := 12
	interval := 70 * time.Millisecond

	go func() {
		for i := 0; i < iterations; i++ {
			item := options[rand.Intn(len(options))]
			fyne.Do(func() {
				e.setResultDisplay("正在帮你选", item, 26, resultTextSize(item))
			})
			time.Sleep(interval)
		}

		finalResult, err := e.manager.PickRandom()
		fyne.Do(func() {
			if err != nil {
				dialog.ShowError(err, e.window)
				return
			}
			e.setResultDisplay("今天吃", finalResult, 30, resultTextSize(finalResult))
		})
	}()
}

func (e *eatWhatApp) restoreState() {
	state, err := loadAppState()
	if err != nil {
		return
	}

	e.manager.ReplaceOptions(state.Options)
	e.menuName = state.MenuName
	e.menuPath = state.MenuPath
	if e.menuName == "" && state.CurrentFile != "" {
		e.menuName = defaultMenuNameFromFile(state.CurrentFile)
	}
	if e.menuPath == "" && e.menuName != "" {
		path, pathErr := managedMenuPath(e.menuName)
		if pathErr == nil {
			e.menuPath = path
		}
	}
}

func (e *eatWhatApp) persistState() {
	state := appState{
		Options:  e.manager.All(),
		MenuName: e.menuName,
		MenuPath: e.menuPath,
	}
	if err := saveAppState(state); err != nil {
		dialog.ShowError(err, e.window)
		return
	}
	if err := saveOptionsToTxt(e.menuPath, state.Options); err != nil {
		dialog.ShowError(err, e.window)
	}
}

func (e *eatWhatApp) currentMenuName() string {
	if e.menuName != "" {
		return e.menuName
	}
	return "未命名"
}

func defaultMenuNameFromFile(fileName string) string {
	dot := strings.LastIndex(fileName, ".")
	if dot <= 0 {
		return fileName
	}
	return fileName[:dot]
}

func (e *eatWhatApp) ensureMenuReady() bool {
	if e.menuName != "" && e.menuPath != "" {
		return true
	}
	e.showCreateMenuDialog()
	return false
}

func (e *eatWhatApp) setManagedMenu(name string) error {
	path, err := managedMenuPath(name)
	if err != nil {
		return err
	}
	if err := ensureMenuFile(path); err != nil {
		return err
	}
	e.menuName = name
	e.menuPath = path
	return nil
}

func (e *eatWhatApp) renameMenu(name string) error {
	newPath, err := managedMenuPath(name)
	if err != nil {
		return err
	}
	renamedPath, err := renameManagedMenuFile(e.menuPath, newPath, e.manager.All())
	if err != nil {
		return err
	}
	e.menuName = name
	e.menuPath = renamedPath
	return nil
}

func (e *eatWhatApp) resetResultText() {
	e.setResultDisplay("", "点击按钮开始随机抽取", 26, 22)
}

func (e *eatWhatApp) setResultDisplay(status string, text string, statusSize float32, textSize float32) {
	e.resultStatus.Text = status
	e.resultStatus.TextSize = statusSize
	e.resultStatus.Refresh()
	e.resultText.Text = text
	e.resultText.TextSize = textSize
	e.resultText.Refresh()
}

func resultTextSize(text string) float32 {
	switch n := utf8.RuneCountInString(text); {
	case n <= 4:
		return 48
	case n <= 8:
		return 40
	case n <= 12:
		return 32
	default:
		return 26
	}
}
