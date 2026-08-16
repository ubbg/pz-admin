---
id: SANDBOX-001
title: SandboxVars bearbeiten — was RCON nicht kann
status: draft
version: 0.1.0
created: 2026-08-16
updated: 2026-08-16
owner: mail@urb.cc
stakeholders:
  - Serveradministrator:in
  - Betrieb (Docker/SFTP)
  - KI-Coding-Assistent
feature: SANDBOX-001.feature
guardrails:
  - PG-16
  - PG-17
  - PG-18
  - PG-19
  - PG-20
---

# SANDBOX-001 — SandboxVars bearbeiten

## User Story

> **Als** Administrator
> **möchte ich** die Sandbox-Einstellungen meines Servers aus derselben Anwendung
> heraus ändern — Tierhaltung, Keller, Loot-Respawn, Zombie-Verhalten —
> **damit** ich für die Hälfte der Stellschrauben von Build 42 nicht in einer Lua-Datei
> auf einem entfernten Rechner editieren muss,
> **ohne** dass ein Tippfehler mir die Datei zerschießt oder ich glaube, eine Änderung
> hätte sofort gewirkt.

## Business Context

`servertest_SandboxVars.lua` enthält in Build 42 **273 Variablen**, inklusive der
Untertabellen `Basement`, `Map`, `ZombieLore`, `ZombieConfig` und `MultiplierConfig`.
pz-admin fasst davon **keine einzige** an — im gesamten Repo kommt das Wort „sandbox"
nicht vor.

Dort steckt der halbe B42-Umfang: `Husbandry`, `Butchering`, `Basement`, `Blacksmith`,
`Pottery`, `Glassmaking`, `FlintKnapping`, `Carving`, `MetalWelding`, `Masonry`, 14×
`Animal*`, `EnableSnowOnGround`, die `Rally*`-Gruppe — und die vier Loot-Optionen, die
Build 42 aus der `servertest.ini` hierher verschoben hat
(`HoursForLootRespawn`, `MaxItemsForLootRespawn`, `ConstructionPreventsLootRespawn`,
`MinutesPerPage`).

**Der Haken:** Es gibt dafür keinen RCON-Befehl. `changeoption` bedient ausschließlich
die `.ini`. Diese Spec verlässt damit als einzige das bisherige Architekturprinzip
„alles über eine RCON-Verbindung" — und das ist der Grund, warum sie zuletzt kommt und
warum ihre Leitplanken vor allem von Vorsicht handeln.

## Technical Requirements

### Ein zweiter Zugangsweg, ausdrücklich konfiguriert

Der Editor braucht **Dateizugriff**. Drei Wege sind denkbar; die Anwendung rät keinen
davon (**PG-20**):

| Weg | Fall |
|---|---|
| lokaler Pfad | Server läuft auf demselben Rechner |
| SFTP | Server läuft auf einem Root-Server |
| Docker-Volume | Container-Setup, z. B. `./data:/home/steam/Zomboid` beim Danixu-Image |

Ohne konfigurierten Pfad ist der Tab **inaktiv mit Begründung** — nicht leer, nicht
mit Beispieldaten gefüllt, und schon gar nicht mit einem geratenen Standardpfad, der
auf einem fremden System eine falsche Datei trifft.

### Geschrieben wird die Datei, nicht der Server

Kein Sandbox-Wert geht über RCON hinaus (**PG-16**). Wer versucht, das über
`changeoption` zu erschleichen, bekommt vom Server eine freundliche Antwort und keine
Wirkung — der schlimmste Fehlerfall, weil er wie Erfolg aussieht.

### Vor dem Schreiben eine Sicherung

Jeder Schreibvorgang legt zuerst eine Kopie neben der Datei an
(`servertest_SandboxVars.lua.bak-<zeitstempel>`) (**PG-17**). Diese Datei ist der
Weltzustand eines laufenden Servers; ein misslungener Schreibvorgang ohne Sicherung
ist ein Datenverlust, den niemand rückgängig macht.

### Struktur erhalten

Die Lua-Datei wird **nicht** neu erzeugt, sondern an Ort und Stelle geändert:
Kommentare, Reihenfolge, Einrückung und die Untertabellen bleiben, wie sie sind
(**PG-18**). Die Datei ist von Hand lesbar und wird auch von Hand gelesen — eine
umsortierte Fassung macht jeden künftigen Diff unbrauchbar.

