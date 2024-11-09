import { build as esbuild } from "esbuild";
import { watch } from "chokidar";
import { sync } from "glob";

/**
 * Build options for esbuild.
 * @type {import("esbuild").BuildOptions}
 */
const buildOptions = {
  bundle: true,
  outdir: "build/js",
  minify: process.env.NODE_ENV === "production",
  sourcemap: process.env.NODE_ENV !== "production",
};

/**
 * Build JavaScript files using esbuild.
 * @returns {Promise<void>}
 */
async function BuildJs() {
  buildOptions.entryPoints = sync("assets/js/*.js");

  await esbuild(buildOptions).catch(() => process.exit(1));

  if (process.argv.includes("--watch")) {
    // Watch mode
    watch("assets/js/**/*.js", {
      interval: 0,
    }).on("all", async (event, path) => {
      console.log(`JS file ${event}: ${path}`);
      try {
        await esbuild(buildOptions);
        console.log("JS rebuilt successfully");
      } catch (error) {
        console.error("Error rebuilding JS:", error);
      }
    });
  }
}

BuildJs();
