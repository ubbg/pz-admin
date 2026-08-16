package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/minio/selfupdate"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Herkunft der Selbstaktualisierung — die einzige Stelle, an der sie steht.
//
// Ein leerer updateRepoOwner heißt: Für diesen Build gibt es keine eigene
// Release-Quelle. Die Anwendung prüft dann gar nicht, statt still dem Upstream zu
// folgen — ein Upstream-Release brächte alle Build-41-Annahmen zurück, die dieser
// Fork gerade beseitigt. Wer eigene Releases veröffentlicht, trägt hier seinen
// GitHub-Benutzer ein; alles Weitere (Prüfung, Download-URL, Anzeige in den
// Einstellungen) folgt daraus.
const (
	updateRepoOwner = ""
	updateRepoName  = "pz-admin"
)

type UpdateInfo struct {
	UpdateAvailable bool   `json:"updateAvailable"`
	CurrentVersion  string `json:"currentVersion"`
	LatestVersion   string `json:"latestVersion"`
	Name            string `json:"name"`
	ReleaseNotes    string `json:"releaseNotes"`
	DownloadUrl     string `json:"downloadUrl"`
	ReleaseUrl      string `json:"releaseUrl"`
}

type Release struct {
	TagName      string `json:"tag_name"`
	Name         string `json:"name"`
	ReleaseNotes string `json:"body"`
	Prerelease   bool   `json:"prerelease"`
}

// UpdateSource beschreibt, welcher Release-Quelle diese Anwendung folgt. Die
// Oberfläche zeigt das an, damit niemand raten muss, woher eine Aktualisierung käme.
type UpdateSource struct {
	Configured bool   `json:"configured"`
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	Url        string `json:"url"`
}

func updateSource() UpdateSource {
	if updateRepoOwner == "" {
		return UpdateSource{Configured: false, Repo: updateRepoName}
	}

	return UpdateSource{
		Configured: true,
		Owner:      updateRepoOwner,
		Repo:       updateRepoName,
		Url:        fmt.Sprintf("https://github.com/%s/%s", updateRepoOwner, updateRepoName),
	}
}

func (app *App) GetUpdateSource() UpdateSource {
	return updateSource()
}

