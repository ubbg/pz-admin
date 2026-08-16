package main

import "testing"

// Die Befehlszeile entsteht aus Vorlage, Argumenten und Spielername — nicht per
// Sprintf an der Aufrufstelle. Spielernamen stehen in Anführungszeichen, weil ein
// Name mit Leerzeichen der Regelfall ist.
func TestBuildCommandFillsTemplate(t *testing.T) {
	cases := []struct {
		name       string
		template   string
		args       []RCONCommandParam
		playerName string
		want       string
	}{
		{
			name:       "Spielername wird eingesetzt und gequotet",
			template:   "kick {name} {reason}",
			args:       []RCONCommandParam{{Name: "reason", Key: "-r", Value: `"AFK"`}},
			playerName: "John Doe",
			want:       `kick "John Doe" -r "AFK"`,
		},
		{
			name:       "fehlendes optionales Argument hinterlässt keine Lücke",
			template:   "kick {name} {reason}",
			args:       []RCONCommandParam{{Name: "reason", Key: "-r", Value: nil}},
			playerName: "rj",
			want:       `kick "rj"`,
		},
		{
			name:       "Argument ohne Schlüssel steht für sich",
			template:   "godmodeplayer {name} {value}",
			args:       []RCONCommandParam{{Name: "value", Value: "-true", Mandatory: true}},
			playerName: "rj",
			want:       `godmodeplayer "rj" -true`,
		},
		{
			name:     "Befehl ohne Spielerbezug",
			template: "worldgen {action}",
			args:     []RCONCommandParam{{Name: "action", Value: "status", Mandatory: true}},
			want:     "worldgen status",
		},
	}

	for _, testCase := range cases {
		got, err := buildCommand(testCase.template, testCase.args, testCase.playerName)
		if err != nil {
			t.Errorf("%s: unerwarteter Fehler: %v", testCase.name, err)
			continue
		}
		if got != testCase.want {
			t.Errorf("%s: buildCommand = %q, want %q", testCase.name, got, testCase.want)
		}
	}
}

// Ein fehlendes Pflichtargument ist ein Fehler, kein halber Befehl.
func TestBuildCommandRefusesMissingMandatoryArgument(t *testing.T) {
	_, err := buildCommand("godmodeplayer {name} {value}", []RCONCommandParam{
		{Name: "value", Value: nil, Mandatory: true},
	}, "rj")

	if err == nil {
		t.Fatal("erwartet: Fehler wegen fehlendem Pflichtargument")
	}
}

// Für die meisten Build-42-Befehle ist der Erfolgstext nicht dokumentiert. Belegt ist
// dagegen, dass der Server bei falscher Benutzung seinen eigenen Hilfetext
// zurückwirft — genau daran wird Misserfolg erkannt.
func TestSucceededUnlessUsage(t *testing.T) {
	cases := []struct {
		name     string
		response string
		want     bool
	}{
		// Hilfetexte im Wortlaut der pzwiki-Seite „Admin commands" (42.20.2)
		{"Hilfetext banip", "Ban IP. Use /banip IP", false},
		{"Hilfetext addsteamid", `Use this command to add a SteamID to a list of allowed SteamIDs on server. Use: /addSteamID "steamid"`, false},
		{"Hilfetext setpassword", `Use this command to change password for a user. Use: /setpassword "username" "newpassword"`, false},
		{"leere Antwort", "", false},
		{"nur Leerraum", "   \n", false},
		{"Quittung des Servers", "SteamID 76561198000000000 added", true},
		{"Quittung ohne Schlüsselwort", "IP 203.0.113.5 is now banned", true},
	}

	for _, testCase := range cases {
		if got := succeededUnlessUsage("", testCase.response); got != testCase.want {
			t.Errorf("%s: succeededUnlessUsage(%q) = %v, want %v", testCase.name, testCase.response, got, testCase.want)
		}
	}
}
