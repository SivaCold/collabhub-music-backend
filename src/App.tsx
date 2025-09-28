import React, { useState } from "react";

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
      ensureDir: (dirPath: string) => Promise<boolean>;
      listDirs: (folder: string) => Promise<string[]>;
    };
  }
}

const getYear = () => new Date().getFullYear();

function App() {
  const [videos, setVideos] = useState<string[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [customDirs, setCustomDirs] = useState<string[]>([]);
  const [newButtonName, setNewButtonName] = useState("");
  const [currentFolder, setCurrentFolder] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<string>("");

  // Load videos and custom directories
  const loadFolder = async () => {
    const folder = await window.electronAPI.selectFolder();
    if (folder) {
      const files = await window.electronAPI.listVideos(folder);
      setVideos(files);
      setCurrentIndex(0);
      setCurrentFolder(folder);
      // Load custom directories (excluding Trash)
      const dirs = await window.electronAPI.listDirs(folder);
      setCustomDirs(dirs.filter((d) => d !== "Trash"));
      setActiveTab(""); // Reset active tab
    }
  };

  // Update custom directories for the current video
  const updateCustomDirs = async (filePath: string) => {
    const videoDir = filePath.substring(0, filePath.lastIndexOf("\\"));
    const dirs = await window.electronAPI.listDirs(videoDir);
    setCustomDirs(dirs.filter((d) => d !== "Trash"));
  };

  // Move video to Trash
  const swipeLeft = async () => {
    const filePath = videos[currentIndex];
    if (!filePath) return;
    const videoDir = filePath.substring(0, filePath.lastIndexOf("\\"));
    const trashDir = `${videoDir}\\Trash`;
    await window.electronAPI.ensureDir(trashDir);
    await window.electronAPI.moveVideo({ filePath, targetDir: trashDir });
    const nextIndex = currentIndex + 1;
    setCurrentIndex(nextIndex);
    if (videos[nextIndex]) updateCustomDirs(videos[nextIndex]);
  };

  // Keep video (move to next)
  const swipeRight = () => {
    const nextIndex = currentIndex + 1;
    setCurrentIndex(nextIndex);
    if (videos[nextIndex]) updateCustomDirs(videos[nextIndex]);
  };

  // Add a custom action (directory) with user-chosen name
  const addCustomAction = async () => {
    if (!newButtonName.trim()) return;
    const filePath = videos[currentIndex];
    if (!filePath) return;
    const videoDir = filePath.substring(0, filePath.lastIndexOf("\\"));
    const customDir = `${videoDir}\\${newButtonName}`;
    await window.electronAPI.ensureDir(customDir);
    // Update the list of directories
    const dirs = await window.electronAPI.listDirs(videoDir);
    setCustomDirs(dirs.filter((d) => d !== "Trash"));
    setNewButtonName("");
  };

  // Move video to a custom directory
  const moveToCustomDir = async (dirName: string) => {
    const filePath = videos[currentIndex];
    if (!filePath) return;
    const videoDir = filePath.substring(0, filePath.lastIndexOf("\\"));
    const targetDir = `${videoDir}\\${dirName}`;
    await window.electronAPI.moveVideo({ filePath, targetDir });
    const nextIndex = currentIndex + 1;
    setCurrentIndex(nextIndex);
    if (videos[nextIndex]) updateCustomDirs(videos[nextIndex]);
  };

  const currentVideo = videos[currentIndex];

  // Update custom directories when the current video changes
  React.useEffect(() => {
    if (currentVideo) updateCustomDirs(currentVideo);
  }, [currentVideo]);

  // All tabs: default actions + custom directories
  const allTabs = [
    { label: "Supprimer", key: "Trash" },
    { label: "Garder", key: "Keep" },
    ...customDirs.map((dir) => ({ label: dir, key: dir })),
  ];

  // Handle tab click (Bootstrap Nav)
  const handleTabClick = (tabKey: string) => {
    setActiveTab(tabKey);
    if (tabKey === "Trash") swipeLeft();
    else if (tabKey === "Keep") swipeRight();
    else moveToCustomDir(tabKey);
  };

  return (
    <div className="d-flex flex-column min-vh-100 bg-light">
      {/* Header with Bootstrap Nav Tabs */}
      <header className="bg-white border-bottom shadow-sm">
        <div className="container-fluid py-3 d-flex align-items-center">
          {/* Logo and title */}
          <img
            src="https://gitlab.com/uploads/-/system/project/avatar/278964/logo.png"
            alt="Logo"
            style={{ width: 40, height: 40, marginRight: 16 }}
          />
          <h1 className="h4 mb-0 fw-bold">Collab Video Slider App</h1>
          <button
            onClick={loadFolder}
            className="btn btn-primary ms-auto"
          >
            Choisir un dossier
          </button>
        </div>
        {/* Bootstrap Nav Tabs */}
        <ul className="nav nav-tabs px-3">
          {allTabs.map((tab) => (
            <li className="nav-item" key={tab.key}>
              <button
                className={`nav-link ${activeTab === tab.key ? "active" : ""}`}
                onClick={() => handleTabClick(tab.key)}
                disabled={!currentVideo}
                style={{ cursor: !currentVideo ? "not-allowed" : "pointer" }}
              >
                {tab.label}
              </button>
            </li>
          ))}
        </ul>
      </header>

      {/* Main content */}
      <main className="container flex-grow-1 d-flex flex-column align-items-center py-4">
        {/* Input for new custom button */}
        <form
          className="d-flex align-items-center gap-2 mb-4"
          onSubmit={e => { e.preventDefault(); addCustomAction(); }}
        >
          <input
            type="text"
            placeholder="Nom du bouton"
            value={newButtonName}
            onChange={(e) => setNewButtonName(e.target.value)}
            className="form-control"
            style={{ width: 180 }}
            disabled={!currentVideo}
          />
          <button
            type="submit"
            className="btn btn-warning"
            disabled={!currentVideo}
          >
            Ajouter un bouton
          </button>
        </form>

        {currentVideo ? (
          <div className="card shadow-lg mb-4" style={{ maxWidth: 700, width: "100%" }}>
            <div className="card-body d-flex flex-column align-items-center">
              {/* Video player */}
              <video
                src={`file:///${currentVideo.replace(/\\/g, "/")}`}
                controls
                className="w-100 rounded mb-3"
                style={{ background: "#222" }}
              />
              {/* Action buttons */}
              <div className="d-flex flex-wrap gap-2 justify-content-center">
                <button
                  onClick={swipeLeft}
                  className="btn btn-danger"
                >
                  ❌ Supprimer
                </button>
                <button
                  onClick={swipeRight}
                  className="btn btn-success"
                >
                  💚 Garder
                </button>
                {customDirs.map((dir) => (
                  <button
                    key={dir}
                    onClick={() => moveToCustomDir(dir)}
                    className="btn btn-warning"
                  >
                    📁 {dir}
                  </button>
                ))}
              </div>
            </div>
          </div>
        ) : (
          <p className="mt-5 fs-5 text-secondary">
            Aucune vidéo à traiter ✅
          </p>
        )}
      </main>

      {/* Footer with copyright */}
      <footer className="bg-white border-top py-3 text-center text-muted small mt-auto">
        © {getYear()} Collab Video Slider. All rights reserved.
      </footer>
    </div>
  );
}

export default App;