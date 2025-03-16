Set objFSO = CreateObject("Scripting.FileSystemObject")
tempFile = "C:\Users\User\go\bin\batchrename_files.txt"

' Check if the file exists and delete it
If objFSO.FileExists(tempFile) Then
    On Error Resume Next
    objFSO.DeleteFile tempFile, True
    If Err.Number <> 0 Then
        WScript.Echo "Error: Cannot delete temporary file. Permission denied."
        WScript.Quit 1
    End If
    On Error GoTo 0
End If

' Create a new temp file
Set objFile = objFSO.CreateTextFile(tempFile, True)

Set args = WScript.Arguments
If args.Count = 0 Then
    WScript.Echo "Error: No files selected!"
    WScript.Quit 1
End If

' Write all files to temp file
For i = 0 To args.Count - 1
    objFile.WriteLine args(i)
Next
objFile.Close

' Run batch rename
Set objShell = CreateObject("WScript.Shell")
objShell.Run "C:\Users\User\go\bin\batchrename.bat", 0, False
