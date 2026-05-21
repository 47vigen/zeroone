import { useState } from "react";
import { KeyRound, LogIn, ShieldAlert } from "lucide-react";
import { login } from "../api/auth";

export default function Login({
  onLoggedIn,
  bootstrapNeeded,
}: {
  onLoggedIn: () => void;
  bootstrapNeeded: boolean;
}) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setPending(true);
    try {
      await login(username.trim(), password);
      onLoggedIn();
    } catch (err: any) {
      setError(err?.message || "Login failed");
    } finally {
      setPending(false);
    }
  }

  return (
    <div className="flex min-h-full items-center justify-center bg-bg p-6 dark:bg-bg-dark">
      <form onSubmit={submit} className="panel panel-pad w-full max-w-sm shadow-lg">
        <div className="mb-4 flex items-center gap-2">
          <div className="rounded-lg bg-accent/10 p-2 text-accent">
            <KeyRound size={18} />
          </div>
          <div>
            <h1 className="text-lg font-semibold leading-tight">Xray Stack</h1>
            <p className="text-xs text-muted dark:text-muted-dark">Sign in to the control panel</p>
          </div>
        </div>

        {bootstrapNeeded && (
          <div className="mb-4 flex items-start gap-2 rounded-md border border-warn/40 bg-warn/10 p-3 text-xs">
            <ShieldAlert size={14} className="mt-0.5 shrink-0 text-warn" />
            <div>
              No admin accounts exist yet. Create the first admin by calling
              <code className="mx-1 rounded bg-bg px-1 py-0.5 font-mono dark:bg-bg-dark">
                POST /api/admins
              </code>
              with an existing panel Bearer token, then return here to sign in.
            </div>
          </div>
        )}

        <label className="mb-3 block">
          <div className="mb-1 text-xs text-muted dark:text-muted-dark">Username</div>
          <input
            className="input"
            autoFocus
            autoComplete="username"
            required
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </label>
        <label className="mb-4 block">
          <div className="mb-1 text-xs text-muted dark:text-muted-dark">Password</div>
          <input
            className="input"
            type="password"
            autoComplete="current-password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </label>

        {error && (
          <div className="mb-3 rounded border border-bad/30 bg-bad/5 p-2 text-xs text-bad dark:text-bad-dark">
            {error}
          </div>
        )}

        <button
          type="submit"
          className="btn btn-primary w-full justify-center"
          disabled={pending || !username || !password}
        >
          <LogIn size={14} /> {pending ? "Signing in…" : "Sign in"}
        </button>
      </form>
    </div>
  );
}
