import { Component, ErrorInfo, ReactNode } from "react";
import { AlertTriangle } from "lucide-react";
import { Button } from "@/components/ui/button";

interface Props {
    children?: ReactNode;
}

interface State {
    hasError: boolean;
    error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false,
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error("Uncaught error:", error, errorInfo);
    }

    public render() {
        if (this.state.hasError) {
            return (
                <div className="flex flex-col items-center justify-center min-h-[50vh] p-8 space-y-4 text-center">
                    <div className="p-4 bg-red-100 dark:bg-red-900/20 rounded-full">
                        <AlertTriangle className="w-8 h-8 text-red-600 dark:text-red-400" />
                    </div>
                    <h1 className="text-xl font-bold">Something went wrong</h1>
                    <p className="text-muted-foreground max-w-md">
                        The application encountered an unexpected error. This might happen if you are running in Web Mode without the Backend.
                    </p>
                    <div className="text-xs text-left bg-muted p-2 rounded overflow-auto max-w-sm max-h-32 font-mono">
                        {this.state.error?.message}
                    </div>
                    <Button onClick={() => window.location.reload()}>Reload Page</Button>
                </div>
            );
        }

        return this.props.children;
    }
}
