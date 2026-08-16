package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PzOptions struct {
	AdminSafehouse                               bool    `json:"AdminSafehouse"`
	AllowCoop                                    bool    `json:"AllowCoop"`
	AllowDestructionBySledgehammer               bool    `json:"AllowDestructionBySledgehammer"`
	AllowNonAsciiUsername                        bool    `json:"AllowNonAsciiUsername"`
	AnnounceDeath                                bool    `json:"AnnounceDeath"`
	AntiCheatProtectionType1                     bool    `json:"AntiCheatProtectionType1"`
	AntiCheatProtectionType2                     bool    `json:"AntiCheatProtectionType2"`
	AntiCheatProtectionType3                     bool    `json:"AntiCheatProtectionType3"`
	AntiCheatProtectionType4                     bool    `json:"AntiCheatProtectionType4"`
	AntiCheatProtectionType5                     bool    `json:"AntiCheatProtectionType5"`
	AntiCheatProtectionType6                     bool    `json:"AntiCheatProtectionType6"`
	AntiCheatProtectionType7                     bool    `json:"AntiCheatProtectionType7"`
	AntiCheatProtectionType8                     bool    `json:"AntiCheatProtectionType8"`
	AntiCheatProtectionType9                     bool    `json:"AntiCheatProtectionType9"`
	AntiCheatProtectionType10                    bool    `json:"AntiCheatProtectionType10"`
	AntiCheatProtectionType11                    bool    `json:"AntiCheatProtectionType11"`
	AntiCheatProtectionType12                    bool    `json:"AntiCheatProtectionType12"`
	AntiCheatProtectionType13                    bool    `json:"AntiCheatProtectionType13"`
	AntiCheatProtectionType14                    bool    `json:"AntiCheatProtectionType14"`
	AntiCheatProtectionType15                    bool    `json:"AntiCheatProtectionType15"`
	AntiCheatProtectionType16                    bool    `json:"AntiCheatProtectionType16"`
	AntiCheatProtectionType17                    bool    `json:"AntiCheatProtectionType17"`
	AntiCheatProtectionType18                    bool    `json:"AntiCheatProtectionType18"`
	AntiCheatProtectionType19                    bool    `json:"AntiCheatProtectionType19"`
	AntiCheatProtectionType20                    bool    `json:"AntiCheatProtectionType20"`
	AntiCheatProtectionType21                    bool    `json:"AntiCheatProtectionType21"`
	AntiCheatProtectionType22                    bool    `json:"AntiCheatProtectionType22"`
	AntiCheatProtectionType23                    bool    `json:"AntiCheatProtectionType23"`
	AntiCheatProtectionType24                    bool    `json:"AntiCheatProtectionType24"`
	AntiCheatProtectionType2ThresholdMultiplier  float64 `json:"AntiCheatProtectionType2ThresholdMultiplier"`
	AntiCheatProtectionType3ThresholdMultiplier  float64 `json:"AntiCheatProtectionType3ThresholdMultiplier"`
	AntiCheatProtectionType4ThresholdMultiplier  float64 `json:"AntiCheatProtectionType4ThresholdMultiplier"`
	AntiCheatProtectionType9ThresholdMultiplier  float64 `json:"AntiCheatProtectionType9ThresholdMultiplier"`
	AntiCheatProtectionType15ThresholdMultiplier float64 `json:"AntiCheatProtectionType15ThresholdMultiplier"`
	AntiCheatProtectionType20ThresholdMultiplier float64 `json:"AntiCheatProtectionType20ThresholdMultiplier"`
	AntiCheatProtectionType22ThresholdMultiplier float64 `json:"AntiCheatProtectionType22ThresholdMultiplier"`
	AntiCheatProtectionType24ThresholdMultiplier float64 `json:"AntiCheatProtectionType24ThresholdMultiplier"`
	AutoCreateUserInWhiteList                    bool    `json:"AutoCreateUserInWhiteList"`
	BackupsCount                                 int     `json:"BackupsCount"`
	BackupsOnStart                               bool    `json:"BackupsOnStart"`
	BackupsOnVersionChange                       bool    `json:"BackupsOnVersionChange"`
	BackupsPeriod                                int     `json:"BackupsPeriod"`
	BanKickGlobalSound                           bool    `json:"BanKickGlobalSound"`
	BloodSplatLifespanDays                       int     `json:"BloodSplatLifespanDays"`
	CarEngineAttractionModifier                  float64 `json:"CarEngineAttractionModifier"`
	ChatStreams                                  string  `json:"ChatStreams"`
	ClientActionLogs                             string  `json:"ClientActionLogs"`
	ClientCommandFilter                          string  `json:"ClientCommandFilter"`
	ConstructionPreventsLootRespawn              bool    `json:"ConstructionPreventsLootRespawn"`
	DefaultPort                                  int     `json:"DefaultPort"`
	DenyLoginOnOverloadedServer                  bool    `json:"DenyLoginOnOverloadedServer"`
	DisableRadioAdmin                            bool    `json:"DisableRadioAdmin"`
	DisableRadioGM                               bool    `json:"DisableRadioGM"`
	DisableRadioInvisible                        bool    `json:"DisableRadioInvisible"`
	DisableRadioModerator                        bool    `json:"DisableRadioModerator"`
	DisableRadioOverseer                         bool    `json:"DisableRadioOverseer"`
	DisableRadioStaff                            bool    `json:"DisableRadioStaff"`
	DisableSafehouseWhenPlayerConnected          bool    `json:"DisableSafehouseWhenPlayerConnected"`
	DiscordEnable                                bool    `json:"DiscordEnable"`
	DiscordToken                                 string  `json:"DiscordToken"`
	DiscordChannel                               string  `json:"DiscordChannel"`
	DiscordChannelID                             string  `json:"DiscordChannelID"`
	DisplayUserName                              bool    `json:"DisplayUserName"`
	DoLuaChecksum                                bool    `json:"DoLuaChecksum"`
	DropOffWhiteListAfterDeath                   bool    `json:"DropOffWhiteListAfterDeath"`
	Faction                                      bool    `json:"Faction"`
	FactionDaySurvivedToCreate                   int     `json:"FactionDaySurvivedToCreate"`
	FactionPlayersRequiredForTag                 int     `json:"FactionPlayersRequiredForTag"`
	FastForwardMultiplier                        float64 `json:"FastForwardMultiplier"`
	GlobalChat                                   bool    `json:"GlobalChat"`
	HidePlayersBehindYou                         bool    `json:"HidePlayersBehindYou"`
	HoursForLootRespawn                          int     `json:"HoursForLootRespawn"`
	ItemNumbersLimitPerContainer                 int     `json:"ItemNumbersLimitPerContainer"`
	KickFastPlayers                              bool    `json:"KickFastPlayers"`
	KnockedDownAllowed                           bool    `json:"KnockedDownAllowed"`
	LoginQueueConnectTimeout                     int     `json:"LoginQueueConnectTimeout"`
	LoginQueueEnabled                            bool    `json:"LoginQueueEnabled"`
	Map                                          string  `json:"Map"`
	MapRemotePlayerVisibility                    int     `json:"MapRemotePlayerVisibility"`
	MaxAccountsPerUser                           int     `json:"MaxAccountsPerUser"`
	MaxItemsForLootRespawn                       int     `json:"MaxItemsForLootRespawn"`
	MaxPlayers                                   int     `json:"MaxPlayers"`
	MinutesPerPage                               float64 `json:"MinutesPerPage"`
	Mods                                         string  `json:"Mods"`
	MouseOverToSeeDisplayName                    bool    `json:"MouseOverToSeeDisplayName"`
	NoFire                                       bool    `json:"NoFire"`
	Open                                         bool    `json:"Open"`
	PVP                                          bool    `json:"PVP"`
	PVPFirearmDamageModifier                     float64 `json:"PVPFirearmDamageModifier"`
	PVPMeleeDamageModifier                       float64 `json:"PVPMeleeDamageModifier"`
	PVPMeleeWhileHitReaction                     bool    `json:"PVPMeleeWhileHitReaction"`
	PauseEmpty                                   bool    `json:"PauseEmpty"`
	PerkLogs                                     bool    `json:"PerkLogs"`
	PingLimit                                    int     `json:"PingLimit"`
	PlayerBumpPlayer                             bool    `json:"PlayerBumpPlayer"`
	PlayerRespawnWithOther                       bool    `json:"PlayerRespawnWithOther"`
	PlayerRespawnWithSelf                        bool    `json:"PlayerRespawnWithSelf"`
	PlayerSafehouse                              bool    `json:"PlayerSafehouse"`
	Public                                       bool    `json:"Public"`
	PublicDescription                            string  `json:"PublicDescription"`
	PublicName                                   string  `json:"PublicName"`
	RemovePlayerCorpsesOnCorpseRemoval           bool    `json:"RemovePlayerCorpsesOnCorpseRemoval"`
	ResetID                                      int     `json:"ResetID"`
	SafeHouseRemovalTime                         int     `json:"SafeHouseRemovalTime"`
	SafehouseAllowFire                           bool    `json:"SafehouseAllowFire"`
	SafehouseAllowLoot                           bool    `json:"SafehouseAllowLoot"`
	SafehouseAllowNonResidential                 bool    `json:"SafehouseAllowNonResidential"`
	SafehouseAllowRespawn                        bool    `json:"SafehouseAllowRespawn"`
	SafehouseAllowTrepass                        bool    `json:"SafehouseAllowTrepass"`
	SafehouseDaySurvivedToClaim                  int     `json:"SafehouseDaySurvivedToClaim"`
	SafetyCooldownTimer                          int     `json:"SafetyCooldownTimer"`
	SafetySystem                                 bool    `json:"SafetySystem"`
	SafetyToggleTimer                            int     `json:"SafetyToggleTimer"`
	SaveWorldEveryMinutes                        int     `json:"SaveWorldEveryMinutes"`
	ServerPlayerID                               int     `json:"ServerPlayerID"`
	ShowFirstAndLastName                         bool    `json:"ShowFirstAndLastName"`
	ShowSafety                                   bool    `json:"ShowSafety"`
	SledgehammerOnlyInSafehouse                  bool    `json:"SledgehammerOnlyInSafehouse"`
	SleepAllowed                                 bool    `json:"SleepAllowed"`
	SleepNeeded                                  bool    `json:"SleepNeeded"`
	SneakModeHideFromOtherPlayers                bool    `json:"SneakModeHideFromOtherPlayers"`
	SpawnItems                                   string  `json:"SpawnItems"`
	SpawnPoint                                   string  `json:"SpawnPoint"`
	SpeedLimit                                   float64 `json:"SpeedLimit"`
	SteamScoreboard                              bool    `json:"SteamScoreboard"`
	SteamVAC                                     bool    `json:"SteamVAC"`
	TrashDeleteAll                               bool    `json:"TrashDeleteAll"`
	UDPPort                                      int     `json:"UDPPort"`
	UPnP                                         bool    `json:"UPnP"`
	Voice3D                                      bool    `json:"Voice3D"`
	VoiceEnable                                  bool    `json:"VoiceEnable"`
	VoiceMaxDistance                             float64 `json:"VoiceMaxDistance"`
	VoiceMinDistance                             float64 `json:"VoiceMinDistance"`
	WorkshopItems                                string  `json:"WorkshopItems"`
	ServerWelcomeMessage                         string  `json:"ServerWelcomeMessage"`
}

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
		value := strings.TrimSpace(parts[1])

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
	pzOptions PzOptions
	// serverOptions ist der zuletzt gelesene Serverstand in generischer Form: alles,
	// was showoptions gemeldet hat, unabhängig davon, ob die Anwendung es kennt.
	serverOptions   []Option
	lastOptionsHash string
)

