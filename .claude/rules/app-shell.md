---
paths:
  - "main.go"
  - "app.go"
  - "app_paths.go"
  - "config.go"
  - "update.go"
  - "logger.go"
  - "dialog.go"
  - "wails.json"
  - "build/**"
  - ".github/workflows/*.yml"
---

# Anwendungsrahmen — Wails, Konfiguration, Selbstaktualisierung, Build

## Die Bindung ist der Vertrag

`main.go` bindet genau ein Objekt: `Bind: []interface{}{app}`. Daraus folgt für jede
neue Funktion, die das Frontend aufrufen soll:

* Sie ist eine **exportierte Methode auf `*App`**. Kleingeschriebene Methoden und
  freie Funktionen bleiben intern — so trennt das Projekt heute `players_update()`
  von `Players()`.
* Ihre Typen brauchen `json:"…"`-Tags; ohne Tag erscheinen sie in
  `frontend/src/wailsjs/go/models.ts` großgeschrieben.
* Nach der Änderung `wails dev` laufen lassen, damit die Bindings neu entstehen.
  `frontend/src/wailsjs/` wird **nie** von Hand bearbeitet.

Lebenszyklus-Haken liegen in `app.go`: `startup`, `domReady`, `beforeClose`,
`shutdown`, `onSecondInstanceLaunch`. Der Single-Instance-Lock hat eine feste UUID in
`main.go` — die bleibt, sie identifiziert die Anwendung auf dem System.

## Konfiguration

`Config` (`config.go`) besteht durchgehend aus **Zeigern** (`*bool`, `*int`,
`*string`). Das ist Absicht: `nil` heißt „nicht gesetzt" und ist von `false`/`0`
unterscheidbar, sodass Standardwerte nachrüstbar bleiben, ohne bestehende
Konfigurationen zu überschreiben. Ein neues Feld also **immer** als Zeiger anlegen und
im Default-Pfad füllen.

Kommentare hinter den Feldern nennen den erlaubten Wertebereich
(`// 0 = Normal, 1 = Maximized, …`). Diese Zuordnung existiert doppelt — im Kommentar
und in der `switch`-Anweisung in `main.go`. Wer sie ändert, ändert beide.

Pfade (Konfiguration, Logs, temporäre Dateien) kommen ausschließlich aus
`app_paths.go`. Kein `os.UserConfigDir()` an zweiter Stelle, kein hartkodiertes
`~/.config`.

## Selbstaktualisierung — eine Konstante, kein verstreutes Wissen

Die Herkunft der Aktualisierung steht an **einer** Stelle, oben in `update.go`:

```go
const (
	updateRepoOwner = ""   // leer = keine eigene Release-Quelle
	updateRepoName  = "pz-admin"
)
```

Solange `updateRepoOwner` leer ist, prüft die Anwendung **gar nicht**: `CheckForUpdate`
kehrt vor dem Netzzugriff zurück, `Update` verweigert das Einspielen, und
`CheckForUpdates` steht in der Konfiguration standardmäßig auf `false`
(`GetDefaultConfig`, zusätzlich erzwungen in `config_init`). Das ist Absicht — ein
Upstream-Release brächte alle Build-41-Annahmen zurück, die dieser Fork gerade
beseitigt.

Wer eigene Releases veröffentlicht, trägt hier seinen GitHub-Benutzer ein; Prüfung,
Download-URL und die Anzeige in den Einstellungen folgen daraus. `GetUpdateSource()`
liefert dieselbe Angabe ans Frontend (`UpdateSourceSetting`), damit sichtbar ist,
welcher Quelle die Anwendung folgt. `go.mod` trägt weiterhin den Upstream-Modulpfad —
das ist der Importpfad, nicht die Release-Quelle.

## Build und Release

```bash
wails build                     # nach build/bin/
wails build -tags webkit2_41    # Linux, so baut auch CI
```

CI (`.github/workflows/build.yml`) baut nur — es gibt **keinen** Lint- und keinen
Testschritt:

* Go **1.23.4**, Node **22.x**, Wails-CLI **v2.11.0** (alle drei gepinnt).
* Matrix: Linux amd64 (zusätzlich `.deb` via `dpkg-deb` und `.rpm` via `alien`),
  Windows amd64 und arm64 (NSIS-Installer).
* Auf `main` entsteht ein **Draft-Release** mit `gh release create`.

Die Version steht an genau **einer** Stelle: `wails.json` → `info.productVersion`.
Der Job `extract-version` liest sie per `jq` und baut daraus alle Dateinamen und den
Tag. Ein Release ohne Bump dort überschreibt den vorherigen Tag.

`build/` enthält Icons, NSIS-Vorlagen und Plattformdateien; `build/bin/` ist
Ausgabeverzeichnis und in `.gitignore`.

## Protokoll

`logger.go` schreibt nach `logs/<zeitstempel>.log`, gesteuert über die
`Enable*`-Schalter der Konfiguration und begrenzt durch `MaxLogFiles`. Log-Level
laufen über `runtime.Log*` aus Wails.

**Nie ins Log:** RCON-Passwort, Serverpasswort, `Credentials`-Strukturen als Ganzes.
Beim Loggen einer Fehlermeldung aus dem Verbindungsaufbau prüfen, ob die Bibliothek
die Zugangsdaten mitschickt.
