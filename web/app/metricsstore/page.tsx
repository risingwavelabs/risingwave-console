"use client"

import { useEffect, useState } from "react"
import { ApiError, DefaultService } from "@/api-gen"
import { MetricsStore } from "@/api-gen/models/MetricsStore"
import { MetricsStoreCreate } from "@/api-gen/models/MetricsStoreCreate"
import { MetricsStoreSpec } from "@/api-gen/models/MetricsStoreSpec"
import { MetricsStorePrometheus } from "@/api-gen/models/MetricsStorePrometheus"
import { MetricsStoreVictoriaMetrics } from "@/api-gen/models/MetricsStoreVictoriaMetrics"
import toast from "react-hot-toast"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

type FormData = {
  name: string;
  type: "prometheus" | "victoriametrics" | "";
  endpoint: string;
}

const initialFormState: FormData = {
  name: "",
  type: "",
  endpoint: ""
}

// Separate form component that manages its own state
function MetricsStoreForm({ 
  initialData = initialFormState, 
  onSubmit 
}: { 
  initialData?: FormData; 
  onSubmit: (data: FormData) => void; 
}) {
  const [formData, setFormData] = useState<FormData>(initialData);

  // Reset form if initialData changes (when editing a different metrics store)
  useEffect(() => {
    setFormData(initialData);
  }, [initialData]);

  const handleChange = (field: keyof FormData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleSubmit = () => {
    if (!formData.name || !formData.type || !formData.endpoint) {
      toast.error("Please fill all required fields");
      return;
    }
    onSubmit(formData);
  };

  return (
    <>
      <div className="grid gap-4 py-4">
        <div className="grid grid-cols-4 items-center gap-4">
          <Label htmlFor="name" className="text-right">Name</Label>
          <Input 
            id="name" 
            value={formData.name} 
            onChange={(e) => handleChange('name', e.target.value)} 
            className="col-span-3" 
          />
        </div>
        <div className="grid grid-cols-4 items-center gap-4">
          <Label htmlFor="type" className="text-right">Type</Label>
          <Select
            value={formData.type}
            onValueChange={(value: "prometheus" | "victoriametrics") => 
              handleChange('type', value)
            }
          >
            <SelectTrigger className="col-span-3">
              <SelectValue placeholder="Select type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="prometheus">Prometheus</SelectItem>
              <SelectItem value="victoriametrics">Victoria Metrics</SelectItem>
            </SelectContent>
          </Select>
        </div>
        {formData.type && (
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="endpoint" className="text-right">Endpoint</Label>
            <Input 
              id="endpoint" 
              value={formData.endpoint} 
              onChange={(e) => handleChange('endpoint', e.target.value)} 
              className="col-span-3" 
              placeholder={`${formData.type} endpoint URL`}
            />
          </div>
        )}
      </div>
      <DialogFooter>
        <Button variant="outline" type="button" onClick={() => onSubmit(initialFormState)}>
          Cancel
        </Button>
        <Button type="button" onClick={handleSubmit}>
          {initialData === initialFormState ? 'Create' : 'Save Changes'}
        </Button>
      </DialogFooter>
    </>
  );
}

export default function MetricsStorePage() {
  const [metricsStores, setMetricsStores] = useState<MetricsStore[]>([])
  const [loading, setLoading] = useState(true)
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editDialogOpen, setEditDialogOpen] = useState(false)
  const [currentMetricsStore, setCurrentMetricsStore] = useState<MetricsStore | null>(null)
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [forceDelete, setForceDelete] = useState(false)
  
  // Form initial data for editing
  const [editFormData, setEditFormData] = useState<FormData>(initialFormState)

  useEffect(() => {
    fetchMetricsStores()
  }, [])

  const fetchMetricsStores = async () => {
    try {
      const data = await DefaultService.listMetricsStores()
      setMetricsStores(data)
    } catch (error) {
      toast.error("Failed to load metrics stores")
      console.error("Error loading metrics stores:", error)
    } finally {
      setLoading(false)
    }
  }

  const openCreateDialog = () => {
    setCreateDialogOpen(true)
  }

  const openEditDialog = (store: MetricsStore) => {
    setCurrentMetricsStore(store)
    
    const formData: FormData = {
      name: store.name,
      type: "",
      endpoint: ""
    }
    
    if (store.spec?.prometheus) {
      formData.type = "prometheus"
      formData.endpoint = store.spec.prometheus.endpoint
    } else if (store.spec?.victoriametrics) {
      formData.type = "victoriametrics"
      formData.endpoint = store.spec.victoriametrics.endpoint
    }
    
    setEditFormData(formData)
    setEditDialogOpen(true)
  }

  const openDeleteConfirm = (store: MetricsStore) => {
    setCurrentMetricsStore(store)
    setForceDelete(false)
    setDeleteConfirmOpen(true)
  }

  const handleCreateSubmit = async (formData: FormData) => {
    // If formData equals initialFormState, we know the Cancel button was clicked
    if (formData === initialFormState) {
      setCreateDialogOpen(false)
      return
    }
    
    try {
      const spec: MetricsStoreSpec = {}
      
      if (formData.type === "prometheus") {
        spec.prometheus = { endpoint: formData.endpoint }
      } else if (formData.type === "victoriametrics") {
        spec.victoriametrics = { endpoint: formData.endpoint }
      }

      const createPayload: MetricsStoreCreate = {
        name: formData.name,
        spec: spec
      }

      const newStore = await DefaultService.createMetricsStore(createPayload)
      setMetricsStores(prev => [...prev, newStore])
      setCreateDialogOpen(false)
      toast.success("Metrics store created successfully")
    } catch (error) {
      toast.error("Failed to create metrics store")
      console.error("Error creating metrics store:", error)
    }
  }

  const handleEditSubmit = async (formData: FormData) => {
    // If formData equals initialFormState, we know the Cancel button was clicked
    if (formData === initialFormState) {
      setEditDialogOpen(false)
      return
    }
    
    if (!currentMetricsStore) return
    
    try {
      const spec: MetricsStoreSpec = {}
      
      if (formData.type === "prometheus") {
        spec.prometheus = { endpoint: formData.endpoint }
      } else if (formData.type === "victoriametrics") {
        spec.victoriametrics = { endpoint: formData.endpoint }
      }

      const updatedStore: MetricsStore = {
        ...currentMetricsStore,
        name: formData.name,
        spec: spec
      }

      await DefaultService.updateMetricsStore(currentMetricsStore.ID, updatedStore)
      setMetricsStores(prev => 
        prev.map(store => store.ID === currentMetricsStore.ID ? updatedStore : store)
      )
      setEditDialogOpen(false)
      toast.success("Metrics store updated successfully")
    } catch (error) {
      toast.error("Failed to update metrics store")
      console.error("Error updating metrics store:", error)
    }
  }

  const handleDeleteMetricsStore = async () => {
    if (!currentMetricsStore) return

    try {
      await DefaultService.deleteMetricsStore(currentMetricsStore.ID, forceDelete)
      setMetricsStores(prev => prev.filter(store => store.ID !== currentMetricsStore.ID))
      setDeleteConfirmOpen(false)
      toast.success("Metrics store deleted successfully")
    } catch (error) {
      if (error instanceof ApiError && error.status === 409) {
        toast.error("This metrics store is in use. Use force delete to remove it anyway.")
        setForceDelete(true)
      } else {
        toast.error("Failed to delete metrics store")
        console.error("Error deleting metrics store:", error)
        setDeleteConfirmOpen(false)
      }
    }
  }

  const renderStoreType = (store: MetricsStore): string => {
    if (store.spec?.prometheus) return "Prometheus"
    if (store.spec?.victoriametrics) return "Victoria Metrics"
    return "Unknown"
  }

  const renderEndpoint = (store: MetricsStore): string => {
    if (store.spec?.prometheus) return store.spec.prometheus.endpoint
    if (store.spec?.victoriametrics) return store.spec.victoriametrics.endpoint
    return ""
  }

  if (loading) {
    return (
      <div className="p-8">
        <div className="flex items-center justify-center h-[200px]">
          <div className="text-muted-foreground">Loading metrics stores...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-semibold">Metrics Stores</h1>
          <p className="text-sm text-muted-foreground">
            Manage your metrics stores for monitoring
          </p>
        </div>
        <Button onClick={openCreateDialog}>Create Metrics Store</Button>
      </div>

      {metricsStores.length === 0 ? (
        <div className="rounded-lg border border-dashed p-8 text-center">
          <h3 className="text-lg font-medium">No metrics stores found</h3>
          <p className="text-sm text-muted-foreground mt-2">
            Create a metrics store to start monitoring your clusters
          </p>
        </div>
      ) : (
        <div className="border rounded-md">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Endpoint</TableHead>
                <TableHead>Created</TableHead>
                <TableHead className="w-[100px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {metricsStores.map((store) => (
                <TableRow key={store.ID}>
                  <TableCell className="font-medium">{store.name}</TableCell>
                  <TableCell>{renderStoreType(store)}</TableCell>
                  <TableCell className="max-w-[300px] truncate" title={renderEndpoint(store)}>
                    {renderEndpoint(store)}
                  </TableCell>
                  <TableCell>{new Date(store.createdAt).toLocaleDateString()}</TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end space-x-2">
                      <Button 
                        variant="outline" 
                        size="sm" 
                        onClick={() => openEditDialog(store)}
                      >
                        Edit
                      </Button>
                      <Button 
                        variant="destructive" 
                        size="sm" 
                        onClick={() => openDeleteConfirm(store)}
                      >
                        Delete
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}

      {/* Create Dialog */}
      <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create Metrics Store</DialogTitle>
            <DialogDescription>
              Add a new metrics store to monitor your clusters.
            </DialogDescription>
          </DialogHeader>
          <MetricsStoreForm onSubmit={handleCreateSubmit} />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Metrics Store</DialogTitle>
            <DialogDescription>
              Make changes to your metrics store.
            </DialogDescription>
          </DialogHeader>
          <MetricsStoreForm initialData={editFormData} onSubmit={handleEditSubmit} />
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Metrics Store</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this metrics store? 
              {forceDelete && (
                <div className="mt-2 text-orange-500">
                  Warning: This metrics store is in use. Deleting it may affect monitoring of your clusters.
                </div>
              )}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteConfirmOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDeleteMetricsStore}>
              {forceDelete ? "Force Delete" : "Delete"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
