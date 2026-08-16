---
paths:
  - "rcon.go"
---

# RCON-Schicht — eine Verbindung, ein Mutex, und der Server hat recht

`rcon.go` ist die einzige Stelle, die mit dem Project-Zomboid-Server spricht. Alles,
was die Oberfläche anzeigt, ist entweder direkt eine Serverantwort oder daraus
abgeleitet. Es gibt keine zweite Quelle, keinen Cache mit eigenem Kopf und keine
Datenbank.

## Die Invariante

> Es gibt genau **eine** Verbindung (`conn`), und jeder Zugriff darauf hält
> `connMutex`.

Konkret:

* `ConnectRcon` schließt eine bestehende Verbindung, bevor es eine neue aufbaut —
  nicht daneben.
* `watchConnection` läuft als **eine** Goroutine, bewacht durch `isWatching` und
  beendet über den Kanal `stopWatching`. Ein zweiter Watcher pollt denselben Server
  doppelt und meldet Verbindungsabbrüche doppelt.
* `RCONCommand.execute()` sperrt selbst. Wer `execute()` aus einer bereits sperrenden
  Funktion aufruft, verklemmt die Anwendung — das ist der wahrscheinlichste Deadlock
  in diesem Projekt.

## Der Server ist die Wahrheit

Nach jedem schreibenden Befehl wird der Zustand **neu geholt**, statt lokal
fortgeschrieben zu werden: `players_update()` liest `players`, `options_update()`
liest `showoptions`. Beide vergleichen vorher (`lastOptionsHash`) und schweigen, wenn
sich nichts geändert hat.

Lokal geführt wird ausschließlich, was RCON gar nicht hergibt: Bann-Status,
Zugriffsstufe und Godmode je Spieler. Diese Felder überleben in der Konfiguration und
werden beim nächsten `players`-Lauf wieder an die Online-Liste geheftet. Wer hier ein
weiteres Feld erfindet, muss sagen können, woher es nach einem Neustart kommt.

## Erfolg wird an der Antwort geprüft, nicht am Fehlerwert

Ein PZ-Server antwortet auf einen unbekannten oder fehlgeschlagenen Befehl mit
freundlichem Text — und ohne Fehler. `err == nil` heißt also **nicht**, dass der
Befehl gewirkt hat.

Deshalb trägt jedes Kommando seine Prüfung mit sich:

```go
command := RCONCommand{
    CommandTemplate: "reloadoptions",
    SuccessCheck: func(name string, response string) bool {
        return response == "Options reloaded"
    },
}
```

* `SuccessCheck` — was der Server sagt, wenn es geklappt hat.
* `ErrorCheck` — die bekannten Fehlantworten, damit „unbekannte Antwort" nicht
  stillschweigend als Erfolg durchgeht.
* `Notifications` — fünf Fälle (`AllSuccess`, `AllFail`, `Partial`, `SingleSuccess`,
  `SingleFail`), weil ein Befehl über mehrere Spieler teilweise gelingen kann.
* `UpdateFunc` / `EmitUpdatePlayers` — den lokalen Spielerzustand nachziehen und die
  Oberfläche benachrichtigen.

Diese Vergleichs-Strings sind **Build-42-Annahmen**. Ändert der Server seine
Antworttexte, bricht nicht der Aufruf, sondern die Erfolgsmeldung — der schlimmere
Fehler, weil er wie Erfolg aussieht. Neue Vergleiche deshalb eng an der Wiki-Formulierung
halten und die Build-Version im Kommentar nennen.

## Spielerbezogene Befehle laufen über `RCONCommand`

`conn.Execute` direkt ist drei Stellen vorbehalten: der Watcher, `players_update`,
und `SendRconCommand` (das Terminal, wo die Nutzerin ausdrücklich Rohzugriff will).
Alles andere geht durch `RCONCommand` — sonst entstehen Kommandos ohne
Erfolgsprüfung, ohne Fortschritt und ohne Benachrichtigung, und genau daran erkennt
man im Code die nachträglich angeflanschten Befehle.

Argumente kommen als `RCONCommandParam` (`Name`, `Key`, `Value`, `Mandatory`) in die
Vorlage, statt per `fmt.Sprintf` zusammengeklebt zu werden. Spielernamen gehören in
`"…"` — ein Name mit Leerzeichen ist der Regelfall, nicht die Ausnahme.

## Fortschritt und Benachrichtigungen

* `runtime.EventsEmit(app.ctx, "setProgress", …)` — 0 am Ende, immer per `defer`
  zurücksetzen, sonst bleibt der Balken stehen.
* `app.SendNotification(Notification{Title: "rcon.xxx", …})` — `Title` ist ein
  **i18n-Schlüssel**, keine Meldung. Der Schlüssel muss in allen fünf Katalogen
  existieren, siehe `.claude/rules/i18n.md`.

## Was hier nicht hineingehört

* **Kein Dateizugriff auf den Server.** SandboxVars, `servertest.ini` und Mods-Ordner
  liegen außerhalb dessen, was RCON kann; ein Editor dafür ist eine eigene Schicht
  (`specs/sandbox/SANDBOX-001.spec.md`), nicht eine Erweiterung dieser Datei.
* **Kein Passwort im Log.** `Credentials.Password` ist verschlüsselt abgelegt und
  taucht weder in `runtime.Log*` noch in einer Benachrichtigung auf.
* **Keine Fachlogik der Oberfläche.** Welche Spalten eine Tabelle zeigt, entscheidet
  das Frontend; `rcon.go` liefert Daten.

## Fehlende B42-Befehle

19 Admin-Befehle aus Build 42 haben hier bis heute keinen Einstiegspunkt
(`addkey`, `addsteamid`, `banip`, `worldgen`, `removezombies`, `setpassword`, …).
Die Liste steht vollständig in `../B42-Luecken.md`, das Arbeitspaket in
`specs/commands/CMD-001.spec.md`. Wer einen davon nachrüstet, folgt dem
`RCONCommand`-Muster und trägt ihn zugleich in `frontend/src/components/Terminal.tsx`
und die Sprachkataloge ein.
