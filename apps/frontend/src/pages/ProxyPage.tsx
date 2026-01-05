
import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Copy, Check, Terminal } from 'lucide-react';
import { useConfigStore } from '@/stores/useConfigStore';
import { toast } from 'sonner';
import { generateSnippet } from '@/utils/snippets';
import { Code2 } from 'lucide-react';


export default function ProxyPage() {
    const { config } = useConfigStore();
    const [copied, setCopied] = useState(false);

    const port = config?.proxy_port || 8888;
    const proxyUrl = `http://127.0.0.1:${port}`;

    const commands = {
        unix: `export HTTP_PROXY=${proxyUrl} HTTPS_PROXY=${proxyUrl} ALL_PROXY=${proxyUrl}`,
        powershell: `$env:HTTP_PROXY="${proxyUrl}"; $env:HTTPS_PROXY="${proxyUrl}"; $env:ALL_PROXY="${proxyUrl}"`,
        cmd: `set HTTP_PROXY=${proxyUrl} && set HTTPS_PROXY=${proxyUrl} && set ALL_PROXY=${proxyUrl}`,
    };

    const handleCopy = (cmd: string) => {
        navigator.clipboard.writeText(cmd);
        setCopied(true);
        toast.success("Command copied to clipboard");
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <div className="p-8 space-y-6">
            <div>
                <h1 className="text-3xl font-bold">API Proxy</h1>
                <p className="text-muted-foreground">Configure local proxy server and providers.</p>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Terminal className="h-5 w-5" />
                        Terminal Configuration
                    </CardTitle>
                    <CardDescription>
                        Route your terminal traffic through Siberia. Run this command in your shell.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <Tabs defaultValue="unix" className="w-full">
                        <TabsList className="grid w-full grid-cols-3">
                            <TabsTrigger value="unix">Bash / Zsh (Mac/Linux)</TabsTrigger>
                            <TabsTrigger value="powershell">PowerShell (Windows)</TabsTrigger>
                            <TabsTrigger value="cmd">CMD (Windows)</TabsTrigger>
                        </TabsList>

                        {Object.entries(commands).map(([key, cmd]) => (
                            <TabsContent key={key} value={key} className="space-y-4 mt-4">
                                <div className="space-y-2">
                                    <Label>Run this command:</Label>
                                    <div className="flex gap-2">
                                        <Input
                                            readOnly
                                            value={cmd}
                                            className="font-mono text-sm bg-muted"
                                        />
                                        <Button
                                            size="icon"
                                            variant="outline"
                                            onClick={() => handleCopy(cmd)}
                                        >
                                            {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                                        </Button>
                                    </div>
                                    <p className="text-xs text-muted-foreground">
                                        This will route HTTP and HTTPS traffic from CLI tools (curl, git, npm, etc.) through Siberia for this session.
                                    </p>
                                </div>
                            </TabsContent>
                        ))}
                    </Tabs>
                </CardContent>

            </Card>

            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Code2 className="h-5 w-5" />
                        Connect to Siberia
                    </CardTitle>
                    <CardDescription>
                        Use Siberia as a drop-in replacement for OpenAI/Claude in your code.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <Tabs defaultValue="python" className="w-full">
                        <TabsList className="grid w-full grid-cols-4">
                            <TabsTrigger value="python">Python</TabsTrigger>
                            <TabsTrigger value="node">Node.js</TabsTrigger>
                            <TabsTrigger value="curl">Curl</TabsTrigger>
                            <TabsTrigger value="env">.env</TabsTrigger>
                        </TabsList>

                        {['python', 'node', 'curl', 'env'].map((lang) => {
                            const snippet = generateSnippet(lang as any, port);
                            return (
                                <TabsContent key={lang} value={lang} className="space-y-4 mt-4">
                                    <div className="relative">
                                        <pre className="p-4 rounded-lg bg-muted overflow-x-auto font-mono text-sm border">
                                            {snippet}
                                        </pre>
                                        <Button
                                            size="icon"
                                            variant="ghost"
                                            className="absolute top-2 right-2"
                                            onClick={() => handleCopy(snippet)}
                                        >
                                            {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                                        </Button>
                                    </div>
                                </TabsContent>
                            );
                        })}
                    </Tabs>
                </CardContent>
            </Card>
        </div >
    );
}
