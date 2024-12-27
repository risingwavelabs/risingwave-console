export default function Home() {
  return (
    <div className="flex flex-col p-8">
      <h1 className="text-2xl font-semibold mb-4">Dashboard</h1>
      <div className="grid gap-4 md:grid-cols-2">
        <div className="p-6 bg-card rounded-lg border">
          <h2 className="text-lg font-medium mb-2">Quick Actions</h2>
          <p className="text-muted-foreground">Access your most used features here.</p>
        </div>
        <div className="p-6 bg-card rounded-lg border">
          <h2 className="text-lg font-medium mb-2">Recent Activity</h2>
          <p className="text-muted-foreground">View your recent activities and updates.</p>
        </div>
      </div>
    </div>
  );
}
