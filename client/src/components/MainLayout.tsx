import { ChannelList } from "./ChannelList";
import { Chat } from "./Chat";
import EditAccount from "./EditAccount";
import { FriendList } from "./FriendList";
import CreateServer from "./CreateServer";
import { Sidebar } from "./Sidebar";
import { useState } from "react";
import { useUser } from "../hooks/useUser";

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
    const { user } = useUser()

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

            {selectedView === "edit" && <EditAccount />}
            {selectedView === "friends" && <FriendList />}

            {isCreatingServer && <CreateServer onServerCreated={onServerCreated} />}

            {selectedServer !== null &&
                !isCreatingServer &&
                !selectedView && (
                    <>
                        {isChannelListVisible && (
                            <ChannelList
                                serverId={selectedServer}
                                serverOwnerId={servers.find((s) => s.id === selectedServer)?.userId ?? 0}
                                currentUserId={user?.id ?? 0}
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