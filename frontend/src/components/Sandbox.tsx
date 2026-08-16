import { useCallback, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Switch } from "./ui/switch";
import { ScrollArea } from "./ui/scroll-area";
import { SettingContent, SettingLabel, SettingsGroup, SettingsItem } from "./ui/settings-group";
import { Loader2, Search, TriangleAlert } from "lucide-react";
import { ReadSandboxVars, SandboxAccess, SaveSandboxVars } from "@/wailsjs/go/main/App";
import { main } from "@/wailsjs/go/models";
import { EventsOff, EventsOn } from "@/wailsjs/runtime/runtime";
import { useConfig } from "@/contexts/config-provider";
import { optionValueToString } from "@/lib/options";

// Der Sandbox-Bereich arbeitet auf der Datei des Servers, nicht über RCON — er
// funktioniert deshalb auch bei getrennter Verbindung. Die Variablen sind nirgends im
// Code deklariert: Angezeigt wird, was in der Datei steht.
export function SandboxTab() {
  const { t } = useTranslation();
  const { config } = useConfig();

  const [access, setAccess] = useState<main.SandboxAccess>();
  const [document, setDocument] = useState<main.SandboxDocument>();
  const [values, setValues] = useState<Record<string, string>>({});
  const [saving, setSaving] = useState(false);
  const [searchText, setSearchText] = useState("");

  const load = useCallback(() => {
    SandboxAccess().then((value) => {
      setAccess(value);
      if (!value.configured) {
        setDocument(undefined);
        return;
      }
      ReadSandboxVars().then((read) => {
        setDocument(read);
        const next: Record<string, string> = {};
        for (const variable of read.vars ?? []) {
          next[variable.key] = variable.value;
        }
        setValues(next);
      });
    });
  }, []);

  useEffect(() => {
    load();
    EventsOn("sandbox-access-changed", load);
    return () => EventsOff("sandbox-access-changed");
  }, [load]);

  const vars = document?.vars ?? [];

  const modified = useMemo(
    () => vars.filter((variable) => values[variable.key] !== variable.value).map((variable) => variable.key),
    [vars, values]
  );

  const groups = useMemo(() => {
    const search = searchText.toLowerCase();
    const filtered = vars.filter((variable) => variable.key.toLowerCase().includes(search));
    const byGroup = new Map<string, main.SandboxVar[]>();
    for (const variable of filtered) {
      const list = byGroup.get(variable.group) ?? [];
      list.push(variable);
      byGroup.set(variable.group, list);
    }
    return [...byGroup.entries()];
  }, [vars, searchText]);

  const save = () => {
    setSaving(true);
    SaveSandboxVars(values).finally(() => {
      setSaving(false);
      load();
    });
  };

  const height = config?.useSystemTitleBar ? "h-[calc(100vh-3.5rem)]" : "h-[calc(100vh-5.5rem)]";

  // Ohne eingerichteten Zugangsweg ist der Bereich inaktiv und nennt den Grund — er
  // ist nicht leer, zeigt keine Beispieldaten und rät keinen Pfad (PG-20).
  if (access && !access.configured) {
    return (
      <div className={`w-full ${height} dark:bg-black/20 bg-white/20 p-6`}>
        <div className="max-w-2xl space-y-2">
          <h1 className="text-2xl font-semibold leading-none tracking-tight">{t("admin_panel.tabs.sandbox.title")}</h1>
          <p className="text-sm text-muted-foreground">{t("admin_panel.tabs.sandbox.not_configured")}</p>
          {access.reason && <p className="text-sm text-destructive">{access.reason}</p>}
        </div>
      </div>
    );
  }

  return (
    <div className={`w-full ${height} dark:bg-black/20 bg-white/20 p-2`}>
      <div className="flex items-center gap-2 m-1">
        <div className="relative w-full">
          <Input
            className="peer ps-9"
            type="search"
            placeholder={t("admin_panel.tabs.sandbox.search")}
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
          />
          <div className="pointer-events-none absolute inset-y-0 start-0 flex items-center justify-center ps-3 text-muted-foreground/80">
            <Search className="w-4 h-4" strokeWidth={2} />
          </div>
        </div>
        <span className="text-xs text-muted-foreground whitespace-nowrap">{document?.file}</span>
      </div>

      {document?.error && <p className="m-1 text-sm text-destructive">{document.error}</p>}

      <ScrollArea className="w-full pr-6" style={{ height: "calc(100% - 6.5rem)" }}>
        {groups.map(([group, variables]) => (
          <div key={group || "root"} className="flex flex-col gap-2">
            <div className="border-b pb-2 mt-6">
              <span className="text-2xl font-bold leading-none">
                {group || t("admin_panel.tabs.sandbox.general_group")}
              </span>
            </div>
            <SettingsGroup className="space-y-0.5">
              {variables.map((variable) => (
                <SettingsItem className="border-none" key={variable.key}>
                  <SettingLabel className="font-mono text-sm">{variable.name}</SettingLabel>
                  <SettingContent>
                    <SandboxValueInput
                      variable={variable}
                      value={values[variable.key] ?? variable.value}
                      onChange={(value) => setValues((previous) => ({ ...previous, [variable.key]: value }))}
                    />
                  </SettingContent>
                </SettingsItem>
              ))}
            </SettingsGroup>
          </div>
        ))}
      </ScrollArea>

      {/* Der Hinweis steht am Speichern-Knopf, nicht in einer Fußnote (PG-19). */}
      <div className="h-16 flex items-center justify-between mr-5 gap-4">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <TriangleAlert className="w-4 h-4 shrink-0" />
          {t("admin_panel.tabs.sandbox.restart_required")}
        </div>
        <div className="flex gap-2 items-center">
          <span className="text-sm text-muted-foreground">
            {t("admin_panel.tabs.sandbox.modified", { count: modified.length })}
          </span>
          <Button onClick={save} disabled={saving || modified.length === 0} className="min-w-28">
            {saving && <Loader2 className="mr-2 animate-spin" />}
            {t("save")}
          </Button>
        </div>
      </div>
    </div>
  );
}

function SandboxValueInput({
  variable,
  value,
  onChange,
}: {
  variable: main.SandboxVar;
  value: string;
  onChange: (value: string) => void;
}) {
  if (variable.kind === "Boolean") {
    return (
      <Switch
        checked={value.toLowerCase() === "true"}
        onCheckedChange={(checked) => onChange(optionValueToString(checked, "Boolean"))}
      />
    );
  }

  if (variable.kind === "Integer" || variable.kind === "Double") {
    return (
      <Input
        className="w-[8rem]"
        type="number"
        lang="en"
        step={variable.kind === "Double" ? 0.1 : 1}
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
    );
  }

  return <Input className="w-[20rem]" type="text" value={value} onChange={(e) => onChange(e.target.value)} />;
}
