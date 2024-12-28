import "./globals.css";

import { SidebarProvider } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/ui/app-sidebar"

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="h-full">
      <body className="h-full w-screen overflow-hidden">
        <SidebarProvider defaultOpen={false} className="h-full w-full">
          <div className="flex h-full w-full">
            <AppSidebar />
            <div className="flex-1 h-full w-full">
              {children}
            </div>
          </div>
        </SidebarProvider>
      </body>
    </html>
  )
}
