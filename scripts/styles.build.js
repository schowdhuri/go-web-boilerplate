import postcss from "postcss";
import tailwindcss from "tailwindcss";
import autoprefixer from "autoprefixer";
import { watch } from "chokidar";
import { sync } from "glob";
import fs from "node:fs";

/**
 * Process CSS files using PostCSS with Tailwind CSS and Autoprefixer.
 * @param {string} inputDir - The input directory containing CSS files.
 * @param {string} outputDir - The output directory to write processed CSS files.
 * @returns {Promise<void>}
 */
function build(inputDir, outputDir, isDev = false) {
  const files = sync(`${inputDir}/*.css`);
  // Create output directory if it doesn't exist
  fs.mkdirSync(outputDir, { recursive: true });

  for (const file of files) {
    const outputFile = file.replace(inputDir, outputDir);
    const css = fs.readFileSync(file, "utf8");
    postcss([tailwindcss, autoprefixer])
      .process(css, { from: undefined })
      .then((result) => {
        if (!isDev) {
          // TODO: Minify CSS in production
          result.css = result.css.replace(/\s+/g, " ");
        }
        fs.writeFileSync(outputFile, result.css);
      });
  }
}

async function main() {
  const watchMode = process.argv.includes("--watch");
  const inputDir = "assets/css";
  const outputDir = "build/css";

  await build(inputDir, outputDir, !!watchMode);

  if (watchMode) {
    watch(inputDir, { ignoreInitial: true }).on("all", async (event, path) => {
      if (event === "change" || event === "add") {
        console.log(`File ${path} was ${event}, rebuilding...`);
        await build(inputDir, outputDir, true);
      }
    });
  }
}

main();
