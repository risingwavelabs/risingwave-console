"use client";

import React, { useEffect, useState, useRef } from 'react';
import { ReactFlow, Background, Controls, Panel, MarkerType, Position, Handle, useNodesState, useEdgesState } from '@xyflow/react';
import dagre from '@dagrejs/dagre';
import '@xyflow/react/dist/style.css';
import { ChevronDown, ChevronRight } from 'lucide-react';
import { useTheme } from 'next-themes';
import { DefaultService } from "@/api-gen"

// RisingWave Schema Types - Single Source of Truth
export interface TableColumn {
  name: string;
  type: string;
  isPrimary?: boolean;
}

export type NodeType = 'source' | 'sink' | 'materialized view' | 'table' | 'system table';

export interface RisingWaveNodeData extends Record<string, unknown> {
  id: number;
  name: string;
  type: NodeType;
  columns: TableColumn[];
  connector?: {
    type: string;
    properties?: Record<string, string>;
  };
  format?: string;
  dependencies?: number[]; // Array of node IDs that this node depends on
  throughput?: number; // Added throughput property to track the node's throughput
}

interface StreamingGraphProps {
  clusterId?: number | null;
  data?: RisingWaveNodeData[];
  className?: string;
  height?: string | number;
}

const COLORS = {
  light: {
    source: '#e3f2fd',
    sink: '#fce4ec',
    materializedView: '#e8f5e9',
    table: '#fff3e0',
    text: '#1a1a1a',
    background: '#fafafa',
    border: '#e2e8f0',
    nodeBackground: '#ffffff',
    systemTable: '#e0f2f1',
  },
  dark: {
    source: '#1e293b',
    sink: '#3f2a3f',
    materializedView: '#1e312b',
    table: '#2c2620',
    text: '#e2e8f0',
    background: '#171717',
    border: '#262626',
    nodeBackground: '#1e1e1e',
    systemTable: '#1e312b',
  }
};

const ICONS: Record<string, string> = {
  source: '📥',
  sink: '📤',
  materializedView: '📊',
  table: '📋'
};

// Custom node component for RisingWave nodes
const TableNodeComponent = React.memo(({ data }: { data: RisingWaveNodeData }) => {
  const [isExpanded, setIsExpanded] = React.useState(false);
  const { theme = 'light' } = useTheme();
  const nodeType = data.type === 'materialized view' ? 'materializedView' :
    data.type === 'system table' ? 'systemTable' : data.type;

  // Memoize static values
  const colors = React.useMemo(() => COLORS[theme === 'dark' ? 'dark' : 'light'], [theme]);
  const backgroundColor = React.useMemo(() => colors[nodeType] || colors.table, [colors, nodeType]);
  const icon = React.useMemo(() => ICONS[nodeType] || ICONS.table, [nodeType]);

  // Memoize event handler
  const handleExpand = React.useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    setIsExpanded(prev => !prev);
  }, []);

  // Memoize the header content
  const headerContent = React.useMemo(() => (
    <span className="flex items-center gap-2 flex-1">
      <span>{icon}</span>
      <span className="text-foreground">{data.name}</span>
    </span>
  ), [icon, data.name]);

  // Format throughput value for display
  const formattedThroughput = React.useMemo(() => {
    // Check explicitly for undefined, not truthiness (to handle zero values)
    if (data.throughput === undefined) {
      return null;
    }
    // Special handling for zero
    if (data.throughput === 0) {
      return "0.00 rows/sec";
    }

    // Format to readable values (e.g., convert to K, M, B if large)
    if (data.throughput >= 1_000_000_000) {
      return `${(data.throughput / 1_000_000_000).toFixed(2)} B rows/sec`;
    } else if (data.throughput >= 1_000_000) {
      return `${(data.throughput / 1_000_000).toFixed(2)} M rows/sec`;
    } else if (data.throughput >= 1_000) {
      return `${(data.throughput / 1_000).toFixed(2)} K rows/sec`;
    } else {
      return `${data.throughput.toFixed(2)} rows/sec`;
    }
  }, [data.throughput, data.id, data.name]);

  return (
    <div className={`rounded-lg shadow-lg border min-w-[100px] drag-handle bg-background`}>
      <Handle
        type="target"
        position={Position.Left}
        id="target"
        className="!bg-muted-foreground"
      />
      <div
        className="px-4 py-2 font-semibold border-b text-sm flex items-center gap-2"
        style={{ backgroundColor }}
      >
        <button
          className="p-0.5 hover:bg-black/5 dark:hover:bg-white/5 rounded"
          onClick={handleExpand}
        >
          {isExpanded ? (
            <ChevronDown className="h-4 w-4 text-muted-foreground" />
          ) : (
            <ChevronRight className="h-4 w-4 text-muted-foreground" />
          )}
        </button>
        {headerContent}
        {data.connector && (
          <span className="ml-auto text-xs text-muted-foreground">({data.connector.type})</span>
        )}
      </div>
      {/* Add throughput display */}
      <div className="px-4 py-1 text-xs border-b bg-opacity-50"
        style={{ backgroundColor: `${backgroundColor}80` }}>
        <div className="flex items-center gap-1">
          <span className="font-semibold">Throughput:</span>
          {(() => {
            return formattedThroughput ? (
              <span className="ml-1 font-mono">{formattedThroughput}</span>
            ) : (
              <span className="ml-1 text-muted-foreground italic">Not available</span>
            );
          })()}
        </div>
      </div>
      {isExpanded && (
        <div className="p-2">
          {data.columns?.map((column: TableColumn, index: number) => (
            <div key={index} className="flex items-center text-sm py-1">
              <span className="mr-2">
                {column.isPrimary && '🔑'}
              </span>
              <span className="font-medium text-foreground">{column.name}</span>
              <span className="ml-2 text-muted-foreground text-xs">({column.type})</span>
            </div>
          ))}
          {data.format && (
            <div className="mt-2 text-xs text-muted-foreground border-t pt-2">
              Format: {data.format}
            </div>
          )}
        </div>
      )}
      <Handle
        type="source"
        position={Position.Right}
        id="source"
        className="!bg-muted-foreground"
      />
    </div>
  );
}, (prevProps, nextProps) => {
  // Custom comparison function for memo - make sure all relevant properties are compared
  const prevData = prevProps.data;
  const nextData = nextProps.data;

  const areEqual =
    prevData.id === nextData.id &&
    prevData.name === nextData.name &&
    prevData.type === nextData.type &&
    prevData.throughput === nextData.throughput;

  // Debug output for memoization
  if (prevData.throughput !== nextData.throughput) {
    console.debug(`Throughput changed for ${nextData.name}: ${prevData.throughput} -> ${nextData.throughput}`);
  }

  return areEqual;
});

