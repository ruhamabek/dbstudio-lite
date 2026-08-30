"use client"

import React, { useRef, useState, useCallback, useEffect } from "react"
import Editor, { OnMount, loader } from "@monaco-editor/react"
import { useDatabase } from "@/context/DatabaseContext"
import { Button } from "@/components/ui/button"
import { DataViewer } from "@/components/viewer/DataViewer"
import { Play, RotateCcw, AlertCircle, Loader2, GripHorizontal } from "lucide-react"

export function QueryConsole() {
  const { query, runQuery, setQuerySql } = useDatabase()
  const editorRef = useRef<any>(null)
  const [monacoReady, setMonacoReady] = useState(false)

   const [editorHeight, setEditorHeight] = useState<number>(200)
  const [isDragging, setIsDragging] = useState<boolean>(false)
  const dragStartY = useRef<number>(0)
  const dragStartHeight = useRef<number>(200)

   useEffect(() => {
    import("monaco-editor").then((monaco) => {
      loader.config({ monaco })
      setMonacoReady(true)
    })
  }, [])

  const handleEditorMount: OnMount = (editor, monaco) => {
    editorRef.current = editor

    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, () => {
      runQuery()
    })
  }

  const handleClear = () => {
    setQuerySql("")
    if (editorRef.current) {
      editorRef.current.setValue("")
      editorRef.current.focus()
    }
  }

   const handleMouseDown = useCallback(
    (e: React.MouseEvent) => {
      e.preventDefault()
      setIsDragging(true)
      dragStartY.current = e.clientY
      dragStartHeight.current = editorHeight
      document.body.style.userSelect = "none"
      document.body.style.cursor = "row-resize"
    },
    [editorHeight]
  )

  useEffect(() => {
    if (!isDragging) return

    const handleMouseMove = (e: MouseEvent) => {
      const deltaY = e.clientY - dragStartY.current
      const maxAllowed = Math.max(200, window.innerHeight - 250)
      const newHeight = Math.max(80, Math.min(maxAllowed, dragStartHeight.current + deltaY))
      setEditorHeight(newHeight)
    }

    const handleMouseUp = () => {
      setIsDragging(false)
      document.body.style.userSelect = ""
      document.body.style.cursor = ""
    }

    window.addEventListener("mousemove", handleMouseMove)
    window.addEventListener("mouseup", handleMouseUp)

    return () => {
      window.removeEventListener("mousemove", handleMouseMove)
      window.removeEventListener("mouseup", handleMouseUp)
    }
  }, [isDragging])

  return (
    <div className="flex flex-col h-full overflow-hidden">
       <div className="flex flex-col bg-card shrink-0">
         <div className="flex items-center justify-between border-b border-border px-3 py-1.5 bg-muted/40 shrink-0">
          <span className="text-xs font-semibold text-muted-foreground">
            SQL Query Editor
          </span>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleClear}
              className="h-7 px-2 text-xs gap-1"
            >
              <RotateCcw className="h-3 w-3" />
              Clear
            </Button>
            <Button
              size="sm"
              onClick={() => runQuery()}
              disabled={query.isRunning || !query.sql.trim()}
              className="h-7 px-3 text-xs gap-1.5"
            >
              {query.isRunning ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : (
                <Play className="h-3 w-3 fill-current" />
              )}
              Run Query
              <kbd className="ml-1 text-[10px] opacity-60 font-mono">
                ⌘↵
              </kbd>
            </Button>
          </div>
        </div>

         <div
          style={{ height: `${editorHeight}px` }}
          className="w-full bg-card overflow-hidden transition-none"
        >
          {monacoReady ? (
            <Editor
              height="100%"
              defaultLanguage="sql"
              language="sql"
              theme="vs-dark"
              value={query.sql}
              onChange={(val) => setQuerySql(val || "")}
              onMount={handleEditorMount}
              loading={
                <div className="flex h-full items-center justify-center text-xs text-muted-foreground">
                  <Loader2 className="h-4 w-4 animate-spin mr-2" />
                  Loading editor...
                </div>
              }
              options={{
                minimap: { enabled: false },
                scrollBeyondLastLine: false,
                fontSize: 14.5,
                fontFamily: "var(--font-mono), monospace",
                automaticLayout: true,
                lineNumbers: "on",
                wordWrap: "on",
                padding: { top: 8, bottom: 8 },
                suggestOnTriggerCharacters: true,
                tabSize: 2,
              }}
            />
          ) : (
            <div className="flex h-full items-center justify-center text-xs text-muted-foreground">
              <Loader2 className="h-4 w-4 animate-spin mr-2" />
              Initializing local editor...
            </div>
          )}
        </div>
      </div>

       <div
        onMouseDown={handleMouseDown}
        className={`relative flex h-2 w-full cursor-row-resize items-center justify-center border-y border-border bg-muted/40 hover:bg-primary/20 transition-colors select-none group z-10 shrink-0 ${
          isDragging ? "bg-primary/30 border-primary" : ""
        }`}
        title="Drag up or down to resize query editor"
      >
        <GripHorizontal className="h-3 w-3 text-muted-foreground/50 group-hover:text-foreground transition-colors" />
      </div>

       {query.error && (
        <div className="flex items-start space-x-2 border-b border-destructive/20 bg-destructive/10 p-3 text-xs text-destructive shrink-0">
          <AlertCircle className="h-4 w-4 shrink-0 mt-0.5" />
          <div className="flex-1 font-mono text-[11px] break-all">
            {query.error}
          </div>
        </div>
      )}

       <div className="flex-1 overflow-hidden min-h-[100px]">
        {query.isRunning ? (
          <div className="flex h-full items-center justify-center space-x-2 text-muted-foreground">
            <Loader2 className="h-5 w-5 animate-spin" />
            <span className="text-xs">Executing query...</span>
          </div>
        ) : (
          <DataViewer result={query.result} />
        )}
      </div>
    </div>
  )
}
