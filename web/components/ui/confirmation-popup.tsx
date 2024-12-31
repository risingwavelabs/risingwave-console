import { Button } from "@/components/ui/button"

interface ConfirmationPopupProps {
  message: string
  onConfirm: () => void
  onCancel: () => void
}

export function ConfirmationPopup({ message, onConfirm, onCancel }: ConfirmationPopupProps) {
  return (
    <div className="absolute right-0 top-full mt-1 bg-background border rounded-lg shadow-lg py-2 px-3 z-50 whitespace-nowrap">
      <p className="text-xs text-muted-foreground mb-2">{message}</p>
      <div className="flex gap-2">
        <Button 
          size="sm" 
          variant="ghost"
          className="h-7 px-2 text-xs"
          onClick={onCancel}
        >
          Cancel
        </Button>
        <Button 
          size="sm"
          variant="destructive"
          className="h-7 px-2 text-xs"
          onClick={onConfirm}
        >
          Delete
        </Button>
      </div>
    </div>
  )
} 