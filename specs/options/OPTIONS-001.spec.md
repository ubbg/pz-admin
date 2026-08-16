---
id: OPTIONS-001
title: Serveroptionen ohne Build-Bindung — der Server sagt, was es gibt
status: draft
version: 0.1.0
created: 2026-08-16
updated: 2026-08-16
owner: mail@urb.cc
stakeholders:
  - Serveradministrator:in
  - KI-Coding-Assistent
feature: OPTIONS-001.feature
guardrails:
  - PG-01
  - PG-02
  - PG-03
  - PG-04
  - PG-05
  - PG-06
---

# OPTIONS-001 — Serveroptionen ohne Build-Bindung

## User Story

> **Als** Administrator eines Project-Zomboid-Servers
> **möchte ich** im Optionen-Tab genau die Optionen sehen und ändern können, die *mein*
> Server hat,
> **damit** ich nach einem Server-Update nicht auf ein neues Release dieser Anwendung
> warten muss,
> **ohne** dass mir Optionen angeboten werden, die es auf meinem Server gar nicht gibt.

## Business Context

Die Anwendung kennt Serveroptionen heute als **statisches Go-Struct**
(`PzOptions`, `rcon_options.go:15`, 135 Felder). `parseOptions()` ordnet die Zeilen aus
`showoptions` per Reflection Feldern zu und **verwirft**, was es nicht kennt
(`LogDebugf("Unknown option: %s")`). `diffOptions()` iteriert über dieselben
Struct-Felder.

Damit ist jede Option, die das Struct nicht kennt, nicht nur unsichtbar, sondern
**prinzipiell nicht änderbar** — und jede Option, die es nur im Struct gibt, erzeugt
beim Speichern ein fehlschlagendes `changeoption`.

Der Ist-Zustand gegen Build 42 (`../B42-Luecken.md`): 94 von 144 Optionen abgedeckt,
**50 fehlen**, **41 Felder sind tot**. Die Anti-Cheat- und Discord-Tabs der Oberfläche
sind auf einem B42-Server vollständig funktionslos.

Man könnte die 50 nachtragen und die 41 löschen. Das repariert diesen einen Build und
stellt dasselbe Problem beim nächsten wieder her — es ist bereits das zweite Mal, dass
die Liste auseinanderläuft. Diese Spec dreht die Richtung um: **Nicht die Anwendung
sagt, welche Optionen es gibt, sondern der Server.**

Das Struct verschwindet damit als *Existenzentscheidung*. Was bleibt, ist eine
**Darstellungstabelle** (`frontend/src/assets/options.ts`): Kategorie, Beschriftung,
Wertebereich, Auswahlwerte, Abhängigkeiten. Fehlt zu einer Option ein Eintrag, ist sie
trotzdem da — nur schlichter dargestellt.

## Technical Requirements

### Ein generischer Rundlauf

```
showoptions ──► []Option{Name, Value, Kind} ──► EventsEmit("update-options") ──► UI
UI ──► UpdatePzOptions([]OptionPair) ──► changeoption je geänderter Option ──► reloadoptions
```

* `showoptions` liefert Zeilen der Form `* Name=Wert`. Der Parser übernimmt **jede**
  davon, ohne Vorwissen (**PG-01**).
* Der Typ wird aus dem Wert abgeleitet: `true`/`false` → Boolean, reine Ziffern →
  Integer, Ziffern mit Punkt → Double, sonst String. Was sich keiner Klasse zuordnen
  lässt, ist ein String und wird als Textfeld bearbeitet (**PG-05**).
* Geschrieben wird nur, was sich gegenüber dem zuletzt gelesenen Serverstand geändert
  hat, und nur für Namen, die der Server selbst gemeldet hat (**PG-03**). Damit
  verschwindet die Klasse „wir schicken eine Option, die es nicht gibt" vollständig.

### `options.ts` wird zur reinen Darstellungstabelle

Der heutige `Option`-Typ (`FieldName`, `Type`, `Default`, `Range`, `DisabledValue`,
`Requirements`, `Choices`) bleibt inhaltlich bestehen, verliert aber seine
Torwächterrolle: Ein Eintrag beschreibt, **wie** eine Option aussieht, nicht **ob** es
sie gibt (**PG-02**). Optionen ohne Eintrag erscheinen in einer Kategorie
„Weitere Optionen", alphabetisch, mit ihrem Rohnamen als Beschriftung.

Das ist zugleich der Melder: Was dort auftaucht, gehört bei nächster Gelegenheit
sauber einsortiert.

### Anti-Cheat ist vierwertig, nicht boolesch

Build 42 ersetzt die 24 Booleans (`AntiCheatProtectionType1…24`) plus 7
Schwellenwert-Multiplikatoren durch **10 benannte Optionen** mit vier Werten:
`1 = Ban`, `2 = Kick`, `3 = Log`, `4 = Disable` (**PG-04**).

