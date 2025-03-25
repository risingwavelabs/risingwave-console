import React from 'react';
import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import toast from "react-hot-toast";
import { 
  useMetricsStoreForm, 
  metricsStoreTypes, 
  STORE_TYPE,
  DynamicFormData,
  MetricsStoreType
} from './metricsstoreform-hook';
import { DefaultLabelsEditor } from './default-labels-editor';

interface MetricsStoreFormProps {
  initialData?: DynamicFormData;
  onSubmit: (data: DynamicFormData) => void;
  onCancel: () => void;
  isEdit?: boolean;
}

export function MetricsStoreForm({ 
  initialData, 
  onSubmit, 
  onCancel, 
  isEdit = false 
}: MetricsStoreFormProps) {
  const { 
    formData, 
    setName, 
    setType, 
    setField, 
    addDefaultLabel, 
    removeDefaultLabel, 
    resetForm 
  } = useMetricsStoreForm(initialData);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Debug log
    console.log("Submitting form:", JSON.stringify(formData, null, 2));
    console.log("Default labels:", formData.defaultLabels);
    
    // Validate form
    if (!formData.name || !formData.type) {
      toast.error("Please fill all required fields");
      return;
    }
    
    // Type-specific validation
    if (formData.type === STORE_TYPE.PROMETHEUS) {
      if (!formData.fields.prometheus?.endpoint) {
        toast.error("Please fill the Prometheus endpoint field");
        return;
      }
    } else if (formData.type === STORE_TYPE.VICTORIA_METRICS) {
      if (!formData.fields.victoriametrics?.endpoint) {
        toast.error("Please fill the Victoria Metrics endpoint field");
        return;
      }
    }
    
    onSubmit(formData);
  };

  const handleCancel = () => {
    resetForm();
    onCancel();
  };

  // Render type-specific fields based on the selected type
  const renderTypeFields = () => {
    if (!formData.type) return null;

    if (formData.type === STORE_TYPE.PROMETHEUS) {
      return (
        <div className="space-y-4 mt-4">
          <div className="grid w-full items-center gap-1.5">
            <Label htmlFor="prometheus-endpoint">Endpoint</Label>
            <Input
              id="prometheus-endpoint"
              value={formData.fields.prometheus?.endpoint || ''}
              onChange={(e) => setField('endpoint', e.target.value)}
              placeholder="http://prometheus:9090"
            />
          </div>
        </div>
      );
    }

    if (formData.type === STORE_TYPE.VICTORIA_METRICS) {
      return (
        <div className="space-y-4 mt-4">
          <div className="grid w-full items-center gap-1.5">
            <Label htmlFor="victoriametrics-endpoint">Endpoint</Label>
            <Input
              id="victoriametrics-endpoint"
              value={formData.fields.victoriametrics?.endpoint || ''}
              onChange={(e) => setField('endpoint', e.target.value)}
              placeholder="http://victoriametrics:8428"
            />
          </div>
        </div>
      );
    }

    return null;
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid w-full items-center gap-1.5">
        <Label htmlFor="name">Name</Label>
        <Input
          id="name"
          value={formData.name}
          onChange={(e) => setName(e.target.value)}
          placeholder="My Metrics Store"
        />
      </div>

      <div className="grid w-full items-center gap-1.5">
        <Label htmlFor="type">Type</Label>
        <Select 
          value={formData.type} 
          onValueChange={(value) => setType(value as MetricsStoreType)}
          disabled={isEdit}
        >
          <SelectTrigger id="type">
            <SelectValue placeholder="Select type" />
          </SelectTrigger>
          <SelectContent>
            {metricsStoreTypes.map((type) => (
              <SelectItem key={type.value} value={type.value}>
                {type.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {renderTypeFields()}

      <div className="border-t pt-4">
        <DefaultLabelsEditor
          labels={formData.defaultLabels}
          onAdd={addDefaultLabel}
          onRemove={removeDefaultLabel}
        />
      </div>

      <DialogFooter>
        <Button type="button" variant="secondary" onClick={handleCancel}>
          Cancel
        </Button>
        <Button type="submit">
          {isEdit ? 'Update' : 'Create'}
        </Button>
      </DialogFooter>
    </form>
  );
} 