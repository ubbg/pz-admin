# pz-admin — RCON-Adminoberfläche für Project-Zomboid-Server (Wails v2 + Go 1.23 + React 19)

**Aktueller Stand (v1.2.0):** Eine Desktop-Anwendung, die sich über **RCON** mit einem
Project-Zomboid-Dedicated-Server verbindet und ihn administriert: Spielerliste mit
Ban/Kick/Whitelist/Zugriffsstufen, Items, Fahrzeuge und XP, Wetter und Horden,
Servernachrichten, ein Optionen-Editor über `showoptions`/`changeoption` und ein
Terminal für alles Übrige. Der gesamte Zustand kommt aus genau einer Quelle — der
RCON-Verbindung —, es gibt **keine** Datenbank und **keinen** Dateizugriff auf den
Server.

**Der wichtigste Satz dieses Dokuments:** Die Anwendung ist durchgehend gegen
**Build 41** gebaut, der Server läuft inzwischen auf **Build 42**. Was genau fehlt,
steht vollständig in `../B42-Luecken.md` (50 fehlende Serveroptionen, 41 tote Felder,
273 nicht abgedeckte SandboxVars, 19 fehlende Admin-Befehle, B41-Item-Katalog). Die
Specs unter `specs/` schneiden daraus die Arbeitspakete. Wer hier etwas ändert, prüft
zuerst, ob die Änderung eine dieser Lücken berührt.

## Projektstruktur (Abhängigkeits-Reihenfolge)

Go-Dateien liegen flach in der Wurzel (Package `main`), das Frontend unter `frontend/`.
Die Grenze zwischen beiden ist die **Wails-Bindung** und der wichtigste Schnitt des
Projekts: Jede exportierte Methode auf `*App` wird zu einer TypeScript-Funktion.

* `main.go` — `wails.Run` mit allen Fensteroptionen, `Bind: []interface{}{app}`,
  eingebettete Assets (`//go:embed all:frontend/dist`), Single-Instance-Lock.
* `app.go` — Lebenszyklus (`startup`, `domReady`, `beforeClose`, `shutdown`),
  Benachrichtigungen, `RestartApplication`.
* `rcon.go` — **das Herz.** Genau eine globale Verbindung (`conn` + `connMutex`),
  ein Watcher-Goroutine (`watchConnection`), die Spielerliste als Prozesszustand
  (`players`) und das Muster `RCONCommand` für alles, was Spieler betrifft.
  Details und Invarianten: `.claude/rules/rcon.md`.
* `rcon_options.go` — Serveroptionen: `showoptions` → `[]Option{Name, Value, Kind}` →
  Event `update-options-list` → UI; zurück über `changeoption` + `reloadoptions`.
  **Kein Optionsname steht in Go** — welche Optionen es gibt, entscheidet der Server,
  wie sie aussehen `frontend/src/assets/options.ts`. Siehe `.claude/rules/options.md`
  und `specs/options/OPTIONS-001.spec.md`.
* `config.go`, `app_paths.go` — Anwendungseinstellungen (`Config`, durchgehend
  **Zeiger**-Felder, damit „nicht gesetzt" von „false" unterscheidbar bleibt) und die
  Pfade dorthin.
* `update.go` — Selbstaktualisierung über `minio/selfupdate`. Die Release-Quelle steht
  als Konstante `updateRepoOwner` am Dateikopf; ist sie leer, prüft die Anwendung
  nicht und spielt nichts ein. Vor dem ersten eigenen Release dort den eigenen
  GitHub-Benutzer eintragen — siehe `.claude/rules/app-shell.md`.
* `logger.go`, `language.go`, `dialog.go`, `utils.go`,
  `notification_{windows,default}.go` — Querschnitt; `language.go` ermittelt die
  Systemsprache über `jeandeaual/go-locale`.
* `frontend/src/components/` — die Oberfläche. `AdminPanel.tsx` ist die Klammer,
  darunter `Players.tsx`, `Management.tsx`, `Options.tsx`, `Tools.tsx`,
  `WeatherControl.tsx`, `Terminal.tsx`, `Settings.tsx`; Formulare liegen als
  `Dialogs/*.tsx`, die shadcn-Bausteine unter `ui/`.
* `frontend/src/contexts/` — je ein Provider für RCON-Zustand, Konfiguration, Theme,
  Farbschema, Fortschritt und Storage. **Kein Redux, kein Zustand-Store.**
