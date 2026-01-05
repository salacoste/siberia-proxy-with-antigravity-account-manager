import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Regex } from "lucide-react";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { FilterOptions } from "@/lib/filterUtils";

interface FilterBarProps {
    options: FilterOptions;
    onChange: (opts: FilterOptions) => void;
}

export function FilterBar({ options, onChange }: FilterBarProps) {
    const handleChange = (key: keyof FilterOptions, value: any) => {
        onChange({ ...options, [key]: value });
    };

    return (
        <div className="flex flex-wrap items-center gap-2">
            <div className="flex-1 min-w-[200px] flex gap-2">
                <div className="relative flex-1">
                    <Input
                        placeholder={options.isRegex ? "Regex Pattern (e.g. ^/api/v1/.*)" : "Filter (e.g. method:POST /api/ status:404 !css)"}
                        value={options.query}
                        onChange={(e) => handleChange("query", e.target.value)}
                        className={`font-mono text-xs pr-8 ${options.isRegex ? "border-orange-500/50" : ""}`}
                    />
                </div>
                <Button
                    variant={options.isRegex ? "default" : "outline"}
                    size="icon"
                    className="h-10 w-10 shrink-0"
                    onClick={() => handleChange("isRegex", !options.isRegex)}
                    title="Toggle Regex Mode"
                >
                    <Regex className="h-4 w-4" />
                </Button>
            </div>

            <Select
                value={options.statusCategory}
                onValueChange={(val) => handleChange("statusCategory", val)}
            >
                <SelectTrigger className="w-[130px] h-10 text-xs font-mono">
                    <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                    <SelectItem value="All">All Status</SelectItem>
                    <SelectItem value="2xx">Success (2xx)</SelectItem>
                    <SelectItem value="3xx">Redirect (3xx)</SelectItem>
                    <SelectItem value="4xx">Client Err (4xx)</SelectItem>
                    <SelectItem value="5xx">Server Err (5xx)</SelectItem>
                    <SelectItem value="Error">Errors (4xx/5xx)</SelectItem>
                </SelectContent>
            </Select>

            <Select
                value={options.methods.length > 0 ? options.methods[0] : "All"}
                onValueChange={(val) => handleChange("methods", val === "All" ? [] : [val])}
            >
                <SelectTrigger className="w-[100px] h-10 text-xs font-mono">
                    <SelectValue placeholder="Method" />
                </SelectTrigger>
                <SelectContent>
                    <SelectItem value="All">All Methods</SelectItem>
                    <SelectItem value="GET">GET</SelectItem>
                    <SelectItem value="POST">POST</SelectItem>
                    <SelectItem value="PUT">PUT</SelectItem>
                    <SelectItem value="DELETE">DELETE</SelectItem>
                    <SelectItem value="PATCH">PATCH</SelectItem>
                    <SelectItem value="OPTIONS">OPTIONS</SelectItem>
                </SelectContent>
            </Select>
        </div>
    );
}
