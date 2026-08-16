package main

// Die Anbindung des Sandbox-Editors an Konfiguration, Dateisystem und Oberfläche.
// Der Editor kommt ohne RCON aus — er ist der erste Teil der Anwendung, für den das
// gilt, und funktioniert auch bei getrennter Verbindung.

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SandboxAccess beschreibt den eingerichteten Zugangsweg und die daraus aufgelöste
// Datei. Ist nichts eingerichtet, trägt `Reason` den Grund — der Bereich bleibt
// inaktiv und begründet das, statt einen Pfad zu raten (PG-20).
type SandboxAccess struct {
	Configured bool   `json:"configured"`
	Mode       string `json:"mode"`
	Path       string `json:"path"`
	File       string `json:"file"`
	Reason     string `json:"reason"`
}

// SandboxDocument ist der gelesene Stand der Datei. `Checksum` hält fest, wie die
// Datei beim Lesen aussah — beim Speichern wird dagegen geprüft, damit eine
// Handänderung in der Zwischenzeit nicht stillschweigend überschrieben wird.
type SandboxDocument struct {
	Success  bool         `json:"success"`
	File     string       `json:"file"`
	Vars     []SandboxVar `json:"vars"`
	Checksum string       `json:"checksum"`
	Error    string       `json:"error"`
}

// SandboxSaveResult meldet, was geschrieben wurde — und wohin die Sicherung ging.
type SandboxSaveResult struct {
	Success bool     `json:"success"`
	Backup  string   `json:"backup"`
	Changed []string `json:"changed"`
	Error   string   `json:"error"`
}

func listDirectory(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names, nil
}

func sandboxAccess() SandboxAccess {
	mode := ""
	if config.SandboxAccessMode != nil {
		mode = *config.SandboxAccessMode
	}
	path := ""
	if config.SandboxPath != nil {
		path = *config.SandboxPath
	}

	access := SandboxAccess{Mode: mode, Path: path}

	file, err := resolveSandboxPath(mode, path, listDirectory)
	if err != nil {
		access.Reason = err.Error()
		return access
	}

	// Erst wenn die Datei wirklich dort liegt, ist der Bereich benutzbar.
	if _, err := os.Stat(file); err != nil {
		access.Reason = err.Error()
		return access
	}

	access.Configured = true
	access.File = file
	return access
}

// SandboxAccess liefert den Zugangsweg an die Oberfläche.
func (app *App) SandboxAccess() SandboxAccess {
	return sandboxAccess()
}

// SetSandboxAccess richtet den Zugangsweg ein und prüft ihn sofort. Die Einstellung
// landet in der Konfiguration und überlebt damit einen Neustart.
func (app *App) SetSandboxAccess(mode string, path string) SandboxAccess {
	config.SandboxAccessMode = &mode
	config.SandboxPath = &path

	access := sandboxAccess()
	if access.Configured {
		runtime.LogInfof(app.ctx, "Sandbox file resolved: %s", access.File)
	} else {
		runtime.LogWarningf(app.ctx, "Sandbox access not usable: %s", access.Reason)
	}

	return access
}

// ChooseSandboxPathDialog fragt nach der Sandbox-Datei bzw. dem Serververzeichnis.
func (app *App) ChooseSandboxPathDialog(mode string) string {
	if mode == sandboxModeLocal {
		path, err := runtime.OpenFileDialog(app.ctx, runtime.OpenDialogOptions{
			Title: "Select the SandboxVars file",
			Filters: []runtime.FileFilter{
				{DisplayName: "SandboxVars", Pattern: "*" + sandboxFileSuffix},
			},
		})
		if err != nil {
			runtime.LogWarning(app.ctx, err.Error())
			return ""
		}
		return path
	}

	path, err := runtime.OpenDirectoryDialog(app.ctx, runtime.OpenDialogOptions{
		Title: "Select the server data directory",
	})
	if err != nil {
		runtime.LogWarning(app.ctx, err.Error())
		return ""
	}
	return path
}

// ReadSandboxVars liest die Datei. Es wird ausschließlich gelesen — die Datei bleibt
// unverändert.
func (app *App) ReadSandboxVars() SandboxDocument {
	access := sandboxAccess()
	if !access.Configured {
		return SandboxDocument{Error: access.Reason}
	}

	content, err := os.ReadFile(access.File)
	if err != nil {
		runtime.LogErrorf(app.ctx, "Error reading the sandbox file: %v", err)
		return SandboxDocument{File: access.File, Error: err.Error()}
	}

	vars, err := parseSandboxLua(string(content))
	if err != nil {
		runtime.LogErrorf(app.ctx, "Error parsing the sandbox file: %v", err)
		return SandboxDocument{File: access.File, Error: err.Error()}
	}

	runtime.LogInfof(app.ctx, "Sandbox file read: %d variables", len(vars))
	return SandboxDocument{Success: true, File: access.File, Vars: vars, Checksum: hashString(string(content))}
}

