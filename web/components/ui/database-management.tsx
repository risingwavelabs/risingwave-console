import { Button } from "@/components/ui/button"
import { X, Settings2, Trash2, Plus, ArrowLeft } from "lucide-react"
import { useCallback, useState } from "react"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

interface Database {
  id: string
  name: string
  host: string
  port: number
  username: string
}

interface DatabaseManagementProps {
  isOpen: boolean
  onClose: () => void
}

export function DatabaseManagement({ isOpen, onClose }: DatabaseManagementProps) {
  if (!isOpen) return null

  // Sample databases - replace with actual data
  const [databases, setDatabases] = useState<Database[]>([
    { id: "db1", name: "Main Database", host: "localhost", port: 5432, username: "admin" },
    { id: "db2", name: "Analytics DB", host: "analytics.example.com", port: 5432, username: "analyst" },
  ])
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [isAddingDatabase, setIsAddingDatabase] = useState(false)
  const [isConfiguring, setIsConfiguring] = useState<string | null>(null)
  const [newDatabase, setNewDatabase] = useState<Omit<Database, 'id'>>({
    name: '',
    host: '',
    port: 5432,
    username: '',
  })
  const [editingDatabase, setEditingDatabase] = useState<Database | null>(null)

  const handleDeleteDatabase = (id: string) => {
    setDatabases(prev => prev.filter(db => db.id !== id))
    setDeleteId(null)
  }

  const handleAddDatabase = () => {
    const id = `db${Date.now()}`
    setDatabases(prev => [...prev, { ...newDatabase, id }])
    setNewDatabase({ name: '', host: '', port: 5432, username: '' })
    setIsAddingDatabase(false)
  }

  const handleConfigureDatabase = (db: Database) => {
    setEditingDatabase(db)
    setIsConfiguring(db.id)
  }

  const handleSaveConfiguration = () => {
    if (!editingDatabase) return
    setDatabases(prev => prev.map(db => 
      db.id === editingDatabase.id ? editingDatabase : db
    ))
    setEditingDatabase(null)
    setIsConfiguring(null)
  }

  const handleBackdropClick = useCallback((e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose()
    }
  }, [onClose])

  const isEditing = isAddingDatabase || isConfiguring
  const currentDatabase = isAddingDatabase ? newDatabase : editingDatabase

  return (
    <div 
      className="fixed inset-0 bg-black/50 z-50"
      onClick={handleBackdropClick}
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
                {databases.map(db => (
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
                          <div className="absolute right-0 top-full mt-1 bg-background border rounded-lg shadow-lg py-2 px-3 z-50 whitespace-nowrap">
                            <p className="text-xs text-muted-foreground mb-2">Delete this database?</p>
                            <div className="flex gap-2">
                              <Button 
                                size="sm" 
                                variant="ghost"
                                className="h-7 px-2 text-xs"
                                onClick={() => setDeleteId(null)}
                              >
                                Cancel
                              </Button>
                              <Button 
                                size="sm"
                                variant="destructive"
                                className="h-7 px-2 text-xs"
                                onClick={() => handleDeleteDatabase(db.id)}
                              >
                                Delete
                              </Button>
                            </div>
                          </div>
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
            <div className="space-y-4">
              <div className="grid gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Database Name</Label>
                  <Input
                    id="name"
                    value={currentDatabase?.name}
                    onChange={e => {
                      if (isAddingDatabase) {
                        setNewDatabase(prev => ({ ...prev, name: e.target.value }))
                      } else {
                        setEditingDatabase(prev => prev ? { ...prev, name: e.target.value } : null)
                      }
                    }}
                    placeholder="e.g. Production Database"
                  />
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="host">Host</Label>
                  <Input
                    id="host"
                    value={currentDatabase?.host}
                    onChange={e => {
                      if (isAddingDatabase) {
                        setNewDatabase(prev => ({ ...prev, host: e.target.value }))
                      } else {
                        setEditingDatabase(prev => prev ? { ...prev, host: e.target.value } : null)
                      }
                    }}
                    placeholder="e.g. localhost or db.example.com"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="port">Port</Label>
                  <Input
                    id="port"
                    type="number"
                    value={currentDatabase?.port}
                    onChange={e => {
                      const port = parseInt(e.target.value) || 5432
                      if (isAddingDatabase) {
                        setNewDatabase(prev => ({ ...prev, port }))
                      } else {
                        setEditingDatabase(prev => prev ? { ...prev, port } : null)
                      }
                    }}
                    placeholder="5432"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="username">Username</Label>
                  <Input
                    id="username"
                    value={currentDatabase?.username}
                    onChange={e => {
                      if (isAddingDatabase) {
                        setNewDatabase(prev => ({ ...prev, username: e.target.value }))
                      } else {
                        setEditingDatabase(prev => prev ? { ...prev, username: e.target.value } : null)
                      }
                    }}
                    placeholder="e.g. admin"
                  />
                </div>
              </div>
            </div>
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
                }}
              >
                Cancel
              </Button>
              <Button 
                onClick={isAddingDatabase ? handleAddDatabase : handleSaveConfiguration}
                disabled={!currentDatabase?.name || !currentDatabase?.host || !currentDatabase?.username}
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