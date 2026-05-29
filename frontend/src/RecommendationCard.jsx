import { Sparkles } from "lucide-react";
import { useState } from "react";
export default function RecommendationCard({
  postText,
  postUri,
  channelUsername,
  stocks,
  date,
  reasoning,
}) {
  const [expanded, setExpanded] = useState(false);
  const words = (reasoning || "").split(/\s+/).filter(Boolean);
  const preview = words.slice(0, 30).join(" ");
  const needsToggle = words.length > 30;
  return (
    <div>
      <div className="bg-[#0f172a] border border-slate-800 rounded-xl p-6 shadow-xl w-full max-w-2xl">
        <div className="flex justify-between">
          <h2 className="text-xl font-bold mb-4 flex items-center gap-1.5 group">
            <span className="bg-gradient-to-br from-[#facc15] to-[#ffe082] bg-clip-text text-transparent">
              AI
            </span>

            <Sparkles
              className="text-[#facc15] size-4 relative -top-0.5 opacity-90 transition-opacity group-hover:opacity-100"
              strokeWidth={2.5}
            />
          </h2>
          <div className="flex justify-end">
            <p className="text-stone-400">
              {new Date(date).toLocaleDateString("ru-RU", {
                day: "numeric",
                month: "long",
                year: "numeric",
                hour: "2-digit",
                minute: "2-digit",
              })}
            </p>
          </div>
        </div>
        <div className="flex flex-wrap gap-2 mb-4">
          {stocks?.map((stock, index) => (
            <span
              key={index}
              className={`text-xs font-bold px-3 py-1 rounded-full uppercase tracking-wider ${
                stock.side === "buy"
                  ? "bg-[#14532d] text-[#4ade80]"
                  : stock.side === "sell"
                    ? "bg-red-900/50 text-red-400 border border-red-800"
                    : "bg-slate-700 text-slate-300"
              }`}
            >
              {stock.stockName} • {stock.side}
            </span>
          ))}
        </div>
        <p className="text-slate-400 leading-relaxed">
          {reasoning
            ? expanded
              ? reasoning
              : needsToggle
                ? preview + "..."
                : reasoning
            : "Текст обоснования отсутствует..."}
        </p>
        {needsToggle && (
          <div className="flex justify-end mt-2">
            <button
              onClick={() => setExpanded(!expanded)}
              className="flex items-center gap-1 text-gray-400 hover:text-cyan-400 text-sm font-medium transition-colors duration-200"
            >
              <span>{expanded ? "Свернуть" : "Развернуть"}</span>
              <svg
                className={`w-4 h-4 transition-transform duration-200 ${expanded ? "rotate-180" : ""}`}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 9l-7 7-7-7"
                />
              </svg>
            </button>
          </div>
        )}
        <div className="flex justify-end">
          <a className="text-cyan-500" href={postUri}>
            Оригинальный пост
          </a>
        </div>
      </div>
    </div>
  );
}
