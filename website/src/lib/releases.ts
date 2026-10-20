export const REPOSITORY_URL = "https://github.com/ruhamabek/dbstudio-lite"
export const RELEASE_URL = `${REPOSITORY_URL}/releases/latest`


export const PLATFORMS = [
  { id: "macos", name: "macOS", arch: "Apple Silicon & Intel", format: "Universal · ZIP", url: `${REPOSITORY_URL}/releases/download/v0.1.0/DBStudio-Lite-macOS-Universal.2.zip` },
  { id: "windows", name: "Windows", arch: "64-bit", format: "x64 · ZIP", url: `${REPOSITORY_URL}/releases/download/v0.1.0/DBStudio-Lite-Windows-x64.zip` },
  { id: "linux", name: "Linux", arch: "64-bit", format: "x86_64 · ZIP", url: `${REPOSITORY_URL}/releases/download/v0.1.0/DBStudio-Lite-Linux-x86_64.1.zip` },
] as const
