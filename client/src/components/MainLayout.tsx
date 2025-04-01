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
    selectedView: "edit" | "friends" | null;
    isCreatingServer: boolean;
    isChannelListVisible: boolean;
    onServerSelect: (id: number | null) => void;
    onCreateNewServer: () => void;
    onChannelSelect: (id: number | null) => void;
    onServerCreated: () => void;
    onIsCreatingServer: (selectedView: "edit" | "friends" | null) => void;
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
    selectedView,
}: Props) => {
    const { user } = useUser();

    return (
        <div className="flex h-screen w-screen bg-gray-900">
            <Sidebar
                servers={servers}
                onServerSelect={(id) => {
                    onServerSelect(id);
                    // Don't change selectedView to null here anymore, since it was handled in App.tsx
                }}
                onCreateNewServer={onCreateNewServer}
                onNavigate={(view) => {
                    onIsCreatingServer(view);
                    onServerSelect(null); // This hides the server channels when navigating to "edit" or "friends"
                }}
            />

            {/* Show the EditAccount or FriendList component based on selectedView */}
            {selectedView === "edit" && <EditAccount />}
            {selectedView === "friends" && <FriendList />}

            {/* Show CreateServer when in creation mode */}
            {isCreatingServer && <CreateServer onServerCreated={onServerCreated} />}

            {/* Only show channels when a server is selected and we're not in edit/friends view */}
            {selectedServer !== null && !isCreatingServer && !selectedView && (
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
