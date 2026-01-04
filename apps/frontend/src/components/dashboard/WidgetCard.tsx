import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { ReactNode } from "react";

interface WidgetCardProps {
    title: string;
    description?: string;
    children: ReactNode;
    className?: string; // For column spanning e.g. "col-span-2"
    contentClassName?: string;
}

export function WidgetCard({ title, description, children, className, contentClassName }: WidgetCardProps) {
    return (
        <Card className={cn("flex flex-col h-full shadow-sm hover:shadow-md transition-shadow", className)}>
            <CardHeader className="pb-2">
                <CardTitle className="text-lg font-medium tracking-tight text-foreground/90">{title}</CardTitle>
                {description && <CardDescription>{description}</CardDescription>}
            </CardHeader>
            <CardContent className={cn("flex-1", contentClassName)}>
                {children}
            </CardContent>
        </Card>
    );
}
