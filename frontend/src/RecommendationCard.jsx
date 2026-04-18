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
      <h2 className="text-[#facc15] text-xl ">Рекомендация ИИ</h2>
      <div className="flex items-center gap-3 mb-4">
        <span className="text-lg font-semibold text-slate-200">
          {stocks?.map((stock, index) => (
            <span
              key={index}
              className={`text-xs px-2 py-1 rounded-full ${
                stock.side === "buy"
                  ? "bg-[#14532d] text-[#4ade80]"
                  : stock.side === "sell"
                    ? "bg-red-600"
                    : "bg-gray-600"
              }`}
            >
              {stock.stockName} ({stock.side})
            </span>
          ))}
        </span>
      </div>
      <p className="text-slate-400 ">
        {reasoning ||
          "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua..."}
      </p>
    </div>
  );
}
