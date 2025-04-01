import { useState } from "react";
import { useUser } from "../hooks/useUser";
import { Sidebar } from "./Sidebar";
import { ChannelList } from "./ChannelList";
import { Chat } from "./Chat";
import EditAccount from "./EditAccount";
import { FriendList } from "./FriendList";
import CreateServer from "./CreateServer";

export interface Server {
    id: number;
    name: string;
    description: string;
    image: string;
    userId: number;
    createdAt: string;
    updatedAt: string;
}

interface Props {
    servers: Server[];
    selectedServer: number | null;
    selectedChannel: number | null;
    isCreatingServer: boolean;
    isChannelListVisible: boolean;
    onServerSelect: (id: number | null) => void;
    onCreateNewServer: () => void;
    onChannelSelect: (id: number | null) => void;
    onServerCreated: () => void;
}

export const MainLayout = ({
    servers,
    selectedServer,
    selectedChannel,
    isCreatingServer,
    isChannelListVisible,
    onServerSelect,
    onCreateNewServer,
    onChannelSelect,
    onServerCreated,
}: Props) => {
    const [selectedView, setSelectedView] = useState<"edit" | "friends" | null>(null);
    const { user } = useUser();
    const currentUserId = user?.id ?? 0;

    const currentServer = servers.find((s) => s.id === selectedServer);
    const serverOwnerId = currentServer?.userId ?? 0;

    return (
        <div className="flex h-screen w-screen bg-gray-900">
            <Sidebar
                servers={servers}
                onServerSelect={(id) => {
                    onServerSelect(id);
                    setSelectedView(null);
                }}
                onCreateNewServer={onCreateNewServer}
                onNavigate={(view) => {
                    setSelectedView(view);
                    onServerSelect(null);
                }}
            />

            {/* Special views */}
            {selectedView === "edit" && <EditAccount />}
            {selectedView === "friends" && <FriendList />}

            {/* Create server view */}
            {isCreatingServer && <CreateServer onServerCreated={onServerCreated} />}

            {/* Normal server view */}
            {selectedServer !== null && !isCreatingServer && !selectedView && (
                <>
                    {isChannelListVisible && (
                        <ChannelList
                            serverId={selectedServer}
                            serverOwnerId={serverOwnerId}
                      
                            onChannelSelect={onChannelSelect}
                        />
                    )}
                    {selectedChannel && (
                        <Chat
                            channelId={selectedChannel}
                            onToggleChannels={() => onChannelSelect(null)}
                        />
                    )}
                </>
            )}
        </div>
    );
};