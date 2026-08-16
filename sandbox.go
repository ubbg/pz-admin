package main

// SandboxVars — der zweite Zugangsweg neben RCON.
//
// Die 273 Sandbox-Variablen von Build 42 sind über RCON grundsätzlich nicht
// erreichbar: `changeoption` bedient ausschließlich die servertest.ini. Dieser Teil
// der Anwendung arbeitet deshalb auf der Datei `<Server>/<name>_SandboxVars.lua` —
// und nur auf ihr. Kein Sandbox-Wert geht je über die RCON-Verbindung (PG-16).
//
// Die Variablen selbst sind hier **nicht** deklariert: Gelesen wird, was in der Datei
// steht — dieselbe Haltung wie bei den Serveroptionen und aus demselben Grund
// (specs/sandbox/SANDBOX-001.spec.md).

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// SandboxVar ist eine Variable, so wie sie in der Datei steht: voller Schlüssel
// (mit Gruppe), Gruppe, Wert und der aus dem Wert abgeleitete Typ. `Quoted` merkt
// sich, ob der Wert in Anführungszeichen stand — beim Schreiben wird dieselbe Form
// wiederhergestellt.
//
// Start und End sind die Byte-Grenzen des Wertes in der Datei. Über sie wird beim
// Speichern genau dieser Ausschnitt ersetzt und sonst nichts — unabhängig davon, ob
// die Zuweisung allein auf einer Zeile steht oder mit anderen zusammen.
type SandboxVar struct {
	Key    string `json:"key"`
	Group  string `json:"group"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Kind   string `json:"kind"` // Boolean, Integer, Double, String
	Quoted bool   `json:"quoted"`
	Line   int    `json:"-"`
	Start  int    `json:"-"`
	End    int    `json:"-"`
}

// parseSandboxLua liest die Datei zeichenweise. Es wird kein Lua ausgewertet, aber
// alles gelesen, was in dieser Datei vorkommen kann: Zuweisungen in beliebiger
// Einrückung, Untertabellen über mehrere Zeilen **und** in einer Zeile, Kommentare
// mit `--` und Zeichenketten mit maskierten Anführungszeichen.
//
// Die Datei wird auch von Hand geschrieben; eine Schreibweise, die Lua versteht, darf
// den Editor nicht lahmlegen.
func parseSandboxLua(text string) ([]SandboxVar, error) {
	var vars []SandboxVar
	var groups []string

	position := 0
	line := 1
	depth := 0
	started := false

	// advance zählt Zeilenumbrüche mit, damit Line stimmt.
	advance := func(to int) {
		for i := position; i < to && i < len(text); i++ {
			if text[i] == '\n' {
				line++
			}
		}
		position = to
	}

	for position < len(text) {
		character := text[position]

		switch {
		case character == '\n':
			line++
			position++

		case character == ' ' || character == '\t' || character == '\r' || character == ',' || character == ';':
			position++

		case character == '-' && position+1 < len(text) && text[position+1] == '-':
			end := strings.IndexByte(text[position:], '\n')
			if end < 0 {
				position = len(text)
			} else {
				position += end
			}

		case character == '}':
			depth--
			if len(groups) > 0 {
				groups = groups[:len(groups)-1]
			}
			position++
			if started && depth == 0 {
				return vars, nil
			}

		case isIdentifierStart(character):
			name, after := readIdentifier(text, position)
			assignment := skipSpaces(text, after)

			if assignment >= len(text) || text[assignment] != '=' {
				// Kein `Name =` — etwa `return` oder `local`; weitergehen.
				advance(after)
				continue
			}

			valueStart := skipSpaces(text, assignment+1)
			if valueStart >= len(text) {
				advance(len(text))
				continue
			}

			if text[valueStart] == '{' {
				advance(valueStart + 1)
				depth++
				if !started {
					// Die äußerste Tabelle ist die Wurzel, keine Gruppe.
					started = true
				} else {
					groups = append(groups, name)
				}
				continue
			}

			literal, valueEnd := readValue(text, valueStart)
			if literal == "" {
				advance(valueEnd)
				continue
			}

			if !started {
				// Zuweisungen außerhalb der Tabelle gehören nicht dazu.
				advance(valueEnd)
				continue
			}

			value := literal
			quoted := false
			if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
				quoted = true
				value = unquoteLuaString(value)
			}

			group := strings.Join(groups, ".")
			key := name
			if group != "" {
				key = group + "." + name
			}

			kind := kindOfValue(value)
			if quoted {
				kind = "String"
			}

			advance(valueStart)
			vars = append(vars, SandboxVar{
				Key:    key,
				Group:  group,
				Name:   name,
				Value:  value,
				Kind:   kind,
				Quoted: quoted,
				Line:   line,
				Start:  valueStart,
				End:    valueEnd,
			})
			advance(valueEnd)

		default:
			position++
		}
	}

	if !started {
		return nil, errors.New("no SandboxVars table found in the file")
	}

	return nil, errors.New("the SandboxVars table is not closed — the file looks truncated")
}

func isIdentifierStart(character byte) bool {
	return character == '_' ||
		(character >= 'a' && character <= 'z') ||
		(character >= 'A' && character <= 'Z')
}

func readIdentifier(text string, position int) (string, int) {
	end := position
	for end < len(text) {
		character := text[end]
		if isIdentifierStart(character) || (character >= '0' && character <= '9') {
			end++
			continue
		}
		break
	}
	return text[position:end], end
}

func skipSpaces(text string, position int) int {
	for position < len(text) && (text[position] == ' ' || text[position] == '\t') {
		position++
	}
	return position
}

// readValue liest den Wert genau so, wie er in der Datei steht — Zeichenketten samt
// Anführungszeichen und Maskierungen, alles andere bis zum nächsten Trenner.
func readValue(text string, position int) (string, int) {
	if text[position] == '"' {
		end := position + 1
		for end < len(text) {
			if text[end] == '\\' {
				end += 2
				continue
			}
			if text[end] == '"' {
				end++
				break
			}
			end++
		}
		return text[position:end], end
	}

	end := position
	for end < len(text) {
		character := text[end]
		if character == ',' || character == '}' || character == '\n' || character == ';' {
			break
		}
		if character == '-' && end+1 < len(text) && text[end+1] == '-' {
			break
		}
		end++
	}

	value := strings.TrimRight(text[position:end], " \t\r")
	return value, position + len(value)
}

// unquoteLuaString nimmt die Anführungszeichen weg und löst die Maskierungen auf,
// die in dieser Datei vorkommen.
func unquoteLuaString(literal string) string {
	inner := literal[1 : len(literal)-1]
	replaced := strings.Builder{}

	for i := 0; i < len(inner); i++ {
		if inner[i] == '\\' && i+1 < len(inner) {
			switch inner[i+1] {
			case '"', '\\':
				replaced.WriteByte(inner[i+1])
				i++
				continue
			case 'n':
				replaced.WriteByte('\n')
				i++
				continue
			}
		}
		replaced.WriteByte(inner[i])
	}

	return replaced.String()
}

// quoteLuaString schreibt eine Zeichenkette so, wie Lua sie wieder liest.
func quoteLuaString(value string) string {
	escaped := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`).Replace(value)
	return `"` + escaped + `"`
}

