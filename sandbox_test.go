package main

import (
	"errors"
	"strings"
	"testing"
)

// Ein Ausschnitt aus servertest_SandboxVars.lua in der Form, die Build 42 schreibt:
// Kopfzeile, Kommentare, flache Werte und Untertabellen.
const sandboxSample = `-- Sandbox settings for the server
SandboxVars = {
    VERSION = 5,
    Zombies = 3,
    Distribution = 1,
    DayLength = 3,          -- 3 = 1 hour
    StartYear = 1,
    WaterShutModifier = 14,
    ZombieAttractionMultiplier = 1.0,
    EnableSnowOnGround = true,
    LootItemRemovalList = "Base.Vest,Base.Sock",
    Husbandry = 3,

    Map = {
        AllowMiniMap = false,
        AllowWorldMap = true,
        MapAllKnown = false,
    },

    ZombieLore = {
        Speed = 2,
        Strength = 3,
        Memory = 2,
    },
}
`

func sandboxVarByKey(vars []SandboxVar, key string) (SandboxVar, bool) {
	for _, variable := range vars {
		if variable.Key == key {
			return variable, true
		}
	}
	return SandboxVar{}, false
}

// Angezeigt wird, was in der Datei steht — auch Variablen, die die Anwendung nicht
// kennt. Kein Variablenname ist im Code deklariert (Ticket 13).
func TestParseSandboxLuaReadsEveryVariable(t *testing.T) {
	vars, err := parseSandboxLua(sandboxSample)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}

	if len(vars) != 16 {
		t.Fatalf("erwartet 16 Variablen, bekommen %d: %v", len(vars), vars)
	}

	cases := []struct {
		key   string
		group string
		value string
		kind  string
	}{
		{"VERSION", "", "5", "Integer"},
		{"DayLength", "", "3", "Integer"},
		{"ZombieAttractionMultiplier", "", "1.0", "Double"},
		{"EnableSnowOnGround", "", "true", "Boolean"},
		{"LootItemRemovalList", "", "Base.Vest,Base.Sock", "String"},
		{"Husbandry", "", "3", "Integer"},
		{"Map.AllowMiniMap", "Map", "false", "Boolean"},
		{"Map.MapAllKnown", "Map", "false", "Boolean"},
		{"ZombieLore.Speed", "ZombieLore", "2", "Integer"},
	}

	for _, testCase := range cases {
		variable, ok := sandboxVarByKey(vars, testCase.key)
		if !ok {
			t.Errorf("%s fehlt in der gelesenen Liste", testCase.key)
			continue
		}
		if variable.Group != testCase.group {
			t.Errorf("%s: Gruppe = %q, erwartet %q", testCase.key, variable.Group, testCase.group)
		}
		if variable.Value != testCase.value {
			t.Errorf("%s: Wert = %q, erwartet %q", testCase.key, variable.Value, testCase.value)
		}
		if variable.Kind != testCase.kind {
			t.Errorf("%s: Typ = %q, erwartet %q", testCase.key, variable.Kind, testCase.kind)
		}
	}
}

// Eine Datei mit ungewohntem Aufbau führt zu einer verständlichen Meldung, nicht zu
// einem Absturz und nicht zu einer stillen leeren Liste (Ticket 13).
func TestParseSandboxLuaRejectsAnUnexpectedFile(t *testing.T) {
	cases := []struct {
		name string
		text string
	}{
		{"leer", ""},
		{"nur Kommentare", "-- nichts als ein Kommentar\n"},
		{"fremde Datei", "return { foo = 1 }\n"},
		{"abgeschnitten", "SandboxVars = {\n    Zombies = 3,\n"},
	}

	for _, testCase := range cases {
		if _, err := parseSandboxLua(testCase.text); err == nil {
			t.Errorf("%s: erwartet eine Fehlermeldung, bekommen keine", testCase.name)
		}
	}
}

