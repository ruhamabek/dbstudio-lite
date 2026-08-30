"use client"

import React, { useState } from "react"
import { useDatabase } from "@/context/DatabaseContext"
import { ConnectionBar } from "@/components/connection/ConnectionBar"
import { ConnectionModal } from "@/components/connection/ConnectionModal"
import { TableSidebar } from "@/components/sidebar/TableSidebar"
import { DataViewer } from "@/components/viewer/DataViewer"
import { StructureViewer } from "@/components/viewer/StructureViewer"
import { QueryConsole } from "@/components/console/QueryConsole"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { Button } from "@/components/ui/button"
import { Database, Table, Columns, Terminal, Plus } from "lucide-react"

function Dashboard() {
  const { activeTab, setActiveTab, schema, query, connection } = useDatabase()
  const [modalOpen, setModalOpen] = useState(false)

  return (
    <div className="flex h-screen w-screen flex-col bg-background text-foreground overflow-hidden">
      {/* Top Connection Bar */}
      <ConnectionBar />

      {/* Main Workspace */}
      <div className="flex flex-1 overflow-hidden">
        {/* Left Table Sidebar */}
        <TableSidebar />

        {/* Right Content Area */}
        <main className="flex flex-1 flex-col overflow-hidden bg-background">
          {connection.isConnected ? (
            <Tabs
              value={activeTab}
              onValueChange={(val) => setActiveTab(val as any)}
              className="flex h-full flex-col"
            >
              {/* Tabs Navigation Header */}
              <div className="flex items-center justify-between border-b border-border bg-card px-4 py-2">
                <TabsList className="h-9">
                  <TabsTrigger value="data" className="gap-2 text-sm px-3 py-1.5">
                    <Table className="h-4 w-4" />
                    Data
                    {schema.selectedTable && (
                      <span className="font-semibold text-foreground">
                        ({schema.selectedTable.Name})
                      </span>
                    )}
                  </TabsTrigger>
                  <TabsTrigger value="structure" className="gap-2 text-sm px-3 py-1.5">
                    <Columns className="h-4 w-4" />
                    Structure
                  </TabsTrigger>
                  <TabsTrigger value="console" className="gap-2 text-sm px-3 py-1.5">
                    <Terminal className="h-4 w-4" />
                    SQL Console
                  </TabsTrigger>
                </TabsList>
              </div>

              {/* Tab 1: Table Data */}
              <TabsContent value="data" className="flex-1 overflow-hidden mt-0">
                <DataViewer
                  result={query.result}
                  primaryKeys={schema.primaryKeys}
                />
              </TabsContent>

              {/* Tab 2: Structure / Columns */}
              <TabsContent value="structure" className="flex-1 overflow-hidden mt-0">
                <StructureViewer />
              </TabsContent>

              {/* Tab 3: SQL Console */}
              <TabsContent value="console" className="flex-1 overflow-hidden mt-0">
                <QueryConsole />
              </TabsContent>
            </Tabs>
          ) : (
            /* Welcome / Disconnected Screen */
            <div className="flex h-full flex-col items-center justify-center p-8 text-center">
              <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-muted/60 text-muted-foreground mb-5 shadow-xs">
                <Database className="h-8 w-8 text-primary" />
              </div>
              <h2 className="text-2xl font-bold tracking-tight mb-2">
                Welcome to DBStudio Lite
              </h2>
              <p className="text-sm text-muted-foreground max-w-md leading-relaxed mb-6">
                Connect to a PostgreSQL database to inspect schemas, explore tables,
                and execute SQL queries.
              </p>
              <Button
                size="lg"
                onClick={() => setModalOpen(true)}
                className="gap-2 text-sm font-semibold h-11 px-6 shadow-sm"
              >
                <Plus className="h-4 w-4" />
                New Connection
              </Button>
            </div>
          )}
        </main>
      </div>

      <ConnectionModal open={modalOpen} onOpenChange={setModalOpen} />
    </div>
  )
}

export default function Home() {
  return <Dashboard />
}
