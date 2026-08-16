package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type OptionPair struct {
	Name  string
	Value string
}

// Option ist eine Serveroption, so wie showoptions sie meldet: Name, Wert und der aus
// dem Wert abgeleitete Typ. Welche Optionen es gibt, entscheidet damit der Server —
// nicht die Anwendung. Siehe specs/options/OPTIONS-001.spec.md.
type Option struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Kind  string `json:"kind"` // Boolean, Integer, Double, String
}

// parseOptionLines liest die Antwort von showoptions Zeile für Zeile. Jede Zeile der
// Form "* Name=Wert" wird übernommen, auch wenn die Anwendung die Option nicht kennt
// (PG-01). Alles andere wird übergangen.
func parseOptionLines(lines []string) []Option {
	options := make([]Option, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "* ") {
			continue
		}

		parts := strings.SplitN(strings.TrimPrefix(line, "* "), "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue
		}
		value := optionValueForDisplay(name, strings.TrimSpace(parts[1]))

		options = append(options, Option{Name: name, Value: value, Kind: kindOfValue(value)})
	}

	return options
}

// optionUpdateSucceeded prüft die Antwort des Servers auf changeoption. Build 42
// antwortet "Option : <Name> is now : <Wert>"; bei Kommazahlen formatiert der Server
// den Wert um ("70" wird zu "70.0"), deshalb wird dort numerisch verglichen.
func optionUpdateSucceeded(option OptionPair, kind string, response string) bool {
	prefix := fmt.Sprintf("Option : %s is now : ", option.Name)
	if !strings.HasPrefix(response, prefix) {
		return false
	}

	confirmed := strings.TrimPrefix(response, prefix)
	if confirmed == option.Value {
		return true
	}

	if kind == "Double" || kind == "Integer" {
		return optionValuesEqual(confirmed, option.Value)
	}

	return false
}

// diffServerOptions vergleicht die Eingaben der Oberfläche mit dem zuletzt gelesenen
// Serverstand. Es entsteht ein Änderungsbefehl je Option, die sich geändert hat — und
// keiner für einen Namen, den der Server nicht gemeldet hat (PG-03). Die Reihenfolge
// folgt der Serverantwort.
func diffServerOptions(server []Option, values map[string]string) []OptionPair {
	var changes []OptionPair

	for _, option := range server {
		value, ok := values[option.Name]
		if !ok {
			continue
		}
		if optionValuesEqual(option.Value, value) {
			continue
		}
		changes = append(changes, OptionPair{Name: option.Name, Value: value})
	}

	return changes
}

// optionValuesEqual vergleicht zwei Werte so, wie der Server sie versteht: "70" und
// "70.0" sind dieselbe Zahl, "false" und "FALSE" derselbe Wahrheitswert.
func optionValuesEqual(a string, b string) bool {
	if a == b {
		return true
	}

	if isBooleanValue(a) && isBooleanValue(b) {
		return strings.EqualFold(a, b)
	}

	numberA, errA := strconv.ParseFloat(a, 64)
	numberB, errB := strconv.ParseFloat(b, 64)
	if errA == nil && errB == nil {
		return numberA == numberB
	}

	return false
}

func isBooleanValue(value string) bool {
	return strings.EqualFold(value, "true") || strings.EqualFold(value, "false")
}

// Build 42 führt Mod-IDs in der Serverkonfiguration mit einem vorangestellten
// Backslash: Mods=\mod1;\mod2. Ohne ihn lädt der Server die Mods nicht. Die
// Umwandlung sitzt hier und nur hier — die Oberfläche zeigt und nimmt blanke IDs.
func optionValueForServer(name string, value string) string {
	if name == "Mods" {
		return modsValueForServer(value)
	}
	return value
}

func optionValueForDisplay(name string, value string) string {
	if name == "Mods" {
		return modsValueForDisplay(value)
	}
	return value
}

