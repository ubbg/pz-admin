import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "./ui/button";
import { SettingContent, SettingDescription, SettingLabel, SettingsGroup, SettingsItem } from "./ui/settings-group";
import { ReleaseSafehouse } from "@/wailsjs/go/main/App";
import { RemoveItemDialog } from "./Dialogs/RemoveItemDialog";

// Zwei Befehle von Build 42 kennen keinen Spielernamen: releasesafehouse und
// removeitem wirken auf die Figur der Konsole. Sie stehen deshalb hier und nicht an
// einer Spielerzeile — und die Beschreibung sagt, was das auf einem Dedicated Server
// bedeutet.
export function ConsoleActions() {
  const { t } = useTranslation();
  const [isRemoveItemDialogOpen, setRemoveItemDialogOpen] = useState(false);

  return (
    <div className="py-4">
      <h1 className="text-2xl font-semibold leading-none tracking-tight">
        {t("admin_panel.tabs.management.console.title")}
      </h1>
      <h2 className="text-sm text-muted-foreground">{t("admin_panel.tabs.management.console.description")}</h2>

      <SettingsGroup>
        <SettingsItem>
          <div className="flex justify-end flex-col">
            <SettingLabel>{t("admin_panel.tabs.management.console.releasesafehouse.title")}</SettingLabel>
            <SettingDescription>
              {t("admin_panel.tabs.management.console.releasesafehouse.description")}
            </SettingDescription>
          </div>
          <SettingContent>
            <Button className="min-w-40" onClick={() => ReleaseSafehouse()}>
              {t("admin_panel.tabs.management.console.releasesafehouse.submit")}
            </Button>
          </SettingContent>
        </SettingsItem>

        <SettingsItem className="border-none">
          <div className="flex justify-end flex-col">
            <SettingLabel>{t("admin_panel.tabs.management.console.removeitem.title")}</SettingLabel>
            <SettingDescription>{t("admin_panel.tabs.management.console.removeitem.description")}</SettingDescription>
          </div>
          <SettingContent>
            <Button className="min-w-40" onClick={() => setRemoveItemDialogOpen(true)}>
              {t("admin_panel.tabs.management.console.removeitem.submit")}
            </Button>
          </SettingContent>
        </SettingsItem>
      </SettingsGroup>

      <RemoveItemDialog isOpen={isRemoveItemDialogOpen} onClose={() => setRemoveItemDialogOpen(false)} />
    </div>
  );
}
