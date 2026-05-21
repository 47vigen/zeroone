import { NavLink } from "react-router-dom";
import {
  Activity,
  BarChart3,
  Box,
  FileCode2,
  FileLock,
  Layers,
  Puzzle,
  ScrollText,
  Settings,
  Shield,
  Users as UsersIcon,
} from "lucide-react";
import clsx from "clsx";

const NAV: { to: string; label: string; icon: typeof Activity }[] = [
  { to: "/", label: "Overview", icon: Activity },
  { to: "/analytics", label: "Analytics", icon: BarChart3 },
  { to: "/users", label: "Users", icon: UsersIcon },
  { to: "/rules", label: "Rules", icon: Shield },
  { to: "/routes", label: "Routes", icon: Layers },
  { to: "/tunnels", label: "Tunnels", icon: Box },
  { to: "/logs", label: "Logs", icon: ScrollText },
  { to: "/xray-config", label: "Xray Config", icon: FileCode2 },
  { to: "/snapshots", label: "Snapshots", icon: FileLock },
  { to: "/plugins", label: "Plugins", icon: Puzzle },
  { to: "/settings", label: "Settings", icon: Settings },
];

export default function Sidebar({ publicIP }: { publicIP?: string }) {
  return (
    <aside className="hidden shrink-0 flex-col border-r border-border bg-panel dark:border-border-dark dark:bg-panel-dark md:flex md:w-60 lg:w-64">
      <div className="flex items-center gap-2 border-b border-border px-5 py-5 dark:border-border-dark">
        <div className="grid h-8 w-8 place-items-center rounded-lg bg-accent">
          <span className="text-sm font-bold text-white">X</span>
        </div>
        <div className="leading-tight">
          <div className="text-sm font-bold tracking-tight">Xray Stack</div>
          <div className="text-xs text-muted dark:text-muted-dark">{publicIP || "—"}</div>
        </div>
      </div>
      <nav className="flex-1 space-y-0.5 overflow-y-auto p-3">
        {NAV.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.to === "/"}
            className={({ isActive }) => clsx("nav-item", isActive && "active")}
          >
            <span className="nav-mark" />
            <item.icon size={16} className="shrink-0" />
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>
      <div className="border-t border-border px-5 py-3 text-xs text-muted dark:border-border-dark dark:text-muted-dark">
        v0.2 · Cloudflare-style
      </div>
    </aside>
  );
}
