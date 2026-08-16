package main

// Serverweite Eingriffe von Build 42: Welt (Zombies, zweite Horden-Variante),
// Weltgenerator und Wartung (Lua neu laden, Protokollstufe, Statistik).
// Alle Befehle laut pzwiki „Admin commands", Fassung 42.20.2.

import (
	"fmt"
	"strings"
)

// RemoveZombies entfernt die geladenen Zombies. Nicht zurücknehmbar — die Rückfrage
// dazu stellt die Oberfläche, bevor sie hier landet.
// Wiki: /removezombies (ohne Aufrufbeispiel, als WIP geführt)
func (app *App) RemoveZombies() bool {
	command := RCONCommand{
		CommandTemplate: "removezombies",
		SuccessCheck:    succeededUnlessUsage,
		Notifications: RCONCommandNotifications{
			SingleSuccess: "rcon.removeZombies.single_success",
			SingleFail:    "rcon.removeZombies.single_fail",
		},
	}

	return command.execute() == 1
}

// CreateHorde2 löst die zweite Horden-Variante aus.
// Wiki: /createhorde2 — als WIP ohne Aufrufbeispiel geführt; die Argumente folgen
// createhorde (Anzahl, Spielername).
func (app *App) CreateHorde2(names []string, count int) {
	command := RCONCommand{
		CommandTemplate: "createhorde2 {count} {name}",
		PlayerNames:     names,
		Args: []RCONCommandParam{
			{Name: "count", Value: func() interface{} {
				if count < 0 {
					return nil
				}
				return count
			}(), Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		ErrorCheck:   playerNotFound,
		Notifications: RCONCommandNotifications{
			AllSuccess:    "rcon.createHorde2.all_success",
			AllFail:       "rcon.createHorde2.all_fail",
			Partial:       "rcon.createHorde2.partial",
			SingleSuccess: "rcon.createHorde2.single_success",
			SingleFail:    "rcon.createHorde2.single_fail",
		},
	}

	command.execute()
}

// WorldGen führt den stufenweisen Weltgenerator von Build 42.
// Wiki: /worldgen start | recheck | stop | status
//
// „stop" ist nicht zurücknehmbar; die Rückfrage stellt die Oberfläche.
func (app *App) WorldGen(action string) bool {
	if !worldGenActionAllowed(action) {
		return false
	}

	command := RCONCommand{
		CommandTemplate: "worldgen {action}",
		Args: []RCONCommandParam{
			{Name: "action", Value: action, Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		Notifications: RCONCommandNotifications{
			SingleSuccess: fmt.Sprintf("rcon.worldGen.%s.single_success", action),
			SingleFail:    fmt.Sprintf("rcon.worldGen.%s.single_fail", action),
		},
	}

	return command.execute() == 1
}

// WorldGenStatus liest den Fortschritt, ohne eine Meldung zu erzeugen — es wird
// gepollt, solange die Ansicht offen ist. Der Zustand wird gelesen, nie angenommen:
// Die Antwort des Servers geht unverändert an die Oberfläche.
func (app *App) WorldGenStatus() RconResponse {
	return app.SendRconCommand("worldgen status")
}

func worldGenActionAllowed(action string) bool {
	switch action {
	case "start", "recheck", "stop", "status":
		return true
	}
	return false
}

// ReloadAllLua lädt alle Lua-Skripte des Servers neu.
// Wiki: /reloadalllua
func (app *App) ReloadAllLua() bool {
	command := RCONCommand{
		CommandTemplate: "reloadalllua",
		SuccessCheck:    succeededUnlessUsage,
		Notifications: RCONCommandNotifications{
			SingleSuccess: "rcon.reloadAllLua.single_success",
			SingleFail:    "rcon.reloadAllLua.single_fail",
		},
	}

	return command.execute() == 1
}

// ReloadLua lädt ein einzelnes Lua-Skript neu.
// Wiki: /reloadlua "filename"
func (app *App) ReloadLua(filename string) bool {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return false
	}

	command := RCONCommand{
		CommandTemplate: "reloadlua {file}",
		Args: []RCONCommandParam{
			{Name: "file", Value: fmt.Sprintf("\"%s\"", filename), Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		Notifications: RCONCommandNotifications{
			SingleSuccess: "rcon.reloadLua.single_success",
			SingleFail:    "rcon.reloadLua.single_fail",
			Parameters:    map[string]string{"file": filename},
		},
	}

	return command.execute() == 1
}

// SetLogLevel setzt die Protokollstufe eines Bereichs.
// Wiki: /log "Type" "Level"
func (app *App) SetLogLevel(logType string, level string) bool {
	logType = strings.TrimSpace(logType)
	level = strings.TrimSpace(level)
	if logType == "" || level == "" {
		return false
	}

	command := RCONCommand{
		CommandTemplate: "log {type} {level}",
		Args: []RCONCommandParam{
			{Name: "type", Value: fmt.Sprintf("\"%s\"", logType), Mandatory: true},
			{Name: "level", Value: fmt.Sprintf("\"%s\"", level), Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		Notifications: RCONCommandNotifications{
			SingleSuccess: "rcon.setLogLevel.single_success",
			SingleFail:    "rcon.setLogLevel.single_fail",
			Parameters:    map[string]string{"type": logType, "level": level},
		},
	}

	return command.execute() == 1
}

// SetStats stellt die Statistikausgabe ein.
// Wiki: /stats none/file/console/all period
func (app *App) SetStats(mode string, period int) bool {
	switch mode {
	case "none", "file", "console", "all":
	default:
		return false
	}

	command := RCONCommand{
		CommandTemplate: "stats {mode} {period}",
		Args: []RCONCommandParam{
			{Name: "mode", Value: mode, Mandatory: true},
			{Name: "period", Value: func() interface{} {
				if mode == "none" {
					return nil
				}
				return period
			}()},
		},
		SuccessCheck: succeededUnlessUsage,
		Notifications: RCONCommandNotifications{
			SingleSuccess: "rcon.setStats.single_success",
			SingleFail:    "rcon.setStats.single_fail",
			Parameters:    map[string]string{"mode": mode, "period": fmt.Sprintf("%d", period)},
		},
	}

	return command.execute() == 1
}
