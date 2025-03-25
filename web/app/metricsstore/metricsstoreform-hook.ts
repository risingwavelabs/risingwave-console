import { useState, useEffect } from "react";
import { MetricsStoreSpec } from "@/api-gen/models/MetricsStoreSpec";
import { MetricsStorePrometheus } from "@/api-gen/models/MetricsStorePrometheus";
import { MetricsStoreVictoriaMetrics } from "@/api-gen/models/MetricsStoreVictoriaMetrics";

// Define metrics store types based on generated API types
export type MetricsStoreType = keyof Omit<MetricsStoreSpec, "extends" | "implements"> | "";

// Type definitions for the available store types
export const STORE_TYPE = {
  PROMETHEUS: "prometheus" as const,
  VICTORIA_METRICS: "victoriametrics" as const
};

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
}

export const initialFormState: DynamicFormData = {
  name: "",
  type: "",
  fields: {}
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
  
  // Reset form to initial state
  const resetForm = () => {
    setFormData(initialFormState);
  };
  
  return {
    formData,
    setName,
    setType,
    setField,
    resetForm
  };
} 