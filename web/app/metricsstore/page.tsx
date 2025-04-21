"use client"

import { useEffect, useState } from "react"
import { ApiError, DefaultService } from "@/api-gen"
import { MetricsStore } from "@/api-gen/models/MetricsStore"
import { MetricsStoreImport } from "@/api-gen/models/MetricsStoreImport"
import { MetricsStoreSpec } from "@/api-gen/models/MetricsStoreSpec"
import toast from "react-hot-toast"
import { Button } from "@/components/ui/button"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { DynamicFormData, STORE_TYPE, initialFormState } from "./metricsstoreform-hook"
import { MetricsStoreForm } from "./metricsstore-form"

export default function MetricsStorePage() {
  const [metricsStores, setMetricsStores] = useState<MetricsStore[]>([])
  const [loading, setLoading] = useState(true)
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editDialogOpen, setEditDialogOpen] = useState(false)
  const [currentMetricsStore, setCurrentMetricsStore] = useState<MetricsStore | null>(null)
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [forceDelete, setForceDelete] = useState(false)
  
  // Form initial data for editing
  const [editFormData, setEditFormData] = useState<DynamicFormData>(initialFormState)

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
    
    console.log("Original store:", store);
    console.log("Default labels from store:", store.defaultLabels);
    
    const formData: DynamicFormData = {
      name: store.name,
      type: "",
      fields: {},
      defaultLabels: store.defaultLabels || []
    }
    
    if (store.spec?.prometheus) {
      formData.type = STORE_TYPE.PROMETHEUS;
      formData.fields = {
        prometheus: store.spec.prometheus
      };
    } else if (store.spec?.victoriametrics) {
      formData.type = STORE_TYPE.VICTORIA_METRICS;
      formData.fields = {
        victoriametrics: store.spec.victoriametrics
      };
    }
    
    console.log("Prepared form data:", formData);
    
    setEditFormData(formData)
    setEditDialogOpen(true)
  }

  const openDeleteConfirm = (store: MetricsStore) => {
    setCurrentMetricsStore(store)
    setForceDelete(false)
    setDeleteConfirmOpen(true)
  }

  const handleCreateSubmit = async (formData: DynamicFormData) => {
    // If formData equals initialFormState, we know the Cancel button was clicked
    if (formData === initialFormState) {
      setCreateDialogOpen(false)
      return
    }
    
    try {
      const spec: MetricsStoreSpec = {};
      
      // Use the fields directly from the form data
      if (formData.type === STORE_TYPE.PROMETHEUS && formData.fields.prometheus) {
        spec.prometheus = formData.fields.prometheus;
      } else if (formData.type === STORE_TYPE.VICTORIA_METRICS && formData.fields.victoriametrics) {
        spec.victoriametrics = formData.fields.victoriametrics;
      }

      // Ensure defaultLabels is an array even if empty
      const defaultLabels = Array.isArray(formData.defaultLabels) ? formData.defaultLabels : [];
      
      console.log("Create formData", formData);
      console.log("Create defaultLabels", defaultLabels);

      const createPayload: MetricsStoreImport = {
        name: formData.name,
        spec: spec,
        defaultLabels: defaultLabels
      }

      const newStore = await DefaultService.importMetricsStore(createPayload)
      setMetricsStores(prev => [...prev, newStore])
      setCreateDialogOpen(false)
      toast.success("Metrics store created successfully")
    } catch (error) {
      toast.error("Failed to create metrics store")
      console.error("Error creating metrics store:", error)
    }
  }

  const handleEditSubmit = async (formData: DynamicFormData) => {
    // If formData equals initialFormState, we know the Cancel button was clicked
    if (formData === initialFormState) {
      setEditDialogOpen(false)
      return
    }
    
    if (!currentMetricsStore) return
    
    // Ensure defaultLabels is an array even if empty
    const defaultLabels = Array.isArray(formData.defaultLabels) ? formData.defaultLabels : [];
    
    console.log("Edit formData", formData);
    console.log("Edit defaultLabels", defaultLabels);
    
    try {
      const spec: MetricsStoreSpec = {};
      
      // Use the fields directly from the form data
      if (formData.type === STORE_TYPE.PROMETHEUS && formData.fields.prometheus) {
        spec.prometheus = formData.fields.prometheus;
      } else if (formData.type === STORE_TYPE.VICTORIA_METRICS && formData.fields.victoriametrics) {
        spec.victoriametrics = formData.fields.victoriametrics;
      }

      const updatedStore: MetricsStore = {
        ...currentMetricsStore,
        name: formData.name,
        spec: spec,
        defaultLabels: defaultLabels
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

  const renderDefaultLabels = (store: MetricsStore): React.ReactNode => {
    if (!store.defaultLabels || store.defaultLabels.length === 0) {
      return <span className="text-muted-foreground text-xs">None</span>
    }

    // Only show up to 3 labels to save space
    const visibleLabels = store.defaultLabels.slice(0, 3)
    const extraCount = store.defaultLabels.length - 3

    return (
      <div className="flex flex-wrap gap-1">
        {visibleLabels.map((label, index) => (
          <span 
            key={index} 
            className="inline-flex items-center rounded-full border border-gray-200 px-2 py-0.5 text-xs"
            title={`${label.key} ${getOperatorDisplay(label.op)} ${label.value}`}
          >
            {label.key} {getOperatorDisplay(label.op)} {label.value}
          </span>
        ))}
        {extraCount > 0 && (
          <span className="inline-flex items-center rounded-full bg-gray-100 px-2 py-0.5 text-xs">
            +{extraCount} more
          </span>
        )}
      </div>
    )
  }

  const getOperatorDisplay = (op: string): string => {
    switch (op) {
      case "EQ": return "="
      case "NEQ": return "!="
      case "RE": return "=~"
      case "NRE": return "!~"
      default: return op
    }
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
        <Button onClick={openCreateDialog}>Import Metrics Store</Button>
      </div>

      {metricsStores.length === 0 ? (
        <div className="rounded-lg border border-dashed p-8 text-center">
          <h3 className="text-lg font-medium">No metrics stores found</h3>
          <p className="text-sm text-muted-foreground mt-2">
            Add a metrics store to start monitoring your clusters
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
                <TableHead>Default Labels</TableHead>
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
                  <TableCell>
                    {renderDefaultLabels(store)}
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
            <DialogTitle>Import Metrics Store</DialogTitle>
            <DialogDescription>
              Add a new metrics store to monitor your clusters.
            </DialogDescription>
          </DialogHeader>
          <MetricsStoreForm 
            onSubmit={handleCreateSubmit} 
            onCancel={() => setCreateDialogOpen(false)} 
          />
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
          <MetricsStoreForm 
            initialData={editFormData} 
            onSubmit={handleEditSubmit} 
            onCancel={() => setEditDialogOpen(false)}
            isEdit={true}
          />
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