import { ServerIcon, CodeIcon, ExternalLinkIcon } from "lucide-react"
import { FaGithub } from "react-icons/fa"
import { FaSlack, FaLinkedin, FaTwitter, FaYoutube } from "react-icons/fa"
import Link from "next/link"

export default function Home() {
  return (
    <div className="flex flex-col p-8">
      <h1 className="text-2xl font-semibold mb-4">WaveKit</h1>
      <div className="grid gap-4 md:grid-cols-2 mb-8">
        <Link 
          href="/clusters" 
          className="p-6 bg-card rounded-lg border transition-colors hover:bg-accent/50"
        >
          <div className="flex items-center gap-2 mb-2">
            <ServerIcon className="h-5 w-5" />
            <h2 className="text-lg font-medium">Clusters</h2>
          </div>
          <p className="text-muted-foreground">Manage your database clusters and monitor their status.</p>
        </Link>
        <Link 
          href="/sqlconsole" 
          className="p-6 bg-card rounded-lg border transition-colors hover:bg-accent/50"
        >
          <div className="flex items-center gap-2 mb-2">
            <CodeIcon className="h-5 w-5" />
            <h2 className="text-lg font-medium">SQL Console</h2>
          </div>
          <p className="text-muted-foreground">Write and execute SQL queries, analyze data in real-time.</p>
        </Link>
      </div>

      <div className="p-6 bg-card rounded-lg border">
        <h2 className="text-lg font-medium mb-4">Resources & Documentation</h2>
        <div className="space-y-4">
          <p className="text-muted-foreground">
            WaveKit is a modern management UI for RisingWave, the cloud-native streaming database. 
            Here are some helpful resources to get you started:
          </p>
          <div className="grid gap-3 text-sm">
            <a 
              href="https://risingwave.com" 
              target="_blank" 
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-primary hover:underline"
            >
              <ExternalLinkIcon className="h-4 w-4" />
              Visit RisingWave's official website
            </a>
            <a 
              href="https://docs.risingwave.com" 
              target="_blank" 
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-primary hover:underline"
            >
              <ExternalLinkIcon className="h-4 w-4" />
              Read the RisingWave documentation
            </a>
            <a 
              href="https://github.com/risingwavelabs/risingwave" 
              target="_blank" 
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-primary hover:underline"
            >
              <FaGithub className="h-4 w-4" />
              RisingWave on GitHub
            </a>
            <a 
              href="https://github.com/risingwavelabs/wavekit" 
              target="_blank" 
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-primary hover:underline"
            >
              <FaGithub className="h-4 w-4" />
              WaveKit on GitHub
            </a>
          </div>
          <div className="border-t pt-4 mt-6">
            <h3 className="text-sm font-medium mb-3">Connect with RisingWave</h3>
            <div className="flex gap-4">
              <a
                href="https://go.risingwave.com/slack"
                target="_blank"
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-primary transition-colors"
                title="Join us on Slack"
              >
                <FaSlack className="h-5 w-5" />
              </a>
              <a
                href="https://go.risingwave.com/linkedin"
                target="_blank"
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-primary transition-colors"
                title="Follow us on LinkedIn"
              >
                <FaLinkedin className="h-5 w-5" />
              </a>
              <a
                href="https://go.risingwave.com/twitter"
                target="_blank"
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-primary transition-colors"
                title="Follow us on X (Twitter)"
              >
                <FaTwitter className="h-5 w-5" />
              </a>
              <a
                href="https://www.youtube.com/channel/UCsHwdyBRxBpmkA5RRd0YNEA"
                target="_blank"
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-primary transition-colors"
                title="Subscribe to our YouTube channel"
              >
                <FaYoutube className="h-5 w-5" />
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
