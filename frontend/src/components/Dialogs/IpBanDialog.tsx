import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { BanIp, UnbanIp } from "@/wailsjs/go/main/App";
import { useEffect, useState } from "react";
import { Input } from "../ui/input";
import { useTranslation } from "react-i18next";

interface IpBanDialogProps {
  isOpen: boolean;
  onClose: () => void;
  mode: "ban" | "unban";
}

// Ein IP-Bann trifft alle, die hinter derselben Adresse spielen. Vor dem Sperren
// steht deshalb eine Rückfrage, die die Adresse und die Wirkung benennt (PG-10).
export function IpBanDialog({ isOpen, onClose, mode }: IpBanDialogProps) {
  const { t } = useTranslation();
  const [ip, setIp] = useState("");
  const [confirming, setConfirming] = useState(false);

  useEffect(() => {
    setIp("");
    setConfirming(false);
  }, [isOpen]);

  const handleSubmit = () => {
    if (mode === "ban" && !confirming) {
      setConfirming(true);
      return;
    }

    onClose();
    if (mode === "ban") {
      BanIp(ip);
    } else {
      UnbanIp(ip);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-[28rem] gap-0">
        <DialogHeader>
          <DialogTitle>{t(`admin_panel.tabs.management.moderation.ip.${mode}.title`)}</DialogTitle>
          <DialogDescription>
            {confirming
              ? t("admin_panel.tabs.management.moderation.ip.ban.confirm", { ip: ip.trim() })
              : t(`admin_panel.tabs.management.moderation.ip.${mode}.description`)}
          </DialogDescription>
        </DialogHeader>
        <div className="py-4 space-y-1">
          <Label htmlFor="ban-ip">{t("admin_panel.tabs.management.moderation.ip.label")}</Label>
          <Input
            id="ban-ip"
            value={ip}
            placeholder="203.0.113.5"
            disabled={confirming}
            onChange={(e) => setIp(e.target.value.replace(/[^0-9a-fA-F:.]/g, ""))}
          />
        </div>
        <DialogFooter>
          <Button variant={"outline"} onClick={confirming ? () => setConfirming(false) : onClose}>
            {confirming ? t("back") : t("cancel")}
          </Button>
          <Button
            type="submit"
            variant={mode === "ban" ? "destructive" : "default"}
            onClick={handleSubmit}
            disabled={ip.trim().length === 0}
          >
            {confirming
              ? t("admin_panel.tabs.management.moderation.ip.ban.confirm_submit")
              : t(`admin_panel.tabs.management.moderation.ip.${mode}.submit`)}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
