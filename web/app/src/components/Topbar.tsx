import { useEffect, useState } from "react";
import { LogOut, Moon, Search, Sun, Zap } from "lucide-react";
import clsx from "clsx";
import { useQueryClient } from "@tanstack/react-query";
import { useApplyPlan } from "../api/hooks";
import { logout, useMe } from "../api/auth";
import { useToast } from "./Toast";
import CommandPalette from "./CommandPalette";
import DiffModal from "./DiffModal";

export default function Topbar({ publicIP }: { publicIP?: string }) {
  const [dark, setDark] = useState(() => document.documentElement.classList.contains("dark"));
  const [paletteOpen, setPaletteOpen] = useState(false);
  const [diffOpen, setDiffOpen] = useState(false);
  const apply = useApplyPlan();
  const me = useMe();
  const qc = useQueryClient();
  const toast = useToast();

  async function doLogout() {
    try {
      await logout();
      qc.invalidateQueries({ queryKey: ["me"] });
      toast.show("Signed out", "ok");
    } catch (e: any) {
      toast.show(`Logout failed: ${e?.message}`, "bad");
    }
  }

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        setPaletteOpen((v) => !v);
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  function toggleTheme() {
    const next = !dark;
    setDark(next);
    document.documentElement.classList.toggle("dark", next);
    try {
      localStorage.setItem("zo-theme", next ? "dark" : "light");
    } catch {
      // ignore: localStorage may be unavailable (private mode, quota)
    }
  }

  const pending = apply.data?.changed === true;
  const allowed = apply.data?.allow_apply !== false;

  return (
    <>
      <header className="border-border bg-panel/90 dark:border-border-dark dark:bg-panel-dark/90 sticky top-0 z-30 flex h-14 items-center gap-3 border-b px-4 backdrop-blur lg:px-6">
        <div className="flex items-center gap-2 md:hidden">
          <div className="bg-accent grid h-7 w-7 place-items-center rounded-md">
            <span className="text-xs font-bold text-white">Z</span>
          </div>
          <span className="font-semibold">{publicIP || "ZeroOne"}</span>
        </div>
        <button
          onClick={() => setPaletteOpen(true)}
          className="border-border bg-bg text-muted hover:text-text dark:border-border-dark dark:bg-bg-dark dark:text-muted-dark dark:hover:text-text-dark hidden min-w-[18rem] items-center gap-2 rounded-lg border px-3 py-1.5 text-sm md:flex"
        >
          <Search size={14} />
          <span>Search users, rules, domains…</span>
          <kbd className="border-border bg-panel dark:border-border-dark dark:bg-panel-dark ml-auto rounded border px-1 py-0.5 font-mono text-[10px]">
            ⌘K
          </kbd>
        </button>
        <div className="flex-1" />
        {pending && (
          <button
            disabled={!allowed}
            onClick={() => setDiffOpen(true)}
            className={clsx("btn btn-primary", !allowed && "btn-danger")}
          >
            <Zap size={14} />
            {allowed ? "Review & deploy" : "Apply locked"}
          </button>
        )}
        <button onClick={toggleTheme} className="btn" aria-label="Toggle theme">
          {dark ? <Sun size={14} /> : <Moon size={14} />}
        </button>
        {me.data?.auth === "session" && (
          <button
            onClick={doLogout}
            className="btn"
            aria-label="Sign out"
            title={`Signed in as ${me.data.username}`}
          >
            <LogOut size={14} />
            <span className="hidden text-xs md:inline">{me.data.username}</span>
          </button>
        )}
      </header>
      <CommandPalette open={paletteOpen} onClose={() => setPaletteOpen(false)} />
      <DiffModal open={diffOpen} onClose={() => setDiffOpen(false)} />
    </>
  );
}
