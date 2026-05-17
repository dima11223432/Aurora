export default function BaseChannelCard({
  children,
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

        {children}
      </li>
    </div>
  );
}
