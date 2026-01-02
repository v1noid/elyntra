import { defineConfig } from "tsup";

export default defineConfig({
  format: ["cjs"],
  entry: ["./src/index.ts"],
  dts: true,
  clean: true,
  external: ["fs", "os", "ora"],
  outDir: "./dist",
  treeshake: true,
  skipNodeModulesBundle: true,
  minify: true,
  minifySyntax: true,
  minifyWhitespace: true,
  banner: {
    js: `#!/usr/bin/env node\n`,
  },
});
