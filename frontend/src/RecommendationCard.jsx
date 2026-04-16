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
      <p>{postText}</p>
      <p>{postUri}</p>
      <p>{channelUsername}</p>
      <p>{date}</p>
      <p>{reasoning}</p>
    </div>
  );
}