Praktisch heißt das: zeilenweise ersetzen, was sich geändert hat, statt die Struktur
zu parsen und zurückzuschreiben.

### Sagen, dass es erst nach dem Neustart gilt

Änderungen an den SandboxVars greifen **nicht** zur Laufzeit. Es gibt kein
`reloadoptions`-Gegenstück (**PG-19**). Die Oberfläche sagt das an der Stelle, an der
gespeichert wird — nicht in einer Fußnote.

## Gate-Contract

| ID | Zusage | Prüfung |
|---|---|---|
| PG-16 | Kein Sandbox-Wert wird über RCON gesendet | `grep -n 'changeoption' rcon*.go` trifft keine Sandbox-Variable; Log eines Speichervorgangs enthält kein `changeoption` |
| PG-17 | Vor jedem Schreiben existiert eine Sicherung | Wert ändern, speichern: `*.bak-*` liegt neben der Datei und entspricht dem Stand davor |
| PG-18 | Kommentare und Reihenfolge überleben | `diff` vor/nach dem Speichern zeigt ausschließlich die geänderten Wertzeilen |
| PG-19 | Der Hinweis auf den Neustart steht am Speichern-Knopf | Sandbox-Tab öffnen: Hinweis ist ohne Scrollen und ohne Aufklappen sichtbar |
| PG-20 | Ohne konfigurierten Zugang wird nichts geraten | Anwendung ohne Sandbox-Pfad starten: Tab inaktiv mit Begründung, keine Datei angefasst |

## Data Model / Registry

Neu in `Config` (als Zeiger, wie alle Felder dort):

* `SandboxAccessMode` — `nil` (nicht konfiguriert), `"local"`, `"sftp"`, `"docker"`
* `SandboxPath` — Pfad zur `*_SandboxVars.lua` bzw. zum Serververzeichnis
* SFTP-Zugangsdaten, falls dieser Weg gewählt wird — Passwort verschlüsselt wie
  `Credentials.Password`, niemals im Log

Die 273 Variablen selbst werden **nicht** im Code deklariert. Der Editor liest, was in
der Datei steht — dieselbe Haltung wie OPTIONS-001 für die `.ini`, und aus demselben
Grund. Eine Darstellungstabelle für Beschriftungen und Wertebereiche darf danebenstehen,
entscheidet aber nichts über Existenz.

## Dependencies

* SFTP braucht eine Bibliothek (`pkg/sftp` + `golang.org/x/crypto/ssh`) — die erste
  neue Laufzeitabhängigkeit des Projekts seit Langem. Wenn der Umfang zu groß wird:
  Erst nur den lokalen Pfad und das Docker-Volume, SFTP später.
* Der Danixu-Container legt die Datei unter `<volume>/Server/` ab; der Servername
  (Standard `servertest`) steckt im Dateinamen.
* Kein Zusammenhang mit der RCON-Verbindung: Der Editor funktioniert auch bei
  getrennter Verbindung — und ist damit der erste Teil der Anwendung, für den das gilt.

## Definition of Done

- [ ] Zugangsweg konfigurierbar (lokal, Docker-Volume, optional SFTP)
- [ ] Lesen und Anzeigen aller Variablen inklusive Untertabellen
- [ ] Strukturerhaltendes Schreiben mit vorheriger Sicherung
- [ ] Hinweis auf die Neustart-Pflicht am Speichern-Knopf
- [ ] Tab inaktiv mit Begründung, solange nichts konfiguriert ist
- [ ] `../B42-Luecken.md` Abschnitt 2 abgehakt

## Anti-Patterns

* **Einen Standardpfad raten.** Auf einem fremden System trifft er die falsche Datei —
  und die falsche Datei ist hier ein Weltzustand.
* **Die Datei neu erzeugen statt zu ändern.** Kommentare und Reihenfolge sind die
  einzige Dokumentation, die diese Datei hat.
* **Ohne Sicherung schreiben.** Ein abgebrochener Schreibvorgang auf einer 273-Zeilen-
  Konfiguration kostet einen Serverabend.
* **Sandbox-Werte über `changeoption` versuchen.** Der Server antwortet freundlich und
  tut nichts.
* **Die 273 Variablen im Code deklarieren.** Derselbe Fehler wie beim `PzOptions`-Struct,
  nur sechsmal so groß.
* **Beim Speichern so tun, als wirke es sofort.** Der Administrator sucht sonst den
  Fehler bei sich.
