import { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Plus, Trash2, FileText } from 'lucide-react';

// import { middleware } from '../../wailsjs/go/models';
import { AddMapLocalRule, DeleteMapLocalRule, GetMapLocalRules } from '../../wailsjs/go/main/App';

// Temporary fix for build: module resolution failing for middleware namespace
export namespace middleware {
    export class MapLocalRule {
        id: string;
        name: string;
        url_regex: string;
        local_path: string;
        content_type: string;
        status: number;
        enabled: boolean;

        constructor(source: any = {}) {
            if ('string' === typeof source) source = JSON.parse(source);
            this.id = source["id"] || "";
            this.name = source["name"] || "";
            this.url_regex = source["url_regex"] || "";
            this.local_path = source["local_path"] || "";
            this.content_type = source["content_type"] || "";
            this.status = source["status"] || 0;
            this.enabled = source["enabled"] || false;
        }
    }
}

export default function MapLocalPage() {
    const [rules, setRules] = useState<middleware.MapLocalRule[]>([]);
    const [isOpen, setIsOpen] = useState(false);

    // Form State
    const [regex, setRegex] = useState('');
    const [path, setPath] = useState('');
    const [contentType, setContentType] = useState('');

    const fetchRules = async () => {
        try {
            const r = await GetMapLocalRules();
            setRules(r || []);
        } catch (error) {
            console.error("Failed to fetch rules", error);
        }
    };

    useEffect(() => {
        fetchRules();
    }, []);

    const handleAdd = async () => {
        const newRule = new middleware.MapLocalRule({
            id: crypto.randomUUID(),
            enabled: true,
            url_regex: regex,
            local_path: path,
            content_type: contentType
        });

        try {
            await AddMapLocalRule(newRule);
            setIsOpen(false);
            setRegex('');
            setPath('');
            setContentType('');
            fetchRules();
        } catch (error) {
            console.error("Failed to add rule", error);
            alert("Failed to add rule: " + error);
        }
    };

    const handleDelete = async (id: string) => {
        try {
            await DeleteMapLocalRule(id);
            fetchRules();
        } catch (error) {
            console.error("Failed to delete rule", error);
        }
    };

    return (
        <div className="container mx-auto p-6 space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Map Local</h1>
                    <p className="text-muted-foreground">Serve local files instead of upstream responses</p>
                </div>
                <Dialog open={isOpen} onOpenChange={setIsOpen}>
                    <DialogTrigger asChild>
                        <Button>
                            <Plus className="mr-2 h-4 w-4" />
                            Add Rule
                        </Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Add Map Local Rule</DialogTitle>
                            <DialogDescription>
                                Intercepts URLs matching the Regex and serves the file at Local Path.
                            </DialogDescription>
                        </DialogHeader>
                        <div className="grid gap-4 py-4">
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="regex" className="text-right">
                                    URL Regex
                                </Label>
                                <Input
                                    id="regex"
                                    placeholder=".*\/api\/v1\/users"
                                    className="col-span-3 font-mono text-sm"
                                    value={regex}
                                    onChange={(e) => setRegex(e.target.value)}
                                />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="path" className="text-right">
                                    Local Path
                                </Label>
                                <Input
                                    id="path"
                                    placeholder="/Users/me/mocks/users.json"
                                    className="col-span-3 font-mono text-sm"
                                    value={path}
                                    onChange={(e) => setPath(e.target.value)}
                                />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="ctype" className="text-right">
                                    Content-Type
                                </Label>
                                <Input
                                    id="ctype"
                                    placeholder="application/json (Optional)"
                                    className="col-span-3 text-sm"
                                    value={contentType}
                                    onChange={(e) => setContentType(e.target.value)}
                                />
                            </div>
                        </div>
                        <DialogFooter>
                            <Button variant="outline" onClick={() => setIsOpen(false)}>Cancel</Button>
                            <Button onClick={handleAdd}>Save Rule</Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Active Mappings</CardTitle>
                    <CardDescription>
                        Requests matching these patterns will be intercepted locally.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    {rules.length === 0 ? (
                        <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
                            <FileText className="h-12 w-12 mb-4 opacity-20" />
                            <p>No map local rules defined.</p>
                        </div>
                    ) : (
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>URL Pattern (Regex)</TableHead>
                                    <TableHead>Local File Path</TableHead>
                                    <TableHead>Content-Type</TableHead>
                                    <TableHead className="w-[100px]">Actions</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {rules.map((rule) => (
                                    <TableRow key={rule.id}>
                                        <TableCell className="font-mono text-xs">{rule.url_regex}</TableCell>
                                        <TableCell className="font-mono text-xs truncate max-w-[300px]" title={rule.local_path}>
                                            {rule.local_path}
                                        </TableCell>
                                        <TableCell className="text-xs">{rule.content_type || 'Auto'}</TableCell>
                                        <TableCell>
                                            <Button variant="ghost" size="icon" onClick={() => handleDelete(rule.id)}>
                                                <Trash2 className="h-4 w-4 text-destructive" />
                                            </Button>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
