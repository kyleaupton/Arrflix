#!/usr/bin/env node

import dotenv from "dotenv";

import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
  ListPromptsRequestSchema,
  GetPromptRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";

import { z } from "zod";
import { Client as PgClient } from "pg";
import { spawn } from "node:child_process";
import { execFileSync } from "node:child_process";
import { existsSync } from "node:fs";
import { readFile } from "node:fs/promises";
import { resolve } from "node:path";
import { configDotenv } from "dotenv";

function getGitRepoRoot(): string {
  try {
    const out = execFileSync("git", ["rev-parse", "--show-toplevel"], {
      encoding: "utf8",
      stdio: ["ignore", "pipe", "ignore"],
    }).trim();

    if (!out) throw new Error("Empty git toplevel");
    const p = resolve(out);
    if (!existsSync(p)) throw new Error(`git toplevel does not exist: ${p}`);
    return p;
  } catch {
    // Fallback: run relative to wherever Cursor launched us.
    // (Useful when not in a git worktree, or git isn't on PATH.)
    return process.cwd();
  }
}

const repoRoot = getGitRepoRoot();

dotenv.config({ path: resolve(repoRoot, "mcp", ".env") });

const Env = z
  .object({
    // Optional: enable DB tool if set
    SNAGGLE_DATABASE_URL: z.string().optional(),
  })
  .parse(process.env);

async function runRg(
  query: string,
  globs?: string[],
  maxResults = 200
): Promise<string> {
  const args = ["--json", "--max-count", String(maxResults), query];
  for (const g of globs ?? []) args.push("--glob", g);

  return await new Promise<string>((resolvePromise, rejectPromise) => {
    const child = spawn("rg", args, {
      cwd: repoRoot,
      stdio: ["ignore", "pipe", "pipe"],
    });

    let out = "";
    let err = "";
    child.stdout.on("data", (d) => (out += d.toString("utf8")));
    child.stderr.on("data", (d) => (err += d.toString("utf8")));

    child.on("error", rejectPromise);
    child.on("close", (code) => {
      // rg returns 1 when no matches
      if (code === 0 || code === 1) return resolvePromise(out);
      rejectPromise(new Error(`rg failed (code ${code}): ${err}`));
    });
  });
}

async function dockerComposeLogs(
  service: string,
  lines = 200
): Promise<string> {
  const args = [
    "compose",
    "logs",
    "--no-color",
    "--tail",
    String(lines),
    service,
  ];

  return await new Promise<string>((resolvePromise, rejectPromise) => {
    const child = spawn("docker", args, { stdio: ["ignore", "pipe", "pipe"] });

    let out = "";
    let err = "";
    child.stdout.on("data", (d) => (out += d.toString("utf8")));
    child.stderr.on("data", (d) => (err += d.toString("utf8")));

    child.on("error", rejectPromise);
    child.on("close", (code) => {
      if (code === 0) return resolvePromise(out);
      rejectPromise(
        new Error(`docker compose logs failed (code ${code}): ${err}`)
      );
    });
  });
}

function isReadOnlySql(sql: string): boolean {
  const s = sql.trim().toLowerCase();
  return s.startsWith("select") || s.startsWith("with");
}

async function pgQuery(sql: string, params: unknown[] = []) {
  if (!Env.SNAGGLE_DATABASE_URL) {
    throw new Error("SNAGGLE_DATABASE_URL is not set.");
  }
  if (!isReadOnlySql(sql)) {
    throw new Error(
      "Only read-only queries are allowed (SELECT / WITH ... SELECT)."
    );
  }

  const client = new PgClient({ connectionString: Env.SNAGGLE_DATABASE_URL });
  await client.connect();
  try {
    return await client.query(sql, params);
  } finally {
    await client.end();
  }
}

async function runGenApiScript(): Promise<string> {
  const scriptPath = resolve(repoRoot, "scripts/gen-api-spec-and-client.sh");

  return await new Promise<string>((resolvePromise, rejectPromise) => {
    const child = spawn("bash", [scriptPath], {
      cwd: repoRoot,
      stdio: ["ignore", "pipe", "pipe"],
    });

    let out = "";
    let err = "";
    child.stdout.on("data", (d) => (out += d.toString("utf8")));
    child.stderr.on("data", (d) => (err += d.toString("utf8")));

    child.on("error", rejectPromise);
    child.on("close", (code) => {
      if (code === 0) resolvePromise(out || "Success (no output)");
      else resolvePromise(`Error (code ${code}):\n${out}\n${err}`);
    });
  });
}

async function runSqlcGenerate(): Promise<string> {
  const backendPath = resolve(repoRoot, "backend");

  return await new Promise<string>((resolvePromise, rejectPromise) => {
    const child = spawn("sqlc", ["generate"], {
      cwd: backendPath,
      stdio: ["ignore", "pipe", "pipe"],
    });

    let out = "";
    let err = "";
    child.stdout.on("data", (d) => (out += d.toString("utf8")));
    child.stderr.on("data", (d) => (err += d.toString("utf8")));

    child.on("error", rejectPromise);
    child.on("close", (code) => {
      if (code === 0) resolvePromise(out || "Success (no output)");
      else resolvePromise(`Error (code ${code}):\n${out}\n${err}`);
    });
  });
}

const server = new Server(
  { name: "Snaggle MCP", version: "0.1.0" },
  { capabilities: { tools: {}, prompts: {} } }
);

