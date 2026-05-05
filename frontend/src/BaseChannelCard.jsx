export default function BaseChannelCard({
  idx,
  isChecked,
  channelName,
  handleCheckboxChange,
}) {
  return (
    <div>
      <li
        key={idx}
        className="py-2 px-3 rounded-lg hover:bg-cyan-500/70 hover:scale-[1.03] hover:shadow-lg text-white cursor-pointer transition-all duration-200 flex items-center gap-2 group"
        style={{ backdropFilter: "blur(1px)" }}
      >
        <input
          type="checkbox"
          id={`channel-${idx}`}
          checked={isChecked}
          onChange={() => handleCheckboxChange(channelName)}
          className="w-4 h-4 text-cyan-600 bg-gray-100 border-gray-300 rounded focus:ring-cyan-500 focus:ring-2"
        />
        <label
          htmlFor={`channel-${idx}`}
          className="flex-1 group-hover:text-cyan-200 transition-colors duration-200 cursor-pointer"
        >
          {channelName}
        </label>
        <div>
          <button className="text-cyan-300 hover:text-cyan hover:bg-cyan-300/20 rounded-2xl transition-colors">
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              />
            </svg>
          </button>
        </div>
      </li>
    </div>
  );
}
