import { Sparkles } from "lucide-react";
export default function RecommendationCard({
  postText,
  postUri,
  channelUsername,
  stocks,
  date,
  reasoning,
}) {
  return (
    <div>
      <div className="bg-[#0f172a] border border-slate-800 rounded-xl p-6 shadow-xl w-full max-w-2xl">
        <h2 className="text-xl font-bold mb-4 flex items-center gap-1.5 group">
          <span className="bg-gradient-to-br from-[#facc15] to-[#ffe082] bg-clip-text text-transparent">
            AI
          </span>

          <Sparkles
            className="text-[#facc15] size-4 relative -top-0.5 opacity-90 transition-opacity group-hover:opacity-100"
            strokeWidth={2.5}
          />
        </h2>
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
          {reasoning || "Текст обоснования отсутствует..."}
        </p>
      </div>
    </div>
  );
}
