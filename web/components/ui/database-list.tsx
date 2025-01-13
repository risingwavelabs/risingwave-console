'use client'

import { useState, useCallback, useEffect } from 'react'
import { ChevronDown, ChevronRight, Database, Table, ArrowDownToLine, ArrowUpFromLine, Workflow, Touchpad, Copy, Check } from 'lucide-react'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "./context-menu"
import { toast } from "sonner"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./tooltip"

export interface Column {
  name: string
  type: string
  isHidden: boolean
  isPrimaryKey: boolean
}

export enum RelationType {
  Table = 'table',
  Source = 'source',
  Sink = 'sink',
  MaterializedView = 'materialized view',
  SystemTable = 'system table',
  View = 'view'
}

export interface Relation {
  id: number
  name: string
  columns: Column[]
  type: RelationType
}

export interface Schema {
  name: string
  relations: Relation[]
}

export interface DatabaseItem {
  id: string
  name: string
  clusterId: string
  clusterName: string
  user: string
  password?: string
  database: string
  schemas?: Schema[]
}

export interface DatabaseListProps {
  databases: DatabaseItem[]
  onSelectTable?: (databaseId: string, tableId: number) => void
  onUseDatabase?: (databaseId: string) => void
  queryHelper?: (name: string, content: string, executeNow?: boolean) => void
  expandedDbs: Set<string>
  onToggleDb: (dbId: string) => void
}

const SELECTED_DB_KEY = 'selectedDatabaseId'

const SYSTEM_SCHEMAS = ['pg_catalog', 'rw_catalog', 'information_schema']

// Helper function to convert relation type
export const convertRelationType = (type: string): RelationType => {
  if (type === 'materializedView') return RelationType.MaterializedView
  if (type === 'systemTable') return RelationType.SystemTable
  return type as RelationType
}

