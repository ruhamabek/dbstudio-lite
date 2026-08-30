"use client"

import React, { useState, useMemo } from "react"
import { useDatabase } from "@/context/DatabaseContext"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Table, Search, Layers, Loader2, Database } from "lucide-react"

export function TableSidebar() {
  const { schema, selectTable, connection } = useDatabase()
  const [search, setSearch] = useState("")

   const filteredTables = useMemo(() => {
    if (!search.trim()) return schema.tables
    const query = search.toLowerCase()
    return schema.tables.filter(
      (t) =>
        t.Name.toLowerCase().includes(query) ||
        t.Schema.toLowerCase().includes(query)
    )
  }, [schema.tables, search])

   const groupedTables = useMemo(() => {
    const map = new Map<string, typeof schema.tables>()
    for (const t of filteredTables) {
      const list = map.get(t.Schema) || []
      list.push(t)
      map.set(t.Schema, list)
    }
    return map
  }, [filteredTables])

  if (!connection.isConnected) {
    return (
      <aside className="flex w-68 flex-col border-r border-border bg-card p-6 text-center shrink-0">
        <div className="my-auto flex flex-col items-center space-y-3 text-muted-foreground">
          <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-muted/60 mb-1">
            <Database className="h-7 w-7 opacity-40 text-primary" />
          </div>
          <p className="text-sm font-semibold text-foreground">No connection active</p>
          <p className="text-xs text-muted-foreground leading-relaxed">
            Connect to a database to inspect tables, view schemas, and query data.
          </p>
        </div>
      </aside>
    )
  }

  return (
    <aside className="flex w-68 flex-col border-r border-border bg-card shrink-0">
       <div className="p-3.5 border-b border-border space-y-2.5">
        <div className="flex items-center justify-between">
          <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
            <Layers className="h-4 w-4" />
            Tables ({schema.tables.length})
          </span>
          {schema.isLoading && (
            <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
          )}
        </div>
        <div className="relative">
          <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Filter tables..."
            className="h-9 pl-9 text-sm bg-background"
          />
        </div>
      </div>

       <div className="flex-1 overflow-y-auto p-2 space-y-4">
        {Array.from(groupedTables.entries()).map(([schemaName, tables]) => (
          <div key={schemaName} className="space-y-1">
            <div className="px-2.5 py-0.5 text-xs font-bold text-muted-foreground uppercase tracking-wider">
              {schemaName}
            </div>
            <div className="space-y-0.5">
              {tables.map((t) => {
                const isSelected =
                  schema.selectedTable?.Schema === t.Schema &&
                  schema.selectedTable?.Name === t.Name

                return (
                  <button
                    key={`${t.Schema}.${t.Name}`}
                    onClick={() => selectTable(t)}
                    type="button"
                    className={`flex w-full items-center justify-between rounded-md px-3 py-2 text-left text-sm transition-colors ${
                      isSelected
                        ? "bg-primary text-primary-foreground font-semibold"
                        : "text-muted-foreground hover:bg-muted hover:text-foreground font-medium"
                    }`}
                  >
                    <span className="flex items-center gap-2.5 truncate">
                      <Table className="h-4 w-4 shrink-0 opacity-70" />
                      <span className="truncate">{t.Name}</span>
                    </span>
                    {t.Schema !== "public" && (
                      <Badge
                        variant="muted"
                        className="text-[10px] px-1.5 py-0.5 h-4.5 uppercase"
                      >
                        {t.Schema}
                      </Badge>
                    )}
                  </button>
                )
              })}
            </div>
          </div>
        ))}

        {groupedTables.size === 0 && (
          <div className="py-8 text-center text-sm text-muted-foreground">
            No tables found.
          </div>
        )}
      </div>
    </aside>
  )
}
