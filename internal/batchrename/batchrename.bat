@echo off
rem Goes to c:\Users\User\Go
setlocal enabledelayedexpansion

:: Define temp file for storing file paths
set tempFile=%TEMP%\batchrename_files.txt

echo tempFile "%tempFile%"

:: If temp file doesn't exist, create it and mark this as the first execution
if not exist "%tempFile%" (
    echo Creating temp file...
    echo. > "%tempFile%"
    set firstExecution=1
) else (
    set firstExecution=0
)

:: Append the current file to the temp file
for %%F in (%*) do (
    echo append "%%F"
    echo "%%F" >> "%tempFile%"
)

:: If this is not the first execution, just exit
if "%firstExecution%"=="0" exit /b

:: Wait to ensure all files are written before processing
ping -n 2 127.0.0.1 >nul

:: Processing batch rename
echo Processing batch rename...

:: Read all collected filenames
set cmd="C:\Users\%USERNAME%\go\bin\img.exe" batchrename
for /f "usebackq tokens=*" %%F in ("%tempFile%") do (
    set cmd=!cmd! %%F
)

:: Show the final command for debugging
echo !cmd!

:: Execute the command
!cmd!

:: Cleanup
del "%tempFile%"
