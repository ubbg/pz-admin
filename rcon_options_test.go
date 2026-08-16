package main

import "testing"

func serverStateForTest() []Option {
	return parseOptionLines([]string{
		"* PublicName=My B42 Server",
		"* MaxPlayers=32",
		"* SpeedLimit=70.0",
		"* War=false",
		"* AntiCheatSafety=1",
	})
}

// Geschrieben wird nur, was sich gegenüber dem gelesenen Serverstand geändert hat —
// und nur für Namen, die der Server selbst gemeldet hat (PG-03).
func TestDiffServerOptionsOnlySendsChangedNamesTheServerReported(t *testing.T) {
	server := serverStateForTest()

	changes := diffServerOptions(server, map[string]string{
		"PublicName":          "My B42 Server", // unverändert
		"MaxPlayers":          "48",            // geändert
		"SpeedLimit":          "70.0",          // unverändert
		"War":                 "true",          // geändert
		"AntiCheatSafety":     "1",             // unverändert
		"HoursForLootRespawn": "12",            // vom Server nie gemeldet -> nie senden
	})

	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d: %v", len(changes), changes)
	}

	// Reihenfolge folgt der Serverantwort, damit der Ablauf nachvollziehbar bleibt.
	if changes[0].Name != "MaxPlayers" || changes[0].Value != "48" {
		t.Errorf("first change = %v, want MaxPlayers=48", changes[0])
	}
	if changes[1].Name != "War" || changes[1].Value != "true" {
		t.Errorf("second change = %v, want War=true", changes[1])
	}
}

// Der Server meldet Kommazahlen als "70.0", die Oberfläche schickt vielleicht "70".
// Das ist derselbe Wert und darf kein changeoption auslösen.
func TestDiffServerOptionsTreatsEqualNumbersAndCasedBooleansAsUnchanged(t *testing.T) {
	server := serverStateForTest()

	changes := diffServerOptions(server, map[string]string{
		"SpeedLimit": "70",
		"MaxPlayers": "32.0",
		"War":        "FALSE",
	})

	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %v", changes)
	}
}

// Erfolg wird an der Antwort geprüft, nicht daran, dass Execute keinen Fehler
// lieferte — ein PZ-Server antwortet auf Unsinn freundlich.
func TestOptionUpdateSucceededChecksTheServerAnswer(t *testing.T) {
	cases := []struct {
		name     string
		option   OptionPair
		kind     string
		response string
		want     bool
	}{
		{
			name:     "wörtliche Bestätigung",
			option:   OptionPair{Name: "MaxPlayers", Value: "48"},
			kind:     "Integer",
			response: "Option : MaxPlayers is now : 48",
			want:     true,
		},
		{
			name:     "Kommazahl wird vom Server anders formatiert",
			option:   OptionPair{Name: "SpeedLimit", Value: "70"},
			kind:     "Double",
			response: "Option : SpeedLimit is now : 70.0",
			want:     true,
		},
		{
			name:     "Text mit Sonderzeichen kommt unverändert zurück",
			option:   OptionPair{Name: "ClientCommandFilter", Value: "-vehicle.*;+vehicle.fixPart"},
			kind:     "String",
			response: "Option : ClientCommandFilter is now : -vehicle.*;+vehicle.fixPart",
			want:     true,
		},
		{
			name:     "Server bestätigt einen anderen Wert",
			option:   OptionPair{Name: "MaxPlayers", Value: "48"},
			kind:     "Integer",
			response: "Option : MaxPlayers is now : 32",
			want:     false,
		},
		{
			name:     "freundliche Antwort auf eine unbekannte Option",
			option:   OptionPair{Name: "KickFastPlayers", Value: "true"},
			kind:     "Boolean",
			response: "Change a server option. Use /changeoption \"optionName\" \"newValue\"",
			want:     false,
		},
	}

	for _, testCase := range cases {
		got := optionUpdateSucceeded(testCase.option, testCase.kind, testCase.response)
		if got != testCase.want {
			t.Errorf("%s: optionUpdateSucceeded = %v, want %v", testCase.name, got, testCase.want)
		}
	}
}

// Die Antwort von showoptions, wie ein Server sie liefert: eine Kopfzeile, dann je
// Option eine Zeile "* Name=Wert". Der Parser übernimmt jede davon, ohne vorher zu
// wissen, welche Optionen es gibt (PG-01).
func TestParseOptionLinesTakesEveryOptionTheServerReports(t *testing.T) {
	lines := []string{
		"List of Server Options:",
		"* PublicName=My B42 Server",
		"* War=true",
		"* AnnounceAnimalDeath=false",
		"* AntiCheatSafety=1",
	}

	options := parseOptionLines(lines)

	if len(options) != 4 {
		t.Fatalf("expected 4 options, got %d: %v", len(options), options)
	}

	byName := map[string]Option{}
	for _, option := range options {
		byName[option.Name] = option
	}

	for _, name := range []string{"PublicName", "War", "AnnounceAnimalDeath", "AntiCheatSafety"} {
		if _, ok := byName[name]; !ok {
			t.Errorf("option %s missing from parsed options", name)
		}
	}

	if got := byName["PublicName"].Value; got != "My B42 Server" {
		t.Errorf("PublicName value = %q, want %q", got, "My B42 Server")
	}
	if got := byName["AntiCheatSafety"].Value; got != "1" {
		t.Errorf("AntiCheatSafety value = %q, want %q", got, "1")
	}
}

