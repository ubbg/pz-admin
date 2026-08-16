import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { useTranslation } from "react-i18next";

interface ConfirmActionDialogProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  description: string;
  confirmText: string;
  onConfirm: () => void;
}

// Rückfrage für nicht zurücknehmbare Befehle. Titel und Beschreibung benennen Wirkung
// und Umfang — ein „Sind Sie sicher?" sagt der Nutzerin nichts (PG-10).
export function ConfirmActionDialog({
  isOpen,
  onClose,
  title,
  description,
  confirmText,
  onConfirm,
}: ConfirmActionDialogProps) {
  const { t } = useTranslation();

  const handleConfirm = () => {
    onClose();
    onConfirm();
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-[28rem]">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant={"outline"} onClick={onClose}>
            {t("cancel")}
          </Button>
          <Button variant={"destructive"} type="submit" onClick={handleConfirm}>
            {confirmText}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
