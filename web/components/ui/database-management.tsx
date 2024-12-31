import { Button } from "@/components/ui/button"
import { X, Settings2, Trash2, Plus, ArrowLeft } from "lucide-react"
import { useCallback, useState } from "react"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import * as yup from "yup"
import { useForm } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import { ConfirmationPopup } from "@/components/ui/confirmation-popup"

interface Database {
  id: string
  name: string
  host: string
  port: number
  username: string
  password?: string
}

const databaseSchema = yup.object().shape({
  name: yup.string().required("Database name is required"),
  host: yup.string().required("Host is required"),
  port: yup.number()
    .typeError("Port must be a number")
    .required("Port is required")
    .min(1, "Port must be greater than 0")
    .max(65535, "Port must be less than 65536"),
  username: yup.string().required("Username is required"),
  password: yup.string().default(""),
})

interface DatabaseManagementProps {
  isOpen: boolean
  onClose: () => void
}

export function DatabaseManagement({ isOpen, onClose }: DatabaseManagementProps) {
  if (!isOpen) return null

  // Sample databases - replace with actual data
  const [databases, setDatabases] = useState<Database[]>([
    { id: "db1", name: "Main Database", host: "localhost", port: 5432, username: "admin", password: "********" },
    { id: "db2", name: "Analytics DB", host: "analytics.example.com", port: 5432, username: "analyst", password: "********" },
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
  } = useForm({
    resolver: yupResolver(databaseSchema),
    mode: "all",
    defaultValues: {
      name: '',
      host: '',
      port: 5432,
      username: '',
      password: '',
    }
  })

  const handleDeleteDatabase = (id: string) => {
    setDatabases(prev => prev.filter(db => db.id !== id))
    setDeleteId(null)
  }

  const onSubmit = (data: Omit<Database, 'id'>) => {
    if (isAddingDatabase) {
      const id = `db${Date.now()}`
      setDatabases(prev => [...prev, { ...data, id }])
      setIsAddingDatabase(false)
    } else if (editingDatabase) {
      setDatabases(prev => prev.map(db => 
        db.id === editingDatabase.id ? { ...data, id: editingDatabase.id } : db
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
    setValue('host', db.host)
    setValue('port', db.port)
    setValue('username', db.username)
    setValue('password', db.password || '')
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
                      <p className="text-sm text-muted-foreground">
                        {db.username}@{db.host}:{db.port}
                      </p>
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
                  <Label htmlFor="name">Database Name</Label>
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
                  <Label htmlFor="host">Host</Label>
                  <Input
                    id="host"
                    {...register('host')}
                    placeholder="e.g. localhost or db.example.com"
                    className={errors.host ? 'border-red-500' : ''}
                  />
                  {errors.host && (
                    <p className="text-sm text-red-500">{errors.host.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="port">Port</Label>
                  <Input
                    id="port"
                    type="number"
                    {...register('port')}
                    placeholder="5432"
                    className={errors.port ? 'border-red-500' : ''}
                  />
                  {errors.port && (
                    <p className="text-sm text-red-500">{errors.port.message}</p>
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