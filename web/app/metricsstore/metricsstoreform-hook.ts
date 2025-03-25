import { useState, useEffect } from "react";
import { MetricsStoreSpec } from "@/api-gen/models/MetricsStoreSpec";
import { MetricsStorePrometheus } from "@/api-gen/models/MetricsStorePrometheus";
import { MetricsStoreVictoriaMetrics } from "@/api-gen/models/MetricsStoreVictoriaMetrics";
import { MetricsStoreLabelMatcherList } from "@/api-gen/models/MetricsStoreLabelMatcherList";
import { MetricsStoreLabelMatcher } from "@/api-gen/models/MetricsStoreLabelMatcher";

// Define metrics store types based on generated API types
export type MetricsStoreType = keyof Omit<MetricsStoreSpec, "extends" | "implements"> | "";

// Define store type constants for easier access
export enum STORE_TYPE {
  PROMETHEUS = "prometheus",
  VICTORIA_METRICS = "victoriametrics"
}

// Available metrics store types with display names
export const metricsStoreTypes = [
  { value: STORE_TYPE.PROMETHEUS, label: "Prometheus" },
  { value: STORE_TYPE.VICTORIA_METRICS, label: "Victoria Metrics" }
];

// Dynamic form data with typed fields
export interface DynamicFormData {
  name: string;
  type: MetricsStoreType;
  fields: {
    prometheus?: MetricsStorePrometheus;
    victoriametrics?: MetricsStoreVictoriaMetrics;
  };
  defaultLabels: MetricsStoreLabelMatcherList;
}

export const initialFormState: DynamicFormData = {
  name: "",
  type: "",
  fields: {},
  defaultLabels: []
}

// Simple hook for form state management
export function useMetricsStoreForm(initialData: DynamicFormData = initialFormState) {
  const [formData, setFormData] = useState<DynamicFormData>(initialData);
  
  // Reset form if initialData changes
  useEffect(() => {
    setFormData(initialData);
  }, [initialData]);
  
  // Update name field
  const setName = (name: string) => {
    setFormData(prev => ({ ...prev, name }));
  };
  
  // Update type field and reset type-specific fields
  const setType = (type: MetricsStoreType) => {
    setFormData(prev => ({ ...prev, type, fields: {} }));
  };
  
  // Update fields based on the selected type
  const setField = (fieldName: string, value: string) => {
    if (!formData.type) return;
    
    setFormData(prev => {
      const newFields = { ...prev.fields };
      
      if (formData.type === STORE_TYPE.PROMETHEUS) {
        newFields.prometheus = {
          ...(newFields.prometheus || { endpoint: '' }),
          [fieldName]: value
        } as MetricsStorePrometheus;
      } else if (formData.type === STORE_TYPE.VICTORIA_METRICS) {
        newFields.victoriametrics = {
          ...(newFields.victoriametrics || { endpoint: '' }),
          [fieldName]: value
        } as MetricsStoreVictoriaMetrics;
      }
      
      return {
        ...prev,
        fields: newFields
      };
    });
  };

  // Add a default label
  const addDefaultLabel = (label: MetricsStoreLabelMatcher) => {
    // Ensure the label has all required properties
    if (!label.op || !label.key || !label.value) {
      console.error("Invalid label format. All properties (op, key, value) are required.");
      return;
    }
    
    // Clone the label to avoid reference issues
    const newLabel: MetricsStoreLabelMatcher = {
      op: label.op,
      key: label.key,
      value: label.value
    };
    
    setFormData(prev => ({
      ...prev,
      defaultLabels: [...prev.defaultLabels, newLabel]
    }));
  };

  // Update a default label
  const updateDefaultLabel = (index: number, label: MetricsStoreLabelMatcher) => {
    setFormData(prev => {
      const newLabels = [...prev.defaultLabels];
      newLabels[index] = label;
      return {
        ...prev,
        defaultLabels: newLabels
      };
    });
  };

  // Remove a default label
  const removeDefaultLabel = (index: number) => {
    setFormData(prev => ({
      ...prev,
      defaultLabels: prev.defaultLabels.filter((_, i) => i !== index)
    }));
  };
  
  // Reset form to initial state
  const resetForm = () => {
    setFormData(initialFormState);
  };
  
  return {
    formData,
    setName,
    setType,
    setField,
    addDefaultLabel,
    updateDefaultLabel,
    removeDefaultLabel,
    resetForm
  };
} 