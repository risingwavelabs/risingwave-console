'use client'

import { useState } from 'react'
import { ChevronRight, Database, Table } from 'lucide-react'
import { cn } from '@/lib/utils'

interface DatabaseItem {
  id: string
  name: string
  tables: { id: string; name: string }[]
}

interface DatabaseListProps {
  databases: DatabaseItem[]
  onSelectTable?: (databaseId: string, tableId: string) => void
}

export function DatabaseList({ databases, onSelectTable }: DatabaseListProps) {
  const [expandedDbs, setExpandedDbs] = useState<Set<string>>(new Set())

  const toggleDatabase = (dbId: string) => {
    const newExpanded = new Set(expandedDbs)
    if (newExpanded.has(dbId)) {
      newExpanded.delete(dbId)
    } else {
      newExpanded.add(dbId)
    }
    setExpandedDbs(newExpanded)
  }

  return (
    <div className="space-y-1">
      {databases.map((db) => (
        <div key={db.id} className="space-y-0.5">
          <button
            onClick={() => toggleDatabase(db.id)}
            className="flex items-center w-full text-sm hover:bg-muted/50 rounded-md"
          >
            <ChevronRight
              className={cn(
                "h-4 w-4 shrink-0 transition-transform",
                expandedDbs.has(db.id) && "transform rotate-90"
              )}
            />
            <Database className="h-4 w-4 shrink-0 mx-2" />
            <span className="truncate">{db.name}</span>
          </button>

          {expandedDbs.has(db.id) && (
            <div className="ml-6 space-y-1">
              {db.tables.map((table) => (
                <button
                  key={table.id}
                  onClick={() => onSelectTable?.(db.id, table.id)}
                  className="flex items-center w-full  text-sm hover:bg-muted/50 rounded-md"
                >
                  <Table className="h-4 w-4 shrink-0 mr-2" />
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