func (app *App) PzOptionsList() []Option {
	return serverOptions
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Bool:
		field.SetBool(strings.ToLower(value) == "true")
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.String:
		field.SetString(value)
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func hashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func pzOptions_update() error {
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
	lines := strings.Split(res, "\n")
	updatedOptions := PzOptions{}

	if err := parseOptions(lines, &updatedOptions); err != nil {
		return fmt.Errorf("error parsing options: %v", err)
	}

	pzOptions = updatedOptions
	serverOptions = parseOptionLines(lines)

	runtime.EventsEmit(app.ctx, "update-options", pzOptions)
	runtime.EventsEmit(app.ctx, "update-options-list", serverOptions)
	runtime.LogDebugf(app.ctx, "Options synced: %d options reported by the server", len(serverOptions))

	return nil
}

func parseOptions(lines []string, target *PzOptions) error {
	v := reflect.ValueOf(target).Elem()

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "* ") {
			continue
		}

		parts := strings.SplitN(strings.TrimPrefix(line, "* "), "=", 2)
		if len(parts) != 2 {
			runtime.LogDebugf(app.ctx, "Invalid option: %s", line)
			continue
		}

		fieldName, fieldValue := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		field := v.FieldByName(fieldName)

		if !field.IsValid() {
			runtime.LogDebugf(app.ctx, "Unknown option: %s", fieldName)
			continue
		}

		if err := setFieldValue(field, fieldValue); err != nil {
			runtime.LogWarningf(app.ctx, "failed to set field %s: %v", fieldName, err)
		}
	}

	return nil
}

