# Photo Date Aligner

## Defender の警告画面で実行する方法／How to Run the App When Defender Shows a Warning

初回起動すると  **「Windows はお使いの PC を保護しました」** という画面が表示されることがあります。  
When you start the app for the first time,  you may see a screen saying **“Windows protected your PC.”**

### なぜ出るのか／Why This Appears
このアプリは個人開発で、**コード署名証明書を付けていません**。
そのため Windows が「発行元が確認できないアプリ」と判断し、安全確認のため警告画面が表示されます。  
This app is individually developed and **does not include a code-signing certificate**.
Windows therefore cannot verify the publisher and shows this warning
as a safety confirmation.

怪しい処理はしていません。心配な場合はコードをご確認ください。  
We do not perform any suspicious operations. If you have concerns, please review the source code.

### 実行手順／How to Run It
1. 警告画面左側の **「詳細情報」** をクリックします。／Click **“More info”** on the left side of the warning screen.  
3. 下部に **「実行」** ボタンが表示されます。／A new button labeled **“Run anyway”** will appear at the bottom.  
3. **「実行」** をクリックするとアプリが起動します。／Click **“Run anyway”** to launch the app.

## Dependencies
- [Fyne](https://github.com/fyne-io/fyne) — Cross-platform GUI framework  
- [fyne-datepicker](https://github.com/sdassow/fyne-datepicker) — DatePicker widget for Fyne (used for date selection UI)
- [sqweek/dialog](https://github.com/sqweek/dialog) — Native file/folder selection dialogs  
- [go-exif](https://github.com/dsoprea/go-exif) — EXIF metadata handler  
- [go-jpeg-image-structure](https://github.com/dsoprea/go-jpeg-image-structure) — JPEG parser/writer

## License
This project is licensed under the MIT License – see the [LICENSE](LICENSE) file for details.
