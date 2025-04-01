import { useEffect, useState } from "react";
import Login from "./components/Login";
import Register from "./components/Register";
import { Route, Routes } from "react-router-dom";
import { ProtectedRoute } from "./components/middlewares/ProtectedRoute";
import { MainLayout, Server } from "./components/MainLayout";

function App() {
  const [selectedServer, setSelectedServer] = useState<number | null>(null);
  const [selectedChannel, setSelectedChannel] = useState<number | null>(null);
  const [selectedView, setSelectedView] = useState<"edit" | "friends" | null>(null);
  const [isChannelListVisible, setIsChannelListVisible] = useState(true);
  const [isCreatingServer, setIsCreatingServer] = useState(false);
  const [servers, setServers] = useState<Server[]>([]);
  const [isLoggedIn, setIsLoggedIn] = useState(true);

  // Handle creation of new server
  const handleCreateNewServer = () => {
    setIsCreatingServer(true);
    setSelectedServer(null);
    setSelectedChannel(null);
    setSelectedView(null);  // Reset selectedView when creating a new server
  };

  // Handle server selection
  const handleServerSelect = (serverId: number | null) => {
    setSelectedServer(serverId);
    setSelectedChannel(null);
    setIsChannelListVisible(true);
    setIsCreatingServer(false);
  };

  // Handle channel selection
  const handleChannelSelect = (channelId: number | null) => {
    setSelectedChannel(channelId);
    setIsChannelListVisible(false);
  };

  // Handle changing views (edit or friends)
  const handleIsCreatingServerSelect = (view: "edit" | "friends" | null) => {
    setSelectedView(view);  // Update selectedView when a view is selected
    setSelectedChannel(null);  // Clear selectedChannel when changing views
    setIsChannelListVisible(false);  // Hide channel list when navigating to a view
  };

  useEffect(() => {
    if (!isLoggedIn) return;

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
      .then((data) =>
        setServers(
          data.map((server: any) => ({
            id: server.ID,
            name: server.Name,
            image: server.Image,
            description: server.Description,
            userId: server.UserID,
            createdAt: server.CreatedAt,
            updatedAt: server.UpdatedAt,
          }))
        )
      )
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
  }, [isLoggedIn]);

  const fetchServers = () => {
    fetch("http://localhost:8080/servers", {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
        setServers(data);
        setIsLoggedIn(true);
      })
      .catch((err) => console.error("Error reloading servers:", err));
  };

  return (
    <Routes>
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <MainLayout
              servers={servers}
              selectedServer={selectedServer}
              selectedChannel={selectedChannel}
              selectedView={selectedView}  // Pass selectedView as a prop
              isCreatingServer={isCreatingServer}
              isChannelListVisible={isChannelListVisible}
              onServerSelect={handleServerSelect}
              onCreateNewServer={handleCreateNewServer}
              onChannelSelect={handleChannelSelect}
              onServerCreated={fetchServers}
              onIsCreatingServer={handleIsCreatingServerSelect}  // Pass the handler to manage view
            />
          </ProtectedRoute>
        }
      />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
    </Routes>
  );
}

export default App;
