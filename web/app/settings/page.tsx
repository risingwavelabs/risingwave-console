"use client"

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const CURRENT_VERSION = "0.1.0";
const MODEL_OPTIONS = [
  { value: "gpt-4", label: "GPT-4" },
  { value: "gpt-3.5-turbo", label: "GPT-3.5 Turbo" },
  { value: "claude-2", label: "Claude 2" },
];

export default function SettingsPage() {
  const handleCheckUpdate = async () => {
    // TODO: Implement update check logic
    console.log("Checking for updates...");
  };

  return (
    <div className="p-6">
      <div className="space-y-6 max-w-2xl">
        <div className="space-y-4">
          <h2 className="text-lg font-semibold">Agent</h2>
          <div className="space-y-2">
            <label className="text-sm font-medium">Model Selection</label>
            <Select defaultValue="gpt-4">
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select a model" />
              </SelectTrigger>
              <SelectContent>
                {MODEL_OPTIONS.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Input type="password" placeholder="Enter OpenAI API key" className="mt-2" />
            <label className="text-sm font-medium block mt-4">Prompt</label>
            <Input placeholder="Enter your prompt" />
          </div>
        </div>

        <div className="space-y-4">
          <h2 className="text-lg font-semibold">About</h2>
          <div className="flex items-center justify-between border rounded-lg p-4 bg-muted/50">
            <div>
              <p className="text-sm font-medium">Current Version</p>
              <p className="text-sm text-muted-foreground">{CURRENT_VERSION}</p>
            </div>
            <Button onClick={handleCheckUpdate} variant="outline">
              Check for Updates
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
