export const REPOSITORY_URL = "https://github.com/ruhamabek/dbstudio-lite"
export const RELEASE_URL = `${REPOSITORY_URL}/releases/tag/v0.1.0`

// ZIP asset names verified against the GitHub v0.1.0 release. Keep destinations
// centralized so hero and platform downloads remain consistent.
export const PLATFORMS = [
  { id: "macos", name: "macOS", arch: "Apple Silicon & Intel", format: "Universal · ZIP", url: `${REPOSITORY_URL}/releases/download/v0.1.0/DBStudio-Lite-macOS-Universal.zip` },
  { id: "windows", name: "Windows", arch: "64-bit", format: "x64 · ZIP", url: `${REPOSITORY_URL}/releases/download/v0.1.0/DBStudio-Lite-Windows-x64.zip` },
  { id: "linux", name: "Linux", arch: "64-bit", format: "x86_64 · ZIP", url: `${REPOSITORY_URL}/releases/download/v0.1.0/DBStudio-Lite-Linux-x86_64.zip` },
] as const
