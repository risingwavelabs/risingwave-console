import React, { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { MetricsStoreLabelMatcher } from '@/api-gen/models/MetricsStoreLabelMatcher';
import { PlusCircle, Trash2 } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent } from '@/components/ui/card';

interface DefaultLabelsEditorProps {
  labels: MetricsStoreLabelMatcher[];
  onAdd: (label: MetricsStoreLabelMatcher) => void;
  onRemove: (index: number) => void;
}

export function DefaultLabelsEditor({ 
  labels, 
  onAdd, 
  onRemove 
}: DefaultLabelsEditorProps) {
  const [newLabel, setNewLabel] = useState<MetricsStoreLabelMatcher>({
    op: MetricsStoreLabelMatcher.op.EQ,
    key: '',
    value: ''
  });

  const resetNewLabel = () => {
    setNewLabel({
      op: MetricsStoreLabelMatcher.op.EQ,
      key: '',
      value: ''
    });
  };

  const handleAddLabel = () => {
    if (newLabel.key && newLabel.value && newLabel.op) {
      // Create a properly formatted label object
      const labelToAdd: MetricsStoreLabelMatcher = {
        op: newLabel.op,
        key: newLabel.key,
        value: newLabel.value
      };
      
      onAdd(labelToAdd);
      resetNewLabel();
    } else {
      console.error("Cannot add label: missing required fields");
    }
  };

  const getOperatorDisplay = (op: MetricsStoreLabelMatcher.op): string => {
    switch (op) {
      case MetricsStoreLabelMatcher.op.EQ: return '=';
      case MetricsStoreLabelMatcher.op.NEQ: return '!=';
      case MetricsStoreLabelMatcher.op.RE: return '=~';
      case MetricsStoreLabelMatcher.op.NRE: return '!~';
      default: return op;
    }
  };

  return (
    <div className="space-y-4">
      <Label>Default Labels</Label>
      
      {/* Display existing labels */}
      {labels.length > 0 ? (
        <div className="space-y-2">
          {labels.map((label, index) => (
            <div key={index} className="flex items-center gap-2">
              <Badge variant="outline" className="px-2 py-1 text-xs">
                {label.key} {getOperatorDisplay(label.op)} {label.value}
              </Badge>
              <Button 
                type="button" 
                variant="ghost" 
                size="sm" 
                onClick={() => onRemove(index)}
                className="h-6 w-6 p-0"
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>
          ))}
        </div>
      ) : (
        <p className="text-sm text-muted-foreground">No default labels set.</p>
      )}

      {/* Add new label form */}
      <Card>
        <CardContent className="pt-4">
          <div className="grid grid-cols-12 gap-2 items-end">
            <div className="col-span-4">
              <Label htmlFor="label-key" className="text-xs">Key</Label>
              <Input
                id="label-key"
                placeholder="e.g. namespace"
                value={newLabel.key}
                onChange={(e) => setNewLabel({ ...newLabel, key: e.target.value })}
                className="h-8"
              />
            </div>
            
            <div className="col-span-3">
              <Label htmlFor="label-op" className="text-xs">Operator</Label>
              <Select 
                value={newLabel.op} 
                onValueChange={(value) => setNewLabel({ ...newLabel, op: value as MetricsStoreLabelMatcher.op })}
              >
                <SelectTrigger id="label-op" className="h-8">
                  <SelectValue placeholder="=" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value={MetricsStoreLabelMatcher.op.EQ}>=</SelectItem>
                  <SelectItem value={MetricsStoreLabelMatcher.op.NEQ}>!=</SelectItem>
                  <SelectItem value={MetricsStoreLabelMatcher.op.RE}>=~</SelectItem>
                  <SelectItem value={MetricsStoreLabelMatcher.op.NRE}>!~</SelectItem>
                </SelectContent>
              </Select>
            </div>
            
            <div className="col-span-4">
              <Label htmlFor="label-value" className="text-xs">Value</Label>
              <Input
                id="label-value"
                placeholder="e.g. risingwave-console"
                value={newLabel.value}
                onChange={(e) => setNewLabel({ ...newLabel, value: e.target.value })}
                className="h-8"
              />
            </div>
            
            <div className="col-span-1">
              <Button 
                type="button" 
                variant="ghost" 
                size="sm" 
                onClick={handleAddLabel}
                disabled={!newLabel.key || !newLabel.value}
                className="h-8 w-8 p-0"
              >
                <PlusCircle className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
} 