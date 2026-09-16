"use client"

import React, { createContext, useContext, useState, useCallback, useEffect } from "react"
import {
  Connect,
  ConnectWithURL,
  Disconnect,
  Ping,
  GetTables,
  GetColumns,
  GetPrimaryKeys,
  ExecuteQuery,
  GetSavedConnections,
  SaveConnection,
  DeleteConnection,
} from "@wailsjs/go/main/App"
import { db } from "@wailsjs/go/models"

export type ActiveTab = "data" | "structure" | "console"

interface ConnectionState {
  isConnected: boolean
  isConnecting: boolean
  activeDsn: string | null
  error: string | null
}

interface SchemaState {
  tables: db.Table[]
  selectedTable: db.Table | null
  columns: db.Column[]
  primaryKeys: string[]
  isLoading: boolean
  error: string | null
}

interface QueryState {
  sql: string
  result: db.QueryResult | null
  isRunning: boolean
  error: string | null
}

interface DatabaseContextValue {
  connection: ConnectionState
  schema: SchemaState
  query: QueryState
  activeTab: ActiveTab
  savedConnections: db.SavedConnection[]
  setActiveTab: (tab: ActiveTab) => void
  connectWithUrl: (url: string) => Promise<boolean>
  connectWithConfig: (cfg: db.Config) => Promise<boolean>
  disconnect: () => Promise<void>
  selectTable: (table: db.Table) => Promise<void>
  runQuery: (customSql?: string) => Promise<void>
  setQuerySql: (sql: string) => void
  refreshTables: () => Promise<void>
  saveCurrentConnection: (name: string, url: string) => Promise<void>
  deleteSavedConnection: (id: string) => Promise<void>
}

const DatabaseContext = createContext<DatabaseContextValue | undefined>(undefined)

const DEFAULT_SQL = "SELECT * FROM information_schema.tables LIMIT 20;"

