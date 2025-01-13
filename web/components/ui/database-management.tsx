import { Button } from "@/components/ui/button"
import { X, Settings2, Trash2, Plus, ArrowLeft, RefreshCw } from "lucide-react"
import { useCallback, useState, useEffect } from "react"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import * as yup from "yup"
import { useForm } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import { ConfirmationPopup } from "@/components/ui/confirmation-popup"
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select"
import { DefaultService } from "@/api-gen/services/DefaultService"
import { Database } from "@/api-gen/models/Database"
import { DatabaseConnectInfo } from "@/api-gen/models/DatabaseConnectInfo"
import { toast } from "sonner"

interface DatabaseManagementProps {
  isOpen: boolean
  onClose: () => void
  onDatabaseChange?: () => void
}

const databaseSchema = yup.object().shape({
  name: yup.string().required("Name is required"),
  clusterID: yup.number().required("Cluster is required"),
  username: yup.string().required("Username is required"),
  password: yup.string().optional(),
  database: yup.string().required("Database name is required"),
})

export function DatabaseManagement({ isOpen, onClose, onDatabaseChange }: DatabaseManagementProps) {
  const [databases, setDatabases] = useState<Database[]>([])
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [isAddingDatabase, setIsAddingDatabase] = useState(false)
  const [isConfiguring, setIsConfiguring] = useState<string | null>(null)
  const [editingDatabase, setEditingDatabase] = useState<Database | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [isTesting, setIsTesting] = useState(false)
  const [testResult, setTestResult] = useState<{ success: boolean; message: string } | null>(null)
  const [selectedClusterId, setSelectedClusterId] = useState<string>('')
  const [clusters, setClusters] = useState<Array<{ id: string, name: string }>>([])
  const [isRefreshingClusters, setIsRefreshingClusters] = useState(false)

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
    setValue,
    getValues,
  } = useForm({
    resolver: yupResolver(databaseSchema),
    mode: "all",
    defaultValues: {
      name: '',
      username: '',
      password: '',
      clusterID: 0,
      database: '',
    }
  })

  // Fetch databases on component mount
  useEffect(() => {
    const fetchDatabases = async () => {
      try {
        setIsLoading(true)
        const data = await DefaultService.listDatabases()
        setDatabases(data)
      } catch (error) {
        toast.error("Failed to fetch databases")
        console.error("Error fetching databases:", error)
      } finally {
        setIsLoading(false)
      }
    }

    if (isOpen) {
      fetchDatabases()
    }
  }, [isOpen])

  // Fetch clusters on component mount and when refreshing
  const fetchClusters = useCallback(async () => {
    try {
      setIsRefreshingClusters(true)
      const data = await DefaultService.listClusters()
      const transformedClusters = data.map(cluster => ({
        id: String(cluster.ID),
        name: cluster.name
      }))
      setClusters(transformedClusters)
      toast.success("Clusters refreshed")
    } catch (error) {
      toast.error("Failed to refresh clusters")
      console.error("Error fetching clusters:", error)
    } finally {
      setIsRefreshingClusters(false)
    }
  }, [])

  useEffect(() => {
    if (isOpen) {
      fetchClusters()
    }
  }, [isOpen, fetchClusters])

  const handleBackdropClick = useCallback((e: React.MouseEvent) => {
    if (e.target === e.currentTarget && e.button === 0) {
      onClose()
    }
  }, [onClose])

  if (!isOpen) return null

  const handleDeleteDatabase = async (id: number) => {
    try {
      await DefaultService.deleteDatabase(id)
      setDatabases(prev => prev.filter(db => db.ID !== id))
      toast.success("Database deleted successfully")
      onDatabaseChange?.()
    } catch (error) {
      toast.error("Failed to delete database")
      console.error("Error deleting database:", error)
    }
    setDeleteId(null)
  }

  const handleTestConnection = async () => {
    try {
      setIsTesting(true)
      setTestResult(null)
      const formData = getValues()
      const result = await DefaultService.testDatabaseConnection({
        clusterID: Number(formData.clusterID),
        username: formData.username,
        password: formData.password,
        database: formData.database,
      })
      setTestResult({
        success: result.success,
        message: result.result
      })
    } catch (error) {
      setTestResult({
        success: false,
        message: "Connection test failed"
      })
      console.error("Error testing connection:", error)
    } finally {
      setIsTesting(false)
    }
  }

  const onSubmit = async (data: { name: string, username: string, password?: string, clusterID: number, database: string }) => {
    const dbConnectInfo: DatabaseConnectInfo = {
      name: data.name,
      username: data.username,
      password: data.password,
      clusterID: data.clusterID,
      database: data.database,
    }

    try {
      if (isAddingDatabase) {
        const newDb = await DefaultService.createDatabase(dbConnectInfo)
        setDatabases(prev => [...prev, newDb])
        toast.success("Database created successfully")
        setIsAddingDatabase(false)
        onDatabaseChange?.()
      } else if (editingDatabase) {
        const updatedDb = await DefaultService.updateDatabase(editingDatabase.ID, dbConnectInfo)
        setDatabases(prev => prev.map(db => 
          db.ID === editingDatabase.ID ? updatedDb : db
        ))
        toast.success("Database updated successfully")
        setEditingDatabase(null)
        setIsConfiguring(null)
        onDatabaseChange?.()
      }
      reset()
    } catch (error) {
      toast.error(isAddingDatabase ? "Failed to create database" : "Failed to update database")
      console.error("Error saving database:", error)
    }
  }

  const handleConfigureDatabase = async (db: Database) => {
    setEditingDatabase(db)
    setIsConfiguring(String(db.ID))
    setValue('name', db.name)
    setValue('database', db.database)
    
    try {
      setValue('username', db.username || '')
      setValue('password', db.password || '')
      setValue('clusterID', db.clusterID)
      setSelectedClusterId(String(db.clusterID))
    } catch (error) {
      toast.error("Failed to fetch cluster information")
      console.error("Error fetching cluster:", error)
    }
  }

  const isEditing = isAddingDatabase || isConfiguring

  // Update the refresh button click handler
  const handleRefreshClusters = () => {
    fetchClusters()
  }

  return (
    <div 
      className="fixed inset-0 bg-black/50 z-50"
      onMouseDown={handleBackdropClick}
    >
      <div className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] bg-background rounded-lg shadow-lg border">
        <div className="flex items-center justify-between border-b p-4">
          <div className="flex items-center gap-2">
            {isEditing && (
              <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0"
                onClick={() => {
                  setIsAddingDatabase(false)
                  setIsConfiguring(null)
                  setEditingDatabase(null)
                  reset()
                }}
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
            )}
            <h2 className="text-lg font-semibold">
              {isAddingDatabase ? "Add New Database" : 
               isConfiguring ? "Configure Database" : 
               "Database Management"}
            </h2>
          </div>
          <Button
            variant="ghost"
            size="sm"
            className="h-auto p-1.5 hover:bg-muted"
            onClick={onClose}
          >
            <X className="h-4 w-4" />
          </Button>
        </div>

        <div className="p-4 space-y-4">
          {!isEditing ? (
            <>
              <div className="flex justify-between items-center">
                <p className="text-sm text-muted-foreground">
                  Manage your database connections here.
                </p>
                <Button size="sm" onClick={() => setIsAddingDatabase(true)}>
                  <Plus className="h-4 w-4 mr-1" />
                  Add Database
                </Button>
              </div>
              
              {isLoading ? (
                <div className="text-center py-8 text-muted-foreground">
                  Loading databases...
                </div>
              ) : (
                <div className="space-y-2">
                  {databases.map((db) => (
                    <div 
                      key={db.ID}
                      className="flex items-center justify-between p-3 border rounded-lg bg-muted/30 group hover:bg-muted/50"
                    >
                      <div className="space-y-1">
                        <h3 className="font-medium">{db.name}</h3>
                        <div className="text-sm text-muted-foreground">
                          <span>Created: {new Date(db.createdAt).toLocaleString()}</span>
                          <span className="mx-1">·</span>
                          <span>Updated: {new Date(db.updatedAt).toLocaleString()}</span>
                        </div>
                      </div>
                      <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                        <Button 
                          variant="ghost" 
                          size="sm" 
                          className="h-8 w-8 p-0"
                          onClick={() => handleConfigureDatabase(db)}
                        >
                          <Settings2 className="h-4 w-4 text-muted-foreground hover:text-foreground" />
                        </Button>
                        <div className="relative">
                          <Button 
                            variant="ghost" 
                            size="sm" 
                            className="h-8 w-8 p-0 hover:text-red-500"
                            onClick={() => setDeleteId(String(db.ID))}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                          {deleteId === String(db.ID) && (
                            <ConfirmationPopup
                              message="Delete this database?"
                              onConfirm={() => handleDeleteDatabase(db.ID)}
                              onCancel={() => setDeleteId(null)}
                            />
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {!isLoading && databases.length === 0 && (
                <div className="text-center py-8 text-muted-foreground">
                  No databases configured. Click &quot;Add Database&quot; to get started.
                </div>
              )}
            </>
          ) : (
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div className="grid gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Name</Label>
                  <Input
                    id="name"
                    {...register('name')}
                    placeholder="e.g. Production Database"
                    className={errors.name ? 'border-red-500' : ''}
                  />
                  {errors.name && (
                    <p className="text-sm text-red-500">{errors.name.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <Label htmlFor="cluster">Cluster</Label>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-8 w-8 p-0"
                      onClick={handleRefreshClusters}
                      disabled={isRefreshingClusters}
                    >
                      <RefreshCw className={`h-4 w-4 ${isRefreshingClusters ? 'animate-spin' : ''}`} />
                    </Button>
                  </div>
                  <Select
                    onValueChange={(value) => {
                      setValue('clusterID', parseInt(value))
                      setSelectedClusterId(value)
                    }}
                    value={selectedClusterId}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select a cluster" />
                    </SelectTrigger>
                    <SelectContent>
                      {clusters.length === 0 ? (
                        <div className="p-2 text-sm text-center">
                          <p className="text-muted-foreground mb-2">No clusters found</p>
                          <Button
                            variant="link"
                            className="h-auto p-0"
                            onClick={() => window.open('/clusters', '_blank')}
                          >
                            Create a new cluster
                          </Button>
                        </div>
                      ) : (
                        clusters.map((cluster) => (
                          <SelectItem key={cluster.id} value={cluster.id}>
                            {cluster.name}
                          </SelectItem>
                        ))
                      )}
                    </SelectContent>
                  </Select>
                  {errors.clusterID && (
                    <p className="text-sm text-red-500">{errors.clusterID.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="username">Username</Label>
                  <Input
                    id="username"
                    {...register('username')}
                    placeholder="e.g. admin"
                    className={errors.username ? 'border-red-500' : ''}
                  />
                  {errors.username && (
                    <p className="text-sm text-red-500">{errors.username.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="password">Password</Label>
                  <Input
                    id="password"
                    type="password"
                    {...register('password')}
                    placeholder="Enter password"
                    className={errors.password ? 'border-red-500' : ''}
                  />
                  {errors.password && (
                    <p className="text-sm text-red-500">{errors.password.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="database">Database Name</Label>
                  <Input
                    id="database"
                    {...register('database')}
                    placeholder="e.g. postgres"
                    className={errors.database ? 'border-red-500' : ''}
                  />
                  {errors.database && (
                    <p className="text-sm text-red-500">{errors.database.message}</p>
                  )}
                </div>

                <Button
                  type="button"
                  variant="outline"
                  onClick={handleTestConnection}
                  disabled={isTesting}
                >
                  {isTesting ? "Testing..." : "Test Connection"}
                </Button>

                {testResult && (
                  <div className={`p-3 rounded-md ${
                    testResult.success 
                      ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' 
                      : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
                  }`}>
                    {testResult.message}
                  </div>
                )}
              </div>
            </form>
          )}
        </div>

        <div className="flex justify-end gap-2 border-t p-4">
          {isEditing ? (
            <>
              <Button 
                variant="outline" 
                onClick={() => {
                  setIsAddingDatabase(false)
                  setIsConfiguring(null)
                  setEditingDatabase(null)
                  reset()
                }}
              >
                Cancel
              </Button>
              <Button 
                onClick={handleSubmit(onSubmit)}
                type="submit"
              >
                {isAddingDatabase ? 'Add Database' : 'Save Changes'}
              </Button>
            </>
          ) : (
            <Button variant="outline" onClick={onClose}>
              Close
            </Button>
          )}
        </div>
      </div>
    </div>
  )
} 