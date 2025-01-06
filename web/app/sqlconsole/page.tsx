'use client'

import { DatabaseList, type DatabaseItem } from "../../components/ui/database-list"
import { useState, useCallback, useEffect, useRef } from 'react'
import { Button } from "../../components/ui/button"
import { Settings } from 'lucide-react'
import { DatabaseManagement } from "../../components/ui/database-management"
import { SQLEditor } from "../../components/ui/sql-editor"
import { RisingWaveNodeData } from "../../components/streaming-graph"

// Rest of the file content remains the same
const sampleDatabases: DatabaseItem[] = [
  {
    id: "db1",
    name: "Main Database",
    clusterId: "cluster-1",
    clusterName: "Production Cluster",
    user: "admin",
    password: "********",
    database: "main",
    tables: [
      {
        id: "t1",
        name: "users",
        type: "table",
        columns: [
          { id: "c1", name: "id", type: "int" },
          { id: "c2", name: "name", type: "varchar" },
          { id: "c3", name: "email", type: "varchar" },
          { id: "c4", name: "created_at", type: "timestamp" }
        ]
      },
      {
        id: "t2",
        name: "products",
        type: "table",
        columns: [
          { id: "c5", name: "id", type: "int" },
          { id: "c6", name: "name", type: "varchar" },
          { id: "c7", name: "price", type: "decimal" },
          { id: "c8", name: "stock", type: "int" }
        ]
      },
      {
        id: "t3",
        name: "orders",
        type: "table",
        columns: [
          { id: "c9", name: "id", type: "int" },
          { id: "c10", name: "user_id", type: "int" },
          { id: "c11", name: "product_id", type: "int" },
          { id: "c12", name: "quantity", type: "int" },
          { id: "c13", name: "total", type: "decimal" },
          { id: "c14", name: "order_date", type: "timestamp" }
        ]
      },
      {
        id: "kafka_source",
        name: "kafka_source",
        type: "source",
        columns: [
          { id: "ks1", name: "event_id", type: "varchar" },
          { id: "ks2", name: "user_id", type: "int" },
          { id: "ks3", name: "event_type", type: "varchar" },
          { id: "ks4", name: "timestamp", type: "timestamp" },
          { id: "ks5", name: "data", type: "jsonb" }
        ]
      },
      {
        id: "user_events_mv",
        name: "user_events_mv",
        type: "materialized_view",
        columns: [
          { id: "mv1", name: "user_id", type: "int" },
          { id: "mv2", name: "event_count", type: "int" },
          { id: "mv3", name: "last_event", type: "timestamp" },
          { id: "mv4", name: "event_types", type: "varchar[]" }
        ]
      }
    ],
  },
  {
    id: "db2",
    name: "Analytics DB",
    clusterId: "cluster-2",
    clusterName: "Analytics Cluster",
    user: "analyst",
    password: "********",
    database: "analytics",
    tables: [
      {
        id: "t4",
        name: "events",
        type: "table",
        columns: [
          { id: "c15", name: "event_id", type: "varchar" },
          { id: "c16", name: "user_id", type: "int" },
          { id: "c17", name: "event_type", type: "varchar" },
          { id: "c18", name: "timestamp", type: "timestamp" },
          { id: "c19", name: "data", type: "jsonb" }
        ]
      },
      {
        id: "t5",
        name: "metrics",
        type: "table",
        columns: [
          { id: "c20", name: "id", type: "int" },
          { id: "c21", name: "name", type: "varchar" },
          { id: "c22", name: "value", type: "float" },
          { id: "c23", name: "recorded_at", type: "timestamp" }
        ]
      },
      {
        id: "clickhouse_sink",
        name: "clickhouse_sink",
        type: "sink",
        columns: [
          { id: "cs1", name: "user_id", type: "int" },
          { id: "cs2", name: "event_count", type: "int" },
          { id: "cs3", name: "last_event", type: "timestamp" },
          { id: "cs4", name: "event_types", type: "array(string)" }
        ]
      }
    ],
  }
]

// Add sample database schema for the streaming graph
const sampleDatabaseSchema: RisingWaveNodeData[] = [
  {
    id: 1,
    name: "kafka_source",
    type: "source",
    connector: {
      type: "kafka",
      properties: {
        topic: "user_events",
        bootstrap_servers: "localhost:9092"
      }
    },
    format: "json",
    columns: [
      { name: "event_id", type: "varchar", isPrimary: true },
      { name: "user_id", type: "int" },
      { name: "event_type", type: "varchar" },
      { name: "timestamp", type: "timestamp" },
      { name: "data", type: "jsonb" }
    ]
  },
  {
    id: 2,
    name: "mysql_source",
    type: "source",
    connector: {
      type: "mysql",
      properties: {
        host: "localhost",
        port: "3306",
        database: "users"
      }
    },
    columns: [
      { name: "id", type: "int", isPrimary: true },
      { name: "name", type: "varchar" },
      { name: "email", type: "varchar" },
      { name: "created_at", type: "timestamp" }
    ]
  },
  {
    id: 3,
    name: "user_events_mv",
    type: "materialized_view",
    dependencies: [1, 2],
    columns: [
      { name: "user_id", type: "int", isPrimary: true },
      { name: "event_count", type: "int" },
      { name: "last_event", type: "timestamp" },
      { name: "event_types", type: "varchar[]" }
    ]
  },
  {
    id: 4,
    name: "clickhouse_sink",
    type: "sink",
    dependencies: [3],
    connector: {
      type: "clickhouse",
      properties: {
        host: "localhost",
        port: "8123",
        database: "analytics"
      }
    },
    columns: [
      { name: "user_id", type: "int" },
      { name: "event_count", type: "int" },
      { name: "last_event", type: "timestamp" },
      { name: "event_types", type: "array(string)" }
    ]
  }
]

// Sample saved queries - replace with real data from your backend
const savedQueries = [
  { id: "q1", name: "Get Active Users" },
  { id: "q2", name: "Monthly Revenue" },
  { id: "q3", name: "Product Inventory" },
]

const MIN_WIDTH = 150 // 9.375rem
const MAX_WIDTH = 400 // 25rem

const sampleClusters = [
  { id: "cluster-1", name: "Production Cluster" },
  { id: "cluster-2", name: "Analytics Cluster" }
]

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
        databaseSchema={sampleDatabaseSchema}
      />

      <DatabaseManagement
        isOpen={showDatabaseManagement}
        onClose={() => setShowDatabaseManagement(false)}
        clusters={sampleClusters}
      />
    </div>
  )
}
