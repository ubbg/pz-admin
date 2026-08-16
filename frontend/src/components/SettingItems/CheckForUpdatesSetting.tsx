import { useTranslation } from "react-i18next";
import { SwitchConfig } from "./Presets/SwitchConfig";
import { useUpdateSource } from "@/lib/update-source";

export function CheckForUpdatesSetting() {
  const { t } = useTranslation();
  const source = useUpdateSource();

  return (
    <SwitchConfig
      configKey="checkForUpdates"
      label={t("settings.setting.check_for_updates.label")}
      description={
        source?.configured
          ? t("settings.setting.check_for_updates.description")
          : t("settings.setting.check_for_updates.no_source_description")
      }
      disabled={source !== undefined && !source.configured}
    />
  );
}