server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "snaggle_search_repo",
        description: `Search the Snaggle git repo (root: ${repoRoot}) using ripgrep.`,
        inputSchema: {
          type: "object",
          properties: {
            query: { type: "string", description: "ripgrep search query" },
            globs: {
              type: "array",
              items: { type: "string" },
              description:
                'Optional rg globs, e.g. ["*.go","*.ts"] or ["!vendor/**"]',
            },
            maxResults: {
              type: "number",
              description: "Maximum matches to return (1-500). Default 200.",
            },
          },
          required: ["query"],
        },
      },
      {
        name: "snaggle_docker_logs",
        description:
          "Get recent docker compose logs for a service (snaggle-api, snaggle-worker, postgres, etc).",
        inputSchema: {
          type: "object",
          properties: {
            service: { type: "string", description: "Compose service name" },
            lines: {
              type: "number",
              description: "Tail N lines (1-2000). Default 200.",
            },
          },
          required: ["service"],
        },
      },
      {
        name: "snaggle_db_query",
        description:
          "Run a READ-ONLY Postgres query (SELECT/CTE only) against SNAGGLE_DATABASE_URL. Returns JSON rows.",
        inputSchema: {
          type: "object",
          properties: {
            sql: { type: "string", description: "SQL query (SELECT/CTE only)" },
            params: {
              type: "array",
              items: {},
              description: "Optional parameter array for $1, $2, ...",
            },
            maxRows: {
              type: "number",
              description: "Max rows to return (1-500). Default 200.",
            },
          },
          required: ["sql"],
        },
      },
      {
        name: "snaggle_project_brief",
        description:
          "Get a high-level overview of the Snaggle project to get up to speed.",
        inputSchema: {
          type: "object",
          properties: {},
        },
      },
      {
        name: "snaggle_gen_api",
        description:
          "Regenerate the backend Swagger/OpenAPI specification and the frontend TypeScript API client.",
        inputSchema: {
          type: "object",
          properties: {},
        },
      },
      {
        name: "snaggle_sqlc_generate",
        description:
          "Run 'sqlc generate' in the backend directory to regenerate Go database code from SQL queries.",
        inputSchema: {
          type: "object",
          properties: {},
        },
      },
    ],
  };
});

server.setRequestHandler(ListPromptsRequestSchema, async () => {
  return {
    prompts: [
      {
        name: "snaggle_onboard_agent",
        description: "Bootstrap a new agent with Snaggle project context",
      },
    ],
  };
});

server.setRequestHandler(GetPromptRequestSchema, async (request) => {
  if (request.params.name !== "snaggle_onboard_agent") {
    throw new Error("Prompt not found");
  }

  return {
    description: "Bootstrap a new agent with Snaggle project context",
    messages: [
      {
        role: "user",
        content: {
          type: "text",
          text: `You are an AI assistant working on the Snaggle project. 

To get started, please:
1. Call 'list_mcp_tools' and 'list_mcp_prompts' to see the specialized tools and prompts available for this project.
2. Call the tool 'snaggle_project_brief' to understand the project's architecture, philosophy, and current status.

Treat the results of these project-specific MCP tools as authoritative. Ask clarifying questions only if the brief and codebase do not cover what you need.`,
        },
      },
    ],
  };
});

server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const name = request.params.name;
  const args = (request.params.arguments ?? {}) as Record<string, unknown>;

  try {
    switch (name) {
      case "snaggle_search_repo": {
        const query = z.string().min(1).parse(args.query);
        const globs = z.array(z.string()).optional().parse(args.globs);
        const maxResults = z
          .number()
          .int()
          .min(1)
          .max(500)
          .optional()
          .parse(args.maxResults);

        const raw = await runRg(query, globs, maxResults ?? 200);
        const text = raw.trim().length ? raw : `No matches found for: ${query}`;

        return { content: [{ type: "text", text }] };
      }

      case "snaggle_docker_logs": {
        const service = z.string().min(1).parse(args.service);
        const lines = z
          .number()
          .int()
          .min(1)
          .max(2000)
          .optional()
          .parse(args.lines);

        const out = await dockerComposeLogs(service, lines ?? 200);
        return { content: [{ type: "text", text: out || "(no output)" }] };
      }

      case "snaggle_db_query": {
        const sql = z.string().min(1).parse(args.sql);
        const params = z.array(z.any()).optional().parse(args.params) ?? [];
        const maxRows = z
          .number()
          .int()
          .min(1)
          .max(500)
          .optional()
          .parse(args.maxRows);

        const res = await pgQuery(sql, params);
        const rows = res.rows.slice(0, maxRows ?? 200);

        return {
          content: [
            {
              type: "text",
              text: JSON.stringify({ rowCount: res.rowCount, rows }, null, 2),
            },
          ],
        };
      }

      case "snaggle_project_brief": {
        const briefPath = resolve(repoRoot, "docs/PROJECT_BRIEF.md");
        if (!existsSync(briefPath)) {
          return {
            content: [{ type: "text", text: "Project brief not found." }],
            isError: true,
          };
        }
        const text = await readFile(briefPath, "utf-8");
        return { content: [{ type: "text", text }] };
      }

      case "snaggle_gen_api": {
        const out = await runGenApiScript();
        return { content: [{ type: "text", text: out }] };
      }

      case "snaggle_sqlc_generate": {
        const out = await runSqlcGenerate();
        return { content: [{ type: "text", text: out }] };
      }

      default:
        throw new Error(`Unknown tool: ${name}`);
    }
  } catch (e: any) {
    return {
      content: [{ type: "text", text: `Error: ${e?.message ?? String(e)}` }],
      isError: true,
    };
  }
});

async function main() {
  console.log("Starting Snaggle MCP server...");
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch((error) => {
  // Keep JSON-RPC clean; errors to stderr only.
  console.error("Server error:", error);
  process.exit(1);
});