export function DatabaseProvider({ children }: { children: React.ReactNode }) {
  const [connection, setConnection] = useState<ConnectionState>({
    isConnected: false,
    isConnecting: false,
    activeDsn: null,
    error: null,
  })

  const [schema, setSchema] = useState<SchemaState>({
    tables: [],
    selectedTable: null,
    columns: [],
    primaryKeys: [],
    isLoading: false,
    error: null,
  })

  const [query, setQuery] = useState<QueryState>({
    sql: DEFAULT_SQL,
    result: null,
    isRunning: false,
    error: null,
  })

  const [savedConnections, setSavedConnections] = useState<db.SavedConnection[]>([])
  const [activeTab, setActiveTab] = useState<ActiveTab>("data")

   const loadSavedConnections = useCallback(async () => {
    try {
      const list = await GetSavedConnections()
      setSavedConnections(list || [])
    } catch (err) {
      console.error("Failed to load saved connections:", err)
    }
  }, [])

  useEffect(() => {
    loadSavedConnections()
  }, [loadSavedConnections])

   const loadTables = useCallback(async () => {
    try {
      setSchema((prev) => ({ ...prev, isLoading: true, error: null }))
      const tables = await GetTables()
      setSchema((prev) => ({
        ...prev,
        tables: tables || [],
        isLoading: false,
      }))
    } catch (err: any) {
      setSchema((prev) => ({
        ...prev,
        isLoading: false,
        error: err?.message || String(err),
      }))
    }
  }, [])

   const connectWithUrl = useCallback(
    async (url: string): Promise<boolean> => {
      setConnection((prev) => ({ ...prev, isConnecting: true, error: null }))
      try {
        await ConnectWithURL(url)
        await Ping()
        setConnection({
          isConnected: true,
          isConnecting: false,
          activeDsn: url.replace(/:[^:@]*@/, ":****@"), // mask password
          error: null,
        })
        await loadTables()
        return true
      } catch (err: any) {
        setConnection({
          isConnected: false,
          isConnecting: false,
          activeDsn: null,
          error: err?.message || String(err),
        })
        return false
      }
    },
    [loadTables]
  )

   const connectWithConfig = useCallback(
    async (cfg: db.Config): Promise<boolean> => {
      setConnection((prev) => ({ ...prev, isConnecting: true, error: null }))
      try {
        await Connect(cfg)
        await Ping()
        const activeDsn =
          cfg.driver === "sqlite" || cfg.filePath
            ? cfg.filePath
            : `${cfg.user}@${cfg.host}:${cfg.port || (cfg.driver === "mysql" ? 3306 : 5432)}/${cfg.database}`

        setConnection({
          isConnected: true,
          isConnecting: false,
          activeDsn: activeDsn,
          error: null,
        })
        await loadTables()
        return true
      } catch (err: any) {
        setConnection({
          isConnected: false,
          isConnecting: false,
          activeDsn: null,
          error: err?.message || String(err),
        })
        return false
      }
    },
    [loadTables]
  )

   const disconnect = useCallback(async () => {
    try {
      await Disconnect()
    } catch {
      // ignore
    } finally {
      setConnection({
        isConnected: false,
        isConnecting: false,
        activeDsn: null,
        error: null,
      })
      setSchema({
        tables: [],
        selectedTable: null,
        columns: [],
        primaryKeys: [],
        isLoading: false,
        error: null,
      })
      setQuery({
        sql: DEFAULT_SQL,
        result: null,
        isRunning: false,
        error: null,
      })
    }
  }, [])

   const selectTable = useCallback(async (table: db.Table) => {
    setSchema((prev) => ({
      ...prev,
      selectedTable: table,
      isLoading: true,
      error: null,
    }))

    try {
      const [columns, primaryKeys] = await Promise.all([
        GetColumns(table.Schema, table.Name),
        GetPrimaryKeys(table.Schema, table.Name),
      ])

      setSchema((prev) => ({
        ...prev,
        columns: columns || [],
        primaryKeys: primaryKeys || [],
        isLoading: false,
      }))

       const tableSql = `SELECT * FROM "${table.Schema}"."${table.Name}" LIMIT 50;`
      setQuery((prev) => ({ ...prev, sql: tableSql, isRunning: true, error: null }))

      const result = await ExecuteQuery(tableSql)
      setQuery((prev) => ({
        ...prev,
        result,
        isRunning: false,
        error: null,
      }))
    } catch (err: any) {
      setSchema((prev) => ({
        ...prev,
        isLoading: false,
        error: err?.message || String(err),
      }))
      setQuery((prev) => ({
        ...prev,
        isRunning: false,
        error: err?.message || String(err),
      }))
    }
  }, [])

   const runQuery = useCallback(
    async (customSql?: string) => {
      const sqlToRun = customSql !== undefined ? customSql : query.sql
      if (!sqlToRun.trim()) return

      setQuery((prev) => ({ ...prev, isRunning: true, error: null }))
      try {
        const result = await ExecuteQuery(sqlToRun)
        setQuery((prev) => ({
          ...prev,
          result,
          isRunning: false,
          error: null,
        }))
      } catch (err: any) {
        setQuery((prev) => ({
          ...prev,
          isRunning: false,
          error: err?.message || String(err),
        }))
      }
    },
    [query.sql]
  )

  const setQuerySql = useCallback((sql: string) => {
    setQuery((prev) => ({ ...prev, sql }))
  }, [])

   const saveCurrentConnection = useCallback(
    async (name: string, url: string) => {
      const saved = await SaveConnection(
        new db.SavedConnection({
          name,
          url,
        })
      )
      setSavedConnections((prev) => {
        const filtered = prev.filter((p) => p.id !== saved.id)
        return [...filtered, saved]
      })
    },
    []
  )

   const deleteSavedConnection = useCallback(async (id: string) => {
    await DeleteConnection(id)
    setSavedConnections((prev) => prev.filter((p) => p.id !== id))
  }, [])

  const value: DatabaseContextValue = {
    connection,
    schema,
    query,
    activeTab,
    savedConnections,
    setActiveTab,
    connectWithUrl,
    connectWithConfig,
    disconnect,
    selectTable,
    runQuery,
    setQuerySql,
    refreshTables: loadTables,
    saveCurrentConnection,
    deleteSavedConnection,
  }

  return (
    <DatabaseContext.Provider value={value}>{children}</DatabaseContext.Provider>
  )
}

export function useDatabase() {
  const context = useContext(DatabaseContext)
  if (!context) {
    throw new Error("useDatabase must be used within a DatabaseProvider")
  }
  return context
}
