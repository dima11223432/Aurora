import BaseChannelCard from "./BaseChannelCard";

export default function UserCustomChannelCard({
  idx,
  isChecked,
  channelName,
  handleCheckboxChange,
  deleteUserCustomParsingChannel,
}) {
  return (
    <BaseChannelCard
      idx={idx}
      isChecked={isChecked}
      channelName={channelName}
      handleCheckboxChange={handleCheckboxChange}
      deleteUserCustomParsingChannel={deleteUserCustomParsingChannel}
    >
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
    </BaseChannelCard>
  );
}