func (app *App) CheckForUpdate() UpdateInfo {
	os := app.GetOs()     // windows, linux, macos
	arch := app.GetArch() // amd64, arm64

	source := updateSource()

	// Ohne eigene Release-Quelle wird nichts geprüft: kein Netzzugriff, keine
	// Fehlermeldung, kein fremdes Release.
	if !source.Configured {
		runtime.LogInfo(app.ctx, "No update source configured, skipping update check")
		return UpdateInfo{
			UpdateAvailable: false,
			CurrentVersion:  version,
		}
	}

	if !(os == "windows" || os == "linux" || os == "macos") {
		return UpdateInfo{
			UpdateAvailable: false,
			CurrentVersion:  version,
			LatestVersion:   "",
			ReleaseNotes:    "",
			DownloadUrl:     "",
			ReleaseUrl:      "",
		}
	}

	var updateInfo UpdateInfo = UpdateInfo{
		UpdateAvailable: false,
		CurrentVersion:  version,
		LatestVersion:   "",
		ReleaseNotes:    "",
		DownloadUrl:     "",
		ReleaseUrl:      "",
	}
	updateInfo.CurrentVersion = version

	// Set last update check
	lastUpdateCheck := int(time.Now().Unix())
	config.LastUpdateCheck = &lastUpdateCheck

	repoOwner := source.Owner
	repoName := source.Repo

	// GitHub API endpoint to fetch latest release
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	runtime.LogDebug(app.ctx, "GitHub API URL: "+apiUrl)

	// Make GET request to GitHub API
	resp, err := http.Get(apiUrl)
	if err != nil {
		runtime.LogError(app.ctx, "Error sending request: "+err.Error())
		app.SendNotification(Notification{
			Title:   "settings.setting.update.failed_to_check_for_updates",
			Variant: "error",
		})
		return updateInfo
	}
	defer resp.Body.Close()
	runtime.LogDebug(app.ctx, fmt.Sprintf("GitHub API response status: %d", resp.StatusCode))

	// Check if response was successful
	if resp.StatusCode != http.StatusOK {
		app.SendNotification(Notification{
			Title:   "settings.setting.update.failed_to_check_for_updates",
			Variant: "error",
		})
		return updateInfo
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		runtime.LogError(app.ctx, "Error reading response: "+err.Error())
		return updateInfo
	}
	runtime.LogTrace(app.ctx, "GitHub API response body: "+string(body))

	// Parse JSON response
	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		runtime.LogError(app.ctx, "Error decoding JSON: "+err.Error())
		return updateInfo
	}

	// Extract release information
	latestVersion := release.TagName
	releaseNotes := release.ReleaseNotes
	prerelease := release.Prerelease
	name := release.Name

	versionNoPrefix := strings.TrimPrefix(latestVersion, "v")

	// Determine the correct download URL based on OS and architecture
	var fileName string
	switch os {
	case "windows":
		if arch == "amd64" {
			fileName = "pz-admin_" + versionNoPrefix + "_windows_amd64.exe"
		} else if arch == "arm64" {
			fileName = "pz-admin_" + versionNoPrefix + "_windows_arm64.exe"
		}
	case "linux":
		if arch == "amd64" {
			fileName = "pz-admin_" + versionNoPrefix + "_linux_amd64"
		} else if arch == "arm64" {
			fileName = "pz-admin_" + versionNoPrefix + "_linux_arm64"
		}
	}

	downloadUrl := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", repoOwner, repoName, latestVersion, fileName)

	// Parse current and latest versions
	parsedVersion, err := semver.ParseTolerant(version)
	if err != nil {
		runtime.LogError(app.ctx, "Error parsing current version: "+err.Error())
		updateInfo.UpdateAvailable = false
		return updateInfo
	}
	parsedLatestVersion, err := semver.ParseTolerant(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(latestVersion, "v", ""), "-", ""), "alpha", ""), "beta", ""))
	if err != nil {
		runtime.LogError(app.ctx, "Error parsing latest version: "+err.Error())
		updateInfo.UpdateAvailable = false
		return updateInfo
	}

	// Log release information
	runtime.LogDebug(app.ctx, fmt.Sprintf("Current version: %s", parsedVersion))
	runtime.LogDebug(app.ctx, fmt.Sprintf("Latest version: %s", parsedLatestVersion))
	runtime.LogDebug(app.ctx, fmt.Sprintf("Prerelease: %t", prerelease))
	runtime.LogDebug(app.ctx, fmt.Sprintf("Release name: %s", name))
	runtime.LogDebug(app.ctx, fmt.Sprintf("Release notes: %s", releaseNotes))
	runtime.LogDebug(app.ctx, fmt.Sprintf("Download URL: %s", downloadUrl))

	// Check if a new version is available
	if parsedVersion.Compare(parsedLatestVersion) < 0 && !prerelease {
		runtime.LogInfo(app.ctx, fmt.Sprintf("A new version (%s) is available.", latestVersion))
		updateInfo.UpdateAvailable = true
		updateInfo.LatestVersion = latestVersion
		updateInfo.Name = name
		updateInfo.ReleaseNotes = releaseNotes
		updateInfo.DownloadUrl = downloadUrl
		updateInfo.ReleaseUrl = fmt.Sprintf("https://github.com/%s/%s/releases/latest", repoOwner, repoName)
	} else {
		runtime.LogInfo(app.ctx, "You have the latest version.")
	}

	return updateInfo
}

func (app *App) Update(downloadUrl string) error {
	// Ohne eigene Release-Quelle gibt es nichts einzuspielen — auch nicht aus einer
	// URL, die noch von einem früheren Lauf im Speicher liegt.
	if !updateSource().Configured {
		runtime.LogWarning(app.ctx, "No update source configured, refusing to apply an update")
		app.SendNotification(Notification{
			Title:   "settings.setting.update.no_update_source",
			Variant: "warning",
		})
		return errors.New("no update source configured")
	}

	// Log the download URL
	runtime.LogInfo(app.ctx, "Starting update download from: "+downloadUrl)

	resp, err := http.Get(downloadUrl)
	if err != nil {
		runtime.LogError(app.ctx, "Error downloading update: "+err.Error())
		app.SendNotification(Notification{
			Title:   "settings.setting.update.failed_to_download_update",
			Variant: "error",
		})
		return err
	}
	defer resp.Body.Close()

	// Log the status code from the response
	runtime.LogDebug(app.ctx, fmt.Sprintf("Download response status: %d", resp.StatusCode))

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		app.SendNotification(Notification{
			Title:   "settings.setting.update.failed_to_download_update",
			Variant: "error",
		})
		return err
	}

	// Apply the update
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		runtime.LogError(app.ctx, "Error applying update: "+err.Error())
		app.SendNotification(Notification{
			Title:   "settings.setting.update.failed_to_apply_update",
			Message: err.Error(),
			Variant: "error",
		})
		return err
	}

	runtime.LogInfo(app.ctx, "Update applied successfully. Restarting.")
	app.SendNotification(Notification{
		Title:   "settings.setting.update.update_applied",
		Message: "settings.setting.update.restarting",
		Variant: "success",
	})

	// Restart the application
	app.RestartApplication([]string{"--goto", "settings__update", "--notify", "settings.setting.update.update_successful", "", "", "success"})

	return nil
}
