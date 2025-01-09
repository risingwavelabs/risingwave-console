'use client'

import { useState, useCallback, useEffect } from 'react'
import { ChevronDown, ChevronRight, Database, Table, ArrowDownToLine, ArrowUpFromLine, Eye } from 'lucide-react'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "./context-menu"

export interface Column {
  name: string
  type: string
}

export interface Relation {
  id: string
  name: string
  columns: Column[]
  type: 'table' | 'source' | 'sink' | 'materialized_view'
}

export interface DatabaseItem {
  id: string
  name: string
  clusterId: string
  clusterName: string
  user: string
  password?: string
  database: string
  tables: Relation[]
}

interface DatabaseListProps {
  databases: DatabaseItem[]
  onSelectTable?: (databaseId: string, tableId: string) => void
  onUseDatabase?: (databaseId: string) => void
}

const SELECTED_DB_KEY = 'selectedDatabaseId'

export function DatabaseList({ databases, onSelectTable, onUseDatabase }: DatabaseListProps) {
  const [expandedDbs, setExpandedDbs] = useState<Set<string>>(new Set())
  const [expandedTables, setExpandedTables] = useState<Set<string>>(new Set())
  const [selectedDbId, setSelectedDbId] = useState<string | null>(null)

  // Initialize selected database from local storage or first available
  useEffect(() => {
    const savedDbId = localStorage.getItem(SELECTED_DB_KEY)

    // If there are databases available
    if (databases.length > 0) {
      // Check if saved ID exists in current database list
      const dbExists = databases.some(db => db.id === savedDbId)

      if (dbExists) {
        setSelectedDbId(savedDbId)
      } else {
        // Use first database if saved one doesn't exist
        setSelectedDbId(databases[0].id)
        localStorage.setItem(SELECTED_DB_KEY, databases[0].id)
      }
    }
  }, [databases])

  const handleUseDatabase = useCallback((dbId: string) => {
    setSelectedDbId(dbId)
    localStorage.setItem(SELECTED_DB_KEY, dbId)
    onUseDatabase?.(dbId)
  }, [onUseDatabase])

  const toggleDb = useCallback((e: React.MouseEvent, dbId: string) => {
    if (e.button === 0) { // Left click only
      const newExpanded = new Set(expandedDbs)
      if (newExpanded.has(dbId)) {
        newExpanded.delete(dbId)
      } else {
        newExpanded.add(dbId)
      }
      setExpandedDbs(newExpanded)
    }
  }, [expandedDbs])

  const toggleTable = useCallback((e: React.MouseEvent, tableId: string) => {
    e.stopPropagation()
    if (e.button === 0) { // Left click only
      const newExpanded = new Set(expandedTables)
      if (newExpanded.has(tableId)) {
        newExpanded.delete(tableId)
      } else {
        newExpanded.add(tableId)
      }
      setExpandedTables(newExpanded)
    }
  }, [expandedTables])

  return (
    <div className="space-y-1">
      {databases.map((db) => (
        <div key={db.id}>
          <ContextMenu>
            <ContextMenuTrigger asChild>
              <button
                onClick={(e) => toggleDb(e, db.id)}
                className={`flex items-center gap-1 w-full hover:bg-muted/50 rounded-sm p-1 text-sm hover:text-foreground ${selectedDbId === db.id ? 'font-semibold text-foreground' : 'text-muted-foreground'
                  }`}
              >
                {expandedDbs.has(db.id) ? (
                  <ChevronDown className="h-4 w-4 shrink-0" />
                ) : (
                  <ChevronRight className="h-4 w-4 shrink-0" />
                )}
                <Database className="h-4 w-4 shrink-0" />
                <span className="truncate">{db.name}</span>
              </button>
            </ContextMenuTrigger>
            <ContextMenuContent className="w-48">
              <ContextMenuItem onClick={() => handleUseDatabase(db.id)}>
                Use Database
              </ContextMenuItem>
            </ContextMenuContent>
          </ContextMenu>
          {expandedDbs.has(db.id) && (
            <div className="ml-6 mt-1 space-y-1">
              <span className="text-xs text-muted-foreground truncate">{db.clusterName}</span>
              {db.tables.map((table) => (
                <div key={table.id} className="space-y-1">
                  <button
                    onClick={(e) => {
                      toggleTable(e, table.id)
                      onSelectTable?.(db.id, table.id)
                    }}
                    className="flex items-center gap-1 w-full hover:bg-muted/50 rounded-sm p-1 text-sm text-muted-foreground hover:text-foreground"
                  >
                    {expandedTables.has(table.id) ? (
                      <ChevronDown className="h-4 w-4 shrink-0" />
                    ) : (
                      <ChevronRight className="h-4 w-4 shrink-0" />
                    )}
                    {table.type === 'table' && <Table className="h-4 w-4 shrink-0" />}
                    {table.type === 'source' && <ArrowDownToLine className="h-4 w-4 shrink-0 rotate-270" />}
                    {table.type === 'sink' && <ArrowUpFromLine className="h-4 w-4 shrink-0 rotate-90" />}
                    {table.type === 'materialized_view' && <Eye className="h-4 w-4 shrink-0" />}
                    <span className="truncate">{table.name}</span>
                  </button>
                  {expandedTables.has(table.id) && (
                    <div className="ml-6 space-y-1">
                      {table.columns.map((column, index) => (
                        <div
                          key={index}
                          className="flex items-center gap-1 w-full rounded-sm p-1 text-sm text-muted-foreground"
                        >
                          <span className="truncate">{column.name}</span>
                          <span className="text-xs text-muted-foreground">({column.type})</span>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      ))}
    </div>
  )
}
