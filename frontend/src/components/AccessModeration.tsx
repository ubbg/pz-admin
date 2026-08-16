import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "./ui/button";
import { SettingContent, SettingDescription, SettingLabel, SettingsGroup, SettingsItem } from "./ui/settings-group";
import { SteamIdDialog, SteamIdMode } from "./Dialogs/SteamIdDialog";
import { IpBanDialog } from "./Dialogs/IpBanDialog";
import { SetPasswordDialog } from "./Dialogs/SetPasswordDialog";

// Zugang und Moderation: Vorgänge, die an keiner Spielerzeile hängen, weil sie eine
// SteamID, eine IP oder ein Konto betreffen — bisher nur über die Serverkonsole
// erreichbar.
export function AccessModeration() {
  const { t } = useTranslation();

  const [steamIdMode, setSteamIdMode] = useState<SteamIdMode>();
  const [ipMode, setIpMode] = useState<"ban" | "unban">();
  const [isPasswordDialogOpen, setPasswordDialogOpen] = useState(false);

  return (
    <div className="py-4">
      <h1 className="text-2xl font-semibold leading-none tracking-tight">
        {t("admin_panel.tabs.management.access.title")}
      </h1>
      <h2 className="text-sm text-muted-foreground">{t("admin_panel.tabs.management.access.description")}</h2>

      <SettingsGroup>
        <SettingsItem>
          <div className="flex justify-end flex-col">
            <SettingLabel>{t("admin_panel.tabs.management.access.steamid.title")}</SettingLabel>
            <SettingDescription>{t("admin_panel.tabs.management.access.steamid.hint")}</SettingDescription>
          </div>
          <SettingContent className="gap-2">
            <Button className="min-w-40" onClick={() => setSteamIdMode("add")}>
              {t("admin_panel.tabs.management.access.steamid.add.submit")}
            </Button>
            <Button className="min-w-40" onClick={() => setSteamIdMode("remove")}>
              {t("admin_panel.tabs.management.access.steamid.remove.submit")}
            </Button>
          </SettingContent>
        </SettingsItem>

        <SettingsItem>
          <div className="flex justify-end flex-col">
            <SettingLabel>{t("admin_panel.tabs.management.access.setpassword.title")}</SettingLabel>
            <SettingDescription>{t("admin_panel.tabs.management.access.setpassword.description")}</SettingDescription>
          </div>
          <SettingContent>
            <Button className="min-w-40" onClick={() => setPasswordDialogOpen(true)}>
              {t("admin_panel.tabs.management.access.setpassword.submit")}
            </Button>
          </SettingContent>
        </SettingsItem>

        <SettingsItem>
          <div className="flex justify-end flex-col">
            <SettingLabel>{t("admin_panel.tabs.management.moderation.ip.title")}</SettingLabel>
            <SettingDescription>{t("admin_panel.tabs.management.moderation.ip.hint")}</SettingDescription>
          </div>
          <SettingContent className="gap-2">
            <Button variant={"destructive"} className="min-w-40" onClick={() => setIpMode("ban")}>
              {t("admin_panel.tabs.management.moderation.ip.ban.submit")}
            </Button>
            <Button className="min-w-40" onClick={() => setIpMode("unban")}>
              {t("admin_panel.tabs.management.moderation.ip.unban.submit")}
            </Button>
          </SettingContent>
        </SettingsItem>

        <SettingsItem className="border-none">
          <div className="flex justify-end flex-col">
            <SettingLabel>{t("admin_panel.tabs.management.moderation.steamid.title")}</SettingLabel>
            <SettingDescription>{t("admin_panel.tabs.management.moderation.steamid.hint")}</SettingDescription>
          </div>
          <SettingContent className="gap-2">
            <Button variant={"destructive"} className="min-w-40" onClick={() => setSteamIdMode("ban")}>
              {t("admin_panel.tabs.management.access.steamid.ban.submit")}
            </Button>
            <Button className="min-w-40" onClick={() => setSteamIdMode("unban")}>
              {t("admin_panel.tabs.management.access.steamid.unban.submit")}
            </Button>
          </SettingContent>
        </SettingsItem>
      </SettingsGroup>

      <SteamIdDialog
        isOpen={steamIdMode !== undefined}
        mode={steamIdMode ?? "add"}
        onClose={() => setSteamIdMode(undefined)}
      />
      <IpBanDialog isOpen={ipMode !== undefined} mode={ipMode ?? "ban"} onClose={() => setIpMode(undefined)} />
      <SetPasswordDialog isOpen={isPasswordDialogOpen} onClose={() => setPasswordDialogOpen(false)} />
    </div>
  );
}
