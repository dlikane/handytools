@echo off
setlocal enabledelayedexpansion

:: Define temp file
set tempFile=C:\Users\User\go\bin\batchrename_files.txt

:: Check if the file exists
if not exist "%tempFile%" (
    echo Error: File list not found!
    exit /b
)

:: Read all collected filenames
set cmd="C:\Users\%USERNAME%\go\bin\img.exe" batchrename
for /f "usebackq tokens=*" %%F in ("%tempFile%") do (
    set cmd=!cmd! "%%F"
)

:: Execute command
!cmd!

:: Cleanup
del "%tempFile%"
