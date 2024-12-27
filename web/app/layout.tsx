import "./globals.css";

import { SidebarProvider } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/ui/app-sidebar"

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="h-full">
      <body className="h-full w-screen">
        <SidebarProvider className="h-full w-full">
          <AppSidebar />
          <main className="h-full w-full">
            {children}
          </main>
        </SidebarProvider>
      </body>
    </html>
  )
}
