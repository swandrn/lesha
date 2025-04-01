import { useEffect, useState } from "react";
import { ServerList } from "./components/ServerList";
import { ChannelList } from "./components/ChannelList";
import { Chat } from "./components/Chat";
import EditAccount from "./components/EditAccount";
import { FriendList } from "./components/FriendList";
import CreateServer from "./components/CreateServer";
import Login from "./components/Login";
import Register from "./components/Register";
import { Route, Routes } from "react-router-dom";

function App() {
  const [selectedServer, setSelectedServer] = useState<number | null>(null);
  const [selectedChannel, setSelectedChannel] = useState<number | null>(null);
  const [isChannelListVisible, setIsChannelListVisible] = useState(true);
  const [isCreatingServer, setIsCreatingServer] = useState(false);

  const [servers, setServers] = useState([
    { id: 0, name: "Edit Account", icon: "🛠️" },
    { id: 1, name: "Friends", icon: "👫" },
    { id: 3, name: "Programming", icon: "💻" },
    { id: 4, name: "Music", icon: "🎵" },
    { id: 5, name: "Movies", icon: "🎬" },
  ]);

  const handleCreateNewServer = () => {
    setIsCreatingServer(true);
    setSelectedServer(null);
    setSelectedChannel(null);
  };

  const handleServerSelect = (serverId: number) => {
    setSelectedServer(serverId);
    setSelectedChannel(null);
    setIsChannelListVisible(true);
    setIsCreatingServer(false);
  };

  const handleChannelSelect = (channelId: number) => {
    setSelectedChannel(channelId);
    setIsChannelListVisible(false);
  };

  const selectedServerName = servers.find(s => s.id === selectedServer)?.name;

  const specialViews: Record<string, React.JSX.Element> = {
    "Edit Account": <EditAccount />,
    "Friends": <FriendList />,
  };

  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;

    fetch("http://localhost:8080/servers", {
      credentials: "include",
      signal,
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to load servers");
        return res.json();
      })
      .then((data) => setServers(data))
      .catch((err) => {
        if (err.name === "AbortError") {
          console.log("Fetch aborted");
          return;
        }
        console.error("Error fetching servers:", err);
      });

    return () => {
      controller.abort();
    };
  }, []);

  const fetchServers = () => {
    fetch("http://localhost:8080/servers", {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => setServers(data))
      .catch((err) => console.error("Error reloading servers:", err));
  };

  return (
    <>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
      </Routes>

      <div className="flex h-screen w-screen bg-gray-900">
        <ServerList
          onServerSelect={handleServerSelect}
          onCreateNewServer={handleCreateNewServer}
        />

        {selectedServerName && specialViews[selectedServerName]}

        {isCreatingServer && <CreateServer onServerCreated={fetchServers} />}

        {selectedServer !== null &&
          !isCreatingServer &&
          !specialViews[selectedServerName!] && (
            <>
              {isChannelListVisible && (
                <ChannelList
                  serverId={selectedServer}
                  onChannelSelect={handleChannelSelect}
                />
              )}
              {selectedChannel && (
                <Chat
                  channelId={selectedChannel}
                  onToggleChannels={() => setIsChannelListVisible(true)}
                />
              )}
            </>
          )}
      </div>
    </>
  );
}

export default App;
