import { useState } from "react";

// Add a type definition for electronAPI on the Window interface
declare global {
  interface Window {
    electronAPI: {
      selectFolder: () => Promise<string>;
      listVideos: (folder: string) => Promise<string[]>;
      moveVideo: (args: {
        filePath: string;
        targetDir: string;
      }) => Promise<void>;
    };
  }
}

function App() {
  const [videos, setVideos] = useState<string[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);

  const loadFolder = async () => {
    const folder = await window.electronAPI.selectFolder();
    if (folder) {
      const files = await window.electronAPI.listVideos(folder);
      setVideos(files);
      setCurrentIndex(0);
    }
  };

  const swipeLeft = async () => {
    const filePath = videos[currentIndex];
    const targetDir = "C:/Users/Public/Videos/Rejected"; // exemple
    await (window as any).electronAPI.moveVideo({ filePath, targetDir });
    setCurrentIndex((i) => i + 1);
  };

  const swipeRight = () => {
    setCurrentIndex((i) => i + 1);
  };

  const currentVideo = videos[currentIndex];
  console.log(currentVideo);
  return (
    <div className="flex flex-col items-center p-6">
      <h1 className="text-2xl font-bold mb-4">
        🎬 Collab Video Slider
      </h1>

      <button
        onClick={loadFolder}
        className="px-4 py-2 bg-blue-500 text-white rounded"
      >
        Choisir un dossier
      </button>

      {currentVideo ? (
        <div className="mt-6 w-[800px]">
          <video
            src={`file:///${currentVideo.replace(/\\/g, "/")}`}
            controls
            className="w-full rounded-xl shadow-lg"
          />
          <div className="flex justify-between mt-4">
            <button
              onClick={swipeLeft}
              className="px-6 py-2 bg-red-500 text-white rounded-xl"
            >
              ⬅️ Supprimer
            </button>
            <button
              onClick={swipeRight}
              className="px-6 py-2 bg-green-500 text-white rounded-xl"
            >
              ➡️ Garder
            </button>
          </div>
        </div>
      ) : (
        <p className="mt-6">Aucune vidéo à traiter ✅</p>
      )}
    </div>
  );
}

export default App;
