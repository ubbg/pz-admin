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
// sich, ob der Wert in der Datei in Anführungszeichen stand — beim Schreiben wird
// dieselbe Form wiederhergestellt.
type SandboxVar struct {
	Key    string `json:"key"`
	Group  string `json:"group"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Kind   string `json:"kind"` // Boolean, Integer, Double, String
	Quoted bool   `json:"quoted"`
	Line   int    `json:"line"` // 1-basiert, für das strukturerhaltende Schreiben
}

// parseSandboxLua liest die Datei zeilenweise. Es wird nicht Lua ausgewertet, sondern
// die Schreibweise gelesen, die der Server erzeugt: `Name = Wert,` in beliebiger
// Einrückung, Untertabellen als `Name = {` … `},`.
func parseSandboxLua(text string) ([]SandboxVar, error) {
	lines := strings.Split(text, "\n")

	started := false
	depth := 0
	var groups []string
	var vars []SandboxVar

	for index, raw := range lines {
		line := strings.TrimSpace(stripLuaComment(raw))
		if line == "" {
			continue
		}

		if !started {
			if strings.HasPrefix(line, "SandboxVars") && strings.Contains(line, "{") {
				started = true
				depth = 1
			}
			continue
		}

		if strings.HasPrefix(line, "}") {
			depth--
			if len(groups) > 0 {
				groups = groups[:len(groups)-1]
			}
			if depth == 0 {
				return vars, nil
			}
			continue
		}

		name, value, ok := splitLuaAssignment(line)
		if !ok {
			continue
		}

		if strings.HasPrefix(value, "{") {
			groups = append(groups, name)
			depth++
			continue
		}

		value = strings.TrimSuffix(value, ",")
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		quoted := false
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			quoted = true
			value = value[1 : len(value)-1]
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

		vars = append(vars, SandboxVar{
			Key:    key,
			Group:  group,
			Name:   name,
			Value:  value,
			Kind:   kind,
			Quoted: quoted,
			Line:   index + 1,
		})
	}

	if !started {
		return nil, errors.New("no SandboxVars table found in the file")
	}

	return nil, errors.New("the SandboxVars table is not closed — the file looks truncated")
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
// zeilenweise, an Ort und Stelle. Kommentare, Reihenfolge, Einrückung und die
// Untertabellen bleiben unangetastet (PG-18): Diese Datei wird von Menschen gelesen,
// eine umsortierte Fassung macht jeden künftigen Vergleich unbrauchbar.
//
// Zurück kommen der neue Dateiinhalt und die Schlüssel, die tatsächlich geändert
// wurden.
func writeSandboxValues(text string, vars []SandboxVar, values map[string]string) (string, []string, error) {
	byKey := make(map[string]SandboxVar, len(vars))
	for _, variable := range vars {
		byKey[variable.Key] = variable
	}

	lines := strings.Split(text, "\n")
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

		index := variable.Line - 1
		if index < 0 || index >= len(lines) {
			return "", nil, fmt.Errorf("line %d for %s is outside the file", variable.Line, variable.Key)
		}

		replaced, ok := replaceValueInLine(lines[index], variable, value)
		if !ok {
			return "", nil, fmt.Errorf("could not find the value of %s in line %d", variable.Key, variable.Line)
		}

		lines[index] = replaced
		changed = append(changed, variable.Key)
	}

	// Namen, die die Datei nicht kennt, werden nicht eingefügt — die Datei bestimmt,
	// welche Variablen es gibt.
	for key := range values {
		if _, ok := byKey[key]; !ok {
			return "", nil, fmt.Errorf("unknown sandbox variable: %s", key)
		}
	}

	return strings.Join(lines, "\n"), changed, nil
}

// replaceValueInLine tauscht nur den Wert aus und lässt alles andere an der Zeile
// stehen: Einrückung, Leerraum um das Gleichheitszeichen, Komma und Kommentar.
func replaceValueInLine(line string, variable SandboxVar, value string) (string, bool) {
	assignment := strings.Index(line, "=")
	if assignment < 0 || strings.TrimSpace(line[:assignment]) != variable.Name {
		return "", false
	}

	rest := line[assignment+1:]
	oldLiteral := sandboxValueLiteral(variable, variable.Value)

	start := strings.Index(rest, oldLiteral)
	if start < 0 {
		return "", false
	}

	newLiteral := sandboxValueLiteral(variable, value)
	rest = rest[:start] + newLiteral + rest[start+len(oldLiteral):]

	return line[:assignment+1] + rest, true
}

// stripLuaComment entfernt einen Kommentar am Zeilenende, ohne einen `--` innerhalb
// einer Zeichenkette anzutasten.
func stripLuaComment(line string) string {
	inString := false

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '"':
			inString = !inString
		case '-':
			if !inString && i+1 < len(line) && line[i+1] == '-' {
				return line[:i]
			}
		}
	}

	return line
}

// splitLuaAssignment zerlegt `Name = Wert` in seine beiden Teile.
func splitLuaAssignment(line string) (string, string, bool) {
	index := strings.Index(line, "=")
	if index <= 0 {
		return "", "", false
	}

	name := strings.TrimSpace(line[:index])
	value := strings.TrimSpace(line[index+1:])

	if name == "" || value == "" || !isLuaIdentifier(name) {
		return "", "", false
	}

	return name, value, true
}

func isLuaIdentifier(name string) bool {
	for i, character := range name {
		switch {
		case character >= 'a' && character <= 'z',
			character >= 'A' && character <= 'Z',
			character == '_':
		case character >= '0' && character <= '9':
			if i == 0 {
				return false
			}
		default:
			return false
		}
	}

	return name != ""
}

// sandboxValueLiteral bringt einen Wert in die Schreibweise, in der er in der Datei
// steht — Zeichenketten in Anführungszeichen, Zahlen und Wahrheitswerte blank.
func sandboxValueLiteral(variable SandboxVar, value string) string {
	if variable.Quoted {
		return fmt.Sprintf("%q", value)
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
		return !strings.ContainsAny(value, "\"\n")
	}
}
