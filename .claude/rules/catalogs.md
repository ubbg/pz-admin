---
paths:
  - "frontend/scripts/*.mjs"
  - "frontend/src/assets/items.json"
  - "frontend/src/assets/vehicles.json"
  - "frontend/public/items/**"
  - "frontend/public/vehicles/**"
---

# Kataloge — generierte Daten, keine Handarbeit

`items.json` und `vehicles.json` sind **Erzeugnisse**. Wer dort eine Zeile von Hand
ändert, verliert sie beim nächsten Lauf des Generators.

| Datei | Erzeugt von | Quelle |
|---|---|---|
| `src/assets/items.json` + `public/items/*.png` | `scripts/items.mjs` (`yarn items`) | pzwiki, per Puppeteer |
| `src/assets/vehicles.json` | `scripts/vehicles.mjs` (`yarn vehicles`) | die Bilder unter `public/vehicles/` |
| `public/locales/<code>/items.json` | `scripts/convert-item-translations.mjs` | die `ItemName_*`-Zeilen der Spiel-Übersetzungsdateien |

`yarn dev` und `yarn build` rufen `vehicles.mjs` **immer** vorher auf — dieser Katalog
ist also aus dem Dateibaum abgeleitet und muss nicht gepflegt werden. `items.mjs`
läuft nur auf Zuruf.

## Die Versionsbindung steht im Skript

`scripts/items.mjs` zeigt auf eine **festgenagelte Wiki-Fassung**:

```js
const baseUrl = "https://pzwiki.net/w/index.php?oldid=509463"; // Build 41.78.16
//const baseUrl = "https://pzwiki.net/wiki/PZwiki:Item_list";
```

Das ist der Grund, warum der Katalog ~2060 Einträge hat und kein einziges
B42-Handwerks- oder Tierhaltungs-Item (`Anvil`, `Crucible`, `Whetstone`, `Kiln`,
`Saddle`: null Treffer). Der auskommentierte Live-Link daneben ist die
naheliegende, aber gefährliche Alternative: Er macht den Katalog vom Tagesstand einer
Wiki-Seite abhängig, und niemand sieht später, gegen welchen Build erzeugt wurde.

**Regel:** Der Generator zeigt auf eine *benannte* Fassung, und die Fassung steht
sowohl im Skript als auch als Metadatum in der erzeugten Datei. Arbeitspaket dazu:
`specs/catalog/ITEMS-001.spec.md` (PG-12).

## Wenn du den Katalog neu erzeugst

* Der Lauf lädt **Bilder mit**. `public/items/` wächst entsprechend; ein Eintrag ohne
  Bilddatei erzeugt in der Oberfläche ein leeres Feld, kein Fallback.
* Puppeteer mit Stealth-Plugin lädt einen kompletten Chromium herunter. Das ist eine
  Entwickler-Abhängigkeit, kein Laufzeitbestandteil — sie gehört nicht in den Build.
* Ein einzelner fehlgeschlagener Bild-Download darf den Lauf nicht abbrechen
  (heute: `downloadImage` fängt ab und loggt). Diese Nachsicht beibehalten.
* Nach dem Lauf: Diff ansehen. Ein Katalogwechsel, der bestehende Item-IDs
  **entfernt**, macht gespeicherte Auswahlen der Nutzerinnen ungültig — das ist eine
  bewusste Entscheidung, kein Nebeneffekt.

## Struktur von `items.json`

Eine Liste von Kategorien, je Kategorie eine Liste von Items:

```json
[{ "name": "Accessory",
   "items": [{ "name": "Helmet - Airforce",
               "itemId": "Base.Hat_SPHhelmet",
               "images": ["/items/Base.Hat_SPHhelmet_0.png"] }] }]
```

`itemId` ist der Wert, der im RCON-Befehl landet (`additem "user" Base.Axe 5`) und
ist **groß-/kleinschreibungsempfindlich** — `Base.Axe` funktioniert, `base.axe` nicht.
`name` ist nur Anzeige und wird über den Namensraum `items` übersetzt.

## Der Notausgang bleibt

Seit v1.2.0 kann im Item-Dialog eine **freie Item-ID** eingegeben werden (Commit
`eaa65fb`). Das ist die einzige Möglichkeit, an B42-Items heranzukommen, solange der
Katalog alt ist — dieses Feld nicht wegoptimieren, auch nicht nach einer
Katalogaktualisierung (PG-14).
