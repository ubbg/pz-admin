import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Combobox } from "./ui/combobox";
import { SettingContent, SettingDescription, SettingLabel, SettingsGroup, SettingsItem } from "./ui/settings-group";
import { ReloadAllLua, ReloadLua, RemoveZombies, SetLogLevel, SetStats, WorldGen, WorldGenStatus } from "@/wailsjs/go/main/App";
import { useRcon } from "@/contexts/rcon-provider";
import { ConfirmActionDialog } from "./Dialogs/ConfirmActionDialog";

// Protokollbereiche und -stufen laut pzwiki „Admin commands", Fassung 42.20.2
// (/log "Type" "Level"). Versionsgebundene Daten, deshalb hier beisammen.
const logTypes = [
  "General", "Network", "Multiplayer", "Voice", "Packet", "NetworkFileDebug", "Lua", "Mod", "Sound",
  "Zombie", "Combat", "Objects", "Fireplace", "Radio", "MapLoading", "Clothing", "Animation", "Asset",
  "Script", "Shader", "Input", "Recipe", "ActionSystem", "IsoRegion", "UniTests", "FileIO", "Ownership",
  "Death", "Damage", "Statistic", "Vehicle", "Checksum",
];
const logLevels = ["Trace", "Debug", "General", "Warning", "Error"];
const statsModes = ["none", "file", "console", "all"];

// Der Fortschritt des Weltgenerators wird abgefragt, solange diese Ansicht offen ist.
const worldGenPollIntervalMs = 5000;

