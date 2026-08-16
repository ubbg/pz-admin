# `specs/` — Spec-Driven Development für pz-admin

Dieser Ordner enthält **ausführbare Spezifikationen**: Jede abgrenzbare Fähigkeit wird
zuerst als Spec beschrieben (maschinenlesbare Metadaten + Gherkin-Szenarien), bevor sie
implementiert wird. Die Spec ist die *Single Source of Truth* für „was gebaut wird und
woran man fertig erkennt".

Abgrenzung zu den anderen Dokumenten:

| Dokument | beantwortet |
|---|---|
| `../B42-Luecken.md` | **Was fehlt** — die Bestandsaufnahme gegen Build 42, normativ für den Ist-Zustand |
| `specs/` | **Woran man fertig erkennt** — abgrenzbare Fähigkeiten mit prüfbaren Leitplanken |
| `.claude/rules/` | **Wie man in einem Pfad schreibt** — Konventionen und Invarianten je Dateibereich |
| `CLAUDE.md` | **Wo was liegt** — Einstieg, Struktur, Befehle |

Widersprechen sich Spec und Lückenliste, gewinnt die **Lückenliste** — sie ist am
Wiki-Stand geprüft, die Spec ist eine Absicht.

## Aufbau einer Spec

Pro Spec ein Unterordner mit zwei Artefakten:

```
specs/<name>/
  <ID>.spec.md     # Markdown-Spec mit YAML-Frontmatter + strukturierten Abschnitten
  <ID>.feature     # Gherkin-Szenarien (Angenommen/Wenn/Dann), je Leitplanke ≥ 1 Szenario
```

### `*.spec.md`

Beginnt mit **YAML-Frontmatter** (maschinenlesbar):

```yaml
---
id: OPTIONS-001
title: …
status: draft | active | superseded
version: 0.1.0
created: 2026-08-16
updated: 2026-08-16
owner: mail@urb.cc
feature: OPTIONS-001.feature
guardrails: [PG-01, PG-02, …]   # alle in dieser Spec definierten, prüfbaren IDs
---
```

… gefolgt von Markdown-Abschnitten: **User Story**, **Business Context**,
**Technical Requirements**, **Gate-Contract**, **Data Model / Registry**,
**Dependencies**, **Definition of Done**, **Anti-Patterns**.

### `*.feature`

Gherkin in deutscher Sprache (`# language: de`). Jedes Szenario ist mit der
zugehörigen Leitplanken-ID getaggt (`@PG-07`), sodass Spec ↔ Szenario ↔ Prüfung über
die ID verbunden sind.

Die `.feature`-Dateien werden **nicht ausgeführt**. Sie sind Verhaltensdokumentation.

## Traceability-Konvention (IDs)

* **Spec-ID:** `OPTIONS-001`, `CMD-001`, … (Großbuchstaben-Präfix + Nummer).
* **Leitplanken:** `PG-01 … PG-NN` (**P**Z-Admin-**G**uardrail), fortlaufend über alle
  Specs hinweg — eine Nummer wird nie zweimal vergeben und nach einer aufgegebenen
  Spec nicht recycelt.
* Eine `PG`-ID erscheint in **(1)** der Spec, **(2)** als `@PG-xx`-Tag im `.feature`
  und **(3)** — wo prüfbar — in der Spalte „Prüfung" des Gate-Contracts.
  `grep -rn "PG-12" .` zeigt alle Berührungspunkte.

### Es gibt keinen Traceability-Test

Anders als in Schwesterprojekten hat pz-admin **keine Testsuite** und kein
Gesamt-Gate — CI baut nur (`.github/workflows/build.yml`). Der Abgleich läuft
deshalb von Hand:

```bash
# Alle Specs prüfen: Frontmatter-IDs gegen Feature-Tags
cd specs
for f in */*.spec.md; do
  d="${f%.spec.md}"
  spec=$(awk '/^---$/{n++; next} n==1' "$f" | grep -oE 'PG-[0-9]+' | sort -u)
  feat=$(grep -oE '@PG-[0-9]+' "$d.feature" | tr -d '@' | sort -u)
  [ "$spec" = "$feat" ] && echo "OK   $d" || { echo "DIFF $d"; diff <(echo "$spec") <(echo "$feat"); }
done

# doppelt vergebene IDs über alle Specs finden (Ausgabe muss leer sein)
grep -rhoE '^  - PG-[0-9]+' --include='*.spec.md' . | grep -oE 'PG-[0-9]+' | sort | uniq -d
```

Weil keine Testdatei die Zusage einlöst, trägt jeder Gate-Contract **ausführbare
Kommandos oder benannte Handgriffe** (z. B. „gegen einen B42-Server verbinden und X
beobachten") — niemals den Verweis auf einen Test, den es nicht gibt.

## Sprache

Prosa und Gherkin sind **deutsch**, Bezeichner, Kommandos, Options- und Befehlsnamen
bleiben **englisch** — dieselbe Aufteilung wie im Code (siehe `CLAUDE.md`).

## Aktuelle Specs

| ID | Titel | Leitplanken | Status |
|----|-------|-------------|--------|
| [OPTIONS-001](options/OPTIONS-001.spec.md) | Serveroptionen ohne Build-Bindung | PG-01 … PG-06 | draft |
| [CMD-001](commands/CMD-001.spec.md) | Die Befehlsfläche von Build 42 | PG-07 … PG-11 | draft |
| [ITEMS-001](catalog/ITEMS-001.spec.md) | Der Item-Katalog aus Build 42 | PG-12 … PG-15 | draft |
| [SANDBOX-001](sandbox/SANDBOX-001.spec.md) | SandboxVars bearbeiten — was RCON nicht kann | PG-16 … PG-20 | draft |

Alle vier sind `draft`: Nichts davon ist implementiert. Die Reihenfolge oben ist
zugleich die empfohlene Umsetzungsreihenfolge — OPTIONS-001 räumt den größten Teil der
Lückenliste in einem Zug ab.