// Geschrieben wird an Ort und Stelle: Kommentare, Reihenfolge, Einrückung und die
// Untertabellen bleiben, wie sie sind — geändert werden nur die betroffenen
// Wertzeilen (PG-18).
func TestWriteSandboxValuesTouchesOnlyChangedValueLines(t *testing.T) {
	vars, err := parseSandboxLua(sandboxSample)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}

	updated, changed, err := writeSandboxValues(sandboxSample, vars, map[string]string{
		"Zombies":             "1",
		"EnableSnowOnGround":  "false",
		"LootItemRemovalList": "Base.Vest",
		"Map.AllowMiniMap":    "true",
		"DayLength":           "3", // unverändert -> keine Änderung
	})
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}

	if len(changed) != 4 {
		t.Fatalf("erwartet 4 geänderte Werte, bekommen %d: %v", len(changed), changed)
	}

	oldLines := strings.Split(sandboxSample, "\n")
	newLines := strings.Split(updated, "\n")
	if len(oldLines) != len(newLines) {
		t.Fatalf("Zeilenanzahl geändert: %d -> %d", len(oldLines), len(newLines))
	}

	var differing []string
	for i := range oldLines {
		if oldLines[i] != newLines[i] {
			differing = append(differing, newLines[i])
		}
	}

	if len(differing) != 4 {
		t.Fatalf("erwartet 4 geänderte Zeilen, bekommen %d: %v", len(differing), differing)
	}

	want := []string{
		"    Zombies = 1,",
		"    EnableSnowOnGround = false,",
		`    LootItemRemovalList = "Base.Vest",`,
		"        AllowMiniMap = true,",
	}
	for index, line := range want {
		if differing[index] != line {
			t.Errorf("geänderte Zeile %d = %q, erwartet %q", index, differing[index], line)
		}
	}

	// Der Kommentar hinter DayLength überlebt unangetastet.
	if !strings.Contains(updated, "DayLength = 3,          -- 3 = 1 hour") {
		t.Error("der Kommentar hinter DayLength wurde verändert")
	}
}

// Ein Wert, der nicht zum Typ passt, erreicht die Datei nicht.
func TestWriteSandboxValuesRefusesAValueOfTheWrongKind(t *testing.T) {
	vars, err := parseSandboxLua(sandboxSample)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}

	if _, _, err := writeSandboxValues(sandboxSample, vars, map[string]string{"Zombies": "viele"}); err == nil {
		t.Error("erwartet: Fehler wegen unpassendem Wert")
	}

	// Ein Name, den die Datei nicht kennt, wird nicht eingefügt.
	if _, _, err := writeSandboxValues(sandboxSample, vars, map[string]string{"ErfundeneVariable": "1"}); err == nil {
		t.Error("erwartet: Fehler wegen unbekannter Variable")
	}
}

// Der Zugangsweg wird ausdrücklich eingerichtet. Ohne Einrichtung wird nichts
// geraten und keine Datei angefasst (PG-20).
func TestResolveSandboxPath(t *testing.T) {
	// Eine gefälschte Verzeichnisauflistung — der Test fasst keine echte Datei an.
	listing := map[string][]string{
		"/srv/zomboid/Server": {"servertest.ini", "servertest_SandboxVars.lua", "servertest_spawnpoints.lua"},
		"/srv/zomboid":        {"Server", "Saves"},
		"/data/Server":        {"myserver_SandboxVars.lua", "otherserver_SandboxVars.lua"},
		"/leer/Server":        {"servertest.ini"},
	}
	list := func(dir string) ([]string, error) {
		entries, ok := listing[dir]
		if !ok {
			return nil, errors.New("no such directory")
		}
		return entries, nil
	}

	cases := []struct {
		name       string
		mode       string
		configured string
		want       string
		wantErr    bool
	}{
		{"nicht eingerichtet", "", "", "", true},
		{"Modus ohne Pfad", "local", "", "", true},
		{"lokale Datei direkt benannt", "local", "/srv/zomboid/Server/servertest_SandboxVars.lua", "/srv/zomboid/Server/servertest_SandboxVars.lua", false},
		{"lokales Serververzeichnis", "local", "/srv/zomboid/Server", "/srv/zomboid/Server/servertest_SandboxVars.lua", false},
		{"Docker-Volume", "docker", "/srv/zomboid", "/srv/zomboid/Server/servertest_SandboxVars.lua", false},
		{"Verzeichnis ohne Sandbox-Datei", "local", "/leer/Server", "", true},
		{"unbekanntes Verzeichnis", "local", "/gibt/es/nicht", "", true},
		{"unbekannter Modus", "sftp", "/srv/zomboid", "", true},
	}

	for _, testCase := range cases {
		got, err := resolveSandboxPath(testCase.mode, testCase.configured, list)
		if testCase.wantErr {
			if err == nil {
				t.Errorf("%s: erwartet eine Fehlermeldung, bekommen %q", testCase.name, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unerwarteter Fehler: %v", testCase.name, err)
			continue
		}
		if got != testCase.want {
			t.Errorf("%s: Pfad = %q, erwartet %q", testCase.name, got, testCase.want)
		}
	}
}

// Mehrere Serverdateien nebeneinander: Die Anwendung rät nicht, sondern sagt es.
func TestResolveSandboxPathRefusesToGuessBetweenServers(t *testing.T) {
	list := func(dir string) ([]string, error) {
		return []string{"myserver_SandboxVars.lua", "otherserver_SandboxVars.lua"}, nil
	}

	if _, err := resolveSandboxPath("local", "/data/Server", list); err == nil {
		t.Error("erwartet: Fehler, weil zwei Sandbox-Dateien nebeneinander liegen")
	}
}
