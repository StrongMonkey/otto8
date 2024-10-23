import { PlusIcon, SearchIcon } from "lucide-react";
import { useState } from "react";
import useSWR, { preload } from "swr";

import { ToolReferenceService } from "~/lib/service/api/toolreferenceService";

import { TypographyH2 } from "~/components/Typography";
import { CreateTool } from "~/components/tools/CreateTool";
import { ToolGrid } from "~/components/tools/toolGrid";
import { Button } from "~/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "~/components/ui/dialog";
import { Input } from "~/components/ui/input";

export async function clientLoader() {
    await Promise.all([
        preload(ToolReferenceService.getToolReferences.key("tool"), () =>
            ToolReferenceService.getToolReferences("tool")
        ),
    ]);
    return null;
}

export default function Tools() {
    const { data: tools, mutate } = useSWR(
        ToolReferenceService.getToolReferences.key("tool"),
        () => ToolReferenceService.getToolReferences("tool")
    );

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [searchQuery, setSearchQuery] = useState("");

    const handleCreateSuccess = () => {
        mutate();
        setIsDialogOpen(false);
    };

    const handleDelete = async (id: string) => {
        await ToolReferenceService.deleteToolReference(id);
        mutate();
    };

    return (
        <div className="h-full p-8 flex flex-col gap-4">
            <div className="flex justify-between items-center">
                <TypographyH2>Tools</TypographyH2>
                <div className="flex items-center space-x-2">
                    <div className="relative">
                        <SearchIcon className="w-5 h-5 text-gray-400 absolute left-3 top-1/2 transform -translate-y-1/2" />
                        <Input
                            type="text"
                            placeholder="Search for tools..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="pl-10 w-64"
                        />
                    </div>
                    <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                        <DialogTrigger asChild>
                            <Button variant="outline">
                                <PlusIcon className="w-4 h-4 mr-2" />
                                Register New Tool
                            </Button>
                        </DialogTrigger>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>
                                    Create New Tool Reference
                                </DialogTitle>
                            </DialogHeader>
                            <CreateTool onSuccess={handleCreateSuccess} />
                        </DialogContent>
                    </Dialog>
                </div>
            </div>

            {tools && (
                <ToolGrid
                    tools={tools}
                    filter={searchQuery}
                    onDelete={handleDelete}
                />
            )}
        </div>
    );
}