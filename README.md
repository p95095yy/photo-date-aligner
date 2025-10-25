# Photo Date Aligner

A simple GUI tool written in Go to copy JPEG files and align their EXIF shooting dates sequentially by one-minute intervals.  
Uses a calendar picker for the start date and supports native folder selection dialogs.

## Features
- Cross-platform GUI (Fyne)
- Native folder selection (sqweek/dialog)
- EXIF date modification using go-exif and go-jpeg-image-structure
- Remembers the last selected folder

## Build

```bash
task build
```

## Usage
1. Run the application.
2. Select a folder containing JPEG files.
3. Pick the starting date from the calendar.
4. Click **Run** to copy all images into a new folder with incremented EXIF times.

## Dependencies
- [Fyne](https://github.com/fyne-io/fyne) — Cross-platform GUI framework  
- [sqweek/dialog](https://github.com/sqweek/dialog) — Native file/folder selection dialogs  
- [go-exif](https://github.com/dsoprea/go-exif) — EXIF metadata handler  
- [go-jpeg-image-structure](https://github.com/dsoprea/go-jpeg-image-structure) — JPEG parser/writer

## License
This project is licensed under the MIT License – see the [LICENSE](LICENSE) file for details.
