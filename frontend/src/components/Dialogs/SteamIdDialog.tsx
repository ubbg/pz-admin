import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { AddSteamId, BanSteamId, RemoveSteamId, UnbanSteamId } from "@/wailsjs/go/main/App";
import { useEffect, useState } from "react";
import { Input } from "../ui/input";
import { useTranslation } from "react-i18next";

export type SteamIdMode = "add" | "remove" | "ban" | "unban";

interface SteamIdDialogProps {
  isOpen: boolean;
  onClose: () => void;
  mode: SteamIdMode;
}

const actions: Record<SteamIdMode, (steamId: string) => Promise<boolean>> = {
  add: AddSteamId,
  remove: RemoveSteamId,
  ban: BanSteamId,
  unban: UnbanSteamId,
};

// Die Erlaubnis- und Sperrlisten des Servers sind über RCON nicht lesbar. Deshalb
// zeigt dieser Dialog keine Liste, sondern führt genau einen Vorgang aus.
export function SteamIdDialog({ isOpen, onClose, mode }: SteamIdDialogProps) {
  const { t } = useTranslation();
  const [steamId, setSteamId] = useState("");

  useEffect(() => {
    setSteamId("");
  }, [isOpen]);

  const handleSubmit = () => {
    onClose();
    actions[mode](steamId);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-[28rem] gap-0">
        <DialogHeader>
          <DialogTitle>{t(`admin_panel.tabs.management.access.steamid.${mode}.title`)}</DialogTitle>
          <DialogDescription>{t(`admin_panel.tabs.management.access.steamid.${mode}.description`)}</DialogDescription>
        </DialogHeader>
        <div className="py-4 space-y-1">
          <Label htmlFor="steam-id">{t("admin_panel.tabs.management.access.steamid.label")}</Label>
          <Input
            id="steam-id"
            value={steamId}
            placeholder="76561198000000000"
            onChange={(e) => setSteamId(e.target.value.replace(/[^0-9A-Za-z:_.-]/g, ""))}
          />
        </div>
        <DialogFooter>
          <Button variant={"outline"} onClick={onClose}>
            {t("cancel")}
          </Button>
          <Button type="submit" onClick={handleSubmit} disabled={steamId.trim().length === 0}>
            {t(`admin_panel.tabs.management.access.steamid.${mode}.submit`)}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
