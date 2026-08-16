import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { AddKey } from "@/wailsjs/go/main/App";
import { useEffect, useState } from "react";
import { Input } from "../ui/input";
import { useTranslation } from "react-i18next";

interface AddKeyDialogProps {
  isOpen: boolean;
  onClose: () => void;
  names: string[];
}

export function AddKeyDialog({ isOpen, onClose, names }: AddKeyDialogProps) {
  const { t } = useTranslation();
  const [keyId, setKeyId] = useState("");
  const [keyName, setKeyName] = useState("");

  useEffect(() => {
    setKeyId("");
    setKeyName("");
  }, [isOpen]);

  const handleSubmit = () => {
    onClose();
    AddKey(names, keyId, keyName);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-[28rem] gap-0">
        <DialogHeader>
          <DialogTitle>{t("admin_panel.tabs.players.dialogs.addkey.title")}</DialogTitle>
          <DialogDescription>
            {t("admin_panel.tabs.players.dialogs.addkey.description", { names: names.join(", ") })}
          </DialogDescription>
        </DialogHeader>
        <div className="py-4 space-y-2">
          <div className="space-y-1">
            <Label htmlFor="add-key-id">{t("admin_panel.tabs.players.dialogs.addkey.key_id")}</Label>
            <Input
              id="add-key-id"
              value={keyId}
              placeholder="7295"
              onChange={(e) => setKeyId(e.target.value.replace(/[\\"']/g, ""))}
            />
          </div>
          <div className="space-y-1">
            <Label htmlFor="add-key-name">{t("admin_panel.tabs.players.dialogs.addkey.key_name")}</Label>
            <Input
              id="add-key-name"
              value={keyName}
              placeholder={t("admin_panel.tabs.players.dialogs.addkey.key_name_placeholder")}
              onChange={(e) => setKeyName(e.target.value.replace(/[\\"']/g, ""))}
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant={"outline"} onClick={onClose}>
            {t("cancel")}
          </Button>
          <Button type="submit" onClick={handleSubmit} disabled={keyId.trim().length === 0}>
            {t("admin_panel.tabs.players.dialogs.addkey.submit")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
