import React, { useState, useEffect } from 'react';
import { X, ChevronRight, ChevronDown } from 'lucide-react';
import { Button } from './button';
import Editor from '@monaco-editor/react';

interface ProgressItem {
  name: string;
  status: 'running' | 'completed' | 'failed';
  progress: number;
  startTime: number; // Unix timestamp in milliseconds
  sql: string; // SQL code that's being executed
}

interface ProgressViewProps {
  items?: ProgressItem[];
  onCancel?: (name: string) => void;
}

export function ProgressView({ items = [], onCancel }: ProgressViewProps) {
  const [cancelId, setCancelId] = useState<string | null>(null);
  const [now, setNow] = useState(Date.now());
  const [expandedItems, setExpandedItems] = useState<Set<string>>(new Set());

  // Update time every second for running items
  useEffect(() => {
    const hasRunningItems = (items.length > 0 ? items : sampleItems).some(
      item => item.status === 'running'
    );

    if (hasRunningItems) {
      const interval = setInterval(() => {
        setNow(Date.now());
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [items]);

  // Sample data for demonstration
  const sampleItems: ProgressItem[] = [
    {
      name: 'user_metrics',
      status: 'completed',
      progress: 100,
      startTime: Date.now() - 3600000, // 1 hour ago
      sql: `CREATE MATERIALIZED VIEW user_metrics AS
SELECT user_id,
       COUNT(*) as event_count,
       MAX(timestamp) as last_seen
FROM user_events
GROUP BY user_id;`
    },
    {
      name: 'product_analytics',
      status: 'running',
      progress: 65,
      startTime: Date.now() - 300000, // 5 minutes ago
      sql: `CREATE MATERIALIZED VIEW product_analytics AS
SELECT product_id,
       COUNT(CASE WHEN event_type = 'view' THEN 1 END) as view_count,
       COUNT(CASE WHEN event_type = 'purchase' THEN 1 END) as purchase_count
FROM product_events
GROUP BY product_id;`
    },
    {
      name: 'user_product_recommendations',
      status: 'running',
      progress: 30,
      startTime: Date.now() - 120000, // 2 minutes ago
      sql: `CREATE MATERIALIZED VIEW user_product_recommendations AS
SELECT u.user_id,
       p.product_id,
       COUNT(*) as interaction_count,
       AVG(CASE WHEN p.event_type = 'purchase' THEN 1 ELSE 0.1 END) as score
FROM user_events u
JOIN product_events p ON u.session_id = p.session_id
GROUP BY u.user_id, p.product_id;`
    }
  ];

  const displayItems = items.length > 0 ? items : sampleItems;

  const toggleExpand = (name: string) => {
    setExpandedItems(prev => {
      const newSet = new Set(prev);
      if (newSet.has(name)) {
        newSet.delete(name);
      } else {
        newSet.add(name);
      }
      return newSet;
    });
  };

  const getStatusColor = (status: ProgressItem['status']) => {
    switch (status) {
      case 'completed':
        return 'bg-green-500';
      case 'failed':
        return 'bg-red-500';
      default:
        return 'bg-gray-400';
    }
  };

  const getStatusIcon = (status: ProgressItem['status']) => {
    switch (status) {
      case 'completed':
        return '✓';
      case 'failed':
        return '✗';
      default:
        return '↻';
    }
  };

  const formatDuration = (startTime: number) => {
    const duration = Math.floor((now - startTime) / 1000); // duration in seconds
    const hours = Math.floor(duration / 3600);
    const minutes = Math.floor((duration % 3600) / 60);
    const seconds = duration % 60;

    if (hours > 0) {
      return `${hours}h ${minutes}m ${seconds}s`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds}s`;
    } else {
      return `${seconds}s`;
    }
  };

  const handleConfirmCancel = (name: string) => {
    if (onCancel) {
      onCancel(name);
    } else {
      console.log('Cancel clicked for:', name);
    }
    setCancelId(null);
  };

  return (
    <div className="p-2 space-y-2">
      {displayItems.map((item, index) => (
        <div key={index} className="bg-white rounded-md shadow-sm">
          <div 
            className="p-2 cursor-pointer hover:bg-gray-50"
            onClick={() => toggleExpand(item.name)}
          >
            <div className="flex items-center justify-between mb-1.5">
              <div className="flex items-center gap-1.5">
                <span className="text-gray-500">
                  {expandedItems.has(item.name) ? (
                    <ChevronDown className="h-3.5 w-3.5" />
                  ) : (
                    <ChevronRight className="h-3.5 w-3.5" />
                  )}
                </span>
                <span className={`w-4 h-4 rounded-full flex items-center justify-center text-white text-[10px] ${getStatusColor(item.status)}`}>
                  {getStatusIcon(item.status)}
                </span>
                <span className="font-medium text-sm">{item.name}</span>
                <span className="text-xs text-muted-foreground">
                  {formatDuration(item.startTime)}
                </span>
              </div>
              {item.status === 'running' && (
                <div className="relative" onClick={e => e.stopPropagation()}>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-6 w-6 p-0 hover:bg-red-50 hover:text-red-500"
                    onClick={() => setCancelId(item.name)}
                  >
                    <X className="h-3.5 w-3.5" />
                  </Button>
                  {cancelId === item.name && (
                    <div className="absolute right-0 top-full mt-1 bg-background border rounded-lg shadow-lg py-1.5 px-2 z-50 whitespace-nowrap">
                      <p className="text-xs text-muted-foreground mb-1.5">Cancel this operation?</p>
                      <div className="flex gap-1.5">
                        <Button
                          size="sm"
                          variant="ghost"
                          className="h-6 px-2 text-xs"
                          onClick={() => setCancelId(null)}
                        >
                          No
                        </Button>
                        <Button
                          size="sm"
                          variant="destructive"
                          className="h-6 px-2 text-xs"
                          onClick={() => handleConfirmCancel(item.name)}
                        >
                          Yes, cancel
                        </Button>
                      </div>
                    </div>
                  )}
                </div>
              )}
            </div>
            <div className="flex items-center gap-2">
              <div className="flex-1 bg-gray-200 rounded-full h-1.5">
                <div
                  className={`h-1.5 rounded-full ${getStatusColor(item.status)} transition-all duration-500`}
                  style={{ width: `${item.progress}%` }}
                />
              </div>
              <div className="w-10 text-right">
                <span className="text-xs font-medium">
                  {item.progress}%
                </span>
              </div>
            </div>
          </div>
          {expandedItems.has(item.name) && (
            <div className="px-2 pb-2 pt-0">
              <div className="border rounded-md overflow-hidden">
                <Editor
                  height="150px"
                  defaultLanguage="sql"
                  value={item.sql}
                  options={{
                    readOnly: true,
                    minimap: { enabled: false },
                    scrollBeyondLastLine: false,
                    lineNumbers: 'off',
                    folding: false,
                    tabSize: 2,
                  }}
                />
              </div>
            </div>
          )}
        </div>
      ))}
    </div>
  );
} 