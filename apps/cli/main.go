package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Default Siberia Proxy Port
const ProxyPort = 19999

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// 1. Prepare Environment
	proxyURL := fmt.Sprintf("http://localhost:%d", ProxyPort)
	env := os.Environ()
	env = append(env, fmt.Sprintf("HTTP_PROXY=%s", proxyURL))
	env = append(env, fmt.Sprintf("HTTPS_PROXY=%s", proxyURL))
	env = append(env, fmt.Sprintf("ALL_PROXY=%s", proxyURL))
	env = append(env, "SIBERIA_ACTIVE=1") // Marker

	// 2. Identify Command
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// 3. LookPath (resolve absolute path of executable)
	binary, err := exec.LookPath(cmdName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: command not found: %s\n", cmdName)
		os.Exit(127)
	}

	// 4. Exec (Replace Process)
	// We use syscall.Exec to replace the current process with the new one,
	// maintaining PID and stdin/stdout/stderr pipes automatically.
	// Argv[0] should usually be the binary name.
	argv := append([]string{cmdName}, cmdArgs...)

	err = syscall.Exec(binary, argv, env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: sib <command> [args...]")
	fmt.Println("Example: sib curl http://example.com")
	fmt.Println("         sib npm install")
	fmt.Println("Wraps the command with Siberia Proxy environment variables.")
}