Ein Schalter kann das nicht abbilden. Die Typwahl `Choice` mit diesen vier Werten ist
Teil dieser Spec, weil eine automatische Typerkennung hier zwangsläufig „Integer"
raten und ein Zahlenfeld anbieten würde — richtig, aber unbrauchbar.

### Fehler nennen ihren Namen

`applyOptions` meldet heute nur eine Anzahl (`failed_to_update_n_options`). Beim
Übergang auf einen anderen Build ist genau die Frage interessant, **welche** Option
abgelehnt wurde (**PG-06**).

## Gate-Contract

| ID | Zusage | Prüfung |
|---|---|---|
| PG-01 | Eine Option, die die Anwendung nicht kennt, erscheint in der Oberfläche und bleibt änderbar | Gegen einen B42-Server verbinden; `War`, `AnnounceAnimalDeath` und `AntiCheatSafety` sind im Optionen-Tab sichtbar und speicherbar |
| PG-02 | Die Existenz einer Option ist nirgends im Code deklariert | `grep -c 'FieldName:' frontend/src/assets/options.ts` darf schrumpfen, ohne dass Optionen verschwinden; kein `PzOptions`-Struct mehr in `rcon_options.go` |
| PG-03 | `changeoption` wird nur für Namen gesendet, die `showoptions` gemeldet hat | Log auf `changeoption` filtern; jeder gesendete Name kommt in der vorherigen `showoptions`-Antwort vor |
| PG-04 | Die zehn B42-Anti-Cheat-Optionen sind vierwertig bedienbar | Optionen-Tab öffnen: `AntiCheatSafety` bietet Ban/Kick/Log/Disable, keinen Schalter |
| PG-05 | Ein Wert unbekannten Typs geht beim Speichern nicht verloren | Eine String-Option mit Sonderzeichen (`ChatStreams`, `ClientCommandFilter`) ändern, speichern, `showoptions` erneut lesen — Wert identisch |
| PG-06 | Eine abgelehnte Option wird namentlich gemeldet | Auf einem B42-Server eine tote Option erzwingen; die Benachrichtigung nennt ihren Namen |

## Data Model / Registry

Kein persistentes Modell. Der Serverstand lebt im Prozess (heute `pzOptions`,
künftig eine Liste oder Map) und wird über `lastOptionsHash` gegen unnötige Events
verglichen — dieser Kurzschluss bleibt erhalten, er spart bei jedem Poll ein
vollständiges Rerender.

Die Darstellungstabelle `options.ts` bleibt die einzige Stelle mit Wissen über
Kategorien, Bereiche und Auswahlwerte. Sie ist Daten, kein Code — Änderungen dort
brauchen keinen Go-Build.

## Dependencies

* Keine neue Bibliothek. Der Parser wird einfacher, nicht komplizierter: `reflect`
  entfällt.
* `frontend/src/wailsjs/go/models.ts` ändert sich (`PzOptions` verschwindet als Typ) —
  betroffen sind `rcon-provider.tsx` und `Options.tsx`.
* Die Sprachkataloge behalten ihre Optionsschlüssel; fehlende Schlüssel führen künftig
  zum Rohnamen statt zum leeren Feld.

## Definition of Done

- [ ] `showoptions` wird generisch geparst, `PzOptions` als Struct entfällt
- [ ] Unbekannte Optionen erscheinen unter „Weitere Optionen" und sind änderbar
- [ ] `changeoption` nur für vom Server gemeldete Namen
- [ ] Die zehn B42-Anti-Cheat-Optionen als vierwertige Auswahl in `options.ts`
- [ ] Tote B41-Einträge aus `options.ts` entfernt, B42-Einträge ergänzt
- [ ] Fehlermeldung nennt die betroffene Option
- [ ] `../B42-Luecken.md` Abschnitt 1 abgehakt, `.claude/rules/options.md` nachgezogen

## Anti-Patterns

* **Die 50 fehlenden Optionen ins Struct nachtragen.** Repariert Build 42 und
  reproduziert das Problem für Build 43.
* **Unbekannte Optionen ausblenden, statt sie roh zu zeigen.** Genau das macht die
  Anwendung heute, und genau deshalb ist die Lücke drei Jahre lang nicht aufgefallen.
* **Den Typ aus `options.ts` als Wahrheit nehmen.** Der Server liefert den Wert; die
  Tabelle sagt nur, wie er aussehen soll. Bei Widerspruch gewinnt der Server.
* **Alle Optionen bei jedem Speichern senden.** Ein `changeoption` je Feld über 144
  Felder ist ein spürbarer Stau auf der RCON-Verbindung — und schreibt Werte neu, die
  niemand angefasst hat.
* **Das RCON-Passwort im Optionen-Tab änderbar machen, ohne die Folge abzufangen.**
  Der nächste Befehl läuft dann ins Leere.
