import { useNavigate } from "react-router-dom";
import { routes } from "./config/api";
import axios from "axios";

const checkIsAdmin = () => {
  try {
    const token = localStorage.getItem("token");
    if (!token) return false;
    const resp = axios.get(routes.isAdmin, {});
  } catch (e) {
    console.log(e);
  }
};

export default function Footer() {
  const navigate = useNavigate();
  return (
    <div className="fixed bottom-0 left-0 right-0 z-50">
      <div className="container mx-auto px-4 py-2">
        <div
          className="rounded-full p-3 border border-cyan-500"
          style={{ backgroundColor: "#151925" }}
        >
          <div className="flex items-center justify-around w-full ">
            <button
              onClick={() => {
                navigate("/feed");
              }}
              className="flex flex-col items-center justify-center py-2 px-1 text-cyan-300 hover:text-cyan hover:bg-cyan-600/20 rounded-2xl transition-colors"
            >
              <div className="w-8 h-8 mb-1 flex items-center justify-center">
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
                    d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                  />
                </svg>
              </div>
              <span className="text-xs">Новости</span>
            </button>

            <button
              onClick={() => {
                navigate("/settings");
              }}
              className="flex flex-col items-center justify-center py-2 px-1 text-cyan-300 hover:text-cyan hover:bg-cyan-600/20 rounded-2xl transition-colors"
            >
              <div className="w-8 h-8 mb-1 flex items-center justify-center">
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
                    d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"
                  />
                </svg>
              </div>
              <span className="text-xs">Настройки</span>
            </button>
            {/* <button className="flex flex-col items-center justify-center py-2 px-1 text-cyan-300 hover:text-cyan hover:bg-cyan-600/20 rounded-2xl transition-colors"> */}
            {/*   <div className="w-8 h-8 mb-1 flex items-center justify-center"> */}
            {/*     <svg */}
            {/*       className="w-6 h-6" */}
            {/*       fill="none" */}
            {/*       stroke="currentColor" */}
            {/*       viewBox="0 0 24 24" */}
            {/*     > */}
            {/*       <path */}
            {/*         strokeLinecap="round" */}
            {/*         strokeLinejoin="round" */}
            {/*         strokeWidth={1.5} */}
            {/*         d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" */}
            {/*       /> */}
            {/*     </svg> */}
            {/*   </div> */}
            {/*   <span className="text-xs">Админка</span> */}
            {/* </button> */}
          </div>
        </div>
      </div>
    </div>
  );
}
