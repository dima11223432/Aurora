export default function ErrorPage() {
  return (
    <div className="flex flex-col bg-gray-900 items-center justify-center h-screen">
      <div className="border-gray-500 bg-gray-800 border-[1px] border-gray-600 rounded-[10%] flex flex-col items-center justify-center w-1/2 h-1/2">
        <p className="text-4xl font-bold text-[#0fd2f5] justify-center">404</p>
        <p className="text-2xl font-bold text-gray-400 justify-center">
          Упс... Кажется, что <br />
          страница не найдена
        </p>
      </div>
    </div>
  );
}