// UpdateOptions schreibt die Eingaben der Oberfläche zurück: ein changeoption je
// Option, die sich gegenüber dem gelesenen Serverstand geändert hat, danach optional
// ein reloadoptions. Namen, die der Server nicht gemeldet hat, werden nie gesendet.
func (app *App) UpdateOptions(values map[string]string, reloadOptions bool) bool {
	defer pzOptions_update()

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

	if err := pzOptions_update(); err != nil {
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

		command := fmt.Sprintf("changeoption %s \"%s\"", option.Name, option.Value)
		res, err := conn.Execute(command)

		if err != nil || !optionUpdateSucceeded(option, kinds[option.Name], res) {
			runtime.LogErrorf(app.ctx, "Failed to update %s: %v (response: %s)", option.Name, err, res)
			failed = append(failed, option.Name)
		}
	}

	return failed
}

func (app *App) UpdatePzOptions(newOptions PzOptions, reloadOptions bool) bool {
	defer pzOptions_update()

	optionsToUpdate := app.diffOptions(newOptions)

	if len(optionsToUpdate) == 0 {
		app.SendNotification(Notification{Title: "rcon.no_options_to_update", Variant: "warning"})
		return false
	}

	successCount := app.applyOptions(optionsToUpdate)

	if successCount != len(optionsToUpdate) {
		app.SendNotification(Notification{Title: "rcon.failed_to_update_n_options", Parameters: map[string]string{
			"n": fmt.Sprintf("%d", len(optionsToUpdate)-successCount),
		}, Variant: "error"})
		return false
	}

	if err := pzOptions_update(); err != nil {
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

		success := command.execute() == 1
		if !success {
			app.SendNotification(Notification{Title: "rcon.reloadOptions.single_fail", Variant: "error"})
			return false
		} else {
			app.SendNotification(Notification{Title: "rcon.options_saved_and_reloaded", Variant: "success"})
			return true
		}
	}

	app.SendNotification(Notification{Title: "rcon.options_updated", Variant: "success"})
	return true
}

