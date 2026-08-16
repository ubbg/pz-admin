---
paths:
  - "rcon_options.go"
  - "frontend/src/assets/options.ts"
---

# Serveroptionen — der Server sagt, was es gibt

Eine Serveroption existiert im Projekt an **zwei** Orten, und nur einer davon
entscheidet über ihre Existenz:

1. **Der Server.** `showoptions` meldet sie; `parseOptionLines()` übernimmt jede
   gemeldete Zeile. Steht eine Option nicht in der Antwort, gibt es sie nicht — im
   Go-Code ist **kein** Optionsname deklariert (PG-02).
2. `frontend/src/assets/options.ts` — die **Darstellungstabelle**: Kategorie, Typ,
   Bereich, Auswahlwerte, Abhängigkeiten. Fehlt hier ein Eintrag, erscheint die Option
   trotzdem, nur schlichter (Kategorie „Other Options", Rohname als Beschriftung).
3. Die vier Sprachkataloge geben ihr einen Namen. Fehlt der Schlüssel, steht der
   Rohname da — kein leeres Feld.

Der Umbau dahin ist `specs/options/OPTIONS-001.spec.md`; vor Änderungen hier lesen.

## Der Rundlauf

```
showoptions ──► parseOptionLines() ──► []Option{Name, Value, Kind}
            ──► EventsEmit("update-options-list") ──► rcon-provider ──► Options-Tab

Options-Tab ──► UpdateOptions(map[string]string, reload)
            ──► diffServerOptions() ──► changeoption je Änderung ──► reloadoptions
```

* **`parseOptionLines`** kennt keine Namen. Der Typ kommt aus dem Wert
  (`kindOfValue`): `true`/`false` → Boolean, reine Ganzzahl → Integer, Zahl mit Punkt
  → Double, alles Übrige → String.
* **`diffServerOptions`** läuft über die **Serverantwort**, nicht über eine Liste im
  Code: Gesendet wird nur, was sich geändert hat, und nur für Namen, die der Server
  gemeldet hat (PG-03). `optionValuesEqual` vergleicht dabei Zahlen numerisch
  („70" = „70.0") und Wahrheitswerte ohne Rücksicht auf Groß-/Kleinschreibung.
* **`optionUpdateSucceeded`** prüft die Antwort `Option : <Name> is now : <Wert>`.
  Ein PZ-Server antwortet auch auf Unsinn freundlich — `err == nil` ist kein Erfolg.
* **`applyServerOptions`** hält `connMutex` und liefert die Namen zurück, die der
  Server abgelehnt hat. Die Meldung nennt sie (PG-06); eine bloße Anzahl ist beim
  Wechsel des Spiel-Builds nutzlos.

## Regeln für Änderungen

* **Kein Optionsname in Go.** Wer einen braucht, hat den falschen Weg gewählt: Die
  Existenz entscheidet der Server, die Darstellung `options.ts`.
* **Feldname = Optionsname des Servers, buchstabengetreu.** Der Abgleich läuft über
  die Zeichenkette; `server_browser_announced_ip` bleibt kleingeschrieben,
  `SafeHouseRemovalTime` behält sein großes `H`.
* **Typwahl in `options.ts` ist Darstellung, keine Wahrheit.** Bei Widerspruch
  gewinnt der Server. Eine vierwertige Option (B42-Anti-Cheat) ist `Type: "Choice"`
  mit vier `Choices` — der Server liefert dazu `1`…`4`.
* **Reihenfolge und Kategorie** stehen ausschließlich in `options.ts`
  (`categories[].options[]`). Was dort fehlt, landet alphabetisch unter
  „Other Options" — das ist zugleich der Melder, dass etwas einzusortieren ist.
* **`Requirements`** blendet abhängige Felder aus (z. B. alles unter `VoiceEnable`).
  Ein ausgeblendetes Feld wird trotzdem gesendet, wenn es sich geändert hat —
  Abhängigkeit ist Darstellung, keine Sperre.
* **Neue Beschriftung = vier Kataloge.** `display_name`, `description`, `keywords` je
  Option; siehe `.claude/rules/i18n.md`.

## Build-42-Besonderheiten

| früher (B41) | Build 42 |
|---|---|
| `AntiCheatProtectionType1…24` + 7 `…ThresholdMultiplier` | 10 benannte Optionen (`AntiCheatSafety`, `AntiCheatHit`, `AntiCheatNoClip`, …), **vierwertig**: 1 = Ban, 2 = Kick, 3 = Log, 4 = Disable |
| `DiscordChannel`, `DiscordChannelID` | `DiscordChatChannel`, `DiscordLogChannel`, `DiscordCommandChannel` |
| `DisableSafehouseWhenPlayerConnected` | `DisableSafehouseWhenOwnerConnected` |
| `HoursForLootRespawn`, `MaxItemsForLootRespawn`, `ConstructionPreventsLootRespawn`, `MinutesPerPage` | in die **SandboxVars** verschoben — gehören nicht mehr hierher |
| `AutoCreateUserInWhiteList`, `KickFastPlayers` | ersatzlos entfallen |

`Mods` trägt in Build 42 je Eintrag ein führendes `\` (`\mod1;\mod2`). Die Umwandlung
sitzt an einer Stelle (`modsValueForServer` / `modsValueForDisplay`), nicht in der
Oberfläche.

## Passwörter

`Password`, `RCONPassword` und Ähnliches meldet der Server mit. Sie dürfen nicht ins
Log, nicht in eine Benachrichtigung und im UI nur maskiert erscheinen. Ein
Optionen-Editor, der das RCON-Passwort ändert, kappt außerdem die eigene Verbindung —
das gehört abgefangen, nicht ausprobiert.