// Zugangswege. SFTP steht in der Spec, kommt aber bewusst später — solange nur
// diese beiden, damit kein halber Weg entsteht.
const (
	sandboxModeLocal  = "local"
	sandboxModeDocker = "docker"
)

// sandboxFileSuffix ist der Namensteil, den der Server allen Sandbox-Dateien gibt;
// davor steht der Servername (Standard „servertest").
const sandboxFileSuffix = "_SandboxVars.lua"

// resolveSandboxPath bestimmt die Datei, die bearbeitet wird. Geraten wird nichts:
// Ohne eingerichteten Zugangsweg gibt es einen Fehler statt eines Standardpfads, und
// liegen mehrere Serverdateien nebeneinander, entscheidet die Nutzerin (PG-20).
//
// `list` liefert die Einträge eines Verzeichnisses — im Betrieb aus dem Dateisystem,
// im Test aus einer Tabelle.
func resolveSandboxPath(mode string, configured string, list func(string) ([]string, error)) (string, error) {
	configured = strings.TrimSpace(configured)

	if mode == "" || configured == "" {
		return "", errors.New("no sandbox access configured")
	}

	var directory string

	switch mode {
	case sandboxModeLocal:
		if strings.HasSuffix(configured, ".lua") {
			return configured, nil
		}
		directory = configured
	case sandboxModeDocker:
		// Der Danixu-Container bindet das Volume auf /home/steam/Zomboid; die
		// Sandbox-Datei liegt darin unter Server/.
		directory = filepath.Join(configured, "Server")
	default:
		return "", fmt.Errorf("unknown sandbox access mode: %s", mode)
	}

	entries, err := list(directory)
	if err != nil {
		return "", fmt.Errorf("cannot read %s: %w", directory, err)
	}

	var candidates []string
	for _, entry := range entries {
		if strings.HasSuffix(entry, sandboxFileSuffix) {
			candidates = append(candidates, entry)
		}
	}

	switch len(candidates) {
	case 0:
		return "", fmt.Errorf("no *%s found in %s", sandboxFileSuffix, directory)
	case 1:
		return filepath.Join(directory, candidates[0]), nil
	default:
		return "", fmt.Errorf("several sandbox files in %s (%s) — name the file itself", directory, strings.Join(candidates, ", "))
	}
}

