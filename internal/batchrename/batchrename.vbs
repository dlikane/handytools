Set objShell = CreateObject("WScript.Shell")
Set args = WScript.Arguments
command = """C:\Users\User\go\bin\batchrename.bat"""

For i = 0 To args.Count - 1
    command = command & " """ & args(i) & """"
Next

objShell.Run command, 0, False ' 0 = No Window, False = Don't Wait
