'use client'

import { useState, useCallback, useEffect } from 'react'
import { ChevronDown, ChevronRight, Database, Table } from 'lucide-react'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "@/components/ui/context-menu"

interface DatabaseItem {
  id: string
  name: string
  tables: { id: string; name: string }[]
}

interface DatabaseListProps {
  databases: DatabaseItem[]
  onSelectTable?: (databaseId: string, tableId: string) => void
  onUseDatabase?: (databaseId: string) => void
}

const SELECTED_DB_KEY = 'selectedDatabaseId'

export function DatabaseList({ databases, onSelectTable, onUseDatabase }: DatabaseListProps) {
  const [expandedDbs, setExpandedDbs] = useState<Set<string>>(new Set())
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

  return (
    <div className="space-y-1">
      {databases.map((db) => (
        <div key={db.id}>
          <ContextMenu>
            <ContextMenuTrigger asChild>
              <button
                onClick={(e) => toggleDb(e, db.id)}
                className={`flex items-center gap-1 w-full hover:bg-muted/50 rounded-sm p-1 text-sm hover:text-foreground ${
                  selectedDbId === db.id ? 'font-semibold text-foreground' : 'text-muted-foreground'
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
              {db.tables.map((table) => (
                <button
                  key={table.id}
                  onClick={() => onSelectTable?.(db.id, table.id)}
                  className="flex items-center gap-1 w-full hover:bg-muted/50 rounded-sm p-1 text-sm text-muted-foreground hover:text-foreground"
                >
                  <Table className="h-4 w-4 shrink-0" />
                  <span className="truncate">{table.name}</span>
                </button>
              ))}
            </div>
          )}
        </div>
      ))}
    </div>
  )
} 