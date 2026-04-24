export default function RecommendationCard({
  postText,
  postUri,
  channelUsername,
  stocks,
  date,
  reasoning,
}) {
  return (
    <div className="bg-[#0f172a] rounded-xl border border-slate-700 p-6 w-full max-w-2xl">
      <h3 className="text-[#facc15] font-bold text-xl mb-4">
        Рекомендация ИИ
      </h3>
      
      <div className="flex flex-wrap gap-2 mb-4">
        {stocks?.map((stock, index) => (
          <span
            key={index}
            className={`px-3 py-1 rounded-md uppercase font-bold text-xs tracking-wide ${
              stock.side === "buy"
                ? "bg-[#14532d] text-[#4ade80]"
                : stock.side === "sell"
                ? "bg-[#7f1d1d] text-[#f87171]"
                : "bg-slate-700 text-slate-300"
            }`}
          >
            {stock.stockName} ({stock.side})
          </span>
        ))}
      </div>
      
      <p className="text-slate-400 leading-relaxed text-sm">
        {reasoning ||
          "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua..."}
      </p>
    </div>
  );
}