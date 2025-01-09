'use client'

import { DatabaseList, type DatabaseItem } from "../../components/ui/database-list"
import { useState, useCallback, useEffect, useRef } from 'react'
import { Button } from "../../components/ui/button"
import { Settings } from 'lucide-react'
import { DatabaseManagement } from "../../components/ui/database-management"
import { SQLEditor } from "../../components/ui/sql-editor"
import { DefaultService } from "@/api-gen"
import { toast } from "sonner"

const SELECTED_DB_KEY = 'selected-database-id'

export default function SQLConsole() {
  const [isManagementOpen, setIsManagementOpen] = useState(false)
  const [databases, setDatabases] = useState<DatabaseItem[]>([])
  const [clusters, setClusters] = useState<Array<{ id: string, name: string }>>([])
  const [isLoading, setIsLoading] = useState(true)
  const [editorWidth, setEditorWidth] = useState(0)
  const containerRef = useRef<HTMLDivElement>(null)
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
          tables: db.relations?.map(relation => ({
            id: relation.ID,
            name: relation.name,
            type: relation.type.toLowerCase() as 'table' | 'source' | 'sink' | 'materialized_view',
            columns: relation.columns.map(col => ({
              // id: col.ID,
              name: col.name,
              type: col.type
            }))
          })) || []
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
        setEditorWidth(totalWidth - 256) // 256px is the width of the sidebar (w-64)
      }
    }

    calculateEditorWidth()
    window.addEventListener('resize', calculateEditorWidth)

    return () => {
      window.removeEventListener('resize', calculateEditorWidth)
    }
  }, [])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  const handleSelectTable = useCallback((databaseId: string, tableId: string) => {
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
      return {
        type: 'success' as const,
        message: `Query executed successfully`,
        rows: result.rows,
      }
    } catch (error) {
      return {
        type: 'error' as const,
        message: error instanceof Error ? error.message : 'Failed to execute query'
      }
    }
  }, [selectedDatabaseId])

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

  return (
    <div ref={containerRef} className="flex h-screen">
      <div className="w-64 border-r bg-muted/30">
        <div className="p-4 border-b">
          <Button
            variant="outline"
            size="sm"
            className="w-full"
            onClick={() => setIsManagementOpen(true)}
          >
            <Settings className="h-4 w-4 mr-2" />
            Manage Databases
          </Button>
        </div>
        <div className="p-2">
          {isLoading ? (
            <div className="text-center py-4 text-muted-foreground">
              Loading databases...
            </div>
          ) : (
            <DatabaseList
              databases={databases}
              onSelectTable={handleSelectTable}
              onUseDatabase={handleUseDatabase}
            />
          )}
        </div>
      </div>

      <div className="flex-1 min-w-0">
        <SQLEditor
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
