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
import { db } from "@wailsjs/go/models"
import { Clock, Hash } from "lucide-react"

interface DataViewerProps {
  result: db.QueryResult | null
  primaryKeys?: string[]
}

export function DataViewer({ result, primaryKeys = [] }: DataViewerProps) {
  if (!result || !result.columns || result.columns.length === 0) {
    return (
      <div className="flex h-full flex-col items-center justify-center p-8 text-center text-muted-foreground">
        <p className="text-sm font-medium">No results to display</p>
        <p className="text-xs">Run a query or select a table from the sidebar.</p>
      </div>
    )
  }

  const pkSet = new Set(primaryKeys)

  const formatCellValue = (val: any) => {
    if (val === null || val === undefined) {
      return (
        <span className="italic text-muted-foreground/60 text-xs font-mono">
          NULL
        </span>
      )
    }
    if (typeof val === "boolean") {
      return (
        <span
          className={`font-mono text-sm ${
            val ? "text-emerald-500 dark:text-emerald-400 font-medium" : "text-rose-500 dark:text-rose-400"
          }`}
        >
          {val ? "true" : "false"}
        </span>
      )
    }
    if (typeof val === "object") {
      return (
        <span className="font-mono text-sm text-muted-foreground">
          {JSON.stringify(val)}
        </span>
      )
    }
    return <span className="font-mono text-sm">{String(val)}</span>
  }

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between border-b border-border bg-muted/30 px-4 py-2 text-sm text-muted-foreground">
        <div className="flex items-center space-x-4">
          <span className="flex items-center gap-1">
            <Hash className="h-4 w-4" />
            <strong className="text-foreground">{result.rows?.length || 0}</strong> rows
          </span>
          <span className="flex items-center gap-1">
            <Clock className="h-4 w-4" />
            <strong className="text-foreground">{result.durationMs}</strong> ms
          </span>
        </div>
      </div>

       <div className="flex-1 overflow-auto">
        <Table>
          <TableHeader className="sticky top-0 bg-card z-10">
            <TableRow>
              <TableHead className="w-12 text-center text-xs text-muted-foreground">
                #
              </TableHead>
              {result.columns.map((col) => {
                const isPk = pkSet.has(col)
                return (
                  <TableHead key={col} className="whitespace-nowrap font-semibold text-sm">
                    <div className="flex items-center gap-1.5">
                      <span>{col}</span>
                      {isPk && (
                        <Badge
                          variant="secondary"
                          className="h-4.5 px-1.5 text-[10px] gap-0.5"
                        >
                          PK
                        </Badge>
                      )}
                    </div>
                  </TableHead>
                )
              })}
            </TableRow>
          </TableHeader>
          <TableBody>
            {result.rows && result.rows.length > 0 ? (
              result.rows.map((row, rowIdx) => (
                <TableRow key={rowIdx} className="hover:bg-muted/40">
                  <TableCell className="text-center text-xs text-muted-foreground/60 font-mono">
                    {rowIdx + 1}
                  </TableCell>
                  {row.map((cell, cellIdx) => (
                    <TableCell key={cellIdx} className="whitespace-nowrap max-w-sm truncate text-sm">
                      {formatCellValue(cell)}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={result.columns.length + 1}
                  className="h-24 text-center text-xs text-muted-foreground"
                >
                  Empty result set (0 rows returned).
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
