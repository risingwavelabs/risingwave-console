"use client"

import "./globals.css";
import { ThemeProvider } from "next-themes"
import { SidebarProvider } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/ui/app-sidebar"
import { usePathname } from "next/navigation"

export default function Layout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const isLoginPage = pathname === "/login"

  return (
    <html lang="en" className="h-full" suppressHydrationWarning>
      <body className="h-full w-full">
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          {isLoginPage ? (
            children
          ) : (
            <SidebarProvider defaultOpen={false} className="h-full w-full">
              <div className="flex h-full w-full">
                <AppSidebar />
                <div className="flex-1 h-full w-full">
                  {children}
                </div>
              </div>
            </SidebarProvider>
          )}
        </ThemeProvider>
      </body>
    </html>
  )
}
