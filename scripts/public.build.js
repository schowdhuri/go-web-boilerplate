import { existsSync, mkdirSync, copyFileSync, statSync } from "fs";
import { dirname } from "path";
import { sync } from "glob";
import { watch } from "chokidar";

/**
 * Copy files from source directory to destination directory.
 * @param {string} srcDir - The source directory.
 * @param {string} destDir - The destination directory.
 * @returns {void}
 */
function copyFiles(srcDir, destDir) {
  sync(`${srcDir}/**/*`)
    .filter((file) => statSync(file).isFile())
    .forEach((file) => {
      const destFile = file.replace(srcDir, destDir);
      const dir = dirname(destFile);

      if (!existsSync(dir)) {
        mkdirSync(dir, { recursive: true });
      }

      copyFileSync(file, destFile);
    });
}

function main() {
  const watchMode = process.argv.includes("--watch");
  const srcDir = "public";
  const destDir = "build";

  copyFiles(srcDir, destDir);

  if (watchMode) {
    watch(srcDir, { ignoreInitial: true }).on("all", async (event, path) => {
      if (event === "change" || event === "add") {
        console.log(`File ${path} was ${event}, copying...`);
        copyFiles(srcDir, destDir);
      }
    });
  }
}

main();