// writeSandboxValues ersetzt in der Datei genau die Werte, die sich geändert haben —
// zeichengenau an ihrer Stelle. Kommentare, Reihenfolge, Einrückung und die
// Untertabellen bleiben unangetastet (PG-18): Diese Datei wird von Menschen gelesen,
// eine umsortierte Fassung macht jeden künftigen Vergleich unbrauchbar.
//
// Ersetzt wird von hinten nach vorn, damit die Grenzen der noch offenen Werte gültig
// bleiben. Zurück kommen der neue Dateiinhalt und die Schlüssel, die tatsächlich
// geändert wurden.
func writeSandboxValues(text string, vars []SandboxVar, values map[string]string) (string, []string, error) {
	byKey := make(map[string]SandboxVar, len(vars))
	for _, variable := range vars {
		byKey[variable.Key] = variable
	}

	// Namen, die die Datei nicht kennt, werden nicht eingefügt — die Datei bestimmt,
	// welche Variablen es gibt.
	for key := range values {
		if _, ok := byKey[key]; !ok {
			return "", nil, fmt.Errorf("unknown sandbox variable: %s", key)
		}
	}

	var pending []SandboxVar
	var changed []string

	// In der Reihenfolge der Datei, damit die Rückmeldung nachvollziehbar bleibt.
	for _, variable := range vars {
		value, ok := values[variable.Key]
		if !ok || value == variable.Value {
			continue
		}
		if !sandboxValueIsWellFormed(variable, value) {
			return "", nil, fmt.Errorf("value %q does not fit the type %s of %s", value, variable.Kind, variable.Key)
		}
		pending = append(pending, variable)
		changed = append(changed, variable.Key)
	}

	updated := text
	for i := len(pending) - 1; i >= 0; i-- {
		variable := pending[i]
		if variable.Start < 0 || variable.End > len(updated) || variable.Start > variable.End {
			return "", nil, fmt.Errorf("the position of %s is outside the file", variable.Key)
		}
		literal := sandboxValueLiteral(variable, values[variable.Key])
		updated = updated[:variable.Start] + literal + updated[variable.End:]
	}

	return updated, changed, nil
}

// sandboxValueLiteral bringt einen Wert in die Schreibweise, in der er in der Datei
// steht — Zeichenketten in Anführungszeichen, Zahlen und Wahrheitswerte blank.
func sandboxValueLiteral(variable SandboxVar, value string) string {
	if variable.Quoted {
		return quoteLuaString(value)
	}
	return value
}

// sandboxValueIsWellFormed prüft, ob ein neuer Wert zum Typ der Variablen passt.
// Ein Tippfehler soll die Datei nicht erreichen.
func sandboxValueIsWellFormed(variable SandboxVar, value string) bool {
	switch variable.Kind {
	case "Boolean":
		return isBooleanValue(value)
	case "Integer":
		_, err := strconv.ParseInt(value, 10, 64)
		return err == nil
	case "Double":
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	default:
		// Zeichenketten werden beim Schreiben maskiert (quoteLuaString), deshalb
		// darf hier alles stehen — auch Anführungszeichen.
		return true
	}
}
