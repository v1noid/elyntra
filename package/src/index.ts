"#!/usr/bin/env node";
import os from "os";
import fs from "fs";
import { Command } from "commander";
import ora from "ora";
import { spawnSync } from "child_process";

const tmp = os.tmpdir() + "/" + crypto.randomUUID();
const CURR_DIR = process.cwd();
const TEMPLATE_DIR = "template";
const prg = new Command();

const downloadCommands = [
  "git init",
  "git remote add origin https://github.com/v1noid/elyntra.git",
  "git config core.sparseCheckout true",
  `echo '${TEMPLATE_DIR}' >> .git/info/sparse-checkout`,
  "git pull origin main",
];

async function run(commands: string[], cwd?: string) {
  for (const command of commands) {
    const [cmd, ...args] = command.split(" ");
    const a = spawnSync(cmd, args, {
      cwd: cwd || tmp,
    });
  }
}

const mora = (txt: string) =>
  ora({
    spinner: "bouncingBall",
    text: txt,
  });

prg
  .command("init <name>")
  .description("Initialize the server folder in just a click")
  .action(async (name: string) => {
    try {
      let d = mora("Downloading template...").start();
      fs.mkdirSync(tmp);

      await run(downloadCommands).catch((e) => {
        d.fail("Failed to download template!");
        console.error(e);
      });
      d.succeed("Downloaded template successfully!");
      d = ora("Copying files...").start();
      fs.rmSync(tmp + "/.git", { recursive: true, force: true });

      fs.cpSync(tmp + "/" + TEMPLATE_DIR, CURR_DIR + "/" + name, {
        recursive: true,
      });
      fs.rmSync(tmp, { recursive: true, force: true });
      d.succeed(`Server (${name}) folder initialized successfully!`);
      d = ora("Installing dependencies...").start();
      await run(["bun i"], CURR_DIR + "/" + name).catch((e) => {
        d.fail("Failed to install dependencies!");
        console.error(e);
      });
      d.succeed("Dependencies installed successfully!");
    } catch (error) {
      console.error("\nFailed to set up elyntra: ", error);
    }
  });

prg.parse(process.argv);
