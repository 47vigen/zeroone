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
    <div className="bg-bg dark:bg-bg-dark flex min-h-full items-center justify-center p-6">
      <form onSubmit={submit} className="panel panel-pad w-full max-w-sm shadow-lg">
        <div className="mb-4 flex items-center gap-2">
          <div className="bg-accent/10 text-accent rounded-lg p-2">
            <KeyRound size={18} />
          </div>
          <div>
            <h1 className="text-lg leading-tight font-semibold">ZeroOne</h1>
            <p className="text-muted dark:text-muted-dark text-xs">Sign in to the control panel</p>
          </div>
        </div>

        {bootstrapNeeded && (
          <div className="border-warn/40 bg-warn/10 mb-4 flex items-start gap-2 rounded-md border p-3 text-xs">
            <ShieldAlert size={14} className="text-warn mt-0.5 shrink-0" />
            <div>
              No admin accounts exist yet. Create the first admin by calling
              <code className="bg-bg dark:bg-bg-dark mx-1 rounded px-1 py-0.5 font-mono">
                POST /api/admins
              </code>
              with an existing panel Bearer token, then return here to sign in.
            </div>
          </div>
        )}

        <label className="mb-3 block">
          <div className="text-muted dark:text-muted-dark mb-1 text-xs">Username</div>
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
          <div className="text-muted dark:text-muted-dark mb-1 text-xs">Password</div>
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
          <div className="border-bad/30 bg-bad/5 text-bad dark:text-bad-dark mb-3 rounded border p-2 text-xs">
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
