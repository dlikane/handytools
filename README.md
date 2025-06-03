## Win installation

explorer hack to show classic menu
```azure
reg add HKCU\Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32 /ve /d "" /f
```
## Running samples
img pinterest -l fit -o output/pin.jpg -u https://au.pinterest.com/dlikane/dancer-photography/

go run ./cmd/img pinterest -l fit -o output/pin.jpg -u https://au.pinterest.com/dlikane/dancer-photography/ -e