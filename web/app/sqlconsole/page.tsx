'use client'

import { DatabaseList, type DatabaseItem, RelationType, convertRelationType } from "../../components/ui/database-list"
import { useState, useCallback, useEffect, useRef } from 'react'
import { Button } from "../../components/ui/button"
import { Settings, RefreshCw } from 'lucide-react'
import { DatabaseManagement } from "../../components/ui/database-management"
import { SQLEditor, type SQLEditorHandle } from "../../components/ui/sql-editor"
import { DefaultService } from "@/api-gen"
import { toast } from "sonner"

const SELECTED_DB_KEY = 'selected-database-id'

export default function SQLConsole() {
  const [isManagementOpen, setIsManagementOpen] = useState(false)
  const [databases, setDatabases] = useState<DatabaseItem[]>([])
  const [clusters, setClusters] = useState<Array<{ id: string, name: string }>>([])
  const [isLoading, setIsLoading] = useState(true)
  const [editorWidth, setEditorWidth] = useState(0)
  const [isResizing, setIsResizing] = useState(false)
  const [expandedDbs, setExpandedDbs] = useState<Set<string>>(new Set())
  const [isRefreshing, setIsRefreshing] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)
  const editorRef = useRef<SQLEditorHandle>(null)
  const [sidebarWidth, setSidebarWidth] = useState(256) // Default width of the sidebar
  const [selectedDatabaseId, setSelectedDatabaseId] = useState<string | null>(() => {
    if (typeof window !== 'undefined') {
      const savedId = localStorage.getItem(SELECTED_DB_KEY)
      if (savedId) {
        const exists = databases.some(db => db.id === savedId)
        if (!exists) {
          localStorage.removeItem(SELECTED_DB_KEY)
          return null
        }
        return savedId
      }
    }
    return null
  })

  const setSelectedDatabase = useCallback((databaseId: string | null) => {
    setSelectedDatabaseId(databaseId)
    if (databaseId) {
      localStorage.setItem(SELECTED_DB_KEY, databaseId)
    } else {
      localStorage.removeItem(SELECTED_DB_KEY)
    }
  }, [])

  useEffect(() => {
    // Select first database if none is selected and databases are loaded
    if (!selectedDatabaseId && databases.length > 0) {
      const savedId = typeof window !== 'undefined' ? localStorage.getItem(SELECTED_DB_KEY) : null
      const idToSelect = savedId && databases.some(db => db.id === savedId)
        ? savedId
        : databases[0].id
      setSelectedDatabase(idToSelect)
    }
  }, [databases, selectedDatabaseId, setSelectedDatabase])

  const fetchData = useCallback(async () => {
    try {
      setIsLoading(true)
      const [dbData, clusterData] = await Promise.all([
        DefaultService.listDatabases(),
        DefaultService.listClusters()
      ])

      // Transform clusters data
      const transformedClusters = clusterData.map(cluster => ({
        id: String(cluster.ID),
        name: cluster.name
      }))
      setClusters(transformedClusters)

      // Transform databases data
      const transformedDatabases = await Promise.all(dbData.map(async (db) => {
        const cluster = clusterData.find(c => c.ID === db.clusterID)
        return {
          id: String(db.ID),
          name: db.name,
          clusterId: String(db.clusterID),
          clusterName: cluster?.name || 'Unknown Cluster',
          user: db.username,
          password: db.password,
          database: db.database,
          schemas: db.schemas?.map(schema => ({
            name: schema.name,
            relations: schema.relations.map(relation => ({
              id: relation.ID,
              name: relation.name,
              type: convertRelationType(relation.type),
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
      setDatabases(transformedDatabases)
    } catch (error) {
      console.error('Error fetching data:', error)
      toast.error('Failed to load databases and clusters')
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    const calculateEditorWidth = () => {
      if (containerRef.current) {
        const totalWidth = containerRef.current.getBoundingClientRect().width
        setEditorWidth(totalWidth - sidebarWidth - 30)
      }
    }

    calculateEditorWidth()
    window.addEventListener('resize', calculateEditorWidth)

    return () => {
      window.removeEventListener('resize', calculateEditorWidth)
    }
  }, [sidebarWidth])

  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    setIsResizing(true)

    const startX = e.clientX
    const startWidth = sidebarWidth

    const handleMouseMove = (e: MouseEvent) => {
      if (containerRef.current) {
        const newWidth = startWidth + (e.clientX - startX)
        const maxWidth = containerRef.current.getBoundingClientRect().width * 0.8 // 80% of container width
        const minWidth = 200 // Minimum sidebar width

        if (newWidth >= minWidth && newWidth <= maxWidth) {
          setSidebarWidth(newWidth)
        }
      }
    }

    const handleMouseUp = () => {
      setIsResizing(false)
      document.removeEventListener('mousemove', handleMouseMove)
      document.removeEventListener('mouseup', handleMouseUp)
      document.body.style.cursor = 'default'
    }

    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
    document.body.style.cursor = 'ew-resize'
  }, [sidebarWidth])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  const handleSelectTable = useCallback((databaseId: string, tableId: number) => {
    // Handle table selection
    console.log('Selected table:', databaseId, tableId)
  }, [])

  const handleUseDatabase = useCallback((databaseId: string) => {
    setSelectedDatabase(databaseId)
  }, [setSelectedDatabase])

  const handleRunQuery = useCallback(async (query: string) => {
    if (!selectedDatabaseId) {
      return {
        type: 'error' as const,
        message: 'Please select a database first'
      }
    }

    try {
      const result = await DefaultService.queryDatabase(Number(selectedDatabaseId), {
        query,
      })
      if (result.error) {
        return {
          type: 'error' as const,
          message: result.error
        }
      }

      // Check if the query is a DDL operation (case insensitive)
      const isDDL = /^(CREATE|DROP|ALTER)\s/i.test(query.trim())
      if (isDDL && expandedDbs.has(selectedDatabaseId)) {
        // Refresh the database details
        try {
          const dbDetails = await DefaultService.getDatabase(Number(selectedDatabaseId))
          setDatabases(prev => prev.map(d => {
            if (d.id === selectedDatabaseId) {
              return {
                ...d,
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
            return d
          }))
        } catch (error) {
          console.error('Error refreshing database details after DDL:', error)
          // Don't show error toast here as the query itself was successful
        }
      }
      
      return {
        type: 'success' as const,
        message: `Query executed successfully`,
        rows: result.rows,
        columns: result.columns.map(col => col.name),
        rowCount: result.rows?.length ?? 0
      }
    } catch (error) {
      return {
        type: 'error' as const,
        message: error instanceof Error ? error.message : 'Failed to execute query'
      }
    }
  }, [selectedDatabaseId, expandedDbs, setDatabases])

  const handleSaveQuery = useCallback((query: string, name: string) => {
    // Handle query saving
    console.log('Saving query:', name, query)
  }, [])

  const handleDatabaseChange = useCallback(() => {
    fetchData()
    setIsManagementOpen(false)
  }, [fetchData])

  const handleCancelDDL = useCallback(async (ddlId: string) => {
    if (!selectedDatabaseId) return
    try {
      await DefaultService.cancelDdlProgress(Number(selectedDatabaseId), ddlId)
    } catch (error) {
      console.error('Error canceling DDL:', error)
      toast.error('Failed to cancel operation')
    }
  }, [selectedDatabaseId])

  const queryHelper = useCallback((name: string, content: string, executeNow: boolean = false) => {
    if (executeNow) {
      // Execute the query immediately
      editorRef.current?.handleRunQuery(content)
      return
    }

    // Create a new tab and set its content
    editorRef.current?.handleNewTab()
    // Wait for the next tick to ensure the tab is created
    setTimeout(() => {
      editorRef.current?.handleEditorChange(content)
    }, 0)
  }, [])

  const handleToggleDb = useCallback(async (dbId: string) => {
    const newExpanded = new Set(expandedDbs)
    if (newExpanded.has(dbId)) {
      newExpanded.delete(dbId)
      setExpandedDbs(newExpanded)
    } else {
      newExpanded.add(dbId)
      setExpandedDbs(newExpanded)
      
      // Load database details if not already loaded
      const db = databases.find(db => db.id === dbId)
      if (db && !db.schemas) {
        try {
          const dbDetails = await DefaultService.getDatabase(Number(dbId))
          setDatabases(prev => prev.map(d => {
            if (d.id === dbId) {
              return {
                ...d,
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
            return d
          }))
        } catch (error) {
          console.error('Error loading database details:', error)
          toast.error('Failed to load database details')
          newExpanded.delete(dbId)
          setExpandedDbs(newExpanded)
        }
      }
    }
  }, [expandedDbs, databases])

  const handleRefresh = useCallback(async () => {
    try {
      setIsRefreshing(true)
      
      // First get the list of databases and clusters
      const [dbData, clusterData] = await Promise.all([
        DefaultService.listDatabases(),
        DefaultService.listClusters()
      ])

      // Transform clusters data
      const transformedClusters = clusterData.map(cluster => ({
        id: String(cluster.ID),
        name: cluster.name
      }))
      setClusters(transformedClusters)

      // Transform databases data
      const transformedDatabases = dbData.map(db => {
        const cluster = clusterData.find(c => c.ID === db.clusterID)
        return {
          id: String(db.ID),
          name: db.name,
          clusterId: String(db.clusterID),
          clusterName: cluster?.name || 'Unknown Cluster',
          user: db.username,
          password: db.password,
          database: db.database
        }
      })

      // For expanded databases, fetch their details
      const expandedDetails = await Promise.all(
        Array.from(expandedDbs).map(async (dbId) => {
          try {
            const dbDetails = await DefaultService.getDatabase(Number(dbId))
            const db = transformedDatabases.find(d => d.id === dbId)
            if (!db) return null

            return {
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
          } catch (error) {
            console.error(`Error refreshing database ${dbId} details:`, error)
            return null
          }
        })
      )

      // Update databases with the refreshed data
      setDatabases(transformedDatabases.map(db => {
        const expandedDb = expandedDetails.find(d => d?.id === db.id)
        return expandedDb || db
      }))

      toast.success('Database list refreshed')
    } catch (error) {
      console.error('Error refreshing databases:', error)
      toast.error('Failed to refresh database list')
    } finally {
      setIsRefreshing(false)
    }
  }, [expandedDbs])

  return (
    <div ref={containerRef} className="flex h-screen overflow-hidden">
      <div style={{ width: sidebarWidth }} className="flex-shrink-0 border-r bg-muted/30 overflow-hidden flex flex-col">
        <div className="p-4 border-b flex-shrink-0 flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            className="flex-1"
            onClick={() => setIsManagementOpen(true)}
          >
            <Settings className="h-4 w-4 mr-2" />
            Manage Databases
          </Button>
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
        <div className="p-2 overflow-auto flex-1">
          {isLoading ? (
            <div className="text-center py-4 text-muted-foreground">
              Loading databases...
            </div>
          ) : (
            <DatabaseList
              databases={databases}
              onSelectTable={handleSelectTable}
              onUseDatabase={handleUseDatabase}
              queryHelper={queryHelper}
              expandedDbs={expandedDbs}
              onToggleDb={handleToggleDb}
            />
          )}
        </div>
      </div>

      <div
        className="w-[3px] hover:bg-primary/20 active:bg-primary/40 cursor-ew-resize flex-shrink-0"
        onMouseDown={handleMouseDown}
      />

      <div className="flex-1 min-w-0">
        <SQLEditor
          ref={editorRef}
          width={editorWidth}
          savedQueries={[]}
          onRunQuery={handleRunQuery}
          onSaveQuery={handleSaveQuery}
          databaseSchema={[]}
          selectedDatabaseId={selectedDatabaseId}
          onCancelProgress={handleCancelDDL}
        />
      </div>

      <DatabaseManagement
        isOpen={isManagementOpen}
        onClose={() => setIsManagementOpen(false)}
        clusters={clusters}
        onDatabaseChange={handleDatabaseChange}
      />
    </div>
  )
}
