"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip"
import * as yup from "yup"
import { useForm, Controller } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import { Loader2, HelpCircle } from "lucide-react"
import { DefaultService } from "@/api-gen"
import { toast } from "react-hot-toast"

const schema = yup.object().shape({
  name: yup.string().required("Name is required"),
  host: yup.string().required("Host is required"),
  sqlPort: yup.number()
    .typeError("SQL Port must be a number")
    .required("SQL Port is required")
    .min(1, "Port must be greater than 0")
    .max(65535, "Port must be less than 65536"),
  metaPort: yup.number()
    .typeError("Meta Node Port must be a number")
    .required("Meta Node Port is required")
    .min(1, "Port must be greater than 0")
    .max(65535, "Port must be less than 65536"),
  httpPort: yup.number()
    .typeError("HTTP Port must be a number")
    .required("HTTP Port is required")
    .min(1, "Port must be greater than 0")
    .max(65535, "Port must be less than 65536"),
  version: yup.string().required("Version is required"),
  prometheusEndpoint: yup.string()
})

export type ClusterFormData = yup.InferType<typeof schema>

interface ClusterDialogProps {
  mode: 'create' | 'edit'
  defaultValues?: ClusterFormData
  trigger?: React.ReactNode
  open?: boolean
  onOpenChange?: (open: boolean) => void
  onSubmit: (data: ClusterFormData) => void
}