* `frontend/src/wailsjs/` — **generiert.** Niemals von Hand ändern; entsteht neu bei
  `wails dev` / `wails build`.
* `frontend/src/assets/` — `items.json` (~2060 IDs), `vehicles.json`, `options.ts`
  (Darstellungs-Metadaten der Serveroptionen). Die beiden JSON-Kataloge werden von
  `frontend/scripts/*.mjs` aus dem Wiki erzeugt, siehe `.claude/rules/catalogs.md`.
* `frontend/public/locales/<sprache>/common.json` — vier Sprachen: `en-US`, `tr-TR`,
  `ru-RU`, `uk-UA`.

## Globale Entwicklungsrichtlinien

* **Sprache:** Bezeichner, Kommandos, i18n-Schlüssel und Tool-Namen sind **englisch**;
  Kommentare, Commit-Messages und Prosa (auch in Specs und Regeln) sind **deutsch**.
  Vom Upstream übernommener englischer Kommentar bleibt englisch — nicht übersetzen,
  nur Neues auf Deutsch.
* **Nach außen sichtbare Strings gehören in die Sprachkataloge**, nie in den Code.
  Das gilt auch für Go: `Notification.Title` ist ein **i18n-Schlüssel**
  (`"rcon.rcon_connection_failed"`), keine Meldung. Ein neuer Schlüssel wird in
  **allen vier** Katalogen angelegt — siehe `.claude/rules/i18n.md`.
* **Eine Verbindung, ein Mutex.** Jeder RCON-Zugriff läuft über `conn` unter
  `connMutex`. Kein zweiter Dial, kein Zugriff ohne Sperre, kein `conn.Execute` aus
  dem Frontend-Pfad heraus außer über `SendRconCommand`.
* **Der Server ist die Wahrheit.** Nach jedem schreibenden Befehl wird der Zustand
  neu vom Server geholt (`players_update`, `options_update`) statt lokal
  fortgeschrieben. Lokal geführt wird nur, was RCON nicht hergibt (Bann-Status,
  Zugriffsstufen).
* **Kein Ratespiel bei Serverantworten.** Erfolg wird an der Antwort geprüft
  (`SuccessCheck`/`ErrorCheck`), nicht daran, dass `Execute` keinen Fehler lieferte.
  Ein PZ-Server antwortet auf Unsinn freundlich mit Text und Exit-Code 0.
* **Build-Annahmen sichtbar machen.** Alles, was an einer bestimmten
  Project-Zomboid-Version hängt (Optionsnamen, Befehlsnamen, Item-IDs, Antwort-Texte),
  gehört an eine benannte Stelle mit Versionsangabe — nicht verstreut in
  Vergleichs-Strings.
