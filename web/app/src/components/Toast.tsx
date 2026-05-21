import { createContext, useCallback, useContext, useEffect, useState } from "react";
import clsx from "clsx";

type Tone = "ok" | "bad" | "warn" | "info";
type ToastItem = { id: number; text: string; tone: Tone };

const ToastContext = createContext<{ show: (text: string, tone?: Tone) => void } | null>(null);

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [items, setItems] = useState<ToastItem[]>([]);
  const show = useCallback((text: string, tone: Tone = "info") => {
    const id = Date.now() + Math.random();
    setItems((prev) => [...prev, { id, text, tone }]);
    setTimeout(() => setItems((prev) => prev.filter((t) => t.id !== id)), 4500);
  }, []);
  return (
    <ToastContext.Provider value={{ show }}>
      {children}
      <div className="fixed right-4 bottom-4 z-50 flex max-w-sm flex-col gap-2">
        {items.map((t) => (
          <ToastView
            key={t.id}
            item={t}
            onDismiss={() => setItems((p) => p.filter((x) => x.id !== t.id))}
          />
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) return { show: (_t: string, _tone?: Tone) => {} };
  return ctx;
}

function ToastView({ item, onDismiss }: { item: ToastItem; onDismiss: () => void }) {
  const [show, setShow] = useState(false);
  useEffect(() => {
    setShow(true);
  }, []);
  return (
    <div
      className={clsx(
        "bg-panel shadow-elev dark:bg-panel-dark flex items-center gap-2 rounded-lg border px-3 py-2 text-sm transition-all",
        show ? "translate-y-0 opacity-100" : "translate-y-2 opacity-0",
        item.tone === "ok" && "border-ok/30 text-ok dark:text-ok-dark",
        item.tone === "bad" && "border-bad/30 text-bad dark:text-bad-dark",
        item.tone === "warn" && "border-warn/30 text-warn dark:text-warn-dark",
        item.tone === "info" && "border-border dark:border-border-dark",
      )}
      onClick={onDismiss}
    >
      {item.text}
    </div>
  );
}
