import "./App.css";
import { Chat } from "./components/Chat";
import { ServerList } from "./components/ServerList";

function App() {
  return (
    <div className="flex h-screen w-screen">
      <ServerList />
      <Chat />
    </div>
  );
}

export default App;
