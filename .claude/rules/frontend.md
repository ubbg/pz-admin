---
paths:
  - "frontend/src/**/*.tsx"
  - "frontend/src/**/*.ts"
  - "frontend/index.html"
  - "frontend/tailwind.config.js"
---

# Oberfläche — React 19, shadcn/ui, und der Zustand kommt aus Providern

Das Frontend ist ein Vite-React-Projekt in TypeScript. Es hat **keinen** Store
(kein Redux, kein Zustand, kein Jotai): Geteilter Zustand liegt in Contexts unter
`src/contexts/`, alles andere ist lokaler Komponentenzustand.

## Der Pfad zu Go

```ts
import { KickUsers } from "@/wailsjs/go/main/App";
import { main } from "@/wailsjs/go/models";
```

* `@/` ist das Alias auf `frontend/src`. Keine relativen Ketten (`../../..`).
* **`src/wailsjs/` ist generiert.** Nie von Hand bearbeiten — der Ordner entsteht bei
  jedem `wails dev` / `wails build` neu. Fehlt dort eine Funktion, fehlt in Go die
  exportierte Methode auf `*App`, nicht hier ein Import.
* Typen kommen aus `main.*` (`main.Player`, `main.Option`, `main.Credentials`,
  `main.RconResponse`). Keine handgeschriebenen Zwillinge dieser Typen anlegen — sie
  laufen sofort auseinander.
* Ereignisse aus Go über `EventsOn` / `EventsOff` aus `@/wailsjs/runtime/runtime`.
  Jedes `EventsOn` braucht sein `EventsOff` im Cleanup, sonst hängen nach einem
  Reconnect zwei Handler an `update-players`.

## Provider statt Prop-Ketten

Unter `src/contexts/` liegt je ein Provider für RCON-Zustand, Konfiguration, Theme,
Farbschema, Fortschritt, Storage, OS und Neustart; `providers.tsx` setzt sie
zusammen.

`rcon-provider.tsx` ist der größte und der einzige, der wirklich Fachlogik trägt:
Verbindung, Spielerliste, die **zwei** Optionsstände (`options` = Serverstand,
`modifiedOptions` = ungespeicherte Eingaben) und deren Vergleich. Wer ein neues
Serverdatum in die Oberfläche bringt, hängt es dort an — nicht in eine Komponente,
die es per `useEffect` selbst holt.

Neue Werte, die aus Go kommen, sind `main.*`-Typen und werden im Provider einmal
entgegengenommen; Komponenten lesen über `useRcon()` / `useConfig()`.

## Komponenten

* `AdminPanel.tsx` ist die Klammer; darunter je ein Bereich pro Tab
  (`Players`, `Management`, `Options`, `Tools`, `WeatherControl`, `Terminal`,
  `Settings`).
* **Jede Aktion mit Parametern ist ein Dialog** unter `components/Dialogs/`, benannt
  nach der Aktion (`KickUserDialog`, `AddItemDialog`, …). Muster: `isOpen`, `onClose`,
  fachliche Props (`names: string[]`); `onClose()` **vor** dem Go-Aufruf, damit der
  Dialog nicht über der Fortschrittsanzeige stehen bleibt; `useEffect` auf `isOpen`
  setzt die Felder zurück.
* `components/ui/` ist **shadcn/ui** — generierte Radix-Bausteine. Dort nur ändern,
  was bewusst projektweit gelten soll; sonst wird die nächste Aktualisierung eines
  Bausteins zum Konfliktfeld.
* Tabellen laufen über `@tanstack/react-table`, Formulare über `react-hook-form` +
  `zod`, Toasts über `sonner`, Icons über `lucide-react`. Keine zweite Bibliothek für
  dasselbe hinzufügen.

## Texte

Kein sichtbarer String im JSX. Immer `const { t } = useTranslation();` und
`t("admin_panel.tabs.players.…")`. Details und Katalogpflege:
`.claude/rules/i18n.md`.

## Styling

Tailwind mit `tailwind-merge` und `class-variance-authority`; Farbschemata liegen in
`src/colorSchemes.json` und laufen über den `color-scheme-provider`. Kein Inline-CSS,
keine eigene `.css`-Datei neben `index.css`. Dunkel- und Hellmodus laufen über
`next-themes` — eine neue Farbe wird als Token im Schema ergänzt, nicht als
Hex-Wert in die Komponente geschrieben.

## Was das Frontend nicht tut

* **Keine Kommandozeilen bauen.** Die Oberfläche ruft `KickUsers(names, reason)` auf,
  nicht `SendRconCommand("kick \"" + name + "\" -r ...")`. Rohzugriff gibt es genau
  einmal, im `Terminal`, und dort ausdrücklich auf Wunsch der Nutzerin.
* **Keine Serverantworten interpretieren.** Erfolg und Misserfolg entscheidet Go
  (`SuccessCheck`) und meldet sie als Benachrichtigung.
* **Keine Build-Annahmen verstecken.** Die Befehlsliste in `Terminal.tsx` und die
  Optionsmetadaten in `assets/options.ts` sind versionsgebundene Daten — sie gehören
  gepflegt wie Daten, siehe `.claude/rules/options.md` und
  `.claude/rules/catalogs.md`.

## Lint

`cd frontend && yarn lint` läuft mit `--max-warnings 0`. Es gibt keine Testsuite und
keinen CI-Lint-Schritt — der Lauf von Hand ist die einzige Prüfung, die es gibt.
