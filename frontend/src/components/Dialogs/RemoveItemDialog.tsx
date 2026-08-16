import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { RemoveItems } from "@/wailsjs/go/main/App";
import { main } from "@/wailsjs/go/models";
import { useEffect, useState } from "react";
import { Input } from "../ui/input";
import { useTranslation } from "react-i18next";

interface RemoveItemDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

// removeitem kennt in Build 42 keinen Spielernamen — der Befehl räumt bei der Figur
// der Konsole auf. Deshalb hängt dieser Dialog nicht an einer Spielerzeile.
export function RemoveItemDialog({ isOpen, onClose }: RemoveItemDialogProps) {
  const { t } = useTranslation();
  const [itemId, setItemId] = useState("");
  const [count, setCount] = useState(1);

  useEffect(() => {
    setItemId("");
    setCount(1);
  }, [isOpen]);

  const handleSubmit = () => {
    onClose();
    RemoveItems([main.ItemRecord.createFrom({ itemId, count })]);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-[28rem] gap-0">
        <DialogHeader>
          <DialogTitle>{t("admin_panel.tabs.management.console.removeitem.title")}</DialogTitle>
          <DialogDescription>{t("admin_panel.tabs.management.console.removeitem.description")}</DialogDescription>
        </DialogHeader>
        <div className="py-4 space-y-2">
          <div className="space-y-1">
            <Label htmlFor="remove-item-id">{t("admin_panel.tabs.management.console.removeitem.item")}</Label>
            <Input
              id="remove-item-id"
              value={itemId}
              placeholder="Base.Axe"
              onChange={(e) => setItemId(e.target.value.replace(/[\\"' ]/g, ""))}
            />
          </div>
          <div className="space-y-1">
            <Label htmlFor="remove-item-count">{t("admin_panel.tabs.management.console.removeitem.count")}</Label>
            <Input
              id="remove-item-count"
              type="number"
              min={0}
              value={count}
              onChange={(e) => setCount(Math.max(0, parseInt(e.target.value, 10) || 0))}
            />
            <p className="text-xs text-muted-foreground">
              {t("admin_panel.tabs.management.console.removeitem.count_hint")}
            </p>
          </div>
        </div>
        <DialogFooter>
          <Button variant={"outline"} onClick={onClose}>
            {t("cancel")}
          </Button>
          <Button type="submit" onClick={handleSubmit} disabled={itemId.trim().length === 0}>
            {t("admin_panel.tabs.management.console.removeitem.submit")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
