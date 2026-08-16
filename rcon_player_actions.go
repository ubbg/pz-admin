package main

// Befehle, die auf eine einzelne Person wirken: Fernwirkung (unsichtbar,
// wanddurchlässig, stummgeschaltet), Safehouse-Verwaltung, Schlüssel und
// Kartensymbole. Alle Befehle laut pzwiki „Admin commands", Fassung 42.20.2.
//
// Build 42 hat für Unverwundbarkeit, Unsichtbarkeit und Versetzen je zwei Formen.
// Die kurze zielt auf die eigene Figur der Konsole und tut auf einem Dedicated
// Server nichts; hier steht durchgehend die Form für andere Spieler.

import (
	"fmt"
	"strings"
)

// toggleValue übersetzt einen Schalter in die Schreibweise der Konsole (-true/-false).
func toggleValue(value bool) string {
	if value {
		return "-true"
	}
	return "-false"
}

// InvisiblePlayers macht Spieler für Zombies unsichtbar.
// Wiki: /invisibleplayer "username" -value
func (app *App) InvisiblePlayers(names []string, value bool) {
	command := RCONCommand{
		CommandTemplate: "invisibleplayer {name} {value}",
		PlayerNames:     names,
		Args: []RCONCommandParam{
			{Name: "value", Value: toggleValue(value), Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		ErrorCheck:   playerNotFound,
		UpdateFunc: func(name string, response string) {
			setPlayerFlag(name, func(player *Player) { player.Invisible = value })
		},
		EmitUpdatePlayers: true,
		Notifications: RCONCommandNotifications{
			AllSuccess:    "rcon.invisible.all_success",
			AllFail:       "rcon.invisible.all_fail",
			Partial:       "rcon.invisible.partial",
			SingleSuccess: "rcon.invisible.single_success",
			SingleFail:    "rcon.invisible.single_fail",
		},
	}

	command.execute()
}

// NoClipPlayers lässt Spieler durch Wände gehen.
// Wiki: /noclip "username" -value
func (app *App) NoClipPlayers(names []string, value bool) {
	command := RCONCommand{
		CommandTemplate: "noclip {name} {value}",
		PlayerNames:     names,
		Args: []RCONCommandParam{
			{Name: "value", Value: toggleValue(value), Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		ErrorCheck:   playerNotFound,
		UpdateFunc: func(name string, response string) {
			setPlayerFlag(name, func(player *Player) { player.NoClip = value })
		},
		EmitUpdatePlayers: true,
		Notifications: RCONCommandNotifications{
			AllSuccess:    "rcon.noclip.all_success",
			AllFail:       "rcon.noclip.all_fail",
			Partial:       "rcon.noclip.partial",
			SingleSuccess: "rcon.noclip.single_success",
			SingleFail:    "rcon.noclip.single_fail",
		},
	}

	command.execute()
}

// VoiceBanPlayers sperrt die Sprachübertragung eines Spielers.
// Wiki: /voiceban "username" -value
func (app *App) VoiceBanPlayers(names []string, value bool) {
	command := RCONCommand{
		CommandTemplate: "voiceban {name} {value}",
		PlayerNames:     names,
		Args: []RCONCommandParam{
			{Name: "value", Value: toggleValue(value), Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		ErrorCheck:   playerNotFound,
		UpdateFunc: func(name string, response string) {
			setPlayerFlag(name, func(player *Player) { player.VoiceBanned = value })
		},
		EmitUpdatePlayers: true,
		Notifications: RCONCommandNotifications{
			AllSuccess:    "rcon.voiceban.all_success",
			AllFail:       "rcon.voiceban.all_fail",
			Partial:       "rcon.voiceban.partial",
			SingleSuccess: "rcon.voiceban.single_success",
			SingleFail:    "rcon.voiceban.single_fail",
		},
	}

	command.execute()
}

// TeleportPlayerToPlayer versetzt Spieler zu einem anderen Spieler.
// Wiki: /teleportplayer "player1" "player2"
func (app *App) TeleportPlayerToPlayer(names []string, targetUser string) {
	targetUser = strings.TrimSpace(targetUser)
	if targetUser == "" {
		return
	}

	command := RCONCommand{
		CommandTemplate: "teleportplayer {name} {target}",
		PlayerNames:     names,
		Args: []RCONCommandParam{
			{Name: "target", Value: fmt.Sprintf("\"%s\"", targetUser), Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		ErrorCheck:   playerNotFound,
		Notifications: RCONCommandNotifications{
			AllSuccess:    "rcon.teleport.all_success",
			AllFail:       "rcon.teleport.all_fail",
			Partial:       "rcon.teleport.partial",
			SingleSuccess: "rcon.teleport.single_success",
			SingleFail:    "rcon.teleport.single_fail",
			Parameters:    map[string]string{"target": targetUser},
		},
	}

	command.execute()
}

// AddToSafehouse nimmt Spieler in ein Safehouse auf.
// Wiki: /addtosafehouse — die Seite führt den Befehl als WIP ohne Aufrufbeispiel;
// die Form mit Spielernamen ist die einzige, die zum übrigen Muster passt.
func (app *App) AddToSafehouse(names []string) {
	command := RCONCommand{
		CommandTemplate: "addtosafehouse {name}",
		PlayerNames:     names,
		SuccessCheck:    succeededUnlessUsage,
		ErrorCheck:      playerNotFound,
		Notifications: RCONCommandNotifications{
			AllSuccess:    "rcon.addToSafehouse.all_success",
			AllFail:       "rcon.addToSafehouse.all_fail",
			Partial:       "rcon.addToSafehouse.partial",
			SingleSuccess: "rcon.addToSafehouse.single_success",
			SingleFail:    "rcon.addToSafehouse.single_fail",
		},
	}

	command.execute()
}

// KickFromSafehouse wirft Spieler aus ihrem Safehouse.
// Wiki: /kickfromsafehouse — ebenfalls als WIP ohne Aufrufbeispiel geführt.
func (app *App) KickFromSafehouse(names []string) {
	command := RCONCommand{
		CommandTemplate: "kickfromsafehouse {name}",
		PlayerNames:     names,
		SuccessCheck:    succeededUnlessUsage,
		ErrorCheck:      playerNotFound,
		Notifications: RCONCommandNotifications{
			AllSuccess:    "rcon.kickFromSafehouse.all_success",
			AllFail:       "rcon.kickFromSafehouse.all_fail",
			Partial:       "rcon.kickFromSafehouse.partial",
			SingleSuccess: "rcon.kickFromSafehouse.single_success",
			SingleFail:    "rcon.kickFromSafehouse.single_fail",
		},
	}

	command.execute()
}

// ReleaseSafehouse gibt das Safehouse der Konsolenfigur frei.
// Wiki: /releasesafehouse — der Befehl kennt keinen Spielernamen und wirkt auf die
// eigene Figur; auf einem Dedicated Server ohne Charakter bleibt er wirkungslos.
func (app *App) ReleaseSafehouse() bool {
	command := RCONCommand{
		CommandTemplate: "releasesafehouse",
		SuccessCheck:    succeededUnlessUsage,
		Notifications: RCONCommandNotifications{
			SingleSuccess: "rcon.releaseSafehouse.single_success",
			SingleFail:    "rcon.releaseSafehouse.single_fail",
		},
	}

	return command.execute() == 1
}

// AddKey gibt Spielern einen Schlüssel.
// Wiki: /addkey "username" "keyId" "name"
func (app *App) AddKey(names []string, keyId string, keyName string) {
	keyId = strings.TrimSpace(keyId)
	if keyId == "" {
		return
	}

	command := RCONCommand{
		CommandTemplate: "addkey {name} {keyId} {keyName}",
		PlayerNames:     names,
		Args: []RCONCommandParam{
			{Name: "keyId", Value: fmt.Sprintf("\"%s\"", keyId), Mandatory: true},
			{Name: "keyName", Value: func() interface{} {
				if strings.TrimSpace(keyName) == "" {
					return nil
				}
				return fmt.Sprintf("\"%s\"", strings.TrimSpace(keyName))
			}()},
		},
		SuccessCheck: succeededUnlessUsage,
		ErrorCheck:   playerNotFound,
		Notifications: RCONCommandNotifications{
			AllSuccess:    "rcon.addKey.all_success",
			AllFail:       "rcon.addKey.all_fail",
			Partial:       "rcon.addKey.partial",
			SingleSuccess: "rcon.addKey.single_success",
			SingleFail:    "rcon.addKey.single_fail",
			Parameters:    map[string]string{"keyId": keyId},
		},
	}

	command.execute()
}

// RemoveMapSymbolsForUser löscht die geteilten Kartensymbole von Spielern.
// Wiki: /removemapsymbolsforuser "username"
func (app *App) RemoveMapSymbolsForUser(names []string) {
	command := RCONCommand{
		CommandTemplate: "removemapsymbolsforuser {name}",
		PlayerNames:     names,
		SuccessCheck:    succeededUnlessUsage,
		ErrorCheck:      playerNotFound,
		Notifications: RCONCommandNotifications{
			AllSuccess:    "rcon.removeMapSymbols.all_success",
			AllFail:       "rcon.removeMapSymbols.all_fail",
			Partial:       "rcon.removeMapSymbols.partial",
			SingleSuccess: "rcon.removeMapSymbols.single_success",
			SingleFail:    "rcon.removeMapSymbols.single_fail",
		},
	}

	command.execute()
}

// RemoveItems entfernt Gegenstände von der Figur der Konsole.
// Wiki: /removeitem "module.item" count — der Befehl kennt keinen Spielernamen
// ("Remove items from yourself"), anders als additem. Auf einem Dedicated Server
// ohne eigene Figur bleibt er deshalb wirkungslos; die Oberfläche sagt das dazu.
func (app *App) RemoveItems(itemRecords []ItemRecord) {
	for _, item := range itemRecords {
		command := RCONCommand{
			CommandTemplate: "removeitem {item} {count}",
			Args: []RCONCommandParam{
				{Name: "item", Value: fmt.Sprintf("\"%s\"", item.ItemId), Mandatory: true},
				{Name: "count", Value: item.Count, Mandatory: true},
			},
			SuccessCheck: succeededUnlessUsage,
			Notifications: RCONCommandNotifications{
				SingleSuccess: "rcon.removeItems.single_success",
				SingleFail:    "rcon.removeItems.single_fail",
				Parameters:    map[string]string{"item": item.ItemId},
			},
		}

		command.execute()
	}
}

// playerNotFound ist die bekannte Fehlantwort auf einen unbekannten Spielernamen.
func playerNotFound(name string, response string) bool {
	return response == fmt.Sprintf("Can't find player %s", name) ||
		response == fmt.Sprintf("User %s not found.", name)
}

// setPlayerFlag schreibt den lokal geführten Zustand eines Spielers fort. Diese
// Zustände gibt RCON nicht her; sie liegen wie Bann-Status und Zugriffsstufe in
// players.json und überleben damit einen Neustart.
func setPlayerFlag(name string, apply func(*Player)) {
	for i := range players {
		if players[i].Name == name {
			apply(&players[i])
			return
		}
	}
}
