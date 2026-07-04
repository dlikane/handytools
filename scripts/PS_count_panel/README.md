# Open Files Counter — Photoshop panel

A small Photoshop CEP panel that shows the number of open documents and a
clickable list of their names. It updates **live**, even while a Batch/action
is closing files one by one.

Tested on Adobe Photoshop 2026 (Windows). Should also work on 2021–2025.

---

## Contents

```
PS_count_panel/
├── extension/               ← the CEP panel
│   ├── index.html           ← panel UI + logic
│   └── CSXS/
│       └── manifest.xml     ← panel registration (Node.js enabled)
├── PS notify panel.jsx      ← script that pushes updates during a Batch
└── README.md
```

---

## How it works

- The **panel** (`extension/`) polls Photoshop every 500 ms via `evalScript`
  and shows the count + file list. Clicking a name switches to that document.
- During a **Batch/action**, Photoshop's main thread is saturated, so
  `evalScript` (and CSXSEvents) are blocked until the batch finishes — the
  panel would otherwise freeze on a stale count.
- To update live during a batch, the **`PS notify panel.jsx`** script is added
  as a step inside the action. When it runs it opens a **local TCP socket**
  (`127.0.0.1:45678`) and sends `COUNT|ACTIVE|name1|name2|...`.
- The panel runs in its **own process** (CEF + Node.js) and listens on that
  port, so it receives the push even while Photoshop itself is busy. This is
  why Node.js is enabled in the manifest (`--enable-nodejs --mixed-context`).

The small `evt: N` badge in the panel's top-right corner counts how many socket
pushes have been received — handy for confirming the action step is firing.

---

## Installation

### 1. Install the panel (CEP extension)

Copy the **`extension`** folder to Photoshop's user extensions folder and name
it `com.openfiles.counter`:

**Windows**
```
%APPDATA%\Adobe\CEP\extensions\com.openfiles.counter\
```
i.e. `C:\Users\<you>\AppData\Roaming\Adobe\CEP\extensions\com.openfiles.counter\`

**macOS**
```
~/Library/Application Support/Adobe/CEP/extensions/com.openfiles.counter/
```

The final layout must be:
```
com.openfiles.counter/
├── index.html
└── CSXS/manifest.xml
```

### 2. Allow unsigned extensions

This panel is unsigned, so enable debug mode once.

**Windows** — run in PowerShell (covers the common CEP versions):
```powershell
'CSXS.11','CSXS.12','CSXS.13' | ForEach-Object {
  New-Item -Path "HKCU:\SOFTWARE\Adobe\$_" -Force | Out-Null
  Set-ItemProperty -Path "HKCU:\SOFTWARE\Adobe\$_" -Name PlayerDebugMode -Value 1 -Type String
}
```

**macOS** — run in Terminal:
```bash
defaults write com.adobe.CSXS.11 PlayerDebugMode 1
defaults write com.adobe.CSXS.12 PlayerDebugMode 1
defaults write com.adobe.CSXS.13 PlayerDebugMode 1
```

### 3. Install the notify script

Copy **`PS notify panel.jsx`** into Photoshop's Scripts folder (needs admin on
Windows — click *Continue* when prompted):

**Windows**
```
C:\Program Files\Adobe\Adobe Photoshop 2026\Presets\Scripts\
```
**macOS**
```
/Applications/Adobe Photoshop 2026/Presets/Scripts/
```

### 4. Restart Photoshop

Open the panel from **Window ▸ Extensions (Legacy) ▸ Open Files Counter**.

---

## Usage

### Everyday
Just keep the panel open (dock it next to Layers). It shows the live count and
list. Click any filename to jump to that document.

### Live updates during a Batch
To make the count/list update *while an action closes files*:

1. **Window ▸ Actions**, open your action.
2. Select the step **before** the final *Close* step.
3. Press **Record**.
4. Run **File ▸ Scripts ▸ PS notify panel** once.
5. Press **Stop**. Drag the new step so it sits right after *Close* (so it
   reports the count *after* each file is closed).
6. Run **File ▸ Automate ▸ Batch** as usual — the panel now updates per file.

Confirm it's working: the `evt:` badge in the panel increments on each file.

---

## Troubleshooting

| Symptom | Cause / fix |
|---|---|
| Panel not in the menu | Extension folder misnamed, or debug mode not set (step 2). Must be `com.openfiles.counter` with `index.html` + `CSXS/manifest.xml` inside. |
| `evt: no node` / `no net` | Node.js not enabled — check `manifest.xml` has the `<CEFCommandLine>` block with `--enable-nodejs`. Restart PS. |
| `evt: port busy` | Another process holds port 45678. Change the port in **both** `index.html` (`PORT`) and `PS notify panel.jsx` (the `127.0.0.1:PORT` string). |
| Count updates during batch but list doesn't | You have an older `PS notify panel.jsx` that only sends the count. Use the version here (sends `COUNT|ACTIVE|names…`). |
| Nothing updates during batch | The script step isn't in the action, or `File ▸ Scripts ▸ PS notify panel` doesn't appear — re-do step 3 and re-record. |

---

## Notes / limitations

- Windows-focused. The socket push mechanism is cross-platform in principle,
  but paths and debug-mode commands above differ per OS (see each step).
- CEP is a legacy Adobe technology. A future rewrite as a **UXP** plugin would
  allow the panel's own button to be recorded directly into an action (via
  `recordAction()`), removing the need for the separate `.jsx` step. That is a
  larger rewrite and not done here.
- Filenames containing a `|` character would break the pipe-delimited payload
  (extremely rare in practice).