* **Keine Geheimnisse im Klartext.** Das RCON-Passwort liegt verschlüsselt in den
  Credentials (`Credentials.Password`, „Encrypted"); es gehört nicht ins Log, nicht in
  eine Fehlermeldung und nicht in eine Spec.

## Befehle

```bash
wails dev                      # Entwicklungsmodus mit Hot-Reload (Go + Vite)
wails build                    # Produktionsbuild nach build/bin/
wails build -tags webkit2_41   # Linux: nötig für WebKit2GTK 4.1 (so baut auch CI)
cd frontend && yarn install    # Frontend-Abhängigkeiten (yarn, nicht npm — wails.json)
cd frontend && yarn lint       # eslint, --max-warnings 0
cd frontend && yarn items      # items.json + Bilder aus dem Wiki neu erzeugen
cd frontend && yarn vehicles   # vehicles.json neu erzeugen (läuft auch vor jedem Build)
go vet ./...                   # das einzige Go-seitige Prüfwerkzeug im Projekt
gofmt -l .                     # muss leer sein
```

**Es gibt keine Testsuite und kein Gesamt-Gate.** CI (`.github/workflows/build.yml`)
baut nur: Go 1.23.4, Node 22, Wails-CLI v2.11.0, Matrix aus Linux amd64 (plus `.deb`
und `.rpm`) sowie Windows amd64/arm64 (NSIS-Installer); auf `main` entsteht daraus ein
**Draft-Release**. Die Versionsnummer kommt aus `wails.json` → `info.productVersion`
und wird von `extract-version` gelesen — sie ist damit die einzige Stelle, an der eine
Version gepflegt wird.

Weil kein Testlauf existiert, tragen die Specs ihre Prüfungen als **ausführbare
Kommandos oder benannte Handgriffe** — nicht als Testdatei, die es nicht gibt.

## Wails-Bindung — der Vertrag zwischen Go und TypeScript

Gebunden wird ausschließlich `app` (`main.go`, `Bind:`). Daraus folgt:

* Eine Funktion wird für das Frontend sichtbar, indem sie **exportierte Methode auf
  `*App`** wird — `func (app *App) DoThing(...)`. Freie Funktionen und
  kleingeschriebene Methoden bleiben unsichtbar; genau so trennt das Projekt heute
  Innenleben (`players_update`, `options_update`) von Oberfläche.
* Parameter- und Rückgabetypen erscheinen in `frontend/src/wailsjs/go/models.ts`. Ein
  Go-Struct ohne `json:"…"`-Tags erzeugt dort Feldnamen in Großschreibung — Tags sind
  Pflicht, nicht Geschmack.
* Ereignisse laufen über `runtime.EventsEmit`: `update-players`, `update-options-list`,
  `setProgress`. Ein neues Ereignis braucht einen Abnehmer in einem Context-Provider,
  sonst ist es tot.

## Referenzen & Regeln

* **B42-Lückenliste:** `../B42-Luecken.md` — die vollständige Bestandsaufnahme gegen
  pzwiki (42.20.0/42.20.2) und den Danixu-Docker-Server. Liegt bewusst außerhalb des
  Repos, weil sie den Ist-Zustand des *Forks* beschreibt, nicht den des Upstreams.
* **Ausführbare Spezifikationen:** `specs/` — Format und ID-Konvention in
  `specs/README.md`; die Leitplanken heißen `PG-xx`.
  * `specs/options/OPTIONS-001.spec.md` (PG-01..PG-06) — Serveroptionen ohne
    Build-Bindung
  * `specs/commands/CMD-001.spec.md` (PG-07..PG-11) — die B42-Befehlsfläche
  * `specs/catalog/ITEMS-001.spec.md` (PG-12..PG-15) — Item-Katalog aus B42
  * `specs/sandbox/SANDBOX-001.spec.md` (PG-16..PG-20) — SandboxVars bearbeiten
* **Pfad-gebundene Regeln** unter `.claude/rules/` laden automatisch, sobald an den
  passenden Dateien gearbeitet wird:
  * `rcon.md` — Verbindung, `RCONCommand`, Antwortprüfung, Polling (`rcon.go`).
  * `options.md` — der Dreiklang Struct ↔ `options.ts` ↔ Sprachkataloge, B42-Namen
    (`rcon_options.go`, `frontend/src/assets/options.ts`).
  * `frontend.md` — React, shadcn/ui, Provider, Dialoge, Tabellen
    (`frontend/src/**`).
  * `i18n.md` — vier Kataloge, Schlüsselkonventionen, Go-seitige Schlüssel
    (`frontend/public/locales/**`, `frontend/src/i18n.ts`, `language.go`).
  * `catalogs.md` — die Wiki-Generatoren und ihre Versionsbindung
    (`frontend/scripts/*.mjs`, `frontend/src/assets/*.json`).
  * `app-shell.md` — Wails-Optionen, Config, Selbstaktualisierung, Build und Release
    (`main.go`, `app.go`, `config.go`, `update.go`, `wails.json`, `build/**`,
    `.github/workflows/*`).

## Agent skills

Konfiguration für die Engineering-Skills aus `mattpocock/skills`. Sie lesen die
Dateien unter `docs/agents/` — das Verzeichnis richtet sich an Agenten, nicht an
Menschen.

### Issue tracker

Issues und Specs liegen als Markdown unter `.scratch/<feature>/` in diesem Repo
(nicht in GitHub-Issues — der Fork trägt die Issue-Liste des Upstreams).
Siehe `docs/agents/issue-tracker.md`.

### Triage labels

Die fünf kanonischen Rollen unter ihren Standardnamen (`needs-triage`, `needs-info`,
`ready-for-agent`, `ready-for-human`, `wontfix`). Siehe `docs/agents/triage-labels.md`.

### Domain docs

Single-Context: ein `CONTEXT.md` und ein `docs/adr/` in der Wurzel.
Siehe `docs/agents/domain.md`.
