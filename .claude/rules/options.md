---
paths:
  - "rcon_options.go"
  - "frontend/src/assets/options.ts"
---

# Serveroptionen — heute drei Orte für eine Option, und genau das ist das Problem

Eine einzelne Serveroption existiert im Projekt derzeit **dreifach**:

1. als Feld im Struct `PzOptions` (`rcon_options.go:15`) — entscheidet, ob sie
   überhaupt gelesen und geschrieben werden kann,
2. als Eintrag in `frontend/src/assets/options.ts` — entscheidet, wie sie aussieht
   (Kategorie, Typ, Bereich, Auswahl, Abhängigkeiten),
3. als Schlüssel in allen vier Sprachkatalogen — entscheidet, wie sie heißt.

Fehlt einer der drei, ist die Option unsichtbar, unbenannt oder nicht änderbar. Das
ist die Ursache der halben B42-Lückenliste — **vor** einer Änderung hier
`specs/options/OPTIONS-001.spec.md` lesen.

## Wie der Rundlauf heute funktioniert

```
showoptions ──► parseOptions() ──► PzOptions ──► EventsEmit("update-options") ──► UI
UI ──► UpdatePzOptions() ──► diffOptions() ──► changeoption … ──► reloadoptions
```

Zwei Stellen darin sind der eigentliche Engpass:

* **`parseOptions`** sucht zu jedem `* Name=Wert` per Reflection ein gleichnamiges
  Struct-Feld. Findet es keins, wird die Zeile mit `LogDebugf("Unknown option: %s")`
  **verworfen**.
* **`diffOptions`** läuft über die Struct-Felder, nicht über die Serverantwort. Was
  nicht im Struct steht, kann also nicht einmal versehentlich geschrieben werden.

Folge: Ein B42-Server meldet 144 Optionen, die Anwendung kennt 135 davon nur zum Teil
(94 gemeinsam), und 41 Struct-Felder gibt es auf dem Server gar nicht mehr.

## Regeln für Änderungen am heutigen Stand

* **Nie nur eine der drei Stellen anfassen.** Ein neues Feld im Struct ohne Eintrag in
  `options.ts` erscheint nirgends; ein Eintrag in `options.ts` ohne Struct-Feld
  erzeugt ein Bedienelement, dessen Wert beim Speichern verschwindet.
* **Feldname = Optionsname des Servers, buchstabengetreu.** Der Abgleich läuft über
  `FieldByName`; `server_browser_announced_ip` bleibt deshalb kleingeschrieben und
  `SafeHouseRemovalTime` behält sein großes `H`.
* **`json:"…"`-Tag identisch zum Feldnamen.** Es landet in
  `frontend/src/wailsjs/go/models.ts` und ist der Name, den das Frontend sieht.
* **Typwahl:** `bool`, `int`, `float64`, `string` — der Parser (`setFieldValue`) kennt
  nichts anderes. Eine vierwertige Option (B42-Anti-Cheat) ist ein `int` mit
  `Type: "Choice"` in `options.ts`, **kein** `bool`.
* **Reihenfolge und Kategorie** stehen ausschließlich in `options.ts`
  (`categories[].options[]`). Das Struct ist alphabetisch — dabei bleiben, es liest
  sich sonst nicht mehr.
* **`Requirements`** in `options.ts` blendet abhängige Felder aus (z. B. alles unter
  `VoiceEnable`). Ein ausgeblendetes Feld wird trotzdem gesendet, wenn es sich
  geändert hat — Abhängigkeit ist Darstellung, keine Sperre.

## B42-Namen, die heute falsch sind

Kurzfassung; vollständig mit Begründung in `../B42-Luecken.md`:

| im Code | in Build 42 |
|---|---|
| `AntiCheatProtectionType1…24` + 7 `…ThresholdMultiplier` | 10 benannte Optionen (`AntiCheatSafety`, `AntiCheatHit`, `AntiCheatNoClip`, …), **vierwertig**: 1 = Ban, 2 = Kick, 3 = Log, 4 = Disable |
| `DiscordChannel`, `DiscordChannelID` | `DiscordChatChannel`, `DiscordLogChannel`, `DiscordCommandChannel` |
| `DisableSafehouseWhenPlayerConnected` | `DisableSafehouseWhenOwnerConnected` |
| `HoursForLootRespawn`, `MaxItemsForLootRespawn`, `ConstructionPreventsLootRespawn`, `MinutesPerPage` | in die **SandboxVars** verschoben — gehören nicht mehr hierher |
| `AutoCreateUserInWhiteList`, `KickFastPlayers` | ersatzlos entfallen |

Ganz fehlen unter anderem die B42-Blöcke **Krieg** (`War*`), **Tiere**
(`AnnounceAnimalDeath`, `UltraSpeedDoesnotAffectToAnimals`), **Verkleidungen**
(`UsernameDisguises`, `HideDisguisedUserName`, `SafehouseDisableDisguises`),
**Chat-Moderation** (`ChatMessage*`, `BadWord*`, `GoodWordListFile`) und
**Abschleppen** (`Disable*Towing`).

## Fehlermeldungen müssen die Option benennen

`applyOptions` zählt heute nur, wie viele `changeoption`-Aufrufe fehlgeschlagen sind
(`failed_to_update_n_options`). Auf einem B42-Server schlagen die toten Felder
zuverlässig fehl — eine Zahl ohne Namen ist dann nutzlos. Wer hier etwas ändert,
nennt die betroffene Option in der Meldung (PG-06).

## Passwörter

`Password`, `RCONPassword` und Ähnliches fehlen im Struct. Falls sie aufgenommen
werden: Sie dürfen nicht ins Log, nicht in eine Benachrichtigung und im UI nur
maskiert erscheinen. Ein Optionen-Editor, der das RCON-Passwort ändert, kappt
außerdem die eigene Verbindung — das gehört abgefangen, nicht ausprobiert.
