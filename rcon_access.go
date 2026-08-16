package main

// Zugang und Moderation — SteamID-Erlaubnisliste, Kontopasswörter, Sperren nach IP
// und SteamID. Alle Befehle laut pzwiki „Admin commands", Fassung 42.20.2.
//
// Bewusst ohne eigene Liste: `addsteamid`/`removesteamid` und die Sperrlisten sind
// über RCON nicht lesbar. Eine in der Anwendung geführte Liste wäre eine Behauptung,
// deshalb quittiert die Oberfläche nur den einzelnen Vorgang (siehe
// specs/commands/CMD-001.spec.md).

import (
	"fmt"
	"strings"
)

// AddSteamId trägt eine SteamID in die Erlaubnisliste des Servers ein.
// Wiki: /addSteamID "steamid"
func (app *App) AddSteamId(steamId string) bool {
	return app.steamIdCommand("addsteamid", steamId, "addSteamId")
}

// RemoveSteamId nimmt eine SteamID aus der Erlaubnisliste.
// Wiki: /removeSteamID "steamid"
func (app *App) RemoveSteamId(steamId string) bool {
	return app.steamIdCommand("removesteamid", steamId, "removeSteamId")
}

// BanSteamId sperrt eine SteamID.
// Wiki: /banid SteamID
func (app *App) BanSteamId(steamId string) bool {
	return app.steamIdCommand("banid", steamId, "banSteamId")
}

// UnbanSteamId hebt die Sperre einer SteamID auf.
// Wiki: /unbanid SteamID
func (app *App) UnbanSteamId(steamId string) bool {
	return app.steamIdCommand("unbanid", steamId, "unbanSteamId")
}

func (app *App) steamIdCommand(commandName string, steamId string, notificationKey string) bool {
	steamId = strings.TrimSpace(steamId)
	if steamId == "" {
		return false
	}

	command := RCONCommand{
		CommandTemplate: fmt.Sprintf("%s {steamId}", commandName),
		Args: []RCONCommandParam{
			{Name: "steamId", Value: fmt.Sprintf("\"%s\"", steamId), Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		Notifications: RCONCommandNotifications{
			SingleSuccess: fmt.Sprintf("rcon.%s.single_success", notificationKey),
			SingleFail:    fmt.Sprintf("rcon.%s.single_fail", notificationKey),
			Parameters:    map[string]string{"steamId": steamId},
		},
	}

	return command.execute() == 1
}

// BanIp sperrt eine IP-Adresse. Der Bann kann mehr als eine Person treffen — die
// Rückfrage dazu stellt die Oberfläche, bevor sie hier landet.
// Wiki: /banip IP
func (app *App) BanIp(ip string) bool {
	return app.ipCommand("banip", ip, "banIp")
}

// UnbanIp hebt die Sperre einer IP-Adresse auf.
// Wiki: /unbanip IP
func (app *App) UnbanIp(ip string) bool {
	return app.ipCommand("unbanip", ip, "unbanIp")
}

func (app *App) ipCommand(commandName string, ip string, notificationKey string) bool {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return false
	}

	command := RCONCommand{
		CommandTemplate: fmt.Sprintf("%s {ip}", commandName),
		Args: []RCONCommandParam{
			{Name: "ip", Value: ip, Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		Notifications: RCONCommandNotifications{
			SingleSuccess: fmt.Sprintf("rcon.%s.single_success", notificationKey),
			SingleFail:    fmt.Sprintf("rcon.%s.single_fail", notificationKey),
			Parameters:    map[string]string{"ip": ip},
		},
	}

	return command.execute() == 1
}

// SetUserPassword ändert das Passwort eines Benutzerkontos.
// Wiki: /setpassword "username" "newpassword"
//
// Das Passwort erscheint weder im Protokoll noch in einer Meldung — deshalb steht in
// den Benachrichtigungen nur der Benutzername.
func (app *App) SetUserPassword(username string, password string) bool {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return false
	}

	command := RCONCommand{
		CommandTemplate: "setpassword {name} {password}",
		PlayerNames:     []string{username},
		Args: []RCONCommandParam{
			{Name: "password", Value: fmt.Sprintf("\"%s\"", password), Mandatory: true},
		},
		SuccessCheck: succeededUnlessUsage,
		Notifications: RCONCommandNotifications{
			SingleSuccess: "rcon.setPassword.single_success",
			SingleFail:    "rcon.setPassword.single_fail",
		},
	}

	return command.execute() == 1
}
