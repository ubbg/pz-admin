---
id: ITEMS-001
title: Der Item-Katalog aus Build 42 — mit benannter Herkunft
status: draft
version: 0.1.0
created: 2026-08-16
updated: 2026-08-16
owner: mail@urb.cc
stakeholders:
  - Serveradministrator:in
  - KI-Coding-Assistent
feature: ITEMS-001.feature
guardrails:
  - PG-12
  - PG-13
  - PG-14
  - PG-15
---

# ITEMS-001 — Der Item-Katalog aus Build 42

## User Story

> **Als** Administrator
> **möchte ich** im Item-Dialog die Gegenstände finden, die es auf meinem Server
> wirklich gibt — auch die aus Build 42,
> **damit** ich nicht für jedes Schmiede- oder Tierhaltungs-Item die ID von Hand
> nachschlagen muss,
> **ohne** dass später jemand raten muss, gegen welchen Spielstand der Katalog erzeugt
> wurde.

## Business Context

`frontend/src/assets/items.json` hat ~2060 Einträge und stammt aus einer
**festgenagelten Wiki-Fassung**:

```js
const baseUrl = "https://pzwiki.net/w/index.php?oldid=509463"; // Build 41.78.16
```

Deshalb enthält der Katalog kein einziges B42-Handwerks- oder Tierhaltungs-Item —
`Anvil`, `Crucible`, `Whetstone`, `Kiln`, `Saddle` ergeben null Treffer. Genau die
Gegenstände, wegen derer man in Build 42 etwas verschenkt, fehlen.

Die Zeile darunter im Skript ist der auskommentierte Live-Link auf
`PZwiki:Item_list`. Ihn zu aktivieren wäre der schnelle Weg und der falsche: Dann
hängt der Katalog am Tagesstand einer Wiki-Seite, und in einem halben Jahr kann
niemand mehr sagen, was drin ist. Die Festnagelung war richtig — nur ist die Nummer
alt und steht als Kommentar an einer Stelle, die niemand liest.

Diese Spec macht die Herkunft zum **Datum der erzeugten Datei**, nicht zur
Kommentarzeile im Generator.

## Technical Requirements

### Die Fassung steht in der Ausgabe

`items.json` bekommt einen Kopf: Quell-URL, Wiki-Revision, Build-Bezeichnung,
Erzeugungsdatum und Anzahl der Einträge. Die heutige Struktur (Liste von Kategorien)
wandert unter einen Schlüssel darunter (**PG-12**).

Damit beantwortet ein Blick in die Datei die Frage „welcher Stand?" — heute
beantwortet sie niemand.

Der Generator liest seine Ziel-Revision aus einer benannten Konstante am Dateikopf und
schreibt genau diese in die Ausgabe. Ein Lauf gegen die Live-Seite ist erlaubt, muss
dann aber die tatsächlich ausgelieferte `oldid` aus der Seite herausziehen und
eintragen — „live" ist kein zulässiger Wert für das Feld.

### Ein fehlendes Bild bricht nichts ab

`downloadImage` fängt Fehler heute ab und protokolliert sie. Diese Nachsicht bleibt,
und sie wird sichtbar: Am Ende des Laufs steht, wie viele Bilder fehlen (**PG-13**).
Ein Katalog mit 3000 Einträgen und 12 fehlenden Bildern ist brauchbar; ein Lauf, der
nach dem zwölften Fehler abbricht, ist es nicht.

### Der Notausgang bleibt

Das Freitextfeld für eigene Item-IDs (seit v1.2.0, Commit `eaa65fb`) bleibt erhalten
und sichtbar — auch nach der Aktualisierung (**PG-14**). Mods bringen eigene IDs mit,
die in keinem Wiki stehen; der Katalog ist eine Bequemlichkeit, kein Gatter.

### Bilder und Einträge bleiben zusammen

Jeder `images`-Pfad in `items.json` zeigt auf eine Datei, die es unter
`frontend/public/items/` auch gibt (**PG-15**). Ein Verweis ins Leere erzeugt in der
Oberfläche eine leere Fläche ohne Erklärung.

Umgekehrt gilt das nicht: Verwaiste Bilder aus früheren Läufen dürfen liegen bleiben,
sie kosten nur Platz. Wer aufräumt, tut das in einem eigenen Schritt und sieht sich
den Diff an.

## Gate-Contract

| ID | Zusage | Prüfung |
|---|---|---|
| PG-12 | Die Herkunft steht in der erzeugten Datei | `jq '.source' frontend/src/assets/items.json` nennt URL, Revision, Build und Datum |
| PG-13 | Ein fehlendes Bild bricht den Lauf nicht ab | `yarn items` mit unterbrochener Netzverbindung für einzelne Dateien: Lauf endet, Bilanz nennt die Fehlzahl |
| PG-14 | Eigene Item-IDs bleiben eingebbar | Item-Dialog öffnen: Freitextfeld vorhanden, `Base.Anvil` lässt sich ohne Katalogeintrag verschenken |
| PG-15 | Kein Bildverweis ohne Datei | `jq -r '..\|.images? // empty\|.[]' items.json \| sed 's\|^/\|frontend/public/\|' \| xargs -I{} test -f {}` läuft fehlerfrei |

## Data Model / Registry

Neue Kopfstruktur von `items.json`:

```json
{
  "source": {
    "url": "https://pzwiki.net/w/index.php?oldid=…",
    "revision": "…",
    "build": "42.20.2",
    "generatedAt": "…",
    "itemCount": 0
  },
  "categories": [ … ]
}
```

Betroffen sind die Leser dieser Datei — `AddItemDialog.tsx` und der
Übersetzungs-Namensraum `items`. Die Item-IDs selbst ändern sich nicht in ihrer Form
(`Base.Axe`, groß-/kleinschreibungsempfindlich).

## Dependencies

* Puppeteer bleibt Entwickler-Abhängigkeit, nicht Laufzeitbestandteil.
* `yarn convert-item-translations` erzeugt weiterhin
  `public/locales/<code>/items.json` — nach einem Katalogwechsel neu laufen lassen,
  sonst zeigen neue Items nur ihre IDs.
* Die Wiki-Seitenstruktur ist eine Fremdannahme: Der Generator liest `h2`-Überschriften
  und `table.wikitable`. Ändert das Wiki sein Layout, bricht der Lauf — das ist
  hinnehmbar, weil er von Hand angestoßen wird.

## Definition of Done

- [ ] `items.mjs` zieht gegen eine benannte B42-Revision
- [ ] `items.json` trägt ihren Herkunftskopf
- [ ] Bilanz am Ende des Laufs (Einträge, Bilder, Fehlende)
- [ ] `AddItemDialog` liest die neue Struktur, Freitextfeld unangetastet
- [ ] Übersetzungen für den Namensraum `items` neu erzeugt
- [ ] `../B42-Luecken.md` Abschnitt 4 abgehakt

## Anti-Patterns

* **Auf den Live-Link umstellen und fertig.** Erzeugt einen Katalog ohne Alter.
* **Den Lauf beim ersten fehlenden Bild abbrechen.** Aus einer Kleinigkeit wird ein
  Blocker.
* **Item-IDs von Hand nachtragen.** Die Datei ist ein Erzeugnis; Handarbeit darin ist
  beim nächsten Lauf weg.
* **Das Freitextfeld nach der Aktualisierung entfernen.** Mods gibt es weiterhin.
* **Verwaiste Bilder im selben Schritt löschen.** Vermischt zwei Diffs, von denen der
  eine geprüft werden will.