func (app *App) diffOptions(newOptions PzOptions) []OptionPair {
	var optionsToUpdate []OptionPair

	newVal := reflect.ValueOf(newOptions)
	oldVal := reflect.ValueOf(pzOptions)

	for i := 0; i < newVal.NumField(); i++ {
		fieldName := newVal.Type().Field(i).Name
		newField := newVal.Field(i).Interface()
		oldField := oldVal.FieldByName(fieldName).Interface()

		if newField != oldField {
			runtime.LogDebugf(app.ctx, "Updating %s from %v to %v", fieldName, oldField, newField)
			optionsToUpdate = append(optionsToUpdate, OptionPair{
				Name:  fieldName,
				Value: fmt.Sprintf("%v", newField),
			})
		}
	}

	return optionsToUpdate
}

func (app *App) applyOptions(options []OptionPair) int {
	defer runtime.EventsEmit(app.ctx, "setProgress", 0)
	runtime.EventsEmit(app.ctx, "setProgress", 10)

	successCount := 0
	optionCount := len(options)

	for _, option := range options {
		runtime.EventsEmit(app.ctx, "setProgress", float64(successCount)/float64(optionCount)*100)

		command := fmt.Sprintf("changeoption %s \"%s\"", option.Name, option.Value)
		res, err := conn.Execute(command)

		if err == nil && isOptionUpdateSuccessful(option, res) {
			successCount++
		} else {
			runtime.LogErrorf(app.ctx, "Failed to update %s: %v", option.Name, err)
		}
	}

	return successCount
}

func isOptionUpdateSuccessful(option OptionPair, res string) bool {
	// Extract the server response value for the given option
	expectedResponse := fmt.Sprintf("Option : %s is now : %s", option.Name, option.Value)

	if res == expectedResponse {
		return true
	}

	// For floats, handle formatting differences
	field := reflect.ValueOf(pzOptions).FieldByName(option.Name)
	if !field.IsValid() {
		runtime.LogErrorf(nil, "Invalid field name: %s", option.Name)
		return false
	}

	switch field.Kind() {
	case reflect.Float64:
		value, err := strconv.ParseFloat(option.Value, 64)
		if err != nil {
			runtime.LogErrorf(nil, "Failed to parse float: %s", option.Value)
			return false
		}

		normalized := fmt.Sprintf("Option : %s is now : %.1f", option.Name, value)
		return res == normalized
	default:
		return false
	}
}
