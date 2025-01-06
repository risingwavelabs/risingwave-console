import { Button } from "@/components/ui/button"
import { X, Settings2, Trash2, Plus, ArrowLeft } from "lucide-react"
import { useCallback, useState } from "react"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import * as yup from "yup"
import { useForm } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import { ConfirmationPopup } from "@/components/ui/confirmation-popup"
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select"

interface Database {
  id: string
  name: string
  clusterId: string
  clusterName: string
  user: string
  password?: string
  database: string
}

const databaseSchema = yup.object().shape({
  name: yup.string().required("Name is required"),
  clusterId: yup.string().required("Cluster is required"),
  user: yup.string().required("Username is required"),
  password: yup.string().optional(),
  database: yup.string().required("Database name is required"),
})

interface DatabaseManagementProps {
  isOpen: boolean
  onClose: () => void
  clusters: Array<{ id: string, name: string }>
}

export function DatabaseManagement({ isOpen, onClose, clusters }: DatabaseManagementProps) {
  if (!isOpen) return null

  // Sample databases - replace with actual data
  const [databases, setDatabases] = useState<Database[]>([
    { 
      id: "db1", 
      name: "Main Database", 
      clusterId: "cluster-1", 
      clusterName: "Production Cluster",
      user: "admin",
      password: "********",
      database: "main"
    },
    { 
      id: "db2", 
      name: "Analytics DB", 
      clusterId: "cluster-2", 
      clusterName: "Analytics Cluster",
      user: "analyst",
      password: "********",
      database: "analytics"
    },
  ])
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [isAddingDatabase, setIsAddingDatabase] = useState(false)
  const [isConfiguring, setIsConfiguring] = useState<string | null>(null)
  const [editingDatabase, setEditingDatabase] = useState<Database | null>(null)

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
      clusterId: '',
      user: '',
      password: '',
      database: '',
    }
  })

  const handleDeleteDatabase = (id: string) => {
    setDatabases(prev => prev.filter(db => db.id !== id))
    setDeleteId(null)
  }

  const onSubmit = (data: { name: string, clusterId: string, user: string, password?: string, database: string }) => {
    const selectedCluster = clusters.find(c => c.id === data.clusterId)
    if (!selectedCluster) return

    const dbData: Omit<Database, 'id'> = {
      name: data.name,
      clusterId: data.clusterId,
      clusterName: selectedCluster.name,
      user: data.user,
      password: data.password || '',
      database: data.database
    }

    if (isAddingDatabase) {
      const id = `db${Date.now()}`
      setDatabases(prev => [...prev, { ...dbData, id }])
      setIsAddingDatabase(false)
    } else if (editingDatabase) {
      setDatabases(prev => prev.map(db => 
        db.id === editingDatabase.id ? { ...dbData, id: editingDatabase.id } : db
      ))
      setEditingDatabase(null)
      setIsConfiguring(null)
    }
    reset()
  }

  const handleConfigureDatabase = (db: Database) => {
    setEditingDatabase(db)
    setIsConfiguring(db.id)
    setValue('name', db.name)
    setValue('clusterId', db.clusterId)
    setValue('user', db.user)
    setValue('password', db.password || '')
    setValue('database', db.database)
  }

  const handleBackdropClick = useCallback((e: React.MouseEvent) => {
    if (e.target === e.currentTarget && e.button === 0) {
      onClose()
    }
  }, [onClose])

  const isEditing = isAddingDatabase || isConfiguring
  
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
              {isAddingDatabase ? 'Add New Database' : 
               isConfiguring ? 'Configure Database' : 
               'Database Management'}
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
              
              <div className="space-y-2">
                {databases.map((db) => (
                  <div 
                    key={db.id}
                    className="flex items-center justify-between p-3 border rounded-lg bg-muted/30 group hover:bg-muted/50"
                  >
                    <div className="space-y-1">
                      <h3 className="font-medium">{db.name}</h3>
                      <div className="text-sm text-muted-foreground">
                        <span>Cluster: {db.clusterName}</span>
                        <span className="mx-1">·</span>
                        <span>User: {db.user}</span>
                        <span className="mx-1">·</span>
                        <span>Database: {db.database}</span>
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
                          onClick={() => setDeleteId(db.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                        {deleteId === db.id && (
                          <ConfirmationPopup
                            message="Delete this database?"
                            onConfirm={() => handleDeleteDatabase(db.id)}
                            onCancel={() => setDeleteId(null)}
                          />
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {databases.length === 0 && (
                <div className="text-center py-8 text-muted-foreground">
                  No databases configured. Click "Add Database" to get started.
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
                  <Label htmlFor="clusterId">Cluster</Label>
                  <Select
                    value={getValues('clusterId')}
                    onValueChange={(value) => setValue('clusterId', value)}
                  >
                    <SelectTrigger className={errors.clusterId ? 'border-red-500' : ''}>
                      <SelectValue placeholder="Select a cluster" />
                    </SelectTrigger>
                    <SelectContent>
                      {clusters.map((cluster) => (
                        <SelectItem key={cluster.id} value={cluster.id}>
                          {cluster.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {errors.clusterId && (
                    <p className="text-sm text-red-500">{errors.clusterId.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="user">Username</Label>
                  <Input
                    id="user"
                    {...register('user')}
                    placeholder="e.g. admin"
                    className={errors.user ? 'border-red-500' : ''}
                  />
                  {errors.user && (
                    <p className="text-sm text-red-500">{errors.user.message}</p>
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
                  <Label htmlFor="database">Database</Label>
                  <Input
                    id="database"
                    {...register('database')}
                    placeholder="e.g. main"
                    className={errors.database ? 'border-red-500' : ''}
                  />
                  {errors.database && (
                    <p className="text-sm text-red-500">{errors.database.message}</p>
                  )}
                </div>
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