// SaveSandboxVars schreibt geänderte Werte in die Datei — strukturerhaltend und erst,
// nachdem eine Sicherung mit Zeitstempel daneben liegt (PG-17). Ein misslungener
// Schreibvorgang ohne Sicherung wäre ein Datenverlust, den niemand rückgängig macht.
//
// Kein Wert geht dabei über RCON (PG-16); die Änderungen greifen erst nach einem
// Serverneustart (PG-19) — das sagt die Oberfläche am Speichern-Knopf.
func (app *App) SaveSandboxVars(checksum string, values map[string]string) SandboxSaveResult {
	access := sandboxAccess()
	if !access.Configured {
		return SandboxSaveResult{Error: access.Reason}
	}

	content, err := os.ReadFile(access.File)
	if err != nil {
		runtime.LogErrorf(app.ctx, "Error reading the sandbox file: %v", err)
		return SandboxSaveResult{Error: err.Error()}
	}

	// Hat jemand die Datei seit dem Lesen von Hand geändert, wird nichts überschrieben:
	// Die Oberfläche kennt nur den alten Stand und würde fremde Änderungen zurückdrehen.
	if checksum != "" && hashString(string(content)) != checksum {
		app.SendNotification(Notification{Title: "sandbox.file_changed_meanwhile", Variant: "error"})
		return SandboxSaveResult{Error: "the file changed since it was read"}
	}

	vars, err := parseSandboxLua(string(content))
	if err != nil {
		return SandboxSaveResult{Error: err.Error()}
	}

	updated, changed, err := writeSandboxValues(string(content), vars, values)
	if err != nil {
		app.SendNotification(Notification{Title: "sandbox.save_failed", Message: err.Error(), Variant: "error"})
		return SandboxSaveResult{Error: err.Error()}
	}

	if len(changed) == 0 {
		app.SendNotification(Notification{Title: "sandbox.nothing_to_save", Variant: "warning"})
		return SandboxSaveResult{Success: true}
	}

	backup, err := backupSandboxFile(access.File, content)
	if err != nil {
		runtime.LogErrorf(app.ctx, "Error writing the sandbox backup: %v", err)
		app.SendNotification(Notification{Title: "sandbox.save_failed", Message: err.Error(), Variant: "error"})
		return SandboxSaveResult{Error: err.Error()}
	}

	info, err := os.Stat(access.File)
	mode := os.FileMode(0o644)
	if err == nil {
		mode = info.Mode()
	}

	if err := writeFileAtomically(access.File, []byte(updated), mode); err != nil {
		// Die Sicherung liegt daneben — ihr Pfad gehört in die Meldung, damit die
		// Nutzerin weiß, woher sie den vorherigen Stand bekommt.
		runtime.LogErrorf(app.ctx, "Error writing the sandbox file: %v", err)
		app.SendNotification(Notification{
			Title:      "sandbox.save_failed_with_backup",
			Message:    err.Error(),
			Variant:    "error",
			Parameters: map[string]string{"backup": backup},
		})
		return SandboxSaveResult{Backup: backup, Error: err.Error()}
	}

	runtime.LogInfof(app.ctx, "Sandbox file written: %d values changed, backup at %s", len(changed), backup)
	// Die Bestätigung behauptet nicht, die Änderung sei übernommen: Sie gilt erst
	// nach einem Serverneustart (PG-19).
	app.SendNotification(Notification{
		Title:   "sandbox.saved",
		Message: "sandbox.saved_restart_hint",
		Variant: "success",
		Parameters: map[string]string{
			"count":  fmt.Sprintf("%d", len(changed)),
			"backup": backup,
		},
	})
	return SandboxSaveResult{Success: true, Backup: backup, Changed: changed}
}

// backupSandboxFile legt den bisherigen Stand mit Zeitstempel neben die Datei.
func backupSandboxFile(path string, content []byte) (string, error) {
	backup := fmt.Sprintf("%s.bak-%s", path, time.Now().Format("20060102-150405"))
	if err := os.WriteFile(backup, content, 0o644); err != nil {
		return "", err
	}
	return backup, nil
}

// writeFileAtomically schreibt neben die Zieldatei und benennt dann um. Ein
// abgebrochener Schreibvorgang hinterlässt damit keine halbe Weltkonfiguration —
// die Datei ist entweder die alte oder die neue.
func writeFileAtomically(path string, content []byte, mode os.FileMode) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp-")
	if err != nil {
		return err
	}
	defer os.Remove(temporary.Name())

	if _, err := temporary.Write(content); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Chmod(temporary.Name(), mode); err != nil {
		return err
	}

	return os.Rename(temporary.Name(), path)
}
