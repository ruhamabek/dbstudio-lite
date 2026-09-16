"use client"

import React, { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useDatabase } from "@/context/DatabaseContext"
import { db } from "@wailsjs/go/models"
import {
  AlertCircle,
  Database,
  Link,
  Sliders,
  Bookmark,
  Trash2,
  FileCode,
  Layers,
} from "lucide-react"

interface ConnectionModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

type DriverType = "postgres" | "mysql" | "sqlite"

export function ConnectionModal({ open, onOpenChange }: ConnectionModalProps) {
  const {
    connectWithUrl,
    connectWithConfig,
    connection,
    savedConnections,
    saveCurrentConnection,
    deleteSavedConnection,
  } = useDatabase()

  const [mode, setMode] = useState<"saved" | "url" | "fields">(
    savedConnections.length > 0 ? "saved" : "url"
  )

  const [selectedDriver, setSelectedDriver] = useState<DriverType>("postgres")
  const [sqlitePath, setSqlitePath] = useState("/home/sapphire/dbstudio-lite/sample.db")

  // URL mode
  const [rawUrl, setRawUrl] = useState(
    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
  )
  const [profileName, setProfileName] = useState("")
  const [shouldSave, setShouldSave] = useState(false)

   const [host, setHost] = useState("localhost")
  const [port, setPort] = useState("5432")
  const [database, setDatabase] = useState("postgres")
  const [user, setUser] = useState("postgres")
  const [password, setPassword] = useState("")
  const [sslMode, setSslMode] = useState("disable")

  const handleDriverChange = (driver: DriverType) => {
    setSelectedDriver(driver)
    if (driver === "sqlite") {
      setRawUrl("sqlite:///home/sapphire/dbstudio-lite/sample.db")
    } else if (driver === "mysql") {
      if (port === "5432" || !port) setPort("3306")
      if (user === "postgres") setUser("root")
      if (database === "postgres") setDatabase("mydb")
      setRawUrl("mysql://root:password@localhost:3306/mydb")
    } else {
      if (port === "3306" || !port) setPort("5432")
      if (user === "root") setUser("postgres")
      if (database === "mydb") setDatabase("postgres")
      setRawUrl("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
    }
  }

  const handleConnect = async (e: React.FormEvent) => {
    e.preventDefault()
    let success = false
    if (mode === "url") {
      success = await connectWithUrl(rawUrl)
      if (success && shouldSave && profileName.trim()) {
        await saveCurrentConnection(profileName.trim(), rawUrl)
      }
    } else {
      let cfg: db.Config
      if (selectedDriver === "sqlite") {
        cfg = new db.Config({
          driver: "sqlite",
          filePath: sqlitePath.trim(),
        })
      } else {
        cfg = new db.Config({
          driver: selectedDriver,
          host: host.trim(),
          port: parseInt(port, 10) || (selectedDriver === "mysql" ? 3306 : 5432),
          database: database.trim(),
          user: user.trim(),
          password: password,
          sslMode: selectedDriver === "postgres" ? sslMode : "",
          filePath: "",
        })
      }
      success = await connectWithConfig(cfg)
      if (success && shouldSave && profileName.trim()) {
        const urlToSave =
          selectedDriver === "sqlite"
            ? `sqlite://${sqlitePath.trim()}`
            : selectedDriver === "mysql"
            ? `mysql://${user}:${password}@${host}:${port}/${database}`
            : `postgres://${user}:${password}@${host}:${port}/${database}?sslmode=${sslMode}`
        await saveCurrentConnection(profileName.trim(), urlToSave)
      }
    }
    if (success) {
      onOpenChange(false)
    }
  }

  const handleConnectSaved = async (saved: db.SavedConnection) => {
    const success = await connectWithUrl(saved.url)
    if (success) {
      onOpenChange(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent onClose={() => onOpenChange(false)}>
        <form onSubmit={handleConnect} className="space-y-4.5">
          <DialogHeader>
            <div className="flex items-center space-x-2.5">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <Database className="h-4.5 w-4.5" />
              </div>
              <DialogTitle className="text-lg font-bold">Connect to Database</DialogTitle>
            </div>
            <DialogDescription className="text-sm">
              Select a saved profile or connect to PostgreSQL, MySQL, or SQLite.
            </DialogDescription>
          </DialogHeader>

          {/* Mode Switcher */}
          <div className="flex rounded-lg bg-muted p-1 text-muted-foreground text-sm font-medium">
            <button
              type="button"
              onClick={() => setMode("saved")}
              className={`flex-1 flex items-center justify-center gap-2 py-2 rounded-md transition-colors ${
                mode === "saved"
                  ? "bg-background text-foreground shadow-xs font-semibold"
                  : "hover:text-foreground"
              }`}
            >
              <Bookmark className="h-4 w-4" />
              Saved ({savedConnections.length})
            </button>
            <button
              type="button"
              onClick={() => setMode("url")}
              className={`flex-1 flex items-center justify-center gap-2 py-2 rounded-md transition-colors ${
                mode === "url"
                  ? "bg-background text-foreground shadow-xs font-semibold"
                  : "hover:text-foreground"
              }`}
            >
              <Link className="h-4 w-4" />
              URL / File
            </button>
            <button
              type="button"
              onClick={() => setMode("fields")}
              className={`flex-1 flex items-center justify-center gap-2 py-2 rounded-md transition-colors ${
                mode === "fields"
                  ? "bg-background text-foreground shadow-xs font-semibold"
                  : "hover:text-foreground"
              }`}
            >
              <Sliders className="h-4 w-4" />
              Parameters
            </button>
          </div>

           {connection.error && (
            <div className="flex items-start space-x-2.5 rounded-lg border border-destructive/20 bg-destructive/10 p-3.5 text-sm text-destructive">
              <AlertCircle className="h-4.5 w-4.5 shrink-0 mt-0.5" />
              <span className="break-all">{connection.error}</span>
            </div>
          )}

           {mode === "saved" && (
            <div className="space-y-2.5 max-h-72 overflow-y-auto pr-1">
              {savedConnections.length > 0 ? (
                savedConnections.map((saved) => (
                  <div
                    key={saved.id}
                    className="flex items-center justify-between p-3 rounded-lg border border-border bg-background hover:border-primary/50 transition-colors"
                  >
                    <div className="flex flex-col min-w-0 pr-3">
                      <span className="text-sm font-semibold text-foreground truncate">
                        {saved.name}
                      </span>
                      <span className="text-xs font-mono text-muted-foreground truncate">
                        {saved.url.replace(/:[^:@]*@/, ":****@")}
                      </span>
                    </div>
                    <div className="flex items-center space-x-2 shrink-0">
                      <Button
                        type="button"
                        size="sm"
                        onClick={() => handleConnectSaved(saved)}
                        disabled={connection.isConnecting}
                        className="h-8 text-xs font-semibold px-3"
                      >
                        Connect
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => deleteSavedConnection(saved.id)}
                        className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                ))
              ) : (
                <div className="py-10 text-center text-sm text-muted-foreground space-y-1.5">
                  <Bookmark className="h-7 w-7 mx-auto opacity-30" />
                  <p className="font-semibold text-foreground">No saved connections</p>
                  <p className="text-xs max-w-xs mx-auto">
                    Switch to the "URL / File" tab and check "Save connection" to store profiles here.
                  </p>
                </div>
              )}
            </div>
          )}

           {mode === "url" && (
            <div className="space-y-3.5">
              <div className="space-y-1.5">
                <label className="text-sm font-medium text-foreground">
                  Connection URI or File Path
                </label>
                <Input
                  value={rawUrl}
                  onChange={(e) => setRawUrl(e.target.value)}
                  placeholder="postgres://..., mysql://..., sqlite:///path/to/db, or /path/to/app.db"
                  className="font-mono text-xs h-10"
                  required
                />
                <p className="text-[11px] text-muted-foreground">
                  Supports PostgreSQL (<code className="text-primary font-mono">postgres://</code>), MySQL (<code className="text-primary font-mono">mysql://</code>), or SQLite (<code className="text-primary font-mono">sqlite:///</code> or direct <code className="text-primary font-mono">.db</code> path).
                </p>
              </div>

               <div className="rounded-lg border border-border bg-muted/20 p-3.5 space-y-2.5">
                <label className="flex items-center space-x-2.5 cursor-pointer text-sm font-medium">
                  <input
                    type="checkbox"
                    checked={shouldSave}
                    onChange={(e) => setShouldSave(e.target.checked)}
                    className="rounded border-input text-primary focus:ring-primary h-4 w-4"
                  />
                  <span className="text-foreground">
                    Save this connection profile
                  </span>
                </label>
                {shouldSave && (
                  <div className="pt-1">
                    <Input
                      value={profileName}
                      onChange={(e) => setProfileName(e.target.value)}
                      placeholder="e.g. Local Postgres, Staging MySQL, SQLite App DB"
                      className="text-sm h-9"
                      required={shouldSave}
                    />
                  </div>
                )}
              </div>
            </div>
          )}

           {mode === "fields" && (
            <div className="space-y-3.5">
               <div className="flex rounded-lg border border-border p-1 bg-muted/30 text-xs font-semibold">
                <button
                  type="button"
                  onClick={() => handleDriverChange("postgres")}
                  className={`flex-1 py-1.5 rounded-md transition-colors flex items-center justify-center gap-1.5 ${
                    selectedDriver === "postgres"
                      ? "bg-primary text-primary-foreground shadow-xs"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Database className="h-3.5 w-3.5" />
                  PostgreSQL
                </button>
                <button
                  type="button"
                  onClick={() => handleDriverChange("mysql")}
                  className={`flex-1 py-1.5 rounded-md transition-colors flex items-center justify-center gap-1.5 ${
                    selectedDriver === "mysql"
                      ? "bg-primary text-primary-foreground shadow-xs"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Layers className="h-3.5 w-3.5" />
                  MySQL
                </button>
                <button
                  type="button"
                  onClick={() => handleDriverChange("sqlite")}
                  className={`flex-1 py-1.5 rounded-md transition-colors flex items-center justify-center gap-1.5 ${
                    selectedDriver === "sqlite"
                      ? "bg-primary text-primary-foreground shadow-xs"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <FileCode className="h-3.5 w-3.5" />
                  SQLite
                </button>
              </div>

              {/* SQLite Specific Form */}
              {selectedDriver === "sqlite" ? (
                <div className="space-y-1.5">
                  <label className="text-sm font-medium text-foreground">
                    SQLite Database File Path
                  </label>
                  <Input
                    value={sqlitePath}
                    onChange={(e) => setSqlitePath(e.target.value)}
                    placeholder="/path/to/database.db or /home/user/app.sqlite"
                    className="font-mono text-xs h-10"
                    required
                  />
                  <p className="text-[11px] text-muted-foreground">
                    Enter the absolute path to your local SQLite (<code className="text-primary font-mono">.db</code>, <code className="text-primary font-mono">.sqlite</code>, <code className="text-primary font-mono">.sqlite3</code>) database file.
                  </p>
                </div>
              ) : (
                 <div className="grid grid-cols-2 gap-3.5 text-sm">
                  <div className="col-span-2 sm:col-span-1 space-y-1.5">
                    <label className="font-medium text-muted-foreground">Host</label>
                    <Input
                      value={host}
                      onChange={(e) => setHost(e.target.value)}
                      placeholder="localhost"
                      className="h-10 text-sm"
                      required
                    />
                  </div>
                  <div className="col-span-2 sm:col-span-1 space-y-1.5">
                    <label className="font-medium text-muted-foreground">Port</label>
                    <Input
                      type="text"
                      inputMode="numeric"
                      pattern="[0-9]*"
                      value={port}
                      onChange={(e) => setPort(e.target.value.replace(/\D/g, ""))}
                      placeholder={selectedDriver === "mysql" ? "3306" : "5432"}
                      className="h-10 text-sm"
                      required
                    />
                  </div>
                  <div className="col-span-2 sm:col-span-1 space-y-1.5">
                    <label className="font-medium text-muted-foreground">Database</label>
                    <Input
                      value={database}
                      onChange={(e) => setDatabase(e.target.value)}
                      placeholder={selectedDriver === "mysql" ? "mydb" : "postgres"}
                      className="h-10 text-sm"
                      required
                    />
                  </div>
                  <div className="col-span-2 sm:col-span-1 space-y-1.5">
                    <label className="font-medium text-muted-foreground">Username</label>
                    <Input
                      value={user}
                      onChange={(e) => setUser(e.target.value)}
                      placeholder={selectedDriver === "mysql" ? "root" : "postgres"}
                      className="h-10 text-sm"
                      required
                    />
                  </div>
                  <div className="col-span-2 space-y-1.5">
                    <label className="font-medium text-muted-foreground">Password</label>
                    <Input
                      type="password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      placeholder="••••••••"
                      className="h-10 text-sm"
                    />
                  </div>
                  {selectedDriver === "postgres" && (
                    <div className="col-span-2 space-y-1.5">
                      <label className="font-medium text-muted-foreground">SSL Mode</label>
                      <select
                        value={sslMode}
                        onChange={(e) => setSslMode(e.target.value)}
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-xs transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring"
                      >
                        <option value="disable">disable</option>
                        <option value="require">require</option>
                        <option value="prefer">prefer</option>
                        <option value="verify-full">verify-full</option>
                      </select>
                    </div>
                  )}
                </div>
              )}

               <div className="rounded-lg border border-border bg-muted/20 p-3.5 space-y-2.5">
                <label className="flex items-center space-x-2.5 cursor-pointer text-sm font-medium">
                  <input
                    type="checkbox"
                    checked={shouldSave}
                    onChange={(e) => setShouldSave(e.target.checked)}
                    className="rounded border-input text-primary focus:ring-primary h-4 w-4"
                  />
                  <span className="text-foreground">
                    Save this connection profile
                  </span>
                </label>
                {shouldSave && (
                  <div className="pt-1">
                    <Input
                      value={profileName}
                      onChange={(e) => setProfileName(e.target.value)}
                      placeholder="e.g. Local SQLite DB, Dev Postgres"
                      className="text-sm h-9"
                      required={shouldSave}
                    />
                  </div>
                )}
              </div>
            </div>
          )}

          <DialogFooter className="pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              className="h-10 px-4 text-sm font-medium"
            >
              Cancel
            </Button>
            {mode !== "saved" && (
              <Button
                type="submit"
                disabled={connection.isConnecting}
                className="h-10 px-5 text-sm font-semibold shadow-xs"
              >
                {connection.isConnecting ? "Connecting..." : "Connect"}
              </Button>
            )}
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
