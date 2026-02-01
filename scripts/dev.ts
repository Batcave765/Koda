import { spawn } from "bun";

console.log("🚀 Starting Koda services...");

const backend = spawn(["docker-compose", "up", "--build"], {
    cwd: "./backend",
    stdout: "inherit",
    stderr: "inherit",
});

const frontend = spawn(["bun", "run", "dev"], {
    cwd: "./frontend",
    stdout: "inherit",
    stderr: "inherit",
});

// Handle termination (e.g., Ctrl+C)
process.on("SIGINT", () => {
    console.log("\n🛑 Shutting down services...");
    backend.kill();
    frontend.kill();
    process.exit();
});

// Wait for processes to exit (though they are long-running)
await Promise.all([backend.exited, frontend.exited]);
