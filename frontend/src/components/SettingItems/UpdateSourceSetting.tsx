import { useTranslation } from "react-i18next";
import { SettingContent, SettingDescription, SettingLabel, SettingsItem } from "@/components/ui/settings-group";
import { Button } from "@/components/ui/button";
import { BrowserOpenURL } from "@/wailsjs/runtime/runtime";
import { useUpdateSource } from "@/lib/update-source";

export function UpdateSourceSetting() {
  const { t } = useTranslation();
  const source = useUpdateSource();

  return (
    <SettingsItem loading={source === undefined}>
      <div>
        <SettingLabel>{t("settings.setting.update_source.label")}</SettingLabel>
        <SettingDescription>
          {source?.configured
            ? t("settings.setting.update_source.description")
            : t("settings.setting.update_source.not_configured_description")}
        </SettingDescription>
      </div>
      <SettingContent>
        {source?.configured ? (
          <Button variant={"link"} className="text-foreground" onClick={() => BrowserOpenURL(source.url)}>
            {source.owner}/{source.repo}
          </Button>
        ) : (
          <span className="text-sm text-muted-foreground">{t("settings.setting.update_source.not_configured")}</span>
        )}
      </SettingContent>
    </SettingsItem>
  );
}
