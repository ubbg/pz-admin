import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { useEffect, useState } from "react";
import { Input } from "../ui/input";
import { CreateHorde, CreateHorde2 } from "@/wailsjs/go/main/App";
import { ToggleGroup, ToggleGroupItem } from "../ui/toggle-group";
import { useTranslation } from "react-i18next";

interface CreateHordeDialogProps {
  isOpen: boolean;
  onClose: () => void;
  names: string[];
}

export function CreateHordeDialog({ isOpen, onClose, names }: CreateHordeDialogProps) {
  const { t } = useTranslation();
  const [count, setCount] = useState("");
  // Build 42 kennt zwei Horden-Befehle; welcher gemeint ist, entscheidet die Nutzerin.
  const [variant, setVariant] = useState<"createhorde" | "createhorde2">("createhorde");

  const handleCreateHorde = () => {
    onClose();

    // Check if count is a number
    if (isNaN(parseInt(count))) {
      return;
    }

    if (variant === "createhorde2") {
      CreateHorde2(names, parseInt(count));
    } else {
      CreateHorde(names, parseInt(count));
    }
  };

  useEffect(() => {
    setCount("");
    setVariant("createhorde");
  }, [isOpen]);

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-[28rem]">
        <DialogHeader>
          <DialogTitle>{t("admin_panel.tabs.players.dialogs.createhorde.title")}</DialogTitle>
          <DialogDescription>
            <p>{t("admin_panel.tabs.players.dialogs.createhorde.players", { players: names.join(", ") })}</p>
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-1 pb-2">
          <Label className="text-right">{t("admin_panel.tabs.players.dialogs.createhorde.variant")}</Label>
          <ToggleGroup type="single" value={variant} className="justify-start">
            <ToggleGroupItem value="createhorde" onClick={() => setVariant("createhorde")}>
              {t("admin_panel.tabs.players.dialogs.createhorde.variant_1")}
            </ToggleGroupItem>
            <ToggleGroupItem value="createhorde2" onClick={() => setVariant("createhorde2")}>
              {t("admin_panel.tabs.players.dialogs.createhorde.variant_2")}
            </ToggleGroupItem>
          </ToggleGroup>
        </div>
        <div className="space-y-1">
          <Label htmlFor="horde-size" className="text-right">
            {t("admin_panel.tabs.players.dialogs.createhorde.horde_size")}
          </Label>
          <Input
            value={count}
            onChange={(e) => {
              const parsedValue = parseInt(e.target.value);

              if (!isNaN(parsedValue)) {
                setCount(Math.max(0, Math.min(parsedValue, 2147483647)).toString());
              } else {
                setCount(e.target.value);
              }
            }}
            min={0}
            max={2147483647}
            id="horde-size"
            type="number"
            placeholder="150"
          />
        </div>
        <DialogFooter>
          <Button
            type="submit"
            onClick={handleCreateHorde}
            disabled={count === "" || isNaN(parseInt(count)) || parseInt(count) < 0 || parseFloat(count) % 1 !== 0}
          >
            {t("admin_panel.tabs.players.dialogs.createhorde.submit")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
