'use client'

import { DatabaseList, type DatabaseItem, convertRelationType } from "../../components/ui/database-list"
import { useState, useCallback, useEffect, useRef } from 'react'
import { Button } from "../../components/ui/button"
import { Settings, RefreshCw } from 'lucide-react'
import { DatabaseManagement } from "../../components/ui/database-management"
import { SQLEditor, type SQLEditorHandle } from "../../components/ui/sql-editor"
import { DefaultService } from "@/api-gen"
import { toast } from "sonner"
import { RisingWaveNodeData, NodeType } from "@/components/streaming-graph"
import type { Relation as APIRelation } from "@/api-gen/models/Relation"
import type { Schema as APISchema } from "@/api-gen/models/Schema"
import type { Database as APIDatabase } from "@/api-gen/models/Database"

const SELECTED_DB_KEY = 'selected-database-id'

// Utility function to convert API schema to UI schema
const convertDatabaseSchema = (dbDetails: APIDatabase): DatabaseItem['schemas'] => {
  return dbDetails.schemas?.map(schema => ({
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

// Utility function to convert API schema to streaming graph data
const convertToStreamingGraph = (schema: APISchema): RisingWaveNodeData[] => {
  return schema.relations.map((relation: APIRelation) => ({
    id: relation.ID,
    name: relation.name,
    type: relation.type.toLowerCase() as NodeType,
    columns: relation.columns.map(col => ({
      name: col.name,
      type: col.type,
      isPrimary: col.isPrimaryKey
    })),
    dependencies: relation.dependencies
  }))
}

// Utility function to fetch and transform databases and clusters data
const fetchDatabasesAndClusters = async () => {
  const [dbData, clusterData] = await Promise.all([
    DefaultService.listDatabases(),
    DefaultService.listClusters()
  ])

  // Transform databases data
  return dbData.map(db => {
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
}

export default function SQLConsole() {
  const [isManagementOpen, setIsManagementOpen] = useState(false)
  const [databases, setDatabases] = useState<DatabaseItem[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [editorWidth, setEditorWidth] = useState(0)
  const [expandedDbs, setExpandedDbs] = useState<Set<string>>(new Set())
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [databaseSchema, setDatabaseSchema] = useState<RisingWaveNodeData[]>([])
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
      const transformedDatabases = await fetchDatabasesAndClusters()

      // Get the database ID to select
      const savedId = typeof window !== 'undefined' ? localStorage.getItem(SELECTED_DB_KEY) : null
      const idToSelect = savedId && transformedDatabases.some(db => db.id === savedId)
        ? savedId
        : transformedDatabases[0]?.id

      // If we have a database to select, fetch its details
      if (idToSelect) {
        try {
          const dbDetails = await DefaultService.getDatabase(Number(idToSelect))
          
          // Update the database with schema details
          const updatedDatabases = transformedDatabases.map(db => {
            if (db.id === idToSelect) {
              return {
                ...db,
                schemas: convertDatabaseSchema(dbDetails)
              }
            }
            return db
          })
          setDatabases(updatedDatabases)

          // Update the streaming graph if there's a public schema
          const publicSchema = dbDetails.schemas?.find(s => s.name === 'public')
          if (publicSchema) {
            setDatabaseSchema(convertToStreamingGraph(publicSchema))
          }

          // Set the selected database ID
          setSelectedDatabase(idToSelect)
        } catch (error) {
          console.error('Error loading selected database details:', error)
          setDatabases(transformedDatabases)
          setSelectedDatabase(idToSelect)
        }
      } else {
        setDatabases(transformedDatabases)
      }
    } catch (error) {
      console.error('Error fetching data:', error)
      toast.error('Failed to load databases and clusters')
    } finally {
      setIsLoading(false)
    }
  }, [setSelectedDatabase, setDatabases, setDatabaseSchema, setIsLoading])

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
    
    // Update streaming graph data for the selected database
    const db = databases.find(db => db.id === databaseId)
    if (db?.schemas) {
      const publicSchema = db.schemas.find(s => s.name === 'public')
      if (publicSchema) {
        const graphData = convertToStreamingGraph(publicSchema as unknown as APISchema)
        setDatabaseSchema(graphData)
      } else {
        setDatabaseSchema([])
      }
    } else {
      setDatabaseSchema([])
    }
  }, [setSelectedDatabase, databases])

  const handleRunQuery = useCallback(async (query: string, backgroundDDL: boolean = false) => {
    if (!selectedDatabaseId) {
      return {
        type: 'error' as const,
        message: 'Please select a database first'
      }
    }

    try {
      const result = await DefaultService.queryDatabase(Number(selectedDatabaseId), {
        query,
        backgroundDDL,
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
                schemas: convertDatabaseSchema(dbDetails)
              }
            }
            return d
          }))

          // Update streaming graph data for the selected database
          const publicSchema = dbDetails.schemas?.find(s => s.name === 'public')
          if (publicSchema) {
            setDatabaseSchema(convertToStreamingGraph(publicSchema))
          } else {
            setDatabaseSchema([])
          }
        } catch (error) {
          console.error('Error refreshing database details after DDL:', error)
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

  const handleDatabaseChange = useCallback(() => {
    fetchData()
    setIsManagementOpen(false)
  }, [fetchData])

  const handleCancelDDL = useCallback(async (ddlId: number) => {
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
      if (dbId === selectedDatabaseId) {
        setDatabaseSchema([])
      }
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
                schemas: convertDatabaseSchema(dbDetails)
              }
            }
            return d
          }))

          // Update streaming graph data if this is the selected database
          if (dbId === selectedDatabaseId) {
            const publicSchema = (dbDetails.schemas as APISchema[])?.find(s => s.name === 'public')
            if (publicSchema) {
              setDatabaseSchema(convertToStreamingGraph(publicSchema))
            }
          }
        } catch (error) {
          console.error('Error loading database details:', error)
          toast.error('Failed to load database details')
          newExpanded.delete(dbId)
          setExpandedDbs(newExpanded)
        }
      }
    }
  }, [expandedDbs, databases, selectedDatabaseId, setDatabases, setDatabaseSchema, setExpandedDbs])

  const handleRefresh = useCallback(async () => {
    try {
      setIsRefreshing(true)
      const transformedDatabases = await fetchDatabasesAndClusters()

      // Create a set of database IDs to fetch details for (expanded + selected)
      const dbsToFetch = new Set([...expandedDbs])
      if (selectedDatabaseId) {
        dbsToFetch.add(selectedDatabaseId)
      }

      // For expanded databases and selected database, fetch their details
      const expandedDetails = await Promise.all(
        Array.from(dbsToFetch).map(async (dbId) => {
          try {
            const dbDetails = await DefaultService.getDatabase(Number(dbId))
            const db = transformedDatabases.find(d => d.id === dbId)
            if (!db) return null

            // If this is the selected database, update the streaming graph
            if (dbId === selectedDatabaseId) {
              const publicSchema = dbDetails.schemas?.find(s => s.name === 'public')
              if (publicSchema) {
                setDatabaseSchema(convertToStreamingGraph(publicSchema))
              } else {
                setDatabaseSchema([])
              }
            }

            return {
              ...db,
              schemas: convertDatabaseSchema(dbDetails)
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
  }, [expandedDbs, selectedDatabaseId, setDatabases, setDatabaseSchema, setIsRefreshing])

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
          onRunQuery={handleRunQuery}
          databaseSchema={databaseSchema}
          selectedDatabaseId={selectedDatabaseId}
          onCancelProgress={handleCancelDDL}
        />
      </div>

      <DatabaseManagement
        isOpen={isManagementOpen}
        onClose={() => setIsManagementOpen(false)}
        onDatabaseChange={handleDatabaseChange}
      />
    </div>
  )
}
