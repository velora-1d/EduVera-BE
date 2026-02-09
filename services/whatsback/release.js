#!/usr/bin/env node
const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const pkgPath = path.join(__dirname, "package.json");
const pkg = JSON.parse(fs.readFileSync(pkgPath, "utf8"));
const currentVersion = pkg.version;

const bumpType = process.argv[2] || "patch";

/**
 * Increment a version string by one unit, according to the specified bump type.
 *
 * @param {string} version
 * @param {"major"|"minor"|"patch"} type
 * @return {string} new version string
 */
function bumpVersion(version, type) {
    const [major, minor, patch] = version
        .split(".")
        .map((num) => parseInt(num, 10));
    if (type === "major") {
        return `${major + 1}.0.0`;
    } else if (type === "minor") {
        return `${major}.${minor + 1}.0`;
    } else {
    // default: patch
        return `${major}.${minor}.${patch + 1}`;
    }
}

const newVersion = bumpVersion(currentVersion, bumpType);

pkg.version = newVersion;
fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));
console.log(`🚀 Updated package.json version to ${newVersion}`);

// Stage the updated package.json
console.log("🔨 Staging package.json...");
execSync("git add package.json", { stdio: "inherit" });

// Commit the version bump
console.log("💾 Committing changes...");
execSync(`git commit -m "chore: bump version to ${newVersion}"`, {
    stdio: "inherit",
});

// Create a tag (without a "v" prefix)
console.log("🏷️  Creating git tag...");
execSync(`git tag v${newVersion}`, { stdio: "inherit" });

// Push commit and tags to the remote
console.log("🚚 Pushing commits to remote...");
execSync("git push origin main", { stdio: "inherit" });
console.log("🚚 Pushing tags to remote...");
execSync("git push --tags", { stdio: "inherit" });

console.log(`🎉 Released version ${newVersion}`);
