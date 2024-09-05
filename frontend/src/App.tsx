import "./App.css";
import YouTube from "react-youtube";
import { IconThumbUp, IconThumbDown } from "@tabler/icons-react";
import { useEffect, useState } from "react";

interface YTVideo {
  id: number;
  title: string;
  description: string;
  video_tag: string;
}

function App() {
  const [vid, setVid] = useState<YTVideo>({
    id: -1,
    title: "rick roll",
    description: "rick roll",
    video_tag: "dQw4w9WgXcQ",
  });
  const [isTopBarShowing, setIsTopBarShowing] = useState(true);

  let topBarInterval: number;
  const hideTopBar = () => {
    setIsTopBarShowing(true);
    clearTimeout(topBarInterval);

    topBarInterval = setTimeout(() => {
      setIsTopBarShowing(false);
    }, 5000);
  };

  const nextVideo = async () => {
    let resp = await fetch("http://localhost:8080/next-vid");
    let data = await resp.json();
    setVid(data);
  };

  useEffect(() => {
    nextVideo();
    hideTopBar();
  }, []);

  useEffect(() => {
    console.log(vid);
  }, [vid]);

  const increaseVolume = (e: any) => {
    setTimeout(
      () => {
        e.target.setVolume(30);
        e.target.playVideo();
      },
      1000,
      e
    );
  };

  const opts = {
    playerVars: {
      // https://developers.google.com/youtube/player_parameters
      autoload: 1,
      mute: 1,
      origin: "http://localhost:5173",
    },
  };
  return (
    <div
      onMouseMoveCapture={() => {
        hideTopBar();
      }}
      className="black-bg"
    >
      <YouTube
        videoId={vid.video_tag}
        onEnd={nextVideo}
        onPlay={(e: any) => {
          increaseVolume(e);
        }}
        opts={opts}
        className="w-full h-full"
      />
      <div
        className={`fixed flex top-0 w-screen justify-center transition-transform ${
          isTopBarShowing ? "" : "-translate-y-20"
        }`}
      >
        <div className="bg-black/50 rounded-b-xl text-white h-fit w-fit flex gap-2 p-5">
          <h1 className="text-lg font-bold border-[1px] border-white/50 rounded-xl self-center p-2">
            https://www.youtube.com/{vid.video_tag}
          </h1>
          <button className="bg-green-700">
            <IconThumbUp size={28} className="text-white bg-green-700" />
          </button>
          <button className="bg-red-800">
            <IconThumbDown size={28} className="text-white" />
          </button>
        </div>
      </div>
    </div>
  );
}

export default App;
