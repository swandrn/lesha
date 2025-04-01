import React from "react";

const baseProfileLink = "https://media.istockphoto.com/id/1495088043/fr/vectoriel/ic%C3%B4ne-de-profil-utilisateur-avatar-ou-ic%C3%B4ne-de-personne-photo-de-profil-symbole-portrait.jpg?s=612x612&w=0&k=20&c=moNRZjYtVpH-I0mAe-ZfjVkuwgCOqH-BRXFLhQkZoP8="
// Exemple de données d'amis, avec des statuts et des avatars
const friends = [
  { id: 1, name: "John", status: "online", avatar: baseProfileLink },
  { id: 2, name: "Jane", status: "offline", avatar: baseProfileLink },
  { id: 3, name: "Alex", status: "online", avatar: baseProfileLink },
  { id: 4, name: "Emily", status: "away", avatar: baseProfileLink },
];

export function FriendList() {
  const getStatusColor = (status: string) => {
    switch (status) {
      case "online":
        return "bg-green-500"; // En ligne
      case "offline":
        return "bg-gray-500"; // Hors ligne
      case "away":
        return "bg-yellow-500"; // En pause
      default:
        return "bg-gray-500";
    }
  };

  return (
    <div className="bg-gray-800 p-6 rounded-lg shadow-lg max-w-xs w-full">
      <h2 className="text-white text-xl font-semibold mb-4">Liste d'amis</h2>

      {/* En ligne */}
      <div>
        <h3 className="text-white text-lg font-semibold mb-2">En ligne</h3>
        <ul>
          {friends
            .filter(friend => friend.status === "online")
            .map(friend => (
              <li key={friend.id} className="flex items-center mb-4">
                <img src={friend.avatar} alt={friend.name} className="w-12 h-12 rounded-full mr-3" />
                <div className="flex-1">
                  <span className="text-white font-medium">{friend.name}</span>
                  <div className={`w-3 h-3 rounded-full ${getStatusColor(friend.status)} mt-1`} />
                </div>
              </li>
            ))}
        </ul>
      </div>

      {/* Hors ligne */}
      <div className="mt-4">
        <h3 className="text-white text-lg font-semibold mb-2">Hors ligne</h3>
        <ul>
          {friends
            .filter(friend => friend.status === "offline")
            .map(friend => (
              <li key={friend.id} className="flex items-center mb-4">
                <img src={friend.avatar} alt={friend.name} className="w-12 h-12 rounded-full mr-3" />
                <div className="flex-1">
                  <span className="text-white font-medium">{friend.name}</span>
                  <div className={`w-3 h-3 rounded-full ${getStatusColor(friend.status)} mt-1`} />
                </div>
              </li>
            ))}
        </ul>
      </div>

      {/* En pause */}
      <div className="mt-4">
        <h3 className="text-white text-lg font-semibold mb-2">En pause</h3>
        <ul>
          {friends
            .filter(friend => friend.status === "away")
            .map(friend => (
              <li key={friend.id} className="flex items-center mb-4">
                <img src={friend.avatar} alt={friend.name} className="w-12 h-12 rounded-full mr-3" />
                <div className="flex-1">
                  <span className="text-white font-medium">{friend.name}</span>
                  <div className={`w-3 h-3 rounded-full ${getStatusColor(friend.status)} mt-1`} />
                </div>
              </li>
            ))}
        </ul>
      </div>
    </div>
  );
}
