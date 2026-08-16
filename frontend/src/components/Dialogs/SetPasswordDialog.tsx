import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { SetUserPassword } from "@/wailsjs/go/main/App";
import { useEffect, useState } from "react";
import { Input } from "../ui/input";
import { useTranslation } from "react-i18next";

interface SetPasswordDialogProps {
  isOpen: boolean;
  onClose: () => void;
  initialName?: string;
}

// Das Passwort verlässt diesen Dialog nur Richtung Go; es steht in keiner Meldung
// und in keinem Protokoll.
export function SetPasswordDialog({ isOpen, onClose, initialName }: SetPasswordDialogProps) {
  const { t } = useTranslation();
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");

  useEffect(() => {
    setName(initialName ?? "");
    setPassword("");
  }, [isOpen, initialName]);

  const handleSubmit = () => {
    onClose();
    SetUserPassword(name, password);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-[28rem] gap-0">
        <DialogHeader>
          <DialogTitle>{t("admin_panel.tabs.management.access.setpassword.title")}</DialogTitle>
          <DialogDescription>{t("admin_panel.tabs.management.access.setpassword.description")}</DialogDescription>
        </DialogHeader>
        <div className="py-4 space-y-2">
          <div className="space-y-1">
            <Label htmlFor="set-password-name">{t("admin_panel.tabs.management.access.setpassword.name")}</Label>
            <Input
              id="set-password-name"
              value={name}
              placeholder="John Doe"
              onChange={(e) => setName(e.target.value.replace(/[\\"']/g, ""))}
            />
          </div>
          <div className="space-y-1">
            <Label htmlFor="set-password-value">{t("admin_panel.tabs.management.access.setpassword.password")}</Label>
            <Input
              id="set-password-value"
              type="password"
              autoComplete="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value.replace(/[\\"']/g, ""))}
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant={"outline"} onClick={onClose}>
            {t("cancel")}
          </Button>
          <Button type="submit" onClick={handleSubmit} disabled={name.trim().length < 3 || password.length < 1}>
            {t("admin_panel.tabs.management.access.setpassword.submit")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
