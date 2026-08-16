import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { SettingContent, SettingDescription, SettingLabel, SettingsItem } from "@/components/ui/settings-group";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { ChooseSandboxPathDialog, SandboxAccess, SetSandboxAccess } from "@/wailsjs/go/main/App";
import { main } from "@/wailsjs/go/models";
import { EventsEmit } from "@/wailsjs/runtime/runtime";

// Der Zugangsweg zur Sandbox-Datei wird ausdrücklich eingerichtet: lokaler Pfad oder
// Docker-Volume. Ohne Einrichtung rät die Anwendung nichts — auf einem fremden System
// träfe ein geratener Pfad eine falsche Datei, und diese Datei ist der Weltzustand
// eines laufenden Servers.
export function SandboxAccessSetting() {
  const { t } = useTranslation();
  const [access, setAccess] = useState<main.SandboxAccess>();
  const [mode, setMode] = useState("");
  const [path, setPath] = useState("");

  useEffect(() => {
    SandboxAccess().then((value) => {
      setAccess(value);
      setMode(value.mode);
      setPath(value.path);
    });
  }, []);

  const apply = (nextMode: string, nextPath: string) => {
    setMode(nextMode);
    setPath(nextPath);
    SetSandboxAccess(nextMode, nextPath).then((value) => {
      setAccess(value);
      // Der Sandbox-Bereich liest den Zugang neu, ohne dass die Anwendung neu starten muss.
      EventsEmit("sandbox-access-changed");
    });
  };

  const choose = () => {
    ChooseSandboxPathDialog(mode || "local").then((chosen) => {
      if (chosen) {
        apply(mode || "local", chosen);
      }
    });
  };

  return (
    <SettingsItem loading={access === undefined} vertical>
      <div className="flex justify-between items-start gap-4">
        <div>
          <SettingLabel>{t("settings.setting.sandbox_access.label")}</SettingLabel>
          <SettingDescription>{t("settings.setting.sandbox_access.description")}</SettingDescription>
        </div>
        <SettingContent className="gap-2 shrink-0">
          <ToggleGroup type="single" value={mode}>
            <ToggleGroupItem value="local" onClick={() => apply("local", path)}>
              {t("settings.setting.sandbox_access.mode_local")}
            </ToggleGroupItem>
            <ToggleGroupItem value="docker" onClick={() => apply("docker", path)}>
              {t("settings.setting.sandbox_access.mode_docker")}
            </ToggleGroupItem>
          </ToggleGroup>
        </SettingContent>
      </div>

      <div className="flex gap-2 items-center">
        <Input
          value={path}
          placeholder={
            mode === "docker"
              ? t("settings.setting.sandbox_access.placeholder_docker")
              : t("settings.setting.sandbox_access.placeholder_local")
          }
          onChange={(e) => setPath(e.target.value)}
          onBlur={() => apply(mode, path)}
        />
        <Button variant={"outline"} className="shrink-0" onClick={choose} disabled={mode === ""}>
          {t("settings.setting.sandbox_access.choose")}
        </Button>
      </div>

      <div className="text-sm">
        {access?.configured ? (
          <span className="text-muted-foreground">
            {t("settings.setting.sandbox_access.resolved", { file: access.file })}
          </span>
        ) : (
          <span className="text-destructive">
            {access?.reason
              ? t("settings.setting.sandbox_access.not_usable", { reason: access.reason })
              : t("settings.setting.sandbox_access.not_configured")}
          </span>
        )}
      </div>
    </SettingsItem>
  );
}
