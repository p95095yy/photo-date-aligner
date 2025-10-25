package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/dsoprea/go-exif/v3"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"

	sqDialog "github.com/sqweek/dialog"
)

// --- Runボタン用ハンドラを生成する関数 ---
func makeRunHandler(selectedDate *time.Time, folderEntry *widget.Entry, selectedMode *string) func() {
	return func() {
		println("Selected date:", selectedDate.Format("2006-01-02"))
		srcFolder := folderEntry.Text
		println("Selected folder:", srcFolder)

		// --- フォルダ存在チェック ---
		if info, err := os.Stat(srcFolder); err != nil || !info.IsDir() {
			sqDialog.Message("The selected folder does not exist.").Title("Error").Error()
			return
		}

		parent := filepath.Dir(srcFolder)
		base := filepath.Base(srcFolder)
		dstFolder := filepath.Join(parent, base+"_update")

		println("Destination folder:", dstFolder)

		if srcFolder == "" {
			sqDialog.Message("No folder selected.").Title("Error").Error()
			return
		}
		if selectedDate.IsZero() {
			sqDialog.Message("No date/time selected.").Title("Error").Error()
			return
		}

		// すでに存在していたらエラー
		if _, err := os.Stat(dstFolder); err == nil {
			sqDialog.Message("Folder '%s' already exists.", dstFolder).Title("Error").Error()
			return
		}
		// フォルダを作成
		if err := os.Mkdir(dstFolder, 0755); err != nil {
			sqDialog.Message("Failed to create folder: %v", err).Title("Error").Error()
			return
		}

		// ---- 元フォルダ内の画像を収集 ----
		files := []string{}
		filepath.Walk(srcFolder, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			ext := strings.ToLower(filepath.Ext(info.Name()))
			if ext == ".jpg" || ext == ".jpeg" {
				files = append(files, path)
			}
			return nil
		})

		// ---- ソート順を選択モードで切り替え ----
		switch *selectedMode {
		case "Increment by filename order (Ascending)":
			sort.Slice(files, func(i, j int) bool {
				return filepath.Base(files[i]) < filepath.Base(files[j])
			})
		case "Increment by reverse filename order (Descending)":
			sort.Slice(files, func(i, j int) bool {
				return filepath.Base(files[i]) > filepath.Base(files[j])
			})
		case "Fix all timestamps":
			// 並び替え不要
		}

		// ---- ファイルを処理 ----
		count := 0
		for i, srcPath := range files {
			info, _ := os.Stat(srcPath)
			dstPath := filepath.Join(dstFolder, info.Name())

			// ファイルをコピー
			if err := copyFile(srcPath, dstPath); err != nil {
				fmt.Println("copy failed:", err)
				continue
			}

			// モードごとに撮影日時を決定
			var offsetTime time.Time
			switch *selectedMode {
			case "Increment by filename order (Ascending)",
				"Increment by reverse filename order (Descending)":
				offsetTime = selectedDate.Add(time.Duration(i) * time.Minute)
			case "Fix all timestamps":
				offsetTime = *selectedDate
			}

			// EXIF更新
			if err := updateExifDate(dstPath, offsetTime); err != nil {
				fmt.Println("EXIF update failed:", err)
			}

			count++
		}

		sqDialog.Message("%d files copied to:\n%s", count, dstFolder).Title("Completed").Info()
	}
}

func main() {
	a := app.NewWithID("io.github.p95095yy.photodatealigner")
	w := a.NewWindow("Photo Date Aligner")

	prefs := a.Preferences()

	var selectedDate time.Time
	selectedDate = time.Now()
	selectedDate = time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(), 0, 0, 0, 0, selectedDate.Location())

	label := widget.NewLabel("Start: " + selectedDate.Format("2006-01-02 15:04"))
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.Refresh()

	cal := widget.NewCalendar(time.Now(), func(t time.Time) {
		selectedDate = t
		selectedDate = time.Date(selectedDate.Year(), selectedDate.Month(), selectedDate.Day(), 0, 0, 0, 0, selectedDate.Location())

		label.SetText("Start: " + selectedDate.Format("2006-01-02 15:04"))
	})

	folderEntry := widget.NewEntry()
	folderEntry.SetPlaceHolder("Enter a folder or use Browse button")

	modeOptions := []string{
		"Fix all timestamps",
		"Increment by filename order (Ascending)",
		"Increment by reverse filename order (Descending)",
	}

	selectedMode := "Fix all timestamps"

	modeRadio := widget.NewRadioGroup(modeOptions, func(value string) {
		selectedMode = value
	})
	modeRadio.SetSelected(selectedMode)

	// --- Runボタン（最初は無効） ---
	runBtn := widget.NewButton("Run", makeRunHandler(&selectedDate, folderEntry, &selectedMode))
	runBtn.Disable()

	// --- フォルダ選択 ---
	selectFolderBtn := widget.NewButton("Browse", func() {
		lastFolder := prefs.StringWithFallback("lastFolder", "")
		path, err := sqDialog.Directory().Title("Select a folder").SetStartDir(lastFolder).Browse()
		if err == nil && path != "" {
			folderEntry.SetText(path)
			prefs.SetString("lastFolder", path)
			runBtn.Enable() // フォルダ選択後に有効化
		}
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Photo Folder:"),
		container.NewBorder(nil, nil, nil, selectFolderBtn, folderEntry),
		widget.NewSeparator(),
		label,
		cal,
		widget.NewSeparator(),
		widget.NewLabel("Timestamp Update Mode:"),
		modeRadio,
		widget.NewSeparator(),
		runBtn,
	))

	w.Resize(fyne.NewSize(600, 260))
	w.ShowAndRun()
}

// ---- ファイルをコピーする ----
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// ---- EXIF日時更新 ----
func updateExifDate(path string, t time.Time) error {
	jmp := jpegstructure.NewJpegMediaParser()
	ec, err := jmp.ParseFile(path)
	if err != nil {
		return fmt.Errorf("failed to parse jpeg: %v", err)
	}
	sl := ec.(*jpegstructure.SegmentList)

	rootBuilder, err := sl.ConstructExifBuilder()
	if err != nil {
		return fmt.Errorf("failed to construct exif builder: %v", err)
	}

	exifBuilder, err := exif.GetOrCreateIbFromRootIb(rootBuilder, "IFD0/Exif0")
	if err != nil {
		return fmt.Errorf("failed to get exif builder: %v", err)
	}

	formatted := t.Format("2006:01:02 15:04:05")
	if err := exifBuilder.SetStandardWithName("DateTimeOriginal", formatted); err != nil {
		if addErr := exifBuilder.AddStandardWithName("DateTimeOriginal", formatted); addErr != nil {
			fmt.Printf("failed to add DateTimeOriginal: %v\n", addErr)
		}
	}

	if err := sl.SetExif(rootBuilder); err != nil {
		return fmt.Errorf("failed to set exif: %v", err)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := sl.Write(out); err != nil {
		return fmt.Errorf("failed to write jpeg: %v", err)
	}

	return nil
}
