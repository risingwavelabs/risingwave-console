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
import * as yup from "yup"
import { useForm } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import { Loader2 } from "lucide-react"

const schema = yup.object().shape({
  name: yup.string().required("Name is required"),
  host: yup.string().required("Host is required"),
  sqlPort: yup.number()
    .typeError("SQL Port must be a number")
    .required("SQL Port is required")
    .min(1, "Port must be greater than 0")
    .max(65535, "Port must be less than 65536"),
  metaNodePort: yup.number()
    .typeError("Meta Node Port must be a number")
    .required("Meta Node Port is required")
    .min(1, "Port must be greater than 0")
    .max(65535, "Port must be less than 65536")
    .test("different-port", "Meta Node Port must be different from SQL Port",
      function (value) {
        return value !== this.parent.sqlPort
      }
    ),
  user: yup.string().required("User is required"),
  password: yup.string().optional(),
  database: yup.string().required("Database name is required")
})

export interface ClusterFormData extends yup.InferType<typeof schema> { }

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

  const isControlled = controlledOpen !== undefined && setControlledOpen !== undefined
  const isOpen = isControlled ? controlledOpen : open
  const setIsOpen = isControlled ? setControlledOpen : setOpen

  const {
    register,
    handleSubmit,
    reset,
    getValues,
    watch,
    formState: { errors, isValid }
  } = useForm<ClusterFormData>({
    mode: "all",
    resolver: yupResolver(schema),
    defaultValues: defaultValues || {
      host: "localhost",
      sqlPort: 8080,
      metaNodePort: 9090,
      user: "admin"
    }
  })

  // Watch all form fields to hide success message on change
  useEffect(() => {
    const subscription = watch(() => {
      setTestSuccess(false)
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
    const values = getValues()
    setTesting(true)
    setTestSuccess(false)
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1500))
      setTestSuccess(true)
    } catch (error) {
      alert("Connection failed!")
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
                <Label htmlFor="metaNodePort">Meta Node Port</Label>
                <Input
                  id="metaNodePort"
                  type="number"
                  {...register("metaNodePort")}
                  className={errors.metaNodePort ? "border-red-500" : ""}
                />
                {errors.metaNodePort && (
                  <p className="text-sm text-red-500">{errors.metaNodePort.message}</p>
                )}
              </div>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="database">Database</Label>
              <Input
                id="database"
                placeholder="myapp_db"
                {...register("database")}
                className={errors.database ? "border-red-500" : ""}
              />
              {errors.database && (
                <p className="text-sm text-red-500">{errors.database.message}</p>
              )}
            </div>
            <div className="grid gap-2">
              <Label htmlFor="user">User</Label>
              <Input
                id="user"
                placeholder="admin"
                {...register("user")}
                className={errors.user ? "border-red-500" : ""}
              />
              {errors.user && (
                <p className="text-sm text-red-500">{errors.user.message}</p>
              )}
            </div>
            <div className="grid gap-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder=""
                {...register("password")}
                className={errors.password ? "border-red-500" : ""}
              />
              {errors.password && (
                <p className="text-sm text-red-500">{errors.password.message}</p>
              )}
            </div>
            {testSuccess && (
              <p className="text-sm text-green-500 font-medium">
                ✓ Successfully connected
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