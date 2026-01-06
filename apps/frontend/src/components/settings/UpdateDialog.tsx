import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
    DialogDescription
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { useUpdateStore } from "@/stores/useUpdateStore";
import { ScrollArea } from "@/components/ui/scroll-area";
import ReactMarkdown from 'react-markdown';

interface UpdateDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export function UpdateDialog({ open, onOpenChange }: UpdateDialogProps) {
    const { updateInfo, openDownloadPage } = useUpdateStore();

    if (!updateInfo) return null;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[600px]">
                <DialogHeader>
                    <DialogTitle>Update Available</DialogTitle>
                    <DialogDescription>
                        A new version ({updateInfo.latest_version}) is available.
                        You are currently on {updateInfo.current_version}.
                    </DialogDescription>
                </DialogHeader>

                <div className="py-4">
                    <h4 className="mb-2 text-sm font-medium">Release Notes</h4>
                    <div className="rounded-md border p-4 bg-muted/50 h-[300px]">
                        <ScrollArea className="h-full">
                            <div className="text-sm prose dark:prose-invert max-w-none">
                                <ReactMarkdown>{updateInfo.release_notes || "No release notes provided."}</ReactMarkdown>
                            </div>
                        </ScrollArea>
                    </div>
                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Cancel
                    </Button>
                    <Button onClick={openDownloadPage}>
                        Download Update
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
