'use client'

import { useState, useCallback, useEffect } from 'react'
import { ChevronDown, ChevronRight, Database, Table, ArrowDownToLine, ArrowUpFromLine, Waves, RefreshCw, Glasses, Workflow, Wallpaper, Monitor, Touchpad } from 'lucide-react'
import { DefaultService } from '@/api-gen'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from "./context-menu"
import { Button } from "./button"
import { toast } from "sonner"

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

interface DatabaseListProps {
  databases: DatabaseItem[]
  onSelectTable?: (databaseId: string, tableId: number) => void
  onUseDatabase?: (databaseId: string) => void
}

const SELECTED_DB_KEY = 'selectedDatabaseId'

const SYSTEM_SCHEMAS = ['pg_catalog', 'rw_catalog', 'information_schema']

// Helper function to convert relation type
export const convertRelationType = (type: string): RelationType => {
  if (type === 'materializedView') return RelationType.MaterializedView
  if (type === 'systemTable') return RelationType.SystemTable
  return type as RelationType
}

export function DatabaseList({ databases, onSelectTable, onUseDatabase }: DatabaseListProps) {
  const [expandedDbs, setExpandedDbs] = useState<Set<string>>(new Set())
  const [expandedSchemas, setExpandedSchemas] = useState<Set<string>>(new Set())
  const [expandedRelations, setExpandedRelations] = useState<Set<string>>(new Set())
  const [selectedDbId, setSelectedDbId] = useState<string | null>(null)
  const [loadedDatabases, setLoadedDatabases] = useState<Record<string, DatabaseItem>>({})
  const [isRefreshing, setIsRefreshing] = useState(false)
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
      const newExpanded = new Set(expandedDbs)
      if (newExpanded.has(dbId)) {
        newExpanded.delete(dbId)
      } else {
        newExpanded.add(dbId)
        // If we haven't loaded this database's details yet, load them
        if (!loadedDatabases[dbId]) {
          try {
            const dbDetails = await DefaultService.getDatabase(Number(dbId))
            setLoadedDatabases(prev => ({
              ...prev,
              [dbId]: {
                ...databases.find(db => db.id === dbId)!,
                schemas: dbDetails.schemas?.map(schema => ({
                  name: schema.name,
                  relations: schema.relations.map(relation => ({
                    id: relation.ID,
                    name: relation.name,
                    type: convertRelationType(relation.type.toLowerCase()),
                    columns: relation.columns.map(col => ({
                      name: col.name,
                      type: col.type,
                      isHidden: col.isHidden,
                      isPrimaryKey: col.isPrimaryKey
                    }))
                  }))
                }))
              }
            }))
          } catch (error) {
            console.error('Error loading database details:', error)
          }
        }
      }
      setExpandedDbs(newExpanded)
    }
  }, [expandedDbs, databases, loadedDatabases])

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
    const displayDb = loadedDatabases[db.id] || db
    if (displayDb.schemas) {
      return {
        ...displayDb,
        schemas: sortSchemas(displayDb.schemas).map(schema => ({
          ...schema,
          relations: sortRelations(schema.relations)
        }))
      }
    }
    return displayDb
  }, [loadedDatabases, sortSchemas, sortRelations])

  const handleRefresh = useCallback(async () => {
    try {
      setIsRefreshing(true)
      
      // First get the list of databases and clusters
      const [dbData, clusterData] = await Promise.all([
        DefaultService.listDatabases(),
        DefaultService.listClusters()
      ])
      
      // Transform the new database data
      const transformedDatabases = dbData.map(db => {
        const cluster = clusterData.find(c => c.ID === db.clusterID)
        return {
          id: String(db.ID),
          name: db.name,
          clusterId: String(db.clusterID),
          clusterName: cluster?.name || 'Unknown Cluster',
          user: db.username,
          password: db.password,
          database: db.database,
          // Don't copy old schemas, we'll fetch fresh ones for expanded DBs
          schemas: undefined
        }
      })

      // Update local databases first
      setLocalDatabases(transformedDatabases)

      // Then immediately fetch details for all expanded databases
      const expandedDetails = await Promise.all(
        Array.from(expandedDbs).map(async (dbId) => {
          try {
            const dbDetails = await DefaultService.getDatabase(Number(dbId))
            const db = transformedDatabases.find(d => d.id === dbId)
            if (!db) return null

            return {
              dbId,
              details: {
                ...db,
                schemas: dbDetails.schemas?.map(schema => ({
                  name: schema.name,
                  relations: schema.relations.map(relation => ({
                    id: relation.ID,
                    name: relation.name,
                    type: convertRelationType(relation.type.toLowerCase()),
                    columns: relation.columns.map(col => ({
                      name: col.name,
                      type: col.type,
                      isHidden: col.isHidden,
                      isPrimaryKey: col.isPrimaryKey
                    }))
                  }))
                }))
              }
            }
          } catch (error) {
            console.error(`Error refreshing database ${dbId} details:`, error)
            return null
          }
        })
      )

      // Update loadedDatabases with fresh details
      const newLoadedDatabases: Record<string, DatabaseItem> = {}
      expandedDetails.forEach(detail => {
        if (detail) {
          newLoadedDatabases[detail.dbId] = detail.details
        }
      })
      setLoadedDatabases(newLoadedDatabases)

      toast.success('Database list refreshed')
    } catch (error) {
      console.error('Error refreshing databases:', error)
      toast.error('Failed to refresh database list')
    } finally {
      setIsRefreshing(false)
    }
  }, [expandedDbs])

  return (
    <div className="space-y-1">
      <div className="p-2 border-b flex items-center justify-between">
        <span className="text-sm font-medium">Databases</span>
        <Button
          variant="ghost"
          size="icon"
          onClick={handleRefresh}
          disabled={isRefreshing}
          className={isRefreshing ? 'animate-spin' : ''}
        >
          <RefreshCw className="h-4 w-4" />
        </Button>
      </div>
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
                {displayDb.schemas?.map((schema) => (
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
                        {schema.relations.map((relation) => (
                          <div key={relation.id} className="space-y-1">
                            <button
                              onClick={(e) => {
                                toggleRelation(e, relation.id)
                                onSelectTable?.(db.id, relation.id)
                              }}
                              className={`flex items-center gap-1 w-full hover:bg-muted/50 rounded-sm p-1 text-sm ${
                                selectedDbId === db.id ? 'text-foreground' : 'text-muted-foreground'
                              } hover:text-foreground`}
                            >
                              {expandedRelations.has(String(relation.id)) ? (
                                <ChevronDown className="h-4 w-4 shrink-0" />
                              ) : (
                                <ChevronRight className="h-4 w-4 shrink-0" />
                              )}
                              {relation.type === RelationType.Table && <Table className="h-4 w-4 shrink-0 text-yellow-500" />}
                              {relation.type === RelationType.Source && <ArrowDownToLine className="h-4 w-4 shrink-0 rotate-270" />}
                              {relation.type === RelationType.Sink && <ArrowUpFromLine className="h-4 w-4 shrink-0 rotate-90" />}
                              {relation.type === RelationType.MaterializedView && <Workflow className="h-4 w-4 shrink-0 text-blue-500" />}
                              {relation.type === RelationType.SystemTable && <Table className="h-4 w-4 shrink-0" />}
                              {relation.type === RelationType.View && <Touchpad className="h-4 w-4 shrink-0" />}
                              <span className="truncate">{relation.name}</span>
                            </button>
                            {expandedRelations.has(String(relation.id)) && (
                              <div className="ml-6 space-y-1">
                                {relation.columns.map((column, index) => (
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
