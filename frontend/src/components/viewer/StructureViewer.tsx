"use client"

import React from "react"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { useDatabase } from "@/context/DatabaseContext"
import { Key, Columns, Loader2 } from "lucide-react"

export function StructureViewer() {
  const { schema } = useDatabase()
  const { selectedTable, columns, primaryKeys, isLoading } = schema

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center p-8 text-muted-foreground">
        <Loader2 className="h-6 w-6 animate-spin mr-2" />
        <span className="text-xs">Loading schema structure...</span>
      </div>
    )
  }

  if (!selectedTable) {
    return (
      <div className="flex h-full flex-col items-center justify-center p-8 text-center text-muted-foreground">
        <Columns className="h-8 w-8 opacity-30 mb-2" />
        <p className="text-sm font-medium">No table selected</p>
        <p className="text-xs">Select a table from the sidebar to inspect its columns.</p>
      </div>
    )
  }

  const pkSet = new Set(primaryKeys)

  return (
    <div className="flex flex-col h-full">
       <div className="flex items-center justify-between border-b border-border bg-muted/30 px-4 py-2 text-xs">
        <div className="flex items-center space-x-2">
          <span className="font-semibold text-foreground">
            {selectedTable.Schema}.{selectedTable.Name}
          </span>
          <span className="text-muted-foreground">
            ({columns.length} columns)
          </span>
        </div>
      </div>

       <div className="flex-1 overflow-auto">
        <Table>
          <TableHeader className="sticky top-0 bg-card z-10">
            <TableRow>
              <TableHead className="w-12 text-center text-[11px] text-muted-foreground">
                #
              </TableHead>
              <TableHead className="text-xs">Column Name</TableHead>
              <TableHead className="text-xs">Data Type</TableHead>
              <TableHead className="text-xs">Nullable</TableHead>
              <TableHead className="text-xs">Default Value</TableHead>
              <TableHead className="text-xs">Key</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {columns.map((col, idx) => {
              const isPk = pkSet.has(col.Name)
              return (
                <TableRow key={col.Name}>
                  <TableCell className="text-center text-[11px] text-muted-foreground/60 font-mono">
                    {idx + 1}
                  </TableCell>
                  <TableCell className="font-mono font-medium text-xs">
                    {col.Name}
                  </TableCell>
                  <TableCell className="font-mono text-xs text-muted-foreground">
                    <Badge variant="muted" className="font-normal font-mono text-[11px]">
                      {col.DataType}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-xs">
                    {col.IsNullable ? (
                      <span className="text-muted-foreground text-xs">NULL</span>
                    ) : (
                      <Badge variant="outline" className="text-[10px] text-muted-foreground">
                        NOT NULL
                      </Badge>
                    )}
                  </TableCell>
                  <TableCell className="font-mono text-xs text-muted-foreground max-w-xs truncate">
                    {col.DefaultVal ? (
                      <code>{col.DefaultVal}</code>
                    ) : (
                      <span className="italic opacity-50">—</span>
                    )}
                  </TableCell>
                  <TableCell className="text-xs">
                    {isPk && (
                      <Badge variant="default" className="gap-1 text-[10px] py-0.5">
                        <Key className="h-2.5 w-2.5" />
                        PRIMARY KEY
                      </Badge>
                    )}
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
