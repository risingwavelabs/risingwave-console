'use client'

import { DatabaseList, type DatabaseItem } from "../../components/ui/database-list"
import { useState, useCallback, useEffect, useRef } from 'react'
import { Button } from "../../components/ui/button"
import { Settings } from 'lucide-react'
import { DatabaseManagement } from "../../components/ui/database-management"
import { SQLEditor } from "../../components/ui/sql-editor"
import { RisingWaveNodeData } from "../../components/streaming-graph"
import { DefaultService } from "@/api-gen"
import { Cluster } from "@/api-gen/models/Cluster"
import { Database } from "@/api-gen/models/Database"
import { toast } from "sonner"

// Sample saved queries - replace with real data from your backend
const savedQueries = [
  { id: "q1", name: "Get Active Users" },
  { id: "q2", name: "Monthly Revenue" },
  { id: "q3", name: "Product Inventory" },
]

export default function SQLConsole() {
  const [isManagementOpen, setIsManagementOpen] = useState(false)
  const [databases, setDatabases] = useState<DatabaseItem[]>([])
  const [clusters, setClusters] = useState<Array<{ id: string, name: string }>>([])
  const [isLoading, setIsLoading] = useState(true)
  const [editorWidth, setEditorWidth] = useState(0)
  const containerRef = useRef<HTMLDivElement>(null)

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
          database: db.name,
          tables: db.relations?.map(relation => ({
            id: relation.ID,
            name: relation.name,
            type: relation.type.toLowerCase() as 'table' | 'source' | 'sink' | 'materialized_view',
            columns: relation.columns.map(col => ({
              id: col.ID,
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
    // Handle database selection
    console.log('Using database:', databaseId)
  }, [])

  const handleRunQuery = useCallback((query: string) => {
    // Handle query execution
    console.log('Running query:', query)
  }, [])

  const handleSaveQuery = useCallback((query: string, name: string) => {
    // Handle query saving
    console.log('Saving query:', name, query)
  }, [])

  const handleDatabaseChange = useCallback(() => {
    fetchData()
    setIsManagementOpen(false)
  }, [fetchData])

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
          savedQueries={savedQueries}
          onRunQuery={handleRunQuery}
          onSaveQuery={handleSaveQuery}
          databaseSchema={[]} // We'll populate this later with real schema data
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
