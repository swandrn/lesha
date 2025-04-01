import { useState } from "react";

function CreateServer({ onServerCreated }: { onServerCreated: () => void }) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [image, setImage] = useState<File | null>(null);
  const [error, setError] = useState<string>("");
  const [imagePreview, setImagePreview] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
  
    // Validation
    if (!name || !description || !image) {
      setError("Tous les champs sont obligatoires.");
      return;
    }
  
    try {
      const formData = new FormData();
      formData.append("name", name);
      formData.append("description", description);
      formData.append("image", image);
  
      const response = await fetch("http://localhost:8080/servers", {
        method: "POST",
        body: formData,
        credentials: "include", // ✅ Include JWT cookie
      });
  
      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || "Erreur lors de la création du serveur");
      }
  
      const data = await response.json();
      console.log("Serveur créé :", data);
  
      // ✅ Call parent callback to refresh list
      if (onServerCreated) onServerCreated();
  
      // Optional: notify the user
      alert("Serveur créé avec succès !");
    } catch (err) {
      console.error(err);
      setError("Erreur lors de la création du serveur.");
      return;
    }
  
    // ✅ Reset form
    setName("");
    setDescription("");
    setImage(null);
    setImagePreview(null);
    setError("");
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files ? e.target.files[0] : null;
    if (file) {
      setImage(file);
      // Créer un aperçu de l'image
      const previewUrl = URL.createObjectURL(file);
      setImagePreview(previewUrl);
    }
  };

  return (
    <div className="w-full h-full flex items-center justify-center">
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
              onChange={handleImageChange}
              required
              className="w-full mt-1 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
            />
          </div>

          {/* Affichage de l'aperçu de l'image si une image est sélectionnée */}
          {imagePreview && (
            <div className="mt-4 flex justify-center">
              <img
                src={imagePreview}
                alt="Aperçu"
                className="w-48 h-48 object-cover rounded-md"
              />
            </div>
          )}

          <div className="mt-6 flex justify-center">
            <input
              type="submit"
              value="Créer le serveur"
              className="px-6 py-2 text-white bg-blue-700 rounded cursor-pointer hover:bg-blue-400 transition"
            />
          </div>
        </form>
      </div>
    </div>
  );
}

export default CreateServer;