export function DatabaseList({ databases, onSelectTable, onUseDatabase, queryHelper, expandedDbs, onToggleDb }: DatabaseListProps) {
  const [expandedSchemas, setExpandedSchemas] = useState<Set<string>>(new Set())
  const [expandedRelations, setExpandedRelations] = useState<Set<string>>(new Set())
  const [selectedDbId, setSelectedDbId] = useState<string | null>(null)
  const [copiedId, setCopiedId] = useState<string | null>(null)
  const [localDatabases, setLocalDatabases] = useState(databases)

  useEffect(() => {
    setLocalDatabases(databases)
  }, [databases])

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

  const toggleDb = useCallback(async (e: React.MouseEvent, dbId: string) => {
    if (e.button === 0) { // Left click only
      onToggleDb(dbId)
    }
  }, [onToggleDb])

  const toggleSchema = useCallback((e: React.MouseEvent, schemaName: string) => {
    e.stopPropagation()
    if (e.button === 0) { // Left click only
      const newExpanded = new Set(expandedSchemas)
      if (newExpanded.has(schemaName)) {
        newExpanded.delete(schemaName)
      } else {
        newExpanded.add(schemaName)
      }
      setExpandedSchemas(newExpanded)
    }
  }, [expandedSchemas])

  const toggleRelation = useCallback((e: React.MouseEvent, relationId: number) => {
    e.stopPropagation()
    if (e.button === 0) { // Left click only
      const newExpanded = new Set(expandedRelations)
      if (newExpanded.has(String(relationId))) {
        newExpanded.delete(String(relationId))
      } else {
        newExpanded.add(String(relationId))
      }
      setExpandedRelations(newExpanded)
    }
  }, [expandedRelations])

  // Schema sorting function
  const sortSchemas = useCallback((schemas: Schema[]) => {
    return [...schemas].sort((a, b) => {
      // 'public' schema always comes first
      if (a.name === 'public') return -1
      if (b.name === 'public') return 1

      // System schemas go last
      const aIsSystem = SYSTEM_SCHEMAS.includes(a.name)
      const bIsSystem = SYSTEM_SCHEMAS.includes(b.name)
      if (aIsSystem && !bIsSystem) return 1
      if (!aIsSystem && bIsSystem) return -1
      if (aIsSystem && bIsSystem) {
        // Sort system schemas among themselves
        return SYSTEM_SCHEMAS.indexOf(a.name) - SYSTEM_SCHEMAS.indexOf(b.name)
      }

      // Everything else is sorted alphabetically
      return a.name.localeCompare(b.name)
    })
  }, [])

  // Sort relations alphabetically
  const sortRelations = useCallback((relations: Relation[]) => {
    return [...relations].sort((a, b) => a.name.localeCompare(b.name))
  }, [])

  // Update the getDisplayDatabase function to sort schemas and relations
  const getDisplayDatabase = useCallback((db: DatabaseItem) => {
    if (db.schemas) {
      return {
        ...db,
        schemas: sortSchemas(db.schemas).map(schema => ({
          ...schema,
          relations: sortRelations(schema.relations)
        }))
      }
    }
    return db
  }, [sortSchemas, sortRelations])

  const handleCopy = useCallback(async (e: React.MouseEvent, text: string, id: string) => {
    e.stopPropagation()
    try {
      await navigator.clipboard.writeText(text)
      setCopiedId(id)
      setTimeout(() => setCopiedId(null), 2000)
      toast.success('Copied to clipboard')
    } catch {
      toast.error('Failed to copy to clipboard')
    }
  }, [])

  const handleDropRelation = useCallback((relation: Relation, schema: Schema, cascade: boolean) => {
    const statement = `DROP ${relation.type.toUpperCase()} ${schema.name}.${relation.name} ${cascade ? 'CASCADE' : ''};`
    queryHelper?.(`Drop ${relation.name}`, statement)
  }, [queryHelper])

  const handleViewData = useCallback((relation: Relation, schema: Schema) => {
    const statement = `SELECT * FROM ${schema.name}.${relation.name} LIMIT 100;`
    queryHelper?.(`View ${relation.name}`, statement, true)
  }, [queryHelper])

  return (
    <div className="space-y-1">
      {localDatabases.map((db) => {
        const displayDb = getDisplayDatabase(db)
        return (
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
                  <span className="truncate">{`${displayDb.name} (${displayDb.clusterName})`}</span>
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
                {displayDb.schemas?.map((schema: Schema) => (
                  <div key={schema.name} className="space-y-1">
                    <button
                      onClick={(e) => toggleSchema(e, schema.name)}
                      className={`flex items-center gap-1 w-full hover:bg-muted/50 rounded-sm p-1 text-sm ${
                        selectedDbId === db.id ? 'text-foreground' : 'text-muted-foreground'
                      } hover:text-foreground`}
                    >
                      {expandedSchemas.has(schema.name) ? (
                        <ChevronDown className="h-4 w-4 shrink-0" />
                      ) : (
                        <ChevronRight className="h-4 w-4 shrink-0" />
                      )}
                      <span className="truncate">{schema.name}</span>
                    </button>
                    {expandedSchemas.has(schema.name) && (
                      <div className="ml-6 space-y-1">
                        {schema.relations.map((relation: Relation) => (
                          <div key={relation.id} className="space-y-1">
                            <ContextMenu>
                              <div className="flex items-center gap-1 w-full group">
                                <ContextMenuTrigger asChild>
                                  <div
                                    onClick={(e) => {
                                      toggleRelation(e, relation.id)
                                      onSelectTable?.(db.id, relation.id)
                                    }}
                                    className={`flex items-center gap-1 flex-1 hover:bg-muted/50 rounded-sm p-1 text-sm cursor-pointer ${
                                      selectedDbId === db.id ? 'text-foreground' : 'text-muted-foreground'
                                    } hover:text-foreground`}
                                  >
                                    {expandedRelations.has(String(relation.id)) ? (
                                      <ChevronDown className="h-4 w-4 shrink-0" />
                                    ) : (
                                      <ChevronRight className="h-4 w-4 shrink-0" />
                                    )}
                                    {relation.type === RelationType.Table && (
                                      <TooltipProvider>
                                        <Tooltip delayDuration={0}>
                                          <TooltipTrigger>
                                            <Table className="h-4 w-4 shrink-0 text-yellow-500" />
                                          </TooltipTrigger>
                                          <TooltipContent>Table</TooltipContent>
                                        </Tooltip>
                                      </TooltipProvider>
                                    )}
                                    {relation.type === RelationType.Source && (
                                      <TooltipProvider>
                                        <Tooltip delayDuration={0}>
                                          <TooltipTrigger>
                                            <ArrowDownToLine className="h-4 w-4 shrink-0 rotate-270" />
                                          </TooltipTrigger>
                                          <TooltipContent>Source</TooltipContent>
                                        </Tooltip>
                                      </TooltipProvider>
                                    )}
                                    {relation.type === RelationType.Sink && (
                                      <TooltipProvider>
                                        <Tooltip delayDuration={0}>
                                          <TooltipTrigger>
                                            <ArrowUpFromLine className="h-4 w-4 shrink-0 rotate-90" />
                                          </TooltipTrigger>
                                          <TooltipContent>Sink</TooltipContent>
                                        </Tooltip>
                                      </TooltipProvider>
                                    )}
                                    {relation.type === RelationType.MaterializedView && (
                                      <TooltipProvider>
                                        <Tooltip delayDuration={0}>
                                          <TooltipTrigger>
                                            <Workflow className="h-4 w-4 shrink-0 text-blue-500" />
                                          </TooltipTrigger>
                                          <TooltipContent>Materialized View</TooltipContent>
                                        </Tooltip>
                                      </TooltipProvider>
                                    )}
                                    {relation.type === RelationType.SystemTable && (
                                      <TooltipProvider>
                                        <Tooltip delayDuration={0}>
                                          <TooltipTrigger>
                                            <Table className="h-4 w-4 shrink-0" />
                                          </TooltipTrigger>
                                          <TooltipContent>System Table</TooltipContent>
                                        </Tooltip>
                                      </TooltipProvider>
                                    )}
                                    {relation.type === RelationType.View && (
                                      <TooltipProvider>
                                        <Tooltip delayDuration={0}>
                                          <TooltipTrigger>
                                            <Touchpad className="h-4 w-4 shrink-0" />
                                          </TooltipTrigger>
                                          <TooltipContent>View</TooltipContent>
                                        </Tooltip>
                                      </TooltipProvider>
                                    )}
                                    <span className="truncate">{relation.name}</span>
                                  </div>
                                </ContextMenuTrigger>
                                <div
                                  onClick={(e) => handleCopy(e, relation.name, String(relation.id))}
                                  className="p-1 opacity-0 group-hover:opacity-100 hover:text-foreground transition-opacity cursor-pointer"
                                >
                                  {copiedId === String(relation.id) ? (
                                    <Check className="h-3 w-3" />
                                  ) : (
                                    <Copy className="h-3 w-3" />
                                  )}
                                </div>
                              </div>
                              <ContextMenuContent className="w-48">
                                <ContextMenuItem onClick={() => handleViewData(relation, schema)}>
                                  View Data (100 rows)
                                </ContextMenuItem>
                                <ContextMenuItem onClick={() => handleDropRelation(relation, schema, false)}>
                                  Drop
                                </ContextMenuItem>
                                <ContextMenuItem onClick={() => handleDropRelation(relation, schema, true)}>
                                  Drop Cascade
                                </ContextMenuItem>
                              </ContextMenuContent>
                            </ContextMenu>
                            {expandedRelations.has(String(relation.id)) && (
                              <div className="ml-6 space-y-1">
                                {relation.columns.map((column: Column, index: number) => (
                                  <div
                                    key={index}
                                    className={`flex items-center gap-1 w-full rounded-sm p-1 text-sm ${
                                      selectedDbId === db.id ? 'text-foreground' : 'text-muted-foreground'
                                    } ${column.isHidden ? 'opacity-50' : ''}`}
                                  >
                                    {column.isPrimaryKey && <span className="text-xs">🔑</span>}
                                    <span className="truncate">{column.name}</span>
                                    <span className={`text-xs ${selectedDbId === db.id ? 'text-muted-foreground/80' : 'text-muted-foreground'}`}>({column.type})</span>
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
            )}
          </div>
        )
      })}
    </div>
  )
}
