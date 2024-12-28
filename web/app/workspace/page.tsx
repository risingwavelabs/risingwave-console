'use client'

import { DatabaseList } from "@/components/ui/database-list"
import { useState, useCallback, useEffect, useRef } from 'react'
import { Button } from "@/components/ui/button"
import { Settings } from 'lucide-react'
import { DatabaseManagement } from "@/components/ui/database-management"
import { SQLEditor } from "@/components/ui/sql-editor"

// Sample data - replace with real data from your backend
const sampleDatabases = [
  {
    id: "db1",
    name: "Main Database",
    tables: [
      { id: "t1", name: "users" },
      { id: "t2", name: "products" },
      { id: "t3", name: "orders" },
    ],
  },
  {
    id: "db2",
    name: "Analytics DB",
    tables: [
      { id: "t4", name: "events" },
      { id: "t5", name: "metrics" },
    ],
  },
]

// Sample saved queries - replace with real data from your backend
const savedQueries = [
  { id: "q1", name: "Get Active Users" },
  { id: "q2", name: "Monthly Revenue" },
  { id: "q3", name: "Product Inventory" },
]

const MIN_WIDTH = 150 // 9.375rem
const MAX_WIDTH = 400 // 25rem

export default function WorkspacePage() {
  const [isResizing, setIsResizing] = useState(false)
  const [panelWidth, setPanelWidth] = useState(192) // 12rem default
  const [editorWidth, setEditorWidth] = useState(0)
  const [showDatabaseManagement, setShowDatabaseManagement] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)
  const panelRef = useRef<HTMLDivElement>(null)

  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    setIsResizing(true)
    document.body.style.cursor = 'ew-resize'
  }, [])

  const handleMouseUp = useCallback(() => {
    setIsResizing(false)
    document.body.style.cursor = 'default'
  }, [])

  const handleMouseMove = useCallback((e: React.MouseEvent<HTMLDivElement> | MouseEvent) => {
    if (!isResizing || !containerRef.current) return

    const containerRect = containerRef.current.getBoundingClientRect()
    const newWidth = e.clientX - containerRect.left

    if (newWidth >= MIN_WIDTH && newWidth <= MAX_WIDTH) {
      setPanelWidth(newWidth)
    }
  }, [isResizing])

  // Add and remove event listeners
  const handleMouseMoveGlobal = useCallback((e: MouseEvent) => {
    handleMouseMove(e)
  }, [handleMouseMove])

  useEffect(() => {
    if (isResizing) {
      window.addEventListener('mousemove', handleMouseMoveGlobal)
      window.addEventListener('mouseup', handleMouseUp)
    }

    return () => {
      window.removeEventListener('mousemove', handleMouseMoveGlobal)
      window.removeEventListener('mouseup', handleMouseUp)
    }
  }, [isResizing, handleMouseMoveGlobal, handleMouseUp])

  const calculateEditorWidth = useCallback(() => {
    if (containerRef.current) {
      const totalWidth = containerRef.current.getBoundingClientRect().width
      setEditorWidth(totalWidth - panelWidth - 5)
    }
  }, [panelWidth])

  useEffect(() => {
    calculateEditorWidth()
    window.addEventListener('resize', calculateEditorWidth)

    return () => {
      window.removeEventListener('resize', calculateEditorWidth)
    }
  }, [calculateEditorWidth])

  const handleTableSelect = (dbId: string, tableId: string) => {
    console.log(`Selected table ${tableId} from database ${dbId}`)
  }

  const handleRunQuery = (query: string) => {
    console.log('Running query:', query)
    // Implement query execution
  }

  const handleSaveQuery = (query: string, name: string) => {
    console.log('Saving query:', name, query)
    // Implement query saving
  }

  return (
    <div 
      ref={containerRef} 
      className="flex h-full w-full relative"
      onMouseMove={handleMouseMove}
    >
      <div
        ref={panelRef}
        style={{ width: panelWidth }}
        className="relative flex-shrink-0 border-r bg-background h-full"
      >
        <div className="p-4 flex items-center justify-between w-full border-b pb-2">
          <h2 className="text-sm">Databases</h2>
          <Button
            variant="ghost"
            size="sm"
            className="h-7 px-2 text-muted-foreground hover:text-foreground"
            onClick={() => setShowDatabaseManagement(true)}
          >
            <Settings className="w-4 h-4 mr-1" />
          </Button>
        </div>
        <div className="p-4">
          <DatabaseList
            databases={sampleDatabases}
            onSelectTable={handleTableSelect}
          />
        </div>
        <div
          className="absolute right-0 top-0 bottom-0 w-[2px] cursor-ew-resize hover:bg-blue-500/20 active:bg-blue-500/40 transition-colors"
          onMouseDown={handleMouseDown}
        />
      </div>

      <SQLEditor
        width={editorWidth}
        savedQueries={savedQueries}
        onRunQuery={handleRunQuery}
        onSaveQuery={handleSaveQuery}
      />

      <DatabaseManagement
        isOpen={showDatabaseManagement}
        onClose={() => setShowDatabaseManagement(false)}
      />
    </div>
  )
} 
