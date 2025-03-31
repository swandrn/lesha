import { useState } from "react";

function CreateServer() {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [image, setImage] = useState<File | null>(null);
  const [userId, setUserId] = useState<number | string>(""); // ID de l'utilisateur
  const [error, setError] = useState<string>("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!name || !description || !image || !userId) {
      setError("Tous les champs sont obligatoires.");
      return;
    }

    // Logique de soumission du formulaire
    const formData = new FormData();
    formData.append("name", name);
    formData.append("description", description);
    formData.append("image", image);
    formData.append("user_id", userId.toString());

    // Vous pouvez remplacer cette ligne par une requête pour envoyer ces données à votre API
    console.log("FormData pour création du serveur :", formData);

    // Réinitialiser le formulaire après soumission
    setName("");
    setDescription("");
    setImage(null);
    setUserId("");
    setError(""); // Réinitialiser les erreurs
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-900 text-white">
      <form
        onSubmit={handleSubmit}
        className="w-full max-w-lg p-6 bg-black border-4 border-blue-400 rounded-lg shadow-md"
      >
        <div className="text-center">
          <h2 className="text-xl font-bold">Créer un serveur</h2>
          <hr className="w-1/2 mx-auto mt-2 border-blue-400" />
        </div>

        {error && (
          <div className="mt-4 text-center text-red-500 text-sm">{error}</div>
        )}

        <div className="mt-4">
          <label
            htmlFor="name"
            className="block text-sm font-bold text-gray-300"
          >
            Nom du serveur
          </label>
          <input
            id="name"
            type="text"
            placeholder="Nom du serveur"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            className="w-full px-3 py-2 mt-1 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
          />
        </div>

        <div className="mt-4">
          <label
            htmlFor="description"
            className="block text-sm font-bold text-gray-300"
          >
            Description
          </label>
          <textarea
            id="description"
            placeholder="Description du serveur"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            required
            className="w-full px-3 py-2 mt-1 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
          />
        </div>

        <div className="mt-4">
          <label
            htmlFor="image"
            className="block text-sm font-bold text-gray-300"
          >
            Image du serveur
          </label>
          <input
            id="image"
            type="file"
            onChange={(e) => {
              const file = e.target.files ? e.target.files[0] : null;
              setImage(file);
            }}
            required
            className="w-full mt-1 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
          />
        </div>

        <div className="mt-4">
          <label
            htmlFor="user_id"
            className="block text-sm font-bold text-gray-300"
          >
            ID de l'utilisateur
          </label>
          <input
            id="user_id"
            type="number"
            placeholder="ID de l'utilisateur"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
            required
            className="w-full px-3 py-2 mt-1 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
          />
        </div>

        <div className="mt-6 flex justify-center">
          <input
            type="submit"
            value="Créer le serveur"
            className="px-6 py-2 text-white bg-blue-700 rounded cursor-pointer hover:bg-blue-400 transition"
          />
        </div>
      </form>
    </div>
  );
}

export default CreateServer;
