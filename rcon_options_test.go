package main

import "testing"

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