func modsValueForServer(value string) string {
	ids := modIds(value)
	for i, id := range ids {
		ids[i] = `\` + id
	}
	return strings.Join(ids, ";")
}

func modsValueForDisplay(value string) string {
	return strings.Join(modIds(value), ";")
}

// modIds zerlegt eine Mod-Liste in ihre IDs — ohne Präfix, ohne Leerraum und ohne
// leere Einträge, damit aus einer leeren Liste kein einzelner Backslash wird.
func modIds(value string) []string {
	var ids []string
	for _, id := range strings.Split(value, ";") {
		id = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(id), `\`))
		if id == "" {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

// kindOfValue leitet den Typ aus dem Wert ab: true/false ist ein Wahrheitswert, eine
// reine Ganzzahl eine Zahl, eine Zahl mit Punkt eine Kommazahl, alles Übrige Text.
func kindOfValue(value string) string {
	switch strings.ToLower(value) {
	case "true", "false":
		return "Boolean"
	}

	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return "Integer"
	}

	if strings.Contains(value, ".") {
		if _, err := strconv.ParseFloat(value, 64); err == nil {
			return "Double"
		}
	}

	return "String"
}

var (
	// serverOptions ist der zuletzt gelesene Serverstand: alles, was showoptions
	// gemeldet hat. Ein Optionsname steht nirgends sonst im Go-Code (PG-02).
	serverOptions   []Option
	lastOptionsHash string
)

func (app *App) ServerOptions() []Option {
	return serverOptions
}

func hashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func options_update() error {
	res, err := conn.Execute("showoptions")
	if err != nil {
		return fmt.Errorf("error getting options: %v", err)
	}

	currentHash := hashString(res)
	if currentHash == lastOptionsHash {
		runtime.LogTracef(app.ctx, "Options unchanged, skipping sync")
		return nil
	}

	lastOptionsHash = currentHash
	serverOptions = parseOptionLines(strings.Split(res, "\n"))

	runtime.EventsEmit(app.ctx, "update-options-list", serverOptions)
	runtime.LogDebugf(app.ctx, "Options synced: %d options reported by the server", len(serverOptions))

	return nil
}

// UpdateOptions schreibt die Eingaben der Oberfläche zurück: ein changeoption je
// Option, die sich gegenüber dem gelesenen Serverstand geändert hat, danach optional
// ein reloadoptions. Namen, die der Server nicht gemeldet hat, werden nie gesendet.
func (app *App) UpdateOptions(values map[string]string, reloadOptions bool) bool {
	defer options_update()

	optionsToUpdate := diffServerOptions(serverOptions, values)

	if len(optionsToUpdate) == 0 {
		app.SendNotification(Notification{Title: "rcon.no_options_to_update", Variant: "warning"})
		return false
	}

	failed := app.applyServerOptions(optionsToUpdate)

	if len(failed) > 0 {
		// Welche Option der Server abgelehnt hat, ist beim Wechsel des Spiel-Builds
		// die einzige interessante Angabe — eine Anzahl allein hilft niemandem (PG-06).
		succeeded := len(optionsToUpdate) - len(failed)
		notification := Notification{
			Title:   "rcon.failed_to_update_options",
			Variant: "error",
			Parameters: map[string]string{
				"f":       fmt.Sprintf("%d", len(failed)),
				"s":       fmt.Sprintf("%d", succeeded),
				"options": strings.Join(failed, ", "),
			},
		}
		if succeeded > 0 {
			notification.Title = "rcon.options_partially_updated"
			notification.Variant = "warning"
		}
		app.SendNotification(notification)
		return false
	}

	if err := options_update(); err != nil {
		runtime.LogErrorf(app.ctx, "Error syncing options after update: %v", err)
		app.SendNotification(Notification{Title: "rcon.options_updated_sync_failed", Variant: "error"})
		return false
	}

	if reloadOptions {
		command := RCONCommand{
			CommandTemplate: "reloadoptions",
			SuccessCheck: func(name string, response string) bool {
				return response == "Options reloaded"
			},
		}

		if command.execute() != 1 {
			app.SendNotification(Notification{Title: "rcon.reloadOptions.single_fail", Variant: "error"})
			return false
		}

		app.SendNotification(Notification{Title: "rcon.options_saved_and_reloaded", Variant: "success"})
		return true
	}

	app.SendNotification(Notification{Title: "rcon.options_updated", Variant: "success"})
	return true
}

// applyServerOptions sendet je Option ein changeoption und liefert die Namen, die der
// Server nicht übernommen hat.
func (app *App) applyServerOptions(options []OptionPair) []string {
	connMutex.Lock()
	defer connMutex.Unlock()

	defer runtime.EventsEmit(app.ctx, "setProgress", 0)
	runtime.EventsEmit(app.ctx, "setProgress", 10)

	kinds := map[string]string{}
	for _, option := range serverOptions {
		kinds[option.Name] = option.Kind
	}

	var failed []string

	for i, option := range options {
		runtime.EventsEmit(app.ctx, "setProgress", float64(i)/float64(len(options))*100)

		// Der Server bekommt den Wert in seiner Schreibweise (Mods mit Backslash),
		// die Oberfläche hat ihn in der lesbaren geschickt.
		sent := OptionPair{Name: option.Name, Value: optionValueForServer(option.Name, option.Value)}

		command := fmt.Sprintf("changeoption %s \"%s\"", sent.Name, sent.Value)
		res, err := conn.Execute(command)

		if err != nil || !optionUpdateSucceeded(sent, kinds[option.Name], res) {
			runtime.LogErrorf(app.ctx, "Failed to update %s: %v (response: %s)", option.Name, err, res)
			failed = append(failed, option.Name)
		}
	}

	return failed
}