TableNodeComponent.displayName = 'TableNodeComponent';

export function StreamingGraph({ clusterId, data = [], className = '', height = '100%' }: StreamingGraphProps) {
  const { theme = 'light' } = useTheme();
  const colors = COLORS[theme === 'dark' ? 'dark' : 'light'];
  const [throughputData, setThroughputData] = useState<Record<string, number>>({});
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  async function getMaterializedViewThroughput() {
    if (clusterId) {
      try {
        const matrix = await DefaultService.getMaterializedViewThroughput(clusterId);

        if (matrix && matrix.length > 0) {
          // Process the throughput data for each materialized view
          const newThroughputData: Record<string, number> = {};
          matrix.forEach(series => {
            // Extract the table_id from the metric to match with node id
            const tableId = series.metric?.['table_id'];

            if (tableId && series.values && series.values.length > 0) {
              // Get the latest value
              const latestValue = series.values[series.values.length - 1];
              const throughputValue = Number(latestValue[1]);

              // Make sure the tableId is a string to match with node.id.toString()
              const tableIdStr = String(tableId);

              // Store the throughput value even if it's zero
              newThroughputData[tableIdStr] = throughputValue;
            } else {
              console.log(`Missing table_id or values in series:`, series);
            }
          });
          setThroughputData(newThroughputData);
        } else {
          console.log(`No throughput data received from API`);
        }
      } catch (error) {
        console.error('Error fetching materialized view throughput:', error);
      }
    }
  }

  useEffect(() => {
    // Initial fetch
    getMaterializedViewThroughput();

    // Set up interval to fetch every 5 seconds
    intervalRef.current = setInterval(() => {
      getMaterializedViewThroughput();
    }, 5000);

    // Clean up interval on unmount
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [clusterId]);  // Only re-run if clusterId changes

  type NodeData = {
    data: RisingWaveNodeData;
    id: string;
    position: { x: number; y: number };
    type?: string;
    dragHandle?: string;
  };

  type EdgeData = {
    id: string;
    source: string;
    target: string;
    sourceHandle?: string;
    targetHandle?: string;
    animated?: boolean;
    style?: Record<string, string | number | boolean>;
    markerEnd?: { type: MarkerType };
    type?: string;
  };

  const [nodes, setNodes, onNodesChange] = useNodesState<NodeData>([]);
  const [edges, setEdges] = useEdgesState<EdgeData>([]);

  const getLayoutedElements = (nodes: NodeData[], edges: EdgeData[]) => {
    const dagreGraph = new dagre.graphlib.Graph();
    dagreGraph.setDefaultEdgeLabel(() => ({}));
    dagreGraph.setGraph({ rankdir: 'LR' });

    nodes.forEach((node) => {
      dagreGraph.setNode(node.id, { width: 250, height: 100 });
    });

    edges.forEach((edge) => {
      dagreGraph.setEdge(edge.source, edge.target);
    });

    dagre.layout(dagreGraph);

    return {
      nodes: nodes.map((node) => {
        const nodeWithPosition = dagreGraph.node(node.id);
        return {
          ...node,
          position: {
            x: nodeWithPosition.x - 125,
            y: nodeWithPosition.y - 50,
          },
        };
      }),
      edges,
    };
  };

  React.useEffect(() => {
    const generatedNodes = data.map((node) => {
      // Add throughput data to materialized view nodes
      let nodeWithThroughput = { ...node };
      // Convert node.id to string for consistent comparison
      const nodeIdStr = node.id.toString();
      // Check if throughputData has the key, not if the value is truthy
      if (nodeIdStr in throughputData) {
        const throughputValue = Number(throughputData[nodeIdStr]);
        // Ensure it's a valid number
        if (!isNaN(throughputValue)) {
          nodeWithThroughput.throughput = throughputValue;
        } else {
          console.log(`Invalid throughput value for node ${nodeIdStr}: ${throughputData[nodeIdStr]}`);
        }
      }


      return {
        id: node.id.toString(),
        position: { x: 0, y: 0 },
        data: nodeWithThroughput,
        type: 'tableNode',
        dragHandle: '.drag-handle'
      };
    }) as NodeData[];

    const generatedEdges = data.flatMap((node) => {
      if (!node.dependencies) return [];

      return node.dependencies.map((depId) => ({
        id: `dep-${depId}-${node.id}`,
        source: depId.toString(),
        target: node.id.toString(),
        sourceHandle: 'source',
        targetHandle: 'target',
        animated: node.type === 'materialized view',
        style: {
          stroke: node.type === 'materialized view' ? '#2196f3' : '#888',
          strokeWidth: node.type === 'materialized view' ? 2 : 1
        },
        markerEnd: {
          type: MarkerType.Arrow
        },
        type: 'default'
      }));
    }) as EdgeData[];

    const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
      generatedNodes,
      generatedEdges
    );

    setNodes(layoutedNodes);
    setEdges(layoutedEdges);
  }, [data, throughputData, setNodes, setEdges]);

  const nodeTypes = React.useMemo(() => ({
    tableNode: TableNodeComponent
  }), []);

  return (
    <div style={{ width: '100%', height: height }} className={className}>
      <div style={{ width: '100%', height: '100%' }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          nodeTypes={nodeTypes}
          fitView
          defaultViewport={{ x: 0, y: 0, zoom: 0.5 }}
          minZoom={0.3}
          maxZoom={1.2}
          defaultEdgeOptions={{
            type: 'default',
            animated: false
          }}
          style={{ background: colors.background }}
          proOptions={{
            hideAttribution: true
          }}
          fitViewOptions={{
            padding: 0.2,
            includeHiddenNodes: true,
            minZoom: 0.3,
            maxZoom: 1.2
          }}
          nodesConnectable={false}
          edgesFocusable={false}
          snapToGrid={true}
          snapGrid={[15, 15]}
          elevateNodesOnSelect={false}
          nodesDraggable={true}
          panOnDrag={true}
          selectionOnDrag={false}
          zoomOnScroll={false}
          panOnScroll={true}
          preventScrolling={true}
          className="dark:react-flow-dark"
        >
          <Background className="dark:bg-background" />
          <Controls className="!bg-background !border-border [&>button]:!bg-background [&>button]:!text-foreground [&>button:hover]:!bg-accent" />
          <Panel position="top-left" className="bg-background border text-foreground p-2 rounded shadow-md">
            <div className="text-sm font-medium mb-2">Legend</div>
            <div className="flex flex-col gap-1 text-xs">
              <div className="flex items-center">
                <div className="w-3 h-3 mr-2" style={{ backgroundColor: colors.source }}></div>
                <span className="mr-1">{ICONS.source}</span>
                Source
              </div>
              <div className="flex items-center">
                <div className="w-3 h-3 mr-2" style={{ backgroundColor: colors.sink }}></div>
                <span className="mr-1">{ICONS.sink}</span>
                Sink
              </div>
              <div className="flex items-center">
                <div className="w-3 h-3 mr-2" style={{ backgroundColor: colors.materializedView }}></div>
                <span className="mr-1">{ICONS.materializedView}</span>
                Materialized View
              </div>
              <div className="flex items-center">
                <div className="w-3 h-3 mr-2" style={{ backgroundColor: colors.table }}></div>
                <span className="mr-1">{ICONS.table}</span>
                Table
              </div>
              <div className="flex items-center mt-2">
                <div className="w-4 h-0 border-t-2 mr-2 border-[#2196f3]"></div>
                Streaming Data Flow
              </div>
            </div>
          </Panel>
        </ReactFlow>
      </div>
    </div>
  );
}