export function ClusterDialog({ 
  mode = 'create',
  defaultValues,
  trigger,
  open: controlledOpen,
  onOpenChange: setControlledOpen,
  onSubmit 
}: ClusterDialogProps) {
  const [open, setOpen] = useState(false)
  const [testing, setTesting] = useState(false)
  const [testSuccess, setTestSuccess] = useState(false)
  const [testMessage, setTestMessage] = useState<string>("")
  const [versions, setVersions] = useState<string[]>([])

  const isControlled = controlledOpen !== undefined && setControlledOpen !== undefined
  const isOpen = isControlled ? controlledOpen : open
  const setIsOpen = isControlled ? setControlledOpen : setOpen

  const {
    register,
    handleSubmit,
    reset,
    watch,
    control,
    formState: { errors, isValid }
  } = useForm<ClusterFormData>({
    mode: "all",
    resolver: yupResolver(schema),
    defaultValues: defaultValues || {
      host: "",
      sqlPort: 4566,
      metaPort: 5690,
      httpPort: 5691,
      version: "latest",
      prometheusEndpoint: undefined
    }
  })

  useEffect(() => {
    const fetchVersions = async () => {
      try {
        const data = await DefaultService.listClusterVersions()
        setVersions(data)
      } catch (error) {
        console.error("Error loading versions:", error)
      }
    }
    void fetchVersions()
  }, [])

  // Watch all form fields to hide success message on change
  useEffect(() => {
    const subscription = watch(() => {
      setTestSuccess(false)
      setTestMessage("")
    })
    return () => subscription.unsubscribe()
  }, [watch])

  const handleFormSubmit = handleSubmit((data) => {
    onSubmit(data)
    setIsOpen(false)
    if (!defaultValues) {
      reset()
    }
  })

  const handleTestConnection = async () => {
    setTesting(true)
    setTestSuccess(false)
    setTestMessage("")
    try {
      const formData = watch()
      const result = await DefaultService.testClusterConnection({
        host: formData.host,
        sqlPort: Number(formData.sqlPort),
        metaPort: Number(formData.metaPort),
        httpPort: Number(formData.httpPort)
      })
      setTestSuccess(result.success)
      setTestMessage(result.result)
    } catch (error) {
      toast.error(`Failed to test connection: ${error instanceof Error ? error.message : "Unknown error"}`)
    } finally {
      setTesting(false)
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        {trigger || <Button>{mode === 'create' ? 'New Cluster' : 'Edit'}</Button>}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleFormSubmit}>
          <DialogHeader>
            <DialogTitle>{mode === 'create' ? 'Create New Cluster' : 'Edit Cluster'}</DialogTitle>
            <DialogDescription>
              {mode === 'create' 
                ? 'Enter the details for your new database cluster.'
                : 'Modify the cluster connection details.'
              }
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                placeholder="Production DB"
                {...register("name")}
                className={errors.name ? "border-red-500" : ""}
              />
              {errors.name && (
                <p className="text-sm text-red-500">{errors.name.message}</p>
              )}
            </div>
            <div className="grid gap-2">
              <Label htmlFor="host">Host</Label>
              <Input
                id="host"
                placeholder="localhost"
                {...register("host")}
                className={errors.host ? "border-red-500" : ""}
              />
              {errors.host && (
                <p className="text-sm text-red-500">{errors.host.message}</p>
              )}
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="sqlPort">SQL Port</Label>
                <Input
                  id="sqlPort"
                  type="number"
                  {...register("sqlPort")}
                  className={errors.sqlPort ? "border-red-500" : ""}
                />
                {errors.sqlPort && (
                  <p className="text-sm text-red-500">{errors.sqlPort.message}</p>
                )}
              </div>
              <div className="grid gap-2">
                <Label htmlFor="metaPort">Meta Node Port</Label>
                <Input
                  id="metaPort"
                  type="number"
                  {...register("metaPort")}
                  className={errors.metaPort ? "border-red-500" : ""}
                />
                {errors.metaPort && (
                  <p className="text-sm text-red-500">{errors.metaPort.message}</p>
                )}
              </div>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="httpPort">HTTP Port</Label>
              <Input
                id="httpPort"
                type="number"
                {...register("httpPort")}
                className={errors.httpPort ? "border-red-500" : ""}
              />
              {errors.httpPort && (
                <p className="text-sm text-red-500">{errors.httpPort.message}</p>
              )}
            </div>
            <div className="grid gap-2">
              <Label htmlFor="prometheusEndpoint">Prometheus Endpoint</Label>
              <Input
                id="prometheusEndpoint"
                placeholder="http://localhost:9090"
                {...register("prometheusEndpoint")}
                className={errors.prometheusEndpoint ? "border-red-500" : ""}
              />
              {errors.prometheusEndpoint && (
                <p className="text-sm text-red-500">{errors.prometheusEndpoint.message}</p>
              )}
            </div>
            <div className="grid gap-2">
              <div className="flex items-center gap-2">
                <Label htmlFor="version">Version</Label>
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger>
                      <HelpCircle className="h-4 w-4 text-muted-foreground" />
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>The version determines which risectl release will be used for meta node operations</p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </div>
              <Controller
                name="version"
                control={control}
                render={({ field }) => (
                  <Select
                    defaultValue={defaultValues?.version || "latest"}
                    value={field.value}
                    onValueChange={field.onChange}
                  >
                    <SelectTrigger className={errors.version ? "border-red-500" : ""}>
                      <SelectValue placeholder="Select version" />
                    </SelectTrigger>
                    <SelectContent className="max-h-[200px]">
                      {versions.map((version) => (
                        <SelectItem key={version} value={version}>
                          {version}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.version && (
                <p className="text-sm text-red-500">{errors.version.message}</p>
              )}
            </div>
            {testMessage && (
              <p className={`text-sm font-medium ${testSuccess ? "text-green-500" : "text-red-500"}`}>
                {testMessage.split("\n").map((line, index) => (
                  <span key={index}>
                    {line}
                    <br />
                  </span>
                ))}
              </p>
            )}
          </div>
          <DialogFooter className="gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={handleTestConnection}
              disabled={!isValid || testing}
            >
              {testing && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Test Connection
            </Button>
            <Button type="submit" disabled={!isValid}>
              {mode === 'create' ? 'Create Cluster' : 'Save Changes'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
} 