// Der Typ kommt aus dem Wert, nicht aus einer Liste im Code. Was sich keiner Klasse
// zuordnen lässt, ist Text und wird als Textfeld bearbeitet (PG-05).
func TestParseOptionLinesDerivesKindFromValue(t *testing.T) {
	cases := []struct {
		line string
		name string
		kind string
	}{
		{"* War=true", "War", "Boolean"},
		{"* AnnounceAnimalDeath=false", "AnnounceAnimalDeath", "Boolean"},
		{"* PVP=TRUE", "PVP", "Boolean"},
		{"* MaxPlayers=32", "MaxPlayers", "Integer"},
		{"* WarStartDelay=-1", "WarStartDelay", "Integer"},
		{"* SpeedLimit=70.0", "SpeedLimit", "Double"},
		{"* PVPMeleeDamageModifier=30.5", "PVPMeleeDamageModifier", "Double"},
		{"* PublicName=My Server", "PublicName", "String"},
		{"* SpawnPoint=0,0,0", "SpawnPoint", "String"},
		{"* Mods=", "Mods", "String"},
		{"* ClientCommandFilter=-vehicle.*;+vehicle.fixPart", "ClientCommandFilter", "String"},
	}

	for _, testCase := range cases {
		options := parseOptionLines([]string{testCase.line})
		if len(options) != 1 {
			t.Fatalf("%q: expected 1 option, got %d", testCase.line, len(options))
		}
		if options[0].Name != testCase.name {
			t.Errorf("%q: name = %q, want %q", testCase.line, options[0].Name, testCase.name)
		}
		if options[0].Kind != testCase.kind {
			t.Errorf("%q: kind = %q, want %q", testCase.line, options[0].Kind, testCase.kind)
		}
	}
}

// Build 42 erwartet Mod-IDs mit führendem Backslash in der Serverkonfiguration
// (\mod1;\mod2). Die Nutzerin sieht und tippt weiterhin nur die IDs.
func TestModsValueForServerAddsThePrefixExactlyOnce(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  string
	}{
		{"IDs ohne Präfix", "mod1;mod2", `\mod1;\mod2`},
		{"schon im B42-Format", `\mod1;\mod2`, `\mod1;\mod2`},
		{"gemischt", `\mod1;mod2`, `\mod1;\mod2`},
		{"mit Leerzeichen", " mod1 ; mod2 ", `\mod1;\mod2`},
		{"leere Liste bleibt leer", "", ""},
		{"nur Trenner", ";;", ""},
		{"einzelne Mod", "mod1", `\mod1`},
	}

	for _, testCase := range cases {
		if got := modsValueForServer(testCase.value); got != testCase.want {
			t.Errorf("%s: modsValueForServer(%q) = %q, want %q", testCase.name, testCase.value, got, testCase.want)
		}
	}
}

func TestModsValueForDisplayDropsThePrefix(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{`\mod1;\mod2`, "mod1;mod2"},
		{"mod1;mod2", "mod1;mod2"},
		{"", ""},
		{`\mod1`, "mod1"},
	}

	for _, testCase := range cases {
		if got := modsValueForDisplay(testCase.value); got != testCase.want {
			t.Errorf("modsValueForDisplay(%q) = %q, want %q", testCase.value, got, testCase.want)
		}
	}
}

// Nur Mods trägt das Präfix — das Workshop-Feld bleibt eine Liste blanker IDs.
func TestOnlyModsIsRewritten(t *testing.T) {
	if got := optionValueForServer("WorkshopItems", "2200148440;2392709985"); got != "2200148440;2392709985" {
		t.Errorf("WorkshopItems wurde verändert: %q", got)
	}
	if got := optionValueForDisplay("WorkshopItems", "2200148440;2392709985"); got != "2200148440;2392709985" {
		t.Errorf("WorkshopItems wurde verändert: %q", got)
	}
	if got := optionValueForServer("Mods", "mod1"); got != `\mod1` {
		t.Errorf("Mods wurde nicht umgeschrieben: %q", got)
	}
	if got := optionValueForDisplay("Mods", `\mod1`); got != "mod1" {
		t.Errorf("Mods wurde nicht umgeschrieben: %q", got)
	}
}