export function WorldMaintenance() {
  const { t } = useTranslation();
  const { isConnected } = useRcon();

  const [confirmRemoveZombies, setConfirmRemoveZombies] = useState(false);
  const [confirmWorldGenStop, setConfirmWorldGenStop] = useState(false);

  const [worldGenStatus, setWorldGenStatus] = useState("");
  const [luaFile, setLuaFile] = useState("");
  const [logType, setLogType] = useState("General");
  const [logLevel, setLogLevel] = useState("General");
  const [statsMode, setStatsMode] = useState("none");
  const [statsPeriod, setStatsPeriod] = useState(10);

  // Der Zustand wird gelesen, nie angenommen: Was hier steht, ist die Antwort des
  // Servers auf „worldgen status". Nach dem Verlassen der Ansicht wird nicht weiter
  // abgefragt — das Aufräumen des Intervalls erledigt das.
  useEffect(() => {
    if (!isConnected) {
      setWorldGenStatus("");
      return;
    }

    let active = true;

    const readStatus = () => {
      WorldGenStatus().then((response) => {
        if (active) {
          setWorldGenStatus(response.error ? response.error : response.response);
        }
      });
    };

    readStatus();
    const interval = setInterval(readStatus, worldGenPollIntervalMs);

    return () => {
      active = false;
      clearInterval(interval);
    };
  }, [isConnected]);

  return (
    <div className="py-4">
      <h1 className="text-2xl font-semibold leading-none tracking-tight">
        {t("admin_panel.tabs.management.world.title")}
      </h1>
      <h2 className="text-sm text-muted-foreground">{t("admin_panel.tabs.management.world.description")}</h2>

      <SettingsGroup>
        <SettingsItem>
          <div className="flex justify-end flex-col">
            <SettingLabel>{t("admin_panel.tabs.management.world.removezombies.title")}</SettingLabel>
            <SettingDescription>
              {t("admin_panel.tabs.management.world.removezombies.description")}
            </SettingDescription>
          </div>
          <SettingContent>
            <Button variant={"destructive"} className="min-w-40" onClick={() => setConfirmRemoveZombies(true)}>
              {t("admin_panel.tabs.management.world.removezombies.submit")}
            </Button>
          </SettingContent>
        </SettingsItem>

        <SettingsItem vertical>
          <div className="flex justify-between items-start gap-4">
            <div className="flex justify-end flex-col">
              <SettingLabel>{t("admin_panel.tabs.management.world.worldgen.title")}</SettingLabel>
              <SettingDescription>{t("admin_panel.tabs.management.world.worldgen.description")}</SettingDescription>
            </div>
            <SettingContent className="gap-2 shrink-0">
              <Button onClick={() => WorldGen("start")}>
                {t("admin_panel.tabs.management.world.worldgen.start")}
              </Button>
              <Button onClick={() => WorldGen("recheck")}>
                {t("admin_panel.tabs.management.world.worldgen.recheck")}
              </Button>
              <Button variant={"destructive"} onClick={() => setConfirmWorldGenStop(true)}>
                {t("admin_panel.tabs.management.world.worldgen.stop")}
              </Button>
            </SettingContent>
          </div>
          <div className="rounded-md border bg-black/10 dark:bg-white/5 p-2 font-mono text-xs whitespace-pre-wrap">
            {worldGenStatus || t("admin_panel.tabs.management.world.worldgen.no_status")}
          </div>
        </SettingsItem>

        <SettingsItem>
          <div className="flex justify-end flex-col">
            <SettingLabel>{t("admin_panel.tabs.management.world.reloadlua.title")}</SettingLabel>
            <SettingDescription>{t("admin_panel.tabs.management.world.reloadlua.description")}</SettingDescription>
          </div>
          <SettingContent className="gap-2">
            <Input
              className="w-52"
              value={luaFile}
              placeholder={t("admin_panel.tabs.management.world.reloadlua.placeholder")}
              onChange={(e) => setLuaFile(e.target.value.replace(/[\\"']/g, ""))}
            />
            <Button className="min-w-28" disabled={luaFile.trim().length === 0} onClick={() => ReloadLua(luaFile)}>
              {t("admin_panel.tabs.management.world.reloadlua.submit")}
            </Button>
            <Button className="min-w-40" onClick={() => ReloadAllLua()}>
              {t("admin_panel.tabs.management.world.reloadlua.submit_all")}
            </Button>
          </SettingContent>
        </SettingsItem>

        <SettingsItem>
          <div className="flex justify-end flex-col">
            <SettingLabel>{t("admin_panel.tabs.management.world.log.title")}</SettingLabel>
            <SettingDescription>{t("admin_panel.tabs.management.world.log.description")}</SettingDescription>
          </div>
          <SettingContent className="gap-2">
            <Combobox
              mandatory
              initialValue={logType}
              elements={logTypes.map((type) => ({ value: type, label: type }))}
              placeholder={t("admin_panel.tabs.management.world.log.type")}
              searchPlaceholder={t("admin_panel.tabs.management.world.log.type")}
              nothingFoundMessage={t("admin_panel.tabs.management.world.log.nothing_found")}
              onChange={(value) => setLogType(value)}
            />
            <Combobox
              mandatory
              disableSearch
              initialValue={logLevel}
              elements={logLevels.map((level) => ({ value: level, label: level }))}
              placeholder={t("admin_panel.tabs.management.world.log.level")}
              nothingFoundMessage={t("admin_panel.tabs.management.world.log.nothing_found")}
              onChange={(value) => setLogLevel(value)}
            />
            <Button className="min-w-28" onClick={() => SetLogLevel(logType, logLevel)}>
              {t("admin_panel.tabs.management.world.log.submit")}
            </Button>
          </SettingContent>
        </SettingsItem>

        <SettingsItem className="border-none">
          <div className="flex justify-end flex-col">
            <SettingLabel>{t("admin_panel.tabs.management.world.stats.title")}</SettingLabel>
            <SettingDescription>{t("admin_panel.tabs.management.world.stats.description")}</SettingDescription>
          </div>
          <SettingContent className="gap-2">
            <Combobox
              mandatory
              disableSearch
              initialValue={statsMode}
              elements={statsModes.map((mode) => ({
                value: mode,
                label: t(`admin_panel.tabs.management.world.stats.modes.${mode}`),
              }))}
              placeholder={t("admin_panel.tabs.management.world.stats.mode")}
              nothingFoundMessage={t("admin_panel.tabs.management.world.stats.nothing_found")}
              onChange={(value) => setStatsMode(value)}
            />
            <Input
              className="w-24"
              type="number"
              min={1}
              value={statsPeriod}
              disabled={statsMode === "none"}
              onChange={(e) => setStatsPeriod(Math.max(1, parseInt(e.target.value, 10) || 1))}
            />
            <Button className="min-w-28" onClick={() => SetStats(statsMode, statsPeriod)}>
              {t("admin_panel.tabs.management.world.stats.submit")}
            </Button>
          </SettingContent>
        </SettingsItem>
      </SettingsGroup>

      <ConfirmActionDialog
        isOpen={confirmRemoveZombies}
        onClose={() => setConfirmRemoveZombies(false)}
        title={t("admin_panel.tabs.management.world.removezombies.confirm_title")}
        description={t("admin_panel.tabs.management.world.removezombies.confirm")}
        confirmText={t("admin_panel.tabs.management.world.removezombies.confirm_submit")}
        onConfirm={() => RemoveZombies()}
      />

      <ConfirmActionDialog
        isOpen={confirmWorldGenStop}
        onClose={() => setConfirmWorldGenStop(false)}
        title={t("admin_panel.tabs.management.world.worldgen.stop_confirm_title")}
        description={t("admin_panel.tabs.management.world.worldgen.stop_confirm")}
        confirmText={t("admin_panel.tabs.management.world.worldgen.stop_confirm_submit")}
        onConfirm={() => WorldGen("stop")}
      />
    </div>
  );
}
